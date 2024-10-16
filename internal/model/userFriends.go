package model

import "time"

// 好友表
type UserFriends struct {
	ID         int       `json:"id" gorm:"column:id"`
	UserID     string    `json:"userId" gorm:"column:user_id"`         // 用户id
	FriendID   string    `json:"friendId" gorm:"column:friend_id"`     // 好友id
	Status     int8      `json:"status" gorm:"column:status"`          // 1:已接受 2:已拒绝 3:待确认
	Remark     string    `json:"remark" gorm:"column:remark"`          // 备注
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"` // 修改时间

}

func (m *UserFriends) TableName() string {
	return "t_user_friends"
}
