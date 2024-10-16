package user

import (
	"context"
	"github.com/spf13/cast"
	"strings"
	"yfapi/core/coreRedis"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/service/auth"
	"yfapi/typedef/enum"
	common_userInfo "yfapi/typedef/redisKey"
	"yfapi/typedef/response"
)

// GetUserLevelPlaque 获取玩家等级铭牌信息
func GetUserLevelPlaque(userId, clientType string, isColor ...bool) (res response.UserPlaqueInfo) {
	// 查询玩家当前在房信息
	if len(isColor) > 1 && isColor[1] {
		roomId := coreRedis.GetChatroomRedis().Get(context.Background(), common_userInfo.UserInWhichRoom(userId, clientType)).Val()
		if len(roomId) > 0 {
			// 玩家在此房间的身份
			_, roleIdList, _ := new(auth.Auth).GetRoleListByRoomIdAndUserId(roomId, userId)
			if len(roleIdList) > 0 {
				for _, info := range roleIdList {
					if len(info.Icon) == 0 {
						continue
					}
					res.HeadList = append(res.HeadList, response.PlaqueInfo{
						PlaqueType: enum.RoleType(info.RoleId).PlaqueTypeName(),
						Content:    "",
						Icon:       helper.FormatImgUrl(info.Icon),
					})
				}
			}
		}
	}

	// 获取用户星光等级铭牌信息
	userStar, err := new(dao.UserLevelStarDao).GetUserStarLevelDTO(userId)
	if err != nil {
		return
	}
	if userStar.ID > 0 {
		res.HeadList = append(res.HeadList, response.PlaqueInfo{
			PlaqueType: "star",
			Content:    cast.ToString(userStar.Level),
			Icon:       helper.FormatImgUrl(userStar.Icon),
		})
	}
	// 获取用户vip等级铭牌信息
	userVip, err := new(dao.UserLevelVipDao).GetUserVipLevelDTO(userId)
	if err != nil {
		return
	}
	res.TailList = append(res.TailList, response.PlaqueInfo{
		PlaqueType: "vip",
		Content:    cast.ToString(userVip.Level),
		Icon:       helper.FormatImgUrl(userVip.Icon),
	})
	// 获取用户lv等级铭牌信息
	userLv, err := new(dao.UserLevelLvDao).GetUserLvLevelDTO(userId)
	if err != nil {
		return
	}
	res.TailList = append(res.TailList, response.PlaqueInfo{
		PlaqueType: "lv",
		Content:    cast.ToString(userLv.Level),
		Icon:       helper.FormatImgUrl(userLv.Icon),
	})
	// 是否解锁彩色昵称特权
	if len(isColor) > 0 && isColor[0] {
		privilegeConfig, _ := new(dao.PrivilegeConfigDao).GetPrivilegeById(8)
		if privilegeConfig.ID > 0 {
			if privilegeConfig.MinLv > 0 && userLv.Level >= privilegeConfig.MinLv {
				res.ColorList = strings.Split(privilegeConfig.ColorList, ",")
			} else if privilegeConfig.MinVip > 0 && userVip.Level >= privilegeConfig.MinVip {
				res.ColorList = strings.Split(privilegeConfig.ColorList, ",")
			} else if privilegeConfig.MinStar > 0 && userStar.Level >= privilegeConfig.MinStar {
				res.ColorList = strings.Split(privilegeConfig.ColorList, ",")
			}
		}
		privilegeConfig, _ = new(dao.PrivilegeConfigDao).GetPrivilegeById(9)
		if privilegeConfig.ID > 0 {
			if privilegeConfig.MinLv > 0 && userLv.Level >= privilegeConfig.MinLv {
				res.ChatColor = strings.Split(privilegeConfig.ColorList, ",")
			} else if privilegeConfig.MinVip > 0 && userVip.Level >= privilegeConfig.MinVip {
				res.ChatColor = strings.Split(privilegeConfig.ColorList, ",")
			} else if privilegeConfig.MinStar > 0 && userStar.Level >= privilegeConfig.MinStar {
				res.ChatColor = strings.Split(privilegeConfig.ColorList, ",")
			}
		}
	}
	return
}
