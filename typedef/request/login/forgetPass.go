package login

type ForgetPassReq struct {
	Id              string `json:"id" validate:"required"`              //选择的对应id
	ChooseUserToken string `json:"chooseUserToken" validate:"required"` //选择用户的token
	Password        string `json:"password" validate:"required,min=8,max=20"`
}
