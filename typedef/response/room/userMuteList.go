package room

import "yfapi/typedef/response"

type UserMuteListRes struct {
	UserId     string                  `json:"userId"`     //用户id
	UserNo     string                  `json:"userNo"`     //用户序号
	Nickname   string                  `json:"nickname"`   //昵称
	Avatar     string                  `json:"avatar"`     //头像
	Sex        int                     `json:"sex"`        //性别;0:保密,1:男,2:女
	UserPlaque response.UserPlaqueInfo `json:"userPlaque"` // 用户铭牌信息
}
