package room

import (
	"yfapi/app/handle"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	error2 "yfapi/i18n/error"
	"yfapi/internal/service/av"
	request_room "yfapi/typedef/request/room"
	"yfapi/typedef/response"
	"yfapi/typedef/response/room"

	"github.com/gin-gonic/gin"
)

// AvToken
//
//	@Description: 获取声网token
func AvToken(c *gin.Context) {
	req := new(request_room.AvTokenReq)
	handle.BindQuery(c, req)
	userId := handle.GetUserId(c)
	service := av.New()
	token, err := service.GetToken(userId, req.ChannelName)
	if err != nil {
		coreLog.LogError("获取声网token失败 %+v", err)
		panic(error2.I18nError{
			Code: error2.ErrCodeAvTokenFailed,
			Msg:  nil,
		})
	}
	resp := room.ShengWangResp{
		Token: token,
		AppId: coreConfig.GetHotConf().Av.AppId,
	}
	response.SuccessResponse(c, resp)
}
