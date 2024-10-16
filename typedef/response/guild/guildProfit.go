package guild

import (
	"yfapi/util/easy"
)

// MemberProfit
// @Description: 公会成员流水
type MemberProfit struct {
	StatDate     string `json:"statDate"`     // 统计日期
	UserId       string `json:"userId"`       // 成员长ID
	UserNo       string `json:"userNo"`       // 成员ID
	Nickname     string `json:"nickname"`     // 成员昵称
	Practitioner string `json:"practitioner"` // 从业身份
	RewardCount  int64  `json:"rewardCount"`  // 被打赏次数
	ProfitAmount string `json:"profitAmount"` // 被打赏钻石
	Income       string `json:"income"`       // 预估收益
}

// RoomProfit
// @Description: 公会房间流水
type RoomProfit struct {
	RoomId        string `json:"roomId"`        // 房间长ID
	RoomNo        string `json:"roomNo"`        // 房间ID
	RoomName      string `json:"roomName"`      // 房间名称
	RoomType      string `json:"roomType"`      // 房间类型
	RoomOwnerNo   string `json:"roomOwnerNo"`   // 房主ID
	RoomOwnerName string `json:"roomOwnerName"` // 房主昵称
	RewardCount   int64  `json:"rewardCount"`   // 打赏次数
	ProfitAmount  string `json:"profitAmount"`  // 打赏钻石
	Income        string `json:"income"`        // 预估收益
}

// RewardDetail
// @Description: 礼物打赏详情
type RewardDetail struct {
	FromUserId   string         `json:"fromUserId"`   // 打赏人长ID
	FromUserNo   string         `json:"fromUserNo"`   // 打赏人ID
	FromNickname string         `json:"fromNickname"` // 打赏人昵称
	ToUserId     string         `json:"toUserId"`     // 被打赏人长ID
	ToUserNo     string         `json:"toUserNo"`     // 被打赏人ID
	ToNickname   string         `json:"toNickname"`   // 被打赏人昵称
	GiftName     string         `json:"giftName"`     // 礼物名称
	GiftCount    int            `json:"giftCount"`    // 礼物数量
	GiftPrice    int            `json:"giftPrice"`    // 礼物价值
	CreateTime   easy.LocalTime `json:"createTime"`   // 打赏时间
}

// AccountBill
// @Description: 账户交易明细
type AccountBill struct {
	OrderId           string                  `json:"orderId"`        // 订单ID
	OrderType         int                     `json:"orderType"`      // 资金类型
	OrderTypeDesc     string                  `json:"orderTypeDesc"`  // 资金类型描述
	Memo              string                  `json:"memo"`           // 备注
	FundFlow          int                     `json:"fundFlow"`       // 1入 2出
	BeforeAmount      string                  `json:"beforeAmount"`   // 交易前余额
	Amount            string                  `json:"amount"`         // 金额（元）
	CurrAmount        string                  `json:"currAmount"`     // 交易后余额
	WithdrawStatus    int                     `json:"withdrawStatus"` // 提现订单状态
	CreateTime        easy.LocalTime          `json:"createTime"`     // 日期
	Withdraw          *WithdrawDetail         `json:"withdraw,omitempty" gorm:"-"`
	SubsidyGuild      *SubsidyGuildMonth      `json:"subsidyGuild,omitempty" gorm:"-"`
	SubsidyGuildValid *SubsidyGuildValidMonth `json:"subsidyGuildValid,omitempty" gorm:"-"`
	SubsidyRoom       *SubsidyRoom            `json:"subsidyRoom,omitempty" gorm:"-"`
}

// WithdrawDetail
// @Description: 提现详情
type WithdrawDetail struct {
	Nickname     string `json:"nickname"`     // 用户昵称
	UserNo       string `json:"userNo"`       // 用户ID
	TrueName     string `json:"trueName"`     // 真实姓名
	Mobile       string `json:"mobile"`       // 手机号
	Amount       string `json:"amount"`       // 提现金额
	PayAmount    string `json:"payAmount"`    // 到账金额
	BankUserName string `json:"bankUserName"` // 收款人
	BankNo       string `json:"bankNo"`       // 银行卡号
}

// SubsidyGuildMonth
// @Description: 公会月结补贴
type SubsidyGuildMonth struct {
	StatDate     string `json:"statDate"`     // 结算日期
	GuildName    string `json:"guildName"`    // 公会名称
	GuildNo      string `json:"guildNo"`      // 公会ID
	ProfitAmount string `json:"profitAmount"` // 打赏流水
	Income       string `json:"income"`       // 发放收益
}

// SubsidyGuildValidMonth
// @Description: 公会有效直播间月结补贴
type SubsidyGuildValidMonth struct {
	StatDate     string `json:"statDate"`     // 结算日期
	GuildName    string `json:"guildName"`    // 公会名称
	GuildNo      string `json:"guildNo"`      // 公会ID
	RoomCount    int    `json:"roomCount"`    // 直播间数量
	ValidCount   int    `json:"validCount"`   // 有效直播间数量
	ProfitAmount string `json:"profitAmount"` // 有效直播间流水
	Income       string `json:"income"`       // 发放收益
}

// SubsidyRoom
// @Description: 房间日/月结补贴
type SubsidyRoom struct {
	StateDate     string `json:"stateDate"`     // 结算日期
	RoomName      string `json:"roomName"`      // 房间名称
	RoomNo        string `json:"roomNo"`        // 房间ID
	RoomOwnerName string `json:"roomOwnerName"` // 房主昵称
	RoomOwnerNo   string `json:"roomOwnerNo"`   // 房主ID
	SettleName    string `json:"settleName"`    // 结算人昵称
	SettleNo      string `json:"settleNo"`      // 结算人ID
	ProfitAmount  string `json:"profitAmount"`  // 打赏流水
	Income        string `json:"income"`        // 发放收益
}
