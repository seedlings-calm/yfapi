package user

type PrivacySettingResp struct {
	DontSeeHeMomentsNum int64         `json:"dontSeeHeMomentsNum"` //不看他动态数量
	DontLetHeSeeMeNum   int64         `json:"dontLetHeSeeMeNum"`   //不让他看我动态数量
	BlacklistNum        int64         `json:"blacklistNum"`        //黑名单数量
	Privilege           PrivilegeData `json:"privilege"`           //特权
}

// 特权
type PrivilegeData struct {
	HiddenRoomStatus           bool `json:"hiddenRoomStatus"`           //隐藏我的聊天室状态
	HiddenMsgReadStatus        bool `json:"hiddenMsgReadStatus"`        //隐藏消息已读状态
	HiddenJoinRoomStatus       bool `json:"hiddenJoinRoom"`             //隐身进厅
	RoomPreventMuteStatus      bool `json:"roomPreventMuteStatus"`      //厅内防禁言
	RoomPreventKickOutStatus   bool `json:"roomPreventKickOutStatus"`   //厅内防踢
	RoomPreventBlacklistStatus bool `json:"roomPreventBlacklistStatus"` //厅内防拉黑
	HiddenRoomListStatus       bool `json:"hiddenRoomListStatus"`       //房间榜单隐身
	HiddenFootprintStatus      bool `json:"hiddenFootprintStatus"`      //隐藏足迹
}
