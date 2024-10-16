package guild

import "yfapi/typedef/request"

// GuildRoomListreq 公会房间列表请求
type GuildRoomListreq struct {
	request.PageInfo
	RoomKeyword string `json:"roomKeyword"`                  //房间名称或ID
	UserKeyword string `json:"userKeyword"`                  //房主昵称或ID
	RoomType    int    `json:"roomType"`                     //房间类型
	Status      int    `json:"status"`                       //状态
	LiveType    int    `json:"liveType" validate:"required"` //直播类型 1:聊天室 2:直播
}
type RoomTypeReq struct {
	LiveType int8 `json:"live_type"` // 房间直播类型 1聊天室 2直播间 3个人房
}
type ChangeRoomParamReq struct {
	UserNo string `json:"userNo"` // 用户ID
	RoomID string `json:"roomId"` // 房间ID
	Op     int    `json:"op"`     // 更新类型 1房主 2日结算 3月结算
}
type UserNoParamReq struct {
	UserNo string `json:"userNo"` // 用户ID
}
type CloseRoomReq struct {
	RoomID string `json:"roomId"` // 房间ID
	Status int    `json:"status"` // 状态 1开启 2关闭
}

// GuildMemberListreq 公会成员列表请求
type GuildMemberListreq struct {
	request.PageInfo
	Status      int    `json:"status"`       //状态
	UserKeyword string `json:"user_keyword"` //成员昵称或ID
	GroupID     int    `json:"group_id"`     //分组
}

// 入会申请列表请求参数
type MemberShipListreq struct {
	request.PageInfo
	Status      int    `json:"status"`      //状态 1=待审核, 2=同意, 3=拒绝, 4=自动拒绝, 5=强制申请自动退出, 6=取消申请
	UserKeyword string `json:"userKeyword"` //成员昵称或ID
}

type GuildRoomApplyReq struct {
	RoomName          string `json:"roomName"  v:"required|max=20#房间名称不能为空且最多20个字"` //房间名称
	UserNo            string `json:"userNo" v:"required#房主不能为空"`                    //房主userNo
	DaySettleUserNo   string `json:"DaySettleUserNo"  v:"required#日收益人不能为空"`        //日收益人UserNo
	MonthSettleUserNo string `json:"monthSettleUserNo" v:"required#月受益人不能为空"`       //月受益人UserNo
	RoomType          int    `json:"roomType"  v:"required#房间类型不能为空"`               //房间类型
	TemplateId        string `json:"templateId" v:"required#房间模板不能为空"`              //房间模板
	RoomDesc          string `json:"roomDesc"  v:"required#房间描述不能为空"`               //房间描述
	CoverImg          string `json:"coverImg"  v:"required#厅图不能为空"`                 //厅图
}

type GuildRoomApplyListReq struct {
	request.PageInfo
	RoomKeyword string `json:"roomKeyword"` //房间名称/roomNo
	UserKeyWord string `json:"userKeyword"` //房主昵称/userNO
	Status      int    `json:"status"`      //审核状态 1-待审核 2-审核通过 3-审核拒绝
}

// GuildMemberApplyReviewReq 公会成员入会申请审核请求参数
type GuildMemberApplyReviewReq struct {
	Id     int    `json:"id"`     //申请ID
	Status int    `json:"status"` //审核状态  2= 同意 3=拒绝
	Reason string `json:"reason"` //拒绝原因
}

// GuildMemberWithdrawReviewReq 公会成员退会申请审核请求参数
type GuildMemberWithdrawReviewReq struct {
	Id     int    `json:"id"`     //申请ID
	Status int    `json:"status"` //审核状态  2= 同意 3=拒绝
	Reason string `json:"reason"` //拒绝原因
}

type GuildMemberListReq struct {
	request.PageInfo
	UserKeyword string   `json:"user_keyword"` //成员昵称或ID
	GroupID     int      `json:"group_id"`     //分组
	IdCard      []string `json:"idCard"`       //（从业者类型 1主持 2音乐 3咨询 4主播）
}
