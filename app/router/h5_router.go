package router

import (
	"sync"
	v1_h5 "yfapi/app/handle/v1/h5"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

func SetH5Router(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/h5")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.POST("/deleteUserApply", BaseRouter(oriEngine, wg, v1_h5.DeleteUserApply))

	}

	{
		// 公会相关
		group.POST("/guild/join", BaseRouter(oriEngine, wg, v1_h5.JoinGuild))                       // 申请加入公会
		group.POST("/guild/info", BaseRouter(oriEngine, wg, v1_h5.GuildInfo))                       // 工会详情
		group.POST("/guild/memberList", BaseRouter(oriEngine, wg, v1_h5.GetGuildMemberList))        // 查询工会成员列表
		group.POST("/guild/quitApply", BaseRouter(oriEngine, wg, v1_h5.QuitGuildApply))             // 退出公会申请
		group.POST("/guild/quitApplyCancel", BaseRouter(oriEngine, wg, v1_h5.QuitGuildApplyCancel)) // 取消退出公会申请
		group.GET("/guild/penaltyDetail", BaseRouter(oriEngine, wg, v1_h5.GetGuildPenaltyDetail))   // 违约金详情
		group.POST("/guild/payPenalty", BaseRouter(oriEngine, wg, v1_h5.PayGuildPenalty))           // 缴纳违约金

		group.GET("/getQuestion/:types", BaseRouter(oriEngine, wg, v1_h5.GetQuestion))   //获取考核题库
		group.POST("/pullAnswer", BaseRouter(oriEngine, wg, v1_h5.PullAnswer))           //上传基础考题答案20
		group.POST("/pullShortAnswer", BaseRouter(oriEngine, wg, v1_h5.PullShortAnswer)) //上传简答考题答案
		group.POST("/pullMusic", BaseRouter(oriEngine, wg, v1_h5.PullMusic))             //接收音乐人考核答案
		group.GET("/applyJoinResult", BaseRouter(oriEngine, wg, v1_h5.ApplyJoinResult))  //申请入驻结果
		group.GET("/cerdAuth/:types", BaseRouter(oriEngine, wg, v1_h5.CerdAuth))         //从业者身份资格
		group.GET("/getUserInfo", BaseRouter(oriEngine, wg, v1_h5.GetUserInfo))          //获取用户基本信息
	}

	{
		group.POST("/level/base", BaseRouter(oriEngine, wg, v1_h5.LevelBaseInfo))        //获取等级基本信息
		group.POST("/level/config", BaseRouter(oriEngine, wg, v1_h5.GetLevelConfigList)) //获取等级配置列表
	}
}
