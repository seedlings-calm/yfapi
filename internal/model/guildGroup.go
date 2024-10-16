package model

import (
	"time"

	"gorm.io/gorm"
)

// GuildGroup 公会成员分组
type GuildGroup struct {
	ID            int       `json:"id" gorm:"column:id;primaryKey"`     //主键id
	GuildID       string    `json:"guildId" gorm:"column:guild_id"`     //工会id
	GroupName     string    `json:"groupName" gorm:"column:group_name"` //分组名称
	Desc          string    `json:"desc" gorm:"column:desc"`            //分组描述
	Status        int8      `json:"status" gorm:"column:status"`        // 1=正常, 2=作废
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time"`
	CreateTimeStr string    `json:"createTimeStr" gorm:"-"` // 创建时间
	UpdateTimeStr string    `json:"updateTimeStr" gorm:"-"` // 更新时间
}

func (g *GuildGroup) TableName() string {
	return "t_guild_group"
}
func (g *GuildGroup) AfterFind(db *gorm.DB) (err error) {
	g.CreateTimeStr = g.CreateTime.Format("2006-01-02 15:04:05")
	g.UpdateTimeStr = g.UpdateTime.Format("2006-01-02 15:04:05")
	return
}
