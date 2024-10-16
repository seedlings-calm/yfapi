package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_callback "yfapi/app/handle/v1/callback"
	"yfapi/internal/engine"
)

func SetCallbackRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1")
	{
		//动态图片审核回调
		group.POST("/momentsImageCallback/:param", BaseRouter(oriEngine, wg, v1_callback.MomentsImageCallback))
		group.POST("/momentsVideoCallback/:param", BaseRouter(oriEngine, wg, v1_callback.MomentsVideoCallback))
		group.POST("/signAudioCallback/:param", BaseRouter(oriEngine, wg, v1_callback.SignAudioCallback))

		//聚合支付回调
		group.POST("/pay/aggregation", BaseRouter(oriEngine, wg, v1_callback.AggregationCallback))

		//代付回调
		group.POST("/pay/another", BaseRouter(oriEngine, wg, v1_callback.AnotherCallback))
	}
}
