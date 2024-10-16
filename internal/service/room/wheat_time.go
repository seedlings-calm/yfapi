package service_room

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	"yfapi/internal/model"
	"yfapi/typedef/redisKey"
)

// redis存储统计字段，统一更改地方
const (
	EnterCount      = "enterCount"      //进房总人数
	EnterTimes      = "enterTimes"      //进房总人次
	RewardCount     = "rewardCount"     //打赏总金额
	RewardTimes     = "rewardTimes"     //打赏次数
	RewardUserCount = "rewardUserCount" //打赏人数
	OnwheatTime     = "onwheatTime"     //开播开始时间
)

//直播间的统计逻辑
/**
 * @description  操作参数，处理存储数据,处理加入房间,打赏，
 * @param userId  string true 触发用户ID
 * @param roomId  string true 房间ID
 * @param num float64 true 触发时间的数字
 * @param keys int strue  修改字段  EnterCount,EnterTimes ...
 */
func DoRoomWheatTimeOperation(userId string, roomId string, num float64, keys ...string) {
	redisCli := coreRedis.GetChatroomRedis()
	ctx := context.Background()
	dataKey := redisKey.RoomWheatTimeCacheKey(roomId)
	res, err := getStoreRoomWHeatTime(roomId)
	if err != nil {
		coreLog.Error("查询主数据失败:", err)
		return
	}
	if _, ok := res[OnwheatTime]; !ok {
		return
	}

	for _, v := range keys {
		var ok = true
		if v == EnterCount { //需要处理去重进入房间用户人数
			ok, _ = isUserIDInSet(redisKey.RoomWheatTimeJoinUser(roomId), userId)
			ok = !ok
		}
		if v == RewardUserCount { //需要处理去重打赏用户人数
			ok, _ = isUserIDInSet(redisKey.RoomWheatTimeGiftUser(roomId), userId)
			ok = !ok
		}
		if ok {
			err := redisCli.HIncrByFloat(ctx, dataKey, v, num).Err()
			if err != nil {
				coreLog.Error(fmt.Sprintf("更新房间开播统计数据出错:%s:%f", v, num))
			}
		}
	}

}

// TODO:如果上主持麦时，房间有用户并且有消费如何处理
func InitStoreRoomWHeatTime(userId string, roomInfo *model.Room) (err error) {
	redisCli := coreRedis.GetChatroomRedis()
	ctx := context.Background()
	dataKey := redisKey.RoomWheatTimeCacheKey(roomInfo.Id)
	res, err := redisCli.HGetAll(ctx, dataKey).Result()
	if err != nil {
		return err
	}
	if len(res) > 0 {
		return
	}
	log.Println("房间开播初始化房间统计信息")
	_, err = storeRoomWheatTime(userId, roomInfo.Id, roomInfo.GuildId, roomInfo.LiveType)
	if err != nil {
		return err
	}
	return
}

/**
 * @description 房间开播统计数据转存到mysql
 * @param roomId string true 房间ID
 */
func StoreRoomWHeatTimeToMysql(roomId string) (err error) {
	data, err := getStoreRoomWHeatTime(roomId)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		coreLog.Error("存储直播数据到mysql时，未查询到记录")
		return
	}
	var (
		by         []byte
		insertData = new(model.DoRoomWheatTime)
	)
	by, err = json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(by, insertData)
	if err != nil {
		return
	}
	OnwheatTime, _ := insertData.OnwheatTime.(string)
	parseTime, _ := time.ParseInLocation(time.DateTime, OnwheatTime, time.Local)

	nowTime := time.Now()
	insertData.UpwheatTime = nowTime.Format(time.DateTime)
	insertData.OnTime = nowTime.Sub(parseTime).Seconds()
	err = coreDb.GetMasterDb().Model(model.RoomWheatTime{}).Create(insertData).Error
	if err != nil {
		return err
	}
	coreRedis.GetChatroomRedis().Del(context.Background(), redisKey.RoomWheatTimeJoinUser(roomId))
	coreRedis.GetChatroomRedis().Del(context.Background(), redisKey.RoomWheatTimeGiftUser(roomId))
	coreRedis.GetChatroomRedis().Del(context.Background(), redisKey.RoomWheatTimeCacheKey(roomId))
	return nil
}

// storeRoomWheatTime 将麦位信息存储到 Redis 哈希中
func storeRoomWheatTime(userId, roomId, guildId string, liveType int) (*model.DoRoomWheatTime, error) {
	nowTime := time.Now()
	wheatRoomTime := &model.DoRoomWheatTime{
		RoomID:          roomId,
		UserID:          userId,
		GuildID:         guildId,
		RoomType:        liveType,
		OnwheatTime:     nowTime.Format(time.DateTime),
		StatDate:        nowTime.Format(time.DateOnly),
		EnterCount:      0,
		EnterTimes:      0,
		RewardCount:     0,
		RewardTimes:     0,
		RewardUserCount: 0,
	}
	redisCli := coreRedis.GetChatroomRedis()
	ctx := context.Background()
	_, err := redisCli.HMSet(ctx, redisKey.RoomWheatTimeCacheKey(roomId), wheatRoomTime).Result()
	return wheatRoomTime, err
}

func getStoreRoomWHeatTime(roomId string) (res map[string]string, err error) {
	redisCli := coreRedis.GetChatroomRedis()
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		res, err = redisCli.HGetAll(ctx, redisKey.RoomWheatTimeCacheKey(roomId)).Result()
		if err != nil {
			continue
		}
	}
	if len(res) == 0 {
		return nil, errors.New("getStoreRoomWHeatTime 未查询到记录")
	}
	return
}

// 判断用户ID是否存在于Redis Set中
func isUserIDInSet(setKey string, userID string) (bool, error) {
	// 使用 SISMEMBER 命令检查ID是否存在
	ctx := context.Background()
	exists, err := coreRedis.GetChatroomRedis().SIsMember(ctx, setKey, userID).Result()
	if err != nil {
		return false, err
	}
	if !exists { //如果不存在，写入进去
		coreRedis.GetChatroomRedis().SAdd(ctx, setKey, userID)
	}
	return exists, nil
}
