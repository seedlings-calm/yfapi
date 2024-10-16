package service_im

import (
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_goods "yfapi/internal/service/goods"
	"yfapi/internal/service/kafka"
	"yfapi/internal/service/rankList"
	"yfapi/internal/service/riskCheck/shumei"
	"yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/message"
	"yfapi/typedef/response/im"
)

type ImPublicService struct {
}

// 发送房间公屏文本消息
func (im *ImPublicService) SendTextMsg(c *gin.Context, content, fromUserId, toUserId, roomId, clientType, extra string) response_im.ResponseMsg {
	mute := new(dao.UserMuteListDao).FindOne(fromUserId, roomId)
	if mute.ID > 0 {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserMuteMsg,
			Msg:  map[string]any{"minute": int(mute.EndTime.Sub(mute.StartTime).Minutes())},
		})
	}
	msgData := message.MsgText{Content: content}
	// 查询用户气泡框
	ses, _ := service_goods.UserGoods{}.GetGoodsByKeys(fromUserId, true, enum.GoodsTypeRCMB)
	if len(ses) > 0 {
		for _, v := range ses {
			if v.GoodsTypeKey == enum.GoodsTypeRCMB {
				msgData.Bubble = *v
			}
		}
	}
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId:    messageId,
		Timestamp:    timestamp,
		MsgType:      enum.MsgText,
		Code:         enum.ROOM_TEXT_MSG,
		MsgData:      msgData,
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RiskLevel:    shumei.RiskLevelPass,
		RoomId:       roomId,
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
		ToRoomId:       roomId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           enum.ROOM_TEXT_MSG,
		RiskLevel:      baseData.RiskLevel,
		RiskData:       riskRes,
	}
	service_kafka.New().PushPublicChat(kafkaData)
	//排行榜首发送公屏消息
	go func() {
		rankList.Instance().Calculate(rankList.CalculateReq{
			FromUserId: fromUserId,
			Types:      "publicMessage",
			RoomId:     roomId,
		})
	}()
	return response_im.ResponseMsg{
		Sn:    0,
		Msg:   baseData,
		Extra: extra,
	}
}

// 发送房间公屏图片消息
func (im *ImPublicService) SendImgMsg(c *gin.Context, url, fromUserId, toUserId, roomId string, width, height int, clientType string, extra string) response_im.ResponseMsg {
	mute := new(dao.UserMuteListDao).FindOne(fromUserId, roomId)
	if mute.ID > 0 {
		panic(error2.I18nError{
			Code: error2.ErrCodeHeMuteMsg,
			Msg:  map[string]any{"minute": int(mute.EndTime.Sub(mute.StartTime).Minutes())},
		})
	}
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId: messageId,
		Timestamp: timestamp,
		MsgType:   enum.MsgImg,
		Code:      enum.ROOM_IMG_MSG,
		MsgData: message.MsgImg{
			Content: helper.FormatImgUrl(url),
			Width:   width,
			Height:  height,
		},
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RiskLevel:    shumei.RiskLevelPass,
		RoomId:       roomId,
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
		ToRoomId:       roomId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           enum.ROOM_IMG_MSG,
		RiskLevel:      baseData.RiskLevel,
		RiskData:       riskRes,
	}
	service_kafka.New().PushPublicChat(kafkaData)
	//排行榜首发送公屏消息
	go func() {
		rankList.Instance().Calculate(rankList.CalculateReq{
			FromUserId: fromUserId,
			Types:      "publicMessage",
			RoomId:     roomId,
		})
	}()
	return response_im.ResponseMsg{
		Sn:    0,
		Msg:   baseData,
		Extra: extra,
	}
}

// 发送动作消息
func (im *ImPublicService) SendActionMsg(c *gin.Context, data any, fromUserId, toUserId, roomId, clientType string, code int) {
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
		RoomId:       roomId,
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.ROOM_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		ToRoomId:       roomId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           code,
	}
	service_kafka.New().PushAction(kafkaData)
	return
}

// 发送礼物消息
func (im *ImPublicService) SendGiftMsg(c *gin.Context, fromUserId, toUserId, roomId, clientType, comboKey string, giftCount, comboCount, totalGiftDiamond int, isBatch bool, giftInfo *model.GiftDTO) {
	messageId := coreSnowflake.GetSnowId()
	timestamp := time.Now().UnixMicro()
	baseData := message.BaseMsg{
		MessageId: messageId,
		Timestamp: timestamp,
		MsgType:   enum.MsgAction,
		Code:      enum.ROOM_GIFT_MSG,
		MsgData: message.MsgGift{
			GiftName:         giftInfo.GiftName,
			GiftImage:        helper.FormatImgUrl(giftInfo.GiftImage),
			AnimationUrl:     helper.FormatImgUrl(giftInfo.AnimationUrl),
			AnimationJsonUrl: helper.FormatImgUrl(giftInfo.AnimationJsonUrl),
			GiftCount:        giftCount,
			ComboCount:       comboCount,
			ComboKey:         comboKey,
			TotalGiftDiamond: totalGiftDiamond,
			IsBatch:          isBatch,
		},
		FromUserInfo: user.GetUserInfo(fromUserId, helper.GetClientType(c)),
		ToUserInfo:   user.GetUserInfo(toUserId, helper.GetClientType(c)),
		RoomId:       roomId,
		ExtraInfo:    isBatch,
	}
	kafkaData := service_kafka.KafkaMessage{
		MessageType:    enum.ROOM_MESSAGE_TYPE,
		FromUserId:     fromUserId,
		FromClientType: clientType,
		ToUserId:       toUserId,
		ToRoomId:       roomId,
		Data:           baseData,
		MessageId:      messageId,
		Timestamp:      timestamp,
		Code:           enum.ROOM_GIFT_MSG,
	}
	service_kafka.New().PushGift(kafkaData)
	return
}

// 发送房间公屏自定义消息
func (im *ImPublicService) SendCustomMsg(roomId string, msgData any, code int) {
	go func() {
		messageId := coreSnowflake.GetSnowId()
		timestamp := time.Now().UnixMicro()
		baseData := message.MsgCustom{
			MessageId: messageId,
			Timestamp: timestamp,
			MsgType:   enum.MsgCustom,
			Code:      code,
			MsgData:   msgData,
			RoomId:    roomId,
		}
		kafkaData := service_kafka.KafkaMessage{
			MessageType: enum.ROOM_MESSAGE_TYPE,
			ToRoomId:    roomId,
			Data:        baseData,
			MessageId:   messageId,
			Timestamp:   timestamp,
			Code:        code,
		}
		service_kafka.New().PushAction(kafkaData)
	}()
}
