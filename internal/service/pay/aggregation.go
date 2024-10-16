package pay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/util/easy"
)

const (
	AggregationPayUrl = "/api/pay"
)

// 聚合支付
type AggregationPay struct {
}

func NewAggregationPay() *AggregationPay {
	return &AggregationPay{}
}

type AggregationPayReq struct {
	MerchantCode string `json:"merchantCode"` //商户号 必须
	BankCode     string `json:"bankCode"`     //通道代码 必须
	Currency     string `json:"currency"`     //币种 必须
	Amount       string `json:"amount"`       //金额 必须
	OrderId      string `json:"orderId"`      //商户订单号 必须
	NotifyUrl    string `json:"notifyUrl"`    //回调地址 必须
	ReturnUrl    string `json:"returnUrl"`    //同步跳转地址 非必须
	OrderDate    string `json:"orderDate"`    //13位时间戳 必须
	Ip           string `json:"ip"`           //客户ip 必须
	GoodsName    string `json:"goodsName"`    //商品名称 必须
	GoodsDetail  string `json:"goodsDetail"`  //商品详情 必须
	Ext          string `json:"ext"`          //扩展字段 非必须
	UserId       string `json:"userId"`
}

type AggregationPayResp struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
	Success    bool   `json:"success"`
	Data       struct {
		Data struct {
			Type     string `json:"infoType"` //url:跳转地址，img:图片地址
			Info     string `json:"info"`
			BankCode string `json:"bankCode"`
		} `json:"data"`
		OrderId      string `json:"orderId"`
		OutTradeNo   string `json:"outTradeNo"`
		MerchantCode string `json:"merchantCode"`
		Date         string `json:"date"`
		Sign         string `json:"sign"`
	} `json:"data"`
}

// 回调通知结构
type AggregationPayNotifyResp struct {
	MerchantCode string `json:"merchantCode" form:"merchantCode"` //商户号
	OrderId      string `json:"orderId" form:"orderId"`           //商户订单号
	OrderDate    string `json:"orderDate" form:"orderDate"`       //提交订单时间戳
	Currency     string `json:"currency" form:"currency"`         //币种
	Amount       string `json:"amount" form:"amount"`             //支付金额
	OutTradeNo   string `json:"outTradeNo" form:"outTradeNo"`     //系统订单号
	BankCode     string `json:"bankCode" form:"bankCode"`         //银行代码
	Time         string `json:"time" form:"time"`                 //回调时间戳
	Remark       string `json:"remark" form:"remark"`             //备注
	Status       string `json:"status" form:"status"`             //状态 0处理中，1成功，2失败
	Sign         string `json:"sign" form:"sign"`                 //签名
	Fee          string `json:"fee" form:"fee"`
	FailReason   string `json:"failReason" form:"failReason"`
}

func (a *AggregationPay) Pay(req *AggregationPayReq) (res AggregationPayResp, err error) {
	if req == nil {
		err = errors.New("req is nil")
		return
	}
	conf := coreConfig.GetHotConf().AggregationPay
	params := map[string]string{
		"merchantCode": req.MerchantCode,
		"bankCode":     req.BankCode,
		"currency":     req.Currency,
		"amount":       req.Amount,
		"orderId":      req.OrderId,
		"notifyUrl":    conf.NotifyUrl,
		"orderDate":    req.OrderDate,
		"ip":           req.Ip,
		"goodsName":    req.GoodsName,
		"goodsDetail":  req.GoodsDetail,
		"userId":       req.UserId,
	}
	if len(req.ReturnUrl) > 0 {
		params["returnUrl"] = req.ReturnUrl
	}
	signature := a.generateSignature(params, conf.Secret)
	formData := url.Values{}
	formData.Set("sign", signature)
	for k, v := range params {
		formData.Set(k, v)
	}
	bodyStr, err := easy.PostForm(conf.Url+AggregationPayUrl, formData)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(bodyStr), &res)
	return
}

// 生成签名
func (a *AggregationPay) generateSignature(params map[string]string, secret string) string {
	// 将参数按照字典序排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 拼接成 key1=value1&key2=value2 的格式
	var paramStr strings.Builder
	for _, key := range keys {
		if paramStr.Len() > 0 {
			paramStr.WriteString("&")
		}
		paramStr.WriteString(fmt.Sprintf("%s=%s", key, params[key]))
	}

	str := paramStr.String() + "&Key=" + secret
	md5Str := easy.Md5(str, 32, false)
	return md5Str
}

// 验证签名
func (a *AggregationPay) VerifySignature(data *AggregationPayNotifyResp) bool {
	params := map[string]string{}
	if len(data.MerchantCode) > 0 {
		params["merchantCode"] = data.MerchantCode
	}
	if len(data.OrderId) > 0 {
		params["orderId"] = data.OrderId
	}
	if len(data.OrderDate) > 0 {
		params["orderDate"] = data.OrderDate
	}
	if len(data.Currency) > 0 {
		params["currency"] = data.Currency
	}
	if len(data.Amount) > 0 {
		params["amount"] = data.Amount
	}
	if len(data.OutTradeNo) > 0 {
		params["outTradeNo"] = data.OutTradeNo
	}
	if len(data.BankCode) > 0 {
		params["bankCode"] = data.BankCode
	}
	if len(data.Time) > 0 {
		params["time"] = data.Time
	}
	if len(data.Remark) > 0 {
		params["remark"] = data.Remark
	}
	if len(data.Status) > 0 {
		params["status"] = data.Status
	}
	if len(data.Fee) > 0 {
		params["fee"] = data.Fee
	}
	if len(data.FailReason) > 0 {
		params["failReason"] = data.FailReason
	}
	// 将参数按照字典序排序
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 拼接成 key1=value1&key2=value2 的格式
	var paramStr strings.Builder
	for _, key := range keys {
		if paramStr.Len() > 0 {
			paramStr.WriteString("&")
		}
		paramStr.WriteString(fmt.Sprintf("%s=%s", key, params[key]))
	}

	str := paramStr.String() + "&Key=" + coreConfig.GetHotConf().AggregationPay.Secret
	md5Str := easy.Md5(str, 32, false)
	coreLog.Info("聚合支付回调 参数:%s 校验签名:%s 原始签名:%s", str, md5Str, data.Sign)
	return md5Str == data.Sign
}
