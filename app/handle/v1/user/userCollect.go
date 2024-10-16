package user

import (
	error2 "yfapi/i18n/error"
	"yfapi/internal/logic"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
)

// @Summary 收藏房间
// @Description
// @Tags 收藏房间功能
// @Accept json
// @Produce json
// @Param  roomId query string  true "房间ID"
// @Success 0 {object} response.Response{}
// @Router /v1/user/collect [get]
func AddCollect(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := new(logic.UserCollect)
	response.SuccessResponse(c, service.AddCollect(roomId, c))
}

// @Summary 取消收藏
// @Description
// @Tags 收藏房间功能
// @Accept json
// @Produce json
// @Param  roomId query string  true "房间ID"
// @Success 0 {object} response.Response{}
// @Router /v1/user/collect [put]
func DelCollect(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := new(logic.UserCollect)
	response.SuccessResponse(c, service.DelCollect(roomId, c))
}
