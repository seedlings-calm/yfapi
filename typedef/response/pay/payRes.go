package pay

type AggregationPayResp struct {
	Types   string `json:"types"` //类型
	Info    string `json:"info"`  //值
	OrderId string `json:"orderId"`
}
