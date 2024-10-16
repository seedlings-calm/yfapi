package pay

import (
	"context"
	"fmt"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
)

const (
	wechatPrivateKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDCEzAuRidPnWDV
SwXT1FCzkFRIDFdLSJC8iX4PoBw5uTret9Onq2Sgeetw/HlyD1UPTciprOBZBVNP
EjU9uSIPAGkxYMaDWNRGhVz3lZspvW24sRQAs52IF9xwcvgLMxGjCL1UWrJnCoOe
+U5+FkiPp1K5j4o/TLHaQwWAhaa85uEM0a6jSW+rfGZGOk0h5A1cRY9FYJ+t/B4g
SnpEhntzvJf4kQht11xSgjJHcP58eiqk+NKja1jt7nYtbzal5yAM3ghVnjroV6ZD
n/l6CvjwDS9ncQjYU9DmDxxGsGduupqgT1vWZl6s0zCINE6fxXHzrbuzKEexH/xd
lwqTXYGVAgMBAAECggEBAJVbNwUlsDMxJsh/SCjRiJnoRTR9auDhI6I1HfLggVhb
GNc2GYk7+eEWcv8tDjmdWxTAA+Gwzac01fxQQvERfQiWhF4f4CZNQnBTSkyvsg3Z
Xot4m7A+ismls1xG3mWgE76ohyxX8FwkV08NBj863vTPDHcJ5Jb6axVR5vYV5VAR
9eRowVXB547+SZQVnkAj1jPbTBN/Gu5r0Aw655q38JGorpOCsTvK7jZdb/yNnicR
MaZKUU3XOtp9hb3fHv0oVPMal4P1fw4/rNyQd9qxEltT9zo8BnU8mxm9bNeKX7m1
sHtGXaya3AB+SpHEGiFouD6fxp4akS+wwxfueEnqxYECgYEA9pAqGOxFM6/zcPRw
rOJa18I2rVqtUz8BmoF7aXJw4Xb6QG9B3X/pF7SUFSWdttci8UTDdUJnxEkOcCV1
W4zTjOgv0bVX5yBbIW0SBbktOvVPURiinJCSsKAWtBEcM5RDgqGgc5pXhqb+cFVj
XeE1OxsgMZ8SgdsemSnjI2dOTgkCgYEAyYC+DST7MVwLvz2cTEQnik6qZ0sFkvqB
won7u6lrfJSecfMUyvw+VS10GfRrdQT9xAHMxBdAKUDCYUvIk4Ignm9Bb+h3hr4v
v25kwV2e4PrPpAwn0Dmwb3FlTI1ARhr15Rcndhi6ZpfrZLFZ2AH5mTeSFLfC8LmC
rtOz0hh9+i0CgYBqbcZOJcalTgD2M/1jEv4Ffibd43NHqL3HdLbRyH1jRVk2cQ/s
TadO/Tqiei8+2lSR8o2wUu65spNR40lqMDqs6xihG/cKpFPR9OO/prQYbAVFyy1+
CYXYSfIi3fPfa9NMUvoQjIHVdMYFtUYEIw84KThXUwPJG810bblG8fPCCQKBgDl0
szxOQi3V4Cecqrd4a4ndWmtvkdxR/7P34kalTVfNjMxTEqe6ew+QkV1hO063qKA+
HyP+uTXKGGLj2AJvhVuHv7HoKETMcBL2qFYWmtntyk0thiCygmOUgtzsHdqfj2PO
UVs0O9pLETy58TNNhN0yYj30E+rOCrxM8yZCA5HVAoGBAJaenmLwKQR+z03k148Q
+uxEVFb5aw9i2DH23r65o0C84AmjXYtoZQWsSb5APyZgdZRwfqE44tUPrbj0dkwG
EvdVWeu2x31+AWiBvjvtQ/c6BI5QZYYbRuLHA94MkCKXrOSNo7cyTqFErM/wJ4q4
sYXPcE9w+s9eNWIlI7eBb8L8
-----END PRIVATE KEY-----`
)

type WechatV3 struct {
	client *wechat.ClientV3
	mchid  string
	appId  string
}

func NewWechatV3() (wechatV3 *WechatV3, err error) {
	wechatV3 = &WechatV3{}
	wxpay := coreConfig.GetHotConf().WxPay
	mchid := wxpay.Mchid       //商户ID
	serialNo := wxpay.SerialNo //商户证书序列号
	apiV3Key := wxpay.ApiV3Key
	client := &wechat.ClientV3{}
	client, err = wechat.NewClientV3(mchid, serialNo, apiV3Key, wechatPrivateKey)
	if err != nil {
		return
	}
	err = client.AutoVerifySign()
	if err != nil {
		return
	}
	//client.DebugSwitch = gopay.DebugOn
	wechatV3.client = client
	wechatV3.appId = wxpay.Appid
	wechatV3.mchid = mchid
	return
}

type V3TransactionH5Req struct {
	Description string `json:"description"`  //商品描述
	OutTradeNo  string `json:"out_trade_no"` //商户订单号
	NotifyUrl   string `json:"notify_url"`   //回调通知地址
	Amount      struct {
		Total    int    `json:"total"` //金额单位：分
		Currency string `json:"currency"`
	} `json:"amount"`
	SceneInfo struct {
		PayerClientIp string `json:"payer_client_ip"` //客户端IP
		H5Info        struct {
			Type string `json:"type"` //iOS, Android, Wap
		} `json:"h5_info"`
	} `json:"scene_info"`
}

// h5v3支付
func (w *WechatV3) V3TransactionH5(data V3TransactionH5Req) (rsp *wechat.H5Rsp, err error) {
	bm := gopay.BodyMap{}.
		Set("appid", w.appId).
		Set("mchid", w.mchid).
		Set("description", data).
		Set("out_trade_no", data.OutTradeNo).
		Set("notify_url", data.NotifyUrl).
		Set("amount", gopay.BodyMap{
			"total": 0,
		}).
		Set("scene_info", gopay.BodyMap{
			"payer_client_ip": data.SceneInfo.PayerClientIp,
			"h5_info": gopay.BodyMap{
				"type": data.SceneInfo.H5Info.Type,
			},
		})
	rsp, err = w.client.V3TransactionH5(context.Background(), bm)
	return
}

type V3TransactionAppReq struct {
	Mchid       string `json:"mchid"`
	OutTradeNo  string `json:"out_trade_no"`
	Appid       string `json:"appid"`
	Description string `json:"description"`
	NotifyUrl   string `json:"notify_url"`
	Amount      struct {
		Total    int64  `json:"total"`
		Currency string `json:"currency"`
	} `json:"amount"`
}

// APPv3支付
func (w *WechatV3) V3TransactionApp(data V3TransactionAppReq) (app *wechat.AppPayParams, err error) {
	bm := gopay.BodyMap{}.
		Set("appid", w.appId).
		Set("mchid", w.mchid).
		Set("description", data.Description).
		Set("out_trade_no", data.OutTradeNo).
		Set("notify_url", data.NotifyUrl).
		Set("amount", gopay.BodyMap{
			"total": data.Amount.Total,
		})
	rsp, err := w.client.V3TransactionApp(context.Background(), bm)
	if err != nil {
		coreLog.Error("wechatPay V3TransactionApp err;%+v", err)
		return
	}
	fmt.Printf("%+v\r\n", rsp)
	app, err = w.client.PaySignOfApp(w.appId, rsp.Response.PrepayId)
	fmt.Println(err)
	if err != nil {
		coreLog.Error("wechatPay V3TransactionApp err;%+v", err)
		return
	}
	return
}
