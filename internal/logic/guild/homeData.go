package guild

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"sort"
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
	"yfapi/typedef/request/user"
	response_guild "yfapi/typedef/response/guild"
	"yfapi/util/easy"
)

type HomeData struct {
}

func (h *HomeData) AccountInfoLogic(c *gin.Context) (resp *response_guild.AccountInfoRes) {
	userId := helper.GetUserId(c)
	guildId := helper.GetGuildId(c)
	// 查询玩家公会账户信息
	accountDao := new(dao.UserAccountDao)
	subsidyInfo, err := accountDao.GetUserAccountGuildSubsidy(userId, guildId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据不存在 初始化
		subsidyInfo = model.UserAccountSubsidy{
			Id:            0,
			UserId:        userId,
			AccountType:   2,
			RoomId:        "0",
			GuildId:       guildId,
			Status:        1,
			SubsidyAmount: "0",
		}
		err = coreDb.GetMasterDb().Create(&subsidyInfo).Error
		if err != nil {
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
	} else if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 查询补贴账户总收益
	totalIncome, _ := accountDao.GetUserGuildHistorySubsidyAmount(userId)
	res := &response_guild.AccountInfoRes{
		UserId:          userId,
		Status:          subsidyInfo.Status,
		CashAmount:      easy.StringFixed(easy.StringToDecimal(subsidyInfo.SubsidyAmount).Div(decimal.NewFromInt(10))),
		TotalCashIncome: easy.StringFixed(easy.StringToDecimal(totalIncome).Div(decimal.NewFromInt(10))),
	}
	//查询手续费和提现说明
	userWithdrawDao := dao.UserWithdrawDao{}
	withdrawInfo, err := userWithdrawDao.UserWithdrawInfoById(1)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if withdrawInfo.ID == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeWithdrawConfigNotExist,
			Msg:  nil,
		})
	}
	res.Desc = withdrawInfo.Desc
	res.SettlementRate = withdrawInfo.SettlementRate
	//查询用户信息
	userInfo := service_user.GetUserBaseInfo(userId)
	res.Mobile = userInfo.Mobile
	res.ReginCode = userInfo.RegionCode
	res.TrueName = userInfo.TrueName
	userBankDao := new(dao.UserBankDao)
	bankInfo, err := userBankDao.GetUserBankListByUserId(userId)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	res.BankList = make([]response_guild.BankInfo, 0)
	if len(bankInfo) > 0 {
		for _, v := range bankInfo {
			bank := response_guild.BankInfo{
				Id:         v.Id,
				BankNo:     v.BankNo,
				BankName:   v.BankName,
				BankHolder: v.BankHolder,
				BankBranch: v.BankBranch,
				IsDefault:  v.IsDefault,
			}
			res.BankList = append(res.BankList, bank)
		}
	}
	return res
}

// ExchangeDiamond
//
//	@Description: 兑换钻石
//	@receiver h
//	@param c *gin.Context -
//	@param req *user.ExchangeDiamondReq -
//	@return res -
func (h *HomeData) ExchangeDiamond(c *gin.Context, req *user.ExchangeDiamondReq) (res response_guild.AccountInfoRes) {
	userId := helper.GetUserId(c)
	guildId := helper.GetGuildId(c)
	// 兑换数量
	exchangeAmount := decimal.NewFromInt(req.ExchangeAmount)

	tx := coreDb.GetMasterDb().Begin()
	orderId := new(accountBook.Order).OrderNum(accountBook.ORDER_SC)
	// 扣除用户星光
	service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
		Tx:        tx,
		UserId:    userId,
		GuildId:   guildId,
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

	// 查询玩家公会账户信息
	accountDao := new(dao.UserAccountDao)
	subsidyInfo, _ := accountDao.GetUserAccountGuildSubsidy(userId, guildId)
	res = response_guild.AccountInfoRes{
		UserId:     userId,
		Status:     1,
		CashAmount: easy.StringFixed(easy.StringToDecimal(subsidyInfo.SubsidyAmount).Div(decimal.NewFromInt(10))),
	}
	return
}

// GetGuildStatInfo 公会数据统计
func (h *HomeData) GetGuildStatInfo(c *gin.Context) (res response_guild.StatGuildInfo) {
	guildId := helper.GetGuildId(c)
	statKey := redisKey.GuildStatInfoKey(guildId)
	cacheData := coreRedis.GetChatroomRedis().Get(c, statKey).Val()
	if len(cacheData) > 0 {
		_ = json.Unmarshal([]byte(cacheData), &res)
		return
	}

	// 聊天室数量
	coreDb.GetSlaveDb().Model(model.Room{}).Where("guild_id=? and live_type=? and status<?", guildId, enum.LiveTypeChatroom, enum.RoomStatusInvalid).Count(&res.Room.ChatroomCount)
	// 直播间数量
	coreDb.GetSlaveDb().Model(model.Room{}).Where("guild_id=? and live_type=? and status<?", guildId, enum.LiveTypeAnchor, enum.RoomStatusInvalid).Count(&res.Room.AnchorRoomCount)
	// 房间总数
	res.Room.RoomCount = res.Room.ChatroomCount + res.Room.AnchorRoomCount

	// 公会资质用户
	coreDb.GetSlaveDb().Table("t_user_practitioner_cred upc").Joins("left join t_guild_member gm on gm.user_id=upc.user_id").Where("gm.guild_id=? and upc.status=1", guildId).
		Group("upc.user_id").Count(&res.Member.CertCount)
	// 公会成员人数
	coreDb.GetSlaveDb().Model(model.GuildMember{}).Where("guild_id=? and status<?", guildId, enum.GuildMemberStatusLeave).Count(&res.Member.MemberCount)
	// 公会普通成员
	res.Member.NormalCount = res.Member.MemberCount - res.Member.CertCount

	// 从业者数量
	var data []struct {
		UserId   string
		TypeList string
	}
	coreDb.GetSlaveDb().Table("t_user_practitioner up").Joins("left join t_guild_member gm on gm.user_id=up.user_id").Where("gm.guild_id=? and up.status=1", guildId).
		Select("up.user_id, GROUP_CONCAT(DISTINCT up.practitioner_type SEPARATOR ',') type_list").Group("up.user_id").Scan(&data)
	res.Practitioners.PractitionersCount = len(data)
	for _, info := range data {
		dst := strings.Split(info.TypeList, ",")
		for _, _type := range dst {
			switch cast.ToInt(_type) {
			case enum.UserPractitionerCompere:
				res.Practitioners.CompereCount++
			case enum.UserPractitionerMusician:
				res.Practitioners.MusicianCount++
			case enum.UserPractitionerCounselor:
				res.Practitioners.CounselorCount++
			case enum.UserPractitionerAnchor:
				res.Practitioners.AnchorCount++
			}
		}
	}

	// 今日房间流水
	coreDb.GetSlaveDb().Model(model.Order{}).Where("guild_id=? and stat_date=? and order_type=?", guildId, time.Now().Format(time.DateOnly), accountBook.ChangeDiamondRewardGift).
		Select("sum(total_amount)").Scan(&res.TodayProfit.TotalProfit)
	// 今日聊天室流水
	coreDb.GetSlaveDb().Table("t_order o").Joins("left join t_room r on r.id=o.room_id").Where("r.guild_id=? and r.live_type=?", guildId, enum.LiveTypeChatroom).
		Where("o.stat_date=? and o.order_type=?", time.Now().Format(time.DateOnly), accountBook.ChangeDiamondRewardGift).Select("sum(o.total_amount)").Scan(&res.TodayProfit.ChatroomProfit)
	// 今日直播间流水
	res.TodayProfit.AnchorProfit = easy.StringFixed(easy.StringToDecimal(res.TodayProfit.TotalProfit).Sub(easy.StringToDecimal(res.TodayProfit.ChatroomProfit)))

	// 本月房间流水
	coreDb.GetSlaveDb().Model(model.Order{}).Where("stat_date between ? and ?", easy.GetCurrMonthStartTime(time.Now()).Format(time.DateOnly), time.Now().Format(time.DateOnly)).
		Where("guild_id=? and order_type=?", guildId, accountBook.ChangeDiamondRewardGift).Select("sum(total_amount)").Scan(&res.MothProfit.TotalProfit)
	// 本月聊天室流水
	coreDb.GetSlaveDb().Table("t_order o").Joins("left join t_room r on r.id=o.room_id").Where("r.guild_id=? and r.live_type=?", guildId, enum.LiveTypeChatroom).
		Where("o.stat_date between ? and ?", easy.GetCurrMonthStartTime(time.Now()).Format(time.DateOnly), time.Now().Format(time.DateOnly)).
		Where("o.order_type=?", accountBook.ChangeDiamondRewardGift).Select("sum(o.total_amount)").Scan(&res.MothProfit.ChatroomProfit)
	// 本月直播间流水
	res.MothProfit.AnchorProfit = easy.StringFixed(easy.StringToDecimal(res.MothProfit.TotalProfit).Sub(easy.StringToDecimal(res.MothProfit.ChatroomProfit)))

	// 增加缓存
	coreRedis.GetChatroomRedis().Set(c, statKey, easy.JSONStringFormObject(res), 5*time.Minute)
	return
}

// GetGuildProfitInfo
//
//	@Description: 首页公会流水信息
//	@receiver h
//	@param c *gin.Context -
//	@return res -
func (h *HomeData) GetGuildProfitInfo(c *gin.Context) (res response_guild.ProfitGuildInfo) {
	guildId := helper.GetGuildId(c)
	statKey := redisKey.GuildProfitInfoKey(guildId)
	cacheData := coreRedis.GetChatroomRedis().Get(c, statKey).Val()
	if len(cacheData) > 0 {
		_ = json.Unmarshal([]byte(cacheData), &res)
		return
	}

	// 最近七天的流水信息
	startTime := time.Now().AddDate(0, 0, -6)
	coreDb.GetSlaveDb().Model(model.Order{}).Where("stat_date between ? and ?", startTime.Format(time.DateOnly), time.Now().Format(time.DateOnly)).
		Where("guild_id=? and order_type=?", guildId, accountBook.ChangeDiamondRewardGift).
		Select("stat_date, sum(total_amount) profit_amount").Group("stat_date").Order("stat_date").Scan(&res.LatestWeek)
	weekMap := make(map[string]struct{})
	for i, info := range res.LatestWeek {
		currDate, _ := time.ParseInLocation(time.RFC3339, info.StatDate, time.Local)
		info.StatDate = currDate.Format(time.DateOnly)
		res.LatestWeek[i].StatDate = info.StatDate
		weekMap[info.StatDate] = struct{}{}
	}
	for i := 0; i < 7; i++ {
		currDate := startTime.AddDate(0, 0, i).Format(time.DateOnly)
		if _, isExist := weekMap[currDate]; !isExist {
			res.LatestWeek = append(res.LatestWeek, response_guild.ProfitInfo{
				StatDate:     currDate,
				ProfitAmount: "0",
			})
		}
	}
	// 排序
	sort.Slice(res.LatestWeek, func(i, j int) bool {
		iDate, _ := time.ParseInLocation(time.DateOnly, res.LatestWeek[i].StatDate, time.Local)
		jDate, _ := time.ParseInLocation(time.DateOnly, res.LatestWeek[j].StatDate, time.Local)
		return iDate.Before(jDate)
	})

	// 房间占比
	var data []struct {
		RoomType  int
		RoomCount int
	}
	coreDb.GetSlaveDb().Model(model.Room{}).Where("guild_id=? and status<?", guildId, enum.RoomStatusInvalid).Select("room_type, COUNT(*) room_count").Group("room_type").Scan(&data)
	for _, info := range data {
		res.RoomCategory = append(res.RoomCategory, response_guild.RoomTypeCount{
			Name:  enum.RoomType(info.RoomType).String(),
			Count: info.RoomCount,
		})
	}
	// 排序
	sort.Slice(res.RoomCategory, func(i, j int) bool {
		return res.RoomCategory[i].Count > res.RoomCategory[j].Count
	})

	// 流水占比
	var profitData []struct {
		RoomId       string
		ProfitAmount string
	}
	beginTime := easy.GetCurrMonthStartTime(time.Now())
	coreDb.GetSlaveDb().Model(model.Order{}).Where("stat_date between ? and ?", beginTime.Format(time.DateOnly), time.Now().Format(time.DateOnly)).
		Where("guild_id=? and order_type=?", guildId, accountBook.ChangeDiamondRewardGift).Select("room_id, sum(total_amount) profit_amount").Group("room_id").Scan(&profitData)
	var roomIdList []string
	for _, info := range profitData {
		roomIdList = append(roomIdList, info.RoomId)
	}
	roomMap, _ := new(dao.RoomDao).GetRoomMapByIdList(roomIdList)
	profitMap := make(map[int]decimal.Decimal)
	for _, info := range profitData {
		profitMap[roomMap[info.RoomId].RoomType] = profitMap[roomMap[info.RoomId].RoomType].Add(easy.StringToDecimal(info.ProfitAmount))
	}
	for roomType, amount := range profitMap {
		res.ProfitCategory = append(res.ProfitCategory, response_guild.RoomTypeProfit{
			Name:         enum.RoomType(roomType).String(),
			ProfitAmount: easy.StringFixed(amount),
		})
	}
	// 排序
	sort.Slice(res.ProfitCategory, func(i, j int) bool {
		return easy.StringToDecimal(res.ProfitCategory[i].ProfitAmount).GreaterThan(easy.StringToDecimal(res.ProfitCategory[j].ProfitAmount))
	})

	// 增加缓存
	coreRedis.GetChatroomRedis().Set(c, statKey, easy.JSONStringFormObject(res), 5*time.Minute)
	return
}

// GetRoomRankList
//
//	@Description: 查询公会房间流水排行榜
//	@receiver h
//	@param c *gin.Context -
//	@param req *request_login.GetRoomRankListReq -
//	@return res -
func (h *HomeData) GetRoomRankList(c *gin.Context, req *request_login.GetRoomRankListReq) (res []response_guild.RoomRank) {
	guildId := helper.GetGuildId(c)
	statKey := redisKey.GuildRoomRankKey(guildId)
	cacheData := coreRedis.GetChatroomRedis().Get(c, statKey).Val()
	if len(cacheData) > 0 {
		_ = json.Unmarshal([]byte(cacheData), &res)
		return
	}
	// 开始时间
	startTime := time.Now().Format(time.DateOnly)
	if len(req.StartTime) > 0 {
		startTime = req.StartTime
	}
	// 结束时间
	endTime := time.Now().Format(time.DateOnly)
	if len(req.EndTime) > 0 {
		endTime = req.EndTime
	}
	var profitData []struct {
		RoomId       string
		ProfitAmount string
	}
	// 查询流水
	coreDb.GetSlaveDb().Model(model.Order{}).Where("stat_date between ? and ?", startTime, endTime).
		Where("guild_id=? and order_type=?", guildId, accountBook.ChangeDiamondRewardGift).Select("room_id, sum(total_amount) profit_amount").Group("room_id").Order("profit_amount desc").Scan(&profitData)
	// 查询房间信息
	var roomIdList []string
	for _, info := range profitData {
		roomIdList = append(roomIdList, info.RoomId)
	}
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
	//　查询房间开播时长
	var onlineData []struct {
		RoomId       string `json:"roomId"`
		OnlineSecond int64  `json:"onlineSecond"`
	}
	coreDb.GetSlaveDb().Model(model.RoomWheatTime{}).Where("stat_date between ? and ?", startTime, endTime).Where("room_id in ?", roomIdList).Select("room_id, sum(on_time) online_second").
		Group("room_id").Scan(&onlineData)
	onlineMap := make(map[string]string)
	for _, info := range onlineData {
		onlineMap[info.RoomId] = easy.SecondFormatString(info.OnlineSecond)
	}
	// 初始化无开播记录的房间
	for _, roomId := range roomIdList {
		if _, isExist := onlineMap[roomId]; !isExist {
			onlineMap[roomId] = "00:00:00"
		}
	}
	// 构造返回信息
	for _, info := range profitData {
		roomInfo := roomMap[info.RoomId]
		res = append(res, response_guild.RoomRank{
			RoomName:      roomInfo.Name,
			RoomOwnerName: userMap[roomInfo.UserId].Nickname,
			OnlineSecond:  onlineMap[info.RoomId],
			ProfitAmount:  info.ProfitAmount,
		})
	}

	// 增加缓存
	coreRedis.GetChatroomRedis().Set(c, statKey, easy.JSONStringFormObject(res), 5*time.Minute)
	return
}
