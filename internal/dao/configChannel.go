package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type ConfigChannelDao struct {
}

func (ConfigChannelDao) First(platform, channel string) (res model.ConfigChannel, err error) {
	err = coreDb.GetMasterDb().Model(model.ConfigChannel{}).Where("platform = ? and channel = ? and status = 1", platform, channel).First(&res).Error
	return
}
