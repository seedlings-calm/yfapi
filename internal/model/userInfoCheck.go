package model

import (
	"time"
)

type UserInfoCheck struct {
	ID           int64      `json:"id" gorm:"column:id"`
	UserID       string     `json:"user_id" gorm:"column:user_id"` // 用户ID
	Content      string     `json:"content" gorm:"column:content"`
	VoiceLength  int        `json:"voice_length" gorm:"column:voice_length"`
	Type         int        `json:"type" gorm:"column:type"`                   // 1头像审核 2语音审核
	AutoStatus   int        `json:"auto_status" gorm:"column:auto_status"`     // 0待审核 1审核通过 2审核不通过 3可疑
	ManualStatus int        `json:"manual_status" gorm:"column:manual_status"` // 0待审核 1审核通过 2审核不通过
	AutoResult   *string    `json:"auto_result" gorm:"column:auto_result"`     // 自动审核检测结果
	ManualResult string     `json:"manual_result" gorm:"column:manual_result"` // 人工审核结果
	Operator     int64      `json:"operator" gorm:"column:operator"`           // 操作人id
	CreateTime   *time.Time `json:"create_time" gorm:"column:create_time"`     // 创建时间
	UpdateTime   *time.Time `json:"update_time" gorm:"column:update_time"`     // 更新时间
}

func (m *UserInfoCheck) TableName() string {
	return "t_user_info_check"
}
