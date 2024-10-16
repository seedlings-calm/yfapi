package login

import (
	"yfapi/app/handle"
	"yfapi/core/coreLog"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/logic"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
	request_login "yfapi/typedef/request/login"

	"yfapi/typedef/response"
	_ "yfapi/typedef/response/login"

	"github.com/gin-gonic/gin"
)

// 登陆前检测
func LoginCheck(c *gin.Context) {
	req := new(request_login.LoginCheckReq)
	handle.BindBody(c, req)
	service := new(logic.Login)
	service.CheckLoginType(c, req)
	response.SuccessResponse(c, true)
}

// LoginByCode
//
//	@Summary	根据验证码登录
//	@Schemes
//	@Description	根据验证码登录
//	@Tags			登录相关
//	@Param			req	body	request_login.LoginCodeReq	true	"登录参数"
//	@Accept			json
//	@Produce		json
//	@Success		200 {object} login.LoginCodeRes
//	@Router			/v1/loginByCode [post]
func LoginByCode(context *gin.Context) {
	req := new(request_login.LoginCodeReq)
	handle.BindBody(context, req)
	handle.RepeatSubmitPost(context)
	service := new(logic.Login)
	response.SuccessResponse(context, service.LoginByCode(req, context))
}

// @Summary 发送验证码
// @Description 发送验证码
// @Tags 登录模块
// @Accept application/json
// @Product application/json
// @Param req body request_login.SendSmsReq	true	"发送验证码参数"
// @Success 0 {object} response.Response{}
// @Router /v1/sendSms [post]
func SendSms(context *gin.Context) {
	req := new(request_login.SendSmsReq)
	handle.BindBody(context, req)
	if req.Type == enum.SmsCodeChangeMobile {
		userModel, _ := new(dao.UserDao).FindOne(&model.User{RegionCode: req.RegionCode, Mobile: req.Mobile})
		if len(userModel.Id) > 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeMobileUsed,
			})
		}
	}
	service := logic.Sms{
		Mobile:     req.Mobile,
		Type:       req.Type,
		RegionCode: req.RegionCode,
	}
	err := service.SendSms(context)
	if err != nil {
		coreLog.Info("sendSms", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeGetCaptcha,
			Msg:  nil,
		})
		return
	}
	response.SuccessResponse(context, "")
}

// LoginByPass
//
//	@Summary	手机号密码登录
//	@Schemes
//	@Description	手机号密码登录
//	@Tags			登录相关
//	@Param			req	body	request_login.LoginPassReq	true	"登录参数"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	login.LoginCodeRes
//	@Router			/v1/loginByPass [post]
func LoginByPass(context *gin.Context) {
	req := new(request_login.LoginPassReq)
	handle.BindBody(context, req)
	handle.RepeatSubmitPost(context)
	service := new(logic.Login)
	response.SuccessResponse(context, service.LoginByPass(req, context))
}

// GetChooseUser
//
//	@Summary	获取选择账号页面信息
//	@Schemes
//	@Description	获取选择账号页面信息
//	@Tags			登录相关
//	@Param			chooseUserToken	query string	true	"选择用户的token"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]login.GetChooseUserRes
//	@Router			/v1/getChooseUser [get]
func GetChooseUser(context *gin.Context) {
	chooseUserToken, _ := context.GetQuery("chooseUserToken")
	service := new(logic.Login)
	response.SuccessResponse(context, service.GetChooseUser(chooseUserToken, context))
}

// ChooseUserLogin
//
//	@Summary	选择账号登陆
//	@Schemes
//	@Description	选择账号登陆
//	@Tags			登录相关
//	@Param			req	body	request_login.ChooseUserLoginReq	true	"选择账号登陆参数"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	login.LoginCodeRes
//	@Router			/v1/chooseUserLogin [post]
func ChooseUserLogin(c *gin.Context) {
	req := new(request_login.ChooseUserLoginReq)
	handle.BindBody(c, req)
	handle.RepeatSubmitPost(c)
	service := new(logic.Login)
	response.SuccessResponse(c, service.ChooseUserLogin(req, c))
}

// ForgetPassCheck
//
//	@Summary	忘记密码检测
//	@Schemes
//	@Description	忘记密码检测
//	@Tags			登录相关
//	@Param			req	body	request_login.LoginCodeReq	true	"参数"
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	login.ForgetPassCheckRes
//	@Router			/v1/forgetPassCheck [post]
func ForgetPassCheck(c *gin.Context) {
	req := new(request_login.LoginCodeReq)
	handle.BindBody(c, req)
	service := new(logic.Login)
	response.SuccessResponse(c, service.ForgetPassCheck(req, c))
}

// ForgetPass
//
//	@Summary	忘记密码
//	@Schemes
//	@Description	忘记密码
//	@Tags			登录相关
//	@Param			req	body	request_login.ForgetPassReq	true	"参数"
//	@Accept			json
//	@Produce		json
//	@Router			/v1/forgetPass [post]
func ForgetPass(c *gin.Context) {
	req := new(request_login.ForgetPassReq)
	handle.BindBody(c, req)
	service := new(logic.Login)
	service.ForgetPass(req, c)
	response.SuccessResponse(c, "")
}
