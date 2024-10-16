package middle

import (
	"net/http"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreJwtToken"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	typedef_enum "yfapi/typedef/enum"
	common_data "yfapi/typedef/redisKey"
	"yfapi/typedef/response"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
)

// Auth
//
//	@Description:	权限认证
//	@return			gin.HandlerFunc
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO 进行token验证
		token := c.GetHeader("Authorization")
		claims, err := coreJwtToken.Decode(token, []byte(coreConfig.GetHotConf().JwtSecret))
		if err != nil {
			response.FailResponse(c, error2.I18nError{
				Code: error2.ErrorCodeToken,
				Msg:  nil,
			})
			c.Abort()
			return
		}
		userId := claims.UserId
		clientType := claims.ClientType
		mobile := claims.Mobile
		if len(userId) == 0 {
			response.FailResponse(c, error2.I18nError{
				Code: error2.ErrorCodeToken,
				Msg:  nil,
			})
			c.Abort()
			return
		}
		redisKey := common_data.UserLoginInfo("app", userId)
		switch clientType {
		case typedef_enum.ClientTypePc:
			redisKey = common_data.UserLoginInfo("pc", userId)
		case typedef_enum.ClientTypeH5:
			redisKey = common_data.UserLoginInfo("h5", userId)
		}
		redisToken := coreRedis.GetUserRedis().Get(c, redisKey).Val()
		if token != redisToken {
			response.FailResponse(c, error2.I18nError{
				Code: error2.ErrorCodeToken,
				Msg:  nil,
			})
			c.Abort()
			return
		}
		//刷新token
		go coreRedis.GetUserRedis().Expire(c, redisKey, 7*24*time.Hour)

		c.Set("userId", userId)
		c.Set("clientType", clientType)
		c.Set("mobile", mobile)
		c.Next()
	}
}

// AuthInner
//
//	@Description: 内网调用验签
//	@return gin.HandlerFunc -
func AuthInner() gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.GetHeader("signature")
		timestamp := c.GetHeader("timestamp")
		currSign := easy.AesDecrypt(sign, coreConfig.GetHotConf().InnerSecret.AESEncryptKey1)
		if timestamp != easy.AesDecrypt(currSign, coreConfig.GetHotConf().InnerSecret.AESEncryptKey2) {
			response.FailResponse(c, error2.I18nError{
				Code: error2.ErrorCodeToken,
				Msg:  nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置允许跨域请求的头部
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, X-Request-Id, X-Forwarded-For, X-Real-IP, appVersion, channel, machineCode, platform")

		// 允许发送Cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 允许前端直接读取
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")

		// 响应OPTIONS请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next() // 继续处理请求
		}
	}
}
