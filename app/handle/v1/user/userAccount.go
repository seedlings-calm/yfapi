package user

import (
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	"yfapi/internal/logic"
	request_orderBill "yfapi/typedef/request/orderBill"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
	_ "yfapi/typedef/response/orderBill"
	response_user "yfapi/typedef/response/user"

	"github.com/gin-gonic/gin"
)

// @Summary 查询用户账户信息
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Success 0 {object} response_user.UserAccountRes{}
// @Router /v1/user/account [get]
func GetUserAccountInfo(c *gin.Context) {
	var res response_user.UserAccountRes
	res = new(logic.UserAccount).GetUserAccountInfo(c)
	response.SuccessResponse(c, res)
}

// @Summary 用户钻石流水
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param  req	query	request_orderBill.DiamondBillReq	true	"钻石流水列表参数"
// @Success 0 {object} response.BasePageRes{}
// @Success 0 {object} []orderBill.DiamondBill{}
// @Router /v1/user/bill/diamond [get]
func GetUserDiamondBill(c *gin.Context) {
	req := new(request_orderBill.DiamondBillReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic.OrderBill).GetUserDiamondBillList(c, req))
}

// @Summary 用户星光流水
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param  req	query	request_orderBill.DiamondBillReq	true	"钻石流水列表参数"
// @Success 0 {object} response.BasePageRes{}
// @Success 0 {object} []orderBill.StarlightBill{}
// @Router /v1/user/bill/starlight [get]
func GetUserStarlightBill(c *gin.Context) {
	req := new(request_orderBill.DiamondBillReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic.OrderBill).GetUserStarlightBillList(c, req))
}

// @Summary 用户星光兑换钻石
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param  req body request_user.ExchangeDiamondReq   true "兑换钻石参数"
// @Success 0 {object} response_user.UserAccountRes{}
// @Router /v1/user/exchangeDiamond [post]
func StarlightExchangeDiamond(c *gin.Context) {
	req := new(request_user.ExchangeDiamondReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic.UserAccount).StarlightExchangeDiamond(c, req))
}

// @Summary 充值钻石页面
// @Description
// @Tags 充值相关
// @Accept json
// @Produce json
// @Param  userNo query int  true "充值用户ID"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_user.RechargeDiamondRes{}
// @Router /v1/user/rechargeDiamond [get]
func RechargeDiamond(c *gin.Context) {
	logicUser := new(logic.UserAccount)
	userNo := c.Query("userNo")
	res := logicUser.RechargeDiamond(c, userNo, "", "")
	response.SuccessResponse(c, res)
}

// @Summary 充值钻石页面-官网使用
// @Description
// @Tags 充值相关
// @Accept json
// @Produce json
// @Param  userNo query int  true "充值用户ID"
// @Param  platform query string  true "平台"
// @Param  channel query string  true "渠道"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_user.RechargeDiamondRes{}
// @Router /v1/website/rechargeDiamondOffical [get]
func RechargeDiamondForOffical(c *gin.Context) {
	logicUser := new(logic.UserAccount)
	userNo := c.Query("userNo")
	platform := c.Query("platform")
	channel := c.Query("channel")
	if userNo == "" || platform == "" || channel == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	res := logicUser.RechargeDiamond(c, userNo, platform, channel)
	response.SuccessResponse(c, res)
}

// @Summary 充值钻石明细
// @Description
// @Tags 充值相关
// @Accept json
// @Produce json
// @Param  req query  request_orderBill.RechargeDiamondReq true "充值钻石日志"
// @Success 0 {object} response.BasePageRes{}
// @Success 0 {object} []orderBill.DiamondBill{}
// @Router /v1/user/bill/rechargeDiamondLog [get]
func RechargeDiamondLog(c *gin.Context) {
	req := new(request_orderBill.RechargeDiamondReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic.OrderBill).GetUserRechargeDiamonLogList(c, req))
}
