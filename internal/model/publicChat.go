package model

type PublicChat struct {
	MessageId  string  `json:"messageId" gorm:"column:message_id"` //消息唯一id
	FromUserId string  `json:"fromUserId" gorm:"column:from_user_id"`
	ToUserId   string  `json:"toUserId" gorm:"column:to_user_id"`
	RoomId     string  `json:"roomId" gorm:"column:room_id"`
	Message    string  `json:"message" gorm:"column:message"`
	Timestamp  int64   `json:"timestamp" gorm:"column:timestamp"` //毫秒时间戳
	RiskData   *string `json:"riskData" gorm:"column:risk_data"`
	Status     int     `json:"status" gorm:"column:status"` //0 正常 1可疑 2违规 3删除
}

func (data *PublicChat) TableName() string {
	return "message_store.public_chat"
}
