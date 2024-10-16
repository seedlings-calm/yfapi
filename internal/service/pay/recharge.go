package pay

import (
	"context"
	"fmt"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	"yfapi/util/easy"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// 苹果内购充值
func appleIapRecharge(userId, order_no, productId string) error {
	success, unlock, err := coreRedis.UserLock(context.Background(), redisKey.UserPayLock(userId), time.Second*10)
	if err != nil || !success {
		return fmt.Errorf("appleIapRecharge get lock err:%+v", err)
	}
	defer unlock()
	orderRecord := new(dao.OrderDao).FindOne(&model.Order{UserId: userId, PayType: accountBook.PAY_TYPE_APPLE_IAP, OrderNo: order_no})
	if orderRecord.ID != 0 { //存在
		return fmt.Errorf("交易已完成 order_no:%s", order_no)
	}
	productInfo := new(dao.ConfigDiamondDao).FindOneByKeys(productId)
	if productInfo.Id == 0 {
		return fmt.Errorf("商品不存在 productId:%s", productId)

	}
	// 根据productId获取钻石数量 金额
	var diamond int64 = cast.ToInt64(productInfo.GotoNums)
	var amount int64 = cast.ToInt64(productInfo.Nums)

	//获取官方账号得账户信息
	OperationAccount, err := new(dao.UserAccountDao).GetUserAccountByUserId(typedef_enum.OperationUserId)
	if err != nil {
		return fmt.Errorf("获取官方账户失败")
	}
	//获取充值用户的账户信息
	userAccount, err := new(dao.UserAccountDao).GetUserAccountByUserId(userId)
	if err != nil {
		return fmt.Errorf("获取用户账户错误")
	}
	diamondDecimal := decimal.NewFromInt(diamond)
	amountDecimal := decimal.NewFromInt(amount)
	// 增加钻石数量
	userBeforeAmount := userAccount.DiamondAmount
	userAccount.DiamondAmount = easy.StringToDecimal(userAccount.DiamondAmount).Add(diamondDecimal).String()
	tx := coreDb.GetMasterDb().Begin()
	// 生成订单
	orderInfo := &model.Order{
		OrderId:         new(accountBook.Order).OrderNum(accountBook.ORDER_CZ),
		UserId:          userId,
		RoomId:          "0",
		GuildId:         "0",
		TotalAmount:     amountDecimal.String(),
		PayAmount:       amountDecimal.String(),
		DiscountsAmount: "0",
		Num:             1,
		Currency:        accountBook.CURRENCY_CNY,
		OrderType:       accountBook.ChangeDiamondRecharge,
		OrderStatus:     accountBook.ORDER_STATUS_COMPLETION,
		PayType:         accountBook.PAY_TYPE_APPLE_IAP,
		PayStatus:       accountBook.PAY_STATUS_COMPLETION,
		Note:            accountBook.EnumOrderType(accountBook.ChangeDiamondRecharge).String(),
		OrderNo:         order_no,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		Gid:             productId,
		StatDate:        time.Now().Format(time.DateOnly),
	}
	err = tx.Create(orderInfo).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("订单生成失败: %+v", err)
	}
	// 记录运营账号流水
	orderBillInfo1 := &model.OrderBill{
		OrderId:      orderInfo.OrderId,
		UserId:       typedef_enum.OperationUserId,
		FromUserId:   "0",
		RoomId:       "0",
		GuildId:      "0",
		Currency:     accountBook.CURRENCY_DIAMOND,
		FundFlow:     accountBook.FUND_OUTFLOW,
		BeforeAmount: OperationAccount.DiamondAmount,
		Amount:       diamondDecimal.String(),
		CurrAmount:   easy.StringToDecimal(OperationAccount.DiamondAmount).Sub(diamondDecimal).String(),
		OrderType:    accountBook.ChangeDiamondRecharge,
		Note:         accountBook.EnumOrderType(accountBook.ChangeDiamondRecharge).String(),
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = tx.Create(orderBillInfo1).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("运营流水生成失败 err:%+v", err)
	}

	// 记录用户流水
	orderBillInfo2 := &model.OrderBill{
		OrderId:      orderInfo.OrderId,
		UserId:       userId,
		FromUserId:   "0",
		RoomId:       "0",
		GuildId:      "0",
		Currency:     accountBook.CURRENCY_DIAMOND,
		FundFlow:     accountBook.FUND_INFLOW,
		BeforeAmount: userBeforeAmount,
		Amount:       diamondDecimal.String(),
		CurrAmount:   userAccount.DiamondAmount,
		OrderType:    accountBook.ChangeDiamondRecharge,
		Note:         accountBook.EnumOrderType(accountBook.ChangeDiamondRecharge).String(),
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = tx.Create(orderBillInfo2).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("用户流水生成失败: %+v", err)
	}
	affectedCount := tx.Model(model.UserAccount{}).Where("user_id=?", userId).Updates(map[string]interface{}{
		"diamond_amount": gorm.Expr("diamond_amount+?", diamond),
		"version":        gorm.Expr("version+?", 1),
	}).RowsAffected
	if affectedCount == 0 {
		tx.Rollback()
		return fmt.Errorf("用户钻石增加失败")
	}
	affectedCount = tx.Model(model.UserAccount{}).Where("user_id=?", typedef_enum.OperationUserId).Update("diamond_amount", gorm.Expr("diamond_amount-?", diamond)).RowsAffected
	if affectedCount == 0 {
		tx.Rollback()
		return fmt.Errorf("运营钻石抽出失败")
	}
	tx.Commit()
	return nil
}

// 创建订单
func CreateOrder(userId string, totalAmount, payAmount string, payType int, productId string) (orderModel *model.Order, err error) {
	orderId := new(accountBook.Order).OrderNum(accountBook.ORDER_CZ)
	// 生成订单
	orderModel = &model.Order{
		OrderId:         orderId,
		UserId:          userId,
		RoomId:          "0",
		GuildId:         "0",
		TotalAmount:     totalAmount,
		PayAmount:       payAmount,
		DiscountsAmount: "0",
		Num:             1,
		Currency:        accountBook.CURRENCY_CNY,
		OrderType:       accountBook.ChangeDiamondRecharge,
		OrderStatus:     accountBook.ORDER_STATUS_UNCOMPLETION,
		PayType:         payType,
		PayStatus:       accountBook.PAY_STATUS_UNPAID,
		Note:            accountBook.EnumOrderType(accountBook.ChangeDiamondRecharge).String(),
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		Gid:             productId,
		StatDate:        time.Now().Format(time.DateOnly),
	}
	err = coreDb.GetMasterDb().Create(orderModel).Error
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	return
}

// 给用户充值
func Recharge(orderId, OutTradeNo string, payType int) error {
	//orderRecord := new(dao.OrderDao).FindOne(&model.Order{PayType: payType, OrderId: orderId, OrderNo: OutTradeNo, PayStatus: accountBook.PAY_STATUS_UNPAID})
	orderRecord := new(dao.OrderDao).FindOneMap(map[string]any{
		"pay_type":   payType,
		"order_id":   orderId,
		"order_no":   OutTradeNo,
		"pay_status": accountBook.PAY_STATUS_UNPAID,
	})
	if orderRecord.ID == 0 || orderRecord.PayStatus > 0 {
		return fmt.Errorf("订单不存在 order_id:%s", orderId)
	}
	productId := orderRecord.Gid
	userId := orderRecord.UserId
	success, unlock, err := coreRedis.UserLock(context.Background(), redisKey.UserPayLock(orderRecord.UserId), time.Second*10)
	if err != nil || !success {
		return fmt.Errorf("获取锁失败 userId:%s", userId)
	}
	defer unlock()
	productInfo := new(dao.ConfigDiamondDao).FindOneByKeys(productId)
	if productInfo.Id == 0 {
		return fmt.Errorf("商品不存在 productId:%s", productId)
	}
	var diamond int64 = cast.ToInt64(productInfo.GotoNums)
	//获取官方账号得账户信息
	OperationAccount, err := new(dao.UserAccountDao).GetUserAccountByUserId(typedef_enum.OperationUserId)
	if err != nil {
		return fmt.Errorf("获取官方账号信息失败：%+v", err)
	}
	//获取充值用户的账户信息
	userAccount, err := new(dao.UserAccountDao).GetUserAccountByUserId(userId)
	if err != nil {
		return fmt.Errorf("用户账号不存在 userId:%s", userId)
	}
	diamondDecimal := decimal.NewFromInt(diamond)
	// 增加钻石数量
	userBeforeAmount := userAccount.DiamondAmount
	userAccount.DiamondAmount = easy.StringToDecimal(userAccount.DiamondAmount).Add(diamondDecimal).String()
	tx := coreDb.GetMasterDb().Begin()
	// 记录运营账号流水
	orderBillInfo1 := &model.OrderBill{
		OrderId:      orderRecord.OrderId,
		UserId:       typedef_enum.OperationUserId,
		FromUserId:   "0",
		RoomId:       "0",
		GuildId:      "0",
		Currency:     accountBook.CURRENCY_DIAMOND,
		FundFlow:     accountBook.FUND_OUTFLOW,
		BeforeAmount: OperationAccount.DiamondAmount,
		Amount:       diamondDecimal.String(),
		CurrAmount:   easy.StringToDecimal(OperationAccount.DiamondAmount).Sub(diamondDecimal).String(),
		OrderType:    accountBook.ChangeDiamondRecharge,
		Note:         accountBook.EnumOrderType(accountBook.ChangeDiamondRecharge).String(),
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = tx.Create(orderBillInfo1).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("创建运营流水失败:%+v", err)
	}
	// 记录用户流水
	orderBillInfo2 := &model.OrderBill{
		OrderId:      orderRecord.OrderId,
		UserId:       userId,
		FromUserId:   "0",
		RoomId:       "0",
		GuildId:      "0",
		Currency:     accountBook.CURRENCY_DIAMOND,
		FundFlow:     accountBook.FUND_INFLOW,
		BeforeAmount: userBeforeAmount,
		Amount:       diamondDecimal.String(),
		CurrAmount:   userAccount.DiamondAmount,
		OrderType:    accountBook.ChangeDiamondRecharge,
		Note:         accountBook.EnumOrderType(accountBook.ChangeDiamondRecharge).String(),
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = tx.Create(orderBillInfo2).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("创建用户流水失败:%+v", err)
	}
	affectedCount := tx.Model(model.UserAccount{}).Where("user_id=?", userId).Updates(map[string]interface{}{
		"diamond_amount": gorm.Expr("diamond_amount+?", diamond),
		"version":        gorm.Expr("version+?", 1),
	}).RowsAffected
	if affectedCount == 0 {
		tx.Rollback()
		return fmt.Errorf("用户增加钻石失败")
	}
	affectedCount = tx.Model(model.UserAccount{}).Where("user_id=?", typedef_enum.OperationUserId).Update("diamond_amount", gorm.Expr("diamond_amount-?", diamond)).RowsAffected
	if affectedCount == 0 {
		tx.Rollback()
		return fmt.Errorf("运营账户扣除钻石失败")
	}
	orderRecord.PayStatus = accountBook.PAY_STATUS_COMPLETION
	orderRecord.OrderStatus = accountBook.ORDER_STATUS_COMPLETION
	err = tx.Save(orderRecord).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新订单状态失败", err)
	}
	tx.Commit()
	return nil
}

// 提现创建流水
func Withdraw(orderId string) error {
	orderRecord := new(dao.OrderDao).FindOne(&model.Order{
		OrderId: orderId,
	})
	if orderRecord.ID == 0 { //不存在
		coreLog.Error("订单不存在 order_id:%s", orderId)
		return nil
	}
	orderRecord.OrderStatus = accountBook.ORDER_STATUS_COMPLETION
	orderRecord.PayStatus = accountBook.PAY_STATUS_COMPLETION
	orderRecord.WithdrawStatus = 3
	err := new(dao.OrderDao).Save(orderRecord)
	return err
}
