package rankList

import "yfapi/typedef/response"

type RankListResp struct {
	List     []UserRankInfo `json:"list"`     //榜单信息
	SelfInfo UserRankInfo   `json:"selfInfo"` //自己榜单信息
	EndTime  int64          `json:"endTime"`  //结束时间
}

type UserRankInfo struct {
	UserId        string                  `json:"userId"`        // 用户id
	UserNo        string                  `json:"userNo"`        // 用户No
	Nickname      string                  `json:"nickname"`      // 昵称
	Avatar        string                  `json:"avatar"`        // 头像
	Sex           int                     `json:"sex"`           // 性别
	UserPlaque    response.UserPlaqueInfo `json:"userPlaque"`    // 用户铭牌信息
	Ranking       int                     `json:"ranking"`       //排名
	Score         int                     `json:"score"`         //得分
	FormatScore   string                  `json:"formatScore"`   //展示得分
	RankingChange int64                   `json:"rankingChange"` //名次变化
	IsOnRanking   bool                    `json:"isOnRanking"`   //是否上榜
	Behind        int                     `json:"behind"`        //与上一名差值
	Ahead         int                     `json:"ahead"`         //与下一名差值
	Distance      int                     `json:"distance"`      //距离上榜
}
