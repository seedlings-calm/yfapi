package dao

import (
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/corePg"
	"yfapi/internal/model"
)

type ChatStoreDao struct {
}

func (c *ChatStoreDao) GetUniteId(userId1, userId2 string) string {
	u1 := cast.ToInt64(userId1)
	u2 := cast.ToInt64(userId2)
	if u2 > u1 {
		return cast.ToString(u1) + cast.ToString(u2)
	}
	return cast.ToString(u2) + cast.ToString(u1)
}

// 获取私聊历史消息
func (c *ChatStoreDao) GetPrivateChatList(selfUserId, otherUserId string, limit, offsetTimestamp int) []model.PrivateChat {
	if limit == 0 {
		limit = 10
	}
	if offsetTimestamp == 0 {
		offsetTimestamp = int(time.Now().UnixMicro())
	}
	list := []model.PrivateChat{}
	params := map[string]any{
		"unite_id": c.GetUniteId(selfUserId, otherUserId),
	}
	timestamp := c.GetPrivateChatClearTimestamp(selfUserId, otherUserId)
	err := corePg.NewDb().Db().Model(&model.PrivateChat{}).
		Where(params).
		Where("timestamp > ? and timestamp < ?", timestamp, offsetTimestamp).
		Not("from_user_id = ? and status = ?", otherUserId, 2).
		Limit(limit).
		Order("timestamp desc").
		Find(&list).Error
	if err != nil {
		coreLog.LogError("GetPrivateChatList err:%+v", err)
	}
	return list
}

// 设置个人私聊消息为已读状态
func (c *ChatStoreDao) SetPrivateChatIsRead(selfUserId, otherUserId string) bool {
	params := map[string]any{
		"unite_id":   c.GetUniteId(selfUserId, otherUserId),
		"type":       1,
		"to_user_id": selfUserId,
		"read":       2,
	}
	err := corePg.NewDb().Db().Model(&model.PrivateChat{}).Where(params).Update("read", 1).Error
	if err != nil {
		coreLog.LogError("SetPrivateChatRead err:%+v", err)
		return false
	}
	return true
}

// 设置系统私聊消息为已读状态
func (c *ChatStoreDao) SetPrivateChatSystemIsRead(selfUserId, otherUserId string) bool {
	params := map[string]any{
		"unite_id":   c.GetUniteId(selfUserId, otherUserId),
		"type":       2,
		"to_user_id": selfUserId,
		"read":       2,
	}
	err := corePg.NewDb().Db().Model(&model.PrivateChat{}).Where(params).Update("read", 1).Error
	if err != nil {
		coreLog.LogError("SetPrivateChatSystemIsRead err:%+v", err)
		return false
	}
	return true
}

// 记录私聊消息到数据库
func (c *ChatStoreDao) WritePrivateChatMessage(message model.PrivateChat) error {
	err := corePg.NewDb().Db().Create(&message).Error
	return err
}

// 记录公屏消息到数据库
func (c *ChatStoreDao) WritePublicChatMessage(message model.PublicChat) error {
	err := corePg.NewDb().Db().Create(&message).Error
	return err
}

// 获取用户会话清除的最后时间戳
func (c *ChatStoreDao) GetPrivateChatClearTimestamp(userId, receiveUserId string) int64 {
	params := map[string]any{
		"unite_id": userId + receiveUserId,
	}
	data := &model.PrivateChatClearRecord{}
	err := corePg.NewDb().Db().Model(&model.PrivateChatClearRecord{}).Where(params).First(data).Error
	if err != nil {
		return 0
	}
	return data.Timestamp
}

// 设置用户会话清除最后时间戳
func (c *ChatStoreDao) SetPrivateChatClearTimestamp(userId, receiveUserId string) {
	data := &model.PrivateChatClearRecord{}
	err := corePg.NewDb().Db().Where("unite_id = ?", userId+receiveUserId).First(&data).Error
	if err == gorm.ErrRecordNotFound {
		data.Timestamp = time.Now().UnixMicro()
		data.UniteId = userId + receiveUserId
		corePg.NewDb().Db().Create(&data)
		return
	}
	if len(data.UniteId) != 0 {
		data.Timestamp = time.Now().UnixMicro()
		corePg.NewDb().Db().Where(&model.PrivateChatClearRecord{
			UniteId: userId + receiveUserId,
		}).Save(data)
	}
	return
}
