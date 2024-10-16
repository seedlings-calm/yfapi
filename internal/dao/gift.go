package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type GiftDao struct {
}

func (g *GiftDao) GetRoomGiftList(liveType, roomType, categoryType int, nowTime string) (result []model.GiftDTO, err error) {
	err = coreDb.GetSlaveDb().Table("t_gift_room gr").Joins("left join t_gift g on g.gift_code=gr.gift_code").Where("gr.room_live_type=? and gr.category_type=? and gr.del_status=0", liveType, categoryType).
		Where("gr.start_time<=? and gr.end_time>=?", nowTime, nowTime).Where("FIND_IN_SET(?, gr.room_type)", roomType).Order("gr.sort_no").
		Select("gr.*, g.gift_name, g.gift_image, g.gift_grade, g.animation_url, g.animation_json_url, g.gift_amount_type, g.gift_diamond, g.gift_revenue_type, g.exp_times, g.send_count_list").
		Scan(&result).Error
	return
}

func (g *GiftDao) GetGiftByCode(giftCode string) (result model.GiftDTO, err error) {
	err = coreDb.GetSlaveDb().Table("t_gift_room gr").Joins("left join t_gift g on g.gift_code=gr.gift_code").Where("gr.gift_code=? and gr.del_status=0", giftCode).
		Select("gr.*, g.gift_name, g.gift_image, g.gift_grade, g.animation_url, g.animation_json_url, g.gift_amount_type, g.gift_diamond, g.gift_revenue_type, g.exp_times, g.send_count_list").Scan(&result).Error
	return
}

func (g *GiftDao) GetGiftSourceList() (result []model.GiftDTO, err error) {
	nowTime := time.Now().Format(time.DateTime)
	err = coreDb.GetSlaveDb().Table("t_gift_room gr").Joins("left join t_gift g on g.gift_code=gr.gift_code").
		Where("gr.start_time<=? and gr.end_time>=? and gr.del_status=0", nowTime, nowTime).
		Select("g.gift_code, g.animation_url, g.animation_json_url").
		Scan(&result).Error
	return
}
