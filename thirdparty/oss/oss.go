package oss

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"yfapi/core/coreConfig"
)

func GenerateUploadPhotoToken() (token sts.Credentials, err error) {
	ossConfig := coreConfig.GetHotConf().Oss
	client, e := sts.NewClientWithAccessKey(ossConfig.DefaultRegionId, ossConfig.AccessKey, ossConfig.SecretKey)
	if e != nil {
		err = e
		return
	}
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = ossConfig.DefaultArn
	request.RoleSessionName = ossConfig.DefaultRoleName
	response, e := client.AssumeRole(request)
	if e != nil {
		err = e
		return
	}
	token = response.Credentials
	res, _ := json.Marshal(response.Credentials)
	fmt.Println(string(res))
	return
}
