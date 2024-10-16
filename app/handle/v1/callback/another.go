package v1_callback

import (
	"encoding/json"
	"net/http"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreLog"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	"yfapi/internal/service/pay"
	"yfapi/typedef/enum"
	"yfapi/typedef/response/orderBill"

	"github.com/gin-gonic/gin"
)

// 代付回调
func AnotherCallback(c *gin.Context) {
	resp := new(pay.AnotherPayNotifyResp)
	handle.BindBody(c, resp)
	if !new(pay.AnotherPay).VerifySignature(resp) {
		coreLog.Error("代付签名验证错误")
		return
	}
	//orderRecord := new(dao.OrderDao).FindOne(&model.Order{
	//	OrderId:     resp.OrderId,
	//	OrderStatus: accountBook.ORDER_STATUS_UNCOMPLETION,
	//})
	orderRecord := new(dao.OrderDao).FindOneMap(map[string]any{
		"order_id":     resp.OrderId,
		"order_status": accountBook.ORDER_STATUS_UNCOMPLETION,
	})
	if orderRecord.ID == 0 { //不存在
		coreLog.Error("AnotherCallback :订单不存在 order_id:%s", resp.OrderId)
		return
	}
	withdrawStatus := orderRecord.WithdrawStatus
	var notes orderBill.WithdrawNote
	err := json.Unmarshal([]byte(orderRecord.Note), &notes)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUnknown,
			Msg:  nil,
		})
	}
	//订单状态 0待处理，1处理中，2完成，3失败
	if resp.Status == "2" {
		orderRecord.OrderStatus = accountBook.ORDER_STATUS_COMPLETION
		orderRecord.PayStatus = accountBook.PAY_STATUS_COMPLETION
		orderRecord.WithdrawStatus = 3
	} else if resp.Status == "3" {
		orderRecord.OrderStatus = accountBook.ORDER_STATUS_UNCOMPLETION
		orderRecord.PayStatus = accountBook.PAY_STATUS_UNPAID
		orderRecord.WithdrawStatus = 4
		//
	} else {
		return
	}
	orderRecord.OrderNo = resp.OutTradeNo
	err = new(dao.OrderDao).Save(orderRecord)
	if err == nil {
		c.JSON(http.StatusOK, "")
	} else {
		coreLog.Error("AnotherCallback 代付错误 order_id:%s", resp.OrderId)
	}
	// 记录操作记录
	subsidyActionRecord := model.SubsidyActionRecord{
		OrderID:      orderRecord.OrderId,
		Action:       enum.WithdrawOrderStatus(orderRecord.WithdrawStatus).String(),
		BeforeStatus: enum.WithdrawOrderStatus(withdrawStatus).String(),
		CurrStatus:   enum.WithdrawOrderStatus(orderRecord.WithdrawStatus).String(),
		Memo:         "",
		StaffName:    notes.StaffName,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = new(dao.OrderDao).Add(&subsidyActionRecord)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
}
