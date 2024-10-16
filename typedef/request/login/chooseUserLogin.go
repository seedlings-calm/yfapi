package login

type ChooseUserLoginReq struct {
	ChooseUserToken string `json:"chooseUserToken"  validate:"required"` //对用的选择用户token
	Id              string `json:"id"  validate:"required"`              //所选择的用户id
}
