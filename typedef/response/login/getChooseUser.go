package login

type GetChooseUserRes struct {
	Id       string `json:"id"`       //对应的用户id
	Nickname string `json:"nickname"` //用户昵称
	Avatar   string `json:"avatar"`   //用户头像
}
