package service_gift

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/typedef/redisKey"
	response_index "yfapi/typedef/response/index"
	"yfapi/util/easy"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

func GetTopMsg(page, pageSize int64) []response_index.TopMsgRes {
	redisCli := coreRedis.GetChatroomRedis()
	ctx := context.Background()
	msgKey := redisKey.TopMsgKey()
	start := (page - 1) * pageSize
	end := start + pageSize - 1
	redisRes := redisCli.LRange(ctx, msgKey, start, end).Val()
	resLen := len(redisRes)
	res := make([]response_index.TopMsgRes, resLen)
	if resLen > 0 {
		for _, v := range redisRes {
			var item map[string]string
			if err := json.Unmarshal([]byte(v), &item); err != nil {
				apdItem := response_index.TopMsgRes{
					UserId:    item["user_id"],
					Nickname:  item["nickname"],
					Avatar:    item["avatar"],
					Types:     item["types"],
					Operate:   item["operate"],
					ToUser:    item["to_user"],
					GiftImg:   item["gift_img"],
					GiftName:  item["gift_name"],
					GiftCount: cast.ToInt(item["gift_count"]) * cast.ToInt(item["combo_count"]),
				}
				if item["to_user"] != "全麦" {
					json.Unmarshal([]byte(item["to_user"]), &apdItem.ToUser)
				}
				res = append(res, apdItem)
			}
		}

	}
	return res
}

// 处理过期的连击数据
func TopMsgQueueExpired() {
	coreLog.Info("开启连击过期逻辑处理")
	redisCli := coreRedis.GetChatroomRedis()
	ctx := context.Background()
	zsetKey := redisKey.TopMsgEqueueKey()

	for {
		// 查找过期的任务
		now := float64(time.Now().Unix())
		tasks, err := redisCli.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
			Min:    "-inf",
			Max:    fmt.Sprintf("%f", now),
			Offset: 0,
			Count:  1,
		}).Result()
		if err != nil {
			log.Fatalf("Failed to get tasks: %v", err)
		}

		if len(tasks) == 0 {
			// 如果没有任务到期，稍等再检查
			time.Sleep(200 * time.Millisecond)
			continue
		}

		taskID := tasks[0]
		// 从 Redis 检查 key 是否存在
		exists, err := redisCli.Exists(ctx, taskID).Result()
		if err != nil {
			log.Fatalf("Failed to check key existence: %v", err)
		}

		if exists == 0 {
			// Key 已过期，处理连击结束逻辑
			// 从 HASH 获取处理回调逻辑的额外信息
			callbackKey := redisKey.TopMsgCallBack(taskID)
			callbackData, err := redisCli.HGetAll(ctx, callbackKey).Result()
			if err != nil {
				log.Printf("Failed to retrieve callback info: %v", err)
				// 移除无效任务
				redisCli.ZRem(ctx, zsetKey, taskID)
				continue
			}
			// 处理任务逻辑
			// 在此处理结束连击的逻辑，例如更新数据库或发送通知
			//TODO: 判断钻石数是否大于等于3000
			if cast.ToInt(callbackData["total_diamond"])*cast.ToInt(callbackData["combo_count"]) >= 3000 {
				msgKey := redisKey.TopMsgKey()
				callbackData["avatar"] = helper.FormatImgUrl(callbackData["avatar"])
				callbackData["gift_img"] = helper.FormatImgUrl(callbackData["gift_img"])
				if callbackData["to_user"] != "全麦" {
					var toUsers []struct {
						Id       string `json:"id"`
						Avatar   string `json:"avatar"`
						Nickname string `json:"nickname"`
					}
					toUserArr := strings.Split(callbackData["to_user"], ",")
					coreDb.GetMasterDb().Model(model.User{}).Where("id in ?", toUserArr).Select("id", "avatar", "nickname").Scan(&toUsers)

					for _, v := range toUsers {
						v.Avatar = helper.FormatImgUrl(v.Avatar)
						callbackData["to_user"] = easy.JSONStringFormObject(v)
						redisCli.LPush(ctx, msgKey, easy.JSONStringFormObject(callbackData))
					}
				} else {
					redisCli.LPush(ctx, msgKey, easy.JSONStringFormObject(callbackData))
				}
				redisCli.LTrim(ctx, msgKey, 0, 199)
			}
			// 删除 HASH 和 ZSET 中的任务
			redisCli.Del(ctx, callbackKey)
			redisCli.ZRem(ctx, zsetKey, taskID)
		}

	}
}
