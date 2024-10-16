package service_room

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"strconv"
	"yfapi/core/coreRedis"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	service_goods "yfapi/internal/service/goods"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	"yfapi/typedef/response/room"

	"github.com/spf13/cast"
)

// 根据房间ID和麦位下标获取指定麦位信息
func GetSeatPositionByRoomIdAndKey(roomId string, seatId int) *room.RoomWheatPosition {
	ctx := context.Background()
	res := getMicPositionsByKey(ctx, roomId, cast.ToString(seatId))
	if res != nil && res.UserInfo.UserId != "" {
		sg := service_goods.UserGoods{}
		sgres, _ := sg.GetGoodsByKeys(res.UserInfo.UserId, true, enum.GoodsTypeSL, enum.GoodsTypeMWK)
		if len(sgres) > 0 {
			for _, v := range sgres {
				if v.GoodsTypeKey == enum.GoodsTypeSL {
					res.UserInfo.Voice = *v
				}
				if v.GoodsTypeKey == enum.GoodsTypeMWK {
					res.UserInfo.Frame = *v
				}
			}
		}
	}
	return res
}

// 获取房间的用户在麦麦位信息，带排序，正序
func GetSeatPositionsToUsersByRoomId(roomId string) []*room.RoomWheatPosition {
	//查询房间的麦位信息
	ctx := context.Background()
	seatList, _ := getMicPositions(ctx, roomId)
	seats := make([]*room.RoomWheatPosition, 0)
	if len(seatList) > 0 {
		for _, sv := range seatList {
			//有用户获取麦位信息
			if sv.Identity != enum.GuestMicSeat {
				seats = append(seats, &sv)
			}
		}
	}
	sort.Slice(seats, func(i, j int) bool {
		return seats[i].Id < seats[j].Id
	})
	return seats
}

// 获取房间的麦位信息 --接口使用
func GetSeatPositionsByRoomId(roomId string) []*room.RoomWheatPosition {
	//查询房间的麦位信息
	ctx := context.Background()
	seatList, _ := getMicPositions(ctx, roomId)
	seats := make([]*room.RoomWheatPosition, len(seatList))
	if len(seatList) > 0 {
		sg := service_goods.UserGoods{}
		for _, sv := range seatList {
			if sv.UserInfo.UserId != "" {
				sgres, _ := sg.GetGoodsByKeys(sv.UserInfo.UserId, true, enum.GoodsTypeSL, enum.GoodsTypeMWK)
				if len(sgres) > 0 {
					for _, v := range sgres {
						if v.GoodsTypeKey == enum.GoodsTypeSL {
							sv.UserInfo.Voice = *v
						}
						if v.GoodsTypeKey == enum.GoodsTypeMWK {
							sv.UserInfo.Frame = *v
						}
					}
				}
			}

			seats[sv.Id] = &sv
		}
	}
	return seats
}

func SetMicPositions(roomId string, nums, isboss, roomType int, mics ...room.RoomWheatPosition) (map[string]room.RoomWheatPosition, error) {
	micPositions := make(map[string]room.RoomWheatPosition)
	var err error
	log.Printf(">>>>%#v", mics)
	//处理主持麦位的增加
	micNum := len(mics)
	if micNum > 0 {
		micPositions["0"] = mics[0]
	}
	for i := 0 + micNum; i < nums; i++ {
		var (
			key            int
			identity, keys string
		)
		key = i
		if i == 0 {
			identity = enum.CompereMicSeat //主持人
		} else {
			switch roomType {
			case enum.RoomTypeEmoMan, enum.RoomTypeEmoWoman:
				identity = enum.CounselorMicSeat //咨询师麦
			case enum.RoomTypeFriend, enum.RoomTypeDating:
				identity = enum.NormalMicSeat //普通麦位
			case enum.RoomTypeSing:
				identity = enum.MusicianMicSeat //音乐人麦
			default:
				identity = enum.NormalMicSeat //普通麦
			}
			if i == 1 && isboss == 1 {
				identity = enum.GuestMicSeat //嘉宾
			}
		}
		keys = strconv.Itoa(key)
		micPositions[keys] = room.RoomWheatPosition{
			Id:       key,
			Identity: identity,
			Status:   1,
			Mute:     false,
		}
	}
	log.Printf("<><><><><><>%#v", micPositions)
	pipe := coreRedis.GetChatroomRedis().Pipeline()
	ctx := context.Background()
	for _, micPosition := range micPositions {
		data, err := json.Marshal(micPosition)
		if err != nil {
			return nil, err
		}
		pipe.HSet(ctx, redisKey.RoomWheatPosition(roomId), micPosition.Id, data)
	}
	_, err = pipe.Exec(ctx)
	return micPositions, err
}

// storeMicPositions 将麦位信息存储到 Redis 哈希中
func storeMicPositions(roomId string) (map[string]room.RoomWheatPosition, error) {
	// micPositions := make(map[string]room.RoomWheatPosition)
	roomDao := &dao.RoomDao{}
	res, err := roomDao.GetRoomPositions(roomId)
	if err != nil {
		return nil, err
	}
	nums := res.Num
	return SetMicPositions(roomId, nums, res.IsBoss, res.RoomType)
	// for i := 0; i < nums; i++ {
	// 	var (
	// 		key            int
	// 		identity, keys string
	// 	)
	// 	key = i
	// 	if i == 0 {
	// 		identity = enum.CompereMicSeat //主持人
	// 	} else {
	// 		switch res.RoomType {
	// 		case enum.RoomTypeEmoMan, enum.RoomTypeEmoWoman:
	// 			identity = enum.CounselorMicSeat //咨询师麦
	// 		case enum.RoomTypeFriend, enum.RoomTypeDating:
	// 			identity = enum.NormalMicSeat //普通麦位
	// 		case enum.RoomTypeSing:
	// 			identity = enum.MusicianMicSeat //音乐人麦
	// 		default:
	// 			identity = enum.NormalMicSeat //普通麦
	// 		}
	// 		if i == 1 && res.IsBoss == 1 {
	// 			identity = enum.GuestMicSeat //嘉宾
	// 		}
	// 	}
	// 	keys = strconv.Itoa(key)
	// 	micPositions[keys] = room.RoomWheatPosition{
	// 		Id:       key,
	// 		Identity: identity,
	// 		Status:   1,
	// 	}
	// }
	// pipe := coreRedis.GetChatroomRedis().Pipeline()
	// ctx := context.Background()
	// for _, micPosition := range micPositions {
	// 	data, err := json.Marshal(micPosition)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	pipe.HSet(ctx, redisKey.RoomWheatPosition(roomId), micPosition.Id, data)
	// }
	// _, err = pipe.Exec(ctx)
	// return micPositions, err
}

// getMicPositions 从 Redis 哈希中获取所有麦位信息
func getMicPositions(ctx context.Context, roomId string) (map[string]room.RoomWheatPosition, error) {
	micPositions := make(map[string]room.RoomWheatPosition)
	fields, err := coreRedis.GetChatroomRedis().HGetAll(ctx, redisKey.RoomWheatPosition(roomId)).Result()
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		log.Println("调取全部麦位，初始化房间麦位")
		micPositions, err = storeMicPositions(roomId)
		if err != nil {
			return nil, err
		}
	} else {
		for id, data := range fields {
			var micPosition = new(room.RoomWheatPosition)
			err := json.Unmarshal([]byte(data), micPosition)
			if err != nil {
				return nil, err
			}
			micPositions[id] = *micPosition
		}
	}

	return micPositions, nil
}

// getMicPositionsByKey 从 Redis 哈希中获取指定麦位信息
func getMicPositionsByKey(ctx context.Context, roomId string, keys string) *room.RoomWheatPosition {
	// log.Println("获取hash存储类型:", coreRedis.GetChatroomRedis().ObjectEncoding(ctx, redisKey.RoomWheatPosition(roomId)))
	micPosition := &room.RoomWheatPosition{}
	field := coreRedis.GetChatroomRedis().HGet(ctx, redisKey.RoomWheatPosition(roomId), keys).Val()
	if len(field) == 0 {
		//如果查询不到，进行初始化麦位
		micPositions, err := storeMicPositions(roomId)
		if err != nil {
			return nil
		}
		micPosition.Id = micPositions[keys].Id
		micPosition.Identity = micPositions[keys].Identity
		micPosition.UserInfo = micPositions[keys].UserInfo
		micPosition.Mute = micPositions[keys].Mute
		micPosition.Status = micPositions[keys].Status
		return micPosition
	}

	err := json.Unmarshal([]byte(field), micPosition)
	if err != nil {
		return nil
	}
	if len(micPosition.Identity) > 0 {
		return micPosition
	}
	return nil
}

// UpdateMicPosition 更新 Redis 哈希中的某个麦位信息
func UpdateMicPosition(ctx context.Context, roomId string, micPosition room.RoomWheatPosition) error {
	data, err := json.Marshal(micPosition)
	if err != nil {
		return err
	}
	_, err = coreRedis.GetChatroomRedis().HSet(ctx, roomId, micPosition.Id, data).Result()
	return err
}

// GetRoomUserMicPositionMap 查询房间在麦玩家
func GetRoomUserMicPositionMap(roomId string) map[string]*room.RoomWheatPosition {
	micPositions := make(map[string]*room.RoomWheatPosition)
	fields := GetSeatPositionsByRoomId(roomId)
	for _, data := range fields {
		if data.Status != enum.MicStatusUsed {
			continue
		}
		micPositions[data.UserInfo.UserId] = data
	}
	return micPositions
}

// 用户上麦存储用户信息和麦位关联
func RoomWheatPositionAddUserInfo(seatInfo *room.RoomWheatPosition, userInfo *model.User, roomInfo *model.Room) *room.RoomWheatPosition {
	userRoomWheat := room.RoomWheatUserInfo{
		UserId:     userInfo.Id,
		UserNo:     cast.ToInt(userInfo.UserNo),
		Uid32:      cast.ToInt32(userInfo.OriUserNo),
		UserName:   userInfo.Nickname,
		UserAvatar: userInfo.Avatar, //头像提前带前缀
	}
	if roomInfo.LiveType == enum.LiveTypeChatroom {
		userRoomWheat.CharmCount = cast.ToInt(coreRedis.GetChatroomRedis().HGet(context.Background(), redisKey.ChatroomUserCharmKey(roomInfo.Id), userInfo.Id).Val())
	} else {
		userRoomWheat.CharmCount = cast.ToInt(coreRedis.GetChatroomRedis().HGet(context.Background(), redisKey.AnchorRoomUserCharmKey(roomInfo.Id), userInfo.Id).Val())
	}
	sg := service_goods.UserGoods{}
	sgres, _ := sg.GetGoodsByKeys(userInfo.Id, true, enum.GoodsTypeSL, enum.GoodsTypeMWK)
	if len(sgres) > 0 {
		for _, v := range sgres {
			if v.GoodsTypeKey == enum.GoodsTypeSL {
				userRoomWheat.Voice = *v
			}
			if v.GoodsTypeKey == enum.GoodsTypeMWK {
				userRoomWheat.Frame = *v
			}
		}
	}
	seatInfo.UserInfo = userRoomWheat
	seatInfo.Status = enum.MicStatusUsed
	seatInfo.Mute = true //麦位静音
	return seatInfo
}
