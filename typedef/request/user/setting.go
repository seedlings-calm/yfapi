package user

type SetPasswordReq struct {
	Password string `json:"password" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

type UserRealNameReq struct {
	RealName string `json:"realName" validate:"required"` //真实姓名
	IdNo     string `json:"idNo" validate:"required"`     //身份证号
	FontUrl  string `json:"fontUrl" validate:"required"`  //身份证正面
	BackUrl  string `json:"backUrl" validate:"required"`  //身份证反面
}

// 验证手机号
type VerifyMobileReq struct {
	Code string `json:"code" validate:"required"`
}

type ChangeMobileReq struct {
	Mobile     string `json:"mobile" validate:"required"`
	Code       string `json:"code" validate:"required"`
	RegionCode string `json:"regionCode" validate:"required"` //区号
}
