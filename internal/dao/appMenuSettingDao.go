package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

// AppMenuSettingDao
// @Description: app菜单配置
type AppMenuSettingDao struct {
}

// GetAppMenuList
//
//	@Description: 获取app菜单列表
//	@receiver a
//	@param moduleType int - 模块类型
//	@param platform string - 平台
//	@return result -
//	@return err -
func (a *AppMenuSettingDao) GetAppMenuList(moduleType int, platform string) (result []model.AppMenuSetting, err error) {
	err = coreDb.GetSlaveDb().Model(model.AppMenuSetting{}).Where("module_type=? and status=1 and FIND_IN_SET(?,platform)", moduleType, platform).Order("sort_no asc,create_time desc").Scan(&result).Error
	return
}
