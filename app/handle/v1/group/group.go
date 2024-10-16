package group

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	service_im "yfapi/internal/service/im"
	request_group "yfapi/typedef/request/group"
	"yfapi/typedef/response"
)

// 发送文本消息
func SendTextMsg(c *gin.Context) {
	req := new(request_group.SendTextMsgReq)
	handle.BindBody(c, req)
	data := handle.GetTokenData(c)
	ser := new(service_im.ImGroupService)
	msg := ser.SendTextMsg(c, req.Content, data.UserId, req.ToUserId, req.RoomId, data.ClientType, req.Extra)
	response.SuccessResponse(c, msg)
}

// 发送图片消息
func SendImgMsg(c *gin.Context) {
	req := new(request_group.SendImgMsgReq)
	handle.BindBody(c, req)
	data := handle.GetTokenData(c)
	ser := new(service_im.ImGroupService)
	msg := ser.SendImgMsg(c, req.Content, data.UserId, req.ToUserId, req.RoomId, req.Width, req.Height, data.ClientType, req.Extra)
	response.SuccessResponse(c, msg)
}
