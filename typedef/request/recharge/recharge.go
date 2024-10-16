package recharge

// UserRechargeTestReq 用户测试充值请求
type UserRechargeTestReq struct {
	Diamond int64 `json:"diamond" validate:"required"`
}

type IosIapReq struct {
	Receipt string `json:"receipt" validate:"required"`
}

type WxAppPayReq struct {
	ProductId string `json:"productId" validate:"required"`
}

type AggregationPayReq struct {
	ProductId string `json:"productId" validate:"required"`
	Payment   string `json:"payment" validate:"required"` //支付方式 支付宝h5 ALIPAY_H5,支付宝扫码 ALIPAY_QR,微信h5 WECHAT_H5,微信扫码 WECHAT_QR
}

type AnotherPayReq struct {
	BankCardNum  string `json:"bankCardNum" validate:"required"`  //银行卡号
	BankCardName string `json:"bankCardName" validate:"required"` //持卡人姓名 必须
	Branch       string `json:"branch" validate:"required"`       //支行名称
	BankCode     string `json:"bankCode" validate:"required"`     //银行卡标识 例如IBCB
	OrderId      string `json:"orderId" validate:"required"`      //订单号
	UserId       string `json:"userId" validate:"required"`       //用户ID
	Amount       string `json:"amount" validate:"required"`       //金额 保留两位小数 1.00
}

type WebsitePayReq struct {
	ProductId string `json:"productId" validate:"required"`
	Payment   string `json:"payment" validate:"required"` //支付方式 支付宝h5 ALIPAY_H5,支付宝扫码 ALIPAY_QR,微信h5 WECHAT_H5,微信扫码 WECHAT_QR
	UserId    string `json:"userId" validate:"required"`
}

// 支付结果查询
type RechargeResultReq struct {
	UserId  string `json:"userId" validate:"required"`
	OrderId string `json:"orderId" validate:"required"`
}
