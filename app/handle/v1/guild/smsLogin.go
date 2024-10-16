package guild

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/logic"
	"yfapi/internal/logic/guild"
	typedef_enum "yfapi/typedef/enum"
	request_login "yfapi/typedef/request/guild"
	"yfapi/typedef/response"
	response_login "yfapi/typedef/response/guild"
)

// SmsCodeLogin
// @Summary 根据手机验证码登录
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body request_login.LoginMobileReq	true	"参数"
// @Success 0 {object} response_login.LoginMobileCodeRes
// @Router /v1/guild/loginAction [post]
func SmsCodeLogin(context *gin.Context) {
	req := new(request_login.LoginMobileReq)
	handle.BindBody(context, req)
	handle.RepeatSubmitPost(context)
	var res response_login.LoginMobileCodeRes
	res = new(guild.GuildLogin).LoginByCode(req, context)
	response.SuccessResponse(context, res)
}

// GetLoginMobileSMS
// @Summary 获取手机验证码
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body request_login.SendMobileCodeReq	true	"参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/sendSms [post]
func GetLoginMobileSMS(context *gin.Context) {
	req := new(request_login.SendMobileCodeReq)
	handle.BindBody(context, req)
	if req.Type == 0 {
		req.Type = typedef_enum.SmsCodeGuildAdminLogin
	}
	switch req.Type {
	case typedef_enum.SmsCodeGuildAdminLogin, typedef_enum.SmsCodeBindBankCard:
	default:
		panic(i18n_err.ErrorCodeParam)
	}
	//通过手机查询用户
	_, _ = new(guild.GuildLogin).SearchByMobile(context, req)
	service := logic.Sms{
		Mobile:     req.Mobile,
		Type:       req.Type,
		RegionCode: req.RegionCode,
	}
	service.SendSms(context)
	response.SuccessResponse(context, "短信发送成功！")
}

// GuildInfo
// @Summary 公会信息
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Success 0 {object} response_login.GuildInfo
// @Router /v1/guild/guildInfo [post]
func GuildInfo(context *gin.Context) {
	service := new(guild.GuildLogin)
	resp := service.GuildInfo(context)
	response.SuccessResponse(context, resp)
}
