package user

import (
	"yfapi/typedef/response"
)

// UserVisitInfo 访客足迹信息
type UserVisitInfo struct {
	UserId     string                  `json:"userId"`     // 用户ID
	Nickname   string                  `json:"nickname"`   // 用户昵称
	Avatar     string                  `json:"avatar"`     // 用户头像
	Sex        int                     `json:"sex"`        // 用户性别
	Introduce  string                  `json:"introduce"`  // 个性签名
	TimeDesc   string                  `json:"timeDesc"`   // 时间描述
	Extra      string                  `json:"extra"`      // 额外描述信息 (你也看过她，他也看过你)
	UserPlaque response.UserPlaqueInfo `json:"userPlaque"` // 用户铭牌信息
}
