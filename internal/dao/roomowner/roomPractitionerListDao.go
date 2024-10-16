package dao

import (
	"github.com/gin-gonic/gin"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/request/roomOwner"
)

type RoomPractitionerListDao struct {
}

// GetRoomPractitionerListPage 获取从业者列表
func (r *RoomPractitionerListDao) GetRoomPractitionerListPage(req *roomOwner.RoomPractitionerListReq, c *gin.Context) (list interface{}, count int64, err error) {
	limit := req.Size
	offset := req.Size * (req.CurrentPage - 1)
	db := coreDb.GetSlaveDb().Model(&model.UserPractitioner{})
	var dataList []model.UserPractitioner
	rid := c.GetHeader("roomId")
	err = db.Where(&model.UserPractitioner{RoomId: rid}).Count(&count).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Where(&model.UserPractitioner{RoomId: rid}).Find(&dataList).Error
	return dataList, count, err
}

// 作废和移除从业者
func (r *RoomPractitionerListDao) GetRoomPractitionerRemove(params *model.UserPractitioner) (err error) {
	err = coreDb.GetMasterDb().Save(params).Error
	return
}
