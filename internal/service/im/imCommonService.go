package service_im

import (
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/core/coreSnowflake"
	"yfapi/internal/helper"
	service_kafka "yfapi/internal/service/kafka"
	"yfapi/internal/service/riskCheck/shumei"
	"yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/message"
)

type ImCommonService struct {
}

//	@Description:通用通知
//	@receiver im
//	@param fromUserId 发送者
//	@param toUserId 接收者
//	@param roomId 房间ID
//	@param msgData 消息内容
//	@param code 消息信道码
//	@param messageType 消息类型
//
// toUserId有值则推送给目标用户，toUserId为空 roomId不为空推送全房间，toUserId为空 roomId为空推送全服，
func (im *ImCommonService) Send(c *gin.Context, fromUserId string, toUserId []string, roomId string, msgType string, msgData any, code int) {
	go func() {
		messageType := enum.ROOM_MESSAGE_TYPE
		switch {
		case len(toUserId) > 0: //推送指定用户
			messageType = enum.USER_MESSAGE_TYPE
			fromUserInfo := user.GetUserInfo(fromUserId, helper.GetClientType(c))
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
					ToRoomId:    roomId,
					Data:        baseData,
					MessageId:   messageId,
					Timestamp:   timestamp,
					Code:        code,
					RiskLevel:   shumei.RiskLevelPass,
				}
				service_kafka.New().PushMessage(kafkaData)
			}
		case len(toUserId) == 0 && roomId != "": //推送全房间
			messageType = enum.ROOM_MESSAGE_TYPE
			messageId := coreSnowflake.GetSnowId()
			timestamp := time.Now().UnixMicro()
			baseData := message.BaseMsg{
				MessageId:    messageId,
				Timestamp:    timestamp,
				MsgType:      msgType,
				Code:         code,
				MsgData:      msgData,
				FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
				RoomId:       roomId,
			}
			kafkaData := service_kafka.KafkaMessage{
				MessageType: messageType,
				FromUserId:  fromUserId,
				ToRoomId:    roomId,
				Data:        baseData,
				MessageId:   messageId,
				Timestamp:   timestamp,
				Code:        code,
			}
			service_kafka.New().PushMessage(kafkaData)
		case len(toUserId) == 0 && len(roomId) == 0: //推送全服在线用户
			messageType = enum.ALL_SERVICE_MESSAGE_TYPE
			messageId := coreSnowflake.GetSnowId()
			timestamp := time.Now().UnixMicro()
			baseData := message.BaseMsg{
				MessageId:    messageId,
				Timestamp:    timestamp,
				MsgType:      msgType,
				Code:         code,
				MsgData:      msgData,
				FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
				RoomId:       roomId,
			}
			kafkaData := service_kafka.KafkaMessage{
				MessageType: messageType,
				FromUserId:  fromUserId,
				ToRoomId:    roomId,
				Data:        baseData,
				MessageId:   messageId,
				Timestamp:   timestamp,
				Code:        code,
			}
			service_kafka.New().PushMessage(kafkaData)
		}
	}()
}

type ToClientUserInfo struct {
	UserId string
	Client string
}

// 向指定得客户端用户推送消息
func (im *ImCommonService) SendClientUser(c *gin.Context, fromUserId string, toUserIds []ToClientUserInfo, msgType string, msgData any, code int) {
	go func() {
		messageType := enum.USER_MESSAGE_CLIENT_TYYPE
		switch {
		case len(toUserIds) > 0: //推送指定用户
			fromUserInfo := user.GetUserInfo(fromUserId, helper.GetClientType(c))
			for _, v := range toUserIds {
				messageId := coreSnowflake.GetSnowId()
				timestamp := time.Now().UnixMicro()
				baseData := message.BaseMsg{
					MessageId:    messageId,
					Timestamp:    timestamp,
					MsgType:      msgType,
					Code:         code,
					MsgData:      msgData,
					FromUserInfo: fromUserInfo,
				}
				kafkaData := service_kafka.KafkaMessage{
					MessageType:  messageType,
					FromUserId:   fromUserId,
					ToUserId:     v.UserId,
					ToClientType: v.Client,
					Data:         baseData,
					MessageId:    messageId,
					Timestamp:    timestamp,
					Code:         code,
					RiskLevel:    shumei.RiskLevelPass,
				}
				service_kafka.New().PushMessage(kafkaData)
			}
		}
	}()
}
