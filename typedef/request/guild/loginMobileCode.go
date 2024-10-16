package guild

// LoginMobileReq 手机号登录
type LoginMobileReq struct {
	Mobile     string `json:"mobile" validate:"required,numeric"` // 手机号
	Code       string `json:"code" validate:"required,numeric"`   // 验证码
	RegionCode string `json:"regionCode" validate:"required"`     // 区号
}

// SendMobileCodeReq 发送手机验证码
type SendMobileCodeReq struct {
	Mobile     string `json:"mobile" validate:"required,numeric" msg:"不能为空"` // 手机号
	RegionCode string `json:"regionCode" validate:"required"`                // 区号
	Type       int    `json:"type"`                                          // 短信类型 9登录 11绑定银行卡
}

// GuildnfoReq 获取公会信息
type GuildnfoReq struct {
	GuildId string `json:"guildId" form:"guildId" validate:"required"` // 公会ID
}
