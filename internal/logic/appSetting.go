package logic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strings"
	"time"
	"yfapi/core/coreDb"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	service_user "yfapi/internal/service/user"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response/orderBill"
	response_user "yfapi/typedef/response/user"
	"yfapi/util/easy"
)

type AppSetting struct {
}

// GetAppSettingInfo
//
//	@Description: 获取用户提现信息
//	@receiver ser
//	@param c
//	@return res
func (ser *AppSetting) GetAppSettingInfo(c *gin.Context) (res response_user.UserWithdrawRes) {
	userId := helper.GetUserId(c)
	// 获取用户信息
	userInfo := service_user.GetUserBaseInfo(userId)
	if len(userInfo.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	// 获取用户账号信息
	accountInfo := service_user.GetUserAccountInfo(userInfo.Id)
	if len(accountInfo.UserId) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	//获取银行卡信息
	bankDao := dao.UserBankDao{}
	bankInfo, err := bankDao.GetUserBankListByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(bankInfo) > 0 {
		for _, v := range bankInfo {
			res.BankList = append(res.BankList, &response_user.BankInfo{
				Id:         v.Id,
				BankName:   v.BankName,
				BankNo:     v.BankNo,
				BankHolder: v.BankHolder,
				IsDefault:  v.IsDefault,
			})
		}
	}
	//查询提现说明
	withdrawInfo, err := new(dao.UserWithdrawDao).UserWithdrawInfoById(1)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}

	res.StarlightAmount = accountInfo.StarlightAmount
	//星光总余额除以10
	res.StarlightAmount = easy.StringFixed(easy.StringToDecimal(accountInfo.StarlightAmount).Div(decimal.NewFromInt(10)))
	res.StarlightWithdraw = easy.StringFixed(easy.StringToDecimal(accountInfo.StarlightWithdraw).Div(decimal.NewFromInt(10)))
	res.StarlightUnWithdraw = easy.StringFixed(easy.StringToDecimal(accountInfo.StarlightUnWithdraw).Div(decimal.NewFromInt(10)))
	res.WithdrawRate = withdrawInfo.RewardRate
	res.StarlightSubsidy = easy.StringFixed(easy.StringToDecimal(accountInfo.StarlightSubsidy).Div(decimal.NewFromInt(10)))
	res.SubsidyRate = withdrawInfo.SettlementRate
	res.Desc = withdrawInfo.Desc
	res.WithdrawDesc = withdrawInfo.WithdrawDesc
	res.UnWithdrawDesc = withdrawInfo.UnWithdrawDesc
	// 将提现日期分割转成星期几返回
	daySlice := strings.Split(withdrawInfo.WithdrawDays, ",")
	// 创建一个空的字符串，用于存储转换后的星期几
	var weekDaysStr string
	// 遍历字符串切片，将每个字符串转换为星期几并添加到新的字符串中
	for _, day := range daySlice {
		// 将字符串转换为整数
		days := DayToString(c, day)
		weekDaysStr = weekDaysStr + days
	}
	if withdrawInfo.WithdrawDays == "1,2,3,4,5,6,0" {
		res.WithdrawDays = ""
	} else {
		res.WithdrawDays = weekDaysStr
	}
	return
}
func DayToString(c *gin.Context, day string) string {
	switch day {
	case "1":
		return i18n_msg.GetI18nMsg(c, i18n_msg.MondayKey)
	case "2":
		return i18n_msg.GetI18nMsg(c, i18n_msg.TuesdayKey)
	case "3":
		return i18n_msg.GetI18nMsg(c, i18n_msg.WednesdayKey)
	case "4":
		return i18n_msg.GetI18nMsg(c, i18n_msg.ThursdayKey)
	case "5":
		return i18n_msg.GetI18nMsg(c, i18n_msg.FridayKey)
	case "6":
		return i18n_msg.GetI18nMsg(c, i18n_msg.SaturdayKey)
	case "0":
		return i18n_msg.GetI18nMsg(c, i18n_msg.SundayKey)
	default:
		return i18n_msg.GetI18nMsg(c, i18n_msg.UnknownKey)
	}
}

// AppSettingApply
//
//	@Description: 用户提现申请
//	@receiver ser
//	@param c
//	@param req
func (ser *AppSetting) AppSettingApply(c *gin.Context, req *request_user.UserWithdrawApplyReq) {
	userId := helper.GetUserId(c)
	// 获取用户信息
	userInfo := service_user.GetUserBaseInfo(userId)
	if len(userInfo.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	//查询提现说明
	withdrawInfo, err := new(dao.UserWithdrawDao).UserWithdrawInfoById(1)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	todayInt := int(time.Now().Weekday())
	// 将字符串转换为切片
	daySlice := strings.Split(withdrawInfo.WithdrawDays, ",")
	// 判断今天是否是提现日
	if !easy.InArray(cast.ToString(todayInt), daySlice) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserWithdrawDay,
			Msg:  nil,
		})
	}
	orderDao := dao.OrderDao{}
	//判断今日是否已提现过
	todayStart := easy.GetCurrDayStartTime(time.Now())
	todayEnd := easy.GetCurrDayEndTime(time.Now())
	withdraw, err := orderDao.IsUserWithdraw(userId, todayStart.Format(time.DateTime), todayEnd.Format(time.DateTime))
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if withdraw.ID > 0 {
		// 今日最新记录非审核拒绝状态不可再次提现
		if withdraw.WithdrawStatus != 1 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserWithdrawDayMaxErr,
				Msg:  nil,
			})
		}
	}
	// 获取用户账号信息
	accountInfo := service_user.GetUserAccountInfo(userInfo.Id)
	if len(accountInfo.UserId) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	bankDao := dao.UserBankDao{}
	//获取银行卡信息
	bankInfo, err := bankDao.GetUserBankById(req.BankId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if bankInfo.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserBankNotExist,
			Msg:  nil,
		})
	}
	starlightWithdraw := easy.StringToDecimal(accountInfo.StarlightWithdraw)
	starlightSubsidy := easy.StringToDecimal(accountInfo.StarlightSubsidy)
	amount := easy.StringToDecimal(cast.ToString(req.Amount))
	//可提现金额需要乘以10等于要扣除的星光数量
	subAmount := amount.Mul(decimal.NewFromInt(10))
	var sumStarlightAmount = starlightWithdraw.Add(starlightSubsidy)
	//1.提现金额不能大于可提现金额，且提现金额必须为100的倍数，如100、200、300，单笔提现金额上限为20000元人民币;
	if sumStarlightAmount.LessThan(subAmount) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserAmountExceed,
			Msg:  nil,
		})
	}
	if req.Amount < 100 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserAmountNotEnough,
			Msg:  nil,
		})
	}
	if req.Amount%100 != 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserAmountInvalid,
			Msg:  nil,
		})
	}
	if req.Amount > 20000 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserAmountMaxErr,
			Msg:  nil,
		})
	}
	var payAmount decimal.Decimal
	var fee decimal.Decimal
	rewardRate := easy.StringToDecimal(cast.ToString(withdrawInfo.RewardRate))
	settlementRate := easy.StringToDecimal(cast.ToString(withdrawInfo.SettlementRate))
	tx := coreDb.GetMasterDb().Begin()
	orderId := new(accountBook.Order).OrderNum(accountBook.ORDER_TX)
	// 用户提现金额优先扣除打赏收入再扣除结算收入
	// 用户提现金额小于等于可提现金额，直接扣除打赏收入
	Note := orderBill.WithdrawNote{
		BankName:   bankInfo.BankName,
		BankNo:     bankInfo.BankNo,
		BankHolder: bankInfo.BankHolder,
		BankCode:   bankInfo.BankCode,
		BankBranch: bankInfo.BankBranch,
		StaffName:  "",
		Reason:     "",
	}
	NoteMarshal, err := json.Marshal(Note)
	if subAmount.LessThanOrEqual(starlightWithdraw) {
		// 计算手续费，使用 Decimal 类型的乘法和除法
		fee = amount.Mul(rewardRate.Div(easy.StringToDecimal(cast.ToString(100))))
		// 提现实际金额
		payAmount = amount.Sub(fee)
		//扣除星光余额
		CurrAmount := easy.StringToDecimal(accountInfo.StarlightWithdraw).Sub(subAmount).String()
		accountInfo.Version++
		affectedCount := tx.Model(model.UserAccount{}).Where("user_id=? and version<?", userId, accountInfo.Version).Update("can_withdraw_amount", gorm.Expr("can_withdraw_amount-?", subAmount)).
			Update("version", accountInfo.Version).RowsAffected
		if affectedCount == 0 {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
		// 记录流水
		orderBillInfo := &model.OrderBill{
			OrderId:      orderId,
			UserId:       userId,
			FromUserId:   "0",
			ToUserIdList: "",
			Gid:          "0",
			Num:          0,
			RoomId:       "0",
			GuildId:      "0",
			Currency:     accountBook.CURRENCY_STARLIGHT_WITHDRAW,
			FundFlow:     2,
			BeforeAmount: starlightWithdraw.String(),
			Amount:       subAmount.String(),
			CurrAmount:   CurrAmount,
			AppId:        "",
			OrderType:    accountBook.ChangeStarlightStarlightWithdrawal,
			Note:         accountBook.EnumOrderType(accountBook.ChangeStarlightStarlightWithdrawal).String(),
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
	} else {
		// 用户提现金额大于可提现金额，先扣除打赏收入，再扣除结算收入
		// 打赏收入手续费，确保使用 Decimal 类型进行计算，首先starlightWithdraw要除以10
		withdrawFee := starlightWithdraw.Div(easy.StringToDecimal(cast.ToString(10))).Mul(rewardRate.Div(easy.StringToDecimal(cast.ToString(100))))
		// 结算收入手续费，确保使用 Decimal 类型进行计算
		subsidyFee := amount.Sub(starlightWithdraw.Div(easy.StringToDecimal(cast.ToString(10)))).Mul(settlementRate.Div(easy.StringToDecimal(cast.ToString(100))))
		fee = withdrawFee.Add(subsidyFee)
		// 提现实际金额
		payAmount = amount.Sub(fee)
		//扣除星光余额
		// 更新星光数量
		accountInfo.Version++
		affectedCount := tx.Model(model.UserAccount{}).Where("user_id=? and version<?", userId, accountInfo.Version).Update("can_withdraw_amount", gorm.Expr("can_withdraw_amount-?", starlightWithdraw)).
			Update("version", accountInfo.Version).RowsAffected
		if affectedCount == 0 {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
		// 记录可提现星光扣除流水
		orderBillInfo := &model.OrderBill{
			OrderId:      orderId,
			UserId:       userId,
			FromUserId:   "0",
			ToUserIdList: "",
			Gid:          "0",
			Num:          0,
			RoomId:       "0",
			GuildId:      "0",
			Currency:     accountBook.CURRENCY_STARLIGHT_WITHDRAW,
			FundFlow:     2,
			BeforeAmount: starlightWithdraw.String(),
			Amount:       starlightWithdraw.String(),
			CurrAmount:   "0",
			AppId:        "",
			OrderType:    accountBook.ChangeStarlightStarlightWithdrawal,
			Note:         accountBook.EnumOrderType(accountBook.ChangeStarlightStarlightWithdrawal).String(),
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
		//扣除完星光，剩余数量用补贴星光扣除
		remainAmount := subAmount.Sub(starlightWithdraw)
		// 扣除补贴星光
		for _, info := range accountInfo.SubsidyList {
			beforeM := info.SubsidyAmount
			// 扣除的星光数量
			currDeduct := info.SubsidyAmount
			starM := easy.StringToDecimal(info.SubsidyAmount)
			if starM.LessThan(remainAmount) {
				currDeduct = info.SubsidyAmount
				remainAmount = remainAmount.Sub(starM)
			} else {
				info.SubsidyAmount = easy.StringFixed(starM.Sub(remainAmount))
				currDeduct = easy.StringFixed(remainAmount)
				remainAmount = decimal.Zero
			}
			// 更新补贴星光
			err = tx.Model(model.UserAccountSubsidy{}).Where("id", info.Id).Updates(info).Error
			if err != nil {
				tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeSystemBusy,
					Msg:  nil,
				})
			}
			// 记录补贴星光扣除流水
			orderBillInfos := &model.OrderBill{
				OrderId:      orderId,
				UserId:       info.UserId,
				FromUserId:   "0",
				ToUserIdList: "",
				Gid:          "0",
				Num:          0,
				RoomId:       "0",
				GuildId:      "0",
				Currency:     accountBook.CURRENCY_STARLIGHT_SUBSIDY,
				FundFlow:     2,
				BeforeAmount: beforeM,
				Amount:       currDeduct,
				CurrAmount:   info.SubsidyAmount,
				AppId:        "",
				OrderType:    accountBook.ChangeStarlightStarlightWithdrawal,
				Note:         accountBook.EnumOrderType(accountBook.ChangeStarlightStarlightWithdrawal).String(),
				CreateTime:   time.Now(),
				UpdateTime:   time.Now(),
			}

			err = tx.Create(orderBillInfos).Error
			if err != nil {
				tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeSystemBusy,
					Msg:  nil,
				})
			}
			// 扣除完成 退出循环
			if remainAmount.IsZero() {
				break
			}
		}
	}
	//生成提现订单
	order := &model.Order{
		OrderId:         orderId,
		UserId:          userId,
		ToUserIdList:    "",
		RoomId:          "0",
		GuildId:         "0",
		Gid:             "0",
		TotalAmount:     easy.StringFixed(amount),
		PayAmount:       easy.StringFixed(payAmount),
		DiscountsAmount: "0",
		Num:             0,
		Currency:        accountBook.CURRENCY_STARLIGHT,
		AppId:           "",
		OrderType:       accountBook.ChangeStarlightStarlightWithdrawal,
		OrderStatus:     0,
		PayType:         1,
		PayStatus:       0,
		WithdrawStatus:  0,
		OrderNo:         "",
		Note:            string(NoteMarshal),
		StatDate:        time.Now().Format(time.DateOnly),
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	err = tx.Create(order).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	bankInfos, err := bankDao.GetUserBankListByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	//修改该用户原有银行卡的状态都为0，除了该用户提现银行卡，提现银行卡状态为1
	if len(bankInfos) > 1 {
		for _, Info := range bankInfos {
			//将所有银行卡都改成0
			err = tx.Model(&model.UserBank{}).Where("user_id", Info.UserId).Update("is_default", 0).Error
			if err != nil {
				tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
			//将该用户提现银行卡改成1
			err = tx.Model(&model.UserBank{}).Where("id =?", req.BankId).Update("is_default", 1).Error
			if err != nil {
				tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
		}
	}

	tx.Commit()
}
