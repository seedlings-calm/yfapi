package model

import "time"

type UserDefaultAvatar struct {
	Id         int64     `json:"id" gorm:"column:id"`
	Avatar     string    `json:"avatar" gorm:"column:avatar"`          // 头像
	Sex        int       `json:"sex" gorm:"column:sex"`                // 性别;0:保密,1:男,2:女
	Status     int       `json:"status" gorm:"column:status"`          // 状态1:开启,2:禁用
	StaffName  string    `json:"staffName" gorm:"column:staff_name"`   // 操作人
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"` // 创建时间
}

func (m *UserDefaultAvatar) TableName() string {
	return "t_user_default_avatar"
}
