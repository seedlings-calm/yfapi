package user

type UserWithdrawApplyReq struct {
	Amount int `json:"amount" form:"amount" validate:"required"` //提现金额
	BankId int `json:"bankId" form:"bankId" validate:"required"` //银行卡信息Id
}
