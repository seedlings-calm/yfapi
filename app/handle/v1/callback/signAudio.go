package v1_callback

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/internal/service/riskCheck/shumei"
	"yfapi/typedef/response"
	"yfapi/util/easy"
)

// 声音签名审核回调
func SignAudioCallback(c *gin.Context) {
	data := shumei.AudioCheckCallbackResult{}
	sign := c.Param("param")
	coreLog.Info("数美声音签名检测回调 %s", sign)
	if len(sign) == 0 {
		coreLog.Error("SignAudioCallback err %s", sign)
		return
	}
	decodeString, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		coreLog.LogError("数美 SignAudioCallback DecodeString err:%+v", err)
		return
	}
	id := easy.AesDecrypt(string(decodeString), coreConfig.GetHotConf().InnerSecret.AESEncryptKey1)
	err = c.ShouldBind(&data)
	if err != nil {
		coreLog.LogError("数美 SignAudioCallback bind err:%+v", err)
		return
	}
	coreLog.Info("数美 SignAudioCallback data:%+v,sign:%s,id:%s", data, sign, id)
	reject := false
	if data.RiskLevel == shumei.RiskLevelReject {
		reject = true
	}
	if reject { //违规处理
		new(dao.UserDao).UpdateById(&model.User{
			Id:          id,
			VoiceStatus: 3,
		})
	} else { //通过处理
		new(dao.UserDao).UpdateById(&model.User{
			Id:          id,
			VoiceStatus: 1,
		})
	}
	response.SuccessResponse(c, "")
}
