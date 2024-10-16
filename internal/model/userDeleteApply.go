package model

import "time"

type UserDeleteApply struct {
	Id         int64     `json:"id" gorm:"column:id"`                  // 主键
	UserId     string    `json:"userId" gorm:"column:user_id"`         // 用户ID
	Status     int       `json:"status" gorm:"column:status"`          // 状态 1注销申请中 2注销撤回 3已注销
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"` // 更新时间
	ExpireTime time.Time `json:"expireTime" gorm:"column:expire_time"` // 到期注销时间
}

func (m *UserDeleteApply) TableName() string {
	return "t_user_delete_apply"
}
