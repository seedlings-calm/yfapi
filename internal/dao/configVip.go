package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type VipConfigDao struct {
}

// GetAllVipConfigList
//
//	@Description: 查询vip等级配置列表
//	@receiver u
//	@return res -
//	@return err -
func (u *VipConfigDao) GetAllVipConfigList() (res []*model.UserVipConfig, err error) {
	err = coreDb.GetMasterDb().Model(&model.UserVipConfig{}).Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// GetAllVipConfigMap 获取所有vip等级配置map
func (u *VipConfigDao) GetAllVipConfigMap() (res map[int]*model.UserVipConfig, err error) {
	res = make(map[int]*model.UserVipConfig)
	dataList, e := u.GetAllVipConfigList()
	if e != nil {
		err = e
		return
	}
	for _, info := range dataList {
		res[info.Level] = info
	}
	return
}

// GetVipConfigByLevel
//
//	@Description: 根据vip等级查询等级配置信息
//	@receiver u
//	@param levelId int -
//	@return res -
//	@return err -
func (u *VipConfigDao) GetVipConfigByLevel(level int) (res *model.UserVipConfig, err error) {
	res = new(model.UserVipConfig)
	err = coreDb.GetMasterDb().Model(&model.UserVipConfig{}).Where("level = ?", level).First(&res).Error
	return
}
