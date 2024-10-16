package dao

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/util/easy"
)

type UserLevelVipDao struct {
}

// 获取用户vip等级信息
func (u *UserLevelVipDao) GetUserVipLevel(userId string) (res *model.UserVipLevel, err error) {
	res = new(model.UserVipLevel)
	err = coreDb.GetMasterDb().Model(&model.UserVipLevel{}).Where("user_id = ?", userId).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = coreDb.GetMasterDb().Create(&model.UserVipLevel{
			UserId:     userId,
			Level:      1,
			CurrExp:    0,
			ExpireTime: easy.GetCurrDayEndTime(time.Now()).AddDate(0, 0, 30),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}).Error
		if err == nil {
			return u.GetUserVipLevel(userId)
		}
	}
	return
}

// 获取用户vip等级信息
func (u *UserLevelVipDao) GetUserVipLevelDTO(userId string) (res *model.UserVipLevelDTO, err error) {
	res = new(model.UserVipLevelDTO)
	err = coreDb.GetSlaveDb().Table("t_user_vip_level vl").Joins("left join t_user_vip_config vc on vc.level=vl.level").Where("vl.user_id=?", userId).
		Select("vl.*, vc.max_experience, vc.min_experience, vc.icon").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// Save 更新用户Vip等级信息
func (u *UserLevelVipDao) Save(data *model.UserVipLevel) (err error) {
	return coreDb.GetMasterDb().Save(data).Error
}
