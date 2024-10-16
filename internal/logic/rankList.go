package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"yfapi/internal/helper"
	rankListSer "yfapi/internal/service/rankList"
	"yfapi/internal/service/user"
	rankListReq "yfapi/typedef/request/rankList"
	rankListResp "yfapi/typedef/response/rankList"
	"yfapi/util/easy"
)

type RankListLogic struct{}

// 获取排行榜数据
func (r *RankListLogic) GetRankList(c *gin.Context, req *rankListReq.RankListReq) (res rankListResp.RankListResp) {
	ser := rankListSer.Instance()
	rankList := ser.GetRankList(req.RoomId, req.Range, req.Types)
	if len(rankList) == 0 {
		return
	}
	selfUserId := helper.GetUserId(c)
	selfRankInfo := rankListResp.UserRankInfo{}
	for k, v := range rankList {
		userId := cast.ToString(v.Member)
		userInfo := user.GetUserInfo(userId, "")
		rankListInfo := rankListResp.UserRankInfo{
			UserId:        userId,
			UserNo:        userInfo.UserNo,
			Nickname:      userInfo.Nickname,
			Avatar:        userInfo.Avatar,
			Sex:           userInfo.Sex,
			UserPlaque:    userInfo.UserPlaque,
			Ranking:       k + 1,
			Score:         cast.ToInt(v.Score),
			FormatScore:   easy.NumberToW(v.Score, 1),
			RankingChange: ser.GetUserRankingChange(userId, req.RoomId, req.Range, req.Types),
			IsOnRanking:   true,
		}
		switch {
		case k == 0:
			//第一名 距离上一名为0
			rankListInfo.Behind = 0
			if len(rankList) > 2 {
				rankListInfo.Ahead = rankListInfo.Score - cast.ToInt(rankList[k+1].Score)
			}
		case k == len(rankList)-1: //最后一名
			rankListInfo.Behind = cast.ToInt(rankList[k-1].Score) - rankListInfo.Score
		default:
			rankListInfo.Behind = cast.ToInt(rankList[k-1].Score) - rankListInfo.Score
			rankListInfo.Ahead = rankListInfo.Score - cast.ToInt(rankList[k+1].Score)
		}
		res.List = append(res.List, rankListInfo)
		if userId == selfUserId { //自己在榜上
			selfRankInfo = rankListInfo
		}
	}
	//计算距离结束差值
	switch req.Range {
	case rankListSer.CycleDay:
		res.EndTime = helper.GetUntilDayTime()
	case rankListSer.CycleWeek:
		res.EndTime = helper.GetUntilSundayTime()
	case rankListSer.CycleMonth:
		res.EndTime = helper.GetUntilMonthTime()
	}
	//未上榜
	if len(selfRankInfo.UserId) == 0 {
		selfInfo := user.GetUserInfo(selfUserId, helper.GetClientType(c))
		score, rank := ser.GetMemberRankList(selfUserId, req.RoomId, req.Range, req.Types)
		res.SelfInfo = rankListResp.UserRankInfo{
			UserId:        selfUserId,
			UserNo:        selfInfo.UserNo,
			Nickname:      selfInfo.Nickname,
			Avatar:        selfInfo.Avatar,
			Sex:           selfInfo.Sex,
			UserPlaque:    selfInfo.UserPlaque,
			Ranking:       rank,
			Score:         score,
			FormatScore:   easy.NumberToW(score, 1),
			RankingChange: 0,
			IsOnRanking:   false,
			Distance:      cast.ToInt(rankList[len(rankList)-1].Score),
		}
	} else {
		res.SelfInfo = selfRankInfo
	}
	return
}
