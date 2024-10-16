package model

import "time"

type UserNoticeFilter struct {
	ID         uint64    `json:"id" gorm:"column:id"`
	UserID     string    `json:"user_id" gorm:"column:user_id"` // 用户ID
	ToID       string    `json:"to_id" gorm:"column:to_id"`     // 目标用户ID
	Types      int       `json:"types" gorm:"column:types"`     // 1动态通知 2直播通知
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

func (m *UserNoticeFilter) TableName() string {
	return "t_user_notice_filter"
}
