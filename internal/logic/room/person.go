package room

import (
	"sort"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/service/accountBook"
	"yfapi/typedef/response/roomOwner"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type PersonDao struct {
	RoomId string
	UserId string
}

// 开播时长结构体
type wheatTimeStruct struct {
	UserId string    `json:"userId"`
	Times  int       `json:"times"` //开播总秒数
	Date   time.Time `json:"date"`
}

// 流水交易结构体
type orderResStruct struct {
	UserId       string    `json:"fromUserId"`
	ProfitAmount int       `json:"profitAmount"`
	RewardCount  int       `json:"rewardCount"`
	Date         time.Time `json:"date"`
}

func (pd *PersonDao) PersonList(c *gin.Context, start, end, sortKey string) (res []*roomOwner.PersonListRes) {
	if pd.RoomId == "" || pd.UserId == "" {
		return
	}
	userPDao := dao.DaoUserPractitioner{}
	userLists, err := userPDao.FindByRoomId(pd.RoomId)
	if err != nil {
		panic(i18n_err.ErrorCodeReadDB)
	}
	if len(userLists) == 0 {
		return
	}
	var userIds []string
	for _, v := range userLists {
		v.Avatar = coreConfig.GetHotConf().ImagePrefix + v.Avatar
		userIds = append(userIds, v.UserId)
	}
	res = append(res, userLists...)
	// 开始时间
	startTime := easy.GetCurrDayStartTime(time.Now())
	// 结束时间
	endTime := easy.GetCurrDayEndTime(time.Now())
	if start != "" && end != "" {
		startTime, _ = time.ParseInLocation(time.DateOnly, start, time.Local)
		startTime = easy.GetCurrDayStartTime(startTime)
		endTime, _ = time.ParseInLocation(time.DateOnly, end, time.Local)
		endTime = easy.GetCurrDayEndTime(endTime)
	}

	//房间相关的订单
	var orderRes []orderResStruct
	coreDb.GetMasterDb().Table("t_order_bill").
		Where("room_id = ? and order_type = ? and from_user_id in ?", pd.RoomId, accountBook.ChangeStarlightRewardIncome, userIds).
		Where("create_time between ? and ?", startTime, endTime).
		Select("from_user_id,IFNULL(sum(diamond),0) as profit_amount,count(*) as reward_count").
		Group("from_user_id").
		Scan(&orderRes)

	var wheatTime []wheatTimeStruct
	coreDb.GetMasterDb().Table("t_room_wheat_time").Where("room_id = ?", pd.RoomId).Where("stat_date between ? and ?", startTime.Format(time.DateOnly), endTime.Format(time.DateOnly)).Select("user_id,sum(on_time) as times").Group("user_id").Scan(&wheatTime)
	var (
		orderResMap  map[string]orderResStruct
		wheatTimeMap map[string]wheatTimeStruct
	)
	if len(orderRes) > 0 {
		for _, v := range orderRes {
			orderResMap[v.UserId] = v
		}
	}
	if len(wheatTime) > 0 {
		for _, v := range wheatTime {
			wheatTimeMap[v.UserId] = v
		}
	}
	for _, v := range res {
		if _, ok := orderResMap[v.UserId]; ok {
			v.ProfitAmount = cast.ToInt(orderResMap[v.UserId].ProfitAmount)
			v.RewardCount = orderResMap[v.UserId].RewardCount
		}
		if _, ok := wheatTimeMap[v.UserId]; ok {
			t := time.Unix(int64(wheatTimeMap[v.UserId].Times), 0)
			v.TimesNum = wheatTimeMap[v.UserId].Times
			v.Times = t.Format(time.TimeOnly)
		}
	}
	sort.Slice(res, func(i, j int) bool {
		switch sortKey {
		case "timesNum":
			return res[i].TimesNum > res[j].TimesNum
		case "rewardCount":
			return res[i].RewardCount > res[j].RewardCount
		default:
			return res[i].ProfitAmount > res[j].ProfitAmount
		}
	})
	return
}

func (pd *PersonDao) PersonListDetail(c *gin.Context, userId, start, end, sortKey string) (res []*roomOwner.PersonListDetailRes) {
	if pd.RoomId == "" || pd.UserId == "" || userId == "" {
		return
	}
	// 开始时间
	startTime := easy.GetCurrDayStartTime(time.Now())
	// 结束时间
	endTime := easy.GetCurrDayEndTime(time.Now())
	if start != "" && end != "" {
		startTime, _ = time.ParseInLocation(time.DateOnly, start, time.Local)
		startTime = easy.GetCurrDayStartTime(startTime)
		endTime, _ = time.ParseInLocation(time.DateOnly, end, time.Local)
		endTime = easy.GetCurrDayEndTime(endTime)
	}
	coreDb.GetMasterDb().Raw(`select 
	tu.user_no,tu.nickname,tu.avatar,
	tob.note as gift_name,tob.user_id,tob.diamond as gift_price,tob.num as gift_num,(tob.diamond * tob.num) as gift_total 
	from t_order_bill tob
	left join t_user tu 
	on tu.id = tob.user_id 
	where (tob.room_id = ? and tob.order_type = ? and tob.from_user_id = ? and (tob.create_time between ? and ?) )`, pd.RoomId, accountBook.ChangeStarlightRewardIncome, userId, startTime, endTime).Scan(&res)
	return
}

func (pd *PersonDao) RoomDashBoard(c *gin.Context) (res roomOwner.RoomDashBoardRes) {
	res = roomOwner.RoomDashBoardRes{
		TodayTimes:    "00:00:00",
		TodayMoneying: 0,
		WeekTimes:     "00:00:00",
		WeekMoneying:  0,
		MonthTimes:    "00:00:00",
		MonthMoneying: 0,
	}
	if pd.RoomId == "" || pd.UserId == "" {
		return
	}
	var (
		todayTimes     = 0
		weekTimes      = 0
		monthTimes     = 0
		nowTime        = time.Now()
		startMonthTime = easy.GetCurrMonthStartTime(nowTime).Format(time.DateOnly)
		endMonthTime   = easy.GetCurrMonthEndTime(nowTime).Format(time.DateOnly)
	)

	//获取本月的开播时长
	var wheatTime []*wheatTimeStruct
	coreDb.GetMasterDb().Table("t_room_wheat_time").Where("room_id = ?", pd.RoomId).Where("stat_date between ? and ?", startMonthTime, endMonthTime).Select("user_id,on_time as times,stat_date as date").Scan(&wheatTime)
	if len(wheatTime) > 0 {
		for _, v := range wheatTime {
			if easy.NowTimeIsToday(v.Date) {
				todayTimes += v.Times
			} else if easy.NowTimeIsThisWeek(v.Date) {
				weekTimes += v.Times
			} else if easy.NowTimeIsThisMonth(v.Date) {
				monthTimes += v.Times
			}
		}
		weekTimes += todayTimes
		monthTimes += weekTimes
	}
	res.TodayTimes = easy.SecondFormatString(int64(todayTimes))
	res.WeekTimes = easy.SecondFormatString(int64(weekTimes))
	res.MonthTimes = easy.SecondFormatString(int64(monthTimes))
	//获取本月的实时流水
	var orderBill []*orderResStruct
	coreDb.GetMasterDb().Table("t_order_bill").
		Where("room_id = ? and order_type = ? ", pd.RoomId, accountBook.ChangeStarlightRewardIncome).
		Where("create_time between ? and ?", startMonthTime, endMonthTime).
		Select("diamond as profit_amount,create_time as date").
		Scan(&orderBill)
	if len(orderBill) > 0 {
		for _, v := range orderBill {
			if easy.NowTimeIsToday(v.Date) {
				res.TodayMoneying += v.ProfitAmount
			} else if easy.NowTimeIsThisWeek(v.Date) {
				res.WeekMoneying += v.ProfitAmount
			} else if easy.NowTimeIsThisMonth(v.Date) {
				res.MonthMoneying += v.ProfitAmount
			}
		}
		res.WeekMoneying += res.TodayMoneying
		res.MonthMoneying += res.WeekMoneying
	}
	return
}

func (pd *PersonDao) RoomDashBoardMoneysChart(c *gin.Context) (res []*roomOwner.RoomDashBoardChartRes) {
	res = make([]*roomOwner.RoomDashBoardChartRes, 0)
	if pd.RoomId == "" || pd.UserId == "" {
		return
	}
	var (
		timeType = c.Query("timeType")
		nowTime  = time.Now()
		timeList []string
	)
	if timeType == "month" {
		timeList = append(timeList, easy.GetMonthDaysByTime(nowTime)...)
	} else {
		timeList = append(timeList, easy.GetWeekDaysByTime(nowTime)...)
	}
	for _, v := range timeList {
		res = append(res, &roomOwner.RoomDashBoardChartRes{
			Date: v,
			Data: 0,
		})
	}
	type orderStruct struct {
		StatDate     time.Time `json:"statDate" gorm:"column:stat_date"`
		ProfitAmount float64   `json:"profitAmount"`
	}
	// 获取本月的实时流水
	var orderRes []orderStruct
	coreDb.GetMasterDb().Table("t_order").
		Where("room_id = ? and order_type = ? and stat_date in ?", pd.RoomId, accountBook.ChangeDiamondRewardGift, timeList).
		Select("IFNULL(sum(pay_amount),0) as profit_amount,stat_date").
		Group("stat_date").
		Find(&orderRes)
	if len(orderRes) > 0 {
		orderResMap := map[string]orderStruct{}
		for _, v := range orderRes {
			orderResMap[v.StatDate.Format(time.DateOnly)] = v
		}
		for _, v := range res {
			if item, ok := orderResMap[v.Date]; ok {
				v.Data = int(item.ProfitAmount)
			}
		}
	}
	return
}

func (pd *PersonDao) RoomDashBoardTimesChart(c *gin.Context) (res []*roomOwner.RoomDashBoardChartRes) {
	res = make([]*roomOwner.RoomDashBoardChartRes, 0)
	if pd.RoomId == "" || pd.UserId == "" {
		return
	}
	var (
		timeType = c.Query("timeType")
		nowTime  = time.Now()
		timeList []string
	)
	if timeType == "month" {
		timeList = append(timeList, easy.GetMonthDaysByTime(nowTime)...)
	} else {
		timeList = append(timeList, easy.GetWeekDaysByTime(nowTime)...)
	}
	for _, v := range timeList {
		res = append(res, &roomOwner.RoomDashBoardChartRes{
			Date: v,
			Data: 0,
		})
	}
	type orderStruct struct {
		StatDate    time.Time `json:"statDate" gorm:"column:stat_date"`
		EnterCounts int       `json:"enterCounts"`
	}
	// 获取本月的实时流水
	var orderRes []orderStruct
	coreDb.GetMasterDb().Table("t_room_wheat_time").
		Where("room_id = ? and stat_date in ?", pd.RoomId, timeList).
		Select("IFNULL(sum(enter_count),0) as enter_counts,stat_date").
		Group("stat_date").
		Find(&orderRes)
	if len(orderRes) > 0 {
		orderResMap := map[string]orderStruct{}
		for _, v := range orderRes {
			orderResMap[v.StatDate.Format(time.DateOnly)] = v
		}
		for _, v := range res {
			if item, ok := orderResMap[v.Date]; ok {
				v.Data = int(item.EnterCounts)
			}
		}
	}
	return
}
