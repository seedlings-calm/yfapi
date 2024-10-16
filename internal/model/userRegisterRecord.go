package model

import "time"

type UserRegisterRecord struct {
	ID               int       `json:"id" gorm:"column:id"`
	UserId           string    `json:"userId" gorm:"column:user_id"`                     //用户id
	RegisterPlatform string    `json:"registerPlatform" gorm:"column:register_platform"` // 注册设备
	SignType         int8      `json:"signType" gorm:"column:sign_type"`                 // 注册方式 1:手机号密码  2:验证码 3:H5手机号或验证码
	RegisterChannel  string    `json:"registerChannel" gorm:"column:register_channel"`   // 注册渠道
	CreateTime       time.Time `json:"createTime" gorm:"column:create_time"`             // 创建时间
	UpdateTime       time.Time `json:"updateTime" gorm:"column:update_time"`             // 修改时间
}

func (m *UserRegisterRecord) TableName() string {
	return "t_user_register_record"
}
