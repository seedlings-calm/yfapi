package orderBill

type DiamondBillReq struct {
	Page     int    `json:"page" form:"page"`                     // 页码 从0开始
	Size     int    `json:"size" form:"size" validate:"required"` // 数量
	TimeKey  string `json:"timeKey" form:"timeKey"`               // 时间年月 2024-08
	FundFlow int    `json:"fundFlow" form:"fundFlow"`             // 资金方向 1入账 2出账
}

// 充值钻石日志
type RechargeDiamondReq struct {
	Page    int    `json:"page" form:"page"`                     // 页码 从0开始
	Size    int    `json:"size" form:"size" validate:"required"` // 数量
	TimeKey string `json:"timeKey" form:"timeKey"`               // 时间年月 2024-08
}
