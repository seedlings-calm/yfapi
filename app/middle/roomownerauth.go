package middle

import (
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreJwtToken"
	"yfapi/core/coreRedis"
	i18n_err "yfapi/i18n/error"
	common_data "yfapi/typedef/redisKey"
	"yfapi/typedef/response"
)

// GuildAuth
//
//	@Description:	权限认证
//	@return			gin.HandlerFunc
func RoomownerAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		claims, err := coreJwtToken.RoomDecode(token, []byte(coreConfig.GetHotConf().JwtSecret))
		if err != nil {
			response.FailResponse(c, i18n_err.ErrorCodeToken)
			c.Abort()
			return
		}

		userId := claims.UserId
		if len(userId) == 0 {
			response.FailResponse(c, i18n_err.ErrorCodeToken)
			c.Abort()
			return
		}
		redisKey := common_data.UserLoginInfo("roomPc", userId)
		redisToken := coreRedis.GetUserRedis().Get(c, redisKey).Val()
		if token != redisToken {
			response.FailResponse(c, i18n_err.ErrorCodeToken)
			c.Abort()
			return
		}
		//刷新token
		go coreRedis.GetUserRedis().Expire(c, redisKey, 7*24*time.Hour)

		c.Set("userId", userId)
		c.Next()
	}
}
