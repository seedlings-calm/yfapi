package guild

type SendSMSRes struct {
}
type GuildInfo struct {
	GuildId       string `json:"guildId"`       // 公会ID
	UserAvatar    string `json:"userAvatar"`    // 会长头像
	UserNo        string `json:"userNo"`        // 会长ID
	UserName      string `json:"userName"`      // 会长昵称
	LogoImg       string `json:"logoImg"`       // 公会头像
	GuildNo       string `json:"guildNo"`       // 公会ID
	LastLoginTime string `json:"lastLoginTime"` // 上次登录时间
}
type LoginMobileCodeRes struct {
	Token    string `json:"token"`    // token
	UserID   string `json:"userId"`   // 会长Id
	UserName string `json:"userName"` // 会长昵称
	Guild    string `json:"guild"`    // 公会名称
	GuildID  string `json:"guildId"`  // 公会ID
}
