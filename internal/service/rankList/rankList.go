package rankList

import (
	context2 "context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	"yfapi/internal/helper"
	"yfapi/internal/service/auth"
	typedef_enum "yfapi/typedef/enum"
)

// 榜单
type RankListService struct {
}

type CalculateReq struct {
	FromUserId string //来源用户
	ToUserId   string //目标用户
	Types      string // freeGift免费礼物，luckGift幸运礼物，gift普通礼物，firstJoinRoom首次进入直播间，publicMessage发送公屏消息，retentionTime观看时长，onMicTime在麦时长
	Diamond    int    //钻石数量
	duration   int    //时长秒数
	RoomId     string
}

func Instance() *RankListService {
	return &RankListService{}
}

// 计算积分
func (r *RankListService) Calculate(req CalculateReq) {
	switch req.Types {
	case "freeGift": //打赏免费礼物
		r.freeGiftScore(req)
	case "luckGift":
		r.luckGiftScore(req)
	case "gift":
		r.luckGiftScore(req)
	case "firstJoinRoom":
		r.firstJoinRoomScore(req)
	case "publicMessage":
		r.publicMessageScore(req)
	case "retentionTime":
		r.retentionTimeScore(req)
	case "onMicTime":
		r.OnMicTimeScore(req)
	}
}

// 在麦时长积分计算
func (r *RankListService) OnMicTimeScore(req CalculateReq) {
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	onMicTimeNumKey := OnMicTimeNumKey(req.FromUserId)
	onMicTimeNum := rd.Get(context, onMicTimeNumKey).Val()
	if cast.ToInt(onMicTimeNum) < OnMicTimeNum {
		rd.Incr(context, onMicTimeNumKey)
		rd.Expire(context, onMicTimeNumKey, time.Hour*24)
		//房间人气榜增加积分
		r.RoomPopularityScore(req.RoomId, req.FromUserId, 1)
		//平台人气榜增加积分
		r.PopularityScore(req.FromUserId, req.RoomId, 1)
	}
}

// 观看时长积分计算
func (r *RankListService) retentionTimeScore(req CalculateReq) {
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	retentionNumKey := RetentionNumKey(req.FromUserId)
	retentionNum := rd.Get(context, retentionNumKey).Val()
	if cast.ToInt(retentionNum) < RetentionNum {
		rd.Incr(context, retentionNumKey)
		rd.Expire(context, retentionNumKey, time.Hour*24)
		//房间贡献榜增加积分
		r.RoomContributorScore(req.RoomId, req.FromUserId, 1)
		//平台贡献榜增加积分
		r.ContributorScore(req.FromUserId, 1)
	}
}

// 免费礼物积分计算
func (r *RankListService) freeGiftScore(req CalculateReq) {
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	freeGiftPresentedNumKey := FreeGiftPresentedNumKey(req.FromUserId)
	freeGiftPresentedNum := rd.Get(context, freeGiftPresentedNumKey).Val()
	if cast.ToInt(freeGiftPresentedNum) < FreeGiftPresentedNum {
		rd.Incr(context, freeGiftPresentedNumKey)
		rd.Expire(context, freeGiftPresentedNumKey, time.Hour*24)
		//房间贡献榜增加积分
		r.RoomContributorScore(req.RoomId, req.FromUserId, 1)
		//平台贡献榜增加积分
		r.ContributorScore(req.FromUserId, 1)
	}
	//人气榜单计算
	freeGiftReceiveNumKey := FreeGiftReceiveNumKey(req.ToUserId)
	freeGiftReceiveNum := rd.Get(context, freeGiftReceiveNumKey).Val()
	if cast.ToInt(freeGiftReceiveNum) < FreeGiftReceiveNum {
		rd.Incr(context, freeGiftReceiveNumKey)
		rd.Expire(context, freeGiftReceiveNumKey, time.Hour*24)
		//房间人气榜增加积分
		r.RoomPopularityScore(req.RoomId, req.ToUserId, 1)
		//平台人气榜增加积分
		r.PopularityScore(req.ToUserId, req.RoomId, 1)
	}
}

// 幸运礼物积分计算
func (r *RankListService) luckGiftScore(req CalculateReq) {
	//房间贡献榜增加积分
	r.RoomContributorScore(req.RoomId, req.FromUserId, req.Diamond)
	//平台贡献榜增加积分
	r.ContributorScore(req.FromUserId, req.Diamond)
	//房间人气榜增加积分
	r.RoomPopularityScore(req.RoomId, req.ToUserId, req.Diamond)
	//平台人气榜增加积分
	r.PopularityScore(req.ToUserId, req.RoomId, req.Diamond)
}

// 首次进入直播间
func (r *RankListService) firstJoinRoomScore(req CalculateReq) {
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	firstJoinRoomKey := FirstJoinRoomKey(req.FromUserId)
	ok := rd.SetNX(context, firstJoinRoomKey, 1, time.Hour*24).Val()
	if ok {
		//房间贡献榜增加积分
		r.RoomContributorScore(req.RoomId, req.FromUserId, 1)
		//平台贡献榜增加积分
		r.ContributorScore(req.FromUserId, 1)
	}
}

// 发送公屏积分计算
func (r *RankListService) publicMessageScore(req CalculateReq) {
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	sendPublicMessageNumKey := SendPublicMessageNumKey(req.FromUserId)
	sendPublicMessageNum := rd.Get(context, sendPublicMessageNumKey).Val()
	if cast.ToInt(sendPublicMessageNum) < SendPublicMessageNum {
		rd.Incr(context, sendPublicMessageNumKey)
		rd.Expire(context, sendPublicMessageNumKey, time.Hour*24)
		//房间贡献榜增加积分
		r.RoomContributorScore(req.RoomId, req.FromUserId, 1)
		//平台贡献榜增加积分
		r.ContributorScore(req.FromUserId, 1)
	}
}

// 房间贡献榜计算
func (r *RankListService) RoomContributorScore(roomId, userId string, score int) {
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	roomDayContributorKey := RoomDayContributorKey(roomId, time.Now().Format(time.DateOnly))
	rd.ZIncrBy(context, roomDayContributorKey, float64(score), userId)
	rd.Expire(context, roomDayContributorKey, time.Hour*24)
	roomWeekContributorKey := RoomWeekContributorKey(roomId, fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
	rd.ZIncrBy(context, roomWeekContributorKey, float64(score), userId)
	rd.Expire(context, roomWeekContributorKey, time.Hour*24*7)
	roomMonthContributorKey := RoomMonthContributorKey(roomId, helper.GetMonthDateTime())
	rd.ZIncrBy(context, roomMonthContributorKey, float64(score), userId)
	rd.Expire(context, roomMonthContributorKey, time.Hour*24*31)
}

// 房间人气榜计算
func (r *RankListService) RoomPopularityScore(roomId, userId string, score int) {
	// 是否为本房间的从业者 主持、音乐人、咨询师、主播
	checkRoleIdList := []int{typedef_enum.CompereRoleId, typedef_enum.MusicianRoleId, typedef_enum.CounselorRoleId, typedef_enum.AnchorRoleId}
	isHave := new(auth.Auth).IsHaveCurrRole(roomId, userId, checkRoleIdList)
	if !isHave {
		return
	}
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	roomDayPopularityKey := RoomDayPopularityKey(roomId, time.Now().Format(time.DateOnly))
	rd.ZIncrBy(context, roomDayPopularityKey, float64(score), userId)
	rd.Expire(context, roomDayPopularityKey, time.Hour*24)
	roomWeekPopularityKey := RoomWeekPopularityKey(roomId, fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
	rd.ZIncrBy(context, roomWeekPopularityKey, float64(score), userId)
	rd.Expire(context, roomWeekPopularityKey, time.Hour*24*7)
	roomMonthPopularityKey := RoomMonthPopularityKey(roomId, helper.GetMonthDateTime())
	rd.ZIncrBy(context, roomMonthPopularityKey, float64(score), userId)
	rd.Expire(context, roomMonthPopularityKey, time.Hour*24*31)
}

// 平台贡献榜计算
func (r *RankListService) ContributorScore(userId string, score int) {
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	dayContributorKey := DayContributorKey(time.Now().Format(time.DateOnly))
	rd.ZIncrBy(context, dayContributorKey, float64(score), userId)
	rd.Expire(context, dayContributorKey, time.Hour*24)
	weekContributorKey := WeekContributorKey(fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
	rd.ZIncrBy(context, weekContributorKey, float64(score), userId)
	rd.Expire(context, weekContributorKey, time.Hour*24*7)
	monthContributorKey := MonthContributorKey(helper.GetMonthDateTime())
	rd.ZIncrBy(context, monthContributorKey, float64(score), userId)
	rd.Expire(context, monthContributorKey, time.Hour*24*31)
}

// 平台人气榜计算
func (r *RankListService) PopularityScore(userId, roomId string, score int) {
	// 是否为本房间的从业者 主持、音乐人、咨询师、主播
	checkRoleIdList := []int{typedef_enum.CompereRoleId, typedef_enum.MusicianRoleId, typedef_enum.CounselorRoleId, typedef_enum.AnchorRoleId}
	isHave := new(auth.Auth).IsHaveCurrRole(roomId, userId, checkRoleIdList)
	if !isHave {
		return
	}
	rd := coreRedis.GetChatroomRedis()
	context := context2.Background()
	dayPopularityKey := DayPopularityKey(time.Now().Format(time.DateOnly))
	rd.ZIncrBy(context, dayPopularityKey, float64(score), userId)
	rd.Expire(context, dayPopularityKey, time.Hour*24)
	weekPopularityKey := WeekPopularityKey(fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
	rd.ZIncrBy(context, weekPopularityKey, float64(score), userId)
	rd.Expire(context, weekPopularityKey, time.Hour*24*7)
	monthPopularityKey := MonthPopularityKey(helper.GetMonthDateTime())
	rd.ZIncrBy(context, monthPopularityKey, float64(score), userId)
	rd.Expire(context, monthPopularityKey, time.Hour*24*31)
}

// 获取榜单数据
// cycleTypes 日榜 day，周榜 week，月榜 month
// rankListTypes 榜单类型 contributor 贡献榜 popularity 人气榜
func (r *RankListService) GetRankList(roomId, cycleTypes, rankListTypes string) []redis.Z {
	var key = RoomDayContributorKey(roomId, time.Now().Format(time.DateOnly))
	switch cycleTypes {
	case CycleDay:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				key = RoomDayContributorKey(roomId, time.Now().Format(time.DateOnly))
			}
			if rankListTypes == RankListPopularity {
				key = RoomDayPopularityKey(roomId, time.Now().Format(time.DateOnly))
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				key = DayContributorKey(time.Now().Format(time.DateOnly))
			}
			if rankListTypes == RankListPopularity {
				key = DayPopularityKey(time.Now().Format(time.DateOnly))
			}
		}
	case CycleWeek:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				key = RoomWeekContributorKey(roomId, fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
			}
			if rankListTypes == RankListPopularity {
				key = RoomWeekPopularityKey(roomId, fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				key = WeekContributorKey(fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
			}
			if rankListTypes == RankListPopularity {
				key = WeekPopularityKey(fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
			}
		}
	case CycleMonth:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				key = RoomMonthContributorKey(roomId, helper.GetMonthDateTime())
			}
			if rankListTypes == RankListPopularity {
				key = RoomMonthPopularityKey(roomId, helper.GetMonthDateTime())
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				key = MonthContributorKey(helper.GetMonthDateTime())
			}
			if rankListTypes == RankListPopularity {
				key = MonthPopularityKey(helper.GetMonthDateTime())
			}
		}
	}
	res := coreRedis.GetChatroomRedis().ZRevRangeWithScores(context2.Background(), key, 0, 19).Val()
	return res
}

// 获取用户的排名变化 与上次比较
func (r *RankListService) GetUserRankingChange(userId, roomId, cycleTypes, rankListTypes string) int64 {
	var thisKey = RoomDayContributorKey(roomId, time.Now().Format(time.DateOnly)) //本期
	var lastKey = RoomDayContributorKey(roomId, time.Now().Format(time.DateOnly)) //上期
	switch cycleTypes {
	case CycleDay:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				thisKey = RoomDayContributorKey(roomId, time.Now().Format(time.DateOnly))
				lastKey = RoomDayContributorKey(roomId, time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
			}
			if rankListTypes == RankListPopularity {
				thisKey = RoomDayPopularityKey(roomId, time.Now().Format(time.DateOnly))
				lastKey = RoomDayPopularityKey(roomId, time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				thisKey = DayContributorKey(time.Now().Format(time.DateOnly))
				lastKey = DayContributorKey(time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
			}
			if rankListTypes == RankListPopularity {
				thisKey = DayPopularityKey(time.Now().Format(time.DateOnly))
				lastKey = DayPopularityKey(time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
			}
		}
	case CycleWeek:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				thisKey = RoomWeekContributorKey(roomId, fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
				lastKey = RoomWeekContributorKey(roomId, fmt.Sprintf("%s-%s", helper.GetLastMondayDateTime(), helper.GetLastSundayDateTime()))
			}
			if rankListTypes == RankListPopularity {
				thisKey = RoomWeekPopularityKey(roomId, fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
				lastKey = RoomWeekPopularityKey(roomId, fmt.Sprintf("%s-%s", helper.GetLastMondayDateTime(), helper.GetLastSundayDateTime()))
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				thisKey = WeekContributorKey(fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
				lastKey = WeekContributorKey(fmt.Sprintf("%s-%s", helper.GetLastMondayDateTime(), helper.GetLastSundayDateTime()))
			}
			if rankListTypes == RankListPopularity {
				thisKey = WeekPopularityKey(fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
				lastKey = WeekPopularityKey(fmt.Sprintf("%s-%s", helper.GetLastMondayDateTime(), helper.GetLastSundayDateTime()))
			}
		}
	case CycleMonth:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				thisKey = RoomMonthContributorKey(roomId, helper.GetMonthDateTime())
				lastKey = RoomMonthContributorKey(roomId, helper.GetLastMonthDateTime())
			}
			if rankListTypes == RankListPopularity {
				thisKey = RoomMonthPopularityKey(roomId, helper.GetMonthDateTime())
				lastKey = RoomMonthPopularityKey(roomId, helper.GetLastMonthDateTime())
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				thisKey = MonthContributorKey(helper.GetMonthDateTime())
				lastKey = MonthContributorKey(helper.GetLastMonthDateTime())
			}
			if rankListTypes == RankListPopularity {
				thisKey = MonthPopularityKey(helper.GetMonthDateTime())
				lastKey = MonthPopularityKey(helper.GetLastMonthDateTime())
			}
		}
	}
	thisRank, _ := coreRedis.GetChatroomRedis().ZRevRank(context2.Background(), thisKey, userId).Result()
	lastRank, err := coreRedis.GetChatroomRedis().ZRevRank(context2.Background(), lastKey, userId).Result()
	if err != nil {
		coreLog.Error("用户上期没有排名 userId:%s,thisKey:%s,lastKey:%s,err:%v", userId, thisKey, lastKey, err)
		return thisRank
	}
	return -(thisRank - lastRank)
}

// 获取用户的排名也信息
func (r *RankListService) GetMemberRankList(userId, roomId, cycleTypes, rankListTypes string) (int, int) {
	var key = RoomDayContributorKey(roomId, time.Now().Format(time.DateOnly)) //本期
	switch cycleTypes {
	case CycleDay:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				key = RoomDayContributorKey(roomId, time.Now().Format(time.DateOnly))
			}
			if rankListTypes == RankListPopularity {
				key = RoomDayPopularityKey(roomId, time.Now().Format(time.DateOnly))
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				key = DayContributorKey(time.Now().Format(time.DateOnly))
			}
			if rankListTypes == RankListPopularity {
				key = DayPopularityKey(time.Now().Format(time.DateOnly))
			}
		}
	case CycleWeek:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				key = RoomWeekContributorKey(roomId, fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
			}
			if rankListTypes == RankListPopularity {
				key = RoomWeekPopularityKey(roomId, fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				key = WeekContributorKey(fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
			}
			if rankListTypes == RankListPopularity {
				key = WeekPopularityKey(fmt.Sprintf("%s-%s", helper.GetMondayDateTime(), helper.GetSundayDateTime()))
			}
		}
	case CycleMonth:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				key = RoomMonthContributorKey(roomId, helper.GetMonthDateTime())
			}
			if rankListTypes == RankListPopularity {
				key = RoomMonthPopularityKey(roomId, helper.GetMonthDateTime())
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				key = MonthContributorKey(helper.GetMonthDateTime())
			}
			if rankListTypes == RankListPopularity {
				key = MonthPopularityKey(helper.GetMonthDateTime())
			}
		}
	}
	score := coreRedis.GetChatroomRedis().ZScore(context2.Background(), key, userId).Val()
	rank := coreRedis.GetChatroomRedis().ZRevRank(context2.Background(), key, userId).Val()
	return cast.ToInt(score), cast.ToInt(rank)
}

// 获取榜单前几名数据
// num 去前几名
func (r *RankListService) GetRankListSettle(num int64, roomId, cycleTypes, rankListTypes string) []redis.Z {
	var lastKey = ""
	switch cycleTypes {
	case CycleDay:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				lastKey = RoomDayContributorKey(roomId, time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
			}
			if rankListTypes == RankListPopularity {
				lastKey = RoomDayPopularityKey(roomId, time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				lastKey = DayContributorKey(time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
			}
			if rankListTypes == RankListPopularity {
				lastKey = DayPopularityKey(time.Now().AddDate(0, 0, -1).Format(time.DateOnly))
			}
		}
	case CycleWeek:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				lastKey = RoomWeekContributorKey(roomId, fmt.Sprintf("%s-%s", helper.GetLastMondayDateTime(), helper.GetLastSundayDateTime()))
			}
			if rankListTypes == RankListPopularity {
				lastKey = RoomWeekPopularityKey(roomId, fmt.Sprintf("%s-%s", helper.GetLastMondayDateTime(), helper.GetLastSundayDateTime()))
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				lastKey = WeekContributorKey(fmt.Sprintf("%s-%s", helper.GetLastMondayDateTime(), helper.GetLastSundayDateTime()))
			}
			if rankListTypes == RankListPopularity {
				lastKey = WeekPopularityKey(fmt.Sprintf("%s-%s", helper.GetLastMondayDateTime(), helper.GetLastSundayDateTime()))
			}
		}
	case CycleMonth:
		if len(roomId) > 0 { //房间榜单
			if rankListTypes == RankListContributor {
				lastKey = RoomMonthContributorKey(roomId, helper.GetLastMonthDateTime())
			}
			if rankListTypes == RankListPopularity {
				lastKey = RoomMonthPopularityKey(roomId, helper.GetLastMonthDateTime())
			}
		} else { //平台榜单
			if rankListTypes == RankListContributor {
				lastKey = MonthContributorKey(helper.GetLastMonthDateTime())
			}
			if rankListTypes == RankListPopularity {
				lastKey = MonthPopularityKey(helper.GetLastMonthDateTime())
			}
		}
	}
	res := []redis.Z{}
	if len(lastKey) > 0 {
		res = coreRedis.GetChatroomRedis().ZRevRangeWithScores(context2.Background(), lastKey, 0, num-1).Val()
	}
	return res
}
