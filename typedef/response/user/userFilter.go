package user

import "yfapi/typedef/response"

type UserFilterList struct {
	UserId     string                  `json:"userId"` //用户id
	UserNo     string                  `json:"userNo"`
	Uid32      int32                   `json:"uid32"`
	Nickname   string                  `json:"nickname"`   //昵称
	Avatar     string                  `json:"avatar"`     //头像
	Sex        int                     `json:"sex"`        //性别
	UserPlaque response.UserPlaqueInfo `json:"userPlaque"` // 用户铭牌信息
}
