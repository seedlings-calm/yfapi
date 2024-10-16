package roomOwner

type PersonListRes struct {
	UserId       string `json:"userId"`       //
	UserNo       string `json:"userNo"`       // 短ID
	Nickname     string `json:"nickname"`     // 昵称
	Avatar       string `json:"avatar"`       //头像
	IdCards      string `json:"idCards"`      //身份
	Times        string `json:"times"`        //开播时间处理
	TimesNum     int    `json:"timesNum"`     //开播时间
	RewardCount  int    `json:"rewardCount"`  // 打赏次数
	ProfitAmount int    `json:"profitAmount"` // 打赏钻石
}

type PersonListDetailRes struct {
	UserId    string `json:"userId"`    //
	UserNo    string `json:"userNo"`    // 短ID
	Nickname  string `json:"nickname"`  // 昵称
	Avatar    string `json:"avatar"`    //头像
	GiftName  string `json:"giftName"`  //礼物
	GiftPrice int    `json:"giftPrice"` //礼物价格
	GiftNum   int    `json:"giftNum"`   //礼物个数
	GiftTotal int    `json:"giftTotal"` //礼物个数
}

type RoomDashBoardRes struct {
	TodayTimes    string `json:"todayTimes"`    // 今日开播时长
	TodayMoneying int    `json:"todayMoneying"` //今日实时流水
	WeekTimes     string `json:"weekTimes"`     //本周开播时长
	WeekMoneying  int    `json:"weekMoneying"`  //本周实时流水
	MonthTimes    string `json:"monthTimes"`    //本月开播时长
	MonthMoneying int    `json:"monthMoneying"` //本月实时流水
}

type RoomDashBoardChartRes struct {
	Date string `json:"date"` //时间节点
	Data int    `json:"data"` //值
}
