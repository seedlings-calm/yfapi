package service_im

import (
	"github.com/gin-gonic/gin"
	"yfapi/internal/dao"
	"yfapi/typedef/enum"
	"yfapi/typedef/message"

	"github.com/spf13/cast"
)

type ImNoticeService struct {
}

// 发送互动消息
func (i *ImNoticeService) SendInteractiveMsg(c *gin.Context, data message.SendInteractiveMsg, toUserId []string) {
	new(ImOneService).SendNotice(c, enum.InteractiveUserId, toUserId, enum.MsgCustom, data, enum.USER_INTERACTIVE_MSG)
	if len(toUserId) > 0 {
		for _, userId := range toUserId {
			//增加会话列表
			ChatListDealWith(data.Msg, userId, enum.InteractiveUserId, enum.MsgListTextColorNormal)
			//处理消息未读数
			AddNotReadNum(userId, enum.InteractiveUserId)
		}
	}
	return
}

// 发送系统消息
func (i *ImNoticeService) SendSystematicMsg(c *gin.Context, title, img, content, link, h5Content string, toUserId []string) {
	msgData := message.SystematicMsg{
		Title:     title,
		Img:       img,
		Content:   content,
		H5Content: h5Content,
		Link:      link,
	}
	new(ImOneService).SendNotice(c, enum.SystematicUserId, toUserId, enum.MsgCustom, msgData, enum.USER_SYSTEM_MSG)
	AddFunc := func(ids []string) {
		for _, userId := range ids {
			//增加会话列表
			ChatListDealWith(content, userId, enum.SystematicUserId, enum.MsgListTextColorNormal)
			//处理消息未读数
			AddNotReadNum(userId, enum.SystematicUserId)
		}
	}
	if len(toUserId) > 0 {
		AddFunc(toUserId)
	} else {
		var offsetUserId int64 = 0
		var limit = 500
		db := new(dao.UserDao)
		for {
			userIds := db.GetUserIdsOffsetId(offsetUserId, limit)
			if len(userIds) > 0 {
				offsetUserId = cast.ToInt64(userIds[len(userIds)-1])
				go AddFunc(userIds)
			}
			if len(userIds) < limit {
				break
			}
		}
	}
}

// 发送官方公告消息
func (i *ImNoticeService) SendOfficialMsg(c *gin.Context, title, img, content, link, h5Content string, toUserId []string) error {
	msgData := message.OfficialMsg{
		Title:     title,
		Img:       img,
		Content:   content,
		H5Content: h5Content,
		Link:      link,
	}
	new(ImOneService).SendNotice(c, enum.OfficialUserId, toUserId, enum.MsgCustom, msgData, enum.USER_OFFICIAL_MSG)
	AddFunc := func(ids []string) {
		for _, userId := range ids {
			//增加会话列表
			ChatListDealWith(content, userId, enum.OfficialUserId, enum.MsgListTextColorNormal)
			//处理消息未读数
			AddNotReadNum(userId, enum.OfficialUserId)
		}
	}
	if len(toUserId) > 0 {
		AddFunc(toUserId)
	} else { //获取所有用户
		var offsetUserId int64 = 0
		var limit = 500
		db := new(dao.UserDao)
		for {
			userIds := db.GetUserIdsOffsetId(offsetUserId, limit)
			if len(userIds) > 0 {
				offsetUserId = cast.ToInt64(userIds[len(userIds)-1])
				go AddFunc(userIds)
			}
			if len(userIds) < limit {
				break
			}
		}
	}
	return nil
}
