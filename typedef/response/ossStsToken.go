package response

import "github.com/aliyun/alibaba-cloud-sdk-go/services/sts"

// UploadPhotoStsTokenRes
// @Description: oss上传图片token
type UploadPhotoStsTokenRes struct {
	Credentials sts.Credentials //对应信息
	Bucket      string          //对应桶
	Region      string          //对应区域
	EndPoint    string          //对应节点
	ImgPrefix   string          //图片域名前缀
}
