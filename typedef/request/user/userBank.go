package user

type UserBankAddReq struct {
	BankNo     string `json:"bankNo" form:"bankNo" validate:"required"`         //银行卡号
	BankName   string `json:"bankName" form:"bankName" validate:"required"`     //开户行
	BankHolder string `json:"bankHolder" form:"bankPhone" validate:"required"`  //开卡人
	BankBranch string `json:"bankBranch" form:"bankBranch" validate:"required"` //支行
	Mobile     string `json:"mobile" validate:"required,numeric"`               //手机号
	Code       string `json:"code" validate:"required,numeric"`                 //验证码
	RegionCode string `json:"regionCode" validate:"required"`                   //区号
}

type UserBankUnBindReq struct {
	Id   int    `json:"id" form:"id" validate:"required"`     //银行卡id
	IdNo string `json:"idNo" form:"idNo" validate:"required"` //身份证号
}
