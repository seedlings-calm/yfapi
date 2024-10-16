package logic

import (
	"github.com/gin-gonic/gin"
	"time"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/auth"
	service_im "yfapi/internal/service/im"
	"yfapi/typedef/enum"
	group2 "yfapi/typedef/request/group"
	"yfapi/typedef/response/group"
)

type Group struct {
}

// 获取世界频道设置
func (g *Group) GetWorldGroupSetting(c *gin.Context) (res group.WorldGroupSettingResp) {
	userId := helper.GetUserId(c)
	res.RoomId = enum.WorldGroupId
	res.NoticeStatus = enum.GroupReceiveAllMsg
	daoSer := new(dao.GroupFilterDao)
	muteRes := daoSer.FindOne(&model.GroupFilter{
		GroupID: enum.WorldGroupId,
		UserID:  "0",
		Types:   enum.GroupMuteSwitch,
	})
	if muteRes.ID > 0 {
		res.MuteSwitch = true
	}
	noticeRes := daoSer.FindOne(&model.GroupFilter{
		GroupID: enum.WorldGroupId,
		UserID:  userId,
	})
	if noticeRes.ID > 0 {
		res.NoticeStatus = noticeRes.Types
	}
	return
}

// 世界频道禁言
func (g *Group) WorldGroupMute(c *gin.Context) {
	userId := helper.GetUserId(c)
	ser := new(auth.Auth)
	if !ser.IsSuperAdminRole(enum.WorldGroupId, userId) {
		panic(i18n_err.I18nError{Code: i18n_err.ErrorCodeUserNoPermissions})

	}
	daoSer := new(dao.GroupFilterDao)
	muteRes := daoSer.FindOne(&model.GroupFilter{
		GroupID: enum.WorldGroupId,
		UserID:  "0",
		Types:   enum.GroupMuteSwitch,
	})
	if muteRes.ID > 0 {
		err := daoSer.Del(muteRes)
		if err != nil {
			panic(i18n_err.I18nError{Code: i18n_err.ErrorCodeSystemBusy})
		}
		new(service_im.ImPublicService).SendCustomMsg(enum.WorldGroupId, "", enum.GROUP_UNMUTE_MSG)
	} else {
		err := daoSer.Add(&model.GroupFilter{
			GroupID:    enum.WorldGroupId,
			UserID:     "0",
			Types:      enum.GroupMuteSwitch,
			CreateTime: time.Now(),
		})
		if err != nil {
			panic(i18n_err.I18nError{Code: i18n_err.ErrorCodeSystemBusy})
		}
		new(service_im.ImPublicService).SendCustomMsg(enum.WorldGroupId, "", enum.GROUP_MUTE_MSG)
	}
	return
}

// 消息通知设置
func (g *Group) WorldGroupNoticeSetting(c *gin.Context, req *group2.WorldGroupNoticeSettingReq) {
	userId := helper.GetUserId(c)
	clientType := helper.GetClientType(c)
	daoSer := new(dao.GroupFilterDao)
	groupFilterModel := daoSer.FindOne(&model.GroupFilter{
		GroupID: enum.WorldGroupId,
		UserID:  userId,
	})
	if groupFilterModel.ID > 0 {
		groupFilterModel.Types = req.Types
		daoSer.Save(groupFilterModel)
	} else {
		err := daoSer.Add(&model.GroupFilter{
			GroupID:    enum.WorldGroupId,
			UserID:     userId,
			Types:      req.Types,
			CreateTime: time.Now(),
		})
		if err != nil {
			panic(i18n_err.I18nError{Code: i18n_err.ErrorCodeSystemBusy})
		}
	}
	switch req.Types {
	case enum.GroupReceiveAllMsg, enum.GroupReceiveAdminMsg:
		new(service_im.ImGroupService).SendActionMsg(c, nil, userId, "", enum.WorldGroupId, clientType, enum.GROUP_SUBSCRIBE)
	case enum.GroupNotReceiveMsg:
		new(service_im.ImGroupService).SendActionMsg(c, nil, userId, "", enum.WorldGroupId, clientType, enum.GROUP_UN_SUBSCRIBE)
	}
	return
}
