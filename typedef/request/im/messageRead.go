package im

type MessageReadReq struct {
	ChatUserId string `json:"chatUserId"  validate:"required"` //会话对象的userId
}
