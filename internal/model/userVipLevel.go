package model

import "time"

type UserVipLevel struct {
	ID         int       `json:"id" gorm:"column:id"`
	UserId     string    `json:"userId" gorm:"column:user_id"`         // 用户id
	Level      int       `json:"level" gorm:"column:level"`            // 等级id
	CurrExp    int       `json:"currExp" gorm:"column:curr_exp"`       // 当前经验
	KeepExp    int       `json:"keepExp" gorm:"column:keep_exp"`       // 保级经验
	ExpireTime time.Time `json:"expireTime" gorm:"column:expire_time"` // 过期时间
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *UserVipLevel) TableName() string {
	return "t_user_vip_level"
}

type UserVipLevelDTO struct {
	ID            int       `json:"id" gorm:"column:id"`
	UserId        string    `json:"userId" gorm:"column:user_id"`               // 用户id
	Level         int       `json:"level" gorm:"column:level"`                  // 等级id
	CurrExp       int       `json:"currExp" gorm:"column:curr_exp"`             // 当前经验
	MinExperience int       `json:"minExperience" gorm:"column:min_experience"` // 最小经验
	MaxExperience int       `json:"maxExperience" gorm:"column:max_experience"` // 最高经验
	Icon          string    `json:"icon" gorm:"icon"`                           // 图标
	ExpireTime    time.Time `json:"expireTime" gorm:"column:expire_time"`       // 过期时间
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time"`
}
