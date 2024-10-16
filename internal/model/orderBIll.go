package model

import "time"

type OrderBill struct {
	ID           int64     `json:"id" gorm:"column:id"`
	OrderId      string    `json:"orderId" gorm:"column:order_id"`             // 订单号
	UserId       string    `json:"userId" gorm:"column:user_id"`               // 用户ID
	FromUserId   string    `json:"fromUserId" gorm:"column:from_user_id"`      // 打赏用户
	ToUserIdList string    `json:"ToUserIdList" gorm:"column:to_user_id_list"` // 打赏用户列表
	Gid          string    `json:"gid" gorm:"column:gid"`                      // 关联礼物id或物品id
	Num          int       `json:"num" gorm:"column:num"`                      // (打赏)数量
	Diamond      int       `json:"diamond" gorm:"column:diamond"`              // (打赏)钻石
	RoomId       string    `json:"roomId" gorm:"column:room_id"`               // 房间ID
	GuildId      string    `json:"guildId" gorm:"column:guild_id"`             // 公会ID
	Currency     string    `json:"currency" gorm:"column:currency"`            // 币种
	FundFlow     int       `json:"fundFlow" gorm:"column:fund_flow"`           // 1入账 2出账
	BeforeAmount string    `json:"beforeAmount" gorm:"column:before_amount"`   // 变动前余额
	Amount       string    `json:"amount" gorm:"column:amount"`                // 变动金额
	CurrAmount   string    `json:"currAmount" gorm:"column:curr_amount"`       // 当前余额
	AppId        string    `json:"appId" gorm:"column:app_id"`                 // 代理APPID
	OrderType    int       `json:"orderType" gorm:"column:order_type"`         // 订单类型
	Note         string    `json:"note" gorm:"column:note"`                    // 备注信息 礼物打赏填写礼物名称
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *OrderBill) TableName() string {
	return "t_order_bill"
}

type OrderBillDTO struct {
	OrderBill
	FromUserId string `json:"fromUserId"`
	RoomName   string `json:"roomName"`
}
