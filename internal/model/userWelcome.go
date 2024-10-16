package model

import "time"

// TUserAutoWelcome 表示用户自动欢迎语的数据库表结构
type UserAutoWelcome struct {
	ID             int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                    // 主键ID
	UserID         string    `gorm:"column:user_id;not null" json:"user_id"`                          // 用户ID
	WelcomeContent string    `gorm:"column:welcome_content;size:128;not null" json:"welcome_content"` // 欢迎语内容
	State          int       `gorm:"column:state;not null" json:"state"`                              // 状态 1正常 2冻结
	StaffName      string    `gorm:"column:staff_name;size:45;not null" json:"staff_name"`            // 操作人
	CreateTime     time.Time `gorm:"column:create_time;not null" json:"create_time"`                  // 创建时间
	UpdateTime     time.Time `gorm:"column:update_time;not null" json:"update_time"`                  // 更新时间
}

// TableName 返回表名
func (UserAutoWelcome) TableName() string {
	return "t_user_auto_welcome"
}
