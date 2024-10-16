package v1_callback

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yfapi/app/handle"
	"yfapi/core/coreLog"
	"yfapi/internal/service/accountBook"
	"yfapi/internal/service/pay"
)

// 第三方聚合支付回调
func AggregationCallback(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			coreLog.Error("AggregationCallback err:%+v", err)
		}
	}()
	resp := new(pay.AggregationPayNotifyResp)
	handle.BindBody(c, resp)
	coreLog.LogInfo("聚合支付回调 参数:%+v", resp)
	if !new(pay.AggregationPay).VerifySignature(resp) {
		coreLog.Error("聚合支付回调 验签失败")
		return
	}
	err := pay.Recharge(resp.OrderId, resp.OutTradeNo, accountBook.PAY_TYPE_AGGREGATION)
	if err != nil {
		coreLog.Error("聚合支付回调 支付钻石失败:%+v", err)
		return
	}
	c.JSON(http.StatusOK, "")
}
