package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type UserWithdrawDao struct {
}

func (u *UserWithdrawDao) UserWithdrawInfoById(id int64) (result model.AppSetting, err error) {
	err = coreDb.GetMasterDb().Model(&model.AppSetting{}).Where("id = ?", id).First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}
