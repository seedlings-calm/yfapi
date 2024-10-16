package model

type AuthRule struct {
	ID           int    `json:"id" gorm:"column:id"`
	RuleID       int    `json:"rule_id" gorm:"column:rule_id"`           // 规则ID
	Name         string `json:"name" gorm:"column:name"`                 // 规则名唯一标识
	Title        string `json:"title" gorm:"column:title"`               // 规则中文名
	ShowName     string `json:"show_name" gorm:"column:show_name"`       // 展示名
	Category     string `json:"category" gorm:"column:category"`         // 规则类型
	CategoryName string `json:"categoryName" gorm:"column:categoryName"` // 规则类型描述
	Status       int    `json:"status" gorm:"column:status"`             // 状态 1正常 0禁用
}

func (m *AuthRule) TableName() string {
	return "t_auth_rule"
}
