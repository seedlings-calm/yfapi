package logic

import (
	"context"
	"encoding/json"
	"math"
	"regexp"
	"time"
	"unicode/utf8"
	"yfapi/app/handle"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/acl"
	service_im "yfapi/internal/service/im"
	"yfapi/internal/service/riskCheck/shumei"
	service_room "yfapi/internal/service/room"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_room "yfapi/typedef/request/room"
	response_room "yfapi/typedef/response/room"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// 房间在线用户列表 redis缓存 5s
func (a *ActionRoom) OnlineUsers(c *gin.Context, roomId string) (resp response_room.OnlineUsersResponse) {
	userId := handle.GetUserId(c)
	//如果用户不再房间内，直接返回错误
	// isOnline := service_room.IsUserInRoom(roomId, userId)
	// if !isOnline {
	// 	return
	// }
	resp.IsShowNums = false
	var cacheKey string

	authDao := new(dao.UserAuthDao)
	//榜单值全部展示给从业者、房主、超管、巡查 ,所以需要检索用户身份
	//ok, _ := authDao.IsRoles(userId, roomId, []string{"1001", "1002", "1004", "1005", "1007", "1008", "1009"})
	ok, _ := authDao.IsRoles(userId, roomId, []string{cast.ToString(enum.SuperAdminRoleId), cast.ToString(enum.PatrolRoleId), cast.ToString(enum.HouseOwnerRoleId), cast.ToString(enum.CompereRoleId), cast.ToString(enum.MusicianRoleId), cast.ToString(enum.CounselorRoleId), cast.ToString(enum.AnchorRoleId)})
	if ok {
		resp.IsShowNums = true
		cacheKey = redisKey.RoomOnlineUsersIdCardCache(roomId)
	} else {
		cacheKey = redisKey.RoomOnlineUsersCache(roomId)
	}
	redisCli := coreRedis.GetChatroomRedis()
	str := redisCli.Get(c, cacheKey).Val()
	onlineUserService := &service_room.RoomUsersOnlie{RoomId: roomId}
	userDao := dao.UserDao{}

	//有缓存
	var cacheResp response_room.OnlineUsersResponse
	if str != "" {
		err := json.Unmarshal([]byte(str), &cacheResp)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrCodeOnlineUserErr,
				Msg:  nil,
			})
		}
	}
	info := onlineUserService.GetMemberWithNeighbors(roomId, userId)
	ownerInfo, err := userDao.FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeOnlineUserErr,
			Msg:  nil,
		})
	}
	resp.OwnerInfo.Avatar = coreConfig.GetHotConf().ImagePrefix + ownerInfo.Avatar
	resp.OwnerInfo.Nickname = ownerInfo.Nickname
	resp.OwnerInfo.Sex = ownerInfo.Sex
	resp.OwnerInfo.UserId = userId
	resp.OwnerInfo.UserNo = ownerInfo.UserNo
	resp.OwnerInfo.Contribution = "0"
	resp.OwnerInfo.OriginalContribution = "0"
	resp.OwnerInfo.UserPlaque = service_user.GetUserLevelPlaque(userId, helper.GetClientType(c))

	if info["owner"] > 0 {
		resp.OwnerInfo.OriginalContribution = cast.ToString(info["owner"])
		resp.OwnerInfo.Contribution = easy.FormatLeaderboardValue(info["owner"])
	}
	if info["first"] > 0 {
		resp.OwnerInfo.UpgradeNum = easy.FormatLeaderboardValue(info["first"] - info["owner"])
	} else {
		resp.OwnerInfo.UpgradeNum = easy.FormatLeaderboardValue(info["owner"] - info["end"])
	}
	if str != "" {
		resp.DayUsersCount = cacheResp.DayUsersCount
		resp.OnlineLists = cacheResp.OnlineLists
		resp.OnlineUsersCount = cacheResp.OnlineUsersCount
		return
	}

	res, err := onlineUserService.GetOnlineUsersIdCard(c, 200)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeOnlineUserErr,
			Msg:  nil,
		})
	}
	var ids []string
	for _, v := range res {
		id := v.Member.(string)
		ids = append(ids, id)
	}
	userRes := userDao.FindByIds(ids)
	var (
		list    []*response_room.RoomUsersBase
		mapList = make(map[string]model.User)
	)
	for _, v := range userRes {
		mapList[v.Id] = v
	}
	for _, v := range res {
		id := v.Member.(string)
		vScore := math.Trunc(v.Score)
		onUsers := &response_room.RoomUsersBase{
			UserId:       id,
			UserNo:       mapList[id].UserNo,
			Uid32:        cast.ToInt32(mapList[id].OriUserNo),
			Nickname:     mapList[id].Nickname,
			Avatar:       coreConfig.GetHotConf().ImagePrefix + mapList[id].Avatar,
			Sex:          mapList[id].Sex,
			Contribution: easy.FormatLeaderboardValue(vScore),
			UserPlaque:   service_user.GetUserLevelPlaque(id, helper.GetClientType(c)),
		}
		list = append(list, onUsers)
	}
	resp.OnlineUsersCount = int(redisCli.ZCount(c, redisKey.RoomUsersOnlineList(roomId), "0", "+inf").Val())
	resp.OnlineLists = list
	resp.DayUsersCount = int(redisCli.ZCount(c, redisKey.RoomUsersDayList(roomId), "1000", "+inf").Val())
	go func() {
		tb, err := json.Marshal(resp)
		if err == nil {
			coreRedis.GetChatroomRedis().Set(c, redisKey.RoomOnlineUsersIdCardCache(roomId), string(tb), 8*time.Second)
		}
	}()
	return
}

// 房间1000贡献榜列表 缓存10s
func (a *ActionRoom) DayUsers(c *gin.Context, roomId string) (resp response_room.DayUsersResponse) {
	userId := handle.GetUserId(c)
	//如果用户不再房间内，直接返回错误
	isOnline := service_room.IsUserInRoom(roomId, userId)
	if !isOnline {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}

	cacheKey := redisKey.RoomDayUsersCache(roomId)
	redisCli := coreRedis.GetChatroomRedis()
	onlineService := &service_room.RoomUsersOnlie{RoomId: roomId}
	userDao := dao.UserDao{}

	str := redisCli.Get(c, cacheKey).Val()
	if str != "" {
		err := json.Unmarshal([]byte(str), &resp)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrCodeRoomPermissionDenied,
				Msg:  nil,
			})
		}
	}
	ownerInfo, err := userDao.FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}
	resp.OwnerInfo.Avatar = coreConfig.GetHotConf().ImagePrefix + ownerInfo.Avatar
	resp.OwnerInfo.Nickname = ownerInfo.Nickname
	resp.OwnerInfo.Sex = ownerInfo.Sex
	resp.OwnerInfo.UserId = userId
	resp.OwnerInfo.UserNo = ownerInfo.UserNo
	score := math.Trunc(redisCli.ZScore(c, redisKey.RoomUsersDayList(roomId), userId).Val())
	resp.OwnerInfo.OriginalContribution = cast.ToString(score)
	resp.OwnerInfo.Contribution = easy.FormatLeaderboardValue(score)
	resp.OwnerInfo.UserPlaque = service_user.GetUserLevelPlaque(userId, helper.GetClientType(c))
	resp.OwnerInfo.UpgradeNum = "0"
	if score < 1000 {
		resp.OwnerInfo.UpgradeNum = easy.FormatLeaderboardValue(1000 - score)
	}
	if str != "" {
		return
	}
	resp.NoOnlineUsers = make([]*response_room.RoomUsersBase, 0)
	resp.OnlineUsers = make([]*response_room.RoomUsersBase, 0)
	userService := &service_room.RoomUserDay{RoomId: roomId}
	res := userService.GetDayUsers(c)
	var (
		ids          []string
		userInfoList = make(map[string]model.User)
		Isonlines    = make(map[string]bool)
	)
	resLen := len(res)
	if resLen != 0 {
		for _, v := range res {
			id := v.Member.(string)
			ids = append(ids, id)
		}
		userRes := userDao.FindByIds(ids)
		for _, v := range userRes {
			userInfoList[v.Id] = v
		}
		//判断列表用户是否在线
		Isonlines = onlineService.IsUserOnline(c, ids)
		for _, v := range res {
			id := v.Member.(string)
			vScore := math.Trunc(v.Score)
			onUsers := &response_room.RoomUsersBase{
				UserId:       id,
				UserNo:       userInfoList[id].UserNo,
				Uid32:        cast.ToInt32(userInfoList[id].OriUserNo),
				Nickname:     userInfoList[id].Nickname,
				Avatar:       coreConfig.GetHotConf().ImagePrefix + userInfoList[id].Avatar,
				Sex:          userInfoList[id].Sex,
				Contribution: easy.FormatLeaderboardValue(vScore),
				UserPlaque:   service_user.GetUserLevelPlaque(id, helper.GetClientType(c)),
			}
			if val, ok := Isonlines[id]; ok && val {
				resp.OnlineUsers = append(resp.OnlineUsers, onUsers)
			} else {
				resp.NoOnlineUsers = append(resp.NoOnlineUsers, onUsers)
			}
		}
	}
	resp.NoOnlineUsersCount = len(resp.NoOnlineUsers)
	resp.OnlineUsersCount = len(resp.OnlineUsers)
	go func() {
		tb, err := json.Marshal(resp)
		if err == nil {
			coreRedis.GetChatroomRedis().Set(c, cacheKey, string(tb), 10*time.Second)
		}
	}()
	return
}

func (a *ActionRoom) RoomLock(c *gin.Context, req *request_room.RoomLockReq) {
	userId := handle.GetUserId(c)
	rule := &acl.RoomAcl{}
	ok, _ := rule.CheckUserRule(userId, req.RoomId, acl.LockRoom)
	if !ok {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}
	roomDao := dao.RoomDao{}
	roomInfo, _ := roomDao.FindOne(&model.Room{Id: req.RoomId})
	if roomInfo.RoomPwd == "" {
		if !regexp.MustCompile(`^\d{4}$`).MatchString(req.Pwd) {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		roomInfo.RoomPwd = easy.Md5(req.Pwd, 16, true)
		roomInfo.UpdateTime = time.Now()
	} else {
		roomInfo.RoomPwd = ""
		roomInfo.UpdateTime = time.Now()
	}
	err := roomDao.Save(*roomInfo)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomUpdateErr,
			Msg:  nil,
		})
	}
}

// 自动欢迎语
func (a *ActionRoom) AutoWelcome(c *gin.Context, content string) {
	userId := handle.GetUserId(c)
	autoDao := dao.AutoWelcome{}
	if content == "" { //删除
		err := autoDao.Del(userId)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrCodeAutoWelcomeErr,
				Msg:  nil,
			})
		}
	} else { //添加，编辑
		shumeiSer := shumei.ShuMei{}
		if !shumeiSer.MomentsCheck(userId, content) {
			panic(error2.I18nError{
				Code: error2.ErrorCodeTextCheckReject,
				Msg:  nil,
			})
		}
		lens := utf8.RuneCountInString(content)
		if lens < 10 || lens > 50 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}

		err := autoDao.Save(&model.UserAutoWelcome{UserID: userId, WelcomeContent: content})
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrCodeAutoWelcomeErr,
				Msg:  nil,
			})
		}
	}

}

// 用户禁言
func (a *ActionRoom) UserMute(c *gin.Context, userId, targetUserId, roomId string, minute int) {
	if !new(acl.RoomAcl).CompareRole(userId, targetUserId, roomId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNoPermissions,
			Msg:  nil,
		})
	}
	muteInfo := new(dao.UserMuteListDao).FindOne(targetUserId, roomId)
	if muteInfo.ID > 0 { //已禁言
		if minute == 0 {
			new(UserMute).UnMute(c, targetUserId, roomId)
			new(service_im.ImCommonService).Send(c, userId, []string{targetUserId}, roomId, enum.MsgCustom, "", enum.MUTE_USER_MSG)
		} else {
			panic(error2.I18nError{
				Code: error2.ErrCodeUserIsMute,
				Msg:  nil,
			})
		}
	} else { //禁言
		new(UserMute).UnMute(c, targetUserId, roomId)
		err := new(UserMute).Mute(c, targetUserId, roomId, minute)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrCodeUserMuteFail,
				Msg:  nil,
			})
		}
		new(service_im.ImCommonService).Send(c, userId, []string{targetUserId}, roomId, enum.MsgCustom, i18n_msg.GetI18nMsg(c, i18n_msg.YouHaveBeenMutedMsgKey, map[string]any{"minute": minute}), enum.MUTE_USER_MSG)
		seatInfo := new(acl.RoomAcl).GetUserMicInfo(roomId, targetUserId)
		if seatInfo != nil && !seatInfo.Mute { //静音关闭状态
			seatInfo.Mute = !seatInfo.Mute
			err := coreRedis.GetChatroomRedis().HSet(context.Background(), redisKey.RoomWheatPosition(roomId), seatInfo.Id, easy.JSONStringFormObject(seatInfo)).Err()
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
		}
		// 推送麦位静音消息
		new(service_im.ImPublicService).SendCustomMsg(roomId, seatInfo, enum.MUTE_SEAT_MSG)
	}
}

func (a *ActionRoom) GetRoomExtraInfo(c *gin.Context, roomId string) (res response_room.ChatroomExtra) {
	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	res.RoomId = roomId
	res.IsFreedMic = roomInfo.FreedMicStatus == 1
	res.IsFreedSpeak = roomInfo.FreedSpeakStatus == 1
	res.IsPublicScreen = roomInfo.PublicScreenStatus == 1
	return
}

func (a *ActionRoom) GetRoomBgs(c *gin.Context, roomId string) (res []*response_room.GetRoomBgsRes) {
	roomDao := dao.RoomDao{}
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	if roomInfo.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}

	bgsrDao := dao.RoomBgsRDao{}
	bgsrRes := bgsrDao.GetBgs(0)
	if len(bgsrRes) == 0 {
		return
	}
	bgsDao := dao.RoomBgsDao{}
	bgsRes, _ := bgsDao.GetRoomBgs(roomId)
	useId := bgsRes.TrbrId //使用中的背景ID
	if useId == 0 {
		//默认背景第一次使用走此处
		itemRes := bgsrDao.GetDefaultBgs()
		tx := coreDb.GetMasterDb().Begin()
		err = bgsDao.UpdateRoomBgs(tx, roomInfo, &itemRes)
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
		useId = int(itemRes.Id)
	}
	for _, v := range bgsrRes {
		item := new(response_room.GetRoomBgsRes)
		item.TrbrId = int(v.Id)
		item.Name = v.Name
		item.Types = v.Types
		item.Backgroud = helper.FormatImgUrl(v.Backgroud)
		item.Icon = helper.FormatImgUrl(v.Icon)
		item.CreateTime = v.CreateTime
		if useId == item.TrbrId {
			item.IsUse = 2
		} else {
			item.IsUse = 1
		}
		res = append(res, item)
	}
	return
}

func (a *ActionRoom) SetRoomBgs(c *gin.Context, roomId string, bgId int) (err error) {
	roomDao := dao.RoomDao{}
	roomInfo, err := roomDao.FindOne(&model.Room{Id: roomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	if roomInfo.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	bgsrDao := dao.RoomBgsRDao{}

	bgres := bgsrDao.GetBgsById(bgId)
	if bgres.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	bgsDao := dao.RoomBgsDao{}
	tx := coreDb.GetMasterDb().Begin()
	err = bgsDao.UpdateRoomBgs(tx, roomInfo, &bgres)
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	imComm := new(service_im.ImCommonService)
	imComm.Send(c, handle.GetUserId(c), nil, roomId, enum.MsgCustom, helper.FormatImgUrl(bgres.Backgroud), enum.Room_BackGroudImg_Update)
	return
}

func (a *ActionRoom) GetBgs(c *gin.Context) (res []model.RoomBgsResource) {
	bgsrDao := dao.RoomBgsRDao{}
	res = bgsrDao.GetBgs(1)
	for k, v := range res {
		res[k].Icon = helper.FormatImgUrl(v.Icon)
		res[k].Backgroud = helper.FormatImgUrl(v.Backgroud)
	}
	return
}

// 房间高等级用户列表 redis缓存 30s
func (a *ActionRoom) HighGradeUsers(c *gin.Context, roomId string) (resp response_room.HighGradeUsersResponse) {
	resp = response_room.HighGradeUsersResponse{
		FirstCount:  0,
		SecondCount: 0,
		ThreeCount:  0,
		FirstList:   make([]*response_room.RoomUsersBase, 0),
		SecondList:  make([]*response_room.RoomUsersBase, 0),
		ThreeList:   make([]*response_room.RoomUsersBase, 0),
	}
	userId := handle.GetUserId(c)
	//如果用户不再房间内，直接返回错误
	isOnline := service_room.IsUserInRoom(roomId, userId)
	if !isOnline {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}
	redisCli := coreRedis.GetChatroomRedis()
	str := redisCli.Get(c, redisKey.RoomHightGradeUsersCache(roomId)).Val()
	onlineUserService := &service_room.RoomUsersOnlie{RoomId: roomId}
	userDao := dao.UserDao{}

	//有缓存
	if str != "" {
		err := json.Unmarshal([]byte(str), &resp)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrCodeHighGradeUsersErr,
				Msg:  nil,
			})
		}
	}
	ownerInfo, err := userDao.FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeHighGradeUsersErr,
			Msg:  nil,
		})
	}
	resp.OwnerInfo.Avatar = coreConfig.GetHotConf().ImagePrefix + ownerInfo.Avatar
	resp.OwnerInfo.Nickname = ownerInfo.Nickname
	resp.OwnerInfo.Sex = ownerInfo.Sex
	resp.OwnerInfo.UserId = userId
	resp.OwnerInfo.UserNo = ownerInfo.UserNo
	resp.OwnerInfo.UserPlaque = service_user.GetUserLevelPlaque(userId, helper.GetClientType(c))

	if str != "" { //使用缓存
		return
	}
	redisRes := onlineUserService.GetHightGradeUsers(c)
	if len(redisRes) == 0 {
		return
	}
	var ids []string
	for _, v := range redisRes {
		id := v.Member.(string)
		ids = append(ids, id)
	}
	userLvLevelDao := dao.UserLevelLvDao{}
	lvLevelIds := userLvLevelDao.GetUserIdsByLvLevel(ids, 31)
	if len(lvLevelIds) == 0 {
		return
	}

	var (
		mapList   = make(map[string]model.User)
		mapLvList = make(map[string]int)
	)
	for _, v := range lvLevelIds {
		mapLvList[v.UserId] = v.Level
	}
	userRes := userDao.FindByIds(ids)
	for _, v := range userRes {
		mapList[v.Id] = v
	}
	for _, v := range redisRes {
		id := v.Member.(string)

		item := &response_room.RoomUsersBase{
			UserId:       id,
			UserNo:       mapList[id].UserNo,
			Uid32:        cast.ToInt32(mapList[id].OriUserNo),
			Nickname:     mapList[id].Nickname,
			Avatar:       coreConfig.GetHotConf().ImagePrefix + mapList[id].Avatar,
			Sex:          mapList[id].Sex,
			UserPlaque:   service_user.GetUserLevelPlaque(id, helper.GetClientType(c)),
			Contribution: cast.ToString(v.Score),
		}
		if mapLvList[id] >= 51 { //51
			resp.FirstCount += 1
			resp.FirstList = append(resp.FirstList, item)
		} else if mapLvList[id] >= 41 && mapLvList[id] <= 50 { //41-50
			resp.SecondCount += 1
			resp.SecondList = append(resp.SecondList, item)
		} else { //31-40
			resp.ThreeCount += 1
			resp.ThreeList = append(resp.ThreeList, item)
		}
	}
	go func() {
		tb, err := json.Marshal(resp)
		if err == nil {
			coreRedis.GetChatroomRedis().Set(c, redisKey.RoomHightGradeUsersCache(roomId), string(tb), 30*time.Second)
		}
	}()
	return
}

// 高等级用户的统计数据-redis缓存30s
func (a *ActionRoom) HighGradeUsersCount(c *gin.Context, roomId string) (resp response_room.HighGradeUsersCountResponse) {
	var LevelList []response_room.HighGradeUsersCountItem
	for i := 0; i < 3; i++ {
		var rule string
		if i == 0 {
			rule = "high"
		} else if i == 1 {
			rule = "mid"
		} else {
			rule = "low"
		}
		LevelList = append(LevelList, response_room.HighGradeUsersCountItem{Rule: rule})
	}
	resp = response_room.HighGradeUsersCountResponse{
		LevelList: LevelList,
	}
	userId := handle.GetUserId(c)
	//如果用户不再房间内，直接返回错误
	isOnline := service_room.IsUserInRoom(roomId, userId)
	if !isOnline {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}
	redisCli := coreRedis.GetChatroomRedis()
	str := redisCli.Get(c, redisKey.RoomHightGradeUsersCountCache()).Val()
	//有缓存
	if str != "" {
		err := json.Unmarshal([]byte(str), &resp)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrCodeHighGradeUsersErr,
				Msg:  nil,
			})
		}
		return
	}
	onlineUserService := &service_room.RoomUsersOnlie{RoomId: roomId}
	redisRes := onlineUserService.GetHightGradeUsers(c)
	if len(redisRes) == 0 {
		return
	}
	var ids []string
	for _, v := range redisRes {
		id := v.Member.(string)
		ids = append(ids, id)
	}
	userLvLevelDao := dao.UserLevelLvDao{}
	lvLevelIds := userLvLevelDao.GetUserIdsByLvLevel(ids, 31)
	if len(lvLevelIds) == 0 {
		return
	}
	for _, v := range lvLevelIds {
		if v.Level >= 51 {
			resp.LevelList[0].Count++
		} else if v.Level <= 50 && v.Level >= 40 {
			resp.LevelList[1].Count++
		} else {
			resp.LevelList[2].Count++
		}
	}
	//获取等级的logoIcon
	lvDao := dao.LvConfigDao{}
	first, _ := lvDao.GetLvConfigByLevel(50)
	second, _ := lvDao.GetLvConfigByLevel(40)
	three, _ := lvDao.GetLvConfigByLevel(30)
	if first != nil {
		resp.LevelList[0].Icon = helper.FormatImgUrl(first.LogoIcon)
	}
	if second != nil {
		resp.LevelList[1].Icon = helper.FormatImgUrl(second.LogoIcon)
	}
	if three != nil {
		resp.LevelList[2].Icon = helper.FormatImgUrl(three.LogoIcon)
	}

	go func() {
		tb, err := json.Marshal(resp)
		if err == nil {
			coreRedis.GetChatroomRedis().Set(c, redisKey.RoomHightGradeUsersCountCache(), string(tb), 30*time.Second)
		}
	}()
	return
}

// 房间广告位
func (a *ActionRoom) RoomAdvertising(c *gin.Context) (resp []response_room.AdvertisingResp) {
	resp = []response_room.AdvertisingResp{
		{
			Id:       1,
			OpenType: 1,
			Ratio:    50,
			Image:    "",
			Url:      "",
			SiteType: 1,
		},
		{
			Id:       2,
			OpenType: 2,
			Ratio:    100,
			Image:    "",
			Url:      "",
			SiteType: 2,
		},
	}
	return
}
