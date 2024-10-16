package login

type ForgetPassCheckRes struct {
	Id              string `json:"id"`              //对应的用户的id
	JumpType        int    `json:"jumpType"`        //跳转页面，1修改密码，2选择账号
	ChooseUserToken string `json:"chooseUserToken"` //选择用户的token,当jumpType=2的时候该值使用
}
