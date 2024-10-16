package v1_callback

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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

// 动态视频审核回调
func MomentsVideoCallback(c *gin.Context) {
	data := shumei.VideoCheckCallbackResult{}
	sign := c.Param("param")
	if len(sign) == 0 {
		coreLog.Error("MomentsVideoCallback err %s", sign)
		return
	}
	decodeString, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		coreLog.LogError("MomentsVideoCallback DecodeString err:%+v", err)
		return
	}
	id := easy.AesDecrypt(string(decodeString), coreConfig.GetHotConf().InnerSecret.AESEncryptKey1)
	err = c.ShouldBind(&data)
	if err != nil {
		coreLog.LogError("MomentsVideoCallback bind err:%+v", err)
		return
	}
	coreLog.Info("MomentsVideoCallback data:%+v,sign:%s,id:%s", data, sign, id)
	if data.RiskLevel == shumei.RiskLevelReject {

		new(dao.TimelineDao).Update(&model.Timeline{
			Id:         cast.ToInt64(id),
			Status:     typedef_enum.TimelineStatusRefuseExamine,
			UpdateTime: time.Now(),
		})
		timelineModel, _ := new(dao.TimelineDao).GetTimelineById(cast.ToInt64(id))
		var tempList response_user.TimelineVideoDTO
		_ = json.Unmarshal([]byte(timelineModel.VideoData), &tempList)
		if len(tempList.VideoUrl) > 0 {
			go new(service_im.ImNoticeService).SendSystematicMsg(c, fmt.Sprintf("您发布的动态违规"), helper.FormatImgUrl(tempList.VideoUrl+typedef_enum.VideoCoverImgSuffix), "", "", "", []string{timelineModel.UserId})
		}
	} else {
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
