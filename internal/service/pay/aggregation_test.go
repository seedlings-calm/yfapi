package pay

import "testing"

// bankCode=1&merchantCode=123&returnUrl=http://127.0.0.1&key=123
func TestAggregationPay_generateSignature(t *testing.T) {
	ser := new(AggregationPay)
	params := map[string]string{
		"merchantCode": "123",
		"bankCode":     "1",
		"returnUrl":    "http://127.0.0.1",
	}
	signature := ser.generateSignature(params, "123")
	t.Log(signature)
}
