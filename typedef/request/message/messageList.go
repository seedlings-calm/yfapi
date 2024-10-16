package message

type MessageListReq struct {
	SelfUserId  string `json:"selfUserId"`
	OtherUserId string `json:"otherUserId"  form:"otherUserId" validate:"required"`
	Limit       int    `json:"limit"  form:"limit" validate:"required"`
	Timestamp   int    `json:"timestamp" form:"timestamp"`
}
