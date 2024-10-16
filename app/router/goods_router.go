package router

import (
	"sync"
	"yfapi/app/handle/v1/h5"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

func SetH5GoodsRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/goods")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.GET("/types", BaseRouter(oriEngine, wg, h5.GoodsTypesList))
		group.GET("/type_list", BaseRouter(oriEngine, wg, h5.GoodsListByTypes))
		group.GET("/list", BaseRouter(oriEngine, wg, h5.GoodsAll))
		group.GET("/goods_center", BaseRouter(oriEngine, wg, h5.GoodsListToUser))
		group.GET("/use_goods", BaseRouter(oriEngine, wg, h5.UseGoodsToUser))
		group.GET("/del_goods", BaseRouter(oriEngine, wg, h5.UserGoodsDel))
		group.POST("/buy", BaseRouter(oriEngine, wg, h5.ByGoodsToUser))
		group.GET("/del_reddot", BaseRouter(oriEngine, wg, h5.DelRedDot))
		group.GET("/goods_center_user", BaseRouter(oriEngine, wg, h5.GetGoodsUserAccounts))

	}
}
