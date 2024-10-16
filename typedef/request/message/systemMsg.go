package message

type SystemMsgReq struct {
	Title     string   `json:"title" form:"title"`
	Img       string   `json:"img"  form:"img"`
	Content   string   `json:"content"  form:"content"`
	Link      string   `json:"link"  form:"link" `
	H5Content string   `json:"h5Content"  form:"h5Content" `
	ToUserId  []string `json:"toUserId" form:"toUserId"`
}

type SendCommonMsgReq struct {
	FromUserId string   `json:"fromUserId" form:"fromUserId"`
	ToUserId   []string `json:"toUserId" form:"toUserId"`
	RoomId     string   `json:"roomId" form:"roomId"`
	MsgType    string   `json:"msgType" form:"msgType"`
	MsgData    any      `json:"msgData" form:"msgData"`
	Code       int      `json:"code" form:"code"`
}
