package model

import "time"

type UserVipConfig struct {
	ID                   int       `json:"id" gorm:"column:id"`
	LevelName            string    `json:"levelName" gorm:"column:level_name"`                       // 等级名称
	Level                int       `json:"level" gorm:"column:level"`                                // 级别
	Icon                 string    `json:"icon" gorm:"column:icon"`                                  // 图标
	LogoIcon             string    `json:"logoIcon" gorm:"column:logo_icon"`                         // logo图标
	MinExperience        int       `json:"minExperience" gorm:"column:min_experience"`               // 最小经验
	MaxExperience        int       `json:"maxExperience" gorm:"column:max_experience"`               // 最高经验
	RelegationExperience int       `json:"relegationExperience" gorm:"column:relegation_experience"` // 保级经验
	FirstOperator        string    `json:"firstOperator" gorm:"column:first_operator"`               // 创建人
	LastOperator         string    `json:"lastOperator" gorm:"column:last_operator"`                 // 最新操作人
	CreateTime           time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime           time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *UserVipConfig) TableName() string {
	return "t_user_vip_config"
}
