package response_goods

import (
	"time"
)

type GoodsTypesListRes struct {
	Id         int    `json:"id"`
	Icon       string `json:"icon"`   //图标
	Name       string `json:"name"`   //名称
	Keys       string `json:"keys"`   //键
	Sort       int    `json:"sort"`   //排序
	Status     int    `json:"status"` // 是否上架：2：下架中，1：上架中
	CreateTime string `json:"createTime"`
}

type GoodsListByTypesRes struct {
	GoodsId          int       `json:"goodsId"    description:"物品ID"`                  //商品ID
	GoodsTypeId      int       `json:"goodsTypeId" `                                   //商品类型ID
	GoodsName        string    `json:"goodsName"`                                      //商品名称
	Icon             string    `json:"icon"             description:"图标"`              //图标
	AnimationUrl     string    `json:"animationUrl"     description:"图片动效"`            //图片动效
	AnimationJsonUrl string    `json:"animationJsonUrl" description:"json文件动效"`        //json文件动效
	Money            int       `json:"money"      description:"单价/天"`                  //单价
	Moneys           string    `json:"moneys"     description:"7/15/30天价格\"7,15,28\""` //批量价格
	CreateTime       time.Time `json:"createTime" description:""`                      //时间
}

type GoodsAllRes struct {
	Id        int                    `json:"id"`
	Icon      string                 `json:"icon"` //图标
	Name      string                 `json:"name"` //名称
	Keys      string                 `json:"keys"` //键
	Sort      int                    `json:"sort"` //排序
	GoodsList []*GoodsListByTypesRes `json:"goodsList"`
}

type GoodsListToUserRes struct {
	Id        int          `json:"id"`
	Icon      string       `json:"icon"`      //图标
	Name      string       `json:"name"`      //名称
	Keys      string       `json:"keys"`      //键
	Sort      int          `json:"sort"`      //排序
	Nums      int          `json:"nums"`      //小于等于100时 展示红点
	Status    int          `json:"status"`    // 是否上架：2：下架中，1：上架中
	GoodsList []*UserGoods `json:"goodsList"` //用户商品
}

type UserGoods struct {
	GoodsId          int       `json:"goodsId"    description:"物品ID"`           //商品ID
	GoodsTypeId      int       `json:"goodsTypeId" `                            //商品类型ID
	Name             string    `json:"goodsName"`                               //商品名称
	Icon             string    `json:"icon"             description:"图标"`       //图标
	AnimationUrl     string    `json:"animationUrl"     description:"图片动效"`     //图片动效
	AnimationJsonUrl string    `json:"animationJsonUrl" description:"json文件动效"` //json文件动效
	ExpireTime       time.Time `json:"expireTime"`                              //有效期
	IsUse            int       `json:"isUse"`                                   //是否使用 1：未使用，2：使用中'
	Nums             int       `json:"nums"`                                    //商品来源次数，小于等于100时 展示红点
	CreateTime       time.Time `json:"createTime" description:""`               //时间
}

// 特效
type SpecialEffects struct {
	GoodsId          string `json:"goodsId" `                                // 商品ID
	GoodsName        string `json:"goodsName"`                               //商品名称
	GoodsIcon        string `json:"goodsIcon"`                               //商品图标
	GoodsTypeKey     string `json:"goodsTypeKey"`                            //物品分类key
	AnimationUrl     string `json:"animationUrl"     description:"图片动效"`     //图片动效
	AnimationJsonUrl string `json:"animationJsonUrl" description:"json文件动效"` //json文件动效
}
