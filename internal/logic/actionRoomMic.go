package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/internal/service/acl"
	"yfapi/internal/service/auth"
	service_im "yfapi/internal/service/im"
	service_room "yfapi/internal/service/room"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_room "yfapi/typedef/request/room"
	"yfapi/typedef/response/room"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// NormalMicCommand 普通麦位操作
func (a *ActionRoom) NormalMicCommand(c *gin.Context, command, roomId, userId, targetUserId string, seat int) {
	seatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomId, seat)
	if seatInfo == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if seatInfo.Identity != enum.NormalMicSeat { // 非普通麦
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	userInfo := service_user.GetUserBaseInfo(userId)
	switch command {
	case acl.UpNormalMic: // 自己上麦
		a.upSeat(c, roomId, seatInfo, userInfo)
	case acl.DownNormalMic: // 下麦
		a.downSeat(c, roomId, userId, seatInfo)
	case acl.HoldUserUpNormalMic: // 抱用户上麦
		targetUserInfo := service_user.GetUserBaseInfo(targetUserId)
		if len(targetUserInfo.Id) == 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		// 当前用户是否在房间内
		isUserInRoom(roomId, targetUserId)
		// 上麦
		a.upSeat(c, roomId, seatInfo, targetUserInfo, true)
	case acl.SwitchNormalMic: // 关闭麦位
		a.closeSeat(roomId, seatInfo)
	case acl.MuteNormalMic: // 静音麦位
		a.muteSeat(userId, roomId, seatInfo)
	case acl.HoldUserDownNormalMic: // 抱用户下麦
		a.downSeat(c, roomId, userId, seatInfo, true)
	case acl.ApplyNormalMic: // 申请上麦
		a.applyUpSeat(roomId, userId, seat)
	}
}

// CompereMicCommand 主持麦位操作
func (a *ActionRoom) CompereMicCommand(c *gin.Context, command, roomId, userId, targetUserId string) {
	seat := 0
	seatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomId, seat)
	if seatInfo == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	userInfo := service_user.GetUserBaseInfo(userId)
	if len(userInfo.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	switch command {
	case acl.UpCompereMic: // 自己上麦
		a.upSeat(c, roomId, seatInfo, userInfo, true)
	case acl.DownCompereMic: // 下麦
		a.downSeat(c, roomId, userId, seatInfo)
	case acl.HoldCompereUpCompereMic: // 抱主持人上麦
		// 是否为主持人的判断
		if !new(auth.Auth).IsHasRoomRole(roomId, targetUserId, 1005) {
			panic(error2.I18nError{
				Code: error2.ErrCodeUserNotCompere,
				Msg:  nil,
			})
		}
		targetUserInfo := service_user.GetUserBaseInfo(targetUserId)
		if len(targetUserInfo.Id) == 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		// 当前用户是否在房间内
		isUserInRoom(roomId, targetUserId)
		// 上麦
		a.upSeat(c, roomId, seatInfo, targetUserInfo, true)
	case acl.SwitchCompereMic: // 关闭麦位
		a.closeSeat(roomId, seatInfo)
	case acl.MuteCompereMic: // 静音麦位
		a.muteSeat(userId, roomId, seatInfo)
	case acl.HoldCompereDownCompereMic: // 抱用户下麦
		a.downSeat(c, roomId, userId, seatInfo, true)
	}
}

// GuestMicCommand 嘉宾麦位操作
func (a *ActionRoom) GuestMicCommand(c *gin.Context, command, roomId, userId, targetUserId string) {
	seat := 1
	seatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomId, seat)
	if seatInfo == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if seatInfo.Identity != enum.GuestMicSeat { // 不是嘉宾麦
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	userInfo := service_user.GetUserBaseInfo(userId)
	switch command {
	case acl.UpGuestMic: // 自己上麦
		a.upSeat(c, roomId, seatInfo, userInfo)
	case acl.DownGuestMic: // 下麦
		a.downSeat(c, roomId, userId, seatInfo)
	case acl.HoldUserUpGuestMic: // 抱嘉宾上麦
		targetUserInfo := service_user.GetUserBaseInfo(targetUserId)
		if len(targetUserInfo.Id) == 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		// 当前用户是否在房间内
		isUserInRoom(roomId, targetUserId)
		// 上麦
		a.upSeat(c, roomId, seatInfo, targetUserInfo, true)
	case acl.SwitchGuestMic: // 关闭麦位
		a.closeSeat(roomId, seatInfo)
	case acl.MuteGuestMic: // 静音麦位
		a.muteSeat(userId, roomId, seatInfo)
	case acl.HoldUserDownGuestMic: // 抱用户下麦
		a.downSeat(c, roomId, userId, seatInfo, true)
	case acl.ApplyGuestMic: // 申请上麦
		a.applyUpSeat(roomId, userId, seat)
	}
}

// CounselorMicCommand 咨询师麦位操作
func (a *ActionRoom) CounselorMicCommand(c *gin.Context, command, roomId, userId, targetUserId string, seat int) {
	seatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomId, seat)
	if seatInfo == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if seatInfo.Identity != enum.CounselorMicSeat { // 非咨询师麦
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	userInfo := service_user.GetUserBaseInfo(userId)
	switch command {
	case acl.UpCounselorMic: // 自己上麦
		a.upSeat(c, roomId, seatInfo, userInfo)
	case acl.DownCounselorMic: // 下麦
		a.downSeat(c, roomId, userId, seatInfo)
	case acl.HoldUserUpCounselorMic: // 抱用户上麦
		targetUserInfo := service_user.GetUserBaseInfo(targetUserId)
		if len(targetUserInfo.Id) == 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		// 当前用户是否在房间内
		isUserInRoom(roomId, targetUserId)
		// 上麦
		a.upSeat(c, roomId, seatInfo, targetUserInfo, true)
	case acl.SwitchCounselorMic: // 关闭麦位
		a.closeSeat(roomId, seatInfo)
	case acl.MuteCounselorMic: // 静音麦位
		a.muteSeat(userId, roomId, seatInfo)
	case acl.HoldUserDownCounselorMic: // 抱用户下麦
		a.downSeat(c, roomId, userId, seatInfo, true)
	case acl.ApplyCounselorMic: // 申请上麦
		a.applyUpSeat(roomId, userId, seat)
	}
}

// MusicianMicCommand 音乐人麦位操作
func (a *ActionRoom) MusicianMicCommand(c *gin.Context, command, roomId, userId, targetUserId string, seat int) {
	seatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomId, seat)
	if seatInfo == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if seatInfo.Identity != enum.MusicianMicSeat { // 非音乐人麦
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	userInfo := service_user.GetUserBaseInfo(userId)
	switch command {
	case acl.UpMusicianMic: // 自己上麦
		a.upSeat(c, roomId, seatInfo, userInfo)
	case acl.DownMusicianMic: // 下麦
		a.downSeat(c, roomId, userId, seatInfo)
	case acl.HoldUserUpMusicianMic: // 抱用户上麦
		targetUserInfo := service_user.GetUserBaseInfo(targetUserId)
		if len(targetUserInfo.Id) == 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		// 当前用户是否在房间内
		isUserInRoom(roomId, targetUserId)
		// 上麦
		a.upSeat(c, roomId, seatInfo, targetUserInfo, true)
	case acl.SwitchMusicianMic: // 关闭麦位
		a.closeSeat(roomId, seatInfo)
	case acl.MuteMusicianMic: // 静音麦位
		a.muteSeat(userId, roomId, seatInfo)
	case acl.HoldUserDownMusicianMic: // 抱用户下麦
		a.downSeat(c, roomId, userId, seatInfo, true)
	case acl.ApplyMusicianMic: // 申请上麦
		a.applyUpSeat(roomId, userId, seat)
	}
}

// 上麦
func (a *ActionRoom) upSeat(c *gin.Context, roomId string, seatInfo *room.RoomWheatPosition, userInfo *model.User, isHold ...bool) {
	if new(acl.RoomAcl).IsOnHiddenMic(roomId, userInfo.Id) {
		panic(error2.I18nError{
			Code: error2.ErrCodeBeforeDownHiddenMic,
		})
	}
	roomInfo, _ := new(dao.RoomDao).GetRoomById(roomId)
	if seatInfo.Status == enum.MicStatusClose {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomSeatClosed,
			Msg:  nil,
		})
	}
	if seatInfo.Status == enum.MicStatusUsed {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomSeatUsed,
			Msg:  nil,
		})
	}
	msgCode := enum.UP_SEAT_MSG
	if len(isHold) == 0 { // 自己上麦

	} else { // 被抱上麦或主持位自己上麦或主持同意上麦
		// 清除玩家之前的麦序申请
		if seatInfo.Identity != enum.CompereMicSeat { // 非主持自己上麦
			removeUpSeatApplyList(roomId, userInfo.Id, true)
			msgCode = enum.HOLD_UP_SEAT_MSG
		}
	}
	currSeatInfo := getUserInSeat(roomId, userInfo.Id)
	if currSeatInfo != nil { // 在座自动离座
		a.downSeat(c, roomId, userInfo.Id, currSeatInfo)
	}
	seatInfo = service_room.RoomWheatPositionAddUserInfo(seatInfo, userInfo, &roomInfo)

	err := coreRedis.GetChatroomRedis().HSet(context.Background(), redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo)).Err()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	//如果上主持麦，表示开播 记录直播统计数据
	if seatInfo.Identity == enum.CompereMicSeat {
		go service_room.InitStoreRoomWHeatTime(userInfo.Id, &roomInfo)
	}
	// 推送上麦消息
	new(service_im.ImPublicService).SendCustomMsg(roomId, seatInfo, msgCode)
	// 推送公屏上麦动作消息
	go noticeActionUpSeat(c, roomId, userInfo.Nickname, seatInfo)
}

// 自由上麦
func (a *ActionRoom) freeUpSeat(c *gin.Context, roomId string, seatInfo *room.RoomWheatPosition, userInfo *model.User, isHold ...bool) {
	if new(acl.RoomAcl).IsOnHiddenMic(roomId, userInfo.Id) {
		panic(error2.I18nError{
			Code: error2.ErrCodeBeforeDownHiddenMic,
		})
	}
	if seatInfo.Status == enum.MicStatusClose {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomSeatClosed,
			Msg:  nil,
		})
	}
	if seatInfo.Status == enum.MicStatusUsed {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomSeatUsed,
			Msg:  nil,
		})
	}
	msgCode := enum.UP_SEAT_MSG
	if len(isHold) == 0 { // 自己上麦
		// 主持麦信息
		compereSeat := service_room.GetSeatPositionByRoomIdAndKey(roomId, 0)
		if compereSeat == nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
	} else { // 被抱上麦或主持位自己上麦或主持同意上麦
		// 清除玩家之前的麦序申请
		if seatInfo.Identity != enum.CompereMicSeat { // 非主持自己上麦
			removeUpSeatApplyList(roomId, userInfo.Id, true)
			msgCode = enum.HOLD_UP_SEAT_MSG
		}
	}
	currSeatInfo := getUserInSeat(roomId, userInfo.Id)
	if currSeatInfo != nil { // 在座自动离座
		a.downSeat(c, roomId, userInfo.Id, currSeatInfo)
	}
	roomInfo, _ := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	seatInfo = service_room.RoomWheatPositionAddUserInfo(seatInfo, userInfo, roomInfo)

	err := coreRedis.GetChatroomRedis().HSet(context.Background(), redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo)).Err()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	//如果上主持麦，表示开播 记录直播统计数据
	if seatInfo.Identity == enum.CompereMicSeat {
		go service_room.InitStoreRoomWHeatTime(userInfo.Id, roomInfo)
	}
	// 推送上麦消息
	new(service_im.ImPublicService).SendCustomMsg(roomId, seatInfo, msgCode)
	// 推送公屏上麦动作消息
	go noticeActionUpSeat(c, roomId, userInfo.Nickname, seatInfo)
}

// 下麦
func (a *ActionRoom) downSeat(c *gin.Context, roomId, userId string, seatInfo *room.RoomWheatPosition, isHold ...bool) {
	if seatInfo.Status != enum.MicStatusUsed {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomMicErr,
			Msg:  nil,
		})
	}
	msgCode := enum.HOLD_DOWN_SEAT_MSG
	if len(isHold) == 0 { // 是否被抱
		if seatInfo.UserInfo.UserId != userId {
			panic(error2.I18nError{
				Code: error2.ErrCodeRoomPermissionDenied,
				Msg:  nil,
			})
		}
		msgCode = enum.DOWN_SEAT_MSG
	}
	nickname := seatInfo.UserInfo.UserName
	seatInfo.UserInfo = room.RoomWheatUserInfo{}
	seatInfo.Status = enum.MicStatusNormal
	err := coreRedis.GetChatroomRedis().HSet(context.Background(), redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo)).Err()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	//如果下主持麦，表示关播 记录关播直播统计数据到mysql
	if seatInfo.Identity == enum.CompereMicSeat {
		go service_room.StoreRoomWHeatTimeToMysql(roomId)
	}
	// 推送下麦消息
	new(service_im.ImPublicService).SendCustomMsg(roomId, seatInfo, msgCode)
	// 推送公屏麦位动作消息
	go noticeActionUpSeat(c, roomId, nickname, seatInfo)
}

// 麦位开始倒计时
func (a *ActionRoom) seatTimerOpen(roomId string, seat int, seconds int64) {
	seatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomId, seat)
	if seatInfo == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	if seatInfo.Identity != enum.GuestMicSeat && seatInfo.Identity != enum.NormalMicSeat && seatInfo.Identity != enum.CounselorMicSeat {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 推送开启倒计时
	data := struct {
		Id      int   `json:"id"`      // 座位ID
		Seconds int64 `json:"seconds"` // 倒计时秒数
	}{
		Id:      seat,
		Seconds: seconds,
	}
	new(service_im.ImPublicService).SendCustomMsg(roomId, data, enum.TIMER_OPEN_MSG)
}

// 重置魅力值
func (a *ActionRoom) ResetCharm(roomId string) {
	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	pipe := coreRedis.GetChatroomRedis().Pipeline()
	ctx := context.Background()
	if roomInfo.LiveType == enum.LiveTypeChatroom {
		pipe.Del(ctx, redisKey.ChatroomUserCharmKey(roomId))
	} else {
		pipe.Del(ctx, redisKey.AnchorRoomUserCharmKey(roomId))
	}
	userSeat := service_room.GetRoomUserMicPositionMap(roomId)
	for _, seatInfo := range userSeat {
		seatInfo.UserInfo.CharmCount = 0
		pipe.HSet(ctx, redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo))
	}
	_, _ = pipe.Exec(ctx)
	// 推送魅力值重置通知
	new(service_im.ImPublicService).SendCustomMsg(roomInfo.Id, nil, enum.RESET_CHARM_MSG)
}

// 隐藏房间
func (a *ActionRoom) hiddenRoom(roomId string) {
	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	isHidden := false
	if roomInfo.HiddenStatus == 2 { // 关闭状态变为开启
		roomInfo.HiddenStatus = 1
		isHidden = true
	} else { // 开启状态变为关闭
		roomInfo.HiddenStatus = 2
	}
	roomInfo.UpdateTime = time.Now()
	err = roomDao.Save(*roomInfo)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 推送房间隐藏状态
	new(service_im.ImPublicService).SendCustomMsg(roomId, isHidden, enum.HIDDEN_ROOM_MSG)
}

// 自由上下麦
func (a *ActionRoom) freedSeat(roomId string) {
	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	isFreedMic := false
	if roomInfo.FreedMicStatus == 2 { // 关闭状态变为开启
		roomInfo.FreedMicStatus = 1
		isFreedMic = true
	} else { // 开启状态变为关闭
		roomInfo.FreedMicStatus = 2
	}
	roomInfo.UpdateTime = time.Now()
	err = roomDao.Save(*roomInfo)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 推送房间自由上下麦状态
	new(service_im.ImPublicService).SendCustomMsg(roomId, isFreedMic, enum.FREED_MIC_MSG)
}

// 自由发言
func (a *ActionRoom) freedSpeak(roomId string) {
	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	isFreedSpeak := false
	if roomInfo.FreedSpeakStatus == 2 { // 关闭状态变为开启
		roomInfo.FreedSpeakStatus = 1
		isFreedSpeak = true
	} else { // 开启状态变为关闭
		roomInfo.FreedSpeakStatus = 2
	}
	roomInfo.UpdateTime = time.Now()
	err = roomDao.Save(*roomInfo)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	if roomInfo.FreedSpeakStatus == 2 { //关闭自由发言
		seatsInfo, _ := coreRedis.GetChatroomRedis().HGetAll(context.Background(), redisKey.RoomWheatPosition(roomId)).Result()
		if len(seatsInfo) > 0 {
			for _, seatInfo := range seatsInfo {
				seatInfoStruct := room.RoomWheatPosition{}
				err = json.Unmarshal([]byte(seatInfo), &seatInfoStruct)
				if err == nil {
					fmt.Println(seatInfoStruct.Id, seatInfoStruct)
					seatInfoStruct.Mute = true
					coreRedis.GetChatroomRedis().HSet(context.Background(), redisKey.RoomWheatPosition(roomId), seatInfoStruct.Id, easy.JSONStringFormObject(seatInfoStruct))
				}
			}
		}
	}
	// 推送房间自由发言状态
	new(service_im.ImPublicService).SendCustomMsg(roomId, isFreedSpeak, enum.FREED_SPEAK_MSG)
}

// 关闭公屏
func (a *ActionRoom) closePublicChat(roomId string) {
	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	isOpen := false
	if roomInfo.PublicScreenStatus == 2 { // 关闭状态变为开启
		roomInfo.PublicScreenStatus = 1
		isOpen = true
	} else { // 开启状态变为关闭
		roomInfo.PublicScreenStatus = 2
	}
	roomInfo.UpdateTime = time.Now()
	err = roomDao.Save(*roomInfo)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 推送房间公屏状态
	new(service_im.ImPublicService).SendCustomMsg(roomId, isOpen, enum.CLOSE_PUBLIC_CHAT_MSG)
}

// 清空公屏
func (a *ActionRoom) clearPublicChat(roomId string) {
	// 推送房间清空公屏通知
	new(service_im.ImPublicService).SendCustomMsg(roomId, nil, enum.CLEAR_PUBLIC_CHAT_MSG)
}

// 关闭麦位
func (a *ActionRoom) closeSeat(roomId string, seatInfo *room.RoomWheatPosition) {
	if seatInfo.Status == enum.MicStatusUsed {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomSeatUsed,
			Msg:  nil,
		})
	}
	if seatInfo.Status == enum.MicStatusClose { // 开启麦位
		seatInfo.Status = enum.MicStatusNormal
	} else if seatInfo.Status == enum.MicStatusNormal { // 关闭麦位
		seatInfo.Status = enum.MicStatusClose
	}
	err := coreRedis.GetChatroomRedis().HSet(context.Background(), redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo)).Err()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	new(service_im.ImPublicService).SendCustomMsg(roomId, seatInfo, enum.CLOSE_SEAT_MSG)
}

// 静音麦位
func (a *ActionRoom) muteSeat(userId, roomId string, seatInfo *room.RoomWheatPosition) {
	seatInfo.Mute = !seatInfo.Mute
	if !seatInfo.Mute { //禁言判断
		mute := new(UserMute).IsMute(seatInfo.UserInfo.UserId, roomId)
		if mute.ID > 0 {
			if userId == seatInfo.UserInfo.UserId {
				panic(error2.I18nError{
					Code: error2.ErrCodeUserMuteMsg,
					Msg:  map[string]any{"minute": int(mute.EndTime.Sub(mute.StartTime).Minutes())},
				})
			} else {
				panic(error2.I18nError{
					Code: error2.ErrCodeHeMuteMsg,
					Msg:  map[string]any{"minute": int(mute.EndTime.Sub(mute.StartTime).Minutes())},
				})
			}
		}
	}
	freedSpeakStatus := new(acl.RoomAcl).RoomFreedSpeakStatus(roomId)
	if freedSpeakStatus == enum.SwitchOff { //自由发言关闭状态
		roleIds, _, _ := new(auth.Auth).GetRoleListByRoomIdAndUserId(roomId, userId)
		hasAuth := false
		isCompere := new(acl.RoomAcl).IsOnCompereMicSeat(userId, roomId)
		for _, roleId := range roleIds {
			if easy.InArray(roleId, []int{enum.PresidentRoleId, enum.HouseOwnerRoleId}) {
				hasAuth = true
			}
			if roleId == enum.CompereRoleId && isCompere {
				hasAuth = true
			}
		}
		if !hasAuth {
			panic(error2.I18nError{
				Code: error2.ErrCodeFreedSpeakClosed,
			})
		}
	}
	err := coreRedis.GetChatroomRedis().HSet(context.Background(), redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo)).Err()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 推送麦位静音消息
	new(service_im.ImPublicService).SendCustomMsg(roomId, seatInfo, enum.MUTE_SEAT_MSG)
}

// 清空全麦
func (a *ActionRoom) clearAllSeat(c *gin.Context, roomId, userId string) {
	userSeatMap := service_room.GetRoomUserMicPositionMap(roomId)
	for _, seatInfo := range userSeatMap {
		if seatInfo.Identity == enum.CompereMicSeat || seatInfo.Identity == enum.GuestMicSeat {
			continue
		}
		a.downSeat(c, roomId, userId, seatInfo, true)
	}
}

// 踢出全麦
func (a *ActionRoom) kickAllSeat(c *gin.Context, roomId, userId string) {
	userSeatMap := service_room.GetRoomUserMicPositionMap(roomId)
	for _, seatInfo := range userSeatMap {
		a.downSeat(c, roomId, userId, seatInfo, true)
	}
}

// 申请上麦
func (a *ActionRoom) applyUpSeat(roomId, userId string, seat int) {
	// 玩家是否在麦
	if _, isExist := service_room.GetRoomUserMicPositionMap(roomId)[userId]; isExist {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomSeatUp,
			Msg:  nil,
		})
	}
	// 清除玩家之前的麦序申请
	removeUpSeatApplyList(roomId, userId)
	redisCli := coreRedis.GetChatroomRedis().Pipeline()
	ctx := context.Background()
	// 将玩家添加到麦序尾部
	redisCli.RPush(ctx, redisKey.RoomUpSeatApplyList(roomId), userId)
	// 保存玩家申请上麦的座位号
	redisCli.HSet(ctx, redisKey.RoomUpSeatApplyInfo(roomId), userId, seat)
	_, _ = redisCli.Exec(ctx)
	// 推送申请上麦列表变动通知
	new(service_im.ImPublicService).SendCustomMsg(roomId, getUpSeatApplyListCount(roomId), enum.APPLY_SEAT_MSG)
}

// 取消上麦申请
func (a *ActionRoom) cancelUpSeatApply(roomId, userId string) {
	// 清除玩家之前的麦序申请
	removeUpSeatApplyList(roomId, userId, true)
}

// 清空申请上麦列表
func (a *ActionRoom) clearUpSeatApply(roomId, userId string) {
	// 麦序操作权限检查
	upSeatApplyCheck(roomId, userId)
	coreRedis.GetChatroomRedis().Del(context.Background(), redisKey.RoomUpSeatApplyList(roomId))
	// 推送申请上麦列表变动通知
	new(service_im.ImPublicService).SendCustomMsg(roomId, 0, enum.APPLY_SEAT_MSG)
}

// 拒绝上麦申请
func (a *ActionRoom) refuseUpSeatApply(c *gin.Context, roomId, userId, targetUserId string) {
	// 麦序操作权限检查
	upSeatApplyCheck(roomId, userId)
	// 移除上麦申请
	removeUpSeatApplyList(roomId, targetUserId, true)
	// 发送上麦申请结果通知
	new(service_im.ImCommonService).Send(c, userId, []string{targetUserId}, "", enum.MsgCustom, false, enum.APPLY_SEAT_RESULT_MSG)
}

// 同意上麦申请
func (a *ActionRoom) acceptUpSeatApply(c *gin.Context, roomId, userId, targetUserId string, seat int) {
	// 麦序操作权限检查
	upSeatApplyCheck(roomId, userId)
	// 目标玩家信息
	targetUserInfo := service_user.GetUserBaseInfo(targetUserId)
	if len(targetUserInfo.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 目标玩家是否在房间内
	if !service_room.IsUserInRoom(roomId, targetUserId) {
		// 清除上麦申请
		removeUpSeatApplyList(roomId, targetUserId, true)
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomUserNotFound,
			Msg:  nil,
		})
	}
	// 目标玩家是否在麦
	isSeat := isUserInSeat(roomId, targetUserId)
	if isSeat {
		// 清除上麦申请
		removeUpSeatApplyList(roomId, targetUserId, true)
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomSeatUp,
			Msg:  nil,
		})
	}
	// 玩家的上麦申请信息
	applySeatId := getUserUpSeatApplyInfo(roomId, targetUserId)
	if applySeatId == 0 {
		// 清除上麦申请
		removeUpSeatApplyList(roomId, targetUserId, true)
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 玩家申请的麦位类型
	applySeatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomId, applySeatId)
	if applySeatInfo == nil || applySeatInfo.Identity == enum.CompereMicSeat {
		// 清除上麦申请
		removeUpSeatApplyList(roomId, targetUserId, true)
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 当前房间的麦位列表
	var currSeatInfo *room.RoomWheatPosition
	seatList := service_room.GetSeatPositionsByRoomId(roomId)
	for _, seatInfo := range seatList {
		if seatInfo.Identity == applySeatInfo.Identity && seatInfo.Status == enum.MicStatusNormal {
			currSeatInfo = seatInfo
			break
		}
	}
	if currSeatInfo == nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomSeatFilled,
			Msg:  nil,
		})
	}
	// 上麦
	a.upSeat(c, roomId, currSeatInfo, targetUserInfo, true)
	// 发送上麦申请结果通知
	new(service_im.ImCommonService).Send(c, userId, []string{targetUserId}, "", enum.MsgCustom, true, enum.APPLY_SEAT_RESULT_MSG)
}

// 当前玩家是否在麦
func isUserInSeat(roomId, userId string) bool {
	if _, isExist := service_room.GetRoomUserMicPositionMap(roomId)[userId]; isExist {
		return true
	}
	return false
}

// 当前玩家在麦信息
func getUserInSeat(roomId, userId string) *room.RoomWheatPosition {
	if seatInfo, isExist := service_room.GetRoomUserMicPositionMap(roomId)[userId]; isExist {
		return seatInfo
	}
	return nil
}

// 当前用户是否在房间内
func isUserInRoom(roomId, userId string) {
	// 当前用户是否在房间内
	if !service_room.IsUserInRoom(roomId, userId) {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomUserNotFound,
			Msg:  nil,
		})
	}
}

// 当前座位是否空闲
func isIdleSeat(roomId string, seat int) *room.RoomWheatPosition {
	// 当前麦位信息
	currSeat := service_room.GetSeatPositionByRoomIdAndKey(roomId, seat)
	if currSeat == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if currSeat.Status != enum.MicStatusNormal { // 非空闲座位
		if currSeat.Status == enum.MicStatusClose {
			panic(error2.I18nError{
				Code: error2.ErrCodeRoomSeatUsed,
				Msg:  nil,
			})
		} else if currSeat.Status == enum.MicStatusUsed {
			panic(error2.I18nError{
				Code: error2.ErrCodeRoomSeatUsed,
				Msg:  nil,
			})
		} else {
			panic(error2.I18nError{
				Code: error2.ErrCodeRoomMicErr,
				Msg:  nil,
			})
		}
	}
	return currSeat
}

// 上麦申请操作权限检查
func upSeatApplyCheck(roomId, userId string) {
	// 超管、巡查、会长、房主、管理员、在麦主持
	checkRoleIdList := []int{enum.SuperAdminRoleId, enum.PatrolRoleId, enum.PresidentRoleId, enum.HouseOwnerRoleId, enum.RoomAdminRoleId}
	isHave := new(auth.Auth).IsHaveCurrRole(roomId, userId, checkRoleIdList)
	if isHave {
		return
	}
	// 当前用户是否为在麦主持人
	isRoomCompere(roomId, userId)
}

// 当前用户是否为在麦主持人
func isRoomCompere(roomId, userId string) {
	// 主持麦信息
	compereSeat := service_room.GetSeatPositionByRoomIdAndKey(roomId, 0)
	if compereSeat == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if compereSeat.UserInfo.UserId != userId {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}
}

// 清除玩家麦序申请
func removeUpSeatApplyList(roomId, userId string, notice ...bool) {
	redisCli := coreRedis.GetChatroomRedis().Pipeline()
	ctx := context.Background()
	cmd := redisCli.LRem(ctx, redisKey.RoomUpSeatApplyList(roomId), 0, userId)
	redisCli.HDel(ctx, redisKey.RoomUpSeatApplyInfo(roomId), userId)
	_, _ = redisCli.Exec(ctx)

	if len(notice) > 0 && notice[0] && cmd.Val() > 0 {
		// 推送申请上麦列表变动通知
		new(service_im.ImPublicService).SendCustomMsg(roomId, getUpSeatApplyListCount(roomId), enum.APPLY_SEAT_MSG)
	}
}

// 查询当前房间的麦位申请人数
func getUpSeatApplyListCount(roomId string) int64 {
	return coreRedis.GetChatroomRedis().LLen(context.Background(), redisKey.RoomUpSeatApplyList(roomId)).Val()
}

func getUserUpSeatApplyInfo(roomId, userId string) int {
	seatId := coreRedis.GetChatroomRedis().HGet(context.Background(), redisKey.RoomUpSeatApplyInfo(roomId), userId).Val()
	return cast.ToInt(seatId)
}

// 公屏上下麦消息通知
func noticeActionUpSeat(c *gin.Context, roomId, nickname string, seatInfo *room.RoomWheatPosition) {
	isUp := seatInfo.Status == enum.MicStatusUsed
	// 查询当前房间的麦位模板信息
	config, _ := new(dao.RoomDao).GetRoomPositions(roomId)
	defaultAdd := 0
	if config.RoomType > 0 && config.IsBoss == 1 {
		defaultAdd = 1
	}

	seatName := ""
	switch seatInfo.Id {
	case 0:
		seatName = i18n_msg.GetI18nMsg(c, i18n_msg.CompereMicMsgKey)
	default:
		if seatInfo.Id == 1 && seatInfo.Identity == enum.GuestMicSeat {
			seatName = i18n_msg.GetI18nMsg(c, i18n_msg.GuestMicMsgKey)
		} else {
			//seatName = fmt.Sprintf("%v号麦", seatInfo.Id+1)
			seatName = i18n_msg.GetI18nMsg(c, i18n_msg.MicSeatMsgKey, map[string]any{"num": seatInfo.Id - defaultAdd})
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

// GetRoomUpSeatApplyList 查询房间的上麦申请列表
func (a *ActionRoom) GetRoomUpSeatApplyList(c *gin.Context, req *request_room.UpSeatApplyListReq) (res room.UpSeatApplyListRes) {
	// 查询房间的所有麦位申请列表
	userIdList := coreRedis.GetChatroomRedis().LRange(c, redisKey.RoomUpSeatApplyList(req.RoomId), 0, 200).Val()
	userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
	for _, userId := range userIdList {
		res.List = append(res.List, room.UpSeatApplyInfo{
			UserId:   userInfoMap[userId].Id,
			UserNo:   userInfoMap[userId].UserNo,
			Uid32:    cast.ToInt32(userInfoMap[userId].OriUserNo),
			Nickname: userInfoMap[userId].Nickname,
			Avatar:   userInfoMap[userId].Avatar,
			Sex:      userInfoMap[userId].Sex,
		})
	}
	res.Count = len(res.List)
	return
}

// 隐藏麦
func (a *ActionRoom) HiddenMic(c *gin.Context, roomId, userId, command string) {
	key := redisKey.RoomHiddenMicKey(roomId)
	id, _ := coreRedis.GetChatroomRedis().Get(context.Background(), key).Result()
	aclService := new(acl.RoomAcl)
	switch command {
	case acl.UpHiddenMic:
		inRoom := aclService.IsInRoom(id, roomId, "")
		if len(id) > 0 && inRoom {
			hiddenMicUserInfo := service_user.GetUserBaseInfo(id)
			panic(error2.I18nError{
				Code: error2.ErrCodeHiddenMicHasUser,
				Msg:  map[string]interface{}{"nickname": hiddenMicUserInfo.Nickname},
			})
		} else {
			seatInfo := aclService.GetUserMicInfo(roomId, userId)
			if seatInfo != nil {
				panic(error2.I18nError{
					Code: error2.ErrCodeOnMicSeat,
				})
				//a.downSeat(c, roomId, userId, seatInfo)
			}
			coreRedis.GetChatroomRedis().Set(context.Background(), key, userId, time.Hour)
		}
	case acl.DownHiddenMic:
		if id == userId {
			coreRedis.GetChatroomRedis().Del(context.Background(), key, userId)
		}
	}
	return
}
