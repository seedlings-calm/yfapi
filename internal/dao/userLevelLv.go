package dao

import (
	"errors"
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"

	"gorm.io/gorm"
)

type UserLevelLvDao struct {
}

// GetUserLvLevel 获取用户等级信息
func (u *UserLevelLvDao) GetUserLvLevel(userId string) (res *model.UserLvLevel, err error) {
	res = new(model.UserLvLevel)
	err = coreDb.GetMasterDb().Model(&model.UserLvLevel{}).Where("user_id = ?", userId).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = coreDb.GetMasterDb().Create(&model.UserLvLevel{
			UserId:     userId,
			Level:      1,
			CurrExp:    0,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}).Error
		if err == nil {
			return u.GetUserLvLevel(userId)
		}
	}
	return
}

// GetUserLvLevelDTO 获取用户等级信息
func (u *UserLevelLvDao) GetUserLvLevelDTO(userId string) (res *model.UserLvLevelDTO, err error) {
	res = new(model.UserLvLevelDTO)
	err = coreDb.GetSlaveDb().Table("t_user_lv_level vl").Joins("left join t_user_lv_config vc on vc.level=vl.level").Where("vl.user_id=?", userId).
		Select("vl.*, vc.max_experience, vc.min_experience, vc.icon").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// Save 更新用户Lv等级信息
func (u *UserLevelLvDao) Save(data *model.UserLvLevel) (err error) {
	return coreDb.GetMasterDb().Save(data).Error
}

func (u *UserLevelLvDao) GetUserIdsByLvLevel(ids []string, lvLevel int) (res []*model.UserLvLevel) {
	coreDb.GetMasterDb().Table("t_user_lv_level").Where("user_id IN ? and level >= ? ", ids, lvLevel).Find(&res)
	return res
}
