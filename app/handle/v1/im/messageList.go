package im

import (
	"yfapi/app/handle"
	"yfapi/internal/logic"
	"yfapi/typedef/request/message"
	"yfapi/typedef/response"
	_ "yfapi/typedef/response/im"

	"github.com/gin-gonic/gin"
)

// GetUserMessageList
//
// @Summary 获取历史消息
// @Schemes
// @Description 获取历史消息
// @Tags 会话相关
// @Produce json
// @Param	req	query	message.MessageListReq	true	"历史消息参数"

// @Router /v1/messageList [get]
func GetUserMessageList(c *gin.Context) {
	req := message.MessageListReq{}
	handle.BindQuery(c, &req)
	tokenData := handle.GetTokenData(c)
	list := new(logic.Message).GetMessageList(tokenData.UserId, req)
	response.SuccessResponse(c, list)
}
