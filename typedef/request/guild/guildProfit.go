package guild

import "yfapi/typedef/request"

// GetGuildMemberProfitReq
// @Description: 公会成员流水列表请求
type GetGuildMemberProfitReq struct {
	UserKeyword       string   `json:"userKeyword" form:"userKeyword"`             // 成员昵称/ID
	PractitionersType int      `json:"practitionersType" form:"practitionersType"` // 从业者身份
	DateRange         []string `json:"dateRange" form:"dateRange"`                 // 查询时间 2024-09-26
	request.PageInfo
}

// GetGuildRoomProfitReq
// @Description: 公会房间流水列表请求
type GetGuildRoomProfitReq struct {
	RoomKeyword string   `json:"roomKeyword" form:"roomKeyword"` // 房间名称/ID
	RoomType    int      `json:"roomType" form:"roomType"`       // 房间类型
	DateRange   []string `json:"dateRange" form:"dateRange"`     // 查询时间 2024-09-26
	request.PageInfo
}

// GetGuildRewardListReq
// @Description: 公会打赏详情列表
type GetGuildRewardListReq struct {
	RewardType int      `json:"rewardType" form:"rewardType"` // 打赏类型 1成员 2房间
	Uid        string   `json:"uid" form:"uid"`               // 用户ID/房间ID
	DateRange  []string `json:"dateRange" form:"dateRange"`   // 查询时间 2024-09-26
	request.PageInfo
}

// GetAccountBillReq
// @Description: 账户交易明细列表
type GetAccountBillReq struct {
	OrderType      int      `json:"orderType" form:"orderType"`           // 资金类型
	FundFlow       int      `json:"fundFlow" form:"fundFlow"`             // 变动类型
	WithdrawStatus string   `json:"withdrawStatus" form:"withdrawStatus"` // 提现订单状态
	DateRange      []string `json:"dateRange" form:"dateRange"`           // 查询时间 2024-09-27
	request.PageInfo
}
