package guild

// AccountInfoRes
// @Description: 账户信息
type AccountInfoRes struct {
	UserId          string     `json:"userId"`
	Status          int        `json:"status"`          // 账户状态 1=正常 2=冻结
	CashAmount      string     `json:"cashAmount"`      // 资产总额
	TotalCashIncome string     `json:"totalCashIncome"` // 累计收益
	SettlementRate  int        `json:"settlementRate"`  // 结算费率
	Desc            string     `json:"desc"`            // 提现说明
	Mobile          string     `json:"mobile"`          // 手机号
	ReginCode       string     `json:"reginCode"`       // 区号
	TrueName        string     `json:"trueName"`        // 真实姓名
	BankList        []BankInfo `json:"bankList"`
}
type BankInfo struct {
	Id         int    `json:"id"`         // 银行卡信息Id
	BankNo     string `json:"bankNo"`     // 银行卡号
	BankName   string `json:"bankName"`   // 银行名称
	BankHolder string `json:"bankHolder"` // 银行账户名
	BankBranch string `json:"bankBranch"` // 银行支行
	IsDefault  int    `json:"isDefault"`  // 是否默认 1=是 2=否
}

// StatGuildInfo
// @Description: 首页公会统计信息
type StatGuildInfo struct {
	Room struct {
		RoomCount       int64 `json:"roomCount"`       // 房间数量
		ChatroomCount   int64 `json:"chatroomCount"`   // 聊天室数量
		AnchorRoomCount int64 `json:"anchorRoomCount"` // 直播间数量
	} `json:"room"`
	Member struct {
		MemberCount int64 `json:"memberCount"` // 公会人数
		NormalCount int64 `json:"normalCount"` // 普通用户数量
		CertCount   int64 `json:"certCount"`   // 资质用户数量
	} `json:"member"`
	Practitioners struct {
		PractitionersCount int `json:"practitionersCount"` // 从业者数量
		CompereCount       int `json:"compereCount"`       // 主持人数量
		MusicianCount      int `json:"musicianCount"`      // 音乐人数量
		CounselorCount     int `json:"counselorCount"`     // 咨询师数量
		AnchorCount        int `json:"anchorCount"`        // 主播数量
	} `json:"practitioners"`
	TodayProfit struct {
		TotalProfit    string `json:"totalProfit"`    // 今日流水
		ChatroomProfit string `json:"chatroomProfit"` // 聊天室流水
		AnchorProfit   string `json:"anchorProfit"`   // 直播间流水
	} `json:"todayProfit"`
	MothProfit struct {
		TotalProfit    string `json:"totalProfit"`    // 月流水
		ChatroomProfit string `json:"chatroomProfit"` // 聊天室流水
		AnchorProfit   string `json:"anchorProfit"`   // 直播间流水
	} `json:"mothProfit"`
}

type ProfitInfo struct {
	StatDate     string `json:"statDate"`     // 统计日期
	ProfitAmount string `json:"profitAmount"` // 流水
}

type RoomTypeCount struct {
	Name  string `json:"name"`  // 厅类型名称
	Count int    `json:"count"` // 厅类型房间数量
}

type RoomTypeProfit struct {
	Name         string `json:"name"`         // 厅名称
	ProfitAmount string `json:"profitAmount"` // 厅类型总流水
}

// ProfitGuildInfo
// @Description: 首页公会流水信息
type ProfitGuildInfo struct {
	LatestWeek     []ProfitInfo     `json:"latestWeek"`     // 最近七天流水
	RoomCategory   []RoomTypeCount  `json:"roomCategory"`   // 房间占比
	ProfitCategory []RoomTypeProfit `json:"profitCategory"` // 流水占比
}

// RoomRank
// @Description: 房间流水排行榜
type RoomRank struct {
	RoomName      string `json:"roomName"`      // 房间名称
	RoomOwnerName string `json:"roomOwnerName"` // 房主名称
	OnlineSecond  string `json:"onlineSecond"`  // 在线时长
	ProfitAmount  string `json:"profitAmount"`  // 流水
}
