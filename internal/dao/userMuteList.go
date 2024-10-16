package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	response_room "yfapi/typedef/response/room"
)

type UserMuteListDao struct {
}

func (u UserMuteListDao) Create(param *model.UserMuteList) (err error) {
	err = coreDb.GetMasterDb().Create(param).Error
	return err
}

func (u UserMuteListDao) Delete(userId, roomId string) (err error) {
	err = coreDb.GetMasterDb().Where("room_id = ? and to_id = ?", roomId, userId).Delete(&model.UserMuteList{}).Error
	return
}

func (u UserMuteListDao) FindOne(userId, roomId string) (res model.UserMuteList) {
	coreDb.GetMasterDb().Where("room_id = ? and to_id = ? and end_time > ?", roomId, userId, time.Now()).First(&res)
	return
}

func (u UserMuteListDao) List(roomId string) (res []response_room.UserMuteListRes) {
	coreDb.GetMasterDb().Table("t_user_mutelist as a").
		Joins("inner join t_user as b on a.to_id = b.id").
		Where("a.room_id = ? and a.end_time > ?", roomId, time.Now()).
		Select("a.to_id as user_id,b.user_no,b.avatar,b.nickname").Order("a.start_time").Scan(&res)
	return
}
