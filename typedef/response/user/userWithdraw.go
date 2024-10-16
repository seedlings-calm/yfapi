package user

type UserWithdrawRes struct {
	StarlightAmount     string      `json:"starlightAmount"`     // 星光总余额
	StarlightUnWithdraw string      `json:"starlightUnWithdraw"` // 不可提现星光余额
	StarlightWithdraw   string      `json:"starlightWithdraw"`   // 可提现星光余额
	WithdrawRate        int         `json:"withdrawRate"`        //打赏收入提现手续费(百分比)
	StarlightSubsidy    string      `json:"starlightSubsidy"`    // 补贴星光余额
	SubsidyRate         int         `json:"subsidyRate"`         // 结算收入提现手续费(百分比)
	WithdrawDesc        string      `json:"withdrawDesc"`        //可提现金额说明
	UnWithdrawDesc      string      `json:"unWithdrawDesc"`      //不可提现金额说明
	Desc                string      `json:"desc"`                //提现说明
	WithdrawDays        string      `json:"withdrawDays"`        //可提现日期
	BankList            []*BankInfo `json:"bankList"`            //银行卡列表
}
type BankInfo struct {
	Id         int    `json:"id"`         //银行卡id
	BankName   string `json:"bankName"`   //银行名称
	BankNo     string `json:"bankNo"`     //银行卡号
	BankHolder string `json:"bankHolder"` //持卡人
	IsDefault  int    `json:"isDefault"`  //是否默认 1:是 2:否
}
