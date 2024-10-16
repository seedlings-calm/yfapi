package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	"yfapi/app/handle/v1/rankList"
	"yfapi/app/middle"
	"yfapi/internal/engine"
)

func SetRankListRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		//排行榜
		group.GET("/rank", BaseRouter(oriEngine, wg, rankList.RankList))
	}
}
