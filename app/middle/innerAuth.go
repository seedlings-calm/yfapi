package middle

import (
	"github.com/gin-gonic/gin"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/util/easy"
)

// 判定内网ip，只允许内网访问
func AuthInnerIp() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIp := c.ClientIP()
		if coreConfig.GetHotConf().ENV != "dev" {
			if !easy.InArray(clientIp, coreConfig.GetHotConf().InnerIp) {
				coreLog.Info("非法访问IP：%s", clientIp)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
