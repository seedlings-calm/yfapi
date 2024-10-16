package handle

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/typedef"
	common_data "yfapi/typedef/redisKey"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindJson 绑定 JSON 数据并验证
func BindJson[T any](c *gin.Context, req *T) {
	// 绑定 JSON 数据到结构体
	if err := c.BindJSON(req); err != nil {
		coreLog.Error("BindJson err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	validate := validator.New()
	// 验证结构体
	if err := validate.Struct(req); err != nil {
		coreLog.Error("BindJson validate err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
}

// BindBody body绑定参数，返回错误码
func BindBody(c *gin.Context, req any) {
	if err := c.ShouldBind(req); err != nil {
		coreLog.Error("BindBody err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		coreLog.Error("BindBody validate err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
}

// BindQuery query绑定参数，返回错误码
func BindQuery(c *gin.Context, req any) {
	if err := c.ShouldBind(req); err != nil {
		coreLog.Error("BindQuery err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		coreLog.Error("BindQuery validate err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
}

func GetTokenData(c *gin.Context) (data typedef.TokenData) {
	userId := c.GetString("userId")
	clientType := c.GetString("clientType")
	mobile := c.GetString("mobile")
	data.UserId = userId
	data.ClientType = clientType
	data.Mobile = mobile
	return
}
func RepeatSubmitPost(c *gin.Context) {
	body := make(map[string]interface{})
	path := c.FullPath()
	_ = c.ShouldBind(&body)
	bodyByte, _ := json.Marshal(body)
	key := common_data.RepeatSubmit(fmt.Sprintf("%x", md5.Sum([]byte(path+string(bodyByte)+c.GetString("Authorization")))))
	isRepeat := coreRedis.GetUserRedis().SetNX(c, key, 1, 2*time.Second).Val()
	if !isRepeat {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRepeatSubmit,
			Msg:  nil,
		})
	}
}

func GetUserId(c *gin.Context) string {
	return c.GetString("userId")
}

func GetClientType(c *gin.Context) string {
	return c.GetString("clientType")
}

func GetMobile(c *gin.Context) string {
	return c.GetString("mobile")
}

func GetHeaderData(c *gin.Context) (data typedef.HeaderData) {
	getData, _ := c.Get("headerData")
	switch v := getData.(type) {
	case typedef.HeaderData:
		data = v
	}
	return
}
