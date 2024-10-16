package room

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	logic_room "yfapi/internal/logic"
	resquest_room "yfapi/typedef/request/room"
	"yfapi/typedef/response"
	response_room "yfapi/typedef/response/room"
)

// @Summary 房间管理员添加
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  data body  resquest_room.RoomAdminAddReq true "房间管理员添加"
// @Success 0 {object} response.Response{}
// @Router /v1/room/admin/add [post]
func RoomAdminAdd(c *gin.Context) {
	req := new(resquest_room.RoomAdminAddReq)
	handle.BindBody(c, req)
	service := new(logic_room.RoomAdmin)
	service.RoomAdminAdd(c, req)
	response.SuccessResponse(c, nil)
}

// @Summary 房间管理员删除
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  data body  resquest_room.RoomAdminDeleteReq true "房间管理员删除"
// @Success 0 {object} response.Response{}
// @Router /v1/room/admin/delete [post]
func RoomAdminDelete(c *gin.Context) {
	req := new(resquest_room.RoomAdminDeleteReq)
	handle.BindBody(c, req)
	service := new(logic_room.RoomAdmin)
	service.RoomAdminDelete(c, req)
	response.SuccessResponse(c, nil)
}

// @Summary 房间管理员列表
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  data body  resquest_room.RoomAdminListReq true "房间管理员列表"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_room.RoomAdminRes{}
// @Router /v1/room/admin/list [get]
func GetRoomAdminList(c *gin.Context) {
	req := new(resquest_room.RoomAdminListReq)
	handle.BindBody(c, req)
	res := new(response_room.RoomAdminRes)
	service := new(logic_room.RoomAdmin)
	res = service.RoomAdminList(c, req)
	response.SuccessResponse(c, res)
}
