package model

import "time"

type SubsidyActionRecord struct {
	ID           string    `json:"id" gorm:"column:id"`
	OrderID      string    `json:"orderId" gorm:"column:order_id"`           // 订单号
	Action       string    `json:"action" gorm:"column:action"`              // 操作类型
	BeforeStatus string    `json:"beforeStatus" gorm:"column:before_status"` // 操作前状态
	CurrStatus   string    `json:"currStatus" gorm:"column:curr_status"`     // 操作后状态
	Memo         string    `json:"memo" gorm:"column:memo"`                  // 备注
	StaffName    string    `json:"staffName" gorm:"column:staff_name"`       // 操作人
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`     // 创建时间
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time"`     // 更新时间
}

func (m *SubsidyActionRecord) TableName() string {
	return "t_subsidy_action_record"
}
