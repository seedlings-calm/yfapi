package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type RoomTemplateDao struct{}

func (RoomTemplateDao) GetBroadcastFirst(num string) (res model.RoomTemplate, err error) {
	err = coreDb.GetMasterDb().Model(model.RoomTemplate{}).Where("live_type = 2 and seat_list_count = ? and status = 1", num).First(&res).Error
	return
}
