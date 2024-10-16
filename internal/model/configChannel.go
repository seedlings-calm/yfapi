package model

import "time"

type ConfigChannel struct {
	Id         uint      `json:"id"         description:""`
	Platform   string    `json:"platform"   description:"平台"`               //平台
	Channel    string    `json:"channel"    description:"渠道"`               //渠道
	PayUrl     string    `json:"payUrl"     description:"充值链接"`             //充值链接
	PayMethods string    `json:"payMethods" description:"支付方法（支付宝,微信）逗号分隔"` //支付方法（支付宝,微信）逗号分隔
	Status     int       `json:"status"     description:"状态 1：有效 2：无效"`     //状态 1：有效 2：无效"
	StaffName  string    `json:"staffName"  description:"操作人"`
	CreateTime time.Time `json:"createTime" description:""`
	UpdateTime time.Time `json:"updateTime" description:""`
}

func (ConfigChannel) TableName() string {
	return "t_config_channel"
}

type ConfigDiamond struct {
	Id         uint      `json:"id"         description:""`
	Platform   string    `json:"platform"   description:"平台"`         //平台
	Keys       string    `json:"keys"       description:"商品ID"`       //商品ID
	Nums       int       `json:"nums"       description:"充值金额"`       //充值金额
	GotoNums   int       `json:"gotoNums"   description:"到账金额"`       //到账金额
	Status     int       `json:"status"     description:"状态：1开启，2关闭"` //
	StaffName  string    `json:"staffName"  description:"操作人"`
	CreateTime time.Time `json:"createTime" description:""`
	UpdateTime time.Time `json:"updateTime" description:""`
}

func (ConfigDiamond) TableName() string {
	return "t_config_diamond"
}
