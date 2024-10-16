package model

import "time"

// GuildMember 结构体  GuildMember
type GuildMember struct {
	ID          int        `json:"id" gorm:"column:id"`                    // 公会成员的唯一标识
	GuildID     string     `json:"guildId" gorm:"column:guild_id"`         // 工会id
	UserID      string     `json:"userId" gorm:"column:user_id"`           // 成员的用户ID
	Status      int8       `json:"status" gorm:"column:status"`            // 成员状态， 1=正常 2=冻结 3=已脱离
	StarExp     int        `json:"starExp" gorm:"column:star_exp"`         // 入会时经验值
	StarLevel   int        `json:"starLevel" gorm:"column:star_level"`     // 入会时等级
	CreateTime  time.Time  `json:"createTime" gorm:"column:create_time"`   // 成员加入公会的时间
	UpdateTime  time.Time  `json:"updateTime" gorm:"column:update_time"`   // 成员信息最后更新时间
	LeaveTime   *time.Time `json:"leaveTime" gorm:"column:leave_time"`     // 成员离开公会的时间，如果未离开则为空
	LeaveReason string     `json:"leaveReason" gorm:"column:leave_reason"` // 成员离开公会的原因，如果未离开则为空
}

// TableName
func (m *GuildMember) TableName() string {
	return "t_guild_member"
}
