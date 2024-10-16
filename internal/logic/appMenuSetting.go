package logic

import (
	"github.com/gin-gonic/gin"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/typedef/enum"
	request_index "yfapi/typedef/request/index"
	"yfapi/typedef/response/index"
)

// AppMenuSetting
// @Description: app菜单配置
type AppMenuSetting struct {
}

// GetAppMenuSettingList
//
//	@Description: 获取app菜单列表
//	@receiver a
//	@param c *gin.Context -
//	@param moduleType int -
//	@return res -
func (a *AppMenuSetting) GetAppMenuSettingList(c *gin.Context, req *request_index.AppMenuSettingReq) (res []index.AppMenuSetting) {
	headerData := helper.GetHeaderData(c)
	platform := headerData.Platform
	userId := helper.GetUserId(c)
	dataList, err := new(dao.AppMenuSettingDao).GetAppMenuList(req.ModuleType, platform)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// TODO 菜单根据个人信息过滤
	// 当前玩家是否有主播资质
	isAnchor := false
	if req.ModuleType == enum.ModuleTypeUserCenter {
		userCerdDao := &dao.DaoUserPractitionerCerd{
			UserId: userId,
		}
		result, _ := userCerdDao.First(enum.UserPractitionerAnchor)
		if result.Id > 0 {
			isAnchor = true
		}
	}

	for _, info := range dataList {
		if info.MenuType == enum.MenuTypeAnchorRoom && !isAnchor {
			// 没有主播资质 不显示我的直播间
			continue
		}
		if info.MenuType == enum.MenuTypeFans && !isAnchor {
			// 没有主播资质 不显示我的粉丝团
			continue
		}
		res = append(res, index.AppMenuSetting{
			MenuName: info.MenuName,
			Icon:     helper.FormatImgUrl(info.Icon),
			LinkUrl:  info.LinkUrl,
		})
	}
	return
}
