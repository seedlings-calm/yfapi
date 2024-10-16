package model

import "time"

type UserLevelMaxConfig struct {
	ID         int       `json:"id" gorm:"column:id"`              // 主键
	MaxLevel   int       `json:"maxLevel" gorm:"column:max_level"` // 最大等级
	Types      int8      `json:"types" gorm:"column:types"`        // 1:vip等级 2:lv等级  3:星光等级
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *UserLevelMaxConfig) TableName() string {
	return "t_user_level_max_config"
}
