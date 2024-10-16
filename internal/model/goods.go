package model

import (
	"time"
)

type Goods struct {
	Id               uint      `json:"id"               description:""`
	Code             string    `json:"code"             description:"物品编号"`
	Name             string    `json:"name"             description:"物品名称"`
	TypeId           int       `json:"typeId"           description:"物品类型"`
	TypeKey          string    `json:"typeKey"           description:"物品类型key"`
	Icon             string    `json:"icon"             description:"图标"`
	AnimationUrl     string    `json:"animationUrl"     description:"图片动效"`
	AnimationJsonUrl string    `json:"animationJsonUrl" description:"json文件动效"`
	Money            int       `json:"money"            description:"基础价格"`
	Status           int       `json:"status"           description:"是否上架 2：下架中，1：上架中"`
	StaffName        string    `json:"staffName"        description:"操作人昵称"`
	CreateTime       time.Time `json:"createTime"       description:"创建时间"`
	UpdateTime       time.Time `json:"updateTime"       description:"修改时间"`
	IsDel            int       `json:"isDel"            description:"是否删除 1:有效 2：已删除"`
}

func (Goods) TableName() string {
	return "t_goods"
}

type GoodsUse struct {
	Id           uint      `json:"id"         description:""`
	GoodsId      int       `json:"goodsId"    description:"物品ID"`
	GoodsTypeId  int       `json:"goods_type_id" description:"物品类型ID"`
	GoodsTypeKey string    `json:"goodsTypeKey" ` // 商品类型key
	Money        int       `json:"money"      description:"单价/天"`
	Moneys       string    `json:"moneys"     description:"7/15/30天价格\"7,15,28\""`
	Status       int       `json:"status"     description:"是否上架：2：下架中，1：上架中"`
	StartTime    time.Time `json:"startTime"  description:"展示开始时间，到期自动上架"`
	EndTime      time.Time `json:"endTime"    description:"展示结束时间，到期自动下架"`
	StaffName    string    `json:"staffName"  description:"操作人昵称"`
	CreateTime   time.Time `json:"createTime" description:""`
	UpdateTime   time.Time `json:"updateTime" description:""`
	IsDel        int       `json:"isDel"      description:"是否删除 1:有效 2：已删除"`
}

func (GoodsUse) TableName() string {
	return "t_goods_use"
}

type GoodsType struct {
	Id         uint      `json:"id"         description:""`
	Name       string    `json:"name"       description:"名称"`
	Keys       string    `json:"keys"       description:"key"`
	Icon       string    `json:"icon"       description:"图标"`
	Sort       int       `json:"sort"       description:"排序"`
	Status     int       `json:"status"     description:"状态：2:关闭中，1：启用中"`
	Remark     string    `json:"remark"     description:"备注消息"`
	StaffName  string    `json:"staffName"  description:"操作人昵称"`
	CreateTime time.Time `json:"createTime" description:"创建时间"`
	UpdateTime time.Time `json:"updateTime" description:"修改时间"`
	IsDel      int       `json:"isDel"      description:"是否删除 1:有效 2：已删除"`
}

func (GoodsType) TableName() string {
	return "t_goods_type"
}

type GoodsGrant struct {
	Id         uint      `json:"id"         description:""`
	UserId     string    `json:"userId"     description:"接收用户ID"`
	GoodsId    int64     `json:"goodsId"    description:"发放商品ID"`
	Day        int       `json:"day"        description:"发放天数"`
	Source     int       `json:"source"     description:"发放来源：1：后台发放，2：vip等级发放，3：用户权益发放"`
	StaffName  string    `json:"staffName"  description:"发放人"`
	CreateTime time.Time `json:"createTime" description:"发放时间"`
}

func (GoodsGrant) TableName() string {
	return "t_goods_grant"
}
