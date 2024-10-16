package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type SmsDao struct {
}

// Create 添加
func (u *SmsDao) Create(data *model.Sms) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// Update 修改
func (u *SmsDao) Update(data model.Sms) (err error) {
	err = coreDb.GetMasterDb().Model(&model.Sms{}).Where("region_code = ? and mobile = ? and types = ? and code = ? and is_use = 0", data.RegionCode, data.Mobile, data.Types, data.Code).Update("is_use", 1).Error
	return
}
