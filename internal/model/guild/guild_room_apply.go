package model

import "time"

type GuildRoomApply struct {
	ID                int64     `json:"id" gorm:"column:id"`                                  // 主键
	GuildID           string    `json:"GuildId" gorm:"column:guild_id"`                       // 公会ID
	RoomID            string    `json:"roomId" gorm:"column:room_id"`                         // 房间ID
	RoomName          string    `json:"roomName" gorm:"column:room_name"`                     // 房间名称
	RoomDesc          string    `json:"roomDesc" gorm:"column:room_desc"`                     // 房间描述
	RoomUserID        string    `json:"roomUserId" gorm:"column:room_user_id"`                // 房主ID
	RoomAvatar        string    `json:"roomAvatar" gorm:"column:room_avatar"`                 // 厅图
	Status            int8      `json:"status" gorm:"column:status"`                          // 申请状态（1申请中，2成功，3被拒绝）
	RoomType          int       `json:"roomType" gorm:"column:room_type"`                     // 类型ID
	TemplateID        string    `json:"templateId" gorm:"column:template_id"`                 // 模板ID
	DaySettleUserID   string    `json:"daySettleUserId" gorm:"column:day_settle_user_id"`     // 日结算人ID
	MonthSettleUserID string    `json:"monthSettleUserId" gorm:"column:month_settle_user_id"` // 月结算人ID
	StaffName         string    `json:"staffName" gorm:"column:staff_name"`                   // 操作人
	Reason            string    `json:"reason" gorm:"column:reason"`                          // 拒绝原因
	CreateTime        time.Time `json:"createTime" gorm:"column:create_time"`                 // 创建时间
	UpdateTime        time.Time `json:"updateTime" gorm:"column:update_time"`                 // 更新时间
}

func (m *GuildRoomApply) TableName() string {
	return "t_guild_room_apply"
}
