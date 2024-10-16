package inner

import (
	"yfapi/app/handle"
	service_im "yfapi/internal/service/im"
	request_message "yfapi/typedef/request/message"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
)

// @Summary 发送系统消息
// @Description
// @Tags 内部调用
// @Accept json
// @Produce json
// @Param req  body request_message.SystemMsgReq true  "参数"
// @Success 0 {object} response.Response{}
// @Router /v1/inner/sendSystemMsg [post]
func SendSystemMsg(c *gin.Context) {
	req := new(request_message.SystemMsgReq)
	handle.BindBody(c, req)
	new(service_im.ImNoticeService).SendSystematicMsg(c, req.Title, req.Img, req.Content, req.Link, req.H5Content, req.ToUserId)
	response.SuccessResponse(c, "")
}

// @Summary 发送官方公告
// @Description
// @Tags 内部调用
// @Accept json
// @Produce json
// @Param req  body request_message.SystemMsgReq true  "参数"
// @Success 0 {object} response.Response{}
// @Router /v1/inner/sendSystemMsg [post]
func SendOfficialMsg(c *gin.Context) {
	req := new(request_message.SystemMsgReq)
	handle.BindBody(c, req)
	err := new(service_im.ImNoticeService).SendOfficialMsg(c, req.Title, req.Img, req.Content, req.Link, req.H5Content, req.ToUserId)
	if err != nil {
		response.FailResponse(c, err.Error())
		return
	}
	response.SuccessResponse(c, "")
}

// @Summary 推送用户im通知
// @Description
// @Tags 内部调用
// @Accept json
// @Produce json
// @Param req  body request_message.SendCommonMsgReq true  "参数"
// @Success 0 {object} response.Response{}
// @Router /v1/inner/sendImNotice [post]
func SendImNotice(c *gin.Context) {
	req := new(request_message.SendCommonMsgReq)
	handle.BindBody(c, req)
	new(service_im.ImCommonService).Send(c, req.FromUserId, req.ToUserId, req.RoomId, req.MsgType, req.MsgData, req.Code)
	response.SuccessResponse(c, "")
}
