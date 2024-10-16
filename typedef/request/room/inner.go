package room

type DoRoomReq struct {
	RoomId     string `json:"roomId"` //房间ID
	UserId     string `json:"userId"` //用户ID
	ClientType string `json:"clientType"`
}

type SendRoomMsgReq struct {
	Code   int    `json:"code"`   // im code
	Msg    any    `json:"msg"`    //消息内容
	RoomId string `json:"roomId"` //房间ID
}
