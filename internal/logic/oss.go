package logic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/thirdparty/oss"
	common_data "yfapi/typedef/redisKey"
	"yfapi/typedef/response"
	"yfapi/util/easy"
)

type OssStsToken struct {
}

// GetOssStsToken
//
//	@Description: oss sts token
//	@receiver o
//	@param ctx *gin.Context -
//	@return res -
func (o *OssStsToken) GetOssStsToken(ctx *gin.Context) (res response.UploadPhotoStsTokenRes) {
	ossConfig := coreConfig.GetHotConf().Oss
	redisClient := coreRedis.NewRedis().GetUserRedis()
	tokenKey := common_data.OssUploadPhotoStsToken()
	cacheToken := redisClient.Get(ctx, tokenKey).Val()
	if len(cacheToken) > 0 {
		err := json.Unmarshal([]byte(cacheToken), &res)
		if err != nil {
			coreLog.Error(err)
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
		return
	}

	token, err := oss.GenerateUploadPhotoToken()
	if err != nil {
		coreLog.Error(err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeCreateToken,
			Msg:  nil,
		})
	}
	res.Credentials = token
	res.Bucket = ossConfig.DefaultBucket
	res.Region = "oss-" + ossConfig.DefaultRegionId
	res.EndPoint = ossConfig.DefaultEndPoint
	res.ImgPrefix = coreConfig.GetHotConf().ImagePrefix

	expireTime, _ := time.ParseInLocation(time.RFC3339, res.Credentials.Expiration, time.Local)
	_ = redisClient.Set(ctx, tokenKey, easy.JSONStringFormObject(res), expireTime.Sub(time.Now().Add(10*time.Minute)))
	return
}
