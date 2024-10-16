package user

// ExchangeDiamondReq 兑换钻石请求
type ExchangeDiamondReq struct {
	ExchangeAmount int64 `json:"exchangeAmount" validate:"required,gt=0"` // 兑换钻石数量
}
