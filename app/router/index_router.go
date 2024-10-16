package router

import (
	"sync"
	v1_index "yfapi/app/handle/v1/index"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

// 设置app首页路由
func SetIndexRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/index")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.GET("/getNavigation", BaseRouter(oriEngine, wg, v1_index.GetNavigation)) //获取首页导航
		group.GET("/getCollect", BaseRouter(oriEngine, wg, v1_index.GetCollect))       //获取首页收藏
		group.GET("/getRecommend", BaseRouter(oriEngine, wg, v1_index.GetRecommend))   // 获取首页推荐
		group.GET("/getRoomsByPC", BaseRouter(oriEngine, wg, v1_index.GetRoomsByPC))   // 获取首页推荐

		group.GET("/top", BaseRouter(oriEngine, wg, v1_index.GetTopRooms))                //
		group.GET("/topMsg", BaseRouter(oriEngine, wg, v1_index.TopMsg))                  //
		group.GET("/menuList", BaseRouter(oriEngine, wg, v1_index.GetAppMenuSettingList)) // app菜单配置列表
	}
}
