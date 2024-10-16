package model

import "time"

type Agent struct {
	ID                         string    `json:"id" gorm:"column:id"`
	AppID                      string    `json:"app_id" gorm:"column:app_id"`                                             // app_id
	Secret                     string    `json:"secret" gorm:"column:secret"`                                             // api密钥
	AgentName                  string    `json:"agent_name" gorm:"column:agent_name"`                                     // 代理名称
	Status                     int       `json:"status" gorm:"column:status"`                                             // 0待审核 1审核通过 2作废
	Amount                     float64   `json:"amount" gorm:"column:amount"`                                             // 账户余额
	Username                   string    `json:"username" gorm:"column:username"`                                         // 账户
	Password                   string    `json:"password" gorm:"column:password"`                                         // 密码
	Salt                       string    `json:"salt" gorm:"column:salt"`                                                 // 盐
	PayDiscount                int       `json:"pay_discount" gorm:"column:pay_discount"`                                 // 支付折扣
	Memo                       string    `json:"memo" gorm:"column:memo"`                                                 // 备注
	CompanyName                string    `json:"company_name" gorm:"column:company_name"`                                 // 公司名称
	LegalRepresentativeName    string    `json:"legal_representative_name" gorm:"column:legal_representative_name"`       // 法人姓名
	LegalRepresentativeContact string    `json:"legal_representative_contact" gorm:"column:legal_representative_contact"` // 法人联系方式
	LegalRepresentativeIDNo    string    `json:"legal_representative_id_no" gorm:"column:legal_representative_id_no"`     // 法人身份证号
	CompanyQualificationImage  string    `json:"company_qualification_image" gorm:"column:company_qualification_image"`   // 公司资质证件图片
	CreateTime                 time.Time `json:"create_time" gorm:"column:create_time"`                                   // 创建时间
	UpdateTime                 time.Time `json:"update_time" gorm:"column:update_time"`                                   // 更新时间
}

func (m *Agent) TableName() string {
	return "t_agent"
}
