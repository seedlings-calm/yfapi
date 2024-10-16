package pay

import (
	"fmt"
	"testing"
)

func TestAliV3_TradeAppPay(t *testing.T) {
	v3, err := NewAliV3()
	if err != nil {
		fmt.Println(err)
		return
	}
	rsp, err := v3.TradeAppPay(AliAppPayReq{
		OutTradeNo:  "CZ20120521201",
		TotalAmount: "0.1",
		Subject:     "aa",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}
