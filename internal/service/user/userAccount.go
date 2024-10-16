package user

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	response_user "yfapi/typedef/response/user"
	"yfapi/util/easy"
)

// 查询用户账户信息
func GetUserAccountInfo(userId string) (result response_user.UserAccountDTO) {
	userAccountDao := new(dao.UserAccountDao)
	accountInfo, err := userAccountDao.GetUserAccountByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	subsidyAmount, err := userAccountDao.GetUserAccountSubsidyTotalCountByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	result.SubsidyList, err = userAccountDao.GetUserAccountSubsidyListByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}

	// 钻石余额
	diamond := easy.StringToDecimal(accountInfo.DiamondAmount)
	// 不可提现星光
	starlightUW := easy.StringToDecimal(accountInfo.StarlightAmount)
	// 可提现星光
	starlightW := easy.StringToDecimal(accountInfo.CanWithdrawAmount)
	// 补贴星光
	subsidy := easy.StringToDecimal(subsidyAmount)

	result.UserId = userId
	result.Status = accountInfo.Status
	result.WithdrawStatus = accountInfo.WithdrawStatus
	result.DiamondAmount = easy.StringFixed(diamond)
	result.StarlightAmount = easy.StringFixed(starlightUW.Add(starlightW).Add(subsidy))
	result.StarlightUnWithdraw = easy.StringFixed(starlightUW)
	result.StarlightWithdraw = easy.StringFixed(starlightW)
	result.StarlightSubsidy = easy.StringFixed(subsidy)
	result.Version = accountInfo.Version
	return
}

// UpdateAccountParam 用户账户变动参数
// 更新用户账户包括 钻石增减，可提现星光、补贴星光、不可提现星光增减
// 不同账户的变动，传入相应参数，必须对应各币种，星光数据库分为了三个币种，但用户看到的是总星光，所以扣除星光需特殊处理
// 钻石增减 - accountBook.CURRENCY_DIAMOND
// 可提现星光增加-accountBook.CURRENCY_STARLIGHT_WITHDRAW、
// 补贴星光增加-accountBook.CURRENCY_STARLIGHT_SUBSIDY、
// 不可提现星光增加-accountBook.CURRENCY_STARLIGHT_UNWITHDRAW
// 星光减少，扣除顺序 不可提现星光>补贴星光>可提现星光 币种统一为accountBook.CURRENCY_STARLIGHT
// 变动账户必须传入MySQL事务
// 一个事务中有多个账户变动，后续变动的账户中必须传入AccountVersion
type UpdateAccountParam struct {
	Tx                *gorm.DB // mysql事务
	UserId            string   // 用户ID
	FromUserId        string   // 打赏用户ID 星光收益赋值
	ToUserIdList      string   // 被打赏用户ID列表 打赏钻石扣除赋值
	Gid               string   // 关联物品id、礼物id
	Num               int      // (打赏)数量
	Diamond           int      // (打赏)钻石
	Currency          string   // 币种
	FundFlow          int      // 资金方向 1入 2出
	Amount            string   // 变动数量
	OrderId           string   // 关联订单ID
	OrderType         int      // 订单类型
	RoomId            string   // 关联房间ID
	GuildId           string   // 关联公会ID
	Note              string   // 备注信息
	SubsidyType       int      // 补贴类型 1房间日补贴 2房间月补贴 4公会月补贴
	SubsidyAmountType int      // 补贴账户类型 1房间 2公会
	AccountVersion    int      // 账户版本 一个事务中有多次账户变动使用
}

// UpdateUserAccount 更新用户账户信息
func UpdateUserAccount(data *UpdateAccountParam) int {
	if data.Tx == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUnknown,
			Msg:  nil,
		})
	}
	switch data.Currency {
	case accountBook.CURRENCY_DIAMOND: // 钻石
		return updateAccountDiamond(data)
	case accountBook.CURRENCY_STARLIGHT_UNWITHDRAW: // 不可提现星光
		if data.FundFlow == accountBook.FUND_INFLOW {
			return addAccountStarlight(data)
		}
	case accountBook.CURRENCY_STARLIGHT_WITHDRAW: // 可提现星光
		if data.FundFlow == accountBook.FUND_INFLOW {
			return addAccountStarlight(data)
		}
	case accountBook.CURRENCY_STARLIGHT_SUBSIDY: // 补贴星光
		if data.FundFlow == accountBook.FUND_INFLOW {
			return addAccountStarlight(data)
		} else if data.FundFlow == accountBook.FUND_OUTFLOW {
			return deductAccountStarlightSub(data)
		}
	case accountBook.CURRENCY_STARLIGHT: // 星光 扣除
		if data.FundFlow == accountBook.FUND_OUTFLOW {
			return deductAccountStarlight(data)
		}
	}
	return 0
}

// 变动用户钻石
func updateAccountDiamond(data *UpdateAccountParam) (accountVersion int) {
	accountInfo, err := new(dao.UserAccountDao).GetUserAccountByUserId(data.UserId)
	if err != nil {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 变动数量
	beforeAmount := accountInfo.DiamondAmount
	amount := data.Amount
	diamondDecimal := easy.StringToDecimal(data.Amount)
	if data.FundFlow == accountBook.FUND_INFLOW {
		accountInfo.DiamondAmount = easy.StringToDecimal(accountInfo.DiamondAmount).Add(diamondDecimal).String()
	} else {
		accountInfo.DiamondAmount = easy.StringToDecimal(accountInfo.DiamondAmount).Sub(diamondDecimal).String()
		amount = "-" + data.Amount
	}
	if data.AccountVersion > 0 {
		accountInfo.Version = data.AccountVersion
	}
	accountInfo.Version++
	// 更新钻石数量
	affectedCount := data.Tx.Model(model.UserAccount{}).Where("user_id=? and version<?", data.UserId, accountInfo.Version).Update("diamond_amount", gorm.Expr("diamond_amount+?", amount)).
		Update("version", accountInfo.Version).RowsAffected
	if affectedCount == 0 {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	if len(data.RoomId) == 0 {
		data.RoomId = "0"
	}
	if len(data.GuildId) == 0 {
		data.GuildId = "0"
	}
	if len(data.FromUserId) == 0 {
		data.FromUserId = "0"
	}
	// 记录流水
	orderBillInfo := &model.OrderBill{
		OrderId:      data.OrderId,
		UserId:       data.UserId,
		FromUserId:   data.FromUserId,
		ToUserIdList: data.ToUserIdList,
		Gid:          data.Gid,
		Num:          data.Num,
		RoomId:       data.RoomId,
		GuildId:      data.GuildId,
		Currency:     accountBook.CURRENCY_DIAMOND,
		FundFlow:     data.FundFlow,
		BeforeAmount: beforeAmount,
		Amount:       diamondDecimal.String(),
		CurrAmount:   accountInfo.DiamondAmount,
		OrderType:    data.OrderType,
		Note:         data.Note,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = data.Tx.Create(orderBillInfo).Error
	if err != nil {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	accountVersion = accountInfo.Version
	return
}

// 增加用户星光
func addAccountStarlight(data *UpdateAccountParam) (accountVersion int) {
	if len(data.RoomId) == 0 {
		data.RoomId = "0"
	}
	if len(data.GuildId) == 0 {
		data.GuildId = "0"
	}
	if len(data.FromUserId) == 0 {
		data.FromUserId = "0"
	}
	accountInfo, err := new(dao.UserAccountDao).GetUserAccountByUserId(data.UserId)
	if err != nil {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 变动数量
	beforeAmount := ""
	amount := data.Amount
	currAmount := ""
	amountDecimal := easy.StringToDecimal(data.Amount)
	changeColumn := ""
	switch data.Currency {
	case accountBook.CURRENCY_STARLIGHT_UNWITHDRAW: // 不可提现星光
		beforeAmount = accountInfo.StarlightAmount
		accountInfo.StarlightAmount = easy.StringToDecimal(accountInfo.StarlightAmount).Add(amountDecimal).String()
		currAmount = accountInfo.StarlightAmount
		changeColumn = "starlight_amount"
	case accountBook.CURRENCY_STARLIGHT_WITHDRAW: // 可提现星光
		beforeAmount = accountInfo.CanWithdrawAmount
		accountInfo.CanWithdrawAmount = easy.StringToDecimal(accountInfo.CanWithdrawAmount).Add(amountDecimal).String()
		currAmount = accountInfo.CanWithdrawAmount
		changeColumn = "can_withdraw_amount"
	case accountBook.CURRENCY_STARLIGHT_SUBSIDY: // 补贴星光
		var subsidyInfo model.UserAccountSubsidy
		if data.SubsidyType == 1 || data.SubsidyType == 2 { // 房间日月补贴
			if data.SubsidyAmountType == 2 { // 补贴人是公会会长
				subsidyInfo, err = new(dao.UserAccountDao).GetUserAccountGuildSubsidy(data.UserId, data.GuildId)
				if err != nil {
					data.Tx.Rollback()
					panic(error2.I18nError{
						Code: error2.ErrorCodeSystemBusy,
						Msg:  nil,
					})
				}
				if subsidyInfo.Id == 0 {
					subsidyInfo = model.UserAccountSubsidy{
						UserId:        data.UserId,
						AccountType:   2,
						RoomId:        "0",
						GuildId:       data.GuildId,
						Status:        1,
						SubsidyAmount: "0",
					}
				}
			} else {
				subsidyInfo, err = new(dao.UserAccountDao).GetUserAccountRoomSubsidy(data.UserId, data.RoomId)
				if err != nil {
					data.Tx.Rollback()
					panic(error2.I18nError{
						Code: error2.ErrorCodeSystemBusy,
						Msg:  nil,
					})
				}
				if subsidyInfo.Id == 0 {
					subsidyInfo = model.UserAccountSubsidy{
						UserId:        data.UserId,
						AccountType:   1,
						RoomId:        data.RoomId,
						GuildId:       "0",
						Status:        1,
						SubsidyAmount: "0",
					}
				}
			}
		} else {
			subsidyInfo, err = new(dao.UserAccountDao).GetUserAccountGuildSubsidy(data.UserId, data.GuildId)
			if err != nil {
				data.Tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeSystemBusy,
					Msg:  nil,
				})
			}
			if subsidyInfo.Id == 0 {
				subsidyInfo = model.UserAccountSubsidy{
					UserId:        data.UserId,
					AccountType:   2,
					RoomId:        "0",
					GuildId:       data.GuildId,
					Status:        1,
					SubsidyAmount: "0",
				}
			}
		}
		// 星光变动
		beforeAmount = subsidyInfo.SubsidyAmount
		subsidyInfo.SubsidyAmount = easy.StringToDecimal(subsidyInfo.SubsidyAmount).Add(amountDecimal).String()
		currAmount = subsidyInfo.SubsidyAmount
		// 更新补贴星光
		err = data.Tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"subsidy_amount": subsidyInfo.SubsidyAmount,
			}),
		}).Create(&subsidyInfo).Error
		if err != nil {
			data.Tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
		// 记录流水
		orderBillInfo := &model.OrderBill{
			OrderId:      data.OrderId,
			UserId:       data.UserId,
			FromUserId:   data.FromUserId,
			RoomId:       data.RoomId,
			GuildId:      data.GuildId,
			Currency:     data.Currency,
			FundFlow:     data.FundFlow,
			BeforeAmount: beforeAmount,
			Amount:       amountDecimal.String(),
			CurrAmount:   currAmount,
			OrderType:    data.OrderType,
			Note:         data.Note,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err = data.Tx.Create(orderBillInfo).Error
		if err != nil {
			data.Tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
		return
	}

	// 更新可提现星光 不可提现星光逻辑
	if data.AccountVersion > 0 {
		accountInfo.Version = data.AccountVersion
	}
	accountInfo.Version++
	// 更新星光数量
	affectedCount := data.Tx.Model(model.UserAccount{}).Where("user_id=? and version<?", data.UserId, accountInfo.Version).Update(changeColumn, gorm.Expr(fmt.Sprintf("%v+?", changeColumn), amount)).
		Update("version", accountInfo.Version).RowsAffected
	if affectedCount == 0 {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 记录流水
	orderBillInfo := &model.OrderBill{
		OrderId:      data.OrderId,
		UserId:       data.UserId,
		FromUserId:   data.FromUserId,
		Gid:          data.Gid,
		Num:          data.Num,
		Diamond:      data.Diamond,
		RoomId:       data.RoomId,
		GuildId:      data.GuildId,
		Currency:     data.Currency,
		FundFlow:     data.FundFlow,
		BeforeAmount: beforeAmount,
		Amount:       amountDecimal.String(),
		CurrAmount:   currAmount,
		OrderType:    data.OrderType,
		Note:         data.Note,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = data.Tx.Create(orderBillInfo).Error
	if err != nil {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	accountVersion = accountInfo.Version
	return
}

// 扣除用户星光
func deductAccountStarlight(data *UpdateAccountParam) (accountVersion int) {
	if len(data.RoomId) == 0 {
		data.RoomId = "0"
	}
	if len(data.GuildId) == 0 {
		data.GuildId = "0"
	}
	if len(data.FromUserId) == 0 {
		data.FromUserId = "0"
	}
	accountInfo := GetUserAccountInfo(data.UserId)
	// 星光余额是否充足
	changeAmount := easy.StringToDecimal(data.Amount)
	if easy.StringToDecimal(accountInfo.StarlightAmount).LessThan(changeAmount) {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeStarlightExchangeNotEnough,
			Msg:  nil,
		})
	}
	if data.AccountVersion > 0 {
		accountInfo.Version = data.AccountVersion
	}
	// 扣除玩家星光 不可提现>补贴>可提现
	// 不可提现星光
	beforeA := accountInfo.StarlightUnWithdraw
	starA := easy.StringToDecimal(accountInfo.StarlightUnWithdraw)
	if starA.GreaterThan(decimal.Zero) { // 不可提现星光余额大于0
		// 扣除的星光数量
		currDeduct := "0"
		if starA.LessThan(changeAmount) { // 不可提现星光不足兑换
			// 剩余待扣除的星光
			changeAmount = changeAmount.Sub(starA)
			// 扣除所有不可提现星光
			currDeduct = accountInfo.StarlightUnWithdraw
			accountInfo.StarlightUnWithdraw = easy.StringFixed(decimal.Zero)
		} else { // 不可提现星光足够兑换
			// 剩余不可提现星光
			accountInfo.StarlightUnWithdraw = easy.StringFixed(starA.Sub(changeAmount))
			currDeduct = easy.StringFixed(changeAmount)
			changeAmount = decimal.Zero
		}
		// 更新星光数量
		accountInfo.Version++
		affectedCount := data.Tx.Model(model.UserAccount{}).Where("user_id=? and version<?", data.UserId, accountInfo.Version).Update("starlight_amount", gorm.Expr("starlight_amount-?", currDeduct)).
			Update("version", accountInfo.Version).RowsAffected
		if affectedCount == 0 {
			data.Tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
		// 记录流水
		orderBillInfo := &model.OrderBill{
			OrderId:      data.OrderId,
			UserId:       data.UserId,
			FromUserId:   data.FromUserId,
			Gid:          data.Gid,
			Num:          data.Num,
			RoomId:       data.RoomId,
			GuildId:      data.GuildId,
			Currency:     accountBook.CURRENCY_STARLIGHT_UNWITHDRAW,
			FundFlow:     data.FundFlow,
			BeforeAmount: beforeA,
			Amount:       currDeduct,
			CurrAmount:   accountInfo.StarlightUnWithdraw,
			OrderType:    data.OrderType,
			Note:         data.Note,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err := data.Tx.Create(orderBillInfo).Error
		if err != nil {
			data.Tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
	}

	// 补贴星光
	starB := easy.StringToDecimal(accountInfo.StarlightSubsidy)
	if changeAmount.GreaterThan(decimal.Zero) { // 兑换的星光还未扣完 使用补贴星光
		if starB.GreaterThan(decimal.Zero) {
			// 扣除补贴星光
			for _, info := range accountInfo.SubsidyList {
				beforeM := info.SubsidyAmount
				// 扣除的星光数量
				currDeduct := info.SubsidyAmount
				starM := easy.StringToDecimal(info.SubsidyAmount)
				if starM.LessThan(changeAmount) {
					currDeduct = info.SubsidyAmount
					changeAmount = changeAmount.Sub(starM)
				} else {
					info.SubsidyAmount = easy.StringFixed(starM.Sub(changeAmount))
					currDeduct = easy.StringFixed(changeAmount)
					changeAmount = decimal.Zero
				}
				// 更新补贴星光
				err := data.Tx.Save(info).Error
				if err != nil {
					data.Tx.Rollback()
					panic(error2.I18nError{
						Code: error2.ErrorCodeSystemBusy,
						Msg:  nil,
					})
				}
				// 记录流水
				orderBillInfo := &model.OrderBill{
					OrderId:      data.OrderId,
					UserId:       data.UserId,
					FromUserId:   data.FromUserId,
					RoomId:       data.RoomId,
					GuildId:      data.GuildId,
					Currency:     accountBook.CURRENCY_STARLIGHT_SUBSIDY,
					FundFlow:     data.FundFlow,
					BeforeAmount: beforeM,
					Amount:       currDeduct,
					CurrAmount:   info.SubsidyAmount,
					OrderType:    data.OrderType,
					Note:         data.Note,
					CreateTime:   time.Now(),
					UpdateTime:   time.Now(),
				}
				err = data.Tx.Create(orderBillInfo).Error
				if err != nil {
					data.Tx.Rollback()
					panic(error2.I18nError{
						Code: error2.ErrorCodeSystemBusy,
						Msg:  nil,
					})
				}
				// 扣除完成 退出循环
				if changeAmount.IsZero() {
					break
				}
			}
		}
	}

	// 可提现星光
	beforeB := accountInfo.StarlightWithdraw
	starC := easy.StringToDecimal(accountInfo.StarlightWithdraw)
	if changeAmount.GreaterThan(decimal.Zero) { // 兑换的星光还未扣完 使用可提现星光
		if starC.GreaterThan(decimal.Zero) { // 可提现星光余额大于0
			// 扣除的星光数量
			currDeduct := "0"
			if starC.LessThan(changeAmount) { // 可提现星光不足兑换
				// 剩余待扣除的星光
				changeAmount = changeAmount.Sub(starC)
				// 扣除所有可提现星光
				currDeduct = accountInfo.StarlightWithdraw
				accountInfo.StarlightWithdraw = easy.StringFixed(decimal.Zero)
			} else { // 可提现星光足够兑换
				// 剩余可提现星光
				accountInfo.StarlightWithdraw = easy.StringFixed(starC.Sub(changeAmount))
				currDeduct = easy.StringFixed(changeAmount)
				changeAmount = decimal.Zero
			}
			// 更新星光数量
			accountInfo.Version++
			affectedCount := data.Tx.Model(model.UserAccount{}).Where("user_id=? and version<?", data.UserId, accountInfo.Version).Update("can_withdraw_amount", gorm.Expr("can_withdraw_amount-?", currDeduct)).
				Update("version", accountInfo.Version).RowsAffected
			if affectedCount == 0 {
				data.Tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeSystemBusy,
					Msg:  nil,
				})
			}
			// 记录流水
			orderBillInfo := &model.OrderBill{
				OrderId:      data.OrderId,
				UserId:       data.UserId,
				FromUserId:   data.FromUserId,
				Gid:          data.Gid,
				Num:          data.Num,
				RoomId:       data.RoomId,
				GuildId:      data.GuildId,
				Currency:     accountBook.CURRENCY_STARLIGHT_WITHDRAW,
				FundFlow:     data.FundFlow,
				BeforeAmount: beforeB,
				Amount:       currDeduct,
				CurrAmount:   accountInfo.StarlightUnWithdraw,
				OrderType:    data.OrderType,
				Note:         data.Note,
				CreateTime:   time.Now(),
				UpdateTime:   time.Now(),
			}
			err := data.Tx.Create(orderBillInfo).Error
			if err != nil {
				data.Tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeSystemBusy,
					Msg:  nil,
				})
			}
		}
	}
	if changeAmount.GreaterThan(decimal.Zero) { // 星光没扣完 回滚
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeStarlightExchangeNotEnough,
			Msg:  nil,
		})
	}
	accountVersion = accountInfo.Version
	return
}

// 扣除用户补贴星光 指定账户扣除，必传房间ID或公会ID
func deductAccountStarlightSub(data *UpdateAccountParam) (accountVersion int) {
	if len(data.RoomId) == 0 {
		data.RoomId = "0"
	}
	if len(data.GuildId) == 0 {
		data.GuildId = "0"
	}
	if len(data.FromUserId) == 0 {
		data.FromUserId = "0"
	}
	var subsidyInfo model.UserAccountSubsidy
	var err error
	if data.RoomId != "0" {
		subsidyInfo, err = new(dao.UserAccountDao).GetUserAccountRoomSubsidy(data.UserId, data.RoomId)
	} else if data.GuildId != "0" {
		subsidyInfo, err = new(dao.UserAccountDao).GetUserAccountGuildSubsidy(data.UserId, data.GuildId)
	}
	// 补贴星光余额是否充足
	changeAmount := easy.StringToDecimal(data.Amount)
	if easy.StringToDecimal(subsidyInfo.SubsidyAmount).LessThan(changeAmount) {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeStarlightExchangeNotEnough,
			Msg:  nil,
		})
	}
	// 扣除补贴星光
	beforeAmount := subsidyInfo.SubsidyAmount
	subsidyInfo.SubsidyAmount = easy.StringFixed(easy.StringToDecimal(subsidyInfo.SubsidyAmount).Sub(changeAmount))
	// 更新补贴星光
	err = data.Tx.Save(subsidyInfo).Error
	if err != nil {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 记录流水
	orderBillInfo := &model.OrderBill{
		OrderId:      data.OrderId,
		UserId:       data.UserId,
		FromUserId:   data.FromUserId,
		RoomId:       data.RoomId,
		GuildId:      data.GuildId,
		Currency:     accountBook.CURRENCY_STARLIGHT_SUBSIDY,
		FundFlow:     data.FundFlow,
		BeforeAmount: beforeAmount,
		Amount:       data.Amount,
		CurrAmount:   subsidyInfo.SubsidyAmount,
		OrderType:    data.OrderType,
		Note:         data.Note,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	err = data.Tx.Create(orderBillInfo).Error
	if err != nil {
		data.Tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	return
}
