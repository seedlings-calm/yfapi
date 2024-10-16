package pay

import (
	"context"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
	"yfapi/core/coreLog"
)

const (
	aliPrivateKey = `MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDCF6IgHFJsiUDiCz81Xa6dGfn3f655qChApoKxhLIp3OuTkoNjoJPCyGTh07+nI+HFK9WJW7xZ++ZD3QoY6GoOxs+yBKcg4CEKnGB2zHO1AAaNh5lle9ZtKcv4JF0dhsnUaQTZ3WFtH0CtVpMvsBbnadED8LfJvM6KLXpv+/vQmkrSetwnXB/CTfqO212NusHAeK9bb/a+1Am553nMDxlkCwzh1B+YauLH4EVmhyWOFhpVZHiFT7O0Vl3dd98gc4a44kGcaopF/jV3rCVD51ihi75HSBjvPy+9LEJc4xlKa3EtGZp5SiOo8gdEtjAK5bl9tlrElxr7a+I6hujfZBqrAgMBAAECggEAYKjgZtlz+vWHyIsNWYhkM30CTc3amF+0XC4QnFOXXt3UvFOU94K606B1DTolEhn+j/E6kQOMk8uta1KjerAUUXOVb/R9PxQfoGcsaz16ykNPACDtteqsaQUNvXBupwu/a/c5IT7tDCkqTqj0+CTb4zeBjlLNVLygp5Pqi+aUC9hMbhDXzMUJNZ3P00xDxx61NT7lh6cZz0G+LU7Ro1v+wREcC7qngwcbyua3Fxm9QzPKo8auOhY9jQzxAN5Maa74lPtN9tRWLoGEqH6P/kgp0+/Abyl1IclepHZBjYFKVP9a7orqNq2QOHJEgxUPyDvAPF3I3I0mFLAjyY+B6kTFsQKBgQDyJz65P7l24VQTOL1JINNJoQt+nWtC82lsy38F1mwroACnvd5mwnJm9eaQwH6X1uGRkQexxhvTT3yCFZtFCVmIakLSCoEg4UdTPpf78yQ+UMusZVbT12yYb8XwhP3KNQGF0cy+eguNW5ngcUfiFmaa0iDYpHUMl4F5Mjk3a1GaHQKBgQDNMNlKyOnJKHR067HQUPYeWFbIbvsrL/WfPqf2ke5mMDiJx35pLbRY+pHrQd0o6vL6gL107pMNyF3x0qDu7UzFWN4RSeOpiCt+U5im5p79azPGp7ZYqTU+o7AozrJHlThZJwQQjW9yIJW3b5nEXitgrmRRtkqKPgd9nc7xxdgtZwKBgHwPbyspYrNtLc3LO+7DSnxmbaUosVNTsadzelhbSn/vMWa+97pd1I67XKy8ch8Ij/gr/W0uugLArmFXAH2WFLC3ABTsHMvjns8fOm5yWxcx/acNJDbUH2bZnOdku0FldqpAmkzb6h851tQONW23XEnlbb6QQwd4d3TILlgeO2hhAoGANZ3vGbiYRmUY0TiJdTrCpTlGLAe8ABP/JcZ0k1yco/0zuOT1Jjy4JIwNNyE1zixeo5CicPyqVm6mBbuZK/W8GtFW7cOWBsW7P75OZEZdAzFRDTbj0hUdAao6LN+d/FCEsd8dE8oxdewH0zAJSSOmSBQpyKROpAMDaBKlcc3V6D8CgYEA37kQZAjaykICij8XzZONvrKH8bOeToyZuM1KMB9biwxyKy2D3f5v9GzpQECPsx+RKAYNqR/SN+nrWOOcG2XJwf/S9q2DwOMIeI6w+tPVvgWLx3obCATJq5uiNYXMGOFPuwxZhEUh90++f7UeY7uhvVhmJxzMO/7D7g5nfqnDXOc=`
	aliPublicKey  = `MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwheiIBxSbIlA4gs/NV2unRn593+ueagoQKaCsYSyKdzrk5KDY6CTwshk4dO/pyPhxSvViVu8WfvmQ90KGOhqDsbPsgSnIOAhCpxgdsxztQAGjYeZZXvWbSnL+CRdHYbJ1GkE2d1hbR9ArVaTL7AW52nRA/C3ybzOii16b/v70JpK0nrcJ1wfwk36jttdjbrBwHivW2/2vtQJued5zA8ZZAsM4dQfmGrix+BFZocljhYaVWR4hU+ztFZd3XffIHOGuOJBnGqKRf41d6wlQ+dYoYu+R0gY7z8vvSxCXOMZSmtxLRmaeUojqPIHRLYwCuW5fbZaxJca+2viOobo32QaqwIDAQAB`
)

type AliV3 struct {
	client *alipay.Client
	appId  string
	isProd bool
}

func NewAliV3() (aliV3 *AliV3, err error) {
	aliV3 = &AliV3{}
	appId := "2021004143688499"
	client, err := alipay.NewClient(appId, aliPrivateKey, true)
	if err != nil {
		return
	}
	//client.DebugSwitch = gopay.DebugOn
	client.SetNotifyUrl("https://api.sdws.shop").
		SetReturnUrl("https://api.sdws.shop").
		SetSignType(alipay.RSA2)
	//SetCharset("GBK")
	aliV3.appId = appId
	aliV3.client = client
	return
}

type AliAppPayReq struct {
	OutTradeNo  string `json:"out_trade_no"` //商户订单号
	TotalAmount string `json:"total_amount"` //支付金额 单位：元
	Subject     string `json:"subject"`      //订单标题
}

// APPv3支付
func (a *AliV3) TradeAppPay(data AliAppPayReq) (rsp string, err error) {
	bm := gopay.BodyMap{}.
		Set("subject", data.Subject).
		Set("out_trade_no", data.OutTradeNo).
		Set("total_amount", data.TotalAmount)
	rsp, err = a.client.TradeAppPay(context.Background(), bm)
	if err != nil {
		coreLog.Error("aliPay TradeAppPay err;%+v", err)
		return
	}
	return
}
