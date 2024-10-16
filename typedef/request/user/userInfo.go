package user

type UserInfoReq struct {
	UserId string `json:"userId" form:"userId" validate:"required"`
}

type UserCreateAccountReq struct {
	Avatar   string `json:"avatar" validate:"required"`
	Nickname string `json:"nickname" validate:"required,min=1,excludesall=@<>/ "`
	Password string `json:"password"`
}

type SwitchUserAccountReq struct {
	Token string `json:"token"`
}

type TimelineFilterReq struct {
	UserId string `json:"userId" validate:"required"`          //目标用户ID
	Types  int    `json:"types" validate:"required,oneof=1 2"` //1不让他看动态 2不看他的动态
}

type GetTimelineFilterListReq struct {
	Page  int `json:"page" form:"page" validate:"omitempty"`            //页码
	Size  int `json:"size" form:"size" validate:"required,min=1"`       //每页条数
	Types int `json:"types" form:"types" validate:"required,oneof=1 2"` //1不让他看动态 2不看他的动态
}

type NoticeFilterReq struct {
	UserId string `json:"userId" validate:"required"`          //目标用户ID
	Types  int    `json:"types" validate:"required,oneof=1 2"` //1动态通知 2直播通知
}
