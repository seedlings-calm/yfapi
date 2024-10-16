package logic

import (
	"log"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	request_inner "yfapi/typedef/request/inner"
	"yfapi/typedef/request/user"
	response_user "yfapi/typedef/response/user"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type UserAccount struct {
}

// GetUserAccountInfo 查询用户账户信息
func (u *UserAccount) GetUserAccountInfo(c *gin.Context) (res response_user.UserAccountRes) {
	userId := helper.GetUserId(c)
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
	// 钻石余额
	diamond := easy.StringToDecimal(accountInfo.DiamondAmount)
	// 不可提现星光
	starlightUW := easy.StringToDecimal(accountInfo.StarlightAmount)
	// 可提现星光
	starlightW := easy.StringToDecimal(accountInfo.CanWithdrawAmount)
	// 补贴星光
	subsidy := easy.StringToDecimal(subsidyAmount)

	res = response_user.UserAccountRes{
		UserId:              userId,
		Status:              accountInfo.Status,
		WithdrawStatus:      accountInfo.WithdrawStatus,
		DiamondAmount:       easy.StringFixed(diamond),
		StarlightAmount:     easy.StringFixed(starlightUW.Add(starlightW).Add(subsidy)),
		StarlightUnWithdraw: easy.StringFixed(starlightUW),
		StarlightWithdraw:   easy.StringFixed(starlightW),
		StarlightSubsidy:    easy.StringFixed(subsidy),
	}
	return
}

// StarlightExchangeDiamond 星光兑换钻石
func (u *UserAccount) StarlightExchangeDiamond(c *gin.Context, req *user.ExchangeDiamondReq) (res response_user.UserAccountRes) {
	userId := helper.GetUserId(c)
	accountInfo := service_user.GetUserAccountInfo(userId)
	// 星光余额是否充足
	exchangeAmount := decimal.NewFromInt(req.ExchangeAmount)
	if easy.StringToDecimal(accountInfo.StarlightAmount).LessThan(exchangeAmount) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeStarlightExchangeNotEnough,
			Msg:  nil,
		})
	}

	// 创建兑换订单
	tx := coreDb.GetMasterDb().Begin()
	orderInfo := &model.Order{
		ID:              0,
		OrderId:         new(accountBook.Order).OrderNum(accountBook.ORDER_SC),
		UserId:          userId,
		ToUserIdList:    "",
		RoomId:          "0",
		GuildId:         "0",
		Gid:             "",
		TotalAmount:     exchangeAmount.String(),
		PayAmount:       exchangeAmount.String(),
		DiscountsAmount: "0",
		Num:             1,
		Currency:        accountBook.CURRENCY_DIAMOND,
		AppId:           "",
		OrderType:       accountBook.ChangeDiamondStarlightExchange,
		OrderStatus:     1,
		PayStatus:       1,
		Note:            "星光兑换钻石",
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		StatDate:        time.Now().Format(time.DateOnly),
	}
	err := tx.Create(orderInfo).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 增加用户钻石
	accountVersion := service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:        tx,
		UserId:    userId,
		Num:       1,
		Currency:  accountBook.CURRENCY_DIAMOND,
		FundFlow:  accountBook.FUND_INFLOW,
		Amount:    cast.ToString(req.ExchangeAmount),
		OrderId:   orderInfo.OrderId,
		OrderType: accountBook.ChangeDiamondStarlightExchange,
		Note:      "星光兑换钻石",
	})
	// 扣除用户星光
	service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:             tx,
		UserId:         userId,
		Num:            1,
		Currency:       accountBook.CURRENCY_STARLIGHT,
		FundFlow:       accountBook.FUND_OUTFLOW,
		Amount:         cast.ToString(req.ExchangeAmount),
		OrderId:        orderInfo.OrderId,
		OrderType:      accountBook.ChangeStarlightStarlightExchange,
		Note:           "星光兑换钻石",
		AccountVersion: accountVersion,
	})
	tx.Commit()

	accountDTO := service_user.GetUserAccountInfo(userId)
	res = response_user.UserAccountRes{
		UserId:              accountDTO.UserId,
		Status:              accountDTO.Status,
		WithdrawStatus:      accountDTO.WithdrawStatus,
		DiamondAmount:       accountDTO.DiamondAmount,
		StarlightAmount:     accountDTO.StarlightAmount,
		StarlightUnWithdraw: accountDTO.StarlightUnWithdraw,
		StarlightWithdraw:   accountDTO.StarlightWithdraw,
		StarlightSubsidy:    accountDTO.StarlightSubsidy,
	}
	return
}

// 社区后台修改用户的钻石账户
func (u *UserAccount) AccountChangeToAdmin(c *gin.Context, req *request_inner.AccountChangeReq) (err error) {
	userDao := dao.UserDao{}
	userInfo, err := userDao.FindOne(&model.User{UserNo: req.UserNo})
	if err != nil {
		return
	}
	if userInfo.Id == "" {
		return
	}
	account := dao.UserAccountDao{}
	resAccount, _ := account.GetUserAccountByUserId(userInfo.Id)

	var totalAmount string
	if req.FundFlow == accountBook.FUND_OUTFLOW && req.Money > cast.ToFloat64(resAccount.DiamondAmount) {
		totalAmount = resAccount.DiamondAmount
	} else {
		totalAmount = cast.ToString(req.Money)
	}
	tx := coreDb.GetMasterDb().Begin()
	orderInfo := &model.Order{
		ID:              0,
		OrderId:         new(accountBook.Order).OrderNum(accountBook.ORDER_SC),
		UserId:          userInfo.Id,
		ToUserIdList:    "",
		RoomId:          "0",
		GuildId:         "0",
		Gid:             "",
		TotalAmount:     totalAmount,
		PayAmount:       totalAmount,
		DiscountsAmount: "0",
		Num:             1,
		Currency:        accountBook.CURRENCY_DIAMOND,
		AppId:           "",
		OrderType:       accountBook.ChangeDiamondOperationGift,
		OrderStatus:     1,
		PayStatus:       1,
		Note:            req.Note + "*" + req.AdminName,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		StatDate:        time.Now().Format(time.DateOnly),
	}
	err = tx.Create(orderInfo).Error
	if err != nil {
		tx.Rollback()
		return
	}
	// 操作用户钻石
	service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:        tx,
		UserId:    userInfo.Id,
		Num:       1,
		Currency:  accountBook.CURRENCY_DIAMOND,
		FundFlow:  req.FundFlow,
		Amount:    totalAmount,
		OrderId:   orderInfo.OrderId,
		OrderType: accountBook.ChangeDiamondOperationGift,
		Note:      orderInfo.Note,
	})
	//操作系统的账户钻石
	var systemFundFlow = accountBook.FUND_INFLOW
	if req.FundFlow == accountBook.FUND_INFLOW {
		systemFundFlow = accountBook.FUND_OUTFLOW
	}
	service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:        tx,
		UserId:    enum.OperationUserId,
		Num:       1,
		Currency:  accountBook.CURRENCY_DIAMOND,
		FundFlow:  systemFundFlow,
		Amount:    totalAmount,
		OrderId:   orderInfo.OrderId,
		OrderType: accountBook.ChangeDiamondOperationGift,
		Note:      orderInfo.Note,
	})
	tx.Commit()
	log.Println("ChangeAccount Success")
	return
}

func (u *UserAccount) RechargeDiamond(c *gin.Context, userNo, platform, channel string) (res response_user.RechargeDiamondRes) {
	userDao := dao.UserDao{}
	var userInfo = new(model.User)
	var err error
	if userNo == "" {
		userInfo.Id = handle.GetUserId(c)
	} else {
		userInfo.UserNo = userNo
	}
	userInfo, err = userDao.FindOne(userInfo)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodePayChannelError,
			Msg:  nil,
		})
	}
	res.UserId = userInfo.Id
	res.UserNickname = userInfo.Nickname
	res.UserNo = userInfo.UserNo
	res.UserAvatar = coreConfig.GetHotConf().ImagePrefix + userInfo.Avatar

	userAccountDao := dao.UserAccountDao{}
	userAccount, err := userAccountDao.GetUserAccountByUserId(userInfo.Id)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodePayChannelError,
			Msg:  nil,
		})
	}
	res.Status = userAccount.Status
	res.DiamondAmount = userAccount.DiamondAmount

	if platform == "" || channel == "" {
		header := handle.GetHeaderData(c)
		platform = header.Platform
		channel = header.Channel
	}
	channelDao := dao.ConfigChannelDao{}
	channelRes, err := channelDao.First(platform, channel)
	res.ChannelInfo = channelRes
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodePayChannelError,
			Msg:  nil,
		})
	}
	if channelRes.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodePayChannelError,
			Msg:  nil,
		})
	}
	diamondDao := dao.ConfigDiamondDao{}
	dis, err := diamondDao.Find(channelRes.Platform)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodePayChannelError,
			Msg:  nil,
		})
	}
	for _, v := range dis {
		item := response_user.NewConfigDiamondRes{
			Keys:     v.Keys,
			Nums:     v.Nums,
			GotoNums: v.GotoNums,
		}
		res.ChannelGoods = append(res.ChannelGoods, item)
	}
	return
}

// OperationChangeAccount
//
//	@Description: 运营后台变动用户账户
//	@receiver u
//	@param c *gin.Context -
//	@param req *request_inner.OperationChangeAccountReq -
//	@return err -
func (u *UserAccount) OperationChangeAccount(c *gin.Context, req *request_inner.OperationChangeAccountReq) (err error) {
	tx := coreDb.GetMasterDb().Begin()
	_ = service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:                tx,
		UserId:            req.UserId,
		FromUserId:        "",
		ToUserIdList:      "",
		Gid:               "",
		Num:               0,
		Currency:          req.Currency,
		FundFlow:          req.FundFlow,
		Amount:            req.Amount,
		OrderId:           req.OrderId,
		OrderType:         req.OrderType,
		RoomId:            req.RoomId,
		GuildId:           req.GuildId,
		Note:              req.Note,
		SubsidyType:       req.SubsidyType,
		SubsidyAmountType: req.SubsidyAmountType,
	})
	tx.Commit()
	return nil
}
