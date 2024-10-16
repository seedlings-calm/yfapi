package model

import "time"

type UserReportPunishRecord struct {
	ID           int       `json:"id" gorm:"column:id"`
	ReportId     int       `json:"reportId" gorm:"column:report_id"`         // 举报id
	DstUserId    int64     `json:"dstUserId" gorm:"column:dst_user_id"`      // 处罚对象id
	Status       int       `json:"status" gorm:"column:status"`              // 1:有效 2:失效
	Object       int       `json:"object" gorm:"column:object"`              // 1:房间 2:个人
	PunishResult string    `json:"punishResult" gorm:"column:punish_result"` // 处罚结果
	PunishNotes  string    `json:"punishNotes" gorm:"column:punish_notes"`   // 处罚说明
	StaffName    string    `json:"staffName" gorm:"column:staff_name"`       // 操作人
	ExpireTime   time.Time `json:"expireTime" gorm:"column:expire_time"`     // 过期时间
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`     // 创建时间
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time"`     // 更新时间
}

func (m *UserReportPunishRecord) TableName() string {
	return "t_user_report_punish_record"
}
