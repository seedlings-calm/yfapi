package model

import "time"

// UserPractitionerAction
// @Description: 从业者行为记录
type UserPractitionerAction struct {
	ID         int       `json:"id" gorm:"column:id"`            // 主键
	UserId     string    `json:"userId" gorm:"column:user_id"`   // 用户id
	GuildId    string    `json:"guildId" gorm:"column:guild_id"` // 公会ID
	Action     string    `json:"action" gorm:"column:action"`    // 用户操作
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

func (m *UserPractitionerAction) TableName() string {
	return "t_user_practitioner_action"
}
