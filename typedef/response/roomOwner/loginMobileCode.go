package roomOwner

type SendSMSRes struct {
}
type RoomInfo struct {
	RoomName   string `json:"roomName"`   // 房间名称
	UserAvatar string `json:"userAvatar"` // 房主头像
	UserNo     string `json:"userNo"`     // 房主ID
	UserName   string `json:"userName"`   // 房主昵称
	LogoImg    string `json:"logoImg"`    // 房间厅图
	RoomNo     string `json:"roomNo"`     // 房间ID
	RoomID     string `json:"roomId"`     // 房间长ID
}
type LoginMobileCodeRes struct {
	Token  string `json:"token"`  // token
	UserID string `json:"userId"` // 用户ID
}

// 聊天室资料
type RoomHomeInfo struct {
	RoomName      string `json:"roomName"`      // 房间名称
	UserAvatar    string `json:"userAvatar"`    // 房主头像
	UserNo        string `json:"userNo"`        // 房主ID
	UserName      string `json:"userName"`      // 房主昵称
	LogoImg       string `json:"logoImg"`       // 房间厅图
	RoomType      string `json:"roomType"`      // 房间类型
	CreateTime    string `json:"createTime"`    // 创建时间
	Notice        string `json:"notice"`        // 房间公告
	RoomNo        string `json:"roomNo"`        // 房间ID
	Welcome       string `json:"welcome"`       // 欢迎语
	RoomID        string `json:"roomId"`        // 房间长ID
	UserID        string `json:"userId"`        // 房主长ID
	LastLoginTime string `json:"lastLoginTime"` // 上次登录时间 TODO
}

// 聊天室概况
type RoomHomeBaseInfo struct {
	OnlineSecond string       `json:"onlineSecond"`        //今日开播时长
	EnterTimes   int          `json:"enterTimes"`          //本日进房人次
	TodayProfit  string       `json:"todayProfit"`         //今日流水（钻）
	Practitioner int          `json:"practitioner"`        //从业者数量
	Host         int          `json:"host"`                //主持人数量
	Musician     int          `json:"musician"`            //音乐人数量
	Counselor    int          `json:"counselor"`           //咨询师数量
	MonthProfit  string       `json:"monthProfit"`         //本月流水（钻）
	LatestWeek   []ProfitInfo `json:"latestWeek" gorm:"-"` // 最近七天流水
}

type ProfitInfo struct {
	StatDate     string `json:"statDate"`     // 统计日期
	ProfitAmount string `json:"profitAmount"` // 流水
}
