package dao

import (
	"errors"
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/util/easy"

	"gorm.io/gorm"
)

type UserLevelStarDao struct {
}

// 获取用户星光等级信息
func (u *UserLevelStarDao) GetUserStarLevel(userId string) (res *model.UserStarLevel, err error) {
	res = new(model.UserStarLevel)
	err = coreDb.GetMasterDb().Model(&model.UserStarLevel{}).Where("user_id = ?", userId).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = coreDb.GetMasterDb().Create(&model.UserStarLevel{
			UserId:     userId,
			Level:      1,
			CurrExp:    0,
			ExpireTime: easy.GetCurrDayEndTime(time.Now()).AddDate(0, 0, 30),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}).Error
		if err == nil {
			return u.GetUserStarLevel(userId)
		}
	}
	return
}

// 获取用户星光等级信息
func (u *UserLevelStarDao) GetUserStarLevelDTO(userId string) (res *model.UserStarLevelDTO, err error) {
	res = new(model.UserStarLevelDTO)
	err = coreDb.GetSlaveDb().Table("t_user_star_level sl").Joins("left join t_user_star_config sc on sc.level=sl.level").Where("sl.user_id=?", userId).
		Select("sl.*, sc.max_experience, sc.min_experience, sc.icon").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// Save 更新用户星光等级信息
func (u *UserLevelStarDao) Save(data *model.UserStarLevel) (err error) {
	return coreDb.GetMasterDb().Save(data).Error
}

func (u *UserLevelStarDao) FirstById(userId string) (res *model.UserStarLevel, err error) {
	err = coreDb.GetMasterDb().Model(&model.UserStarLevel{}).Where("user_id = ?", userId).First(&res).Error
	return
}
