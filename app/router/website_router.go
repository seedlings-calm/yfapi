package router

import (
	"sync"
	v1_recharge "yfapi/app/handle/v1/recharge"
	"yfapi/app/handle/v1/user"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

func SetWebsiteRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/website")
	{
		group.POST("/pay", BaseRouter(oriEngine, wg, v1_recharge.WebsitePay))
		group.GET("/rechargeDiamondOffical", BaseRouter(oriEngine, wg, user.RechargeDiamondForOffical))
		group.GET("/userInfoByUserNo", BaseRouter(oriEngine, wg, user.SearchUserInfoByUserNo))

		group.POST("/rechargeResult", BaseRouter(oriEngine, wg, v1_recharge.RechargeResult)) //支付结果查询
	}
}
