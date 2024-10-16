package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"sync"
	"yfapi/app/middle"
	"yfapi/core/coreConfig"
	"yfapi/i18n"
	"yfapi/internal/engine"

	swaggerfiles "github.com/swaggo/files"
)

func SetupRouter(oriEngine *engine.OriEngine, wg *sync.WaitGroup) *gin.Engine {
	router := gin.New() //实例化gin框架
	router.Use(middle.Cors())
	router.Use(middle.Recover())
	//router.Use(middle.RequestId())
	router.Use(middle.Logger())
	router.Use(i18n.I18nHandler())
	group := router.Group("/api")
	if coreConfig.GetHotConf().ENV != "pro" {
		router.Static("/docs", "./docs")
		sysSwaggerRouter(router)
	}
	SetLoginRouter(oriEngine, group, wg)
	SetLUserRouter(oriEngine, group, wg)
	SetImRouter(oriEngine, group, wg)
	SetInnerRouter(oriEngine, group, wg)
	SetResourceRouter(oriEngine, group, wg)
	SetRoomRouter(oriEngine, group, wg)
	SetGiftRouter(oriEngine, group, wg)
	SetCallbackRouter(oriEngine, group, wg)
	SetH5Router(oriEngine, group, wg)
	SetH5GoodsRouter(oriEngine, group, wg)
	SetIndexRouter(oriEngine, group, wg)
	SetAgentRouter(oriEngine, group, wg)
	SetRechargeRouter(oriEngine, group, wg)
	SetBlacklistRouter(oriEngine, group, wg)
	SetCronRouter(oriEngine, group, wg)
	//SexIndexRouter(oriEngine, group, wg)
	//SetPractitionerRouter(oriEngine, group, wg)
	SetGuildRouter(oriEngine, group, wg)
	SetRoomOwnerRouter(oriEngine, group, wg)
	SetGroupRouter(oriEngine, group, wg)
	SetWebsiteRouter(oriEngine, group, wg)
	SetRankListRouter(oriEngine, group, wg)
	return router
}
func BaseRouter(oriEngine *engine.OriEngine, wg *sync.WaitGroup, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		wg.Add(1)
		defer wg.Done()
		handlerFunc(context)
	}
}

func sysSwaggerRouter(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler(), ginSwagger.URL("/docs/swagger.json"), ginSwagger.InstanceName("yfapi")))
}
