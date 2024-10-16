package room

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"strings"
	"time"
	"yfapi/core/coreDb"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	service_logic "yfapi/internal/logic"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	service_bank "yfapi/internal/service/bank"
	service_user "yfapi/internal/service/user"
	typedef_enum "yfapi/typedef/enum"
	request_roomowner "yfapi/typedef/request/roomOwner"
	"yfapi/typedef/request/user"
	"yfapi/typedef/response/orderBill"
	"yfapi/util/easy"
)

type RoomBank struct{}

func (u *RoomBank) BindBank(c *gin.Context, req *request_roomowner.RoomBankBindReq) (code i18n_err.ErrCode) {
	userId := helper.GetUserId(c)
	userInfo := service_user.GetUserBaseInfo(userId)
	if req.BankHolder != userInfo.TrueName || userInfo.RealNameStatus != typedef_enum.UserRealNameAuthenticated {
		return i18n_err.ErrorCodeIDCardAuth
	}
	if req.Mobile != userInfo.Mobile {
		return i18n_err.ErrCodeMobileOrRegionCode
	}
	//	校验验证码
	sms := &service_logic.Sms{
		Mobile:     req.Mobile,
		Code:       req.Code,
		RegionCode: req.RegionCode,
		Type:       typedef_enum.SmsCodeBindBankCard,
	}
	err := sms.CheckSms(c)
	if err != nil {
		return i18n_err.ErrorCodeCaptchaInvalid
	}
	code = service_bank.Bank{}.BankBind(userId, req.BankNo, req.BankName, req.BankHolder, req.BankBranch)
	return code
}

// 房主提现
func (u *RoomBank) RoomWithdrawApply(c *gin.Context, req *user.UserWithdrawApplyReq) {
	userId := helper.GetUserId(c)
	roomId := helper.GetRoomId(c)
	if req.Amount < 100 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserAmountNotEnough,
			Msg:  nil,
		})
	}
	if req.Amount%100 != 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserAmountInvalid,
			Msg:  nil,
		})
	}
	if req.Amount > 20000 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserAmountMaxErr,
			Msg:  nil,
		})
	}
	amount := easy.StringToDecimal(cast.ToString(req.Amount))
	//提现金额乘以10等于要扣除的补贴星光
	subAmount := amount.Mul(easy.StringToDecimal(cast.ToString(10)))
	//查询房主补贴星光信息
	subsidyInfo, err := new(dao.UserAccountDao).GetUserAccountRoomSubsidy(userId, roomId)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if subsidyInfo.Id == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeSubsidyInfoNotExist,
			Msg:  nil,
		})
	}
	subsidyAmount := easy.StringToDecimal(subsidyInfo.SubsidyAmount)
	if subsidyAmount.LessThan(subAmount) {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserAmountExceed,
			Msg:  nil,
		})
	}
	//查询提现说明
	withdrawInfo, err := new(dao.UserWithdrawDao).UserWithdrawInfoById(1)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	todayInt := int(time.Now().Weekday())
	// 将字符串转换为切片
	daySlice := strings.Split(withdrawInfo.WithdrawDays, ",")
	// 判断今天是否是提现日
	if !easy.InArray(cast.ToString(todayInt), daySlice) {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserWithdrawDay,
			Msg:  nil,
		})
	}

	orderDao := dao.OrderDao{}
	//判断今日是否已提现过
	todayStart := easy.GetCurrDayStartTime(time.Now())
	todayEnd := easy.GetCurrDayEndTime(time.Now())
	withdraw, err := orderDao.IsUserWithdraw(userId, todayStart.Format(time.DateTime), todayEnd.Format(time.DateTime))
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if withdraw.ID > 0 {
		// 今日最新记录非审核拒绝状态不可再次提现
		//if withdraw.WithdrawStatus != 1 {
		//	panic(i18n_err.I18nError{
		//		Code: i18n_err.ErrorCodeUserWithdrawDayMaxErr,
		//		Msg:  nil,
		//	})
		//}
	}
	var payAmount decimal.Decimal
	var fee decimal.Decimal
	settlementRate := easy.StringToDecimal(cast.ToString(withdrawInfo.SettlementRate))
	// 提现手续费
	fee = amount.Mul(settlementRate.Div(easy.StringToDecimal(cast.ToString(100))))
	// 提现实际金额
	payAmount = amount.Sub(fee)
	tx := coreDb.GetMasterDb().Begin()
	orderId := new(accountBook.Order).OrderNum(accountBook.ORDER_TX)
	// 扣除房主补贴星光
	service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:        tx,
		UserId:    userId,
		RoomId:    roomId,
		Num:       1,
		Currency:  accountBook.CURRENCY_STARLIGHT_SUBSIDY,
		FundFlow:  accountBook.FUND_OUTFLOW,
		Amount:    cast.ToString(subAmount),
		OrderId:   orderId,
		OrderType: accountBook.ChangeStarlightStarlightWithdrawal,
		Note:      accountBook.EnumOrderType(accountBook.ChangeStarlightStarlightWithdrawal).String(),
	})
	bankDao := dao.UserBankDao{}
	//获取银行卡信息
	bankInfo, err := bankDao.GetUserBankById(req.BankId)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if bankInfo.Id == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserBankNotExist,
			Msg:  nil,
		})
	}
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
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 创建提现订单
	orderInfo := &model.Order{
		OrderId:         orderId,
		UserId:          userId,
		ToUserIdList:    "",
		RoomId:          roomId,
		GuildId:         "0",
		Gid:             "",
		TotalAmount:     easy.StringFixed(amount),
		PayAmount:       easy.StringFixed(payAmount),
		DiscountsAmount: "0",
		Num:             0,
		Currency:        accountBook.CURRENCY_STARLIGHT_SUBSIDY,
		AppId:           "",
		OrderType:       accountBook.ChangeStarlightStarlightWithdrawal,
		OrderStatus:     0,
		PayType:         1,
		PayStatus:       0,
		WithdrawStatus:  0,
		OrderNo:         "",
		Note:            string(NoteMarshal),
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		StatDate:        time.Now().Format(time.DateOnly),
	}
	bankInfos, err := bankDao.GetUserBankListByUserId(userId)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
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
				panic(i18n_err.I18nError{
					Code: i18n_err.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
			//将该用户提现银行卡改成1
			err = tx.Model(&model.UserBank{}).Where("id =?", req.BankId).Update("is_default", 1).Error
			if err != nil {
				tx.Rollback()
				panic(i18n_err.I18nError{
					Code: i18n_err.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
		}
	}

	err = tx.Create(orderInfo).Error
	if err != nil {
		tx.Rollback()
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	tx.Commit()
	return
}
