package roomOwner

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/logic/room"
	resquest_room "yfapi/typedef/request/room"
	request_room "yfapi/typedef/request/roomOwner"
	request_users "yfapi/typedef/request/user"
	"yfapi/typedef/response"
	response_login "yfapi/typedef/response/roomOwner"
)

// RoomInfo
// @Summary 获取单个房间信息
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Success 0 {object} response_login.RoomHomeInfo{}
// @Router /v1/roomOwner/room/detail [post]
func RoomInfo(c *gin.Context) {
	var res response_login.RoomHomeInfo
	res = new(room.RoomLogin).RoomInfo(c)
	response.SuccessResponse(c, res)
}

// RoomBase
// @Summary 聊天室概况
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Success 0 {object} response_login.RoomHomeBaseInfo{}
// @Router /v1/roomOwner/room/base [post]
func RoomBase(c *gin.Context) {
	service := new(room.RoomLogin).RoomBase(c)
	response.SuccessResponse(c, service)
}

// RoomAdminList
// @Summary 管理员列表
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_room.RoomAdminListReq	true	"参数"
// @Success 0 {object} []response_login.RoomAdminInfo{}
// @Router /v1/roomOwner/admin/list [post]
func RoomAdminList(c *gin.Context) {
	req := new(request_room.RoomAdminListReq)
	handle.BindBody(c, req)
	service := new(room.RoomLogin).RoomAdminList(req, c)
	response.SuccessResponse(c, service)
}

// RoomAdminRemove
// @Summary 管理员移除
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_room.RoomCommonReq	true	"参数"
// @Success 0 {object} response.Response{}
// @Router /v1/roomOwner/admin/remove [post]
func RoomAdminRemove(c *gin.Context) {
	req := new(request_room.RoomCommonReq)
	handle.BindBody(c, req)
	service := new(room.RoomAdmin)
	service.RoomAdminRemove(c, req)
	response.SuccessResponse(c, nil)
}

// RoomAdminAdd
// @Summary 管理员添加
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body resquest_room.RoomAdminAddReq	true	"参数"
// @Success 0 {object} []response_login.RoomAdminInfo{}
// @Router /v1/roomOwner/admin/add [post]
func RoomAdminAdd(c *gin.Context) {
	req := new(resquest_room.RoomAdminAddReq)
	handle.BindBody(c, req)
	req.RoomId = c.GetString("roomId")
	service := new(room.RoomAdmin)
	service.RoomAdminAdd(c, req)
	response.SuccessResponse(c, nil)
}

// RoomPractitionerList
// @Summary 从业者列表
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_room.RoomPractitionerListReq	true	"参数"
// @Success 0 {object} []response_login.RoomPractitionerInfo{}
// @Router /v1/roomOwner/practitioner/list [post]
func RoomPractitionerList(c *gin.Context) {
	req := new(request_room.RoomPractitionerListReq)
	handle.BindBody(c, req)
	service := new(room.RoomLogin).RoomPractitionerList(req, c)
	response.SuccessResponse(c, service)
}

// RoomPractitionerAdd
// @Summary 从业者添加
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_room.RoomPractitionerAddReq	true	"参数"
// @Success 0 {object} response.Response{}
// @Router /v1/roomOwner/practitioner/add [post]
func RoomPractitionerAdd(c *gin.Context) {
	req := new(request_room.RoomPractitionerAddReq)
	handle.BindBody(c, req)
	service := new(room.RoomAdmin)
	service.RoomPractitionerAdd(c, req)
	response.SuccessResponse(c, nil)
}

// RoomPractitionerRemove
// @Summary 从业者移除
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_room.RoomPractitionerAddReq	true	"参数"
// @Success 0 {object} response.Response{}
// @Router /v1/roomOwner/practitioner/remove [post]
func RoomPractitionerRemove(c *gin.Context) {
	req := new(request_room.RoomPractitionerAddReq)
	handle.BindBody(c, req)
	service := new(room.RoomAdmin)
	service.RoomPractitionerRemove(c, req)
	response.SuccessResponse(c, nil)
}

// RoomPractitionerReSave
// @Summary 从业者申请再次提交
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_room.RoomPractitionerUpdateReq	true	"参数"
// @Success 0 {object} response.Response{}
// @Router /v1/roomOwner/practitioner/reSave [post]
func RoomPractitionerReSave(c *gin.Context) {
	req := new(request_room.RoomPractitionerUpdateReq)
	handle.BindBody(c, req)
	service := new(room.RoomAdmin)
	service.RoomPractitionerReSave(c, req)
	response.SuccessResponse(c, nil)
}

// RoomPractitionerInvalid
// @Summary 从业者申请作废
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_room.RoomPractitionerUpdateReq	true	"参数"
// @Success 0 {object} response.Response{}
// @Router /v1/roomOwner/practitioner/invalid [post]
func RoomPractitionerInvalid(c *gin.Context) {
	req := new(request_room.RoomPractitionerUpdateReq)
	handle.BindBody(c, req)
	service := new(room.RoomAdmin)
	service.RoomPractitionerInvalid(c, req)
	response.SuccessResponse(c, nil)
}

// BindBank
// @Summary 房主绑定银行卡
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  req body request_room.RoomBankBindReq  true "房主绑定银行卡参数"
// @Success 0 {object} response.Response{}
// @Router /v1/roomOwner/bindBank [post]
func BindBank(c *gin.Context) {
	req := new(request_room.RoomBankBindReq)
	handle.BindBody(c, req)
	service := new(room.RoomBank)
	code := service.BindBank(c, req)
	if code != i18n_err.SuccessCode {
		panic(i18n_err.I18nError{
			Code: code,
			Msg:  nil,
		})
	}
	response.SuccessResponse(c, "")
}

// RoomWithdrawApply
// @Summary 房主提现
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  req body request_users.UserWithdrawApplyReq  true "房主提现参数"
// @Success 0 {object} response.Response{}
// @Router /v1/roomOwner/withdrawApply [post]
func RoomWithdrawApply(c *gin.Context) {
	req := new(request_users.UserWithdrawApplyReq)
	handle.BindBody(c, req)
	new(room.RoomBank).RoomWithdrawApply(c, req)
	response.SuccessResponse(c, "")
}
