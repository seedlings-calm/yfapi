package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type DaoRoomPractitioner struct {
}

func (du *DaoRoomPractitioner) Find(userId string, roomId string) (res []model.UserPractitioner, err error) {
	err = coreDb.GetMasterDb().Model(model.UserPractitioner{}).Where("room_id = ? and user_id = ? and status=1", roomId, userId).Find(&res).Error
	return
}
func (u *DaoRoomPractitioner) Create(data *model.UserPractitioner) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(&data).Error
	return
}
func (u *DaoRoomPractitioner) Update(data *model.UserPractitioner) (err error) {
	err = coreDb.GetMasterDb().Model(data).Save(&data).Error
	return
}
