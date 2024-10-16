package rankList

import (
	"fmt"
	"time"
)

const (
	FreeGiftPresentedNum = 5  //免费礼物赠送计算次数
	FreeGiftReceiveNum   = 5  //免费记录接收次数
	SendPublicMessageNum = 5  //公屏消息次数
	RetentionNum         = 5  //累计观看积分数
	OnMicTimeNum         = 10 //累计在麦分数
	CycleDay             = "day"
	CycleWeek            = "week"
	CycleMonth           = "month"
	RankListContributor  = "contributor" //贡献榜
	RankListPopularity   = "popularity"  //人气榜
)

// 房间贡献日榜
var RoomDayContributorKey = func(roomId string, key string) string {
	return fmt.Sprintf("rankList:contributor:roomDay:%s:%s", roomId, key)
}

// 房间贡献周榜
var RoomWeekContributorKey = func(roomId string, key string) string {
	return fmt.Sprintf("rankList:contributor:roomWeek:%s:%s", roomId, key)
}

// 房间贡献月榜
var RoomMonthContributorKey = func(roomId string, key string) string {
	return fmt.Sprintf("rankList:contributor:roomMonth:%s:%s", roomId, key)
}

// 平台贡献日榜
var DayContributorKey = func(key string) string {
	return fmt.Sprintf("rankList:contributor:day:%s", key)
}

// 平台贡献周榜
var WeekContributorKey = func(key string) string {
	return fmt.Sprintf("rankList:contributor:week:%s", key)
}

// 平台贡献月榜
var MonthContributorKey = func(key string) string {
	return fmt.Sprintf("rankList:contributor:month:%s", key)
}

// 房间人气日榜
var RoomDayPopularityKey = func(roomId string, key string) string {
	return fmt.Sprintf("rankList:popularity:roomDay:%s:%s", roomId, key)
}

// 房间人气周榜
var RoomWeekPopularityKey = func(roomId string, key string) string {
	return fmt.Sprintf("rankList:popularity:roomWeek:%s:%s", roomId, key)
}

// 房间人气月榜
var RoomMonthPopularityKey = func(roomId string, key string) string {
	return fmt.Sprintf("rankList:popularity:roomMonth:%s:%s", roomId, key)
}

// 平台人气日榜
var DayPopularityKey = func(key string) string {
	return fmt.Sprintf("rankList:popularity:day:%s", key)
}

// 平台人气周榜
var WeekPopularityKey = func(key string) string {
	return fmt.Sprintf("rankList:popularity:week:%s", key)
}

// 平台人气月榜
var MonthPopularityKey = func(key string) string {
	return fmt.Sprintf("rankList:popularity:month:%s", key)
}

// 免费礼物接受次数
var FreeGiftReceiveNumKey = func(userId string) string {
	return fmt.Sprintf("rankList:freeGiftReceiveNum:%s:%s", userId, time.Now().Format(time.DateOnly))
}

// 免费礼物赠送次数
var FreeGiftPresentedNumKey = func(userId string) string {
	return fmt.Sprintf("rankList:freeGiftPresentedNum:%s:%s", userId, time.Now().Format(time.DateOnly))
}

// 首次进入直播间
var FirstJoinRoomKey = func(userId string) string {
	return fmt.Sprintf("rankList:FirstJoinRoomKey:%s:%s", userId, time.Now().Format(time.DateOnly))
}

// 发送公屏消息次数限制
var SendPublicMessageNumKey = func(userId string) string {
	return fmt.Sprintf("rankList:sendPublicMessageNumKey:%s:%s", userId, time.Now().Format(time.DateOnly))
}

// 观看直播积分数限制
var RetentionNumKey = func(userId string) string {
	return fmt.Sprintf("rankList:retentionNumKey:%s:%s", userId, time.Now().Format(time.DateOnly))
}

// 在麦积分数限制
var OnMicTimeNumKey = func(userId string) string {
	return fmt.Sprintf("rankList:onMicTimeNumKey:%s:%s", userId, time.Now().Format(time.DateOnly))
}
