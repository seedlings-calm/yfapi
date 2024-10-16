package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type StarConfigDao struct {
}

// GetAllStarConfigList
//
//	@Description: 查询星光等级配置列表
//	@receiver u
//	@return res -
//	@return err -
func (u *StarConfigDao) GetAllStarConfigList() (res []*model.UserStarConfig, err error) {
	err = coreDb.GetMasterDb().Model(&model.UserStarConfig{}).Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// GetAllStarConfigMap 获取所有星光等级配置map
func (u *StarConfigDao) GetAllStarConfigMap() (res map[int]*model.UserStarConfig, err error) {
	res = make(map[int]*model.UserStarConfig)
	dataList, e := u.GetAllStarConfigList()
	if e != nil {
		err = e
		return
	}
	for _, info := range dataList {
		res[info.Level] = info
	}
	return
}

// GetStarConfigByLevel
//
//	@Description: 根据星光等级查询等级配置信息
//	@receiver u
//	@param levelId int -
//	@return res -
//	@return err -
func (u *StarConfigDao) GetStarConfigByLevel(level int) (res *model.UserStarConfig, err error) {
	res = new(model.UserStarConfig)
	err = coreDb.GetMasterDb().Model(&model.UserStarConfig{}).Where("level = ?", level).First(&res).Error
	return
}
