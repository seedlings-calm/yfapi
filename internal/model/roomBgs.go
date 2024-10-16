package model

import "time"

type RoomBgsResource struct {
	Id         uint      `json:"id"         description:""`
	Name       string    `json:"name"       description:"背景名称"`         //背景名称
	Icon       string    `json:"icon"       description:"房间背景缩略图"`      //房间背景缩略图
	Backgroud  string    `json:"backgroud"  description:"背景图"`          //背景图
	Status     int       `json:"status"     description:"1: 有效 ，2 ：无效"` //1: 有效 ，2 ：无效
	Types      int       `json:"types"      description:"1:默认，2：活动"`    //1:默认，2：活动
	StaffName  string    `json:"staffName"  description:"操作人"`
	CreateTime time.Time `json:"createTime" description:""`
	UpdateTime time.Time `json:"updateTime" description:""`
}

func (RoomBgsResource) TableName() string {
	return "t_room_bgs_resource"
}

type RoomBgs struct {
	Id         uint      `json:"id"         description:""`
	RoomId     string    `json:"roomId"     description:"房间ID"`
	TrbrId     int       `json:"trbrId"     description:"房间背景ID"`
	Day        int       `json:"day"        description:"发放天数"`
	ExpireTime time.Time `json:"expireTime" description:"到期时间"`
	StaffName  string    `json:"staffName"  description:"操作人"`
	IsUse      int       `json:"isUse"      description:"是否使用 1：未使用 2：使用中"`
	IsCron     int       `json:"isCron"      description:"是否过期通知 1：未通知 2：已通知"`
	CreateTime time.Time `json:"createTime" description:""`
	UpdateTime time.Time `json:"updateTime" description:""`
}

func (RoomBgs) TableName() string {
	return "t_room_bgs"
}
