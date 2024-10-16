package guild

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	loggic_guild "yfapi/internal/logic/guild"
	request_room "yfapi/typedef/request/guild"
	"yfapi/typedef/response"
	_ "yfapi/typedef/response/guild"
)

// RoomTypeList
// @Summary 房间类型列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data	body request_room.RoomTypeReq{}  true "房间申请参数"
// @Success 0 {object} guild.RoomTypeInfo{}
// @Router /v1/guild/chatRoomType [post]
func RoomTypeList(c *gin.Context) {
	req := new(request_room.RoomTypeReq)
	handle.BindBody(c, req)
	service := new(loggic_guild.GuildRoom)
	response.SuccessResponse(c, service.RoomTypeList(req))
}

// RoomList
// @Summary 公会房间列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body request_room.GuildRoomListreq  true "公会房间列表"
// @Success 0 {object} 	[]guild.GuildRoomListResp{}
// @Router /v1/guild/roomList [post]
func RoomList(c *gin.Context) {
	req := new(request_room.GuildRoomListreq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(loggic_guild.GuildRoom).RoomList(c, req))
}

// UpdateRoom
// @Summary 公会房间人员管理
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body request_room.ChangeRoomParamReq  true "参数"
// @Success 0 {object} 	response.Response{}
// @Router /v1/guild/updateRoom [post]
func UpdateRoom(c *gin.Context) {
	req := new(request_room.ChangeRoomParamReq)
	handle.BindBody(c, req)
	service := new(loggic_guild.GuildRoom).UpdateRoom(c, req)
	response.SuccessResponse(c, service)
}

// GetUserBaseInfoByUserNo
// @Summary 查询用户信息
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body request_room.UserNoParamReq  true "参数"
// @Success 0 {object} 	guild.GuildUserBaseInfo{}
// @Router /v1/guild/getUserBaseInfo [post]
func GetUserBaseInfoByUserNo(c *gin.Context) {
	req := new(request_room.UserNoParamReq)
	handle.BindBody(c, req)
	service := new(loggic_guild.GuildRoom).GetUserInfoByUserNo(c, req)
	response.SuccessResponse(c, service)
}

// CloseRoom
// @Summary 公会房间关闭
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body request_room.CloseRoomReq  true "参数"
// @Success 0 {object} 	response.Response{}
// @Router /v1/guild/closeRoom [post]
func CloseRoom(c *gin.Context) {
	req := new(request_room.CloseRoomReq)
	handle.BindBody(c, req)
	service := new(loggic_guild.GuildRoom).CloseRoom(c, req)
	response.SuccessResponse(c, service)
}

// RoomApply
// @Summary 公会房间申请
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body request_room.GuildRoomApplyReq{}  true "房间申请参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/guildRoomApply [post]
func RoomApply(c *gin.Context) {
	req := new(request_room.GuildRoomApplyReq)
	handle.BindBody(c, req)
	service := new(loggic_guild.GuildRoom).RoomApply(c, req)
	response.SuccessResponse(c, service)
}

// RoomApplyList
// @Summary 公会房间申请列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body request_room.GuildRoomApplyListReq{}  true "房间申请列表参数"
// @Success 0 {object}  []guild.GuildRoomApplyListResp{}
// @Router /v1/guild/guildRoomApplyList [post]
func RoomApplyList(c *gin.Context) {
	req := new(request_room.GuildRoomApplyListReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(loggic_guild.GuildRoom).RoomApplyList(c, req))
}
