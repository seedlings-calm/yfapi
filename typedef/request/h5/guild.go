package h5

type JoinGuildReq struct {
	GuildId string `json:"guildId" validate:"required"` //公会id
}
type GuildInfoReq struct {
	GuildId string `json:"guildId"` //公会id
}
type GuildMemberListReq struct {
	GuildId string `json:"guildId" validate:"required"` //公会id
}

// QuitGuildApplyReq
// @Description: 退出公会申请
type QuitGuildApplyReq struct {
	GuildId  string `json:"guildId" validate:"required"` // 公会ID
	IsForced bool   `json:"isForced"`                    // 是否强制退出
}
