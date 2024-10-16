package im

import (
	"encoding/json"
	"yfapi/app/handle"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/helper"
	service_im "yfapi/internal/service/im"
	"yfapi/internal/service/user"
	"yfapi/typedef/message"
	common_data "yfapi/typedef/redisKey"
	"yfapi/typedef/request/im"
	"yfapi/typedef/response"
	response_im "yfapi/typedef/response/im"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// GetSessionList
//
// @Summary 获取会话列表
// @Schemes
// @Description 获取会话列表
// @Tags 会话相关
// @Produce json
// @Success 200 {object} []response_im.GetSessionListRes
// @Router /v1/sessionList [get]
func GetSessionList(c *gin.Context) {
	userId := helper.GetUserId(c)
	redisClient := coreRedis.GetImRedis()
	var result []response_im.GetSessionListRes
	//根据排序获取会话列表
	sessionIdList := redisClient.ZRevRange(c, common_data.ImOneSessionSortId(userId), 0, 100).Val()
	if len(sessionIdList) > 0 {
		redisVal := redisClient.HMGet(c, common_data.ImOneSessionList(userId), sessionIdList...).Val()
		for _, item := range redisVal {
			var redisModel message.OneMsgListModel
			if item == nil {
				continue
			}
			err := json.Unmarshal([]byte(item.(string)), &redisModel)
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeSystemBusy,
					Msg:  nil,
				})
			}
			notReadNum := redisClient.HGet(c, common_data.ImOneMsgNotReadNum(userId), redisModel.ToUserId).Val()
			result = append(result, response_im.GetSessionListRes{
				Timestamp:   redisModel.Timestamp,
				TextColor:   redisModel.TextColor,
				ShowContent: redisModel.ShowContent,
				UserInfo:    user.GetSessionListUserInfo(userId, redisModel.ToUserId, helper.GetClientType(c)),
				NotReadNum:  cast.ToInt(notReadNum),
				IsTop:       redisModel.IsTop,
				Types:       service_im.GetSessionListTypes(redisModel.ToUserId),
			})
		}
	}
	response.SuccessResponse(c, result)
}

// MessageRead
//
//	@Description: 消息已读
func MessageRead(c *gin.Context) {
	req := im.MessageReadReq{}
	handle.BindBody(c, &req)
	tokenData := handle.GetTokenData(c)
	service_im.ClearNotReadNum(tokenData.UserId, req.ChatUserId)
	response.SuccessResponse(c, "")
}

// 所有消息已读
func MessageReadAll(c *gin.Context) {
	userId := handle.GetUserId(c)
	redisClient := coreRedis.GetImRedis()
	sessionIdList := redisClient.ZRevRange(c, common_data.ImOneSessionSortId(userId), 0, 100).Val()
	if len(sessionIdList) > 0 {
		redisVal := redisClient.HMGet(c, common_data.ImOneSessionList(userId), sessionIdList...).Val()
		for _, item := range redisVal {
			var redisModel message.OneMsgListModel
			if item == nil {
				continue
			}
			err := json.Unmarshal([]byte(item.(string)), &redisModel)
			if err != nil {
				continue
			}
			service_im.ClearNotReadNum(userId, redisModel.ToUserId)
		}
	}
	response.SuccessResponse(c, "")
}

// 置顶会话
func TopSession(c *gin.Context) {
	req := im.MessageReadReq{}
	handle.BindBody(c, &req)
	tokenData := handle.GetTokenData(c)
	service_im.TopSession(tokenData.UserId, req.ChatUserId)
	response.SuccessResponse(c, "")
}

// 删除会话
func DelSession(c *gin.Context) {
	req := im.MessageReadReq{}
	handle.BindBody(c, &req)
	tokenData := handle.GetTokenData(c)
	service_im.DelSession(tokenData.UserId, req.ChatUserId)
	response.SuccessResponse(c, "")
}

// 取消置顶
func UnTopSession(c *gin.Context) {
	req := im.MessageReadReq{}
	handle.BindBody(c, &req)
	tokenData := handle.GetTokenData(c)
	service_im.UnTopSession(tokenData.UserId, req.ChatUserId)
	response.SuccessResponse(c, "")
}

// 清空会话消息
func ClearChatHistory(c *gin.Context) {
	req := im.MessageReadReq{}
	handle.BindBody(c, &req)
	tokenData := handle.GetTokenData(c)
	service_im.ClearChatHistory(tokenData.UserId, req.ChatUserId)
	service_im.ClearNotReadNum(tokenData.UserId, req.ChatUserId)
	response.SuccessResponse(c, "")
}
