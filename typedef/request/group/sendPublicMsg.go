package group

type SendTextMsgReq struct {
	Content  string `json:"content"  validate:"required,min=1,max=200"` //发送内容
	ToUserId string `json:"toUserId"`                                   //接受方用户id
	RoomId   string `json:"roomId" validate:"required"`
	Extra    string `json:"extra"`
}

type SendImgMsgReq struct {
	Content  string `json:"content" validate:"required"` //发送内容
	ToUserId string `json:"toUserId"`                    //接收方用户id
	Width    int    `json:"width" validate:"required"`   //宽度
	Height   int    `json:"height" validate:"required"`  //长度
	RoomId   string `json:"roomId" validate:"required"`
	Extra    string `json:"extra"`
}
