package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_agent "yfapi/app/handle/v1/login"
	"yfapi/internal/engine"
)

func SetAgentRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/agent")
	{
		//获取授权token
		group.POST("/h5/token", BaseRouter(oriEngine, wg, v1_agent.H5Token))
	}
}
