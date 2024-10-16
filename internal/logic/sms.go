package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/typedef/redisKey"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var (
	keyInfo = map[string]map[string]string{
		//国内
		"local": {
			"templ":  "BU2B52",
			"appId":  "106024",
			"appKey": "394ed18f023ddaab1e5318156b875c43",
			"url":    "https://api-v4.mysubmail.com/sms/xsend",
		},
		//国外 TODO:
		"foreign": {
			"templ":  "VLE0V2",
			"appId":  "63888",
			"appKey": "ed18e58afb3eafb4ff921d1b10a0d32f",
			"url":    "https://api-v4.mysubmail.com/internationalsms/xsend",
		},
	}
)

type Sms struct {
	Mobile     string //手机号
	Type       int    //验证码类型
	Code       string //验证码
	RegionCode string //区域码
	SignType   string //签名类型，md5，sha1，默认md5
}

// 短信发送req
type ReqSms struct {
	AppId       string `json:"appid"`   //必选
	To          string `json:"to"`      //必选
	Project     string `json:"project"` //必选
	Vars        string `json:"vars"`    //自定义短信内容时必选， 短信模板里自定义变量的覆盖值，格式"{"key":"value","key1":"value1"}"
	Tag         string `json:"tag"`
	Timestamp   string `json:"timestamp"`
	SignType    string `json:"sign_type"`
	SignVersion string `json:"sign_version"`
	Signature   string `json:"signature"` //必选
}

// 发送短信
func (s *Sms) SendSms(c *gin.Context) error {
	redisCli := coreRedis.NewRedis().GetUserRedis()

	s.Code = getRandCode()
	key := redisKey.UserSmsCode(s.RegionCode, s.Mobile, s.Type)

	var isDev bool
	if coreConfig.GetHotConf().ENV != "pro" {
		isDev = true
		s.Code = "8888"
	}

	addMsg := func() (err error) {
		if !isDev {
			// redis 设置成功，表示之前没有发送过短信
			err = s.SendMsg()
			if err == nil {
				//写入数据库
				smsDao := new(dao.SmsDao)
				sms := &model.Sms{
					RegionCode: s.RegionCode,
					Mobile:     s.Mobile,
					Types:      s.Type,
					Code:       s.Code,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				}
				err = smsDao.Create(sms)
				if err == nil {
					redisCli.Set(c, key, s.Code, 5*time.Minute)
					return nil
				}
			}
		}
		return err
	}
	ok, _ := redisCli.SetNX(c, key, s.Code, 5*time.Minute).Result()
	if ok {
		return addMsg()
	} else { //如果前端1分钟可以选择发送一次，需要覆盖之前的redis，重新写入一条新的mysql记录
		//TODO:
		expr, _ := redisCli.TTL(c, key).Result()
		coreLog.Info("TTL", expr.Seconds())
		if expr <= time.Minute*4 { // 五分钟有效期，如果过去1分钟，就可以重复获取验证码
			return addMsg()
		}
	}
	panic(error2.I18nError{
		Code: error2.ErrorCodeCaptchaExpre,
		Msg:  nil,
	})
}

// 校验 验证码
func (s *Sms) CheckSms(c context.Context) (err error) {
	redisCli := coreRedis.NewRedis().GetUserRedis()
	key := redisKey.UserSmsCode(s.RegionCode, s.Mobile, s.Type)
	coreLog.Info("redis-key", key)
	var value string
	value, err = redisCli.Get(c, key).Result()

	if err == redis.Nil || err != nil {
		return errors.New("验证码不存在，或者还没有获取验证码")
	}
	if value != s.Code {
		return errors.New("验证码错误")
	}
	if coreConfig.GetHotConf().ENV == "pro" {
		smsDao := new(dao.SmsDao)
		sms := model.Sms{
			RegionCode: s.RegionCode,
			Mobile:     s.Mobile,
			Types:      s.Type,
			Code:       value,
		}
		err = smsDao.Update(sms)
		if err != nil {
			return err
		}
	}
	redisCli.Del(c, key)
	return
}

// Send 发送验证码
func (s *Sms) SendMsg() error {
	if len(s.Mobile) == 0 || len(s.RegionCode) == 0 {
		panic("手机号或者区域有误")
	}
	reqMsg := new(ReqSms)
	if len(s.SignType) == 0 {
		s.SignType = "md5"
	}
	//reqMsg.SignType = s.SignType
	reqMsg.SignVersion = "2"
	//reqMsg.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)

	//TODO: 根据短信类型，修改vars
	reqMsg.Vars = fmt.Sprintf(`{"code":%s}`, s.Code)

	if s.RegionCode != "+86" {
		reqMsg.Project = keyInfo["foreign"]["templ"]
		reqMsg.AppId = keyInfo["foreign"]["appId"]
		reqMsg.To = s.RegionCode + s.Mobile
		return reqMsg.sendForeign()
	}
	reqMsg.Project = keyInfo["local"]["templ"]
	reqMsg.AppId = keyInfo["local"]["appId"]
	reqMsg.To = s.Mobile
	return reqMsg.send()
}

// 获取验证码
func getRandCode() string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 4; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 国内
func (rs *ReqSms) send() (err error) {
	var params = make(map[string]string)
	params, err = rs.generateSubmailSignature(keyInfo["local"]["appKey"])
	if err != nil {
		return
	}
	params["vars"] = rs.Vars
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, value)
	}
	var res string
	res, err = easy.PostForm(keyInfo["local"]["url"], formData)
	if err != nil {
		return
	}
	submailRes := new(SubmailResponse)
	json.Unmarshal([]byte(res), &submailRes)
	log.Printf("%+v", submailRes)
	if submailRes.Status == "success" {
		return nil
	}
	err = errors.New("国内短信发送失败" + submailRes.Msg)
	return
}

// SubmailResponse 定义 Submail 响应的结构体
type SubmailResponse struct {
	Status   string `json:"status"`
	SendID   string `json:"send_id,omitempty"`
	Fee      int    `json:"fee,omitempty"`
	SMSCount int    `json:"sms_credits,omitempty"`
	Code     int    `json:"code,omitempty"`
	Msg      string `json:"msg,omitempty"`
}

// 用Submail发送短信(国际)
func (rs *ReqSms) sendForeign() (err error) {
	var params = map[string]string{}
	params, err = rs.generateSubmailSignature(keyInfo["foreign"]["appKey"])
	if err != nil {
		return
	}
	params["vars"] = rs.Vars
	formData := url.Values{}
	for key, value := range params {
		formData.Set(key, value)
	}

	var res string
	res, err = easy.PostForm(keyInfo["foreign"]["url"], formData)
	if err != nil {
		return
	}
	submailRes := new(SubmailResponse)
	json.Unmarshal([]byte(res), &submailRes)
	if submailRes.Status == "success" {
		return nil
	}
	err = errors.New("国外短信发送失败")
	return
}

// GenerateSubmailSignature 生成 Submail 平台的数字签名
func (rs *ReqSms) generateSubmailSignature(appKey string) (map[string]string, error) {
	params := make(map[string]string)
	params["appid"] = rs.AppId
	params["to"] = rs.To
	params["project"] = rs.Project
	//params["sign_type"] = rs.SignType
	params["sign_version"] = rs.SignVersion
	//params["timestamp"] = rs.Timestamp

	// 按 key 升序排列
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 创建签名字符串
	var signStrings []string
	for _, k := range keys {
		signStrings = append(signStrings, fmt.Sprintf("%s=%s", k, url.QueryEscape(params[k])))
	}
	signStr := strings.Join(signStrings, "&")

	// 拼接 APPID 和 APPKEY
	signStr = fmt.Sprintf("%s%s%s%s%s", rs.AppId, appKey, signStr, rs.AppId, appKey)
	// 根据算法生成签名
	//var signature string
	//switch rs.SignType {
	//case "md5":
	//	hash := md5.Sum([]byte(signStr))
	//	signature = hex.EncodeToString(hash[:])
	//case "sha1":
	//	hash := sha1.Sum([]byte(signStr))
	//	signature = hex.EncodeToString(hash[:])
	//default:
	//	return nil, fmt.Errorf("unsupported algorithm: %s", rs.SignType)
	//}
	params["signature"] = appKey
	return params, nil
}
