package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	"yfapi/util/easy"

	"github.com/spf13/cast"
)

type OrderBillDao struct {
}

// GetUserDiamondBillList 查询玩家钻石流水
func (o OrderBillDao) GetUserDiamondBillList(userId string, page, size, fundFlow int, timeKey string, orderType ...int) (result []model.OrderBillDTO, count int64) {
	startTime := ""
	endTime := ""
	if len(timeKey) > 0 {
		searchTime, _ := time.ParseInLocation("2006-01", timeKey, time.Local)
		startTime = easy.GetCurrMonthStartTime(searchTime).AddDate(-2, 0, 0).Format(time.DateTime)
		endTime = easy.GetCurrMonthEndTime(searchTime).Format(time.DateTime)
	}
	tx := coreDb.GetMasterDb().Table("t_order_bill ob").Joins("left join t_room r on r.id=ob.room_id").Where("ob.user_id=? and ob.currency=?", userId, accountBook.CURRENCY_DIAMOND)
	if len(startTime) > 0 {
		tx = tx.Where("ob.create_time between ? and ?", startTime, endTime)
	}
	if fundFlow > 0 {
		tx = tx.Where("ob.fund_flow=?", fundFlow)
	}
	if len(orderType) > 0 {
		tx = tx.Where("ob.order_type in ?", orderType)
	}
	tx.Count(&count)
	tx.Select("ob.*, r.name room_name").Order("ob.create_time desc").Limit(size).Offset(page * size).Scan(&result)
	return
}

// GetUserStarlightBillList 查询玩家星光流水
func (o OrderBillDao) GetUserStarlightBillList(userId string, page, size, fundFlow int, timeKey string) (result []model.OrderBillDTO, count int64) {
	startTime := ""
	endTime := ""
	if len(timeKey) > 0 {
		searchTime, _ := time.ParseInLocation("2006-01", timeKey, time.Local)
		startTime = easy.GetCurrMonthStartTime(searchTime).AddDate(-2, 0, 0).Format(time.DateTime)
		endTime = easy.GetCurrMonthEndTime(searchTime).Format(time.DateTime)
	}
	currencyList := []string{cast.ToString(accountBook.CURRENCY_STARLIGHT_UNWITHDRAW), cast.ToString(accountBook.CURRENCY_STARLIGHT_WITHDRAW), cast.ToString(accountBook.CURRENCY_STARLIGHT_SUBSIDY)}
	tx := coreDb.GetMasterDb().Table("t_order_bill ob").Joins("left join t_room r on r.id=ob.room_id").Where("ob.user_id=?", userId).Where("ob.currency in " + easy.GetInSql(currencyList))
	if len(startTime) > 0 {
		tx = tx.Where("ob.create_time between ? and ?", startTime, endTime)
	}
	if fundFlow > 0 {
		tx = tx.Where("ob.fund_flow=?", fundFlow)
	}
	tx.Count(&count)
	tx.Select("ob.*, r.name room_name").Order("ob.create_time desc").Limit(size).Offset(page * size).Scan(&result)
	return
}
