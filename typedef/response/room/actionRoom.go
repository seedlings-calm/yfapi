package room

import "yfapi/typedef/response"

type CheckRoomResponse struct {
	IsPwd       bool   `json:"isPwd"`       //房间是否有密码
	IsBlacklist bool   `json:"isBlacklist"` //是否被拉黑
	IsKickOut   bool   `json:"isKickOut"`   //是否被踢出
	Msg         string `json:"msg"`         // 检查房间不能进入的描述信息
	IsMulti     bool   `json:"isMulti"`     //是否多端登录
}

type CheckIsRoomResponse struct {
	IsInRoom bool `json:"isInRoom"` //是否在此房间
}

type ExecCommandRes struct {
}

type UpSeatApplyInfo struct {
	UserId     string                  `json:"userId"` //用户id
	UserNo     string                  `json:"userNo"`
	Uid32      int32                   `json:"uid32"`
	Nickname   string                  `json:"nickname"`   //昵称
	Avatar     string                  `json:"avatar"`     //头像
	Sex        int                     `json:"sex"`        //性别;0:保密,1:男,2:女
	UserPlaque response.UserPlaqueInfo `json:"userPlaque"` // 用户铭牌信息
}

type UpSeatApplyListRes struct {
	List  []UpSeatApplyInfo `json:"list"`  // 上麦申请列表
	Count int               `json:"count"` // 列表数量
}

type HoldUpSeatUserListRes struct {
	List  []UpSeatApplyInfo `json:"list"`  // 可抱上麦的用户列表
	Count int               `json:"count"` // 列表数量
}
