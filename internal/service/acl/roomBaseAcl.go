package acl

import (
	"context"
	"encoding/json"
	"strconv"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/service/auth"
	service_goods "yfapi/internal/service/goods"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	"yfapi/typedef/response/room"
	"yfapi/util/easy"

	"github.com/spf13/cast"
)

// 根据房间ID用户ID获取拥有角色
func (r *RoomAcl) GetRolesByRoomIdAndUserId(roomId, userId string) ([]int, error) {
	result, _, err := new(auth.Auth).GetRoleListByRoomIdAndUserId(roomId, userId)
	return result, err
}

// 检测用户是否拥有此权限
func (r *RoomAcl) CheckUserRule(userId, roomId, ruleName string) (bool, error) {
	switch ruleName { // 自定义命令不做权限校验
	case ClearUpSeatApply, RefuseUpSeatApply, AcceptUpSeatApply, CancelUpSeatApply:
		return true, nil
	}

	rules, err := r.getVerifyAuth(userId, roomId, ruleName)
	if err != nil {
		return false, err
	}
	return rules[ruleName], err
}

// 获取当前麦位用户信息
func (r *RoomAcl) MicUserInfo(roomId string, seat int) room.RoomWheatPosition {
	micInfo := r.GetMicInfoBySeat(roomId, seat)
	return *micInfo
}

// 判断用户是否在主持麦
func (r *RoomAcl) IsOnCompereMicSeat(userId, roomId string) bool {
	micInfo := r.GetMicInfoBySeatName(roomId, enum.CompereMicSeat)
	if micInfo.UserInfo.UserId == userId {
		return true
	}
	return false
}

// 判断房间隐藏状态
func (r *RoomAcl) RoomHiddenStatus(roomId string) int {
	roomInfo, err := new(dao.RoomDao).GetRoomById(roomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	return roomInfo.HiddenStatus
}

// 判断房间锁定状态
func (r *RoomAcl) RoomLockStatus(roomId string) int {
	roomInfo, err := new(dao.RoomDao).GetRoomById(roomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	if len(roomInfo.RoomPwd) > 0 {
		return enum.SwitchOpen
	}
	return enum.SwitchOff
}

// 判断房间自由上下麦状态
func (r *RoomAcl) RoomFreedMicStatus(roomId string) int {
	roomInfo, err := new(dao.RoomDao).GetRoomById(roomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	return roomInfo.FreedMicStatus
}

// 判断房间自由发言状态
func (r *RoomAcl) RoomFreedSpeakStatus(roomId string) int {
	roomInfo, err := new(dao.RoomDao).GetRoomById(roomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	return roomInfo.FreedSpeakStatus
}

// 判断房间关闭声音状态
func (r *RoomAcl) RoomMuteStatus(roomId, userId string) int {
	result, err := coreRedis.GetChatroomRedis().Get(context.Background(), redisKey.UserRoomSwitchStatus(userId, roomId)).Result()
	if err != nil {
		coreLog.Error("RoomMuteStatus err:%+v", err)
		return enum.SwitchOff
	}
	if len(result) == 0 {
		return enum.SwitchOff
	}
	roomUserStatus := room.UserRoomSwitchStatus{}
	err = json.Unmarshal([]byte(result), &roomUserStatus)
	if err != nil {
		coreLog.Error("RoomMuteStatus err:%+v", err)
		return enum.SwitchOff
	}
	return roomUserStatus.RoomMute
}

// 判断动效开启状态
func (r *RoomAcl) RoomSpecialEffectsStatus(roomId, userId string) int {
	result, err := coreRedis.GetChatroomRedis().Get(context.Background(), redisKey.UserRoomSwitchStatus(userId, roomId)).Result()
	if err != nil {
		coreLog.Error("RoomSpecialEffectsStatus err:%+v", err)
		return enum.SwitchOff
	}
	if len(result) == 0 {
		return enum.SwitchOff
	}
	roomUserStatus := room.UserRoomSwitchStatus{}
	err = json.Unmarshal([]byte(result), &roomUserStatus)
	if err != nil {
		coreLog.Error("RoomSpecialEffectsStatus err:%+v", err)
		return enum.SwitchOff
	}
	return roomUserStatus.RoomSpecialEffects
}

// 判断公屏状态
func (r *RoomAcl) RoomPublicChatStatus(roomId string) int {
	roomInfo, err := new(dao.RoomDao).GetRoomById(roomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	return roomInfo.PublicScreenStatus
}

// 判断麦位开启关闭状态
func (r *RoomAcl) MicSwitchStatus(roomId string, seat int, micSeat string) int {
	micInfo := r.GetMicInfoBySeatName(roomId, micSeat, seat)
	if micInfo.Status == 3 {
		return enum.SwitchOpen
	}
	return enum.SwitchOff
}

// 判断麦位静音状态
func (r *RoomAcl) MicMuteStatus(roomId string, seat int, micSeat string) int {
	micInfo := r.GetMicInfoBySeatName(roomId, micSeat, seat)
	if micInfo.Mute {
		return enum.SwitchOpen
	}
	return enum.SwitchOff
}

// 根据座位号获取麦位信息
func (r *RoomAcl) GetMicInfoBySeat(roomId string, seat int) *room.RoomWheatPosition {
	micInfo := &room.RoomWheatPosition{}
	result, err := coreRedis.GetChatroomRedis().HGet(context.Background(), redisKey.RoomWheatPosition(roomId), cast.ToString(seat)).Result()
	if err != nil {
		coreLog.Error("GetMicInfoBySeat err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomMicErr,
			Msg:  nil,
		})
	}
	err = json.Unmarshal([]byte(result), micInfo)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomMicErr,
			Msg:  nil,
		})
	}
	return micInfo
}

// 根据座位名获取用户
func (r *RoomAcl) GetMicInfoBySeatName(roomId, seatName string, seat ...int) *room.RoomWheatPosition {
	micInfo := &room.RoomWheatPosition{}
	result, err := coreRedis.GetChatroomRedis().HGetAll(context.Background(), redisKey.RoomWheatPosition(roomId)).Result()
	if err != nil {
		coreLog.Error("GetMicInfoBySeatName err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomMicErr,
			Msg:  nil,
		})
	}
	for _, v := range result {
		err = json.Unmarshal([]byte(v), micInfo)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrCodeRoomMicErr,
				Msg:  nil,
			})
		}
		if len(seat) == 0 {
			if micInfo.Identity == seatName {
				return micInfo
			}
		} else {
			if micInfo.Identity == seatName && micInfo.Id == seat[0] {
				return micInfo
			}
		}
	}
	return nil
}

// 获取验证的权限
func (r *RoomAcl) getVerifyAuth(userId, roomId string, rules ...string) (map[string]bool, error) {
	isCompere := r.IsOnCompereMicSeat(userId, roomId)
	//authMap, err := coreAuth.Instance().VerifyAuth(cast.ToInt(userId), cast.ToInt(roomId), isCompere, rules...)
	authMap, err := new(auth.Auth).VerifyAuth(userId, roomId, isCompere, rules...)
	roles, _, _ := new(auth.Auth).GetRoleListByRoomIdAndUserId(roomId, userId)
	roleOk := easy.InArray(enum.CompereRoleId, roles)
	if !isCompere && roleOk {
		authMap[UpCompereMic] = true
	}
	return authMap, err
}

// GetUserMicInfo
//
//	@Description: 获取用户所在麦位信息
//	@receiver r
//	@param roomId
//	@return *room.RoomWheatPosition
func (r *RoomAcl) GetUserMicInfo(roomId, userId string) *room.RoomWheatPosition {
	// micInfo := &room.RoomWheatPosition{}
	wheats := GetRoomUserMicPositionMap(roomId)
	if _, ok := wheats[userId]; ok {
		return wheats[userId]
	}
	return nil
}

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

// getMicPositions 从 Redis 哈希中获取所有麦位信息
func getMicPositions(ctx context.Context, roomId string) (map[string]room.RoomWheatPosition, error) {
	micPositions := make(map[string]room.RoomWheatPosition)
	fields, err := coreRedis.GetChatroomRedis().HGetAll(ctx, redisKey.RoomWheatPosition(roomId)).Result()
	if err != nil {
		return nil, err
	}
	if len(fields) == 0 {
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
}

func SetMicPositions(roomId string, nums, isboss, roomType int, mics ...room.RoomWheatPosition) (map[string]room.RoomWheatPosition, error) {
	micPositions := make(map[string]room.RoomWheatPosition)
	var err error
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

// 判断用户是否在房间内
func (r *RoomAcl) IsInRoom(userId, roomId, clientType string) bool {
	if len(clientType) > 0 {
		key := redisKey.UserInWhichRoom(userId, clientType)
		result, _ := coreRedis.GetChatroomRedis().Get(context.Background(), key).Result()
		if result == roomId {
			return true
		}
	} else {
		for _, clientStr := range enum.ClientTypeArray {
			key := redisKey.UserInWhichRoom(userId, clientStr)
			result, _ := coreRedis.GetChatroomRedis().Get(context.Background(), key).Result()
			if result == roomId {
				return true
			}
		}
	}
	return false
}

// 比较身份大小
func (r *RoomAcl) CompareRole(userId, targetUserId, roomId string) bool {
	userRoleIds, _ := r.GetRolesByRoomIdAndUserId(roomId, userId)
	userSmallRoleId := enum.NormalRoleId
	for _, v := range userRoleIds {
		if v == enum.AnchorRoleId {
			v = enum.HouseOwnerRoleId
		}
		if v < userSmallRoleId {
			userSmallRoleId = v
		}
	}
	targetRoleIds, _ := r.GetRolesByRoomIdAndUserId(roomId, targetUserId)
	targetSmallRoleId := enum.NormalRoleId
	for _, v := range targetRoleIds {
		if v == enum.AnchorRoleId {
			v = enum.HouseOwnerRoleId
		}
		if v < targetSmallRoleId {
			targetSmallRoleId = v
		}
	}
	if userSmallRoleId < targetSmallRoleId {
		return true
	}
	return false
}

// 判断用户是否在隐藏麦
func (r *RoomAcl) IsOnHiddenMic(roomId, userId string) bool {
	key := redisKey.RoomHiddenMicKey(roomId)
	id, _ := coreRedis.GetChatroomRedis().Get(context.Background(), key).Result()
	if id == userId {
		return true
	}
	return false
}
