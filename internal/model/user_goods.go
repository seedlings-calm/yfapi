package model

import "time"

type UserGoods struct {
	Id           uint64    `json:"id" `           // 自增主键
	UserId       string    `json:"userId" `       // 用户ID
	GoodsId      string    `json:"goodsId" `      // 商品ID
	GoodsTypeId  string    `json:"goodsTypeId" `  // 商品类型ID
	GoodsTypeKey string    `json:"goodsTypeKey" ` // 商品类型key
	ExpireTime   time.Time `json:"expireTime" `   // 商品有效期
	OrderId      string    `json:"orderId" `      // 订单ID，来源为1时，需要绑定订单
	IsUse        int       `json:"isUse" `        // 是否使用：1：未使用，2：使用中
	Nums         int       `json:"nums" `         // 商品来源次数，小于等于100时 展示红点
	CreateTime   time.Time `json:"createTime" `   // 创建时间
	UpdateTime   time.Time `json:"updateTime" `   // 更新时间
}

func (UserGoods) TableName() string {
	return "t_user_goods"
}
