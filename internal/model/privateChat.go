package model

type PrivateChat struct {
	MessageId  string  `json:"messageId" gorm:"column:message_id"` //消息唯一id
	UniteId    string  `json:"uniteId" gorm:"union_id"`            //两个用户ID拼接 小的在前
	FromUserId string  `json:"fromUserId" gorm:"column:from_user_id"`
	ToUserId   string  `json:"toUserId" gorm:"column:to_user_id"`
	Message    string  `json:"message" gorm:"column:message"`
	Read       int     `json:"read" gorm:"column:read"`           //1已读 2未读
	Timestamp  int64   `json:"timestamp" gorm:"column:timestamp"` //毫秒时间戳
	Type       int     `json:"type" gorm:"column:type"`           //1普通消息
	Status     int     `json:"status" gorm:"column:status"`       //0 正常 1可疑 2违规
	RiskData   *string `json:"riskData" gorm:"column:risk_data"`  //风控数据
}

func (data *PrivateChat) TableName() string {
	return "message_store.private_chat"
}
