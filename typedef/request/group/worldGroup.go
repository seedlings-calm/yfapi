package group

type WorldGroupNoticeSettingReq struct {
	Types int `json:"types" validate:"required,oneof=2 3 4"` //2接收全部消息 3接收管理消息 3不接收任何消息
}
