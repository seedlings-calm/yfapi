package model

import (
	"gorm.io/gorm"
	"yfapi/core/coreConfig"
	i18n_err "yfapi/i18n/error"
	typedef_enum "yfapi/typedef/enum"
)

type GuildUser struct {
	Id       string `json:"id" gorm:"column:id"`
	UserNo   string `json:"user_no" gorm:"column:user_no"`     // 展示的用户id
	Nickname string `json:"nickname" gorm:"column:nickname"`   // 昵称
	Mobile   string `json:"mobile" gorm:"column:mobile"`       // 手机号
	Status   int    `json:"status" gorm:"column:status"`       // 用户状态 1正常 2冻结 3申请注销 4已注销
	Avatar   string `json:"avatar" gorm:"column:avatar"`       // 头像
	TrueName string `json:"true_name" gorm:"column:true_name"` // 真实姓名
}

func (m *GuildUser) AfterFind(tx *gorm.DB) (err error) {
	if m.Avatar != "" {
		m.Avatar = coreConfig.GetHotConf().ImagePrefix + m.Avatar
	}
	return
}
func (m *GuildUser) TableName() string {
	return "t_user"
}

func (m *GuildUser) CheckUserStatus() (errCode i18n_err.ErrCode) {
	switch m.Status {
	case typedef_enum.UserStatusNormal:
	case typedef_enum.UserStatusFreezing:
		errCode = i18n_err.ErrorCodeUserFreezing
	case typedef_enum.UserStatusApplyInvalid:
		errCode = i18n_err.ErrorCodeUserApplyInvalid
	case typedef_enum.UserStatusInvalid:
		errCode = i18n_err.ErrorCodeUserInvalid
	default:
		errCode = i18n_err.ErrorCodeUserNotFound
	}
	return
}
func (m *GuildUser) GetUserInfo(db *gorm.DB, uid string) (resp GuildUser) {
	var user = GuildUser{Id: uid}
	db.First(&user)
	return user
}
