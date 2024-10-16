package router

import (
	"sync"
	v1_room "yfapi/app/handle/v1/room"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

func SetRoomRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.GET("/list", BaseRouter(oriEngine, wg, v1_room.RoomList))
		group.GET("/room/info", BaseRouter(oriEngine, wg, v1_room.RoomInfo))

		group.GET("/searchAll", BaseRouter(oriEngine, wg, v1_room.SearchAll)) // 搜索用户、聊天室、直播间信息

		group.POST("/room/update", BaseRouter(oriEngine, wg, v1_room.RoomUpdate))
		// group.POST("/room/lock", BaseRouter(oriEngine, wg, v1_room.RoomLock))

		//检查是否可以加入房间
		group.POST("/checkRoom", BaseRouter(oriEngine, wg, v1_room.CheckRoom))
		group.GET("/room/checkIsRoom", BaseRouter(oriEngine, wg, v1_room.CheckIsRoom))

		//加入房间
		group.POST("/joinRoom", BaseRouter(oriEngine, wg, v1_room.JoinRoom))
		//离开房间
		group.POST("/leaveRoom", BaseRouter(oriEngine, wg, v1_room.LeaveRoom))
		//获取权限菜单
		group.GET("/authMenu", BaseRouter(oriEngine, wg, v1_room.GetAuthMenu))

		group.POST("/room/getInformationCard", BaseRouter(oriEngine, wg, v1_room.GetInformationCardByUserId))
		// 房间操作
		group.POST("/room/execCommand", BaseRouter(oriEngine, wg, v1_room.ExecCommand))
		// group.POST("/room/doBlackout", BaseRouter(oriEngine, wg, v1_room.DoBlackOut))
		group.GET("/room/blacklist", BaseRouter(oriEngine, wg, v1_room.BlackList))
		// group.POST("/room/delBlackout", BaseRouter(oriEngine, wg, v1_room.DelBlackOut))
		group.GET("/room/upSeatApplyList", BaseRouter(oriEngine, wg, v1_room.GetRoomUpSeatApplyList))

		group.GET("/room/onlineUsers", BaseRouter(oriEngine, wg, v1_room.OnlineUsers))
		group.GET("/room/dayUsers", BaseRouter(oriEngine, wg, v1_room.DayUsers))
		group.GET("/room/highGradeUsers", BaseRouter(oriEngine, wg, v1_room.HighGradeUsers))
		group.GET("/room/highGradeUsersCount", BaseRouter(oriEngine, wg, v1_room.HighGradeUsersCount))

		group.GET("/room/avToken", BaseRouter(oriEngine, wg, v1_room.AvToken))

		// group.POST("/room/kickOut", BaseRouter(oriEngine, wg, v1_room.KickOut))
		group.POST("/reporting", BaseRouter(oriEngine, wg, v1_room.ReportingCenter))
		group.GET("/reportingTypes", BaseRouter(oriEngine, wg, v1_room.ReportingCenterTypes))

		group.GET("/room/freeUpSeat", BaseRouter(oriEngine, wg, v1_room.FreeUpSeat))
		group.GET("/room/muteLocalSeat", BaseRouter(oriEngine, wg, v1_room.MuteLocalSeat))

		group.GET("/room/userMuteList", BaseRouter(oriEngine, wg, v1_room.UserMuteList))
		group.GET("/room/holdUpSeatList", BaseRouter(oriEngine, wg, v1_room.GetHoldUpSeatUserList))
		group.POST("/room/apply", BaseRouter(oriEngine, wg, v1_room.ApplyRoom))
		group.GET("/room/applyInfo", BaseRouter(oriEngine, wg, v1_room.ApplyRoomInfo))
		group.GET("/room/extra", BaseRouter(oriEngine, wg, v1_room.GetRoomExtraInfo))

		group.POST("/room/admin/add", BaseRouter(oriEngine, wg, v1_room.RoomAdminAdd))
		group.POST("/room/admin/delete", BaseRouter(oriEngine, wg, v1_room.RoomAdminDelete))
		group.GET("/room/admin/list", BaseRouter(oriEngine, wg, v1_room.GetRoomAdminList))

		group.GET("/room/getBgs", BaseRouter(oriEngine, wg, v1_room.GetBgs))
		group.GET("/room/getRoomBgs", BaseRouter(oriEngine, wg, v1_room.GetRoomBgs))
		group.GET("/room/setRoomBgs", BaseRouter(oriEngine, wg, v1_room.SetRoomBgs))

		group.GET("/room/advertising", BaseRouter(oriEngine, wg, v1_room.RoomAdvertising))

		//直播间观看时长上报
		group.POST("/room/retentionTime", BaseRouter(oriEngine, wg, v1_room.RoomRetentionTime))
		//直播间在麦时长上报
		group.POST("/room/onMicTime", BaseRouter(oriEngine, wg, v1_room.RoomOnMicTime))
	}
}
