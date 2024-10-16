package model

import "time"

type SubsidyConfigGuild struct {
	ID          int       `json:"id" gorm:"column:id"`
	SubsidyType int       `json:"subsidy_type" gorm:"column:subsidy_type"` // 1公会月流水补贴 2公会有效直播间数量补贴 3公会直播间考核有效开播时长 4公会直播间考核有效开播天数 5公会直播间考核月流水
	ProfitNum   int       `json:"profit_num" gorm:"column:profit_num"`     // 考核数量
	ProfitRate  string    `json:"profit_rate" gorm:"column:profit_rate"`   // 补贴比例(百分比)
	StaffName   string    `json:"staff_name" gorm:"column:staff_name"`     // 操作人
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime  time.Time `json:"update_time" gorm:"column:update_time"`
}

func (m *SubsidyConfigGuild) TableName() string {
	return "t_subsidy_config_guild"
}
