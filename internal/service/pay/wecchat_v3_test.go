package pay

import (
	"fmt"
	"testing"
	"yfapi/internal/service/accountBook"
)

func TestWechatV3_V3TransactionApp(t *testing.T) {
	v3, err := NewWechatV3()
	if err != nil {
		panic(err)
	}
	rsp, err := v3.V3TransactionApp(V3TransactionAppReq{
		Mchid:       "1661666483",
		OutTradeNo:  "CZ202409110000635",
		Appid:       "wxc6978cf5aefcb8d3",
		Description: "测试商品",
		NotifyUrl:   "https://api.sdwsweb.com/api",
		Amount: struct {
			Total    int64  `json:"total"`
			Currency string `json:"currency"`
		}{
			Total:    1,
			Currency: accountBook.CURRENCY_CNY,
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\r\n", rsp)
}
