package service_room

import (
	"yfapi/core/coreRedis"
	"yfapi/typedef/redisKey"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// 1000贡献榜
type RoomUserDay struct {
	RoomId string
}

// 用户进入房间，初始化以天为单位的1000贡献榜单 TODO:未初始化
func (r *RoomUserDay) Add(c *gin.Context, userID string) error {
	redisCli := coreRedis.GetChatroomRedis()
	keys := redisKey.RoomUsersDayList(r.RoomId)
	err := redisCli.ZAddNX(c, keys, redis.Z{
		Score:  0,
		Member: userID,
	}).Err()
	if err != nil {
		return err
	}
	return nil
}

// 增加用户的贡献值
func (r *RoomUserDay) AddInc(c *gin.Context, userId string, score float64) {
	redisCli := coreRedis.GetChatroomRedis()
	keys := redisKey.RoomUsersDayList(r.RoomId)
	redisCli.ZIncrBy(c, keys, score, userId)
}

// 获取大于等于1000贡献榜的用户
func (r *RoomUserDay) GetDayUsers(c *gin.Context) []redis.Z {
	redisCli := coreRedis.GetChatroomRedis()
	keys := redisKey.RoomUsersDayList(r.RoomId)
	res, err := redisCli.ZRangeByScoreWithScores(c, keys, &redis.ZRangeBy{
		Min: "1000",
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil
	}
	return res
}
