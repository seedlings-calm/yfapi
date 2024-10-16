package model

import (
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreConfig"
)

type GuildMember struct {
	ID            int          `json:"id" gorm:"column:id"`                    // 公会成员的唯一标识
	GuildID       string       `json:"guildId" gorm:"column:guild_id"`         // 工会id
	UserID        string       `json:"userId" gorm:"column:user_id"`           // 成员的用户ID
	Status        int8         `json:"status" gorm:"column:status"`            // 成员状态， 1=正常 2=冻结 3=已脱离
	CreateTime    time.Time    `json:"createTime" gorm:"column:create_time"`   // 成员加入公会的时间
	UpdateTime    time.Time    `json:"updateTime" gorm:"column:update_time"`   // 成员信息最后更新时间
	LeaveTime     time.Time    `json:"leaveTime" gorm:"column:leave_time"`     // 成员离开公会的时间，如果未离开则为空
	LeaveReason   string       `json:"leaveReason" gorm:"column:leave_reason"` // 成员离开公会的原因，如果未离开则为空
	Users         UserBaseInfo `gorm:"foreignKey:UserID"`
	CreateTimeStr string       `json:"createTimeStr"` // 成员加入公会的时间
	UpdateTimeStr string       `json:"updateTimeStr"` // 成员加入公会的时间
}

type UserBaseInfo struct {
	Id         string `json:"id" gorm:"column:id"`
	UserNo     string `json:"user_no" gorm:"column:user_no"`         // 展示的用户id
	Nickname   string `json:"nickname" gorm:"column:nickname"`       // 昵称
	RegionCode string `json:"region_code" gorm:"column:region_code"` // 手机区号
	Mobile     string `json:"mobile" gorm:"column:mobile"`           // 手机号
	Status     int    `json:"status" gorm:"column:status"`           // 用户状态 1正常 2冻结 3申请注销 4已注销
	Avatar     string `json:"avatar" gorm:"column:avatar"`           // 头像
	TrueName   string `json:"true_name" gorm:"column:true_name"`     // 真实姓名
}

func (m *GuildMember) TableName() string {
	return "t_guild_member"
}
func (m *GuildMember) AfterFind(tx *gorm.DB) (err error) {
	if m.Users.Avatar != "" {
		m.Users.Avatar = coreConfig.GetHotConf().ImagePrefix + m.Users.Avatar
	}
	m.CreateTimeStr = m.CreateTime.Format("2006-01-02 15:04:05")
	m.UpdateTimeStr = m.UpdateTime.Format("2006-01-02 15:04:05")
	return
}

func (m *UserBaseInfo) TableName() string {
	return "t_user"
}
