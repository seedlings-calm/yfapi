package redisKey

import (
	"fmt"
	"time"
	"yfapi/util/easy"

	"github.com/spf13/cast"
)

var (
	// GetGiftListCacheKey 礼物列表缓存
	GetGiftListCacheKey = func(categoryType int) string {
		return fmt.Sprintf("GiftListCache:%v", categoryType)
	}

	// GetGiftVersionKey 礼物版本号
	GetGiftVersionKey = func(categoryType int) string {
		return fmt.Sprintf("GiftVersion:%v", categoryType)
	}

	SendGiftLockKey = func(userId string) string {
		return fmt.Sprintf("SendGiftLockKey:%v", userId)
	}
	//头条中心数据key
	TopMsgKey = func() string {
		return "TopMsgLog"
	}
	//连击存储的处理逻辑数据
	TopMsgCallBack = func(key string) string {
		return fmt.Sprintf("TopMsg:%s", key)
	}
	//延迟队列key
	TopMsgEqueueKey = func() string {
		return "TopMsg:Equeue"
	}
	// 房间麦位存储
	RoomWheatPosition = func(roomId string) string {
		return fmt.Sprintf("roomWheatPosition:%v", roomId)
	}
	// RoomUpSeatApplyList 房间上麦申请列表
	RoomUpSeatApplyList = func(roomId string) string {
		return fmt.Sprintf("roomSeatApplyList:%v", roomId)
	}
	RoomUpSeatApplyInfo = func(roomId string) string {
		return fmt.Sprintf("roomSeatApplyInfo:%v", roomId)
	}
	//房间用户在线列表
	RoomUsersOnlineList = func(roomId string) string {
		return fmt.Sprintf("roomUsers:%s:online", roomId)
	}
	//房间用户今天列表 单位：天
	RoomUsersDayList = func(roomId string) string {
		day := time.Now().Format("20060102")
		return fmt.Sprintf("roomUsers:%s:%s", roomId, day)
	}

	//房间在线用户列表缓存,无身份查看
	RoomOnlineUsersCache = func(roomId string) string {
		return fmt.Sprintf("roomOnlineUsersCache:%s", roomId)
	}
	//房间在线用户列表缓存,有身份查看
	RoomOnlineUsersIdCardCache = func(roomId string) string {
		return fmt.Sprintf("roomOnlineUsersIdCardCache:%s", roomId)
	}

	//房间1000贡献榜缓存key
	RoomDayUsersCache = func(roomId string) string {
		return fmt.Sprintf("roomDayUsersCache:%s", roomId)
	}

	//房间高等级用户缓存key
	RoomHightGradeUsersCache = func(roomId string) string {
		return fmt.Sprintf("roomHightGradeCache:%s", roomId)
	}
	//房间高等级用户统计缓存key
	RoomHightGradeUsersCountCache = func() string {
		return "roomHightGradeCountCache"
	}

	GiftComboHitCountKey = func(roomId, userId, giftCode, toUserId string, giftCount int) string {
		md5Key := easy.Md5(roomId+userId+giftCode+toUserId+cast.ToString(giftCount), 0, false)
		return fmt.Sprintf("GiftComboHitCount:%v", md5Key)
	}

	InformationCardKey = func(userId string, roomId string) string {
		return fmt.Sprintf("informationCard:%s:%s", roomId, userId)
	}

	UserRoomSwitchStatus = func(userId string, roomId string) string {
		return fmt.Sprintf("userRoomSwitchStatus:%s:%s", roomId, userId)
	}
	//房间踢出的用户存储key
	RoomKickOutKey = func(userId string, roomId string) string {
		return fmt.Sprintf("roomKickOut:%s:%s", roomId, userId)
	}
	ChatroomUserCharmKey = func(roomId string) string {
		today := time.Now().Format("20060102")
		return fmt.Sprintf("roomUserCharm:%v:%v", roomId, today)
	}
	AnchorRoomUserCharmKey = func(roomId string) string {
		return fmt.Sprintf("anchorRoomUserCharm:%v", roomId)
	}
	//房间自动欢迎语key
	RoomAutoWelcomeKey = func(roomId, userId string) string {
		return fmt.Sprintf("roomAutoWelcome:%s:%s", roomId, userId)
	}

	//房间隐藏麦
	RoomHiddenMicKey = func(roomId string) string {
		return fmt.Sprintf("roomHiddenMic:%s", roomId)
	}

	//房间热度值排行
	RoomHotKey = func(liveType int, timeKey ...string) string {
		nowTimeKey := time.Now().Format("2006010215")
		if len(timeKey) > 0 {
			nowTimeKey = timeKey[0]
		}
		return fmt.Sprintf("roomHot:timer:%v:%v", liveType, nowTimeKey)
	}

	RoomHotTotalKey = func(liveType int, timeKey ...string) string {
		nowTimeKey := time.Now().Format("2006010215")
		if len(timeKey) > 0 {
			nowTimeKey = timeKey[0]
		}
		return fmt.Sprintf("roomHot:rank:%v:%v", liveType, nowTimeKey)
	}

	//进厅热度值缓存
	RoomHotJoinRoomKey = func(roomId string) string {
		return fmt.Sprintf("roomHot:joinRoom:%v:%v", roomId, time.Now().Format("2006010215"))
	}

	// 房间发言热度值缓存
	RoomHotPublicChat = func(roomId string) string {
		return fmt.Sprintf("roomHot:publicChat:%v:%v", roomId, time.Now().Format("2006010215"))
	}
	//房间上下播统计数据缓存
	RoomWheatTimeCacheKey = func(roomId string) string {
		return fmt.Sprintf("roomWheatTime:room:%s", roomId)
	}
	//房间上下播统计数据缓存-加入房间的用户存储
	RoomWheatTimeJoinUser = func(roomId string) string {
		return fmt.Sprintf("roomWheatTime:users:%s", roomId)
	}
	//房间上下播统计数据缓存-打赏的用户存储
	RoomWheatTimeGiftUser = func(roomId string) string {
		return fmt.Sprintf("roomWheatTime:giftUsers:%s", roomId)
	}

	//公会后台相关缓存

	// GuildStatInfoKey 首页公会统计信息
	GuildStatInfoKey = func(guildId string) string {
		return fmt.Sprintf("guildStatInfo:%v", guildId)
	}
	// GuildProfitInfoKey 首页公会流水信息
	GuildProfitInfoKey = func(guildId string) string {
		return fmt.Sprintf("guildProfitInfo:%v", guildId)
	}
	// GuildRoomRankKey 首页公会房间排行榜信息
	GuildRoomRankKey = func(guildId string) string {
		return fmt.Sprintf("guildRoomRank:%v", guildId)
	}
	// GetChatroomOpeningListKey 在线聊天室列表
	GetChatroomOpeningListKey = func() string {
		return fmt.Sprintf("ChatroomOpeningList")
	}
	// GuildChatroomMonthProfitKey 公会聊天室月流水
	GuildChatroomMonthProfitKey = func(guildId, statDate string) string {
		return fmt.Sprintf("GuildChatroomMonthProfit:%v:%v", guildId, statDate)
	}

	//清理房间不在线用户列表锁
	ClearRoomNotOnlineUserLock = func(roomId string) string {
		return fmt.Sprintf("ClearRoomNotOnlineUserLock:%s", roomId)
	}

	//用户在房持续时间锁
	UserInRoomRetentionTime = func(userId string) string {
		return fmt.Sprintf("UserInRoomRetentionTime:%s", userId)
	}

	//用户在房在麦时长
	UserInRoomOnMicTime = func(userId string) string {
		return fmt.Sprintf("UserInRoomOnMicTime:%s", userId)
	}
)
