package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/request/roomOwner"
)

type RoomAdminListDao struct {
}

// GetRoomListPage 获取房间管理员列表
func (r *RoomAdminListDao) GetRoomListPage(req *roomOwner.RoomAdminListReq, roomId string) (list interface{}, count int64, err error) {
	limit := req.Size
	offset := req.Size * (req.CurrentPage - 1)
	db := coreDb.GetSlaveDb().Model(&model.RoomAdmin{}).Where(&model.RoomAdmin{RoomId: roomId})
	var dataList []model.RoomAdmin
	err = db.Count(&count).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&dataList).Error
	return dataList, count, err
}
