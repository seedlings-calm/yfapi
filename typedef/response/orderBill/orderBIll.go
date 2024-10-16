package orderBill

// DiamondBill 钻石流水
type DiamondBill struct {
	Title      string      `json:"title"`                // 标题描述
	Amount     string      `json:"amount"`               // 变动金额
	FundFlow   int         `json:"fundFlow"`             // 1入账 2出账
	RoomName   string      `json:"roomName,omitempty"`   // 房间名称
	ToUserList []*UserInfo `json:"toUserList,omitempty"` // 礼物打赏人数列表
	CreateTime string      `json:"createTime"`           // 交易时间
	TimeKey    string      `json:"timeKey"`              // 时间年月
}

// UserInfo 用户信息
type UserInfo struct {
	UserId   string `json:"userId"`   // 用户ID
	Nickname string `json:"nickname"` // 用户昵称
	Avatar   string `json:"avatar"`   // 用户头像
}

// StarlightBill 星光流水
type StarlightBill struct {
	Title      string    `json:"title"`              // 标题描述
	Amount     string    `json:"amount"`             // 变动金额
	FundFlow   int       `json:"fundFlow"`           // 1入账 2出账
	RoomName   string    `json:"roomName,omitempty"` // 房间名称
	FormUser   *UserInfo `json:"formUser,omitempty"` // 打赏人信息
	CreateTime string    `json:"createTime"`         // 交易时间
	TimeKey    string    `json:"timeKey"`            // 时间年月
}

// WithdrawNote
// @Description: 提现备注
type WithdrawNote struct {
	BankName   string `json:"bankName"`
	BankNo     string `json:"bankNo"`
	BankHolder string `json:"bankHolder"`
	BankCode   string `json:"bankCode"`
	BankBranch string `json:"bankBranch"`
	StaffName  string `json:"staffName"`
	Reason     string `json:"reason"`
}

// SubsidyChatroomNote
// @Description: 房间日结备注信息
type SubsidyChatroomNote struct {
	ProfitAmount string `json:"profitAmount"` // 礼物总流水
	ProfitRate   string `json:"profitRate"`   // 流水补贴比例
	OnlineSecond int64  `json:"onlineSecond"` // 开播时长
	OnlineRate   string `json:"onlineRate"`   // 开播时长补贴比例
	ValidDay     int    `json:"validDay"`     // 有效开播天数
}

// SubsidyGuildNote
// @Description: 公会月结备注信息
type SubsidyGuildNote struct {
	ProfitAmount string           `json:"profitAmount"` // 公会总流水
	ProfitRate   string           `json:"profitRate"`   // 公会补贴比例
	List         []RoomProfitInfo `json:"list"`         // 公会房间流水详情
}

type RoomProfitInfo struct {
	RoomId  string `json:"roomId"`  // 房间ID
	Amount  string `json:"amount"`  // 房间流水
	IsValid bool   `json:"isValid"` // 是否有效
}

// SubsidyGuildValidNote
// @Description: 公会有效直播间备注信息
type SubsidyGuildValidNote struct {
	ProfitAmount string            `json:"profitAmount"` // 有效直播间流水
	ValidCount   int               `json:"validCount"`   // 有效直播间个数
	ValidRate    string            `json:"validRate"`    // 有效直播间补贴比例
	List         []GuildAnchorInfo `json:"list"`         // 公会直播间详情
}

type GuildAnchorInfo struct {
	RoomId    string `json:"roomId"`    // 房间ID
	OnlineDay int    `json:"onlineDay"` // 开播天数
	ValidDay  int    `json:"validDay"`  // 有效开播天数
	Amount    string `json:"amount"`    // 房间流水
	IsValid   bool   `json:"isValid"`   // 是否有效
}
