package model

import "time"

// Guild表 结构体  Guild
type Guild struct {
	ID           string    `json:"id" gorm:"column:id"`                       // 主键 Guild
	GuildNo      string    `json:"guildNo" gorm:"column:guild_no"`            // 公会号
	Name         string    `json:"name" gorm:"column:name"`                   // 公会名称
	LogoImg      string    `json:"logoImg" gorm:"column:logo_img"`            // 公会LOGO
	Status       int8      `json:"status" gorm:"column:status"`               // 公会状态 1正常 2待审核 3审核拒绝 4废弃
	BriefDesc    string    `json:"briefDesc" gorm:"column:brief_desc"`        // 公会介绍
	UserID       string    `json:"userId" gorm:"column:user_id"`              // 会长id
	SortNo       int       `json:"sortNo" gorm:"column:sort_no"`              // 排序序号
	ViewForApply int8      `json:"viewForApply" gorm:"column:view_for_apply"` // 显示在申请列表 1=显示  0=不显示
	RoomMax      int       `json:"roomMax" gorm:"column:room_max"`            // 公会房间数上限
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`      // 创建时间
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time"`      // 更新时间
}

// TableName
func (g *Guild) TableName() string {
	return "t_guild"
}
