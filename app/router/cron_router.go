package router

import (
	"sync"
	v1_cron "yfapi/app/handle/v1/cron"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

func SetCronRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/cron")
	group.Use(middle.AuthInnerIp())
	{
		//测试定时任务
		group.POST("/test", BaseRouter(oriEngine, wg, v1_cron.TestTask))
		group.POST("/wheatTime", BaseRouter(oriEngine, wg, v1_cron.WheatTimeCron))
	}
}
