package service_im

import (
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	service_kafka "yfapi/internal/service/kafka"
	"yfapi/internal/service/riskCheck/shumei"
	"yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/message"
	response_im "yfapi/typedef/response/im"

	"github.com/spf13/cast"
)

type ImOneService struct {
}

// 发送文本消息
func (im *ImOneService) SendTextMsg(c *gin.Context, content, fromUserId, toUserId, clientType, extra string) response_im.ResponseMsg {
	if helper.OfficialAccount(toUserId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 黑名单拦截
	if user.IsBlacklist(toUserId, fromUserId, "0", enum.BlacklistTypeUser) {
		panic(error2.I18nError{
			Code: error2.ErrorCodePrivateChatIsBlacklist,
			Msg:  nil,
		})
	}
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId:    messageId,
		Timestamp:    timestamp,
		MsgType:      enum.MsgText,
		Code:         enum.USER_TEXT_MSG,
		MsgData:      message.MsgText{Content: content},
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RiskLevel:    shumei.RiskLevelPass,
	}
	riskRes, ok := new(shumei.ShuMei).PrivateChatCheck(fromUserId, toUserId, content)
	if !ok {
		baseData.RiskLevel = shumei.RiskLevelReject
		coreLog.LogError("SendTextMsg reject %s", riskRes)
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.USER_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		RiskLevel:      baseData.RiskLevel,
		RiskData:       riskRes,
		Code:           enum.USER_TEXT_MSG,
	}
	service_kafka.New().PushPrivateChat(kafkaData)
	if baseData.RiskLevel == shumei.RiskLevelPass {
		//增加会话列表
		ChatListDealWith(content, fromUserId, toUserId, enum.MsgListTextColorNormal)
		ChatListDealWith(content, toUserId, fromUserId, enum.MsgListTextColorNormal)
		//处理消息未读数
		AddNotReadNum(toUserId, fromUserId)
	} else {
		ChatListDealWith(content, fromUserId, toUserId, enum.MsgListTextColorNormal)
	}
	return response_im.ResponseMsg{
		Sn:    0,
		Msg:   baseData,
		Extra: extra,
	}
}

// SendImgMsg 发送图片消息
func (im *ImOneService) SendImgMsg(c *gin.Context, url, fromUserId, toUserId string, width, height int, clientType, extra string) response_im.ResponseMsg {
	if helper.OfficialAccount(toUserId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 黑名单拦截
	if user.IsBlacklist(toUserId, fromUserId, "0", enum.BlacklistTypeUser) {
		panic(error2.I18nError{
			Code: error2.ErrorCodePrivateChatIsBlacklist,
			Msg:  nil,
		})
	}
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId: messageId,
		Timestamp: timestamp,
		MsgType:   enum.MsgImg,
		Code:      enum.USER_IMG_MSG,
		MsgData: message.MsgImg{
			Content: helper.FormatImgUrl(url),
			Width:   width,
			Height:  height,
		},
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RiskLevel:    shumei.RiskLevelPass,
	}
	riskRes, ok := new(shumei.ShuMei).OneChatImageSyncCheck(fromUserId, helper.FormatImgUrl(url))
	if !ok {
		baseData.RiskLevel = shumei.RiskLevelReject
		coreLog.LogError("SendImgMsg reject")
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.USER_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		RiskLevel:      baseData.RiskLevel,
		RiskData:       riskRes,
		Code:           enum.USER_IMG_MSG,
	}

	service_kafka.New().PushPrivateChat(kafkaData)
	if baseData.RiskLevel == shumei.RiskLevelPass {
		ChatListDealWith("[图片]", fromUserId, toUserId, enum.MsgListTextColorNormal)
		ChatListDealWith("[图片]", toUserId, fromUserId, enum.MsgListTextColorNormal)
		//处理消息未读数
		AddNotReadNum(toUserId, fromUserId)
	} else {
		ChatListDealWith("[图片]", fromUserId, toUserId, enum.MsgListTextColorNormal)
	}
	return response_im.ResponseMsg{
		Sn:    0,
		Msg:   baseData,
		Extra: extra,
	}
}

// SendAudioMsg 发送音频消息
func (im *ImOneService) SendAudioMsg(c *gin.Context, url, fromUserId, toUserId string, length int, clientType, extra string) response_im.ResponseMsg {
	if helper.OfficialAccount(toUserId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 黑名单拦截
	if user.IsBlacklist(toUserId, fromUserId, "0", enum.BlacklistTypeUser) {
		panic(error2.I18nError{
			Code: error2.ErrorCodePrivateChatIsBlacklist,
			Msg:  nil,
		})
	}
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId: messageId,
		Timestamp: timestamp,
		MsgType:   enum.MsgAudio,
		Code:      enum.USER_AUDIO_MSG,
		MsgData: message.MsgAudio{
			Content: helper.FormatImgUrl(url),
			Length:  length,
		},
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RiskLevel:    shumei.RiskLevelPass,
	}
	riskRes, ok := new(shumei.ShuMei).PrivateChatAudioCheck(fromUserId, helper.FormatImgUrl(url), coreSnowflake.GetSnowId())
	if !ok {
		baseData.RiskLevel = shumei.RiskLevelReject
		coreLog.LogError("SendAudioMsg reject %s", riskRes)
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.USER_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           enum.USER_AUDIO_MSG,
		RiskData:       riskRes,
		RiskLevel:      baseData.RiskLevel,
	}
	service_kafka.New().PushPrivateChat(kafkaData)

	if baseData.RiskLevel == shumei.RiskLevelPass {
		ChatListDealWith("[语音]", fromUserId, toUserId, enum.MsgListTextColorNormal)
		ChatListDealWith("[语音]", toUserId, fromUserId, enum.MsgListTextColorNormal)
		//处理消息未读数
		AddNotReadNum(toUserId, fromUserId)
	} else {
		ChatListDealWith("[语音]", fromUserId, toUserId, enum.MsgListTextColorNormal)
	}
	return response_im.ResponseMsg{
		Sn:    0,
		Msg:   baseData,
		Extra: extra,
	}
}

// SendCustomMsg 发送自定义消息
func (im *ImOneService) SendCustomMsg(c *gin.Context, fromUserId, toUserId, clientType string, msgData any) response_im.ResponseMsg {
	if helper.OfficialAccount(toUserId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId:    messageId,
		Timestamp:    timestamp,
		MsgType:      enum.MsgCustom,
		Code:         enum.USER_CUSTOM_MSG,
		MsgData:      msgData,
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.USER_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           enum.USER_CUSTOM_MSG,
	}
	service_kafka.New().PushPrivateChat(kafkaData)
	return response_im.ResponseMsg{
		Sn:  0,
		Msg: baseData,
	}
}

// 系统相关消息
func (im *ImOneService) SendNotice(c *gin.Context, fromUserId string, toUserId []string, msgType string, msgData any, code int) {

	go func() {
		messageType := enum.USER_MESSAGE_TYPE
		fromUserInfo := user.GetUserInfo(fromUserId, helper.GetClientType(c))
		if len(toUserId) > 0 {
			for _, toId := range toUserId {
				messageId := coreSnowflake.GetSnowId()
				timestamp := time.Now().UnixMicro()
				baseData := message.BaseMsg{
					MessageId:    messageId,
					Timestamp:    timestamp,
					MsgType:      msgType,
					Code:         code,
					MsgData:      msgData,
					FromUserInfo: fromUserInfo,
					ToUserInfo:   user.GetUserInfo(toId, helper.GetClientType(c)),
				}
				kafkaData := service_kafka.KafkaMessage{
					MessageType: messageType,
					FromUserId:  fromUserId,
					ToUserId:    toId,
					Data:        baseData,
					MessageId:   messageId,
					Timestamp:   timestamp,
					Code:        code,
					RiskLevel:   shumei.RiskLevelPass,
				}
				service_kafka.New().PushPrivateChat(kafkaData)
			}
		} else {
			//全平台用户推送消息
			var offsetUserId int64 = 0
			var limit = 500
			for {
				userIds := new(dao.UserDao).GetUserIdsOffsetId(offsetUserId, limit)
				if len(userIds) > 0 {
					offsetUserId = cast.ToInt64(userIds[len(userIds)-1])
				}
				for _, userId := range userIds {
					messageId := coreSnowflake.GetSnowId()
					timestamp := time.Now().UnixMicro()
					baseData := message.BaseMsg{
						MessageId:    messageId,
						Timestamp:    timestamp,
						MsgType:      msgType,
						Code:         code,
						MsgData:      msgData,
						FromUserInfo: fromUserInfo,
						ToUserInfo:   user.GetUserInfo(userId, helper.GetClientType(c)),
					}
					kafkaData := service_kafka.KafkaMessage{
						MessageType: messageType,
						FromUserId:  fromUserId,
						ToUserId:    userId,
						Data:        baseData,
						MessageId:   messageId,
						Timestamp:   timestamp,
						Code:        code,
						RiskLevel:   shumei.RiskLevelPass,
					}
					service_kafka.New().PushPrivateChat(kafkaData)
				}
				if len(userIds) < limit {
					break
				}
			}
		}
	}()
}
