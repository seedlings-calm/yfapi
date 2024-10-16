package room

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreDb"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	service_user "yfapi/internal/service/user"
	request_login "yfapi/typedef/request/guild"
	"yfapi/typedef/request/user"
	"yfapi/typedef/response"
	response_guild "yfapi/typedef/response/guild"
	"yfapi/typedef/response/orderBill"
	response_login "yfapi/typedef/response/roomOwner"
	"yfapi/util/easy"
)

type RoomHome struct {
}

// ExchangeDiamond
//
//	@Description: 兑换钻石
//	@receiver h
//	@param c *gin.Context -
//	@param req *user.ExchangeDiamondReq -
//	@return res -
func (r *RoomHome) ExchangeDiamond(c *gin.Context, req *user.ExchangeDiamondReq) (res response_guild.AccountInfoRes) {
	userId := helper.GetUserId(c)
	roomId := helper.GetRoomId(c)
	// 兑换数量
	exchangeAmount := decimal.NewFromInt(req.ExchangeAmount)

	tx := coreDb.GetMasterDb().Begin()
	orderId := new(accountBook.Order).OrderNum(accountBook.ORDER_SC)
	// 扣除用户星光
	service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:        tx,
		UserId:    userId,
		RoomId:    roomId,
		Num:       1,
		Currency:  accountBook.CURRENCY_STARLIGHT_SUBSIDY,
		FundFlow:  accountBook.FUND_OUTFLOW,
		Amount:    cast.ToString(req.ExchangeAmount),
		OrderId:   orderId,
		OrderType: accountBook.ChangeStarlightStarlightExchange,
		Note:      "后台兑换钻石",
	})
	// 增加用户钻石
	_ = service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:        tx,
		UserId:    userId,
		Num:       1,
		Currency:  accountBook.CURRENCY_DIAMOND,
		FundFlow:  accountBook.FUND_INFLOW,
		Amount:    cast.ToString(req.ExchangeAmount),
		OrderId:   orderId,
		OrderType: accountBook.ChangeDiamondStarlightExchange,
		Note:      "后台兑换钻石",
	})
	// 创建兑换订单
	orderInfo := &model.Order{
		ID:              0,
		OrderId:         orderId,
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
		Note:            "后台兑换钻石",
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		StatDate:        time.Now().Format(time.DateOnly),
	}
	err := tx.Create(orderInfo).Error
	if err != nil {
		tx.Rollback()
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	tx.Commit()

	// 查询玩家房间账户信息
	accountDao := new(dao.UserAccountDao)
	subsidyInfo, _ := accountDao.GetUserAccountRoomSubsidy(userId, roomId)
	res = response_guild.AccountInfoRes{
		UserId:     userId,
		Status:     1,
		CashAmount: easy.StringFixed(easy.StringToDecimal(subsidyInfo.SubsidyAmount).Div(decimal.NewFromInt(10))),
	}
	return
}

// GetAccountBillList
//
//	@Description: 查询账户交易明细列表
//	@receiver g
//	@param c *gin.Context -
//	@param req *request_login.GetAccountBillReq -
//	@return res -
func (r *RoomHome) GetAccountBillList(c *gin.Context, req *request_login.GetAccountBillReq) (res response.AdminPageRes) {
	userId := helper.GetUserId(c)
	roomId := helper.GetRoomId(c)
	// 查询会长补贴账户交易记录
	// 房间流水日结 房间流水月结 提现 兑换钻石
	// 开始时间
	startTime := easy.GetCurrMonthStartTime(time.Now())
	// 结束时间
	endTime := easy.GetCurrDayEndTime(time.Now())
	if len(req.DateRange) > 1 {
		startTime, _ = time.ParseInLocation(time.DateOnly, req.DateRange[0], time.Local)
		startTime = easy.GetCurrDayStartTime(startTime)
		endTime, _ = time.ParseInLocation(time.DateOnly, req.DateRange[1], time.Local)
		endTime = easy.GetCurrDayEndTime(endTime)
	}
	// 查询流水记录
	var result []*response_guild.AccountBill
	tx := coreDb.GetSlaveDb().Table("t_order_bill ob").Joins("left join t_order o on o.order_id=ob.order_id").Where("ob.user_id=? and ob.currency=? and ob.room_id=?", userId, accountBook.CURRENCY_STARLIGHT_SUBSIDY, roomId).
		Where("ob.create_time between ? and ?", startTime.Format(time.DateTime), endTime.Format(time.DateTime))
	if len(req.WithdrawStatus) > 0 { // 单独查询提现订单
		tx = tx.Where("ob.order_type=?", accountBook.ChangeStarlightStarlightWithdrawal).Where("o.withdraw_status=?", req.WithdrawStatus)
	}
	if req.FundFlow > 0 {
		tx = tx.Where("ob.fund_flow", req.FundFlow)
	}
	if req.OrderType > 0 {
		tx = tx.Where("ob.order_type", req.OrderType)
	}
	_ = tx.Count(&res.Total)
	err := tx.Select("ob.order_id, ob.order_type, ob.note memo, ob.fund_flow, ob.before_amount, ob.amount, ob.curr_amount, ob.create_time, o.withdraw_status").
		Order("ob.create_time desc").Limit(req.Size).Offset((req.CurrentPage - 1) * req.Size).Scan(&result).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 查询订单信息
	orderIdMap := make(map[string]struct{})
	for _, info := range result {
		orderIdMap[info.OrderId] = struct{}{}
	}
	var orderIdList []string
	for key := range orderIdMap {
		orderIdList = append(orderIdList, key)
	}
	var orderList []model.Order
	_ = coreDb.GetSlaveDb().Model(model.Order{}).Where("order_id in ?", orderIdList).Scan(&orderList).Error
	orderMap := make(map[string]model.Order)
	// 用户ID map
	userIdMap := make(map[string]struct{})
	// 房间ID map
	roomIdMap := make(map[string]struct{})
	for _, info := range orderList {
		orderMap[info.OrderId] = info
		userIdMap[info.UserId] = struct{}{}
		if len(info.RoomId) > 0 {
			roomIdMap[info.RoomId] = struct{}{}
		}
	}
	// 查询房间信息
	var roomIdList []string
	for key := range roomIdMap {
		roomIdList = append(roomIdList, key)
	}
	roomMap := make(map[string]model.Room)
	if len(roomIdList) > 0 {
		roomMap, _ = new(dao.RoomDao).GetRoomMapByIdList(roomIdList)
	}
	// 查询用户信息
	for _, info := range roomMap {
		userIdMap[info.UserId] = struct{}{}
	}
	var userIdList []string
	for key := range userIdMap {
		userIdList = append(userIdList, key)
	}
	userMap := service_user.GetUserBaseInfoMap(userIdList)
	// 处理返回结果
	for _, info := range result {
		info.OrderTypeDesc = accountBook.EnumOrderType(info.OrderType).String()
		if info.OrderType != accountBook.ChangeStarlightStarlightWithdrawal {
			// TODO 非提现订单的状态
			info.WithdrawStatus = 99
		}
		info.BeforeAmount = easy.StringFixed(easy.StringToDecimal(info.BeforeAmount).Div(decimal.NewFromInt(10)))
		info.Amount = easy.StringFixed(easy.StringToDecimal(info.Amount).Div(decimal.NewFromInt(10)))
		info.CurrAmount = easy.StringFixed(easy.StringToDecimal(info.CurrAmount).Div(decimal.NewFromInt(10)))

		switch info.OrderType {
		case accountBook.ChangeStarlightRoomFlowDailySettlement, accountBook.ChangeStarlightRoomFlowMonthlySettlement: // 房间日结、月结
			// 订单信息
			orderInfo := orderMap[info.OrderId]
			statTime, _ := time.ParseInLocation(time.RFC3339, orderInfo.StatDate, time.Local)
			statDate := statTime.Format(time.DateOnly)
			if info.OrderType == accountBook.ChangeStarlightRoomFlowMonthlySettlement {
				statDate = statTime.Format("2006-01")
			}
			// 房间信息
			roomInfo := roomMap[orderInfo.RoomId]
			// 订单备注信息
			note := orderBill.SubsidyChatroomNote{}
			_ = json.Unmarshal([]byte(orderInfo.Note), &note)
			info.SubsidyRoom = &response_guild.SubsidyRoom{
				StateDate:     statDate,
				RoomName:      roomInfo.Name,
				RoomNo:        roomInfo.RoomNo,
				RoomOwnerName: userMap[roomInfo.UserId].Nickname,
				RoomOwnerNo:   userMap[roomInfo.UserId].UserNo,
				SettleName:    userMap[orderInfo.UserId].Nickname,
				SettleNo:      userMap[orderInfo.UserId].UserNo,
				ProfitAmount:  note.ProfitAmount,
				Income:        orderInfo.TotalAmount,
			}
		case accountBook.ChangeStarlightStarlightWithdrawal: // 提现
			// 订单信息
			orderInfo := orderMap[info.OrderId]
			// 用户信息
			userInfo := userMap[orderInfo.UserId]
			// 订单备注信息
			note := orderBill.WithdrawNote{}
			_ = json.Unmarshal([]byte(orderInfo.Note), &note)
			info.Withdraw = &response_guild.WithdrawDetail{
				Nickname:     userInfo.Nickname,
				UserNo:       userInfo.UserNo,
				TrueName:     userInfo.TrueName,
				Mobile:       userInfo.Mobile,
				Amount:       orderInfo.TotalAmount,
				PayAmount:    orderInfo.PayAmount,
				BankUserName: note.BankHolder,
				BankNo:       note.BankNo,
			}
		}
	}

	res.CurrentPage = req.CurrentPage
	res.Size = req.Size
	res.Data = result
	return
}

// RoomSearchUser
//
//	@Description: 查询用户信息
//	@receiver r
//	@param c *gin.Context -
//	@return res -
func (r *RoomHome) RoomSearchUser(c *gin.Context) (res response_login.SearchUser) {
	roomId := helper.GetRoomId(c)
	keyword, _ := c.GetQuery("userNo")
	if len(keyword) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeParam,
			Msg:  nil,
		})
	}
	info, err := new(dao.UserDao).FindUserByUserNo(keyword)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeCheckUserID,
			Msg:  nil,
		})
	} else if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	if err != nil || room.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	isExist := new(dao.GuildDao).GetCheckUserInGuild(room.GuildId, info.Id)
	if isExist == false {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeNotGuildMember,
			Msg:  nil,
		})
	}
	res = response_login.SearchUser{
		UserId:   info.Id,
		UserNo:   info.UserNo,
		Nickname: info.Nickname,
		Avatar:   helper.FormatImgUrl(info.Avatar),
	}
	return
}
