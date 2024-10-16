package im

type SendOneTextMsgReq struct {
	Content  string `json:"content"  validate:"required,min=1,max=500"` //发送内容
	ToUserId string `json:"toUserId"  validate:"required"`              //接受方用户id
	Extra    string `json:"extra"`
}

type SendOneImgMsgReq struct {
	Content  string `json:"content" validate:"required"`  //发送内容
	ToUserId string `json:"toUserId" validate:"required"` //接收方用户id
	Width    int    `json:"width" validate:"required"`    //宽度
	Height   int    `json:"height" validate:"required"`   //长度
	Extra    string `json:"extra"`
}

type SendOneAudioReq struct {
	Content  string `json:"content" validate:"required"`  //发送内容
	ToUserId string `json:"toUserId" validate:"required"` //接收方用户id
	Length   int    `json:"length" validate:"required"`   //音频秒数
	Extra    string `json:"extra"`
}
