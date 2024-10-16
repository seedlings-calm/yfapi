package model

import "time"

// 用户从业者身份
type UserPractitionerCred struct {
	Id               int64     `json:"id"               description:""`
	UserId           string    `json:"userId"           description:"用户id"`
	PractitionerType int       `json:"practitionerType" description:"从业者类型"`
	Status           int       `json:"status"           description:"状态 1正常,2冻结,3none,4取消"`
	Exp              int64     `json:"exp"              description:"当前经验"`
	FrozenTime       string    `json:"frozenTime"       description:"冻结时间"`
	FrozenBeginTime  time.Time `json:"frozenBeginTime"  description:"冻结开始时间"`
	FrozenEndTime    time.Time `json:"frozenEndTime"    description:"冻结结束时间"`
	FrozenCreateBy   string    `json:"frozenCreateBy"   description:"冻结操作人"`
	FrozenReason     string    `json:"frozenReason"     description:"冻结时间"`
	AbolishReason    string    `json:"abolishReason"    description:"取消原因"`
	CreateTime       time.Time `json:"createTime"       description:""`
	UpdateTime       time.Time `json:"updateTime"       description:""`
}

func (m *UserPractitionerCred) TableName() string {
	return "t_user_practitioner_cred"
}
