package accountBook

import (
	"context"
	"fmt"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	"yfapi/core/coreSnowflake"
	"yfapi/internal/model"
)

type Order struct {
}

// 生成订单号
func (o *Order) OrderNum(orderType int) string {
	key := lastOrderAutoIncrNum()
	autoNum, err := coreRedis.GetUserRedis().Incr(context.Background(), key).Result()
	orderNum := coreSnowflake.GetSnowId()
	if err != nil {
		coreLog.LogError("OrderNum err:%+v", err)
		return orderNum
	}
	switch orderType {
	case ORDER_CZ:
		orderNum = fmt.Sprintf("%s%s%s", "CZ", time.Now().Format("20060102"), fmt.Sprintf("%0*d", 7, autoNum))
	case ORDER_SC:
		orderNum = fmt.Sprintf("%s%s%s", "SC", time.Now().Format("20060102"), fmt.Sprintf("%0*d", 7, autoNum))
	case ORDER_TX:
		orderNum = fmt.Sprintf("%s%s%s", "TX", time.Now().Format("20060102"), fmt.Sprintf("%0*d", 7, autoNum))
	case ORDER_PW:
		orderNum = fmt.Sprintf("%s%s%s", "PW", time.Now().Format("20060102"), fmt.Sprintf("%0*d", 7, autoNum))
	case ORDER_DS:
		orderNum = fmt.Sprintf("%s%s%s", "DS", time.Now().Format("20060102"), fmt.Sprintf("%0*d", 7, autoNum))
	default:
		orderNum = fmt.Sprintf("%s%s%s", "OR", time.Now().Format("20060102"), fmt.Sprintf("%0*d", 7, autoNum))
	}
	orderModel := &model.Order{}
	coreDb.GetMasterDb().Model(&model.Order{}).Where(&model.Order{OrderId: orderNum}).First(orderModel)
	if orderModel.ID == 0 {
		coreRedis.GetUserRedis().Expire(context.Background(), key, time.Hour*24)
	} else {
		orderNum = coreSnowflake.GetSnowId()
	}
	return orderNum
}
