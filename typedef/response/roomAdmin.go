package response

import (
	"time"
)

type RoomAdminListRes struct {
	ID         int       `json:"id" gorm:"column:id"`
	UserID     int64     `json:"userId" gorm:"column:user_id"` // 管理员id
	RoomID     int64     `json:"roomId" gorm:"column:room_id"` // 房间id
	RoomNo     string    `json:"roomNo" gorm:"column:room_no"` // 房间号
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
}
