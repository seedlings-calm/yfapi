package request_goods

type BuyGoodsToUserReq struct {
	GoodsId int `json:"goodsId" form:"goodsId" validate:"required"`        // 商品ID
	Days    int `json:"days" form:"days" validate:"eq=1|eq=7|eq=15|eq=30"` //有效天数
	Num     int `json:"num" form:"num" validate:"min=1"`                   // 购买次数
}
