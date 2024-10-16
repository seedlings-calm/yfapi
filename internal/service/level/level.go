package level

import (
	"fmt"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
	"yfapi/typedef/response/h5"
)

// GetLvConfigWithPrivilege
//
//	@Description: 根据lv等级查询当前等级的所有配置包含特权权益、物品
//	@param level int -
//	@return res -
//	@return err -
func GetLvConfigWithPrivilege(level int) (res h5.LevelConfig, err error) {
	// 最大等级配置
	lvConfigDao := new(dao.LvConfigDao)
	maxLevel, err := lvConfigDao.GetLevelMaxConfig(enum.UserLevelTypeLv)
	if err != nil {
		return
	}
	// 超过最大等级
	if level > maxLevel.MaxLevel {
		return
	}
	// 查询当前等级配置
	config, err := lvConfigDao.GetLvConfigByLevel(level)
	if err != nil {
		return
	}
	res = h5.LevelConfig{
		LevelName: config.LevelName,
		Level:     config.Level,
		Icon:      helper.FormatImgUrl(config.Icon),
		LogoIcon:  helper.FormatImgUrl(config.LogoIcon),
		MinExp:    config.MinExperience,
		MaxExp:    config.MaxExperience,
		MaxLevel:  maxLevel.MaxLevel,
	}
	// 当前等级特权物品
	res.PrivilegeItemList, err = getPrivilegeItemList(level, enum.UserLevelTypeLv)
	if err != nil {
		return
	}
	// 当前等级权益
	res.PrivilegeList, err = getPrivilegeList(level, enum.UserLevelTypeLv)
	if err != nil {
		return
	}
	// 当前等级权益数量
	res.PrivilegeCount = len(res.PrivilegeList)
	return
}

// GetVipConfigWithPrivilege
//
//	@Description: 根据vip等级查询当前等级的所有配置包含特权权益、物品
//	@param level int -
//	@return res -
//	@return err -
func GetVipConfigWithPrivilege(level int) (res h5.LevelConfig, err error) {
	// 最大等级配置
	lvConfigDao := new(dao.LvConfigDao)
	maxLevel, err := lvConfigDao.GetLevelMaxConfig(enum.UserLevelTypeVip)
	if err != nil {
		return
	}
	// 超过最大等级
	if level > maxLevel.MaxLevel {
		return
	}
	// 查询当前等级配置
	config, err := new(dao.VipConfigDao).GetVipConfigByLevel(level)
	if err != nil {
		return
	}
	res = h5.LevelConfig{
		LevelName: config.LevelName,
		Level:     config.Level,
		Icon:      helper.FormatImgUrl(config.Icon),
		LogoIcon:  helper.FormatImgUrl(config.LogoIcon),
		MinExp:    config.MinExperience,
		MaxExp:    config.MaxExperience,
		MaxLevel:  maxLevel.MaxLevel,
	}
	// 当前等级特权物品
	res.PrivilegeItemList, err = getPrivilegeItemList(level, enum.UserLevelTypeVip)
	if err != nil {
		return
	}
	// 当前等级权益
	res.PrivilegeList, err = getPrivilegeList(level, enum.UserLevelTypeVip)
	if err != nil {
		return
	}
	// 当前等级权益数量
	res.PrivilegeCount = len(res.PrivilegeList)
	return
}

// GetStarConfigWithPrivilege
//
//	@Description: 根据星光等级查询当前等级的所有配置包含特权权益、物品
//	@param level int -
//	@return res -
//	@return err -
func GetStarConfigWithPrivilege(level int) (res h5.LevelConfig, err error) {
	// 最大等级配置
	lvConfigDao := new(dao.LvConfigDao)
	maxLevel, err := lvConfigDao.GetLevelMaxConfig(enum.UserLevelTypeStarlight)
	if err != nil {
		return
	}
	// 超过最大等级
	if level > maxLevel.MaxLevel {
		return
	}
	// 查询当前等级配置
	config, err := new(dao.StarConfigDao).GetStarConfigByLevel(level)
	if err != nil {
		return
	}
	res = h5.LevelConfig{
		LevelName: config.LevelName,
		Level:     config.Level,
		Icon:      helper.FormatImgUrl(config.Icon),
		LogoIcon:  helper.FormatImgUrl(config.LogoIcon),
		MinExp:    config.MinExperience,
		MaxExp:    config.MaxExperience,
		MaxLevel:  maxLevel.MaxLevel,
	}
	// 当前等级特权物品
	res.PrivilegeItemList, err = getPrivilegeItemList(level, enum.UserLevelTypeStarlight)
	if err != nil {
		return
	}
	// 当前等级权益
	res.PrivilegeList, err = getPrivilegeList(level, enum.UserLevelTypeStarlight)
	if err != nil {
		return
	}
	// 当前等级权益数量
	res.PrivilegeCount = len(res.PrivilegeList)
	return
}

// 查询等级特权权益物品列表
func getPrivilegeItemList(level, levelType int) (res []*h5.PrivilegeItemConfig, err error) {
	// 当前等级特权物品
	itemList, err := new(dao.PrivilegeItemConfigDao).GetLevelPrivilegeItemDTO(level, levelType)
	if err != nil {
		return res, err
	}
	for _, info := range itemList {
		expireDate := info.ExpireDate
		if expireDate == 36500 {
			expireDate = -1
		}
		res = append(res, &h5.PrivilegeItemConfig{
			GoodsId:          info.GoodsId,
			Name:             info.GoodsName,
			Icon:             helper.FormatImgUrl(info.GoodsIcon),
			AnimationUrl:     helper.FormatImgUrl(info.GoodsAnimationUrl),
			AnimationJsonUrl: helper.FormatImgUrl(info.GoodsAnimationJsonUrl),
			TypeName:         info.GoodsTypeName,
			TypeKey:          info.GoodsTypeKey,
			Explain:          info.Explain,
			ExpirationDate:   expireDate,
		})
	}
	return
}

// 查询等级特权权益列表
func getPrivilegeList(level, levelType int) (res []*h5.PrivilegeConfig, err error) {
	// 当前等级权益
	var currPrivilegeList []model.UserLevelPrivilege
	switch levelType {
	case enum.UserLevelTypeLv:
		currPrivilegeList, err = new(dao.PrivilegeConfigDao).GetLvPrivilegeConfig(level)
	case enum.UserLevelTypeVip:
		currPrivilegeList, err = new(dao.PrivilegeConfigDao).GetVipPrivilegeConfig(level)
	case enum.UserLevelTypeStarlight:
		currPrivilegeList, err = new(dao.PrivilegeConfigDao).GetStarlightPrivilegeConfig(level)
	default:
		err = fmt.Errorf("getPrivilegeList unknown level type:%v", levelType)
	}
	if err != nil {
		return
	}
	for _, info := range currPrivilegeList {
		res = append(res, &h5.PrivilegeConfig{
			Name:        info.Name,
			Icon:        helper.FormatImgUrl(info.Icon),
			LightEffect: helper.FormatImgUrl(info.LightEffect),
			MinLv:       info.MinLv,
			MinVip:      info.MinVip,
			MinStar:     info.MinStar,
			Explain:     info.Explain,
		})
	}
	return
}
