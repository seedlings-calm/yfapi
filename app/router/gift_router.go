package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_gift "yfapi/app/handle/v1/gift"
	"yfapi/app/middle"
	"yfapi/internal/engine"
)

func SetGiftRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		//礼物列表
		group.GET("/giftList", BaseRouter(oriEngine, wg, v1_gift.GiftList))
		group.GET("/giftSourceList", BaseRouter(oriEngine, wg, v1_gift.GiftSourceList))
		group.POST("/sendGift", BaseRouter(oriEngine, wg, v1_gift.SendGift))
	}
}
