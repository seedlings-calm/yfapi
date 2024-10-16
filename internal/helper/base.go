package helper

import (
	"context"
	"crypto/sha1"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
	"yfapi/core/coreConfig"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/typedef"
	"yfapi/typedef/enum"

	"github.com/gin-gonic/gin"
)

// GetHeaderData 获取头部参数
func GetHeaderData(c *gin.Context) (data typedef.HeaderData) {
	getData, _ := c.Get("headerData")
	switch v := getData.(type) {
	case typedef.HeaderData:
		data = v
	}
	return
}

// EncodeUserPass 加密用户密码
func EncodeUserPass(password string) string {

	return fmt.Sprintf("%x", sha1.Sum([]byte(password)))
}

func GetUserId(c *gin.Context) string {
	return c.GetString("userId")
}

func GetRoomId(c *gin.Context) string {
	return c.GetHeader("roomId")
}

func GetGuildId(c *gin.Context) string {
	return c.GetString("guildId")
}

func GetClientType(c *gin.Context) string {
	return c.GetString("clientType")
}

func GetMobile(c *gin.Context) string {
	return c.GetString("mobile")
}

// FormatImgUrl 获取图片全链接
func FormatImgUrl(url string) string {
	if len(url) == 0 {
		return ""
	}
	config := coreConfig.GetHotConf()
	return config.ImagePrefix + url
}

// 把斜杠转义符改成斜杠
func RemoveEscapedSlashes(s string) string {
	// 使用正则表达式匹配并替换掉转义斜杠
	re := regexp.MustCompile(`\\/`)
	return re.ReplaceAllString(s, "/")
}

// 移除图片地址头
func RemovePrefixImgUrl(url string) string {
	if len(url) == 0 {
		return ""
	}
	return strings.TrimPrefix(url, coreConfig.GetHotConf().ImagePrefix)
}

// GeneUserNo 生成userNo
func GeneUserNo() string {
	dicUserDao := new(dao.DicUserNoDao)
	data := dicUserDao.GetRandData()
	if data.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	data.Status = 1
	dicUserDao.Update(data)
	return data.UserNo
}

// 随机生成用户昵称
func GeneUserNickname() string {
	dicUserDao := new(dao.DicUserNoDao)
	data := dicUserDao.GetRandUserNickName()
	if data.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	data.Status = 1
	err := dicUserDao.UpdateNickNameStatus(data)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	return data.UserNickName
}

// GetUserDefaultAvatar
//
//	@Description: 根据性别获取用户默认头像
//	@param sex int -
//	@return string -
func GetUserDefaultAvatar(sex int) string {
	data, _ := new(dao.UserDefaultAvatarDao).GetUserDefaultAvatar(sex)
	if data.Id > 0 {
		return data.Avatar
	}
	return ""
}

func CheckPasswordLever(password string) error {
	if len(password) < 8 || len(password) > 20 {
		return fmt.Errorf("password need 8-20")
	}
	hasDigit := false
	hasLetter := false
	hasSpecialChar := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
		} else if unicode.IsLetter(char) {
			hasLetter = true
		} else {
			hasSpecialChar = true
		}
	}
	if hasSpecialChar {
		return fmt.Errorf("password has special char")
	}
	if hasDigit && hasLetter {
		return nil
	}
	return fmt.Errorf("password is illegal")
}

// 密码强度必须为字⺟⼤⼩写+数字
func CheckPasswordLever2(ps string) error {
	if len(ps) < 8 || len(ps) > 20 {
		return fmt.Errorf("password need 8-20")
	}
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	//A_Z := `[A-Z]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return fmt.Errorf("password need num :%v", err)
	}
	if b, err := regexp.MatchString(a_z, ps); !b || err != nil {
		return fmt.Errorf("password need a-z :%v", err)
	}
	//if b, err := regexp.MatchString(A_Z, ps); !b || err != nil {
	//	return fmt.Errorf("password need A_Z :%v", err)
	//}
	return nil
}

// 频率限制
func ReqRateLimit(key string, t time.Duration) bool {
	ok, err := coreRedis.GetUserRedis().SetNX(context.Background(), key, 1, t).Result()
	if err != nil || !ok {
		return false
	}
	return true
}

// 脱敏手机号
func PrivateMobile(mobile string) string {
	masked := ""
	if len(mobile) > 5 {
		masked = mobile[:3] + strings.Repeat("*", len(mobile)-5) + mobile[len(mobile)-2:]
	}
	return masked
}

// 脱敏名字
func PrivateRealName(name string) string {
	masked := ""
	runeName := []rune(name)
	switch {
	case len(runeName) > 3:
		masked = string(runeName[0]) + strings.Repeat("*", len(runeName)-2) + string(runeName[len(runeName)-1:])
	case len(runeName) == 3:
		masked = string(runeName[0]) + strings.Repeat("*", 1) + string(runeName[len(runeName)-1:])
	case len(runeName) == 2:
		masked = string(runeName[0]) + strings.Repeat("*", 1)
	default:
		masked = name
	}
	return masked
}

// 脱敏身份证号
func PrivateIdNo(idNo string) string {
	masked := ""
	switch {
	case len(idNo) > 14:
		masked = idNo[0:4] + strings.Repeat("*", len(idNo)-6) + idNo[len(idNo)-2:]
	case len(idNo) > 6 && len(idNo) <= 14:
		masked = idNo[0:3] + strings.Repeat("*", len(idNo)-5) + idNo[len(idNo)-2:]
	case len(idNo) <= 6 && len(idNo) > 3:
		masked = strings.Repeat("*", len(idNo)-2) + idNo[len(idNo)-2:]
	default:
		masked = strings.Repeat("*", len(idNo))
	}
	return masked
}

// 是否官方账户
func OfficialAccount(userId string) bool {
	switch userId {
	case enum.OperationUserId:
		return true
	case enum.SystematicUserId:
		return true
	case enum.OfficialUserId:
		return true
	case enum.InteractiveUserId:
		return true
	default:
		return false
	}
}

// 根据用户id查询当前手机关联的所有用户id
func GetUserIdsByUserId(userId string) (res []string) {
	userDao := new(dao.UserDao)
	data, _ := userDao.FindOne(&model.User{Id: userId})
	if len(data.Id) > 0 {
		res = userDao.GetUserIdsByMobile(data.RegionCode, data.Mobile)
	}
	return
}
