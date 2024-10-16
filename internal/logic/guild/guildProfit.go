package guild

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"strings"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_login "yfapi/typedef/request/guild"
	"yfapi/typedef/response"
	response_guild "yfapi/typedef/response/guild"
	"yfapi/typedef/response/orderBill"
	"yfapi/util/easy"
)

type GuildProfit struct {
}

// GetGuildMemberProfitList
//
//	@Description: 公会成员礼物流水列表
//	@param c *gin.Context -
//	@param req *request_login.GetGuildMemberProfitReq -
//	@return res -
func (g *GuildProfit) GetGuildMemberProfitList(c *gin.Context, req *request_login.GetGuildMemberProfitReq) (res response.AdminPageRes) {
	guildId := helper.GetGuildId(c)
	// 开始时间
	startTime := easy.GetCurrWeekStartTime(time.Now())
	// 结束时间
	endTime := easy.GetCurrWeekEndTime(time.Now())
	if len(req.DateRange) > 1 {
		startTime, _ = time.ParseInLocation(time.DateOnly, req.DateRange[0], time.Local)
		startTime = easy.GetCurrDayStartTime(startTime)
		endTime, _ = time.ParseInLocation(time.DateOnly, req.DateRange[1], time.Local)
		endTime = easy.GetCurrDayEndTime(endTime)
	}

	// 公会成员
	var result []*response_guild.MemberProfit
	tx := coreDb.GetSlaveDb().Table("t_order_bill ob").Joins("left join t_guild_member gm on gm.user_id=ob.user_id").Where("ob.create_time between ? and ?", startTime.Format(time.DateTime), endTime.Format(time.DateTime)).
		Where("gm.guild_id=? and ob.guild_id=? and ob.order_type=?", guildId, guildId, accountBook.ChangeStarlightRewardIncome).
		Group("ob.user_id")
	if len(req.UserKeyword) > 0 {
		tx = tx.Joins("left join t_user u on u.id=gm.user_id").Where("u.nickname like ? or u.user_no like ?", easy.GenLikeSql(req.UserKeyword), easy.GenLikeSql(req.UserKeyword))
	}
	if req.PractitionersType > 0 {
		tx = tx.Joins("left join t_user_practitioner up on up.user_id=gm.user_id").Where("up.status=1 and up.practitioner_type=?", req.PractitionersType)
	}
	tx.Count(&res.Total)
	err := tx.Select("ob.user_id, sum(ob.diamond) profit_amount, count(*) reward_count").Order("profit_amount desc").Limit(req.Size).Offset((req.CurrentPage - 1) * req.Size).Scan(&result).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 成员id列表
	var userIdList []string
	for _, info := range result {
		userIdList = append(userIdList, info.UserId)
	}
	// 用户信息map
	userMap := service_user.GetUserBaseInfoMap(userIdList)
	// 查询成员从业者身份
	var data []struct {
		UserId   string
		TypeList string
	}
	err = coreDb.GetSlaveDb().Table("t_user_practitioner up").Where("up.user_id in ? and up.status=1", userIdList).
		Select("up.user_id, GROUP_CONCAT(DISTINCT up.practitioner_type SEPARATOR ',') type_list").Group("up.user_id").Scan(&data).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	pracMap := make(map[string]string)
	for _, info := range data {
		// 从业者身份处理
		typeList := strings.Split(info.TypeList, ",")
		var pracList []string
		for _, _type := range typeList {
			pracList = append(pracList, enum.PractitionerType(cast.ToInt(_type)).String())
		}
		pracMap[info.UserId] = strings.Join(pracList, "、")
	}
	// 是否计算预估收益
	isCalcIncome := false
	// 公会流水补贴比例
	incomeRate := decimal.Zero
	if startTime.Year() == endTime.Year() && startTime.Month() == endTime.Month() {
		isCalcIncome = true
		// 查询公会收益配置
		var profitConfigList []model.SubsidyConfigGuild
		err = coreDb.GetSlaveDb().Model(model.SubsidyConfigGuild{}).Where("subsidy_type=1").Order("profit_num desc").Scan(&profitConfigList).Error
		if err != nil {
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		// 查询公会聊天室月流水
		profitAmount := easy.StringToDecimal(getGuildChatroomMonthProfit(guildId, startTime))
		for _, info := range profitConfigList {
			if profitAmount.GreaterThanOrEqual(decimal.NewFromInt(int64(info.ProfitNum))) {
				incomeRate = easy.StringToDecimal(info.ProfitRate)
				break
			}
		}
	}
	// 处理返回结果
	// 统计日期
	statDate := fmt.Sprintf("%v ~ %v", startTime.Format(time.DateOnly), endTime.Format(time.DateOnly))
	for _, info := range result {
		info.StatDate = statDate
		info.UserNo = userMap[info.UserId].UserNo
		info.Nickname = userMap[info.UserId].Nickname
		info.Practitioner = pracMap[info.UserId]
		if isCalcIncome {
			// 计算预估收益
			info.Income = easy.StringFixed(easy.StringToDecimal(info.ProfitAmount).Mul(incomeRate).Div(decimal.NewFromInt(100)))
		}
	}
	res.CurrentPage = req.CurrentPage
	res.Size = req.Size
	res.Data = result
	return
}

func getGuildChatroomMonthProfit(guildId string, statMonth time.Time) (profitAmount string) {
	ctx := context.Background()
	starTime := easy.GetCurrMonthStartTime(statMonth).Format(time.DateOnly)
	endTime := easy.GetCurrMonthEndTime(statMonth).Format(time.DateOnly)
	statKey := redisKey.GuildChatroomMonthProfitKey(guildId, starTime)
	profitAmount = coreRedis.GetChatroomRedis().Get(ctx, statKey).Val()
	if len(profitAmount) > 0 {
		return
	}
	// 月聊天室流
	coreDb.GetSlaveDb().Table("t_order o").Joins("left join t_room r on r.id=o.room_id").Where("r.guild_id=? and r.live_type=?", guildId, enum.LiveTypeChatroom).
		Where("o.stat_date between ? and ?", starTime, endTime).
		Where("o.order_type=?", accountBook.ChangeDiamondRewardGift).Select("sum(o.total_amount)").Scan(&profitAmount)
	// 增加缓存
	coreRedis.GetChatroomRedis().Set(ctx, statKey, profitAmount, 10*time.Minute)
	return
}

// GetGuildRoomProfitList
//
//	@Description: 公会房间礼物流水列表
//	@param c *gin.Context -
//	@param req *request_login.GetGuildRoomProfitReq -
//	@return res -
func (g *GuildProfit) GetGuildRoomProfitList(c *gin.Context, req *request_login.GetGuildRoomProfitReq) (res response.AdminPageRes) {
	guildId := helper.GetGuildId(c)
	// 开始时间
	startTime := easy.GetCurrDayStartTime(time.Now())
	// 结束时间
	endTime := easy.GetCurrDayEndTime(time.Now())
	if len(req.DateRange) > 1 {
		startTime, _ = time.ParseInLocation(time.DateOnly, req.DateRange[0], time.Local)
		startTime = easy.GetCurrDayStartTime(startTime)
		endTime, _ = time.ParseInLocation(time.DateOnly, req.DateRange[1], time.Local)
		endTime = easy.GetCurrDayEndTime(endTime)
	}

	// 房间流水
	var result []*response_guild.RoomProfit
	tx := coreDb.GetSlaveDb().Table("t_order ob").Joins("left join t_room r on r.id=ob.room_id").Where("ob.create_time between ? and ?", startTime.Format(time.DateTime), endTime.Format(time.DateTime)).
		Where("r.guild_id=? and ob.order_type=?", guildId, accountBook.ChangeDiamondRewardGift).Group("ob.room_id")
	if len(req.RoomKeyword) > 0 {
		tx = tx.Where("r.name like ? or r.room_no like ?", easy.GenLikeSql(req.RoomKeyword), easy.GenLikeSql(req.RoomKeyword))
	}
	if req.RoomType > 0 {
		tx = tx.Where("r.room_type=?", req.RoomType)
	}
	tx.Count(&res.Total)
	err := tx.Select("ob.room_id, sum(ob.total_amount) profit_amount, count(*) reward_count").Order("profit_amount desc").Limit(req.Size).Offset((req.CurrentPage - 1) * req.Size).Scan(&result).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 房间id列表
	var roomIdList []string
	for _, info := range result {
		roomIdList = append(roomIdList, info.RoomId)
	}
	// 房间信息map
	roomMap, _ := new(dao.RoomDao).GetRoomMapByIdList(roomIdList)
	// 查询房主信息
	userIdMap := make(map[string]struct{})
	for _, info := range roomMap {
		userIdMap[info.UserId] = struct{}{}
	}
	var userIdList []string
	for key := range userIdMap {
		userIdList = append(userIdList, key)
	}
	userMap := service_user.GetUserBaseInfoMap(userIdList)
	// 是否计算预估收益
	isCalcIncome := false
	// 公会流水补贴比例
	incomeRate := decimal.Zero
	if startTime.Year() == endTime.Year() && startTime.Month() == endTime.Month() {
		isCalcIncome = true
		// 查询公会收益配置
		var profitConfigList []model.SubsidyConfigGuild
		err = coreDb.GetSlaveDb().Model(model.SubsidyConfigGuild{}).Where("subsidy_type=1").Order("profit_num desc").Scan(&profitConfigList).Error
		if err != nil {
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		// 查询公会聊天室月流水
		getAmount := getGuildChatroomMonthProfit(guildId, startTime)
		profitAmount := easy.StringToDecimal(getAmount)
		for _, info := range profitConfigList {
			if profitAmount.GreaterThanOrEqual(decimal.NewFromInt(int64(info.ProfitNum))) {
				incomeRate = easy.StringToDecimal(info.ProfitRate)
				break
			}
		}
	}
	// 处理结果返回
	for _, info := range result {
		roomInfo := roomMap[info.RoomId]
		info.RoomNo = roomInfo.RoomNo
		info.RoomName = roomInfo.Name
		info.RoomType = enum.RoomType(roomInfo.RoomType).String()
		info.RoomOwnerNo = userMap[roomInfo.UserId].UserNo
		info.RoomOwnerName = userMap[roomInfo.UserId].Nickname
		if isCalcIncome {
			// 计算预估收益
			info.Income = easy.StringFixed(easy.StringToDecimal(info.ProfitAmount).Mul(incomeRate).Div(decimal.NewFromInt(100)))
		}
	}
	res.CurrentPage = req.CurrentPage
	res.Size = req.Size
	res.Data = result
	return
}

// GetGuildRewardList
//
//	@Description: 公会礼物打赏详情列表
//	@param c *gin.Context -
//	@param req *request_login.GetGuildRewardListReq -
//	@return res -
func (g *GuildProfit) GetGuildRewardList(c *gin.Context, req *request_login.GetGuildRewardListReq) (res response.AdminPageRes) {
	// 开始时间
	startTime := easy.GetCurrDayStartTime(time.Now())
	// 结束时间
	endTime := easy.GetCurrDayEndTime(time.Now())
	if len(req.DateRange) > 1 {
		startTime, _ = time.ParseInLocation(time.DateOnly, req.DateRange[0], time.Local)
		startTime = easy.GetCurrDayStartTime(startTime)
		endTime, _ = time.ParseInLocation(time.DateOnly, req.DateRange[1], time.Local)
		endTime = easy.GetCurrDayEndTime(endTime)
	}
	var result []*response_guild.RewardDetail
	tx := coreDb.GetSlaveDb().Table("t_order_bill ob").Where("ob.create_time between ? and ?", startTime.Format(time.DateTime), endTime.Format(time.DateTime)).
		Where("ob.order_type=?", accountBook.ChangeStarlightRewardIncome)
	if req.RewardType == 1 {
		tx = tx.Where("ob.user_id=?", req.Uid)
	} else {
		tx = tx.Where("ob.room_id=?", req.Uid)
	}
	_ = tx.Count(&res.Total)
	err := tx.Select("ob.user_id to_user_id, ob.from_user_id, ob.note gift_name, ob.num gift_count, ob.diamond gift_price, ob.create_time").Order("create_time desc").
		Limit(req.Size).Offset((req.CurrentPage - 1) * req.Size).Scan(&result).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 查询用户信息
	userIdMap := make(map[string]struct{})
	for _, info := range result {
		userIdMap[info.FromUserId] = struct{}{}
		userIdMap[info.ToUserId] = struct{}{}
	}
	var userIdList []string
	for key := range userIdMap {
		userIdList = append(userIdList, key)
	}
	userMap := service_user.GetUserBaseInfoMap(userIdList)
	// 处理返回结果
	for _, info := range result {
		info.FromUserNo = userMap[info.FromUserId].UserNo
		info.FromNickname = userMap[info.FromUserId].Nickname
		info.ToUserNo = userMap[info.ToUserId].UserNo
		info.ToNickname = userMap[info.ToUserId].Nickname
		info.GiftPrice /= info.GiftCount
	}
	res.CurrentPage = req.CurrentPage
	res.Size = req.Size
	res.Data = result
	return
}

// GetAccountBillList
//
//	@Description: 查询账户交易明细列表
//	@receiver g
//	@param c *gin.Context -
//	@param req *request_login.GetAccountBillReq -
//	@return res -
func (g *GuildProfit) GetAccountBillList(c *gin.Context, req *request_login.GetAccountBillReq) (res response.AdminPageRes) {
	userId := helper.GetUserId(c)
	//guildId := helper.GetGuildId(c)
	// 查询会长补贴账户交易记录
	// 房间流水日结 房间流水月结 公会直播间补贴月结 公会流水补贴月结 提现 兑换钻石
	//var findOrderList = []int{accountBook.ChangeStarlightRoomFlowDailySettlement, accountBook.ChangeStarlightRoomFlowMonthlySettlement, accountBook.ChangeStarlightGuildLiveRoomSubsidyMonthlySettlement,
	//	accountBook.ChangeStarlightGuildFlowSubsidyMonthlySettlement, accountBook.ChangeStarlightWithdrawalReturn, accountBook.ChangeStarlightStarlightExchange}
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
	tx := coreDb.GetSlaveDb().Table("t_order_bill ob").Joins("left join t_order o on o.order_id=ob.order_id").Where("ob.user_id=? and ob.currency=?", userId, accountBook.CURRENCY_STARLIGHT_SUBSIDY).
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
	// 公会ID map
	guildIdMap := make(map[string]struct{})
	for _, info := range orderList {
		orderMap[info.OrderId] = info
		userIdMap[info.UserId] = struct{}{}
		if len(info.RoomId) > 0 {
			roomIdMap[info.RoomId] = struct{}{}
		}
		if len(info.GuildId) > 0 {
			guildIdMap[info.GuildId] = struct{}{}
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
	// 查询公会信息
	var guildIdList []string
	for key := range guildIdMap {
		guildIdList = append(guildIdList, key)
	}
	guildMap := make(map[string]model.Guild)
	if len(guildIdList) > 0 {
		guildMap, _ = new(dao.GuildDao).GetGuildMapByIdList(guildIdList)
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
		case accountBook.ChangeStarlightGuildLiveRoomSubsidyMonthlySettlement: // 公会有效直播间月结
			// 订单信息
			orderInfo := orderMap[info.OrderId]
			statTime, _ := time.ParseInLocation(time.RFC3339, orderInfo.StatDate, time.Local)
			statDate := statTime.Format("2006-01")
			// 公会信息
			guildInfo := guildMap[orderInfo.GuildId]
			// 订单备注信息
			note := orderBill.SubsidyGuildValidNote{}
			_ = json.Unmarshal([]byte(orderInfo.Note), &note)
			info.SubsidyGuildValid = &response_guild.SubsidyGuildValidMonth{
				StatDate:     statDate,
				GuildName:    guildInfo.Name,
				GuildNo:      guildInfo.GuildNo,
				RoomCount:    len(note.List),
				ValidCount:   note.ValidCount,
				ProfitAmount: note.ProfitAmount,
				Income:       orderInfo.TotalAmount,
			}
		case accountBook.ChangeStarlightGuildFlowSubsidyMonthlySettlement: // 公会流水月结
			// 订单信息
			orderInfo := orderMap[info.OrderId]
			statTime, _ := time.ParseInLocation(time.RFC3339, orderInfo.StatDate, time.Local)
			statDate := statTime.Format("2006-01")
			// 公会信息
			guildInfo := guildMap[orderInfo.GuildId]
			// 订单备注信息
			note := orderBill.SubsidyGuildNote{}
			_ = json.Unmarshal([]byte(orderInfo.Note), &note)
			info.SubsidyGuild = &response_guild.SubsidyGuildMonth{
				StatDate:     statDate,
				GuildName:    guildInfo.Name,
				GuildNo:      guildInfo.GuildNo,
				ProfitAmount: note.ProfitAmount,
				Income:       orderInfo.TotalAmount,
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
