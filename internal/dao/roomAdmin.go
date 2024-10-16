package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type RoomAdminDao struct {
}

// Create 添加
func (r *RoomAdminDao) Create(data *model.RoomAdmin) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// Delete 删除
func (r *RoomAdminDao) Delete(userId string, roomId string) error {
	err := coreDb.GetMasterDb().Where("room_id=? and user_id=?", roomId, userId).Delete(model.RoomAdmin{}).Error
	return err
}

// List 获取列表
func (r *RoomAdminDao) List(roomId string) (res []*model.RoomAdmin, err error) {
	err = coreDb.GetMasterDb().Model(model.RoomAdmin{}).Find(model.RoomAdmin{RoomId: roomId}).Error
	return
}

// IsAdmin 判断用户是否已是管理员
func (r *RoomAdminDao) IsAdmin(userId, roomId string) bool {
	roomAdmin := new(model.RoomAdmin)
	err := coreDb.GetMasterDb().Where(model.RoomAdmin{UserId: userId, RoomId: roomId}).First(roomAdmin).Error
	if err != nil {
		return false
	}
	return true
}
