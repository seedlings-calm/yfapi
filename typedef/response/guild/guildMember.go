package guild

import (
	"yfapi/util/easy"
)

type MemberList struct {
	Id            int      `json:"Id"`
	UserId        string   `json:"UserId"`
	UserNo        string   `json:"UserNo"`
	UserNickname  string   `json:"UserNickname"`
	UserAvatar    string   `json:"UserAvatar"`
	IsChairman    string   `json:"IsChairman"`
	UserRole      []string `json:"UserRole"`
	UserSkillList []string `json:"UserSkillList"`
	GroupId       int      `json:"GroupId"`
	GroupName     string   `json:"GroupName"`
	RewardCount   string   `json:"RewardCount"`
	RewardSum     string   `json:"RewardSum"`
	Status        string   `json:"Status"`
	LeaveTime     string   `json:"LeaveTime"`
	KickoutReason string   `json:"KickoutReason"`
	CreateTime    string   `json:"CreateTime"`
}

// LeaveMemberShipListRsp
// @Description: 退会申请信息
type LeaveMemberShipListRsp struct {
	Id             int64          `json:"id"`             // 主键ID
	UserId         string         `json:"userId"`         // 用户长ID
	UserNo         string         `json:"userNo"`         // 用户ID
	Nickname       string         `json:"nickname"`       // 用户昵称
	Avatar         string         `json:"avatar"`         // 用户头像
	Force          int            `json:"force"`          // 是否强制退会 0否 1是
	Status         int            `json:"status"`         // 状态 1=待审核, 2=同意, 3=拒绝, 4=自动拒绝, 5=强制申请自动退出, 6=取消申请
	RewardNum      int            `json:"rewardNum"`      // 被打赏次数
	RewardDiamonds int            `json:"rewardDiamonds"` // 被打赏钻石
	CreateTime     easy.LocalTime `json:"createTime"`     // 申请时间
	JoinTime       easy.LocalTime `json:"joinTime"`       // 入会时间
	UpdateTime     easy.LocalTime `json:"updateTime"`     // 操作时间
}

type MemberGroup struct {
	Count int               `json:"count"`
	List  []MemberGroupList `json:"list"`
}
type MemberGroupList struct {
	Id          int    `json:"Id"`
	GuildId     string `json:"GuildId"`
	Name        string `json:"Name"`
	Desc        string `json:"Desc"`
	State       string `json:"State"`
	CreateAt    string `json:"CreateAt"`
	UpdateAt    string `json:"UpdateAt"`
	MemberCount int    `json:"MemberCount"`
}

type GuildMemberListRes struct {
	Id           int64           `json:"id"`
	UserId       string          `json:"userId"`       //
	UserNo       string          `json:"userNo"`       //
	Nickname     string          `json:"nickname"`     //
	Avatar       string          `json:"avatar"`       //
	Identity     string          `json:"identity"`     //职位
	IdCards      string          `json:"idCards"`      //从业者身份
	GroupName    string          `json:"groupName"`    //分组名
	RewardCount  int64           `json:"rewardCount"`  // 打赏次数
	ProfitAmount string          `json:"profitAmount"` // 打赏钻石
	JoinTime     *easy.LocalTime `json:"joinTime"`
}

type MemberIdcardsInfoRes struct {
	IdCards    string         `json:"idCards"`  //从业者身份
	RoomNo     string         `json:"RoomNo"`   //房间ID
	RoomName   string         `json:"roomName"` //房间名称
	UserId     string         `json:"userId"`   //房主ID
	UserNo     string         `json:"userNo"`   //房主ID
	Nickname   string         `json:"nickname"` //房主昵称
	CreateTime easy.LocalTime `json:"createTime"`
}

// MemberJoinApplyInfo
// @Description: 申请入会信息
type MemberJoinApplyInfo struct {
	Id         int64          `json:"id"`         // 主键ID
	UserId     string         `json:"userId"`     // 用户长ID
	UserNo     string         `json:"userNo"`     // 用户ID
	Nickname   string         `json:"nickname"`   // 用户昵称
	Avatar     string         `json:"avatar"`     // 用户头像
	Status     int            `json:"status"`     // 状态 1=待审核, 2=同意, 3=拒绝, 4=自动拒绝, 5=强制申请自动退出, 6=取消申请
	Reason     string         `json:"reason"`     // 原因
	CreateTime easy.LocalTime `json:"createTime"` // 申请时间
	UpdateTime easy.LocalTime `json:"updateTime"` // 审核时间
}

// UserPractitionerAction
// @Description: 从业者行为记录
type UserPractitionerAction struct {
	UserId     string         `json:"userId"`     // 用户长ID
	UserNo     string         `json:"userNo"`     // 用户ID
	Nickname   string         `json:"nickname"`   // 用户昵称
	Avatar     string         `json:"avatar"`     // 用户头像
	Action     string         `json:"action"`     // 用户操作
	CreateTime easy.LocalTime `json:"createTime"` // 记录时间
}
