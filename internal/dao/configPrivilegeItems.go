package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
)

// PrivilegeItemConfigDao
// @Description: 等级特权物品
type PrivilegeItemConfigDao struct {
}

// GetLvPrivilegeItemConfig
//
//	@Description: 根据lv等级获取特权物品列表
//	@receiver p
//	@param level int -
//	@return res -
//	@return err -
func (p *PrivilegeItemConfigDao) GetLvPrivilegeItemConfig(level int) (res []model.UserLevelPrivilegeItems, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserLevelPrivilegeItems{}).Where("level=? and level_type=?", level, enum.UserLevelTypeLv).Scan(&res).Error
	return
}

// GetLevelPrivilegeItemDTO
//
//	@Description: 根据lv等级获取特权物品全部信息列表
//	@receiver p
//	@param level int -
//	@return res -
//	@return err -
func (p *PrivilegeItemConfigDao) GetLevelPrivilegeItemDTO(level, levelType int) (res []model.LevelPrivilegeItemDTO, err error) {
	err = coreDb.GetSlaveDb().Table("t_user_level_privilege_items pi").Joins("left join t_goods g on g.id=pi.goods_id").
		Joins("left join t_goods_type gt on gt.id=g.type_id").Where("pi.level=? and pi.level_type=?", level, levelType).
		Select("pi.id, pi.level, pi.goods_id, pi.level_type, pi.expire_date, pi.explain, g.name goods_name, g.type_id goods_type, g.icon goods_icon, g.type_key goods_type_key, g.animation_url goods_animation_url, g.animation_json_url goods_animation_json_url, gt.name goods_type_name").
		Scan(&res).Error
	return
}

// GetVipPrivilegeItemConfig
//
//	@Description: 根据vip等级获取特权物品类型
//	@receiver p
//	@param level int -
//	@return res -
//	@return err -
func (p *PrivilegeItemConfigDao) GetVipPrivilegeItemConfig(level int) (res []model.UserLevelPrivilegeItems, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserLevelPrivilegeItems{}).Where("level=? and level_type=?", level, enum.UserLevelTypeVip).Scan(&res).Error
	return
}
