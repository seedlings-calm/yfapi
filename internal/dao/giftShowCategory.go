package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type GiftShowCategory struct {
}

// GetListByLiveType
//
//	@Description: 礼物显示类目列表
//	@receiver g
//	@return result -
//	@return err -
func (g *GiftShowCategory) GetListByLiveType(liveType int) (result []*model.GiftShowCategory, err error) {
	err = coreDb.GetSlaveDb().Model(model.GiftShowCategory{}).Where("room_live_type=?", liveType).Order("sort_no").Scan(&result).Error
	return
}
