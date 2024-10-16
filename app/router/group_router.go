package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_group "yfapi/app/handle/v1/group"
	"yfapi/app/middle"
	"yfapi/internal/engine"
)

func SetGroupRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		//世界频道设置
		group.GET("group/worldGroupSetting", BaseRouter(oriEngine, wg, v1_group.WorldGroupSetting))
		//世界频道禁言
		group.POST("group/worldGroupMute", BaseRouter(oriEngine, wg, v1_group.WorldGroupMute))
		//通知设置
		group.POST("group/worldNoticeSetting", BaseRouter(oriEngine, wg, v1_group.WorldGroupNoticeSetting))

		//发送群聊文本消息
		group.POST("group/sendTextMsg", BaseRouter(oriEngine, wg, v1_group.SendTextMsg))
		//发送群聊图片消息
		group.POST("group/sendImgMsg", BaseRouter(oriEngine, wg, v1_group.SendImgMsg))
	}
}
