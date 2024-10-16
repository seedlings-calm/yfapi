package model

import "time"

// 实名认证表
type UserRealName struct {
	Id         int       `json:"id" gorm:"column:id"`
	UserId     string    `json:"useId" gorm:"column:user_id"`          // 用户id
	TrueName   string    `json:"trueName" gorm:"column:true_name"`     // 真实姓名
	IdNo       string    `json:"idNo" gorm:"column:id_no"`             // 身份证号
	FontUrl    string    `json:"fontUrl" gorm:"column:font_url"`       // 正面照片
	BackUrl    string    `json:"backUrl" gorm:"column:back_url"`       // 反面照片
	Status     int8      `json:"status" gorm:"column:status"`          // 审核状态 1:待审核 2:审核通过 3:审核拒绝
	StaffName  string    `json:"StaffName" gorm:"column:staff_name"`   // 审核人
	Reason     string    `json:"reason" gorm:"column:reason"`          // 拒绝原因
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"` // 修改时间
}

func (m *UserRealName) TableName() string {
	return "t_user_real_name"
}
