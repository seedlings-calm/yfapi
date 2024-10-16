package service_kafka

import (
	"encoding/json"
	"sync"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/core/coreMq/kafka"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/riskCheck/shumei"
	typedef_enum "yfapi/typedef/enum"
	im3 "yfapi/typedef/response/im"
)

var imProducerOnce sync.Once
var imProducer *ImProducer

// 消息结构
type KafkaMessage struct {
	MessageId      string `json:"messageId"`
	MessageType    int    `json:"messageType"`    //消息类型 0 全服，1 房间 2 用户
	FromUserId     string `json:"fromUserId"`     //当前用户ID
	FromClientType string `json:"fromClientType"` //当前用户客户端
	ToClientType   string `json:"toClientType"`   //目标客户端类型
	ToUserId       string `json:"toUserId"`       //目标用户
	ToRoomId       string `json:"toRoomId"`       //目标房间
	Code           int    `json:"code"`           //信道码
	Data           any    `json:"data"`           //推送的消息体
	Timestamp      int64  `json:"timestamp"`
	RiskData       string `json:"riskData"`
	RiskLevel      string `json:"riskLevel"` //拦截等级
}

type ImProducer struct {
	action      *kafka.Kafka //动作消息推送
	gift        *kafka.Kafka //礼物消息生产者
	privateChat *kafka.Kafka //私聊消息生产者
	publicChat  *kafka.Kafka //公屏消息生产者
}

func New() *ImProducer {
	imProducerOnce.Do(func() {
		action, err := kafka.NewAsyncProducer(kafka.KafkaConf{
			Addr:  coreConfig.GetHotConf().Kafka.Action.Addr,
			Topic: coreConfig.GetHotConf().Kafka.Action.Topic,
		})
		if err != nil {
			panic(err)
		}
		gift, err := kafka.NewAsyncProducer(kafka.KafkaConf{
			Addr:  coreConfig.GetHotConf().Kafka.Gift.Addr,
			Topic: coreConfig.GetHotConf().Kafka.Gift.Topic,
		})
		if err != nil {
			panic(err)
		}
		privateChat, err := kafka.NewAsyncProducer(kafka.KafkaConf{
			Addr:  coreConfig.GetHotConf().Kafka.PrivateChat.Addr,
			Topic: coreConfig.GetHotConf().Kafka.PrivateChat.Topic,
		})
		if err != nil {
			panic(err)
		}
		publicChat, err := kafka.NewAsyncProducer(kafka.KafkaConf{
			Addr:  coreConfig.GetHotConf().Kafka.PublicChat.Addr,
			Topic: coreConfig.GetHotConf().Kafka.PublicChat.Topic,
		})
		if err != nil {
			panic(err)
		}
		imProducer = &ImProducer{
			action:      action,
			gift:        gift,
			privateChat: privateChat,
			publicChat:  publicChat,
		}
	})
	return imProducer
}

// 推送动作消息
func (i *ImProducer) PushAction(message ...KafkaMessage) {
	if len(message) == 0 {
		return
	}
	msg := []string{}
	for _, v := range message {
		marshal, err := json.Marshal(v)
		if err != nil {
			coreLog.LogError("Action err:%+v", err)
			continue
		}
		msg = append(msg, string(marshal))
	}
	i.action.AsyncPush(msg)
}

// 推送礼物消息
func (i *ImProducer) PushGift(message ...KafkaMessage) {
	if len(message) == 0 {
		return
	}
	msg := []string{}
	for _, v := range message {
		marshal, err := json.Marshal(v)
		if err != nil {
			coreLog.LogError("PushGift err:%+v", err)
			continue
		}
		msg = append(msg, string(marshal))
		msgStore := im3.ResponseMsg{
			Sn:  0,
			Msg: v.Data,
		}
		messageBytes, _ := json.Marshal(msgStore)
		new(dao.ChatStoreDao).WritePublicChatMessage(model.PublicChat{
			MessageId:  v.MessageId,
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			RoomId:     v.ToRoomId,
			Message:    string(messageBytes),
			RiskData:   nil,
			Timestamp:  v.Timestamp,
			Status:     typedef_enum.ChatMessagePass,
		})
	}
	i.gift.AsyncPush(msg)
}

// 推送私聊消息
func (i *ImProducer) PushPrivateChat(message ...KafkaMessage) {
	if len(message) == 0 {
		return
	}
	msg := []string{}
	for _, v := range message {
		msgStore := im3.ResponseMsg{
			Sn:  0,
			Msg: v.Data,
		}
		messageBytes, _ := json.Marshal(msgStore)
		status := typedef_enum.ChatMessagePass
		switch v.RiskLevel {
		case shumei.RiskLevelPass:
			marshal, err := json.Marshal(v)
			if err != nil {
				coreLog.LogError("PushPrivateChat err:%+v", err)
				continue
			}
			msg = append(msg, string(marshal))
		case shumei.RiskLevelReview:
			marshal, err := json.Marshal(v)
			if err != nil {
				coreLog.LogError("PushPrivateChat err:%+v", err)
				continue
			}
			msg = append(msg, string(marshal))
			status = typedef_enum.ChatMessageReview
		case shumei.RiskLevelReject:
			status = typedef_enum.ChatMessageReject
		}
		var riskData *string
		if len(v.RiskData) != 0 {
			riskData = &v.RiskData
		}
		new(dao.ChatStoreDao).WritePrivateChatMessage(model.PrivateChat{
			MessageId:  v.MessageId,
			UniteId:    helper.GetUniteId(v.FromUserId, v.ToUserId),
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			Message:    string(messageBytes),
			Read:       typedef_enum.PrivateChatMessageUnRead,
			Timestamp:  v.Timestamp,
			Status:     status,
			RiskData:   riskData,
		})
	}
	i.privateChat.AsyncPush(msg)
}

// 推送公屏消息
func (i *ImProducer) PushPublicChat(message ...KafkaMessage) {
	if len(message) == 0 {
		return
	}
	msg := []string{}
	for _, v := range message {
		msgStore := im3.ResponseMsg{
			Sn:  0,
			Msg: v.Data,
		}
		messageBytes, _ := json.Marshal(msgStore)
		status := typedef_enum.ChatMessagePass
		switch v.RiskLevel {
		case shumei.RiskLevelPass:
			marshal, err := json.Marshal(v)
			if err != nil {
				coreLog.LogError("PushPublicChat err:%+v", err)
				continue
			}
			msg = append(msg, string(marshal))
		case shumei.RiskLevelReview:
			marshal, err := json.Marshal(v)
			if err != nil {
				coreLog.LogError("PushPublicChat err:%+v", err)
				continue
			}
			msg = append(msg, string(marshal))
			status = typedef_enum.ChatMessageReview
		case shumei.RiskLevelReject:
			status = typedef_enum.ChatMessageReject
		}
		var riskData *string
		if len(v.RiskData) != 0 {
			riskData = &v.RiskData
		}
		new(dao.ChatStoreDao).WritePublicChatMessage(model.PublicChat{
			MessageId:  v.MessageId,
			FromUserId: v.FromUserId,
			ToUserId:   v.ToUserId,
			RoomId:     v.ToRoomId,
			Message:    string(messageBytes),
			Timestamp:  v.Timestamp,
			Status:     status,
			RiskData:   riskData,
		})
	}
	i.publicChat.AsyncPush(msg)
}

// 推送消息不
func (i *ImProducer) PushMessage(message ...KafkaMessage) {
	if len(message) == 0 {
		return
	}
	msg := []string{}
	for _, v := range message {
		marshal, _ := json.Marshal(v)
		msg = append(msg, string(marshal))
	}
	i.privateChat.AsyncPush(msg)
}
