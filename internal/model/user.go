package model

import (
	"database/sql"
	"time"
	i18n_err "yfapi/i18n/error"
	typedef_enum "yfapi/typedef/enum"
)

type User struct {
	Id                  string         `json:"id" gorm:"column:id"`
	UserNo              string         `json:"user_no" gorm:"column:user_no"`                             // 展示的用户id
	OriUserNo           string         `json:"ori_user_no" gorm:"column:ori_user_no"`                     // 原始用户id
	Nickname            string         `json:"nickname" gorm:"column:nickname"`                           // 昵称
	RegionCode          string         `json:"region_code" gorm:"column:region_code"`                     // 手机区号
	Mobile              string         `json:"mobile" gorm:"column:mobile"`                               // 手机号
	Password            string         `json:"password" gorm:"column:password"`                           // 密码
	UserSalt            string         `json:"user_salt" gorm:"column:user_salt"`                         // 加密盐
	Status              int            `json:"status" gorm:"column:status"`                               // 用户状态 1正常 2冻结 3申请注销 4已注销
	Avatar              string         `json:"avatar" gorm:"column:avatar"`                               // 头像
	Sex                 int            `json:"sex" gorm:"column:sex"`                                     // 性别;0:保密,1:男,2:女
	SexEditNum          int            `json:"sex_edit_num" gorm:"column:sex_edit_num"`                   //性别修改次数
	TrueName            string         `json:"true_name" gorm:"column:true_name"`                         // 真实姓名
	BornDate            sql.NullString `json:"born_date" gorm:"column:born_date"`                         // 出生日期
	VoiceUrl            string         `json:"voice_url" gorm:"column:voice_url"`                         // 语音地址
	VoiceLength         int            `json:"voice_length" gorm:"column:voice_length"`                   // 语音长度
	Introduce           string         `json:"introduce" gorm:"column:introduce"`                         // 个人简介
	RegisterPlatform    string         `json:"register_platform" gorm:"column:register_platform"`         // 注册平台 ios,android,windows
	RegisterMachineCode string         `json:"register_machine_code" gorm:"column:register_machine_code"` // 注册设备码
	RegisterChannel     string         `json:"register_channel" gorm:"column:register_channel"`           // 注册渠道
	CreateTime          time.Time      `json:"create_time" gorm:"column:create_time"`                     // 创建时间
	Guide               int            `json:"guide" gorm:"column:guide"`                                 //1显示引导页面 2不显示引导
	AppId               string         `json:"app_id" gorm:"app_id"`
	Uid                 string         `json:"uid" gorm:"uid"`                           //第三方平台用户id
	Source              int            `json:"source" grom:"source"`                     //用户来源 0站内用户 1三方用户 2特殊用户
	VoiceStatus         int            `json:"voiceStatus" gorm:"column:voice_status"`   //语音 1正常 2待审核 3审核不通过
	RealNameStatus      int            `json:"real_name_status" gorm:"real_name_status"` //实名状态 1未认证 2已认证 3认证中
}

func (m *User) TableName() string {
	return "t_user"
}

func (m *User) CheckUserStatus() (errCode i18n_err.ErrCode) {
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
