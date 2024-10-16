package model

type UserPrivilegeSetting struct {
	ID                         int64  `json:"id" gorm:"column:id"`
	UserID                     string `json:"user_id" gorm:"column:user_id"`
	HiddenRoomStatus           int    `json:"hidden_room_status" gorm:"column:hidden_room_status"`                       // 隐藏房间状态 0关闭 1开启
	HiddenMsgReadStatus        int    `json:"hidden_msg_read_status" gorm:"column:hidden_msg_read_status"`               // 隐藏消息已读状态 0关闭 1开启
	HiddenJoinRoom             int    `json:"hidden_join_room" gorm:"column:hidden_join_room"`                           // 隐身进厅 0关闭 1开启
	RoomPreventMuteStatus      int    `json:"room_prevent_mute_status" gorm:"column:room_prevent_mute_status"`           // 厅内防禁言 0关闭 1开启
	RoomPreventKickoutStatus   int    `json:"room_prevent_kickout_status" gorm:"column:room_prevent_kickout_status"`     // 厅内防踢 0关闭 1开启
	RoomPreventBlacklistStatus int    `json:"room_prevent_blacklist_status" gorm:"column:room_prevent_blacklist_status"` // 厅内防拉黑 0关闭 1开启
	HiddenRoomlistStatus       int    `json:"hidden_roomlist_status" gorm:"column:hidden_roomlist_status"`               // 房间榜单隐身 0关闭 1开启
	HiddenFootprintStatus      int    `json:"hidden_footprint_status" gorm:"column:hidden_footprint_status"`             // 隐藏足迹 0关闭 1开启
}

func (m *UserPrivilegeSetting) TableName() string {
	return "t_user_privilege_setting"
}
