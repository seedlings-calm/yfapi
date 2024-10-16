package group

type WorldGroupSettingResp struct {
	RoomId       string `json:"roomId"`       //群ID
	MuteSwitch   bool   `json:"muteSwitch"`   //世界频道禁言状态
	NoticeStatus int    `json:"noticeStatus"` //消息通知设置 2接收全部消息 3只接收管理消息 4不接收任何消息
}
