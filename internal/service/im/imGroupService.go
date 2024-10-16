package service_im

import (
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/service/auth"
	"yfapi/internal/service/kafka"
	"yfapi/internal/service/riskCheck/shumei"
	"yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/message"
	"yfapi/typedef/response/im"
)

type ImGroupService struct {
}

// 发送群聊文本消息
func (im *ImGroupService) SendTextMsg(c *gin.Context, content, fromUserId, toUserId, groupId, clientType, extra string) response_im.ResponseMsg {
	isSuperAdmin := new(auth.Auth).IsSuperAdminRole(groupId, fromUserId)
	if new(dao.GroupFilterDao).MuteStatus(groupId) && !isSuperAdmin {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGroupIsMute,
		})
	}
	code := enum.GROUP_TEXT_MSG
	if isSuperAdmin {
		code = enum.GROUP_ADMIN_TEXT_MSG
	}
	msgData := message.MsgText{Content: content}
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId:    messageId,
		Timestamp:    timestamp,
		MsgType:      enum.MsgText,
		Code:         code,
		MsgData:      msgData,
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RiskLevel:    shumei.RiskLevelPass,
		RoomId:       groupId,
	}
	riskRes, ok := new(shumei.ShuMei).PrivateChatCheck(fromUserId, toUserId, content)
	if !ok {
		baseData.RiskLevel = shumei.RiskLevelReject
		coreLog.LogError("SendTextMsg reject %s", riskRes)
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.ROOM_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		ToRoomId:       groupId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           code,
		RiskLevel:      baseData.RiskLevel,
		RiskData:       riskRes,
	}
	service_kafka.New().PushPublicChat(kafkaData)
	return response_im.ResponseMsg{
		Sn:    0,
		Msg:   baseData,
		Extra: extra,
	}
}

// 发送群聊图片消息
func (im *ImGroupService) SendImgMsg(c *gin.Context, url, fromUserId, toUserId, groupId string, width, height int, clientType string, extra string) response_im.ResponseMsg {
	isSuperAdmin := new(auth.Auth).IsSuperAdminRole(groupId, fromUserId)
	if new(dao.GroupFilterDao).MuteStatus(groupId) && !isSuperAdmin {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGroupIsMute,
		})
	}
	code := enum.GROUP_IMG_MSG
	if isSuperAdmin {
		code = enum.GROUP_ADMIN_IMG_MSG
	}
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId: messageId,
		Timestamp: timestamp,
		MsgType:   enum.MsgImg,
		Code:      code,
		MsgData: message.MsgImg{
			Content: helper.FormatImgUrl(url),
			Width:   width,
			Height:  height,
		},
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RiskLevel:    shumei.RiskLevelPass,
		RoomId:       groupId,
	}
	riskRes, ok := new(shumei.ShuMei).OneChatImageSyncCheck(fromUserId, helper.FormatImgUrl(url))
	if !ok {
		baseData.RiskLevel = shumei.RiskLevelReject
		coreLog.LogError("SendImgMsg reject")
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.ROOM_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		ToRoomId:       groupId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           code,
		RiskLevel:      baseData.RiskLevel,
		RiskData:       riskRes,
	}
	service_kafka.New().PushPublicChat(kafkaData)
	return response_im.ResponseMsg{
		Sn:    0,
		Msg:   baseData,
		Extra: extra,
	}
}

// 发送动作消息
func (im *ImGroupService) SendActionMsg(c *gin.Context, data any, fromUserId, toUserId, groupId, clientType string, code int) {
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId:    messageId,
		Timestamp:    timestamp,
		MsgType:      enum.MsgAction,
		Code:         code,
		MsgData:      data,
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RoomId:       groupId,
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.ROOM_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		ToRoomId:       groupId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           code,
	}
	service_kafka.New().PushAction(kafkaData)
	return
}
