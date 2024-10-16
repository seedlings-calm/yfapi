package model

import "time"

type UserCollect struct {
	Id         int       `json:"id" gorm:"column:id"`
	UserID     string    `json:"userId" gorm:"column:user_id"`       // 用户id
	RoomId     string    `json:"roomId" gorm:"column:room_id"`       // 收藏房间ID
	IsCollect  int       `json:"isCollect" gorm:"column:is_collect"` // 是否收藏 1：收藏 2：否
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"UpdateTime" gorm:"column:update_time"`
}

func (m *UserCollect) TableName() string {
	return "t_user_collect"
}
