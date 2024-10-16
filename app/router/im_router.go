package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_im "yfapi/app/handle/v1/im"
	"yfapi/app/middle"
	"yfapi/internal/engine"
)

func SetImRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.POST("/sendOneTextMsg", BaseRouter(oriEngine, wg, v1_im.SendOneTextMsg))
		group.POST("/sendOneImgMsg", BaseRouter(oriEngine, wg, v1_im.SendOneImgMsg))
		group.POST("/sendOneAudioMsg", BaseRouter(oriEngine, wg, v1_im.SendOneAudioMsg))
		group.GET("/sessionList", BaseRouter(oriEngine, wg, v1_im.GetSessionList))
		group.GET("/messageList", BaseRouter(oriEngine, wg, v1_im.GetUserMessageList))
		group.POST("/sendPublicTextMsg", BaseRouter(oriEngine, wg, v1_im.SendPublicTextMsg))
		group.POST("/sendPublicImgMsg", BaseRouter(oriEngine, wg, v1_im.SendPublicImgMsg))
		group.POST("/messageRead", BaseRouter(oriEngine, wg, v1_im.MessageRead))
		group.POST("/topSession", BaseRouter(oriEngine, wg, v1_im.TopSession))
		group.POST("/unTopSession", BaseRouter(oriEngine, wg, v1_im.UnTopSession))
		group.POST("/delSession", BaseRouter(oriEngine, wg, v1_im.DelSession))
		group.POST("/clearChatHistory", BaseRouter(oriEngine, wg, v1_im.ClearChatHistory))
		group.POST("/allMessageRead", BaseRouter(oriEngine, wg, v1_im.MessageReadAll)) //所有消息已读
	}
}
