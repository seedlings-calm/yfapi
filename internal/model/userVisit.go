package model

import (
	"time"
)

type UserVisit struct {
	ID           int64     `json:"id" gorm:"column:id"`
	UserId       string    `json:"userId" gorm:"column:user_id"`              // 访问人
	TargetUserId string    `json:"targetUserId" gorm:"column:target_user_id"` // 被访问人
	IsVisit      bool      `json:"isVisit" gorm:"column:is_visit"`            // 互相访问
	VisitCount   int       `json:"visitCount" gorm:"column:visit_count"`      // 访问次数
	VisitHidden  bool      `json:"visitHidden" gorm:"column:visit_hidden"`    // 隐身访问
	ClearVisit   bool      `json:"clearVisit" gorm:"column:clear_visit"`      // 个人足迹清除
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`      // 创建时间
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time"`      // 最新访问时间
}

func (m *UserVisit) TableName() string {
	return "t_user_visit"
}

type UserVisitDTO struct {
	ID           int64     `json:"id" gorm:"column:id"`
	UserId       string    `json:"userId" gorm:"column:user_id"`              // 访问人
	TargetUserId string    `json:"targetUserId" gorm:"column:target_user_id"` // 被访问人
	Nickname     string    `json:"nickname" gorm:"column:nickname"`           // (被)访问人昵称 根据查询类型 足迹是被访问人 访客是访问人
	Avatar       string    `json:"avatar" gorm:"column:avatar"`               // (被)访问人头像 根据查询类型 足迹是被访问人 访客是访问人
	Sex          int       `json:"sex" gorm:"column:sex"`                     // (被)访问人性别 根据查询类型 足迹是被访问人 访客是访问人
	Introduce    string    `json:"introduce"`                                 // 个性签名
	IsVisit      bool      `json:"isVisit" gorm:"column:is_visit"`            // 互相访问
	VisitCount   int       `json:"visitCount" gorm:"column:visit_count"`      // 访问次数
	VisitHidden  bool      `json:"visitHidden" gorm:"column:visit_hidden"`    // 隐身访问
	ClearVisit   bool      `json:"clearVisit" gorm:"column:clear_visit"`      // 个人足迹清除
	CreateTime   time.Time `json:"createTime" gorm:"column:create_time"`      // 创建时间
	UpdateTime   time.Time `json:"updateTime" gorm:"column:update_time"`      // 最新访问时间
}
