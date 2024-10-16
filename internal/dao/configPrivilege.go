package dao

import (
	"errors"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
)

// PrivilegeConfigDao
// @Description: 等级特权权益
type PrivilegeConfigDao struct {
}

// GetLvPrivilegeConfig
//
//	@Description: 根据lv等级查询特权权益列表
//	@receiver p
//	@param level int -
//	@return res -
//	@return err -
func (p *PrivilegeConfigDao) GetLvPrivilegeConfig(level int) (res []model.UserLevelPrivilege, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserLevelPrivilege{}).Where("min_lv>? and min_lv<=?", 0, level).Order("min_lv,id").Scan(&res).Error
	return
}

// GetVipPrivilegeConfig
//
//	@Description: 根据vip等级查询特权权益列表
//	@receiver p
//	@param level int -
//	@return res -
//	@return err -
func (p *PrivilegeConfigDao) GetVipPrivilegeConfig(level int) (res []model.UserLevelPrivilege, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserLevelPrivilege{}).Where("min_vip>? and min_vip<=?", 0, level).Order("min_vip,id").Scan(&res).Error
	return
}

// GetStarlightPrivilegeConfig
//
//	@Description: 根据星光等级查询特权权益列表
//	@receiver p
//	@param level int -
//	@return res -
//	@return err -
func (p *PrivilegeConfigDao) GetStarlightPrivilegeConfig(level int) (res []model.UserLevelPrivilege, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserLevelPrivilege{}).Where("min_star>? and min_star<=?", 0, level).Order("min_star,id").Scan(&res).Error
	return
}

// GetAllPrivilegeList
//
//	@Description: 根据等级类型查询所有等级权益列表
//	@receiver p
//	@param levelType int -
//	@return res -
//	@return err -
func (p *PrivilegeConfigDao) GetAllPrivilegeList(levelType int) (res []model.UserLevelPrivilege, err error) {
	searchSql := ""
	switch levelType {
	case enum.UserLevelTypeLv:
		searchSql = "min_lv"
	case enum.UserLevelTypeVip:
		searchSql = "min_vip"
	case enum.UserLevelTypeStarlight:
		searchSql = "min_star"
	default:
		err = errors.New("unknown levelType")
		return
	}
	err = coreDb.GetSlaveDb().Model(model.UserLevelPrivilege{}).Where(searchSql + ">0").Order(searchSql + ",id").Scan(&res).Error
	return
}

func (p *PrivilegeConfigDao) GetPrivilegeById(id int) (res model.UserLevelPrivilege, err error) {
	err = coreDb.GetSlaveDb().Model(res).Where("id", id).Scan(&res).Error
	return
}
