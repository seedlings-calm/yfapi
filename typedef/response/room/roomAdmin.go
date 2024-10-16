package room

import (
	"time"
)

type RoomAdminRes struct {
	RoomAdmin []RoomAdminInfo
}
type RoomAdminInfo struct {
	UserId     string    `json:"userId" gorm:"column:user_id"`       // 用户id
	RoomId     string    `json:"roomId" gorm:"column:room_id"`       // 房间id
	RoomNo     string    `json:"roomNo" gorm:"column:room_no"`       // 房间号
	StaffName  string    `json:"staffName" gorm:"column:staff_name"` // 操作人
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}
