package user

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	i18n_err "yfapi/i18n/error"
	logic_user "yfapi/internal/logic"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
	response_user "yfapi/typedef/response/user"
)

// @Summary 获取提现详情
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_user.UserWithdrawRes{}
// @Router /v1/user/appSettingInfo [get]
func AppSettingInfo(context *gin.Context) {
	service := new(logic_user.AppSetting)
	res := response_user.UserWithdrawRes{}
	res = service.GetAppSettingInfo(context)
	response.SuccessResponse(context, res)
}

// @Summary 提现申请
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param  req body request_user.UserWithdrawApplyReq{}   true "提现参数"
// @Success 0 {object} response.Response{}
// @Router /v1/user/appSettingApply [post]
func AppSettingApply(context *gin.Context) {
	req := new(request_user.UserWithdrawApplyReq)
	handle.BindBody(context, req)
	service := new(logic_user.AppSetting)
	service.AppSettingApply(context, req)
	response.SuccessResponse(context, "")
}

// @Summary 绑定银行卡
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param  req body request_user.UserBankAddReq{}   true "绑定银行卡参数"
// @Success 0 {object} response.Response{}
// @Router /v1/user/bankAdd [post]
func BankAdd(context *gin.Context) {
	req := new(request_user.UserBankAddReq)
	handle.BindBody(context, req)
	code := new(logic_user.UserBank).Add(context, req)
	if code != i18n_err.SuccessCode {
		panic(i18n_err.I18nError{
			Code: code,
			Msg:  nil,
		})
	}
	response.SuccessResponse(context, "")
}

// @Summary 解绑银行卡
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param  req body request_user.UserBankUnBindReq{}   true "解绑银行卡参数"
// @Success 0 {object} response.Response{}
// @Router /v1/user/bankUnBind [post]
func BankUnBind(context *gin.Context) {
	req := new(request_user.UserBankUnBindReq)
	handle.BindBody(context, req)
	service := new(logic_user.UserBank)
	service.UnBind(context, req)
	response.SuccessResponse(context, "")
}

// @Summary 获取银行卡列表
// @Description
// @Tags	用户相关
// @Accept	json
// @Produce	json
// @Success	0 {object} response.Response{}
// @Success 0 {object} response_user.UserBankInfo{}
// @Router	/v1/user/bankList [get]
func BankList(context *gin.Context) {
	service := new(logic_user.UserBank)
	response.SuccessResponse(context, service.GetBankList(context))
}
