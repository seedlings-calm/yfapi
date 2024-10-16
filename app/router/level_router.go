package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_level "yfapi/app/handle/v1/level"
	"yfapi/app/middle"
	"yfapi/internal/engine"
)

func SetLevelRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/level")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		//用户lv列表
		group.GET("/lvList", BaseRouter(oriEngine, wg, v1_level.LvList))
	}
}
