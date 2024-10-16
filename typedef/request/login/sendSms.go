package login

type SendSmsReq struct {
	Mobile     string `json:"mobile" validate:"required"`     //手机号
	Type       int    `json:"type" validate:"required"`       //验证码类型 1登录验证码,2重置密码,3忘记密码 4设置密码 5账号注销
	RegionCode string `json:"regionCode" validate:"required"` //区号
}
