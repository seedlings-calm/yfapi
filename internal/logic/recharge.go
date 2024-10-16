package logic

import (
	"errors"
	"fmt"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	"yfapi/internal/service/pay"
	request_recharge "yfapi/typedef/request/recharge"
	pay2 "yfapi/typedef/response/pay"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type Recharge struct {
}

// 苹果内购支付
func (r *Recharge) IosIap(c *gin.Context, req *request_recharge.IosIapReq) {
	userId := helper.GetUserId(c)
	applePay := pay.NewApplePay()
	err := applePay.VerifyReceipt(userId, req.Receipt, "")
	if err != nil {
		coreLog.Error("苹果支付失败:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeAppleIapError,
		})
	}
	return
}

// 聚合支付
func (r *Recharge) AggregationPay(c *gin.Context, req *request_recharge.AggregationPayReq) (resp pay2.AggregationPayResp) {
	userId := helper.GetUserId(c)
	productInfo := new(dao.ConfigDiamondDao).FindOneByKeys(req.ProductId)
	if productInfo.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeProductNotFound,
		})
	}
	amount := cast.ToString(productInfo.Nums)
	orderModel, err := pay.CreateOrder(userId, amount, amount, accountBook.PAY_TYPE_AGGREGATION, req.ProductId)
	if err != nil {
		panic(err)
	}
	newAmount := easy.StringToDecimalFixed(amount)
	paySer := pay.NewAggregationPay()
	res, err := paySer.Pay(&pay.AggregationPayReq{
		MerchantCode: coreConfig.GetHotConf().AggregationPay.MerchantCode,
		BankCode:     req.Payment,
		Currency:     accountBook.CURRENCY_CNY,
		Amount:       newAmount,
		OrderId:      orderModel.OrderId,
		OrderDate:    cast.ToString(time.Now().UnixMilli()),
		Ip:           c.ClientIP(),
		GoodsName:    i18n_msg.GetI18nMsg(c, i18n_msg.Diamond),
		GoodsDetail:  i18n_msg.GetI18nMsg(c, i18n_msg.RechargeDiamond),
		UserId:       userId,
	})
	coreLog.LogInfo("聚合支付:%+v", res)
	if err != nil {
		panic(err)
	}
	if !res.Success {
		panic(res.ResultMsg)
	}

	resp.Types = res.Data.Data.Type
	resp.Info = res.Data.Data.Info
	resp.OrderId = res.Data.OrderId
	return
}

// 微信app原生支付
func (r *Recharge) WxAppPay(c *gin.Context, productId string) any {
	//userId := helper.GetUserId(c)
	v3, err := pay.NewWechatV3()
	if err != nil {
		coreLog.Error("微信支付失败:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeWxPayError,
		})
	}
	wxpayConf := coreConfig.GetHotConf().WxPay
	app, err := v3.V3TransactionApp(pay.V3TransactionAppReq{
		Mchid:       wxpayConf.Mchid,
		OutTradeNo:  new(accountBook.Order).OrderNum(accountBook.ORDER_CZ),
		Appid:       wxpayConf.Appid,
		Description: "充值钻石",
		NotifyUrl:   "https://api.sdwsweb.com/api",
		Amount: struct {
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		}{
			Total:    1,
			Currency: accountBook.CURRENCY_CNY,
		},
	})
	if err != nil {
		coreLog.Error("微信支付失败:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeWxPayError,
		})
	}
	return app
}

// 支付宝app原生支付
func (r *Recharge) AliAppPay(c *gin.Context, productId string) any {
	//userId := helper.GetUserId(c)
	v3, err := pay.NewAliV3()
	if err != nil {
		coreLog.Error("支付宝支付失败:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeAliPayError,
		})
	}
	rsp, err := v3.TradeAppPay(pay.AliAppPayReq{
		OutTradeNo:  new(accountBook.Order).OrderNum(accountBook.ORDER_CZ),
		TotalAmount: "0.01",
		Subject:     "测试支付",
	})
	if err != nil {
		coreLog.Error("支付宝支付失败:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeAliPayError,
		})
	}
	return rsp
}

// UserRechargeTest 用户测试充值
func (r *Recharge) UserRechargeTest(c *gin.Context, req *request_recharge.UserRechargeTestReq) {
	if req.Diamond > 100000 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	diamondDecimal := decimal.NewFromInt(req.Diamond)
	userId := helper.GetUserId(c)
	accountInfo, err := new(dao.UserAccountDao).GetUserAccountByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 增加钻石数量
	beforeAmount := accountInfo.DiamondAmount
	accountInfo.DiamondAmount = easy.StringToDecimal(accountInfo.DiamondAmount).Add(decimal.NewFromInt(req.Diamond)).String()
	accountInfo.Version++
	tx := coreDb.GetMasterDb().Begin()
	// 生成订单
	orderInfo := &model.Order{
		OrderId:         new(accountBook.Order).OrderNum(accountBook.ORDER_CZ),
		UserId:          userId,
		RoomId:          "0",
		GuildId:         "0",
		TotalAmount:     diamondDecimal.String(),
		PayAmount:       diamondDecimal.Div(decimal.NewFromInt(10)).String(),
		DiscountsAmount: "0",
		Num:             1,
		Currency:        accountBook.CURRENCY_CNY,
		OrderType:       accountBook.ChangeDiamondRecharge,
		OrderStatus:     1,
		PayStatus:       1,
		Note:            "测试充值",
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		StatDate:        time.Now().Format(time.DateOnly),
	}
	err = tx.Create(orderInfo).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 记录流水
	orderBillInfo := &model.OrderBill{
		OrderId:      orderInfo.OrderId,
		UserId:       userId,
		FromUserId:   "0",
		RoomId:       "0",
		GuildId:      "0",
		Currency:     accountBook.CURRENCY_DIAMOND,
		FundFlow:     accountBook.FUND_INFLOW,
		BeforeAmount: beforeAmount,
		Amount:       diamondDecimal.String(),
		CurrAmount:   accountInfo.DiamondAmount,
		OrderType:    accountBook.ChangeDiamondRecharge,
		Note:         "测试充值",
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = tx.Create(orderBillInfo).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	affectedCount := tx.Model(model.UserAccount{}).Where("user_id=? and version<?", userId, accountInfo.Version).Update("diamond_amount", gorm.Expr("diamond_amount+?", req.Diamond)).
		Update("version", accountInfo.Version).RowsAffected
	if affectedCount == 0 {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	tx.Commit()
	return
}

// 代付
func (r *Recharge) AnotherPay(c *gin.Context, req *request_recharge.AnotherPayReq) (bool, error) {
	paySer := pay.NewAnotherPay()
	//orderRecord := new(dao.OrderDao).FindOne(&model.Order{
	//	OrderId:        req.OrderId,
	//	UserId:         req.UserId,
	//	OrderType:      accountBook.ChangeStarlightStarlightWithdrawal,
	//	OrderStatus:    accountBook.ORDER_STATUS_UNCOMPLETION,
	//	WithdrawStatus: 3,
	//})
	orderRecord := new(dao.OrderDao).FindOneMap(map[string]any{
		"order_id":        req.OrderId,
		"user_id":         req.UserId,
		"order_type":      accountBook.ChangeStarlightStarlightWithdrawal,
		"order_status":    accountBook.ORDER_STATUS_UNCOMPLETION,
		"withdraw_status": 2,
	})
	if orderRecord.ID == 0 { //不存在
		msg := fmt.Sprintf("订单不存在 order_id:%s", req.OrderId)
		coreLog.Info(msg)
		return false, errors.New(msg)
	}
	amount := easy.StringToDecimalFixed(req.Amount)
	res, err := paySer.Pay(&pay.AnotherPayReq{
		MerchantCode: coreConfig.GetHotConf().AnotherPay.MerchantCode,
		OrderId:      req.OrderId,
		BankCardNum:  req.BankCardNum,
		BankCardName: req.BankCardName,
		Branch:       req.Branch,
		BankCode:     req.BankCode,
		Amount:       amount,
		NotifyUrl:    coreConfig.GetHotConf().AnotherPay.NotifyUrl,
		OrderDate:    cast.ToString(time.Now().UnixMilli()),
		Currency:     accountBook.CURRENCY_CNY,
	})
	if err != nil {
		return false, err
	}
	if !res.Success {
		errMsg, _ := res.ResultMsg.(string)
		return false, errors.New(errMsg)
	}
	return true, nil
}

// 官网聚合支付
func (r *Recharge) WebsiteAggregationPay(c *gin.Context, req *request_recharge.WebsitePayReq) (resp pay2.AggregationPayResp) {
	userId := req.UserId
	userModel, _ := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if len(userModel.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
		})
	}
	productInfo := new(dao.ConfigDiamondDao).FindOneByKeys(req.ProductId)
	if productInfo.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeProductNotFound,
		})
	}
	amount := cast.ToString(productInfo.Nums)
	orderModel, err := pay.CreateOrder(userId, amount, amount, accountBook.PAY_TYPE_AGGREGATION, req.ProductId)
	if err != nil {
		panic(err)
	}
	newAmount := easy.StringToDecimalFixed(amount)
	paySer := pay.NewAggregationPay()
	res, err := paySer.Pay(&pay.AggregationPayReq{
		MerchantCode: coreConfig.GetHotConf().AggregationPay.MerchantCode,
		BankCode:     req.Payment,
		Currency:     accountBook.CURRENCY_CNY,
		Amount:       newAmount,
		OrderId:      orderModel.OrderId,
		OrderDate:    cast.ToString(time.Now().UnixMilli()),
		Ip:           c.ClientIP(),
		GoodsName:    i18n_msg.GetI18nMsg(c, i18n_msg.Diamond),
		GoodsDetail:  i18n_msg.GetI18nMsg(c, i18n_msg.RechargeDiamond),
		UserId:       userId,
	})
	if err != nil {
		panic(err)
	}
	if !res.Success {
		panic(res.ResultMsg)
	}
	orderModel.OrderNo = res.Data.OutTradeNo
	err = coreDb.GetMasterDb().Save(orderModel).Error
	if err != nil {
		panic(err)
	}
	resp.Types = res.Data.Data.Type
	resp.Info = res.Data.Data.Info
	resp.OrderId = res.Data.OrderId
	return
}

// 支付结果查询接口
func (r *Recharge) GetRechargeResult(c *gin.Context, req *request_recharge.RechargeResultReq) bool {
	userId := req.UserId
	userModel, _ := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if len(userModel.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
		})
	}
	order := new(dao.OrderDao).FindOneMap(map[string]any{"order_id": req.OrderId, "user_id": userId, "order_status": accountBook.ORDER_STATUS_COMPLETION, "pay_status": accountBook.PAY_STATUS_COMPLETION, "order_type": accountBook.ChangeDiamondRecharge})
	if order.ID > 0 {
		return true
	}
	return false
}
