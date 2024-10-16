package roomOwner

type RoomAccountInfoRes struct {
	UserId          string     `json:"userId"`
	Status          int        `json:"status"`          // 账户状态 1=正常 2=冻结
	CashAmount      string     `json:"cashAmount"`      // 资产总额
	TotalCashIncome string     `json:"totalCashIncome"` // 累计收益
	SettlementRate  int        `json:"settlementRate"`  // 结算费率
	Desc            string     `json:"desc"`            // 提现说明
	Mobile          string     `json:"mobile"`          // 手机号
	TrueName        string     `json:"trueName"`        // 真实姓名
	BankList        []BankInfo `json:"bankList"`
}
type BankInfo struct {
	Id         int    `json:"id"`         // 银行卡信息Id
	BankNo     string `json:"bankNo"`     // 银行卡号
	BankName   string `json:"bankName"`   // 银行名称
	BankHolder string `json:"bankHolder"` // 银行账户名
	BankBranch string `json:"bankBranch"` // 银行支行
	IsDefault  int    `json:"isDefault"`  // 是否默认 1=是 2=否
}
