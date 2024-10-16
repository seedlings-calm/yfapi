package agent

import (
	"encoding/base64"
	"fmt"
	"sort"
	"yfapi/util/easy"
)

type Oauth struct {
	AppId  string
	Secret string
}

// 计算签名
func (o *Oauth) Sign(data map[string]any) string {
	delete(data, "signature")
	// 按键名升序排序
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 将值转换为 Base64 编码并拼接
	var encodedValues []string
	for _, key := range keys {
		value := data[key]
		base64Value := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", value)))
		encodedValues = append(encodedValues, base64Value)
	}
	// 拼接字符串
	joinedString := ""
	for _, encodedValue := range encodedValues {
		joinedString += encodedValue
	}
	joinedString += o.Secret
	// 使用 SHA1 算法计算签名
	signature := easy.Sha1(joinedString)
	return signature
}
