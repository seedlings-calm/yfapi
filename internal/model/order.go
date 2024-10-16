package model

import "time"

type Order struct {
	ID              int64     `json:"id" gorm:"column:id"`
	OrderId         string    `json:"orderId" gorm:"column:order_id"`                 // 订单号
	UserId          string    `json:"userId" gorm:"column:user_id"`                   // 用户ID
	ToUserIdList    string    `json:"toUserIdList" gorm:"column:to_user_id_list"`     // 打赏用户列表 以逗号分隔
	RoomId          string    `json:"roomId" gorm:"column:room_id"`                   // 房间ID
	GuildId         string    `json:"guildId" gorm:"column:guild_id"`                 // 公会ID
	Gid             string    `json:"gid" gorm:"column:gid"`                          // 关联ID 物品id 礼物id等
	TotalAmount     string    `json:"totalAmount" gorm:"column:total_amount"`         // 总金额
	PayAmount       string    `json:"payAmount" gorm:"column:pay_amount"`             // 实际金额
	DiscountsAmount string    `json:"discountsAmount" gorm:"column:discounts_amount"` // 优惠金额
	Num             int       `json:"num" gorm:"column:num"`                          // 购买数量
	Currency        string    `json:"currency" gorm:"column:currency"`                // 使用币种
	AppId           string    `json:"appId" gorm:"column:app_id"`                     // 代理app_id
	OrderType       int       `json:"orderType" gorm:"column:order_type"`             // 订单类型
	OrderStatus     int       `json:"orderStatus" gorm:"column:order_status"`         // 订单状态 0未完成订单 1已完成订单
	PayType         int       `json:"payType" gorm:"column:pay_type"`                 // 支付方式
	PayStatus       int       `json:"payStatus" gorm:"column:pay_status"`             // 支付状态 0待支付 1支付完成 2退款中 3退款完成
	WithdrawStatus  int       `json:"withdrawStatus" gorm:"column:withdraw_status"`   // 提现状态 0待审核 1审核拒绝 2审核通过 3打款成功 4打款失败 5退还成功 6退还失败
	OrderNo         string    `json:"orderNo" gorm:"column:order_no"`                 // 三方订单号
	Note            string    `json:"note" gorm:"column:note"`                        // 备注信息 礼物打赏填写礼物名称
	StatDate        string    `json:"statDate" gorm:"column:stat_date"`               // 统计日期
	CreateTime      time.Time `json:"createTime" gorm:"column:create_time"`           // 创建时间
	UpdateTime      time.Time `json:"updateTime" gorm:"column:update_time"`           // 更新时间
}

func (m *Order) TableName() string {
	return "t_order"
}
