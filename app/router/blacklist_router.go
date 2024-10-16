package router

import (
	"sync"
	"yfapi/app/handle/v1/user"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

func SetBlacklistRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/blacklist")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.GET("/add", BaseRouter(oriEngine, wg, user.AddBlacklist))
		group.GET("/del", BaseRouter(oriEngine, wg, user.DelBlacklist))
		group.GET("/userList", BaseRouter(oriEngine, wg, user.GetUserBlacklist))
	}
}
