package room

import (
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/logic"
	"yfapi/internal/model"
	"yfapi/internal/service/rankList"
	"yfapi/typedef/redisKey"
	resquest_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	response_room "yfapi/typedef/response/room"
)

//房间相关的资料卡功能

// @Summary 用户资料卡
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  data body resquest_user.UserPractitionerReq{}  true "资料卡参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.InformationCardResponse{}
// @Router /v1/room/getInformationCard [post]
func GetInformationCardByUserId(c *gin.Context) {
	req := new(resquest_user.UserPractitionerReq)
	handle.BindBody(c, req)
	var res response_room.InformationCardResponse
	ac := &logic.ActionRoom{}
	res = ac.GetInformationCardByUserId(c, req)
	response.SuccessResponse(c, res)
}

// @Summary 黑名单列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response_room.BlackListResponse{}
// @Success 0 {object} response.Response{}
// @Router /v1/room/blacklist [get]
func BlackList(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	ac := &logic.ActionRoom{}
	resp := ac.BlackList(c, roomId)

	response.SuccessResponse(c, resp)
}

// @Summary 房间的在线用户列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response_room.OnlineUsersResponse{}
// @Router /v1/room/onlineUsers [get]
func OnlineUsers(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := &logic.ActionRoom{}
	resp := service.OnlineUsers(c, roomId)
	response.SuccessResponse(c, resp)
}

// @Summary 1000贡献榜列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.DayUsersResponse{}
// @Router /v1/room/dayUsers [get]
func DayUsers(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := &logic.ActionRoom{}
	resp := service.DayUsers(c, roomId)
	response.SuccessResponse(c, resp)
}

// @Summary 禁言列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response_room.UserMuteListRes{}
// @Success 0 {object} response.Response{}
// @Router /v1/room/userMuteList [get]
func UserMuteList(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	ac := &logic.UserMute{}
	resp := ac.GetMuteList(c, roomId)
	response.SuccessResponse(c, resp)
}

// @Summary 房间公屏状态检查
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response_room.ChatroomExtra{}
// @Router /v1/room/extra [get]
func GetRoomExtraInfo(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	response.SuccessResponse(c, new(logic.ActionRoom).GetRoomExtraInfo(c, roomId))
}

// @Summary 房间内背景列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_room.GetRoomBgsRes{}
// @Router /v1/room/getRoomBgs [get]
func GetRoomBgs(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	response.SuccessResponse(c, new(logic.ActionRoom).GetRoomBgs(c, roomId))

}

// @Summary 全局房间背景列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} []model.RoomBgsResource{}
// @Router /v1/room/getBgs [get]
func GetBgs(c *gin.Context) {
	var res []model.RoomBgsResource
	res = append(res, new(logic.ActionRoom).GetBgs(c)...)
	response.SuccessResponse(c, res)
}

// @Summary 更换房间背景
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomBgId  query string  true "房间背景ID"
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response.Response{}
// @Router /v1/room/setRoomBgs [get]
func SetRoomBgs(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	bgId := c.Query("roomBgId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	response.SuccessResponse(c, new(logic.ActionRoom).SetRoomBgs(c, roomId, cast.ToInt(bgId)))
}

// @Summary 房间内高等级用户
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response_room.HighGradeUsersResponse{}
// @Router /v1/room/highGradeUsers [get]
func HighGradeUsers(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := &logic.ActionRoom{}
	resp := service.HighGradeUsers(c, roomId)
	response.SuccessResponse(c, resp)
}

// @Summary 高等级用户的统计数据
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param roomId query string  true "房间ID"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.HighGradeUsersCountResponse{}
// @Router /v1/room/highGradeUsersCount [get]
func HighGradeUsersCount(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := &logic.ActionRoom{}
	resp := service.HighGradeUsersCount(c, roomId)
	response.SuccessResponse(c, resp)
}

// 房间广告位
func RoomAdvertising(c *gin.Context) {
	service := &logic.ActionRoom{}
	response.SuccessResponse(c, service.RoomAdvertising(c))
}

// 在房间观看时长上报
// 每分钟上报一次
func RoomRetentionTime(c *gin.Context) {
	userId := handle.GetUserId(c)
	key := redisKey.UserInRoomRetentionTime(userId)
	rd := coreRedis.GetChatroomRedis()
	ok := rd.SetNX(c, key, 1, time.Second*59).Val()
	if ok {
		//排行榜观看时长
		roomId := rd.Get(c, redisKey.UserInWhichRoom(userId, handle.GetClientType(c))).Val()
		if roomId != "" {
			go func() {
				rankList.Instance().Calculate(rankList.CalculateReq{
					FromUserId: userId,
					Types:      "retentionTime",
					RoomId:     roomId,
				})
			}()
		}
	}
	response.SuccessResponse(c, "")
}

// 在麦时长上报
func RoomOnMicTime(c *gin.Context) {
	userId := handle.GetUserId(c)
	key := redisKey.UserInRoomOnMicTime(userId)
	rd := coreRedis.GetChatroomRedis()
	ok := rd.SetNX(c, key, 1, time.Second*179).Val()
	if ok {
		//排行榜观看时长
		roomId := rd.Get(c, redisKey.UserInWhichRoom(userId, handle.GetClientType(c))).Val()
		if roomId != "" {
			go func() {
				rankList.Instance().Calculate(rankList.CalculateReq{
					FromUserId: userId,
					Types:      "onMicTime",
					RoomId:     roomId,
				})
			}()
		}
	}
	response.SuccessResponse(c, "")
}
