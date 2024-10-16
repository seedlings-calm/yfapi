package model

type PrivateChatClearRecord struct {
	UniteId   string `json:"uniteId" gorm:"column:unite_id"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"` //毫秒时间戳
}

func (data *PrivateChatClearRecord) TableName() string {
	return "message_store.private_chat_clear_record"
}
