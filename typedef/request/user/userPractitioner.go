package user

type PractitionerAwnserReq struct {
	Types   int      `json:"types" form:"types" validate:"required,min=1"` // "考核类别"：1：主持人，2：咨询师，3：主播
	Answers []Answer `json:"awnsers" form:"awnsers" validate:"dive,required"`
}
type Answer struct {
	Id      int    `json:"id" form:"id" validate:"required"`         //题目ID
	Awnsers string `json:"awnser" form:"awnser" validate:"required"` //题目答案 单选直接赋值，多选，使用英文逗号分隔答案
}

type PractitionerShortAwnserReq struct {
	Types   int           `json:"types" form:"types" validate:"required,min=1"` // "考核类别"：1：主持人，2：咨询师，3：主播
	Answers []ShortAnswer `json:"awnsers" form:"awnsers" validate:"dive,required"`
}

type ShortAnswer struct {
	QuestionName string `json:"questionName" form:"questionName" validate:"required"` //题目
	Awnsers      string `json:"awnser" form:"awnser" validate:"required"`             //题目答案 :语音路径
}

type PractitionerMusicianReq struct {
	Image       string `json:"image" form:"image" validate:"required"`             //认证资料
	Description string `json:"description" form:"description" validate:"required"` //自我描述
	Audio       string `json:"audio" form:"audio" validate:"required"`             //语音介绍
}

type UserPractitionerReq struct {
	UserId string `json:"userId" form:"userId" validate:"required"`
	RoomId string `json:"roomId" form:"roomId" validate:"required"`
}

type RoomBlackOutReq struct {
	UserId string `json:"userId" form:"userId" validate:"required"`
	RoomId string `json:"roomId" form:"roomId" validate:"required"`
}
