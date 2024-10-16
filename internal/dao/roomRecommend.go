package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
)

type RoomRecommendDao struct {
}

//
/**
 * @description 导航栏为你推荐接口专用
 * @param  liveType int true 如果为0 不区分房间类型
 */
func (u *RoomRecommendDao) GetRooms(nowDate, nowTime string, liveType int) (res []model.Room) {
	db := coreDb.GetMasterDb().Table("t_room_recommend as trr").
		Joins("left join t_room as tr on trr.room_id = tr.id").
		Where("trr.start_date <= ? and trr.end_date >= ?", nowDate, nowDate).
		Where("trr.start_time < ? and trr.end_time > ?", nowTime, nowTime).
		Where("trr.status = 1").
		Where("tr.status = 1")
	if liveType != 0 {
		db = db.Where("trr.live_type = ?", liveType)
	}

	db = db.Select("tr.id", "tr.user_id", "tr.room_no", "tr.room_type", "tr.live_type", "tr.template_id", "tr.cover_img", "tr.background_img", "tr.name", "tr.status", "tr.room_pwd").
		Order("trr.sort desc").Order("trr.update_time")
	var err error
	if liveType != enum.LiveTypeAnchor && liveType != 0 { //推荐的聊天室展示3条
		err = db.Limit(3).Scan(&res).Error
		if err != nil || len(res) < 3 {
			return nil
		}
	} else {
		err = db.Scan(&res).Error
		if err != nil {
			return nil
		}
	}

	return
}
