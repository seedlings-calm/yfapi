package roomOwner

import (
	"github.com/gin-gonic/gin"
	"yfapi/internal/logic/room"
	"yfapi/typedef/response"
	response_roomowner "yfapi/typedef/response/roomOwner"
)

// RoomAccountInfo
// @Summary 首页账户信息
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Success 0 {object} response_roomowner.RoomAccountInfoRes{}
// @Router /v1/roomOwner/accountInfo [get]
func RoomAccountInfo(context *gin.Context) {
	var res response_roomowner.RoomAccountInfoRes
	res = new(room.RoomAccountInfo).RoomAccountInfoLogic(context)
	response.SuccessResponse(context, res)
}
