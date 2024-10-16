package user

type PerfectInfoReq struct {
	Sex      int    `json:"sex"`                                                   //性别 1男,2女
	Avatar   string `json:"avatar"`                                                //头像地址
	Nickname string `json:"nickname" validate:"omitempty,min=1,excludesall=@<>/ "` //昵称
	BornDate string `json:"bornDate"`                                              //出生日期
}

type EditUserInfoReq struct {
	Avatar    string     `json:"avatar"`
	Nickname  string     `json:"nickname" validate:"omitempty,min=1,excludesall=@<>/ "` //昵称
	Sex       int        `json:"sex" validate:"omitempty,oneof=1 2 0"`                  //性别 1男,2女 0保密
	BornDate  string     `json:"bornDate" validate:"omitempty"`
	Voice     *UserVoice `json:"voice"`
	Introduce string     `json:"introduce" validate:"omitempty,max=50"`
}

type UserVoice struct {
	Url    string `json:"url" validate:"omitempty"`
	Length int    `json:"length" validate:"max=60"` //秒数
}
