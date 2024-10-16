package model

import "time"

type GroupFilter struct {
	ID         int64     `json:"id" gorm:"column:id"`
	GroupID    string    `json:"group_id" gorm:"column:group_id"` // 群id
	UserID     string    `json:"user_id" gorm:"column:user_id"`   // 用户ID
	Types      int       `json:"types" gorm:"column:types"`       // 1全群禁言 2接收全部消息 3只接受管理消息 4不接收任何消息
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

func (m *GroupFilter) TableName() string {
	return "t_group_filter"
}
