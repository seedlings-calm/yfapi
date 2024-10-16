package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_recharge "yfapi/app/handle/v1/recharge"
	"yfapi/app/middle"
	"yfapi/internal/engine"
)

func SetRechargeRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/recharge")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.POST("/test", BaseRouter(oriEngine, wg, v1_recharge.UserRechargeTest)) // 测试充值
		group.POST("/iosIap", BaseRouter(oriEngine, wg, v1_recharge.IosIap))         // 苹果内购充值
		group.POST("/wxAppPay", BaseRouter(oriEngine, wg, v1_recharge.WxAppPay))     // 微信App充值
		group.POST("/aliAppPay", BaseRouter(oriEngine, wg, v1_recharge.AliAppPay))   // 支付宝App充值

		group.POST("/aggregationPay", BaseRouter(oriEngine, wg, v1_recharge.AggregationPay)) // 第三方聚合支付接口

		group.POST("/rechargeResult", BaseRouter(oriEngine, wg, v1_recharge.RechargeResult)) // 支付结果查询
	}
}
