package model

type DicRoomNo struct {
	Id     int64  `json:"id" gorm:"column:id"`
	RoomNo string `json:"roomNo" gorm:"column:room_no"` // 生成的房间id
	Status int8   `json:"status" gorm:"column:status"`  // 是否使用
}

func (m *DicRoomNo) TableName() string {
	return "t_dic_room_no"
}
