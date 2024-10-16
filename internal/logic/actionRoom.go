package logic

import (
	"context"
	"log"
	"math"
	"sort"
	"strconv"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/acl"
	"yfapi/internal/service/auth"
	service_goods "yfapi/internal/service/goods"
	service_im "yfapi/internal/service/im"
	"yfapi/internal/service/rankList"
	service_room "yfapi/internal/service/room"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	"yfapi/typedef/request/room"
	response_im "yfapi/typedef/response/im"
	response_room "yfapi/typedef/response/room"
	"yfapi/util/easy"

	"github.com/spf13/cast"

	ginI18n "github.com/gin-contrib/i18n"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ActionRoom struct {
}

func (a *ActionRoom) CheckRoom(c *gin.Context, userId, roomId string, isJoinRoom bool) (resp response_room.CheckRoomResponse, roomInfo *model.Room) {
	//房间密码，房间黑名单，用户是否被踢出过
	roomDao := dao.RoomDao{}
	roomInfo, _ = roomDao.FindOne(&model.Room{Id: roomId})
	if roomInfo.Status != 1 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	} else {
		if roomInfo.LiveType == enum.LiveTypeAnchor {
			seatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomInfo.Id, 0)
			if roomInfo.UserId != userId {
				if seatInfo.UserInfo.UserId == "" {
					panic(error2.I18nError{
						Code: error2.ErrCodeRoomAnchorErr,
						Msg:  nil,
					})
				}
			}
		}
	}
	if roomInfo.RoomPwd != "" {
		resp.IsPwd = true
	}
	blacklistDao := dao.UserBlackListDao{}
	count := blacklistDao.IsLog(&model.UserBlacklist{
		RoomID:      roomId,
		ToID:        userId,
		IsEffective: true,
	})
	if count {
		resp.IsBlacklist = true
		resp.Msg = ginI18n.MustGetMessage(
			c,
			&i18n.LocalizeConfig{
				MessageID: strconv.Itoa(error2.ErrCodeBlackErr),
				TemplateData: map[string]interface{}{
					"roomName": roomInfo.Name,
				},
			})

	}
	val := coreRedis.GetChatroomRedis().TTL(c, redisKey.RoomKickOutKey(userId, roomId)).Val()
	if val.Seconds() > 0 {
		resp.IsKickOut = true
		if resp.Msg == "" {
			resp.Msg = ginI18n.MustGetMessage(
				c,
				&i18n.LocalizeConfig{
					MessageID:    strconv.Itoa(error2.ErrCodeBlackOutTimesErr),
					TemplateData: map[string]interface{}{"times": int(math.Ceil(val.Minutes()))},
				})
		}
	}
	if isJoinRoom {
		service_room.MultiLeaveRoom(c, userId)
	} else {
		//进行多端房间进入检测，如果其他端在房间内则提示
		service_room.MultiJoinRoom(c, userId, &resp)
	}
	return
}

func (a *ActionRoom) CheckIsRoom(c *gin.Context, roomId string) (resp response_room.CheckIsRoomResponse) {
	ownerId := handle.GetUserId(c)
	redisCli := coreRedis.GetChatroomRedis()
	redisRoomId := redisCli.Get(c, redisKey.UserInWhichRoom(ownerId, helper.GetClientType(c))).Val()
	if redisRoomId == "" || redisRoomId == "0" {
		resp.IsInRoom = false
		return
	} else {
		resp.IsInRoom = redisRoomId == roomId
		return
	}
}

// 房间公共处理逻辑方法
func (a *ActionRoom) roomInfo(roomInfo *model.Room, userId, roomId string, resp *response_room.ChatroomDTO) {
	if roomInfo.RoomPwd != "" {
		resp.IsLocked = true
	}
	resp.RoomTypeDesc = enum.RoomType(resp.RoomType).String()
	resp.IsCollect = new(dao.DaoUserCollect).IsRoomCollect(userId, roomId)
	resp.IsPublicScreen = roomInfo.PublicScreenStatus == 1
	resp.IsFreedMic = roomInfo.FreedMicStatus == 1
	resp.IsFreedSpeak = roomInfo.FreedSpeakStatus == 1
	mute := new(UserMute).IsMute(userId, roomId) //是否被禁言
	if mute.ID > 0 {
		resp.UserMute = true
	}
	//判断个播房间是否在连麦
	if roomInfo.LiveType == enum.LiveTypeAnchor {
		if _, ok := enum.RoomTemplates[roomInfo.TemplateId]; ok {
			resp.IsRelateWheat = true
		}
	}

	autoWelcomeDao := dao.AutoWelcome{}
	resp.AutoWelcomeContent = autoWelcomeDao.FirstContent(userId)

	seats := service_room.GetSeatPositionsByRoomId(roomId)
	resp.SeatList = append(resp.SeatList, seats...)

	// 礼物类目列表
	giftCategoryList, _ := new(dao.GiftShowCategory).GetListByLiveType(roomInfo.LiveType)
	for _, info := range giftCategoryList {
		resp.GiftCategoryList = append(resp.GiftCategoryList, response_room.GiftShowCategory{
			CategoryId:   info.ID,
			CategoryName: info.CategoryName,
		})
	}
	// 公会名称
	if roomInfo.GuildId != "0" {
		guildInfo, _ := new(dao.GuildDao).FindById(&model.Guild{ID: roomInfo.GuildId})
		resp.GuildName = guildInfo.Name
	}
	// 房间申请上麦人数
	resp.UpSeatApplyCount = getUpSeatApplyListCount(roomId)
	// 当前房间的身份
	resp.RoleIdList, _, _ = new(auth.Auth).GetRoleListByRoomIdAndUserId(roomId, userId)
	if len(resp.RoleIdList) == 0 {
		resp.RoleIdList = []int{enum.NormalRoleId}
	}
	// 当前房间热度值
	_, resp.HotStr = service_room.GetRoomHot(roomId, roomInfo.LiveType)
	// 房主信息
	userInfo := service_user.GetUserBaseInfo(roomInfo.UserId)
	resp.OwnerUserInfo = response_room.UpSeatApplyInfo{
		UserId:   userInfo.Id,
		UserNo:   userInfo.UserNo,
		Uid32:    cast.ToInt32(userInfo.OriUserNo),
		Nickname: userInfo.Nickname,
		Avatar:   userInfo.Avatar,
		Sex:      userInfo.Sex,
	}
	onHiddenMicUserId, _ := coreRedis.GetChatroomRedis().Get(context.Background(), redisKey.RoomHiddenMicKey(roomId)).Result()
	if onHiddenMicUserId == userId {
		resp.IsOnHiddenMic = true
	}
	return
}

func (a *ActionRoom) JoinRoom(c *gin.Context, userId, clientType string, req *room.ActionRoomReq) (resp response_room.ChatroomDTO) {
	roomId := req.RoomId
	pwd := req.Pwd
	//TODO: 检查用户的角色， 如果是超管或者其他的需要特殊对待 ，返回房间信息和麦位信息
	res, roomInfo := a.CheckRoom(c, userId, roomId, true)
	if res.IsPwd {
		//房间密码加密使用md5
		if roomInfo.RoomPwd != easy.Md5(pwd, 16, true) {
			panic(error2.I18nError{
				Code: error2.ErrCodeRoomPwdErr,
				Msg:  nil,
			})
		}
	}
	if res.IsKickOut {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomKickOutErr,
			Msg:  nil,
		})
	}
	if res.IsBlacklist {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomBlacklistErr,
			Msg:  nil,
		})
	}
	//设置房间用户列表
	serviceRoom := &service_room.RoomUsersOnlie{
		RoomId: roomId,
	}
	serviceRoom.ClearOtherRoom(c, userId)
	//个播房直接上麦
	if userId == roomInfo.UserId && roomInfo.RoomType == enum.RoomTypeAnchorVoice {
		seatInfo := new(acl.RoomAcl).GetMicInfoBySeat(roomId, 0)
		userInfo := service_user.GetUserBaseInfo(userId)
		//如果是个播房直接上麦,
		//如果麦位有信息，不执行上麦,杜绝杀死app，检测断开没到时间，进不去房间问题
		if seatInfo.UserInfo.UserId != userId {
			a.upSeat(c, roomId, seatInfo, userInfo, true)
		}
		go new(Notice).LivePublishNotice(c, userId, roomId)
	}

	resp = roomInfo.ToChatroomDTO()
	a.roomInfo(roomInfo, userId, roomId, &resp)

	err := serviceRoom.AddUserToRoom(c, userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeJoinRoomErr,
			Msg:  nil,
		})
	}

	isFollow := false
	if len(req.FollowUserId) > 0 { // 跟随进房
		ctx := context.Background()
		followKey := redisKey.UserFollowJoinRoom(userId)
		// 跟随进房公屏信息 每个用户每天在每个房间触发一次
		isExist, _ := coreRedis.GetUserRedis().HExists(ctx, followKey, roomId).Result()
		if !isExist {
			isFollow = true
			redisCli := coreRedis.GetUserRedis().Pipeline()
			redisCli.HSet(ctx, followKey, roomId, time.Now().Unix())
			redisCli.Expire(ctx, followKey, 24*time.Hour)
			_, _ = redisCli.Exec(ctx)
		}
	}
	//逻辑走完后通知im加入房间
	msg := i18n_msg.GetI18nMsg(c, i18n_msg.JoinRoomMsgKey)
	if isFollow {
		followUser := service_user.GetUserBaseInfo(req.FollowUserId)
		msg = i18n_msg.GetI18nMsg(c, i18n_msg.FollowJoinRoomMsgKey, map[string]any{"nickname": followUser.Nickname})
		//msg = fmt.Sprintf("踩着%v的小尾巴进入房间~", followUser.Nickname)
	}
	imRes := response_im.JoinRoomImResponse{
		Content: msg,
	}
	ses, _ := service_goods.UserGoods{}.GetGoodsByKeys(userId, true, enum.GoodsTypeERS)
	if len(ses) > 0 {
		for _, v := range ses {
			if v.GoodsTypeKey == enum.GoodsTypeERS {
				imRes.JoinSE = *v
			}
		}
	} else { // 没有进场特效 根据等级查询横幅颜色
		userLv, _ := new(dao.UserLevelLvDao).GetUserLvLevel(userId)
		if userLv.ID > 0 {
			// 进场横幅颜色
			imRes.Color, imRes.StrokeColor = userLv.GetColor()
		}
	}
	//判断用户是否再隐藏麦
	onHiddenMicUserId, _ := coreRedis.GetChatroomRedis().Get(context.Background(), redisKey.RoomHiddenMicKey(roomId)).Result()
	if onHiddenMicUserId == userId {
		resp.IsOnHiddenMic = true
	}
	go func() {
		time.Sleep(time.Second)
		new(service_im.ImPublicService).SendActionMsg(c, imRes, userId, "", roomId, clientType, enum.JOIN_ROOM_MSG)
	}()
	//排行榜首次进入直播间
	go func() {
		rankList.Instance().Calculate(rankList.CalculateReq{
			FromUserId: userId,
			Types:      "firstJoinRoom",
			RoomId:     roomId,
		})
	}()
	//直播数据统计
	go service_room.DoRoomWheatTimeOperation(userId, roomId, 1, service_room.EnterCount, service_room.EnterTimes)
	//加入，退出，变成贡献值，刷新榜单前三用户信息
	go serviceRoom.OnlineChangeToThree(roomId)
	//进入房间，欢迎语逻辑
	go a.getOnlinePractitionAutoWelcome(c, userId, roomId)
	// 增加房间热度值
	go service_room.UpdateRoomHotByJoinRoom(roomId, userId, roomInfo.LiveType)
	return
}

func (a *ActionRoom) LeaveRoom(c *gin.Context, userId, roomId, clientType string) any {
	log.Println("离开房间：", userId, roomId, clientType)
	//TODO 进行离开房间的相关逻辑操作
	roomDao := dao.RoomDao{}
	roomInfo, _ := roomDao.FindOne(&model.Room{Id: roomId})
	if roomInfo.UserId == "" {
		panic(error2.I18nError{
			Code: error2.ErrCodeLeaveRoomErr,
			Msg:  nil,
		})
	}
	serviceRoom := &service_room.RoomUsersOnlie{
		RoomId: roomId,
	}
	err := serviceRoom.RemoveUserToRoom(c, userId, helper.GetClientType(c), roomInfo)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeLeaveRoomErr,
			Msg:  nil,
		})
	}
	//新增个播间的博主操作离开房间，如果处在连麦中，断开连麦
	if roomInfo.LiveType == enum.LiveTypeAnchor && userId == roomInfo.UserId {
		if _, ok := enum.RoomTemplates[roomInfo.TemplateId]; ok {
			//踢出所有普通麦
			a.clearAllSeat(c, roomId, userId)
			//退出连麦
			a.updateRoomTemplate(c, roomInfo, cast.ToString(enum.RoomTemplateOne), 1)
		}
	}
	//下隐藏麦
	a.HiddenMic(c, roomId, userId, acl.DownHiddenMic)
	//通知im离开房间，如果是个播房房主退出，前端需要处理接收到im信息的用户自动退出房间
	new(service_im.ImPublicService).SendActionMsg(c, map[string]string{
		"content": i18n_msg.GetI18nMsg(c, i18n_msg.LeaveRoomMsgKey),
	}, userId, "", roomId, clientType, enum.LEAVE_ROOM_MSG)
	//加入，退出，变成贡献值，刷新榜单前三用户信息
	go serviceRoom.OnlineChangeToThree(roomId)
	return nil
}

// 发送文本消息
func (a *ActionRoom) SendText(c *gin.Context, content, userId, roomId, clientType, extra string) any {
	msg := new(service_im.ImPublicService).SendTextMsg(c, content, userId, "", roomId, clientType, extra)
	return msg
}

// 发送图片
func (a *ActionRoom) SendImage(c *gin.Context, url, userId, roomId, clientType, extra string, width, height int) any {
	msg := new(service_im.ImPublicService).SendImgMsg(c, url, userId, "", roomId, width, height, clientType, extra)
	return msg
}

// 发送礼物
func (a *ActionRoom) SendGift(c *gin.Context, giftId, fromUserId, toUserId, roomId, clientType, extra string, width, height int) any {
	//new(im.ImPublicService).SendGiftMsg(giftId, fromUserId, toUserId, roomId, clientType)
	return nil
}

func (a *ActionRoom) ExecCommand(c *gin.Context, req *room.ExecCommandReq) (res response_room.ExecCommandRes) {
	userId := helper.GetUserId(c)
	aclInstance := &acl.RoomAcl{
		UserId:       userId,
		TargetUserId: req.TargetUserId,
		RoomId:       req.RoomId,
		Seat:         req.Seat,
	}
	isOk, err := aclInstance.CheckUserRule(userId, req.RoomId, req.Command)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	if !isOk { // 权限拒绝
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}
	switch req.Command {
	case acl.UpNormalMic, acl.HoldUserUpNormalMic, acl.SwitchNormalMic, acl.ApplyNormalMic, acl.MuteNormalMic, acl.DownNormalMic, acl.HoldUserDownNormalMic: //普通麦位
		a.NormalMicCommand(c, req.Command, req.RoomId, userId, req.TargetUserId, req.Seat)
	case acl.UpCompereMic, acl.DownCompereMic, acl.HoldCompereUpCompereMic, acl.SwitchCompereMic, acl.MuteCompereMic, acl.HoldCompereDownCompereMic: // 主持麦位
		a.CompereMicCommand(c, req.Command, req.RoomId, userId, req.TargetUserId)
	case acl.UpGuestMic, acl.DownGuestMic, acl.HoldUserUpGuestMic, acl.SwitchGuestMic, acl.MuteGuestMic, acl.HoldUserDownGuestMic, acl.ApplyGuestMic: // 嘉宾麦
		a.GuestMicCommand(c, req.Command, req.RoomId, userId, req.TargetUserId)
	case acl.UpCounselorMic, acl.DownCounselorMic, acl.HoldUserUpCounselorMic, acl.SwitchCounselorMic, acl.MuteCounselorMic, acl.HoldUserDownCounselorMic, acl.ApplyCounselorMic: // 咨询师麦
		a.CounselorMicCommand(c, req.Command, req.RoomId, userId, req.TargetUserId, req.Seat)
	case acl.ShutUp: // 禁言
		a.UserMute(c, userId, req.TargetUserId, req.RoomId, cast.ToInt(req.Content))
	case acl.AddRoomBlacklist: //拉黑
		a.DoBlackout(c, req.TargetUserId, req.RoomId, req.Command, req.Content)
	case acl.OutRoom: //踢出房间
		a.KickOut(c, &room.KickOutReq{
			RoomId: req.RoomId,
			UserId: req.TargetUserId,
			Times:  req.Content,
		})
	case acl.Greeting: //自动欢迎语
		a.AutoWelcome(c, req.Content)
	case acl.LockRoom: //锁定房间
		a.RoomLock(c, &room.RoomLockReq{RoomId: req.RoomId, Pwd: req.Content})
	case acl.ClearUpSeatApply: // 清空麦位申请列表
		a.clearUpSeatApply(req.RoomId, userId)
	case acl.AcceptUpSeatApply: // 同意上麦申请
		a.acceptUpSeatApply(c, req.RoomId, userId, req.TargetUserId, req.Seat)
	case acl.RefuseUpSeatApply: // 拒绝上麦申请
		a.refuseUpSeatApply(c, req.RoomId, userId, req.TargetUserId)
	case acl.CancelUpSeatApply: // 取消上麦申请
		a.cancelUpSeatApply(req.RoomId, userId)
	case acl.RoomOutAllMic: // 踢出全麦 全部麦
		a.kickAllSeat(c, req.RoomId, userId)
	case acl.RoomClearMic: // 清空全麦 陪陪麦
		a.clearAllSeat(c, req.RoomId, userId)
	case acl.FreedMic: // 自由上下麦
		a.freedSeat(req.RoomId)
	case acl.FreedSpeak: // 自由发言
		a.freedSpeak(req.RoomId)
	case acl.HiddenRoom: // 隐藏房间
		a.hiddenRoom(req.RoomId)
	case acl.RoomClosePublicChat: // 关闭公屏
		a.closePublicChat(req.RoomId)
	case acl.RoomClearPublicChat: // 清空公屏
		a.clearPublicChat(req.RoomId)
	case acl.ResetGlamour: // 重置魅力值
		a.ResetCharm(req.RoomId)
	case acl.RoomWarningMessage: // 警告房间
		a.roomWarningMessage(c, req.RoomId, userId, req.Content)
	case acl.TimerGuestMic, acl.TimerNormalMic, acl.TimerCounselorMic:
		a.seatTimerOpen(req.RoomId, req.Seat, cast.ToInt64(req.Content))
	case acl.UpHiddenMic, acl.DownHiddenMic: //上下隐藏麦
		a.HiddenMic(c, req.RoomId, userId, req.Command)
	case acl.RoomRelateWheat: //连麦
		a.doRoomRelateWheat(c, req.RoomId, userId, req.Content)
	default:
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	return
}

// 警告房间
func (a *ActionRoom) roomWarningMessage(c *gin.Context, roomId, userId, content string) {
	if len(content) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}

	var toUserIdList []string
	// 展示给在房间的房主，在麦主持，管理员
	roomInfo, _ := new(dao.RoomDao).GetRoomById(roomId)
	roomOwnerUid := ""
	if len(roomInfo.UserId) > 0 && service_room.IsUserInRoom(roomId, roomInfo.UserId) {
		toUserIdList = append(toUserIdList, roomInfo.UserId)
		roomOwnerUid = roomInfo.UserId
	}
	seatInfo := service_room.GetSeatPositionByRoomIdAndKey(roomId, 0)
	if seatInfo != nil && seatInfo.Status == enum.MicStatusUsed && roomOwnerUid != seatInfo.UserInfo.UserId {
		toUserIdList = append(toUserIdList, seatInfo.UserInfo.UserId)
	}
	if len(toUserIdList) > 0 {
		// 推送给房间特定人群警告信息
		new(service_im.ImCommonService).Send(c, userId, toUserIdList, roomId, enum.MsgCustom, content, enum.ROOM_WARNING_MSG)
	}
}

// 自由上麦
func (a *ActionRoom) FreeUpSeat(c *gin.Context, req *room.FreeUpSeatReq) {
	userId := helper.GetUserId(c)
	seatInfo := service_room.GetSeatPositionByRoomIdAndKey(req.RoomId, *req.Seat)
	userInfo := service_user.GetUserBaseInfo(userId)
	if *req.Seat == 0 {
		isOk, err := new(acl.RoomAcl).CheckUserRule(userId, req.RoomId, acl.UpCompereMic)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
		if !isOk { // 权限拒绝
			panic(error2.I18nError{
				Code: error2.ErrCodeRoomPermissionDenied,
				Msg:  nil,
			})
		}
		a.freeUpSeat(c, req.RoomId, seatInfo, userInfo, true)
	} else {
		a.freeUpSeat(c, req.RoomId, seatInfo, userInfo)
	}
	return
}

// MuteLocalSeat
//
//	@Description: 本地静音麦位
//	@receiver a
//	@param c
//	@param req
func (a *ActionRoom) MuteLocalSeat(c *gin.Context, req *room.MutLocalSeatReq) bool {
	userId := helper.GetUserId(c)
	if !helper.ReqRateLimit(redisKey.MuteLocalSeatReqRate(userId, req.RoomId), time.Second*1) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFrequent,
			Msg:  nil,
		})
	}
	roomDao := new(dao.RoomDao)
	_, err := roomDao.FindOne(&model.Room{Id: req.RoomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	micInfo := new(acl.RoomAcl).GetUserMicInfo(req.RoomId, userId)
	if micInfo == nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserNotOnSeat,
			Msg:  nil,
		})
	}
	a.muteSeat(userId, req.RoomId, micInfo)
	micInfo = new(acl.RoomAcl).GetUserMicInfo(req.RoomId, userId)
	return micInfo.Mute
}

// 查询在麦用户身份是从业者的用户，并且设置自动欢迎语的用户
func (a *ActionRoom) getOnlinePractitionAutoWelcome(c *gin.Context, userId, roomId string) {
	// 超管，巡查,不推送欢迎语
	authService := auth.Auth{}
	ok := authService.IsSuperAdminRole(roomId, userId)
	if ok {
		return
	}
	//判断当前用户的身份，如果是从业者，不需要欢迎语
	cerdDao := dao.DaoUserPractitionerCerd{UserId: userId}
	cerds, _ := cerdDao.Find()
	if len(cerds) > 0 {
		return
	}
	ctx := context.Background()
	res := coreRedis.GetChatroomRedis().Get(ctx, redisKey.RoomAutoWelcomeKey(roomId, userId)).Val()
	if res != "" {
		return
	}

	serRes := service_room.GetRoomUserMicPositionMap(roomId)
	if len(serRes) == 0 {
		return
	}
	var ids []string
	for _, v := range serRes {
		ids = append(ids, v.UserInfo.UserId)
	}
	autoWe := dao.AutoWelcome{}
	autoIds := autoWe.FindToUserIds(ids)
	if len(autoIds) == 0 {
		return
	}
	go func() {
		time.Sleep(time.Second)
		new(service_im.ImCommonService).Send(c, userId, autoIds, roomId, enum.MsgCustom, "", enum.AUTO_WELCOME_MSG)
	}()
	//正式一小时有效期
	times := 1 * time.Hour
	coreRedis.GetChatroomRedis().Set(ctx, redisKey.RoomAutoWelcomeKey(roomId, userId), autoIds, times)
}

// GetHoldUpSeatUserList 查询可抱用户上麦列表
func (a *ActionRoom) GetHoldUpSeatUserList(c *gin.Context, req *room.UpSeatUserListReq) (res response_room.HoldUpSeatUserListRes) {
	if req.SeatType == 1 {
		res.List = getHoldUpSeatCompereList(c, req.RoomId)
	} else {
		res.List = getHoldUpSeatAllList(c, req.RoomId)
	}
	res.Count = len(res.List)
	return
}

// 查询可以抱上麦的主持人列表
func getHoldUpSeatCompereList(c *gin.Context, roomId string) (list []response_room.UpSeatApplyInfo) {
	currUserId := helper.GetUserId(c)
	// 查询房间主持人列表
	userIdList := new(auth.Auth).GetRoomRoleListByRoleId(roomId, enum.CompereRoleId)
	var dstList []string
	for _, userId := range userIdList {
		// 过滤自己
		if currUserId == userId {
			continue
		}
		// 当前用户是否在房间内
		if !service_room.IsUserInRoom(roomId, userId) {
			continue
		}
		// 玩家是否在麦
		if isUserInSeat(roomId, userId) {
			continue
		}
		dstList = append(dstList, userId)
	}
	if len(dstList) == 0 {
		return
	}
	userInfoList := service_user.GetUserBaseInfoList(dstList)
	for _, info := range userInfoList {
		list = append(list, response_room.UpSeatApplyInfo{
			UserId:     info.Id,
			UserNo:     info.UserNo,
			Uid32:      cast.ToInt32(info.OriUserNo),
			Nickname:   info.Nickname,
			Avatar:     info.Avatar,
			Sex:        info.Sex,
			UserPlaque: service_user.GetUserLevelPlaque(info.Id, helper.GetClientType(c)),
		})
	}
	return
}

// 查询可以抱上麦的用户列表
func getHoldUpSeatAllList(c *gin.Context, roomId string) (list []response_room.UpSeatApplyInfo) {
	currUserId := helper.GetUserId(c)
	// 查询在线观众列表前50
	onlineUserService := &service_room.RoomUsersOnlie{RoomId: roomId}
	res, err := onlineUserService.GetOnlineUsersIdCard(c, 60) // 需要过滤掉在麦人员和房管 多取十个
	if err != nil {
		return
	}
	var ids []string
	for _, v := range res {
		id, _ := v.Member.(string)
		ids = append(ids, id)
	}
	authService := new(auth.Auth)
	var userIdList []string
	for _, userId := range ids {
		// 过滤自己
		if currUserId == userId {
			continue
		}
		// 玩家是否在麦
		if isUserInSeat(roomId, userId) {
			continue
		}
		// 是否为超管、巡查
		if authService.IsSuperAdminRole(roomId, userId) {
			continue
		}
		userIdList = append(userIdList, userId)
		if len(userIdList) >= 50 {
			break
		}
	}
	if len(userIdList) == 0 {
		return
	}
	userInfoList := service_user.GetUserBaseInfoList(userIdList)
	for _, info := range userInfoList {
		list = append(list, response_room.UpSeatApplyInfo{
			UserId:     info.Id,
			UserNo:     info.UserNo,
			Uid32:      cast.ToInt32(info.OriUserNo),
			Nickname:   info.Nickname,
			Avatar:     info.Avatar,
			Sex:        info.Sex,
			UserPlaque: service_user.GetUserLevelPlaque(info.Id, helper.GetClientType(c)),
		})
	}
	return
}

// 执行更换房间连麦模板
func (a *ActionRoom) updateRoomTemplate(c *gin.Context, roomInfo *model.Room, templateId string, seatListCount int) map[string]response_room.RoomWheatPosition {
	//更新房间信息
	tx := coreDb.GetMasterDb().Begin()
	err := tx.Model(model.Room{}).Where("id = ?", roomInfo.Id).Update("template_id", templateId).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	seat0 := service_room.GetSeatPositionByRoomIdAndKey(roomInfo.Id, 0)
	//更新麦位信息
	err = coreRedis.GetChatroomRedis().Del(c, redisKey.RoomWheatPosition(roomInfo.Id)).Err()
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomRelateWheat,
			Msg:  nil,
		})
	}
	seatList, err := service_room.SetMicPositions(roomInfo.Id, seatListCount, 0, roomInfo.RoomType, *seat0)
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomRelateWheat,
			Msg:  nil,
		})
	}
	tx.Commit()
	return seatList
}

// 执行连麦
func (a *ActionRoom) doRoomRelateWheat(c *gin.Context, roomId, userId, content string) {
	if content == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	//获取连麦模板
	temDao := dao.RoomTemplateDao{}
	temRes, err := temDao.GetBroadcastFirst(content)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomRelateWheat,
			Msg:  nil,
		})
	}
	//踢出所有普通麦
	a.clearAllSeat(c, roomId, userId)

	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	if roomInfo.LiveType != enum.LiveTypeAnchor {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	seatList := a.updateRoomTemplate(c, roomInfo, temRes.Id, temRes.SeatListCount)

	seats := make([]*response_room.RoomWheatPosition, 0)
	if len(seatList) > 0 {
		for _, sv := range seatList {
			seats = append(seats, &sv)
		}
	}

	sort.Slice(seats, func(i, j int) bool {
		return seats[i].Id < seats[j].Id
	})
	new(service_im.ImPublicService).SendCustomMsg(roomId, seats, enum.Room_Relate_Wheat_MSG)
}
