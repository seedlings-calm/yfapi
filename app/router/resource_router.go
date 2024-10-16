package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_resource "yfapi/app/handle/v1/resource"
	"yfapi/app/middle"
	"yfapi/internal/engine"
)

func SetResourceRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/resource")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.GET("/list", BaseRouter(oriEngine, wg, v1_resource.ResourceList))
		group.GET("/getUploadToken", BaseRouter(oriEngine, wg, v1_resource.GetUploadToken))
	}
}
