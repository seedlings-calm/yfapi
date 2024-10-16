package user

// 登录设备日志
type LoginRecordResponse struct {
	LoginPlatform string `json:"loginPlatform" gorm:"column:login_platform"` // 登录设备
	LoginModel    string `json:"loginModel" gorm:"column:login_model"`       // 登录设备
	ClientVersion string `json:"clientVersion" gorm:"column:client_version"` // 客户端版本号
	DeviceID      string `json:"deviceId" gorm:"column:device_id"`           // 设备号
	Address       string `json:"address" gorm:"column:address"`              // 登录地址
	CreateTime    string `json:"createTime" gorm:"column:create_time"`
}
