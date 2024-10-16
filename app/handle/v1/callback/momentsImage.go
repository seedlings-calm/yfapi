package v1_callback

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"strings"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/logic"
	"yfapi/internal/model"
	service_im "yfapi/internal/service/im"
	"yfapi/internal/service/riskCheck/shumei"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/response"
	response_user "yfapi/typedef/response/user"
	"yfapi/util/easy"
)

// 动态图片审核回调
func MomentsImageCallback(c *gin.Context) {
	data := shumei.ImagesCheckCallBackResp{}
	sign := c.Param("param")
	if len(sign) == 0 {
		coreLog.Error("MomentsImageCallback err %s", sign)
		return
	}
	decodeString, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		coreLog.LogError("MomentsImageCallback DecodeString err:%+v", err)
		return
	}
	id := easy.AesDecrypt(string(decodeString), coreConfig.GetHotConf().InnerSecret.AESEncryptKey1)
	err = c.ShouldBind(&data)
	if err != nil {
		coreLog.LogError("MomentsImageCallback bind err:%+v", err)
		return
	}
	coreLog.Info("MomentsImageCallback data:%+v,sign:%s,id:%s", data, sign, id)
	reject := false
	var examineResult []string
	for _, v := range data.Imgs {
		if v.RiskLevel == shumei.RiskLevelReject {
			reject = true
			examineResult = append(examineResult, v.RiskDescription)
		}
	}
	if reject { //违规处理
		new(dao.TimelineDao).Update(&model.Timeline{
			Id:            cast.ToInt64(id),
			Status:        typedef_enum.TimelineStatusRefuseExamine,
			ShumeiExamine: strings.Join(examineResult, ";"),
			UpdateTime:    time.Now(),
		})
		//违规通知
		timelineModel, _ := new(dao.TimelineDao).GetTimelineById(cast.ToInt64(id))
		var tempList []response_user.TimelineImgDTO
		_ = json.Unmarshal([]byte(timelineModel.ImgList), &tempList)
		if len(tempList) > 0 {
			go new(service_im.ImNoticeService).SendSystematicMsg(c, fmt.Sprintf("您发布的动态违规"), helper.FormatImgUrl(tempList[0].ImgPhotoKey), strings.Join(examineResult, ";"), "", "", []string{timelineModel.UserId})
		}
	} else { //通过处理
		new(dao.TimelineDao).Update(&model.Timeline{
			Id:         cast.ToInt64(id),
			Status:     typedef_enum.TimelineStatusNormal,
			UpdateTime: time.Now(),
		})
		timelineModel, _ := new(dao.TimelineDao).GetTimelineById(cast.ToInt64(id))
		go new(logic.Notice).MomentsPublishNotice(c, timelineModel.UserId, cast.ToInt64(id))
	}
	response.SuccessResponse(c, "")
}
