package model

import "time"

type UserLoginRecord struct {
	ID            int64     `json:"id" gorm:"column:id"`
	UserId        string    `json:"userId" gorm:"column:user_id"`               //用户id
	LoginPlatform string    `json:"loginPlatform" gorm:"column:login_platform"` // 登录设备
	LoginModel    string    `json:"loginModel" gorm:"column:login_model"`       // 登录设备
	ClientVersion string    `json:"clientVersion" gorm:"column:client_version"` // 客户端版本号
	DeviceID      string    `json:"deviceId" gorm:"column:device_id"`           // 设备号
	LoginIp       string    `json:"loginIp" gorm:"column:login_ip"`             // 登录ip
	Address       string    `json:"address" gorm:"column:address"`              // 登录地址
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *UserLoginRecord) TableName() string {
	return "t_user_login_record"
}
