package model

import "time"

// 用户银行卡
type UserBank struct {
	Id         int       `json:"id" gorm:"column:id"`
	UserId     string    `json:"userId" gorm:"column:user_id"`         // 用户id
	BankName   string    `json:"bankName" gorm:"column:bank_name"`     // 银行卡名称
	BankNo     string    `json:"bankNo" gorm:"column:bank_no"`         // 银行卡号
	BankHolder string    `json:"bankHolder" gorm:"column:bank_holder"` // 持卡人
	BankBranch string    `json:"bankBranch" gorm:"column:bank_branch"` // 支行
	BankCode   string    `json:"bankCode" gorm:"column:bank_code"`     // 银行缩写
	IsDefault  int       `json:"isDefault" gorm:"column:is_default"`   // 是否默认 0-否 1-是
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *UserBank) TableName() string {
	return "t_user_bank"
}
