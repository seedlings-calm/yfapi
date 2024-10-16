package roomOwner

import "yfapi/typedef/request"

// RoomAdminListReq 房间管理员列表请求
type RoomAdminListReq struct {
	request.PageInfo
}

// RoomPractitionerListReq 房间从业者列表请求
type RoomPractitionerListReq struct {
	request.PageInfo
	UserKeyword      string `json:"userKeyword"`      // 从业者昵称/ID
	PractitionerType int    `json:"practitionerType"` // 从业者身份
	Status           int    `json:"status"`           // 状态
}
type RoomPractitionerAddReq struct {
	UserId            string `json:"userId" form:"userId"`                        //用户ID
	PractitionerType  int    `json:"practitionerType" form:"practitioner_type"`   //从业者类型 1主持 2音乐 3咨询 4主播
	PractitionerBrief string `json:"practitionerBrief" form:"practitioner_brief"` //从业者简介
}
type RoomPractitionerUpdateReq struct {
	Id                int    `json:"id" form:"id"`                                //记录ID
	PractitionerBrief string `json:"practitionerBrief" form:"practitioner_brief"` //从业者简介
}
type RoomPractitionerApply struct {
	UserId        string `json:"userId" form:"userId"`                //用户ID
	Status        int    `json:"status" form:"status"`                // 1正常,2审核中,3审核拒绝
	AbolishReason string `json:"abolishReason" form:"abolish_reason"` //拒绝理由
}

// 搞个通用的入参Userid
type RoomCommonReq struct {
	UserId string `json:"userId" form:"userId"`
}
