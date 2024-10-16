package model

import "time"

// GiftShowCategory
// @Description: 礼物显示类目
type GiftShowCategory struct {
	ID           int64     `json:"id" gorm:"column:id"`                         // 主键
	CategoryName string    `json:"category_name" gorm:"column:category_name"`   // 类目名称
	RoomLiveType int       `json:"room_live_type" gorm:"column:room_live_type"` // 使用位置(房间直播类型) 1聊天室 2直播 3个人
	SortNo       int       `json:"sort_no" gorm:"column:sort_no"`               // 排序
	StaffName    string    `json:"staff_name" gorm:"column:staff_name"`         // 操作人
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time"`       // 创建时间
	UpdateTime   time.Time `json:"update_time" gorm:"column:update_time"`       // 更新时间
}

func (m *GiftShowCategory) TableName() string {
	return "t_gift_show_category"
}
