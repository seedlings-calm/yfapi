package model

import "time"

type UserMuteList struct {
	ID        int       `json:"id" gorm:"column:id"`                 // 主键编码
	FromID    string    `json:"from_id" gorm:"column:from_id"`       // 操作用户ID
	ToID      string    `json:"to_id" gorm:"column:to_id"`           // 被禁言用户ID
	RoomID    string    `json:"room_id" gorm:"column:room_id"`       // 房间ID
	UnsealID  string    `json:"unseal_id" gorm:"column:unseal_id"`   // 解封操作ID
	StartTime time.Time `json:"start_time" gorm:"column:start_time"` // 禁言开始时间
	EndTime   time.Time `json:"end_time" gorm:"column:end_time"`     // 禁言结束时间
}

func (m *UserMuteList) TableName() string {
	return "t_user_mutelist"
}
