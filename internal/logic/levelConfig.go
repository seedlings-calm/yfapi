package logic

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	service_level "yfapi/internal/service/level"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	request_h5 "yfapi/typedef/request/h5"
	response_h5 "yfapi/typedef/response/h5"
)

type LevelConfig struct {
}

// GetLevelBaseInfo
//
//	@Description: 等级基本信息
//	@receiver g
//	@param c
//	@return res
func (g *LevelConfig) GetLevelBaseInfo(c *gin.Context, req *request_h5.LevelBaseInfoReq) (res *response_h5.LevelBaseInfoRes) {
	res = new(response_h5.LevelBaseInfoRes)
	userId := handle.GetUserId(c)
	user := service_user.GetUserBaseInfo(userId)
	if len(user.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	res.Nickname = user.Nickname
	res.Avatar = user.Avatar
	// 查询用户当前lv等级
	lvDao := new(dao.UserLevelLvDao)
	userLv, err := lvDao.GetUserLvLevel(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 当前等级、权益、物品配置信息
	res.LvBase.LevelType = enum.UserLevelTypeLv
	res.LvBase.CurrLevel, err = service_level.GetLvConfigWithPrivilege(userLv.Level)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	res.LvBase.CurrLevel.CurrExp = userLv.CurrExp
	// 下一级等级、权益、物品配置信息
	res.LvBase.NextLevel, _ = getShowNextLevelConfig(userLv.Level, res.LvBase.CurrLevel.MaxLevel, enum.UserLevelTypeLv)
	// 全部等级权益
	allPrivilegeList, err := new(dao.PrivilegeConfigDao).GetAllPrivilegeList(enum.UserLevelTypeLv)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 当前等级的权益列表显示所有的权益信息
	res.LvBase.CurrLevel.PrivilegeList = []*response_h5.PrivilegeConfig{}
	for _, info := range allPrivilegeList {
		res.LvBase.CurrLevel.PrivilegeList = append(res.LvBase.CurrLevel.PrivilegeList, &response_h5.PrivilegeConfig{
			Name:        info.Name,
			Icon:        helper.FormatImgUrl(info.Icon),
			LightEffect: helper.FormatImgUrl(info.LightEffect),
			MinLv:       info.MinLv,
			MinVip:      info.MinVip,
			MinStar:     info.MinStar,
			Explain:     info.Explain,
		})
	}

	// 查询用户当前vip等级
	vipDao := new(dao.UserLevelVipDao)
	userVip, err := vipDao.GetUserVipLevel(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 当前等级、权益、物品配置信息
	res.VipBase.LevelType = enum.UserLevelTypeVip
	res.VipBase.CurrLevel, err = service_level.GetVipConfigWithPrivilege(userVip.Level)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	res.VipBase.CurrLevel.CurrExp = userVip.CurrExp
	// 下一级等级、权益、物品配置信息
	res.VipBase.NextLevel, _ = getShowNextLevelConfig(userVip.Level, res.VipBase.CurrLevel.MaxLevel, enum.UserLevelTypeVip)
	// 全部等级权益
	allPrivilegeList, err = new(dao.PrivilegeConfigDao).GetAllPrivilegeList(enum.UserLevelTypeVip)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 当前等级的权益列表显示所有的权益信息
	res.VipBase.CurrLevel.PrivilegeList = []*response_h5.PrivilegeConfig{}
	for _, info := range allPrivilegeList {
		res.VipBase.CurrLevel.PrivilegeList = append(res.VipBase.CurrLevel.PrivilegeList, &response_h5.PrivilegeConfig{
			Name:        info.Name,
			Icon:        helper.FormatImgUrl(info.Icon),
			LightEffect: helper.FormatImgUrl(info.LightEffect),
			MinLv:       info.MinLv,
			MinVip:      info.MinVip,
			MinStar:     info.MinStar,
			Explain:     info.Explain,
		})
	}

	// 查询用户当前星光等级
	starDao := new(dao.UserLevelStarDao)
	userStarlight, err := starDao.GetUserStarLevelDTO(userId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if userStarlight.ID > 0 {
		//当前等级、权益、物品配置信息
		res.StarBase = new(response_h5.LevelBase)
		res.StarBase.LevelType = enum.UserLevelTypeStarlight
		res.StarBase.CurrLevel, err = service_level.GetStarConfigWithPrivilege(userStarlight.Level)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		res.StarBase.CurrLevel.CurrExp = userStarlight.CurrExp
		// 下一级等级、权益、物品配置信息
		res.StarBase.NextLevel, _ = getShowNextLevelConfig(userStarlight.Level, res.StarBase.CurrLevel.MaxLevel, enum.UserLevelTypeStarlight)
		// 全部等级权益
		allPrivilegeList, err = new(dao.PrivilegeConfigDao).GetAllPrivilegeList(enum.UserLevelTypeStarlight)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		// 当前等级的权益列表显示所有的权益信息
		res.StarBase.CurrLevel.PrivilegeList = []*response_h5.PrivilegeConfig{}
		for _, info := range allPrivilegeList {
			res.StarBase.CurrLevel.PrivilegeList = append(res.StarBase.CurrLevel.PrivilegeList, &response_h5.PrivilegeConfig{
				Name:        info.Name,
				Icon:        helper.FormatImgUrl(info.Icon),
				LightEffect: helper.FormatImgUrl(info.LightEffect),
				MinLv:       info.MinLv,
				MinVip:      info.MinVip,
				MinStar:     info.MinStar,
				Explain:     info.Explain,
			})
		}
	}
	return
}

// 查询显示下一级的特权物品，如果下一级没有，就自增查询下下级，直到有值为止
func getShowNextLevelConfig(level, maxLevel, levelType int) (res response_h5.LevelConfig, err error) {
	if level <= maxLevel {
		isNotFind := true
		for isNotFind {
			level++
			if level <= maxLevel {
				var newPrivilegeList []*response_h5.PrivilegeConfig
				switch levelType {
				case enum.UserLevelTypeLv: // lv
					res, err = service_level.GetLvConfigWithPrivilege(level)
					if err != nil {
						return
					}
					// 过滤已解锁的权益
					for _, info := range res.PrivilegeList {
						if info.MinLv < level {
							continue
						}
						newPrivilegeList = append(newPrivilegeList, info)
					}
				case enum.UserLevelTypeVip: // vip
					res, err = service_level.GetVipConfigWithPrivilege(level)
					if err != nil {
						return
					}
					// 过滤已解锁的权益
					for _, info := range res.PrivilegeList {
						if info.MinVip < level {
							continue
						}
						newPrivilegeList = append(newPrivilegeList, info)
					}
				case enum.UserLevelTypeStarlight: // 星光
					res, err = service_level.GetStarConfigWithPrivilege(level)
					if err != nil {
						return
					}
					// 过滤已解锁的权益
					for _, info := range res.PrivilegeList {
						if info.MinStar < level {
							continue
						}
						newPrivilegeList = append(newPrivilegeList, info)
					}
				default:
					err = fmt.Errorf("getShowNextLevelConfig unknown level type:%v", levelType)
				}
				res.PrivilegeList = newPrivilegeList
				if len(res.PrivilegeItemList) > 0 || len(res.PrivilegeList) > 0 || res.Level >= res.MaxLevel {
					isNotFind = false
				}
			} else {
				isNotFind = false
			}
		}
	}
	return
}

// GetLevelConfigList
//
//	@Description: 获取等级配置列表
//	@receiver g
//	@param c *gin.Context -
//	@param req *request_h5.LevelBaseInfoReq -
//	@return res -
func (g *LevelConfig) GetLevelConfigList(c *gin.Context, req *request_h5.LevelBaseInfoReq) (res response_h5.LevelConfigListRes) {
	lvList, err := new(dao.LvConfigDao).GetAllLvConfigList()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFail,
			Msg:  nil,
		})
	}
	for _, info := range lvList {
		res.LvList = append(res.LvList, response_h5.LevelInfo{
			LevelName: info.LevelName,
			Icon:      helper.FormatImgUrl(info.Icon),
			MinExp:    info.MinExperience,
		})
	}
	starList, err := new(dao.StarConfigDao).GetAllStarConfigList()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFail,
			Msg:  nil,
		})
	}
	for _, info := range starList {
		res.StarList = append(res.StarList, response_h5.LevelInfo{
			LevelName: info.LevelName,
			Icon:      helper.FormatImgUrl(info.Icon),
			MinExp:    info.MinExperience,
			KeepExp:   info.RelegationExperience,
		})
	}
	vipList, err := new(dao.VipConfigDao).GetAllVipConfigList()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFail,
			Msg:  nil,
		})
	}
	for _, info := range vipList {
		res.VipList = append(res.VipList, response_h5.LevelInfo{
			LevelName: info.LevelName,
			Icon:      helper.FormatImgUrl(info.Icon),
			MinExp:    info.MinExperience,
			KeepExp:   info.RelegationExperience,
		})
	}
	return
}
