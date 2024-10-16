package av

import "yfapi/core/coreConfig"

type AvIn interface {
	GetToken(userId, channelName string) (string, error) //获取token方法
}

func New() AvIn {
	avConf := coreConfig.GetHotConf().Av
	switch avConf.Supplier {
	case "shengwang":
		return new(ShengWangAv)
	default:
		return new(ShengWangAv)
	}
}
