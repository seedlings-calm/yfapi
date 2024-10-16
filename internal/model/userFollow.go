package model

import (
	"time"
)

// User关注表
type UserFollow struct {
	Id             int       `json:"id" gorm:"column:id"`
	UserID         string    `json:"userId" gorm:"column:user_id"`                  // 用户id
	FocusUserID    string    `json:"focusUserId" gorm:"column:focus_user_id"`       // 被关注用户id
	IsMutualFollow bool      `json:"isMutualFollow" gorm:"column:is_mutual_follow"` // 是否互相关注，可以用来表示是否双方都关注对方(0:否，1:是)
	CreateTime     time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime     time.Time `json:"UpdateTime" gorm:"column:update_time"`
}

func (m *UserFollow) TableName() string {
	return "t_user_follow"
}
