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
	AnotherPayUrl = "/api/withdraw"
)

type AnotherPay struct {
}

func NewAnotherPay() *AnotherPay {
	return &AnotherPay{}
}

type AnotherPayReq struct {
	MerchantCode string `json:"merchantCode"` //商户号 必须
	OrderId      string `json:"orderId"`      //商户订单号 必须
	BankCardNum  string `json:"bankCardNum"`  //银行卡号 必须
	BankCardName string `json:"bankCardName"` //持卡人姓名 必须
	Branch       string `json:"branch"`       //支行名称
	BankCode     string `json:"bankCode"`     //通道代码 必须
	Amount       string `json:"amount"`       //金额 必须
	NotifyUrl    string `json:"notifyUrl"`    //回调地址 必须
	OrderDate    string `json:"orderDate"`    //13位时间戳 必须
	Currency     string `json:"currency"`
}

type AnotherPayResp struct {
	ResultCode string      `json:"resultCode"`
	ResultMsg  interface{} `json:"resultMsg"`
	Success    bool        `json:"success"`
	Data       struct {
		Data struct {
			OrderId string      `json:"orderId"`
			Status  int         `json:"status"`
			Money   float64     `json:"money"`
			Message interface{} `json:"message"`
		} `json:"data"`
		OrderId      string `json:"orderId"`
		MerchantCode string `json:"merchantCode"`
		Date         string `json:"date"`
		Sign         string `json:"sign"`
	} `json:"data"`
}

// 回调通知结构
type AnotherPayNotifyResp struct {
	Status       string `json:"Status"`
	OutTradeNo   string `json:"OutTradeNo"`
	MerchantCode string `json:"MerchantCode"`
	Amount       string `json:"Amount"`
	Time         string `json:"Time"`
	Sign         string `json:"Sign"`
	OrderId      string `json:"OrderId"`
	Fee          string `json:"Fee"`
}

func (a *AnotherPay) Pay(req *AnotherPayReq) (res AnotherPayResp, err error) {
	if req == nil {
		err = errors.New("req is nil")
		return
	}
	conf := coreConfig.GetHotConf().AnotherPay
	params := map[string]string{
		"merchantCode": req.MerchantCode,
		"orderId":      req.OrderId,
		"bankCardNum":  req.BankCardNum,
		"bankCardName": req.BankCardName,
		"branch":       req.Branch,
		"bankCode":     req.BankCode,
		"amount":       req.Amount,
		"notifyUrl":    conf.NotifyUrl,
		"orderDate":    req.OrderDate,
		"currency":     req.Currency,
	}
	signature := a.generateSignature(params, conf.Secret)
	formData := url.Values{}
	formData.Set("sign", signature)
	for k, v := range params {
		formData.Set(k, v)
	}
	bodyStr, err := easy.PostForm(conf.Url+AnotherPayUrl, formData)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(bodyStr), &res)
	return
}

// 生成签名
func (a *AnotherPay) generateSignature(params map[string]string, secret string) string {
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
func (a *AnotherPay) VerifySignature(data *AnotherPayNotifyResp) bool {
	params := map[string]string{}
	if len(data.MerchantCode) > 0 {
		params["MerchantCode"] = data.MerchantCode
	}
	if len(data.OrderId) > 0 {
		params["OrderId"] = data.OrderId
	}
	if len(data.Amount) > 0 {
		params["Amount"] = data.Amount
	}
	if len(data.Fee) > 0 {
		params["Fee"] = data.Fee
	}
	if len(data.OutTradeNo) > 0 {
		params["OutTradeNo"] = data.OutTradeNo
	}
	if len(data.Time) > 0 {
		params["Time"] = data.Time
	}
	if len(data.Status) > 0 {
		params["Status"] = data.Status
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
	str := paramStr.String() + "&Key=" + coreConfig.GetHotConf().AnotherPay.Secret
	coreLog.Info("代付签名：%s", str)
	md5Str := easy.Md5(str, 32, false)
	return md5Str == data.Sign
}
