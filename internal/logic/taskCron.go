package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/redisKey"

	"github.com/gin-gonic/gin"
)

// TaskCron
// @Description: 定时任务
type TaskCron struct {
}

// AutoDeleteUser
//
//	@Description: 自动删除 申请注销账号到期的
//	@receiver t
//	@param c *gin.Context -
func (t *TaskCron) AutoDeleteUser(c *gin.Context) {
	deleteApply := &dao.UserDeleteApplyDao{}
	applyList, _ := deleteApply.GetCanDeleteUserList()
	for _, apply := range applyList {
		_ = deleteApply.UpdateUserDeleteApply(apply.Id, typedef_enum.UserDeleteStatusDelete)
		_ = new(dao.UserDao).UpdateById(&model.User{Id: apply.UserId, Status: typedef_enum.UserStatusInvalid})
	}
}

/**
 * @description 房间统计数据定时处理跨天数据
 */
func (t *TaskCron) WheatTimeCron() {
	var (
		logKey  = "cron-room:"
		rdb     = coreRedis.GetChatroomRedis()
		cursor  uint64
		keys    []string
		pattern = redisKey.RoomWheatTimeCacheKey("*")
		ctx     = context.Background()
	)

	// 使用 SCAN 查找所有符合模式的键
	for {
		var err error
		var scannedKeys []string
		scannedKeys, cursor, err = rdb.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			coreLog.Warn(logKey, err)
			return
		}

		keys = append(keys, scannedKeys...)

		// 如果 cursor 为 0，表示遍历完成
		if cursor == 0 {
			break
		}
	}

	// 输出获取到的所有键
	coreLog.Info(logKey, "执行的所有Key:", keys)

	var doForFunc = func(rkey string) {
		data, err := rdb.HGetAll(ctx, rkey).Result()
		if err != nil {
			coreLog.Warn(logKey, fmt.Sprintf("获取%s:时出错:", rkey), err)
			return
		}
		var (
			by         []byte
			insertData = new(model.DoRoomWheatTime)
		)
		by, err = json.Marshal(data)
		if err != nil {
			coreLog.Warn(logKey, fmt.Sprintf("%s:Marshal出错:", rkey), err)
			return
		}
		err = json.Unmarshal(by, insertData)
		if err != nil {
			coreLog.Warn(logKey, fmt.Sprintf("%s:Unmarshal出错:", rkey), err)
			return
		}
		onwheatTime, _ := insertData.OnwheatTime.(string)
		parseTime, _ := time.Parse(time.DateTime, onwheatTime)
		endOfDay := time.Date(
			parseTime.Year(), parseTime.Month(), parseTime.Day(),
			23, 59, 59, 0, parseTime.Location(),
		)
		insertData.UpwheatTime = endOfDay.Format(time.DateTime)
		insertData.OnTime = endOfDay.Sub(parseTime).Seconds()
		tx := coreDb.GetMasterDb().Begin()
		err = tx.Model(model.RoomWheatTime{}).Create(insertData).Error
		if err != nil {
			coreLog.Warn(logKey, fmt.Sprintf("%s:写入mysql:", rkey), err)
			tx.Rollback()
			return
		}
		roomId, _ := insertData.RoomID.(string)
		rdb.Del(context.Background(), redisKey.RoomWheatTimeJoinUser(roomId))
		rdb.Del(context.Background(), redisKey.RoomWheatTimeGiftUser(roomId))
		//重置数据存储
		// 获取第二天的 00:00:00 时间
		nextDayStart := endOfDay.Add(1 * time.Second)
		newData := &model.DoRoomWheatTime{
			RoomID:          insertData.RoomID,
			UserID:          insertData.UserID,
			GuildID:         insertData.GuildID,
			RoomType:        insertData.RoomType,
			OnwheatTime:     nextDayStart.Format(time.DateTime),
			StatDate:        nextDayStart.Format(time.DateOnly),
			EnterCount:      0,
			EnterTimes:      0,
			RewardCount:     0,
			RewardTimes:     0,
			RewardUserCount: 0,
		}
		_, err = rdb.HMSet(ctx, redisKey.RoomWheatTimeCacheKey(roomId), newData).Result()
		if err != nil {
			coreLog.Error(logKey, fmt.Sprintf("%s:重置redis失败:", rkey), err)
			tx.Rollback()
		}
		tx.Commit()
	}
	// 遍历所有键并获取它们的值
	for _, key := range keys {
		go doForFunc(key)
	}
	coreLog.Info(logKey, "执行完毕")
}
