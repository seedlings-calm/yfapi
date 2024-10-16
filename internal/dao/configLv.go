package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type LvConfigDao struct {
}

// GetLvConfigByLevel 根据lv等级获取lv等级配置
func (u *LvConfigDao) GetLvConfigByLevel(level int) (res *model.UserLvConfig, err error) {
	res = new(model.UserLvConfig)
	err = coreDb.GetMasterDb().Model(&model.UserLvConfig{}).Where("level = ?", level).First(&res).Error
	return
}

// 根据id获取权益信息
func (u *LvConfigDao) GetUserPrivilegeById(id int) (res *model.UserLevelPrivilege, err error) {
	res = new(model.UserLevelPrivilege)
	err = coreDb.GetMasterDb().Model(&model.UserLevelPrivilege{}).Where("id=?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// 获取lv等级配置信息
func (u *LvConfigDao) GetAllLvConfigList() (res []*model.UserLvConfig, err error) {
	err = coreDb.GetMasterDb().Model(&model.UserLvConfig{}).Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// GetAllLvConfigMap 获取所有lv等级配置map
func (u *LvConfigDao) GetAllLvConfigMap() (res map[int]*model.UserLvConfig, err error) {
	res = make(map[int]*model.UserLvConfig)
	dataList, e := u.GetAllLvConfigList()
	if e != nil {
		err = e
		return
	}
	for _, info := range dataList {
		res[info.Level] = info
	}
	return
}

// 获取lv最大等级配置信息
func (u *LvConfigDao) GetLevelMaxConfig(types int) (res *model.UserLevelMaxConfig, err error) {
	err = coreDb.GetMasterDb().Model(&model.UserLevelMaxConfig{}).Where("types=?", types).First(&res).Error
	return
}
