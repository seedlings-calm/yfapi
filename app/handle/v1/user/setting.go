package user

import (
	"yfapi/app/handle"
	"yfapi/internal/logic"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
)

// SetPassword
//
//	@Description: 设置密码
func SetPassword(c *gin.Context) {
	req := new(request_user.SetPasswordReq)
	handle.BindBody(c, req)
	new(logic.UserSetting).SetPassword(req, c)
	response.SuccessResponse(c, "")
}

// 实名认证
func RealName(c *gin.Context) {
	req := new(request_user.UserRealNameReq)
	handle.BindBody(c, req)
	new(logic.UserSetting).RealName(c, req)
	response.SuccessResponse(c, "")
}

// 实名认证状态
func RealNameStatus(c *gin.Context) {
	response.SuccessResponse(c, new(logic.UserSetting).RealNameStatus(c))
}

// 获取用户关联的账号
func GetUserAccounts(c *gin.Context) {
	response.SuccessResponse(c, new(logic.UserSetting).GetAccountByMobile(c))
}

// 创建账号检测
func CheckCreateAccount(c *gin.Context) {
	response.SuccessResponse(c, new(logic.UserSetting).CheckCreateNewAccount(c))
}

// 创建新账号
func CreateNewAccount(c *gin.Context) {
	req := new(request_user.UserCreateAccountReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(logic.UserSetting).CreateNewAccount(c, req))
}

// 选择账号
func SwitchAccount(c *gin.Context) {
	req := new(request_user.SwitchUserAccountReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(logic.UserSetting).SwitchAccount(c, req))
}

// 获取用户隐私资料
func GetUserPrivateInfo(c *gin.Context) {
	response.SuccessResponse(c, new(logic.UserSetting).GetUserPrivateInfo(c))
}

// 发送校验手机号验证码
func SendVerifyCode(c *gin.Context) {
	new(logic.UserSetting).SendVerifyMobileCode(c)
	response.SuccessResponse(c, true)
}

// 校验旧手机号
func VerifyMobile(c *gin.Context) {
	req := new(request_user.VerifyMobileReq)
	handle.BindBody(c, req)
	new(logic.UserSetting).VerifyMobile(c, req)
	response.SuccessResponse(c, true)
}

// 绑定新手机号
func ChangeMobile(c *gin.Context) {
	req := new(request_user.ChangeMobileReq)
	handle.BindBody(c, req)
	new(logic.UserSetting).ChangeMobile(c, req)
	response.SuccessResponse(c, true)
}

// @Summary 获取登录设备日志
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} []user.LoginRecordResponse{}
// @Router /v1/user/loginRecord [get]
func GetLoginLog(c *gin.Context) {
	response.SuccessResponse(c, new(logic.UserSetting).GetLoginLog(c))
}

// 获取用户隐私设置信息
func GetPrivacySetting(c *gin.Context) {
	response.SuccessResponse(c, new(logic.UserPrivacy).GetPrivacySetting(c))
}
