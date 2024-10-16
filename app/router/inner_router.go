package router

import (
	"sync"
	v1_inner "yfapi/app/handle/v1/inner"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

func SetInnerRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/inner")
	group.Use(middle.AuthInnerIp())
	{
		group.GET("/getUploadToken", BaseRouter(oriEngine, wg, v1_inner.GetUploadToken))
		group.POST("/sendSms", BaseRouter(oriEngine, wg, v1_inner.SendSms))
		group.POST("/doRoom", BaseRouter(oriEngine, wg, v1_inner.DoRoom))
		group.POST("/sendSystemMsg", BaseRouter(oriEngine, wg, v1_inner.SendSystemMsg))
		group.POST("/sendImRoomMsg", BaseRouter(oriEngine, wg, v1_inner.SendImRoomNotice))
		group.POST("/sendOfficialMsg", BaseRouter(oriEngine, wg, v1_inner.SendOfficialMsg))
		//CommonSendMsg公共发送消息
		group.POST("/sendImNotice", BaseRouter(oriEngine, wg, v1_inner.SendImNotice))

		// 定时任务相关
		group.GET("/task/autoDeleteUser", BaseRouter(oriEngine, wg, v1_inner.AutoDeleteUser))

		group.POST("/accountChange", BaseRouter(oriEngine, wg, v1_inner.AccountChange))
		group.POST("/operationChangeAccount", BaseRouter(oriEngine, wg, v1_inner.OperationChangeAccount)) // 运营后台变动账户

		//代付接口 提现
		group.POST("/withdraw", BaseRouter(oriEngine, wg, v1_inner.AnotherPay))
		//排行榜结算
		group.POST("/rankListDay", BaseRouter(oriEngine, wg, v1_inner.RankListDay))
		group.POST("/rankListWeek", BaseRouter(oriEngine, wg, v1_inner.RankListWeek))
		group.POST("/rankListMonth", BaseRouter(oriEngine, wg, v1_inner.RankListMonth))
	}
}
