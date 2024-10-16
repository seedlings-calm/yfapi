package im

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	service_im "yfapi/internal/service/im"
	service_room "yfapi/internal/service/room"
	request_im "yfapi/typedef/request/im"
	"yfapi/typedef/response"
)

// SendOneTextMsg
//
//	@Summary	发送单聊文本消息
//	@Schemes
//	@Description	发送单聊文本消息
//	@Tags			发送消息
//	@Param			req	body	request_im.SendOneTextMsgReq	true	"单聊参数"
//	@Accept			json
//	@Produce		json
//	@Router			/v1/sendOneTextMsg [post]
func SendOneTextMsg(c *gin.Context) {
	req := new(request_im.SendOneTextMsgReq)
	handle.BindBody(c, req)
	data := handle.GetTokenData(c)
	ser := new(service_im.ImOneService)
	msg := ser.SendTextMsg(c, req.Content, data.UserId, req.ToUserId, data.ClientType, req.Extra)
	response.SuccessResponse(c, msg)
}

// SendOneImgMsg
//
//	@Summary	发送单聊图片消息
//	@Schemes
//	@Description	发送单聊图片消息
//	@Tags			发送消息
//
// @Param req body request_im.SendOneImgMsgReq	true	"单聊参数"
//
//	@Accept			json
//	@Produce		json
//	@Router			/v1/sendOneImgMsg [post]
func SendOneImgMsg(c *gin.Context) {
	req := new(request_im.SendOneImgMsgReq)
	handle.BindBody(c, req)
	data := handle.GetTokenData(c)
	ser := new(service_im.ImOneService)
	msg := ser.SendImgMsg(c, req.Content, data.UserId, req.ToUserId, req.Width, req.Height, data.ClientType, req.Extra)
	response.SuccessResponse(c, msg)
}

// SendOneAudioMsg
//
//	@Summary	发送单聊音频消息
//	@Schemes
//	@Description	发送单聊音频消息
//	@Tags			发送消息
//	@Param			req	body	request_im.SendOneAudioReq	true	"单聊参数"
//	@Accept			json
//	@Produce		json
//	@Router			/v1/sendOneAudioMsg [post]
func SendOneAudioMsg(c *gin.Context) {
	req := new(request_im.SendOneAudioReq)
	handle.BindBody(c, req)
	data := handle.GetTokenData(c)
	ser := new(service_im.ImOneService)
	msg := ser.SendAudioMsg(c, req.Content, data.UserId, req.ToUserId, req.Length, data.ClientType, req.Extra)
	response.SuccessResponse(c, msg)
}

// SendPublicTextMsg
//
//	@Description: 发送公屏文本消息
func SendPublicTextMsg(c *gin.Context) {
	req := new(request_im.SendPublicTextMsgReq)
	handle.BindBody(c, req)
	data := handle.GetTokenData(c)
	ser := new(service_im.ImPublicService)
	msg := ser.SendTextMsg(c, req.Content, data.UserId, req.ToUserId, req.RoomId, data.ClientType, req.Extra)
	// 增加房间热度
	if len(req.RoomId) > 0 {
		go service_room.UpdateRoomHotByChat(req.RoomId, data.UserId)
	}
	response.SuccessResponse(c, msg)
}

// SendPublicImgMsg
//
//	@Description: 发送公屏图片消息
func SendPublicImgMsg(c *gin.Context) {
	req := new(request_im.SendPublicImgMsgReq)
	handle.BindBody(c, req)
	data := handle.GetTokenData(c)
	ser := new(service_im.ImPublicService)
	msg := ser.SendImgMsg(c, req.Content, data.UserId, req.ToUserId, req.RoomId, req.Width, req.Height, data.ClientType, req.Extra)
	response.SuccessResponse(c, msg)
}
