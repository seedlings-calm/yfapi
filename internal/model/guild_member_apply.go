package model

import "time"

// GuildMember 表示公会成员申请的结构体，包含成员在公会中的申请信息。
type GuildMemberApply struct {
	ID         int       `json:"id" gorm:"column:id;primaryKey"`     //主键id
	GuildID    string    `json:"guildId" gorm:"column:guild_id"`     //工会id
	UserID     string    `json:"userId" gorm:"column:user_id"`       //用户ID
	ApplyType  int       `json:"applyType" gorm:"column:apply_type"` // 1=入会申请, 2=退会申请
	Force      int       `json:"force" gorm:"column:force"`          // 1=强制申请, 0=非强制申请
	Status     int       `json:"status" gorm:"column:status"`        // 1=待审核, 2=同意, 3=拒绝, 4=自动拒绝, 5=强制申请自动退出, 6=取消申请
	Reason     string    `json:"reason" gorm:"column:reason"`        // 拒绝原因
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (gma *GuildMemberApply) TableName() string {
	return "t_guild_member_apply"
}
