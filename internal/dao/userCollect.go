package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type DaoUserCollect struct {
	UserId string
}

// 如果where条件查不到，使用where+attrs属性创建， 查到不操作插入
func (d *DaoUserCollect) Create(roomId string) error {
	err := coreDb.GetMasterDb().Model(model.UserCollect{}).
		Where(&model.UserCollect{
			UserID:    d.UserId,
			RoomId:    roomId,
			IsCollect: 1,
		}).Attrs(&model.UserCollect{
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}).FirstOrCreate(&model.UserCollect{}).Error
	return err
}

// 修改关注为未关注
func (d *DaoUserCollect) Update(roomId string) error {
	err := coreDb.GetMasterDb().Model(model.UserCollect{}).
		Where("user_id = ?  and is_collect = 1 and room_id = ?", d.UserId, roomId).
		Updates(map[string]interface{}{
			"is_collect":  2,
			"update_time": time.Now(),
		}).Error
	return err
}

// 获取用户收藏的所有房间ID
func (d *DaoUserCollect) GetRoomIds() (res []string, err error) {
	err = coreDb.GetMasterDb().Model(model.UserCollect{}).Where("user_id = ? and is_collect = 1", d.UserId).Pluck("room_id", &res).Error
	return
}

func (d *DaoUserCollect) IsRoomCollect(userId, roomId string) (isCollect bool) {
	var count int64
	err := coreDb.GetMasterDb().Model(model.UserCollect{}).Where("user_id=? and room_id=? and is_collect=1", userId, roomId).Count(&count).Error
	if err != nil {
		return false
	}
	isCollect = count > 0
	return
}
