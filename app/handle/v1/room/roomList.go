package room

import (
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	"yfapi/internal/logic"
	"yfapi/typedef/enum"
	request_index "yfapi/typedef/request/index"
	request_room "yfapi/typedef/request/room"
	"yfapi/typedef/response"
	"yfapi/typedef/response/index"
	response_room "yfapi/typedef/response/room"

	"github.com/gin-gonic/gin"
)

// list
//
//	@Summary	房间列表
//	@Schemes
//	@Description	房间列表
//	@Tags			房间相关
//	@Param			req	query	request_room.RoomListReq	true	"房间列表参数"
//	@Accept			json
//	@Produce		json
//	@Success		200 {object} response.BasePageRes
//	@Router			/v1/list [get]
func RoomList(c *gin.Context) {
	req := new(request_room.RoomListReq)
	handle.BindQuery(c, req)
	roomList := new(logic.Room)
	response.SuccessResponse(c, roomList.Page(c, req))
}

// @Summary 房间信息
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  roomId query string  true "房间ID"
// @Success 0 {object} response_room.ChatroomDTO{}
// @Router /v1/room/info [get]
func RoomInfo(c *gin.Context) {
	roomId := c.Query("roomId")
	if roomId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := new(logic.Room)
	var resp response_room.ChatroomDTO
	resp = service.FindOne(c, roomId)

	response.SuccessResponse(c, resp)
}

// @Summary 编辑房间
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param  req body request_room.RoomUpdateReq   true "编辑字段"
// @Success 0 {object} response.Response{}
// @Router /v1/room/update [post]
func RoomUpdate(c *gin.Context) {
	req := new(request_room.RoomUpdateReq)
	handle.BindBody(c, req)
	roomList := new(logic.Room)

	response.SuccessResponse(c, roomList.Update(c, req))
}

// @Summary 举报类型
// @Description
// @Tags 举报中心
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Router /v1/reportingTypes [get]
func ReportingCenterTypes(c *gin.Context) {
	t := enum.ReportingType
	tkeys := enum.ReportingTypeKey

	type value struct {
		Key   int    `json:"key"`
		Value string `json:"value"`
	}
	// 创建一个新的数组，数组长度等于 map 的长度
	tres := make([]value, len(t))

	// 填充数组，按照排序后的键顺序
	for i, k := range tkeys {
		var vals value
		vals.Key = k
		vals.Value = t[k]
		tres[i] = vals
	}

	response.SuccessResponse(c, tres)
}

// @Summary 执行举报房间
// @Description
// @Tags 举报中心
// @Accept json
// @Produce json
// @Param req  body request_room.ReportingCenterReq{}   true "举报参数"
// @Success 0 {object} response.Response{}
// @Router /v1/reporting [post]
func ReportingCenter(c *gin.Context) {
	req := new(request_room.ReportingCenterReq)
	handle.BindJson(c, req)
	service := &logic.Room{}
	service.ReportingCenter(c, req)
	response.SuccessResponse(c, nil)
}

// @Summary 申请创建房间
// @Description
// @Tags 房间相关
// @Accept json
// @Produce json
// @Param req  body request_room.ApplyAnchorRoomReq{}   true "举报参数"
// @Success 0 {object} response_room.ApplyAnchorRoomRes{}
// @Router /v1/room/apply [post]
func ApplyRoom(c *gin.Context) {
	req := new(request_room.ApplyAnchorRoomReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(logic.Room).ApplyRoom(c, req))
}

// 申请房间信息回显
func ApplyRoomInfo(c *gin.Context) {
	response.SuccessResponse(c, new(logic.Room).ApplyRoomInfo(c))
}

// @Summary	搜索用户聊天室直播间
// @Schemes
// @Description
// @Tags	房间相关
// @Param	req	query	request_index.SearchAllReq	true	"搜索参数"
// @Accept	json
// @Produce	json
// @Success	200 {object} index.SearchAllRes{}
// @Router	/v1/searchAll [get]
func SearchAll(c *gin.Context) {
	req := new(request_index.SearchAllReq)
	handle.BindQuery(c, req)
	res := index.SearchAllRes{}
	res = new(logic.Room).SearchAll(c, req)
	response.SuccessResponse(c, res)
}
