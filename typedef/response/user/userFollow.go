package user

import "yfapi/typedef/response"

// 获取关注列表||粉丝列表||好友列表
type GetUserFollowingList struct {
	UserId         string                  `json:"userId"` //用户id
	UserNo         string                  `json:"userNo"`
	Uid32          int32                   `json:"uid32"`
	Nickname       string                  `json:"nickname"`       //昵称
	Avatar         string                  `json:"avatar"`         //头像
	Sex            int                     `json:"sex"`            //性别
	Introduce      string                  `json:"introduce"`      // 个性签名
	IsMutualFollow bool                    `json:"isMutualFollow"` //是否互相关注
	FollowedType   int                     `json:"followedType"`   //关注状态 0未关注对方 1已关注对方 2互相关注 3对方关注你
	UserPlaque     response.UserPlaqueInfo `json:"userPlaque"`     // 用户铭牌信息
}

// AddFollowRes
// @Description: 关注/取消关注返回
type AddFollowRes struct {
	FollowedType int `json:"followedType"` //0未关注对方 1已关注对方 2互相关注
	FollowedNum  int `json:"followedNum"`
	FansNum      int `json:"fansNum"`
}
