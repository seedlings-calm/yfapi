package router

import (
	"github.com/gin-gonic/gin"
	"sync"
	v1_login "yfapi/app/handle/v1/login"
	"yfapi/app/middle"
	"yfapi/internal/engine"
)

func SetLoginRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1")
	group.Use(middle.CheckHeaderData())
	{
		group.POST("/loginByCode", BaseRouter(oriEngine, wg, v1_login.LoginByCode))
		group.POST("/sendSms", BaseRouter(oriEngine, wg, v1_login.SendSms))
		group.POST("/loginByPass", BaseRouter(oriEngine, wg, v1_login.LoginByPass))
		group.GET("/getChooseUser", BaseRouter(oriEngine, wg, v1_login.GetChooseUser))
		group.POST("/chooseUserLogin", BaseRouter(oriEngine, wg, v1_login.ChooseUserLogin))
		group.POST("/forgetPassCheck", BaseRouter(oriEngine, wg, v1_login.ForgetPassCheck))
		group.POST("/forgetPass", BaseRouter(oriEngine, wg, v1_login.ForgetPass))
		group.POST("/loginCheck", BaseRouter(oriEngine, wg, v1_login.LoginCheck)) //登录检测
	}
}
