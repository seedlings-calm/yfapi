package roomOwner

import "yfapi/util/easy"

// RoomAdminInfo
// @Description: 房间管理员信息
type RoomAdminInfo struct {
	UserId     string         `json:"userId"`     // 管理员长ID
	UserNo     string         `json:"userNo"`     // 管理员ID
	Nickname   string         `json:"nickname"`   // 管理员昵称
	Avatar     string         `json:"avatar"`     // 管理员头像
	StaffName  string         `json:"staffName"`  // 操作人昵称
	CreateTime easy.LocalTime `json:"createTime"` // 添加时间
}

// RoomPractitionerInfo
// @Description: 房间从业者信息
type RoomPractitionerInfo struct {
	ID                   int64          `json:"Id"`                   // 主键ID
	UserId               string         `json:"userId"`               // 从业者长ID
	UserNo               string         `json:"userNo"`               // 从业者ID
	Nickname             string         `json:"nickname"`             // 从业者昵称
	Avatar               string         `json:"avatar"`               // 从业者头像
	PractitionerType     int            `json:"practitionerType"`     // 从业者身份
	PractitionerTypeDesc string         `json:"practitionerTypeDesc"` // 从业者身份描述
	Status               int            `json:"status"`               // 状态 1正常,2审核中,3审核拒绝
	AbolishReason        string         `json:"abolishReason"`        // 拒绝/取消原因
	CreateTime           easy.LocalTime `json:"createTime"`           // 添加时间
	UpdateTime           easy.LocalTime `json:"updateTime"`           // 通过时间
}

type SearchUser struct {
	UserId   string `json:"userId"`   // 长ID
	UserNo   string `json:"userNo"`   // ID
	Nickname string `json:"nickname"` // 昵称
	Avatar   string `json:"avatar"`   // 头像
}
