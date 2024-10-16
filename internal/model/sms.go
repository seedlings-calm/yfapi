package model

import "time"

type Sms struct {
	Id         int       `json:"id" gorm:"primaryKey;autoIncrement;comment:主键编码"`
	RegionCode string    `json:"region_code" gorm:"size:100;column:region_code;comment:手机区号"` // 手机区号
	Mobile     string    `json:"mobile" gorm:"size:100;column:mobile;comment:手机号"`            // 手机号
	Types      int       `json:"types" gorm:"column:types;comment:短信类型"`                      // 短信类型
	Code       string    `json:"code" gorm:"size:100;column:code;comment:短信码"`                // 短信码
	IsUse      bool      `json:"isUse" gorm:"column:is_use;default:0;comment:是否使用 >=1 已使用"`   //是否使用  >=1 已使用
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`                        // 创建时间
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`                        // 更新时间
}

func (s Sms) TableName() string {
	return "t_sms"
}
