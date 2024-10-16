package user

import "time"

type UserBankInfo struct {
	Id         int       `json:"id"`
	UserId     string    `json:"userId"`      // 用户id
	BankName   string    `json:"bankName"`    // 银行卡名称
	BankNo     string    `json:"bankNo"`      // 银行卡号
	BankHolder string    `json:"bankHolder" ` // 持卡人
	BankBranch string    `json:"bankBranch"`  // 支行
	BankCode   string    `json:"bankCode"`    // 银行缩写
	IsDefault  int       `json:"isDefault"`   // 是否默认 0-否 1-是
	CreateTime time.Time `json:"createTime"`  // 创建时间
	UpdateTime time.Time `json:"updateTime"`
}
