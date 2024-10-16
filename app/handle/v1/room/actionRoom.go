package room

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	"yfapi/internal/logic"
	request_room "yfapi/typedef/request/room"
	"yfapi/typedef/response"
	response_room "yfapi/typedef/response/room"
)

// @Summary 进入房间
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  data body  request_room.ActionRoomReq true "加入房间参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.ChatroomDTO{}
// @Router /v1/joinRoom [post]
func JoinRoom(c *gin.Context) {
	req := new(request_room.ActionRoomReq)
	handle.BindBody(c, req)
	tokenData := handle.GetTokenData(c)
	actionRoom := new(logic.ActionRoom)
	response.SuccessResponse(c, actionRoom.JoinRoom(c, tokenData.UserId, tokenData.ClientType, req))
}

// @Summary 离开房间
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  data body  request_room.ActionRoomReq true "加入房间参数"
// @Success 0 {object} response.Response{}
// @Router /v1/leaveRoom [post]
func LeaveRoom(c *gin.Context) {
	req := new(request_room.ActionRoomReq)
	handle.BindBody(c, req)
	tokenData := handle.GetTokenData(c)
	actionRoom := new(logic.ActionRoom)
	response.SuccessResponse(c, actionRoom.LeaveRoom(c, tokenData.UserId, req.RoomId, tokenData.ClientType))
}

// @Summary 进入房间之前的检查
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  data body  request_room.ActionRoomReq   true "进入的房间ID"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.CheckRoomResponse{}
// @Router /v1/checkRoom [post]
func CheckRoom(c *gin.Context) {
	req := new(request_room.ActionRoomReq)
	handle.BindBody(c, req)
	tokenData := handle.GetTokenData(c)
	actionRoom := new(logic.ActionRoom)
	var resp response_room.CheckRoomResponse
	resp, _ = actionRoom.CheckRoom(c, tokenData.UserId, req.RoomId, false)
	response.SuccessResponse(c, resp)
}

// @Summary 检测是否在此房间
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  roomId query string  true "房间ID"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.CheckIsRoomResponse{}
// @Router /v1/room/checkIsRoom [get]
func CheckIsRoom(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		response.FailResponse(c, error2.I18nError{Code: error2.ErrorCodeParam})
	}
	actionRoom := new(logic.ActionRoom)
	resp := actionRoom.CheckIsRoom(c, roomId)
	response.SuccessResponse(c, resp)
}

// 获取房间的权限菜单
func GetAuthMenu(c *gin.Context) {
	req := new(request_room.RoomAuthMenuReq)
	handle.BindQuery(c, req)
	userId := handle.GetUserId(c)
	roomAuth := new(logic.RoomAuth)
	response.SuccessResponse(c, roomAuth.GetAuthMenu(c, userId, req))
}

// @Summary 房间操作
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  data body  request_room.ExecCommandReq true "房间操作参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.ExecCommandRes{}
// @Router /v1/room/execCommand [post]
func ExecCommand(c *gin.Context) {
	req := new(request_room.ExecCommandReq)
	handle.BindBody(c, req)
	actionRoom := new(logic.ActionRoom)
	response.SuccessResponse(c, actionRoom.ExecCommand(c, req))
}

// @Summary 申请上麦列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  req	query	request_room.UpSeatApplyListReq	true	"申请上麦列表参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.UpSeatApplyListRes{}
// @Router /v1/room/upSeatApplyList [get]
func GetRoomUpSeatApplyList(c *gin.Context) {
	req := new(request_room.UpSeatApplyListReq)
	handle.BindQuery(c, req)
	actionRoom := new(logic.ActionRoom)
	response.SuccessResponse(c, actionRoom.GetRoomUpSeatApplyList(c, req))
}

// FreeUpSeat
//
//	@Description: 自由上麦
func FreeUpSeat(c *gin.Context) {
	req := new(request_room.FreeUpSeatReq)
	handle.BindQuery(c, req)
	actionRoom := new(logic.ActionRoom)
	actionRoom.FreeUpSeat(c, req)
	response.SuccessResponse(c, "")
}

// MuteLocalSeat
//
//	@Description: 本地静音
func MuteLocalSeat(c *gin.Context) {
	req := new(request_room.MutLocalSeatReq)
	handle.BindQuery(c, req)
	actionRoom := new(logic.ActionRoom)
	ok := actionRoom.MuteLocalSeat(c, req)
	response.SuccessResponse(c, ok)
}

// @Summary 查询可抱用户上麦列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  req	query	request_room.UpSeatUserListReq	true	"申请上麦列表参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.HoldUpSeatUserListRes{}
// @Router /v1/room/holdUpSeatList [get]
func GetHoldUpSeatUserList(c *gin.Context) {
	req := new(request_room.UpSeatUserListReq)
	handle.BindQuery(c, req)
	actionRoom := new(logic.ActionRoom)
	response.SuccessResponse(c, actionRoom.GetHoldUpSeatUserList(c, req))
}
