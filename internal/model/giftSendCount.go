package model

import "time"

// GiftSendCount
// @Description: 礼物赠送数组
type GiftSendCount struct {
	ID         int64     `json:"id" gorm:"column:id"`                   // 主键
	SendCount  int       `json:"send_count" gorm:"column:send_count"`   // 赠送数量
	Desc       string    `json:"desc" gorm:"column:desc"`               // 描述
	StaffName  string    `json:"staff_name" gorm:"column:staff_name"`   // 操作人
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"` // 更新时间
}

func (m *GiftSendCount) TableName() string {
	return "t_gift_send_count"
}
