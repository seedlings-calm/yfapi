package model

import "time"

type AppMenuSetting struct {
	ID         int       `json:"id" gorm:"column:id"`
	ModuleType int       `json:"module_type" gorm:"column:module_type"` // 模块类型
	Platform   string    `json:"platform" gorm:"column:platform"`       // 使用平台 android ios pc h5
	MenuType   int       `json:"menu_type" gorm:"column:menu_type"`     // 菜单类型
	MenuName   string    `json:"menu_name" gorm:"column:menu_name"`     // 菜单名称
	Icon       string    `json:"icon" gorm:"column:icon"`               // 菜单icon
	LinkUrl    string    `json:"link_url" gorm:"column:link_url"`       // 跳转地址
	SortNo     int       `json:"sort_no" gorm:"column:sort_no"`         // 排序 正序
	Status     int8      `json:"status" gorm:"column:status"`           // 状态 1开启 2禁用
	StaffName  string    `json:"staff_name" gorm:"column:staff_name"`   // 操作人
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time"`
}

func (m *AppMenuSetting) TableName() string {
	return "t_app_menu_setting"
}
