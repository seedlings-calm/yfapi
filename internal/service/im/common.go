package service_im

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/typedef/enum"
	"yfapi/typedef/message"
	common_data "yfapi/typedef/redisKey"
)

// 会话列表处理
func ChatListDealWith(showContent, fromUserId, toUserId, textColor string) {
	redisClient := coreRedis.GetImRedis()
	timestamp := time.Now().UnixMilli()
	sessionId := GeneSessionId(fromUserId, toUserId)
	result, _ := redisClient.HGet(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId).Result()
	fromModelData := message.OneMsgListModel{}
	if len(result) == 0 {
		//记录双方的会话列表
		fromModelData = message.OneMsgListModel{
			ShowContent: showContent,
			Timestamp:   timestamp,
			TextColor:   textColor,
			ToUserId:    toUserId,
		}
	} else {
		err := json.Unmarshal([]byte(result), &fromModelData)
		if err != nil {
			coreLog.LogError("ChatListDealWith unmarshal err fromUserId:%s,toUserId:%s", fromUserId, toUserId)
			return
		}
		fromModelData.ShowContent = showContent
		fromModelData.TextColor = textColor
		fromModelData.Timestamp = timestamp
		if fromModelData.IsTop {
			timestamp *= 2
		}
	}
	//会话列表id
	redisClient.ZAdd(context.Background(), common_data.ImOneSessionSortId(fromUserId), redis.Z{
		Member: sessionId,
		Score:  float64(timestamp),
	})
	fromModelDataByte, _ := json.Marshal(fromModelData)
	redisClient.HSet(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId, fromModelDataByte)
}

// 增加消息未读数
func AddNotReadNum(userId, chatUserId string) {
	redisClient := coreRedis.GetImRedis()
	redisClient.HIncrBy(context.Background(), common_data.ImOneMsgNotReadNum(userId), chatUserId, 1)
}

// @Description: 清空消息未读
// @receiver im
// @param userId
// @param chatUserId
func ClearNotReadNum(userId, chatUserId string) {
	redisClient := coreRedis.GetImRedis()
	redisClient.HSet(context.Background(), common_data.ImOneMsgNotReadNum(userId), chatUserId, 0)
}

// 生成会话列表id
func GeneSessionId(fromUserId, toUserId string) string {
	return fromUserId + toUserId
}

// @Description: 置顶会话
// @receiver im
// @param fromUserId
// @param toUserId
func TopSession(fromUserId, toUserId string) {
	redisClient := coreRedis.GetImRedis()
	timestamp := time.Now().UnixMilli() * 2
	sessionId := GeneSessionId(fromUserId, toUserId)
	//会话列表id
	redisClient.ZAdd(context.Background(), common_data.ImOneSessionSortId(fromUserId), redis.Z{
		Member: sessionId,
		Score:  float64(timestamp),
	})
	result, _ := redisClient.HGet(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId).Result()
	if len(result) == 0 {
		coreLog.LogError("TopSession err Hget not found fromUserId:%s,toUserId:%s", fromUserId, toUserId)
		return
	}
	resultData := message.OneMsgListModel{}
	err := json.Unmarshal([]byte(result), &resultData)
	if err != nil {
		coreLog.LogError("TopSession unmarshal err fromUserId:%s,toUserId:%s", fromUserId, toUserId)
		return
	}
	resultData.IsTop = true
	fromModelDataByte, _ := json.Marshal(resultData)
	redisClient.HSet(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId, fromModelDataByte)
}

// @Description: 删除会话
// @receiver im
// @param fromUserId
// @param toUserId
// @return error
func DelSession(fromUserId, toUserId string) {
	redisClient := coreRedis.GetImRedis()
	sessionId := GeneSessionId(fromUserId, toUserId)
	_, err := redisClient.ZRem(context.Background(), common_data.ImOneSessionSortId(fromUserId), sessionId).Result()
	if err != nil {
		coreLog.Error("DelSession ZRem err:%+v,fromUserId:%s,toUserId", err, fromUserId, toUserId)
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	_, err = redisClient.HDel(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId).Result()
	if err != nil {
		coreLog.Error("DelSession HDel err:%+v,fromUserId:%s,toUserId", err, fromUserId, toUserId)
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	new(dao.ChatStoreDao).SetPrivateChatClearTimestamp(fromUserId, toUserId)
	return
}

// @Description: 取消置顶
// @receiver im
// @param fromUserId
// @param toUserId
func UnTopSession(fromUserId, toUserId string) {
	redisClient := coreRedis.GetImRedis()
	sessionId := GeneSessionId(fromUserId, toUserId)
	result, _ := redisClient.HGet(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId).Result()
	if len(result) == 0 {
		coreLog.LogError("UnTopSession err Hget not found fromUserId:%s,toUserId:%s", fromUserId, toUserId)
		return
	}
	resultData := message.OneMsgListModel{}
	err := json.Unmarshal([]byte(result), &resultData)
	if err != nil {
		coreLog.LogError("UnTopSession unmarshal err fromUserId:%s,toUserId:%s", fromUserId, toUserId)
		return
	}
	resultData.IsTop = false
	fromModelDataByte, _ := json.Marshal(resultData)
	redisClient.HSet(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId, fromModelDataByte)
	//会话列表id
	redisClient.ZAdd(context.Background(), common_data.ImOneSessionSortId(fromUserId), redis.Z{
		Member: sessionId,
		Score:  float64(resultData.Timestamp),
	})
}

// 返回用户对应得会话类型
func GetSessionListTypes(userId string) int {
	res := 0
	switch userId {
	case enum.SystematicUserId:
		res = 1
	case enum.OfficialUserId:
		res = 2
	case enum.InteractiveUserId:
		res = 3
	}
	return res
}

// 清空聊天历史记录
func ClearChatHistory(fromUserId, toUserId string) {
	redisClient := coreRedis.GetImRedis()
	sessionId := GeneSessionId(fromUserId, toUserId)
	result, _ := redisClient.HGet(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId).Result()
	fromModelData := message.OneMsgListModel{}
	if len(result) == 0 {
		return
	}
	err := json.Unmarshal([]byte(result), &fromModelData)
	if err != nil {
		coreLog.LogError("ClearChatHistory err fromUserId:%s,toUserId:%s", fromUserId, toUserId)
		return
	}
	fromModelData.ShowContent = ""
	fromModelData.TextColor = ""
	fromModelDataByte, _ := json.Marshal(fromModelData)
	redisClient.HSet(context.Background(), common_data.ImOneSessionList(fromUserId), sessionId, fromModelDataByte)
	new(dao.ChatStoreDao).SetPrivateChatClearTimestamp(fromUserId, toUserId)
	return
}
