package roomOwner

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/logic"
	"yfapi/internal/logic/room"
	typedef_enum "yfapi/typedef/enum"
	request_login "yfapi/typedef/request/roomOwner"
	"yfapi/typedef/response"
	response_login "yfapi/typedef/response/roomOwner"
)

// SmsCodeLogin
// @Summary 手机验证码登录
// @Description 根据验证码登录
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_login.LoginMobileReq	true	"参数"
// @Success 0 {object} response_login.LoginMobileCodeRes{}
// @Router /v1/roomOwner/login [post]
func SmsCodeLogin(context *gin.Context) {
	req := new(request_login.LoginMobileReq)
	handle.BindBody(context, req)
	handle.RepeatSubmitPost(context)
	service := new(room.RoomLogin)
	response.SuccessResponse(context, service.LoginByCode(req, context))
}

// GetLoginMobileSMS
// @Summary 获取手机验证码
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body request_login.SendMobileCodeReq	true	"参数"
// @Success 0 {object} response.Response{}
// @Router /v1/roomOwner/sendSms [post]
func GetLoginMobileSMS(context *gin.Context) {
	req := new(request_login.SendMobileCodeReq)
	handle.BindBody(context, req)
	//通过手机查询用户
	if req.Type == 0 {
		req.Type = typedef_enum.SmsCodeRoomAdminLogin
	}
	switch req.Type {
	case typedef_enum.SmsCodeRoomAdminLogin, typedef_enum.SmsCodeBindBankCard:
	default:
		panic(i18n_err.ErrorCodeParam)
	}
	_ = new(room.RoomLogin).SearchByMobile(req.Mobile, context)
	service := logic.Sms{
		Mobile:     req.Mobile,
		Type:       req.Type,
		RegionCode: req.RegionCode,
	}
	service.SendSms(context)
	response.SuccessResponse(context, "")
}

// RoomListInfo
// @Summary 登录页房间列表
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Success 0 {object} []response_login.RoomInfo{}
// @Router /v1/roomOwner/room/list [post]
func RoomListInfo(context *gin.Context) {
	service := new(room.RoomLogin)
	var resp []*response_login.RoomInfo
	resp = service.RoomList(context)
	response.SuccessResponse(context, resp)
}
