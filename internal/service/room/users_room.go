package service_room

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/acl"
	service_im "yfapi/internal/service/im"
	"yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	"yfapi/typedef/response/room"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// 房间在线用户列表类 包含，加入房间，退出房间，查看是否在此房间，查询房间的所有成员功能
type RoomUsersOnlie struct {
	RoomId string
}

// 加入 房间在线用户列表
func (r *RoomUsersOnlie) AddUserToRoom(c *gin.Context, userID string) error {
	redisCli := coreRedis.GetChatroomRedis()
	//读取天为单位的房间用户列表的值
	score := redisCli.ZScore(c, redisKey.RoomUsersDayList(r.RoomId), userID).Val()
	if score == 0 {
		nowTime := time.Now().UnixMilli()
		fractionalPart := 1 - (float64(nowTime) / 1e13)
		// 计算总分数
		score += fractionalPart
	}

	err := redisCli.ZAdd(c, redisKey.RoomUsersOnlineList(r.RoomId), redis.Z{
		Score:  score,
		Member: userID,
	}).Err()
	if err != nil {
		return err
	}
	//设置用户所在房间
	err = redisCli.Set(context.Background(), redisKey.UserInWhichRoom(userID, helper.GetClientType(c)), r.RoomId, time.Second*300).Err()
	if err != nil {
		return err
	}
	return nil
}

// 增加繁荣值
func (r *RoomUsersOnlie) EditUserToRoom(c *gin.Context, userId string, score float64) error {
	err := coreRedis.GetChatroomRedis().ZIncrBy(c, redisKey.RoomUsersOnlineList(r.RoomId), score, userId).Err()
	if err != nil {
		return err
	}
	return nil
}

// 退出房间
func (r *RoomUsersOnlie) RemoveUserToRoom(c *gin.Context, userID, client string, roomInfo *model.Room) error {
	err := coreRedis.GetChatroomRedis().ZRem(c, redisKey.RoomUsersOnlineList(r.RoomId), userID).Err()
	if err != nil {
		return err
	}
	if roomInfo == nil {
		roomInfo, _ = new(dao.RoomDao).FindOne(&model.Room{Id: r.RoomId})
	}
	//删除用户所在房间信息
	if len(client) == 0 {
		for _, clientStr := range enum.ClientTypeArray {
			inRoomId, _ := coreRedis.GetChatroomRedis().Get(c, redisKey.UserInWhichRoom(userID, clientStr)).Result()
			if inRoomId == roomInfo.Id {
				coreRedis.GetChatroomRedis().Del(c, redisKey.UserInWhichRoom(userID, clientStr))
			}
		}
	} else {
		coreRedis.GetChatroomRedis().Del(c, redisKey.UserInWhichRoom(userID, client))
	}
	//如果用户在麦位则清除麦位信息
	r.DownMic(c, r.RoomId, userID)
	// 直播间下播清除魅力值
	if roomInfo != nil && roomInfo.LiveType == enum.LiveTypeAnchor && userID == roomInfo.UserId {
		pipe := coreRedis.GetChatroomRedis().Pipeline()
		ctx := context.Background()
		roomId := roomInfo.Id
		pipe.Del(ctx, redisKey.AnchorRoomUserCharmKey(roomId))
		userSeat := GetRoomUserMicPositionMap(roomId)
		for _, seatInfo := range userSeat {
			seatInfo.UserInfo.CharmCount = 0
			pipe.HSet(ctx, redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo))
		}
		_, _ = pipe.Exec(ctx)
		// 推送魅力值重置通知
		new(service_im.ImPublicService).SendCustomMsg(roomInfo.Id, nil, enum.RESET_CHARM_MSG)
	}
	return nil
}

// 判断多个用户是否在此房间
func (r *RoomUsersOnlie) IsUserOnline(c *gin.Context, users []string) map[string]bool {
	// 创建一个管道
	pipe := coreRedis.GetChatroomRedis().Pipeline()
	key := redisKey.RoomUsersOnlineList(r.RoomId)
	// 添加 ZScore 命令到管道
	var cmds []*redis.FloatCmd
	for _, member := range users {
		cmd := pipe.ZScore(c, key, member)
		cmds = append(cmds, cmd)
	}

	// 执行管道
	_, err := pipe.Exec(c)
	if err != nil && err != redis.Nil {
		return nil
	}

	// 处理命令结果
	result := make(map[string]bool)
	for i, cmd := range cmds {
		err := cmd.Err()
		if err == redis.Nil {
			result[users[i]] = false
		} else if err != nil {
			result[users[i]] = false
		} else {
			result[users[i]] = true
		}
	}
	return result
}

// 获取房间的指定数量的在线用户和贡献值
func (r *RoomUsersOnlie) GetOnlineUsersIdCard(c *gin.Context, num int64) ([]redis.Z, error) {
	//检查房间用户是否真实在线，如果不在线则进行清理
	go func() {
		success, unlock, err := coreRedis.ChatroomLock(c, redisKey.ClearRoomNotOnlineUserLock(r.RoomId), time.Second*20)
		if err != nil || !success {
			coreLog.Error("ClearRoomNotOnlineUserLock err :%+v", err)
			return
		}
		defer unlock()
		var cursor uint64
		for {
			keys, cursor, err := coreRedis.GetChatroomRedis().ZScan(c, redisKey.RoomUsersOnlineList(r.RoomId), cursor, "", 100).Result()
			if err != nil {
				coreLog.Error("GetOnlineUsersIdCard err :%+v", err)
				return
			}
			if len(keys) > 0 {
				deleteUserId := []string{}
				for offset, userId := range keys {
					if offset%2 == 1 {
						continue
					}
					//inRoomId := coreRedis.GetChatroomRedis().Get(c, redisKey.UserInWhichRoom(userId)).Val()
					isInRoom := new(acl.RoomAcl).IsInRoom(userId, r.RoomId, "")
					if !isInRoom && userId != helper.GetUserId(c) {
						deleteUserId = append(deleteUserId, userId)
					}
				}
				if len(deleteUserId) > 0 {
					coreLog.Info("需要移除的不在房间用户 %+v", deleteUserId)
					coreRedis.GetChatroomRedis().ZRem(c, redisKey.RoomUsersOnlineList(r.RoomId), deleteUserId)
				}
			}
			if cursor == 0 || len(keys) == 0 {
				return
			}
		}

	}()
	users, err := coreRedis.GetChatroomRedis().ZRevRangeWithScores(c, redisKey.RoomUsersOnlineList(r.RoomId), 0, num-1).Result()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 获取贡献值大于0的用户,并且根据score倒序排序
func (r *RoomUsersOnlie) GetHightGradeUsers(c *gin.Context) []redis.Z {
	redisCli := coreRedis.GetChatroomRedis()
	keys := redisKey.RoomUsersOnlineList(r.RoomId)
	res, err := redisCli.ZRevRangeByScoreWithScores(c, keys, &redis.ZRangeBy{
		Min: "1",
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil
	}
	return res
}

func IsUserInRoom(roomId, userId string) bool {
	_, err := coreRedis.GetChatroomRedis().ZScore(context.Background(), redisKey.RoomUsersOnlineList(roomId), userId).Result()
	return err == nil
}

// 获取房间的指定数量的用户的id集合
func (r *RoomUsersOnlie) GetOnlineUsersMembers(c *gin.Context, num int64) ([]string, error) {
	users, err := coreRedis.GetChatroomRedis().ZRange(c, redisKey.RoomOnlineUsersCache(r.RoomId), 0, num-1).Result()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 清除用户再其他房间信息
func (r *RoomUsersOnlie) ClearOtherRoom(c *gin.Context, userId string) {
	roomId, err := coreRedis.GetChatroomRedis().Get(context.Background(), redisKey.UserInWhichRoom(userId, helper.GetClientType(c))).Result()
	if err != nil {
		coreLog.LogError("用户退出其他房间错误%+v", err)
		return
	}
	if roomId != r.RoomId {
		err := r.RemoveUserToRoom(c, userId, helper.GetClientType(c), nil)
		if err != nil {
			coreLog.LogError("用户退出其他房间错误%+v", err)
		}
	}
}

// 把用户从房间下麦
func (a *RoomUsersOnlie) DownMic(c *gin.Context, roomId, userId string) {
	result, err := coreRedis.GetChatroomRedis().HGetAll(context.Background(), redisKey.RoomWheatPosition(roomId)).Result()
	if err != nil {
		coreLog.Error("DownMic err:%+v", err)
		return
	}
	for _, v := range result {
		seatInfo := room.RoomWheatPosition{}
		err = json.Unmarshal([]byte(v), &seatInfo)
		if err != nil {
			continue
		}
		if seatInfo.UserInfo.UserId == userId {
			nickname := seatInfo.UserInfo.UserName
			seatInfo.UserInfo = room.RoomWheatUserInfo{}
			seatInfo.Status = enum.MicStatusNormal
			err = coreRedis.GetChatroomRedis().HSet(context.Background(), redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo)).Err()
			if err != nil {
				coreLog.Error("DownMic err:%+v", err)
			}
			//如果是主持麦位：停止直播间统计数据
			if seatInfo.Identity == enum.CompereMicSeat {
				go StoreRoomWHeatTimeToMysql(roomId)
			}
			new(service_im.ImPublicService).SendCustomMsg(roomId, seatInfo, enum.DOWN_SEAT_MSG)
			// 推送公屏麦位动作消息
			isUp := seatInfo.Status == enum.MicStatusUsed
			seatName := ""
			switch seatInfo.Id {
			case 0:
				seatName = i18n_msg.GetI18nMsg(c, i18n_msg.CompereMicMsgKey)
			default:
				if seatInfo.Id == 1 && seatInfo.Identity == enum.GuestMicSeat {
					seatName = i18n_msg.GetI18nMsg(c, i18n_msg.GuestMicMsgKey)
				} else {
					//seatName = fmt.Sprintf("%v号麦", seatInfo.Id+1)
					seatName = i18n_msg.GetI18nMsg(c, i18n_msg.MicSeatMsgKey, map[string]any{"num": seatInfo.Id + 1})
				}
			}
			action := ""
			if isUp {
				action = i18n_msg.GetI18nMsg(c, i18n_msg.UpMsgKey)
			} else {
				action = i18n_msg.GetI18nMsg(c, i18n_msg.DownMsgKey)
			}
			noticeMsg := fmt.Sprintf("%v %v %v", nickname, action, seatName)
			new(service_im.ImPublicService).SendActionMsg(c, map[string]string{
				"content": noticeMsg,
			}, "", "", roomId, "", enum.SEAT_ACTION_MSG)
		}
	}
}

// 加入，退出，变成贡献值，刷新榜单前三用户信息
func (a *RoomUsersOnlie) OnlineChangeToThree(roomId string) {
	ctx := context.Background()
	redisCli := coreRedis.GetChatroomRedis()
	three, err := redisCli.ZRevRangeWithScores(ctx, redisKey.RoomUsersOnlineList(roomId), 0, 2).Result()
	if err != nil {
		return
	}
	usersInfo := make([]interface{}, 0)
	for _, v := range three {
		newId, _ := v.Member.(string)
		userInfo := user.GetUserBaseInfo(newId)
		item := map[string]interface{}{
			"userId":   newId,
			"avatar":   userInfo.Avatar,
			"sex":      userInfo.Sex,
			"nickname": userInfo.Nickname,
		}
		usersInfo = append(usersInfo, item)
	}
	new(service_im.ImPublicService).SendCustomMsg(roomId, usersInfo, enum.ONLINE_CHANGE_THREE_MSG)
}

// 获取指定用户的信息，和相邻的两个用户信息的贡献值
func (a *RoomUsersOnlie) GetMemberWithNeighbors(roomId string, userId string) map[string]float64 {
	// 获取成员在有序集合中的排名
	var (
		ctx      = context.Background()
		redisCli = coreRedis.GetChatroomRedis()
		keys     = redisKey.RoomUsersOnlineList(roomId)
	)
	var res map[string]float64 = map[string]float64{
		"owner": math.Trunc(redisCli.ZScore(ctx, keys, userId).Val()),
		"first": 0,
		"end":   0,
	}
	rank, err := redisCli.ZRank(ctx, keys, userId).Result()
	if err != nil {
		return res
	}

	// 获取相邻两个数据的范围
	start := int(rank) - 1
	end := int(rank) + 1
	// 处理边界条件
	if start < 0 {
		start = 0
	}
	info := redisCli.ZRangeWithScores(ctx, keys, int64(start), int64(end)).Val()
	if len(info) == 3 {
		res["first"] = math.Trunc(info[2].Score)
		res["owner"] = math.Trunc(info[1].Score)
		res["end"] = math.Trunc(info[0].Score)
	} else if len(info) == 2 {
		res["owner"] = math.Trunc(info[1].Score)
		res["end"] = math.Trunc(info[0].Score)
	} else {
		res["owner"] = math.Trunc(info[0].Score)
	}
	return res
}
