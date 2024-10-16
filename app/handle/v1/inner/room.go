package inner

import (
	"yfapi/app/handle"
	"yfapi/internal/logic"
	service_im "yfapi/internal/service/im"
	request_room "yfapi/typedef/request/room"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
)

// @Summary 处理用户退出房间逻辑
// @Description
// @Tags 内部调用
// @Accept json
// @Produce json
// @Param req  body request_room.DoRoomReq true  "参数"
// @Success 0 {object} response.Response{}
// @Router /v1/inner/doRoom [post]
func DoRoom(c *gin.Context) {
	req := new(request_room.DoRoomReq)
	handle.BindBody(c, req)
	//tokenData := handle.GetTokenData(c)
	new(logic.ActionRoom).LeaveRoom(c, req.UserId, req.RoomId, req.ClientType)
	response.SuccessResponse(c, nil)
}

// @Summary 推送im房间信息
// @Description
// @Tags 内部调用
// @Accept json
// @Produce json
// @Param req  body request_room.SendRoomMsgReq true  "参数"
// @Success 0 {object} response.Response{}
// @Router /v1/inner/sendImRoomMsg [post]
func SendImRoomNotice(c *gin.Context) {
	req := new(request_room.SendRoomMsgReq)
	handle.BindBody(c, req)
	new(service_im.ImPublicService).SendCustomMsg(req.RoomId, req.Msg, req.Code)
	response.SuccessResponse(c, nil)
}
