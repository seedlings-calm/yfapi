package guild

import (
	"encoding/json"
	"errors"
	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreJwtToken"
	"yfapi/core/coreRedis"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/logic"
	"yfapi/internal/model"
	typedef_enum "yfapi/typedef/enum"
	common_data "yfapi/typedef/redisKey"
	request_login "yfapi/typedef/request/guild"
	response_login "yfapi/typedef/response/guild"
)

type GuildLogin struct {
}

// 手机验证码登录
func (g *GuildLogin) LoginByCode(req *request_login.LoginMobileReq, context *gin.Context) (res response_login.LoginMobileCodeRes) {
	//检测验证码
	sms := &logic.Sms{
		Mobile:     req.Mobile,
		Code:       req.Code,
		RegionCode: req.RegionCode,
		Type:       typedef_enum.SmsCodeGuildAdminLogin,
	}
	//data := helper.GetHeaderData(context)
	//校验手机验证码
	err := sms.CheckSms(context)
	if err != nil {
		panic(i18n_err.ErrorCodeCaptchaInvalid)
	}

	userInfo, guildInfo := g.SearchByMobile(context, &request_login.SendMobileCodeReq{
		Mobile:     req.Mobile,
		RegionCode: req.RegionCode,
	})

	//生成token
	token, err := coreJwtToken.GuildEncode(coreJwtToken.GuildClaims{
		UserId:         guildInfo.UserID,
		Mobile:         req.Mobile,
		GuildId:        guildInfo.ID,
		StandardClaims: gojwt.StandardClaims{},
	}, []byte(coreConfig.GetHotConf().JwtSecret), time.Now().Add(time.Hour*24*7).Unix())
	if err != nil {
		panic(i18n_err.ErrorCodeSystemBusy)
	}
	//将token放入登录缓存
	//需要根据客户端类型进行判定
	coreRedis.GetUserRedis().Set(context, common_data.UserLoginInfo("guildPc", guildInfo.UserID), token, 7*24*time.Hour)

	// 记录本次登录时间
	key := common_data.GuildAdminLoginTime(guildInfo.ID)
	var recordList []string
	now := time.Now().Format(time.DateTime)
	recordStr := coreRedis.GetUserRedis().Get(context, key).Val()
	if len(recordStr) == 0 {
		recordList = append(recordList, now)
	} else {
		_ = json.Unmarshal([]byte(recordStr), &recordList)
		recordList = append(recordList, now)
		if len(recordList) > 2 {
			recordList = recordList[1:]
		}
	}
	recordByte, _ := json.Marshal(recordList)
	coreRedis.GetUserRedis().Set(context, key, string(recordByte), 0)

	return response_login.LoginMobileCodeRes{
		Token:    token,
		UserID:   guildInfo.UserID,
		Guild:    guildInfo.Name,
		UserName: userInfo.Nickname,
		GuildID:  guildInfo.ID,
	}
}

// 获取手机验证码
func (g *GuildLogin) SearchByMobile(c *gin.Context, req *request_login.SendMobileCodeReq) (userInfo *model.User, guildInfo *model.Guild) {
	// 根据手机区号和手机号查询用户信息
	var err error
	userInfo, err = new(dao.UserDao).FindOne(&model.User{
		RegionCode: req.RegionCode,
		Mobile:     req.Mobile,
	})
	// 用户信息不存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeCheckMobile,
			Msg:  nil,
		})
	}
	// 数据读取失败
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 公会信息查询
	guildInfo, err = new(dao.GuildDao).FindById(&model.Guild{
		UserID: userInfo.Id,
		Status: typedef_enum.GuildStatusNormal,
	})
	// 不是会长
	if errors.Is(err, gorm.ErrRecordNotFound) {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeCheckMobile,
			Msg:  nil,
		})
	}
	// 数据查询失败
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	return
}

// GuildInfo
//
//	@Description: 获取公会信息
//	@receiver g
//	@param c *gin.Context -
//	@return resp -
func (g *GuildLogin) GuildInfo(c *gin.Context) (resp response_login.GuildInfo) {
	guildId := helper.GetGuildId(c)
	var err error
	// 公会信息查询
	resp, err = new(dao.GuildDao).GetGuildInfo(guildId)
	// 公会不存在
	if errors.Is(err, gorm.ErrRecordNotFound) {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeNotExistGuild,
			Msg:  nil,
		})
	}
	// 数据查询失败
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	resp.UserAvatar = helper.FormatImgUrl(resp.UserAvatar)
	resp.LogoImg = helper.FormatImgUrl(resp.LogoImg)

	// 上次登录时间
	lastLoginTime := ""
	key := common_data.GuildAdminLoginTime(guildId)
	var recordList []string
	recordStr := coreRedis.GetUserRedis().Get(c, key).Val()
	if len(recordStr) > 0 {
		_ = json.Unmarshal([]byte(recordStr), &recordList)
		if len(recordList) > 1 {
			lastLoginTime = recordList[0]
		}
	}
	resp.LastLoginTime = lastLoginTime
	return
}
