package recharge

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	"yfapi/internal/logic"
	request_recharge "yfapi/typedef/request/recharge"
	"yfapi/typedef/response"
)

// @Summary 测试充值钻石
// @Description
// @Tags 充值相关
// @Accept json
// @Produce json
// @Param  req body request_recharge.UserRechargeTestReq   true "充值参数"
// @Success 0 {object} response.Response{}
// @Router /v1/recharge/test [post]
func UserRechargeTest(c *gin.Context) {
	req := new(request_recharge.UserRechargeTestReq)
	handle.BindBody(c, req)
	new(logic.Recharge).UserRechargeTest(c, req)
	response.SuccessResponse(c, true)
}

// 苹果充值
func IosIap(c *gin.Context) {
	req := new(request_recharge.IosIapReq)
	handle.BindBody(c, req)
	new(logic.Recharge).IosIap(c, req)
	response.SuccessResponse(c, true)
}

// 微信app支付
func WxAppPay(c *gin.Context) {
	req := new(request_recharge.WxAppPayReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(logic.Recharge).WxAppPay(c, req.ProductId))
}

// 支付宝支付
func AliAppPay(c *gin.Context) {
	req := new(request_recharge.WxAppPayReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(logic.Recharge).AliAppPay(c, req.ProductId))
}

// 第三方聚合支付
func AggregationPay(c *gin.Context) {
	req := new(request_recharge.AggregationPayReq)
	handle.BindBody(c, req)
	res := new(logic.Recharge).AggregationPay(c, req)
	response.SuccessResponse(c, res)
}

// 官网支付
func WebsitePay(c *gin.Context) {
	req := new(request_recharge.WebsitePayReq)
	handle.BindBody(c, req)
	res := new(logic.Recharge).WebsiteAggregationPay(c, req)
	response.SuccessResponse(c, res)
}

// 支付结果查询
func RechargeResult(c *gin.Context) {
	req := new(request_recharge.RechargeResultReq)
	handle.BindBody(c, req)
	res := new(logic.Recharge).GetRechargeResult(c, req)
	response.SuccessResponse(c, res)
}
