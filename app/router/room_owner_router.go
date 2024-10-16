package router

import (
	"sync"
	v1_room "yfapi/app/handle/v1/roomOwner"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

// 房主后台的专属路由
// SetRoomOwnerRouter 房主后台路由
func SetRoomOwnerRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/roomOwner")
	//手机验证码登录
	group.POST("/login", BaseRouter(oriEngine, wg, v1_room.SmsCodeLogin))
	//发送短信
	group.POST("/sendSms", BaseRouter(oriEngine, wg, v1_room.GetLoginMobileSMS))
	//国家区号
	//group.GET("/regionCode", BaseRouter(oriEngine, wg, v1_login.RegionCodeList))
	group.Use(middle.RoomownerAuth())
	{
		group.GET("/user/search", BaseRouter(oriEngine, wg, v1_room.RoomSearchUser)) //搜索用户信息
		//房主房间列表信息
		group.POST("/room/list", BaseRouter(oriEngine, wg, v1_room.RoomListInfo))                      //登录后房间列表数据
		group.POST("/room/detail", BaseRouter(oriEngine, wg, v1_room.RoomInfo))                        //获取单个房间信息
		group.POST("/room/base", BaseRouter(oriEngine, wg, v1_room.RoomBase))                          //聊天室概况
		group.GET("/room/dashboard", BaseRouter(oriEngine, wg, v1_room.RoomDashBoard))                 //房间数据
		group.GET("/room/dashboardMoney", BaseRouter(oriEngine, wg, v1_room.RoomDashBoardMoneysChart)) //房间数据-房间流水
		group.GET("/room/dashboardTimes", BaseRouter(oriEngine, wg, v1_room.RoomDashBoardTimesChart))  //房间数据-有效观看人数
		//管理员管理
		group.POST("/admin/list", BaseRouter(oriEngine, wg, v1_room.RoomAdminList))     //管理员列表
		group.POST("/admin/remove", BaseRouter(oriEngine, wg, v1_room.RoomAdminRemove)) //移除管理员
		group.POST("/admin/add", BaseRouter(oriEngine, wg, v1_room.RoomAdminAdd))       //添加管理员
		// 账户相关
		group.POST("/exchange", BaseRouter(oriEngine, wg, v1_room.ExchangeDiamond))        //兑换钻石
		group.POST("/accountBill", BaseRouter(oriEngine, wg, v1_room.GetAccountBillList))  //账户交易明细
		group.POST("/bindBank", BaseRouter(oriEngine, wg, v1_room.BindBank))               //房主绑定银行卡
		group.GET("/accountInfo", BaseRouter(oriEngine, wg, v1_room.RoomAccountInfo))      //玩家房间账户信息
		group.POST("/withdrawApply", BaseRouter(oriEngine, wg, v1_room.RoomWithdrawApply)) //房主提现申请
		//从业者
		group.POST("/practitioner/list", BaseRouter(oriEngine, wg, v1_room.RoomPractitionerList))       //从业者列表
		group.POST("/practitioner/add", BaseRouter(oriEngine, wg, v1_room.RoomPractitionerAdd))         //添加从业者
		group.POST("/practitioner/remove", BaseRouter(oriEngine, wg, v1_room.RoomPractitionerRemove))   //移除从业者
		group.POST("/practitioner/reSave", BaseRouter(oriEngine, wg, v1_room.RoomPractitionerReSave))   //再次提交
		group.POST("/practitioner/invalid", BaseRouter(oriEngine, wg, v1_room.RoomPractitionerInvalid)) //申请作废

		group.GET("/member/list", BaseRouter(oriEngine, wg, v1_room.PersonList))         //成员列表
		group.GET("/member/detail", BaseRouter(oriEngine, wg, v1_room.PersonListDetail)) //成员受赏详情
	}
}
