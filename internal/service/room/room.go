package service_room

import (
	"context"
	"time"
	"yfapi/core/coreRedis"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	service_im "yfapi/internal/service/im"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	response_room "yfapi/typedef/response/room"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

// GetRoomHot 获取房间热度值
func GetRoomHot(roomId string, liveType int) (hot int, hotStr string) {
	if liveType == 0 {
		return
	}
	// 初始化房间热度值
	initRoomHot(roomId, liveType)
	redisCli := coreRedis.GetChatroomRedis()
	ctx := context.Background()
	// 房间每小时热度值排行总热度值(当前小时+上个小时)
	totalKey := redisKey.RoomHotTotalKey(liveType)
	// 当前热度值
	nowHot := redisCli.ZScore(ctx, totalKey, roomId).Val()
	hot = cast.ToInt(nowHot)
	hotStr = easy.NumberToW(hot)
	return
}

// UpdateRoomHotByJoinRoom 增加房间热度值进厅人数
func UpdateRoomHotByJoinRoom(roomId, userId string, liveType int) {
	ctx := context.Background()
	redisCli := coreRedis.GetChatroomRedis()
	key := redisKey.RoomHotJoinRoomKey(roomId)
	isExist := redisCli.HExists(ctx, key, userId).Val()
	if isExist {
		return
	}
	// 初始化房间热度值
	initRoomHot(roomId, liveType)
	// 本次增加的热度值
	addScore := float64(10)
	pipe := redisCli.Pipeline()
	// 房间每小时热度值排行总热度值(当前小时+上个小时)
	totalKey := redisKey.RoomHotTotalKey(liveType)
	// 房间每小时热度值缓存(当前小时)
	hotKey := redisKey.RoomHotKey(liveType)
	// 更新当前小时的热度值
	pipe.ZIncrBy(ctx, hotKey, addScore, roomId)
	// 更新当前小时总的热度值排行榜
	pipe.ZIncrBy(ctx, totalKey, addScore, roomId)
	// 设置当前小时的进房缓存
	pipe.HSet(ctx, key, userId, time.Now().Unix())
	pipe.Expire(ctx, key, 24*time.Hour)
	_, _ = pipe.Exec(ctx)
	// 通知房间热度值更新
	noticeRoomHotUpdate(roomId, liveType)
}

// UpdateRoomHotByChat 增加房间热度值发言次数
func UpdateRoomHotByChat(roomId, userId string) {
	ctx := context.Background()
	redisCli := coreRedis.GetChatroomRedis()
	key := redisKey.RoomHotPublicChat(roomId)
	countStr := redisCli.HGet(ctx, key, userId).Val()
	if cast.ToInt(countStr) >= 10 {
		return
	}
	// 房间信息
	roomInfo, err := new(dao.RoomDao).GetRoomById(roomId)
	if err != nil {
		return
	}
	if len(roomInfo.Id) > 0 {
		// 初始化房间热度值
		initRoomHot(roomId, roomInfo.LiveType)
		// 本次增加的热度值
		addScore := float64(1)
		pipe := redisCli.Pipeline()
		// 房间每小时热度值排行总热度值(当前小时+上个小时)
		totalKey := redisKey.RoomHotTotalKey(roomInfo.LiveType)
		// 房间每小时热度值缓存(当前小时)
		hotKey := redisKey.RoomHotKey(roomInfo.LiveType)
		// 更新当前小时的热度值
		pipe.ZIncrBy(ctx, hotKey, addScore, roomId)
		// 更新当前小时总的热度值排行榜
		pipe.ZIncrBy(ctx, totalKey, addScore, roomId)
		pipe.HIncrBy(ctx, key, userId, 1)
		pipe.Expire(ctx, key, 24*time.Hour)
		_, _ = pipe.Exec(ctx)
		// 通知房间热度值更新
		noticeRoomHotUpdate(roomId, roomInfo.LiveType)
	}

}

// UpdateRoomHotBySendGift 增加房间热度值礼物打赏
func UpdateRoomHotBySendGift(roomId string, exp float64, liveType int) {
	ctx := context.Background()
	// 初始化房间热度值
	initRoomHot(roomId, liveType)
	pipe := coreRedis.GetChatroomRedis().Pipeline()
	// 房间每小时热度值排行总热度值(当前小时+上个小时)
	totalKey := redisKey.RoomHotTotalKey(liveType)
	// 房间每小时热度值缓存(当前小时)
	hotKey := redisKey.RoomHotKey(liveType)
	// 更新当前小时的热度值
	pipe.ZIncrBy(ctx, hotKey, exp, roomId)
	// 更新当前小时总的热度值排行榜
	pipe.ZIncrBy(ctx, totalKey, exp, roomId)
	_, _ = pipe.Exec(ctx)
	// 通知房间热度值更新
	noticeRoomHotUpdate(roomId, liveType)
}

// 初始化当前小时的房间热度值
func initRoomHot(roomId string, liveType int) {
	ctx := context.Background()
	redisCli := coreRedis.GetChatroomRedis()
	// 当前时间是否有热度值
	totalKey := redisKey.RoomHotTotalKey(liveType)
	err := redisCli.ZScore(ctx, totalKey, roomId).Err()
	if err == redis.Nil {
		// 上个小时的热度值
		lastTimeKey := redisKey.RoomHotKey(liveType, time.Now().Add(-1*time.Hour).Format("2006010215"))
		lastScore := redisCli.ZScore(ctx, lastTimeKey, roomId).Val()
		// 当前小时热度值初始化
		currScore := float64(200)
		// 每小时热度值缓存
		hotKey := redisKey.RoomHotKey(liveType)
		pipe := redisCli.Pipeline()
		// 初始化每小时热度值缓存
		pipe.ZIncrBy(ctx, hotKey, currScore, roomId)
		pipe.Expire(ctx, hotKey, 24*time.Hour)
		// 初始化每小时总热度值排行
		pipe.ZIncrBy(ctx, totalKey, lastScore+currScore, roomId)
		pipe.Expire(ctx, totalKey, 24*time.Hour)
		_, _ = pipe.Exec(ctx)
	}
}

// 通知房间热度值更新
func noticeRoomHotUpdate(roomId string, liveType int) {
	_, hotStr := GetRoomHot(roomId, liveType)
	new(service_im.ImPublicService).SendCustomMsg(roomId, hotStr, typedef_enum.ROOM_HOT_UPDATE_MSG)
}

// 多端进房检测
func MultiJoinRoom(c *gin.Context, userId string, resp *response_room.CheckRoomResponse) {
	userClient := helper.GetClientType(c)
	for _, client := range typedef_enum.ClientTypeArray {
		if userClient != client {
			inRoomId, _ := coreRedis.GetChatroomRedis().Get(c, redisKey.UserInWhichRoom(userId, client)).Result()
			if len(inRoomId) > 0 {
				//用户在该房间内
				roomInfo, _ := new(dao.RoomDao).GetRoomById(inRoomId)
				if len(roomInfo.Id) > 0 {
					resp.Msg = i18n_msg.GetI18nMsg(c, i18n_msg.OtherClientInRoom, map[string]interface{}{"client": client, "roomName": roomInfo.Name, "liveTypeName": GetRoomLiveTypeName(c, roomInfo.LiveType)})
					resp.IsMulti = true
					return
				}
			}
		}
	}
}

// 多端离房操做
func MultiLeaveRoom(c *gin.Context, userId string) {
	userClient := helper.GetClientType(c)
	for _, client := range typedef_enum.ClientTypeArray {
		if userClient != client {
			inRoomId, _ := coreRedis.GetChatroomRedis().Get(c, redisKey.UserInWhichRoom(userId, client)).Result()
			if len(inRoomId) > 0 {
				//用户在该房间内
				roomInfo, _ := new(dao.RoomDao).GetRoomById(inRoomId)
				if len(roomInfo.Id) > 0 {
					ser := &RoomUsersOnlie{
						RoomId: inRoomId,
					}
					if err := ser.RemoveUserToRoom(c, userId, client, &roomInfo); err == nil {
						//推送消息
						userArr := []service_im.ToClientUserInfo{
							{
								UserId: userId,
								Client: client,
							},
						}
						msg := i18n_msg.GetI18nMsg(c, i18n_msg.JoinRoomFromOtherClient, map[string]interface{}{"client": userClient, "roomName": roomInfo.Name, "liveTypeName": GetRoomLiveTypeName(c, roomInfo.LiveType), "liveTypeName2": GetRoomLiveTypeName(c, roomInfo.LiveType)})
						new(service_im.ImCommonService).SendClientUser(c, typedef_enum.SystematicUserId, userArr, typedef_enum.MsgCustom, msg, typedef_enum.USER_JOIN_ROOM_FROM_OTHER_CLIENT)
					}
				}
			}
		}
	}
}

// 获取房间直播类型名
func GetRoomLiveTypeName(c *gin.Context, liveType int) (liveTypeName string) {
	switch liveType {
	case typedef_enum.LiveTypeChatroom:
		liveTypeName = i18n_msg.GetI18nMsg(c, i18n_msg.ChatRoomKey)
	case typedef_enum.LiveTypeAnchor:
		liveTypeName = i18n_msg.GetI18nMsg(c, i18n_msg.CastRoomKey)
	case typedef_enum.LiveTypePersonal:
		liveTypeName = i18n_msg.GetI18nMsg(c, i18n_msg.PersonalRoomKey)
	}
	return
}

// 根据热度排序获取指定数量的房间和热度值
func GetRoomHotsList(c *gin.Context, liveType int, num int64) ([]redis.Z, error) {
	users, err := coreRedis.GetChatroomRedis().ZRevRangeWithScores(c, redisKey.RoomHotTotalKey(liveType), 0, num-1).Result()
	if err != nil {
		return nil, err
	}
	return users, nil
}
