package pay

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pay/gopay/apple"
	"strings"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
)

type ApplePay struct {
	Url    string
	isProd bool //是否正式环境
}

// 初始化
func NewApplePay() *ApplePay {
	a := &ApplePay{
		Url: apple.UrlSandbox,
	}
	env := coreConfig.GetHotConf().ENV
	if env == "pro" {
		a.Url = apple.UrlProd
		a.isProd = true
	}
	return a
}

func (a *ApplePay) Client() (*apple.Client, error) {
	return apple.NewClient("", "", "", "", a.isProd)
}

func (a *ApplePay) VerifyReceipt(userId, receipt, pwd string) error {
	rsp, err := apple.VerifyReceipt(context.Background(), a.Url, pwd, receipt)
	if err != nil {
		return err
	}
	/**
	 * 21000  App Store无法读取您提供的JSON对象。
	 * 21002 该receipt-data属性中的数据格式错误或丢失。
	 * 21003 收据无法认证。
	 * 21004 您提供的共享密码与您帐户的文件共享密码不匹配。
	 * 21005 收据服务器当前不可用。
	 * 21006 该收据有效，但订阅已过期。当此状态代码返回到您的服务器时，收据数据也会被解码并作为响应的一部分返回。仅针对自动续订的iOS 6样式交易收据返回。
	 * 21007 该收据来自测试环境，但已发送到生产环境以进行验证。而是将其发送到测试环境。
	 * 21008 该收据来自生产环境，但是已发送到测试环境以进行验证。而是将其发送到生产环境。
	 * 21010 此收据无法授权。就像从未进行过购买一样对待。
	 * 21100-21199 内部数据访问错误。
	 */
	if rsp.Status != 0 {
		return fmt.Errorf("苹果支付 验证票据失败,status:%d", rsp.Status)
	}
	if len(rsp.Environment) == 0 {
		return errors.New("environment is empty")
	}
	coreLog.Info("苹果支付 苹果凭证验证结果:%+v", rsp)
	iosBundleIdList := strings.Split(coreConfig.GetHotConf().Ios.BundleId, ",")
	isBundle := false
	for _, item := range iosBundleIdList {
		if item == rsp.Receipt.BundleId {
			isBundle = true
			break
		}
	}
	if !isBundle {
		return errors.New("苹果支付 伪造票据")
	}
	transactions := rsp.Receipt.InApp
	if len(transactions) == 0 {
		return fmt.Errorf("苹果支付 缺少票据信息 %+v", rsp.Receipt.InApp)
	}
	for _, v := range transactions {
		err = appleIapRecharge(userId, v.TransactionId, v.ProductId)
		if err != nil {
			return err
		}
	}
	return nil
}
