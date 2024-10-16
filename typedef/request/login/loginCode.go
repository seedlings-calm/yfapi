package login

type LoginCodeReq struct {
	Mobile     string `json:"mobile" validate:"required"`     //手机号
	Code       string `json:"code" validate:"required"`       //验证码
	RegionCode string `json:"regionCode" validate:"required"` //区号
}

type LoginCheckReq struct {
	Mobile     string `json:"mobile"`     //手机号
	RegionCode string `json:"regionCode"` //区号
	Password   string `json:"password"`
	UserId     string `json:"userId"`
}
