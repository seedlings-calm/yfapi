package model

import "time"

type RoomType struct {
	ID         int       `json:"id" gorm:"column:id"`                  //主键id
	TypeId     int       `json:"typeId" gorm:"column:type_id"`         // 房间类型ID
	TypeName   string    `json:"typeName" gorm:"column:type_name"`     // 类型名称
	LiveType   int8      `json:"liveType" gorm:"column:live_type"`     // 房间直播类型 1=聊天室 2=个播 3=个人
	SortNum    int8      `json:"sortNum" gorm:"column:sort_num"`       // 排序
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"` // 更新时间
}

func (m *RoomType) TableName() string {
	return "t_room_type"
}
