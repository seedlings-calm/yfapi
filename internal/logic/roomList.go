package logic

import (
	"context"
	"errors"
	"sort"
	"time"
	"unicode/utf8"
	"yfapi/app/handle"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/acl"
	service_im "yfapi/internal/service/im"
	service_room "yfapi/internal/service/room"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_index "yfapi/typedef/request/index"
	request_room "yfapi/typedef/request/room"
	"yfapi/typedef/response"
	"yfapi/typedef/response/index"
	response_room "yfapi/typedef/response/room"
	"yfapi/util/easy"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Room struct {
}

//
/**
 * @description  组装房间返回到前端的信息处理， TODO: hot字段热度没处理
 * @param list []model.Room true 房间信息
 * @param  isOnanchor bool true  个播房：是否过滤未开播的个播房间
 */
func ForRoomsAction(list []model.Room, isOnanchor bool) (result []*response_room.RoomInfo) {
	for _, roomModel := range list {
		// 房间热度值
		hot, hotStr := service_room.GetRoomHot(roomModel.Id, roomModel.LiveType)
		item := &response_room.RoomInfo{
			RoomShowBaseRes: roomModel.ToShowBase(),
			Hot:             hot,
			HotStr:          hotStr,
		}
		// TODO
		seat := make([]*response_room.RoomWheatPosition, 0)
		newSeat := make([]*response_room.RoomWheatPosition, 0)
		switch roomModel.LiveType {
		case enum.LiveTypeChatroom:
			seat = append(seat, service_room.GetSeatPositionsToUsersByRoomId(roomModel.Id)...)
		case enum.LiveTypeAnchor:
			vSeat := service_room.GetSeatPositionByRoomIdAndKey(roomModel.Id, 0)
			if isOnanchor && vSeat.UserInfo.UserId == "" { //如果个播房 主播不在线过滤掉
				continue
			}
			seat = append(seat, vSeat)
		default:
			seat = append(seat, service_room.GetSeatPositionsToUsersByRoomId(roomModel.Id)...)
		}
		for _, v := range seat {
			if v.Status == enum.MicStatusUsed {
				newSeat = append(newSeat, v)
			}
		}
		if len(newSeat) == 0 {
			newSeat = append(newSeat, &response_room.RoomWheatPosition{})
		}
		if len(newSeat) > 5 {
			newSeat = newSeat[:5]
		}
		item.SeatList = newSeat
		logicRoom := Room{}
		item.ShowLabel = logicRoom.RoomWebStatus(roomModel)
		result = append(result, item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Hot > result[j].Hot
	})
	return
}

func (l *Room) Page(c *gin.Context, req *request_room.RoomListReq) (res response.BasePageRes) {
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.RoomType = req.RoomType
	roomDao := &dao.RoomDao{}
	modelList, count, err := roomDao.GetRoomList(req.Page, req.Size, req.RoomType)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	res.Data = ForRoomsAction(modelList, true)
	res.Total = count
	//TODO: banner，取后台配置，无数据，该板块隐藏不展示
	res.Banner = nil
	res.CalcHasNext()
	return
}

func (l *Room) FindOne(c *gin.Context, roomId string) (resp response_room.ChatroomDTO) {
	userId := handle.GetUserId(c)
	roomDao := dao.RoomDao{}
	roomInfo, _ := roomDao.FindOne(&model.Room{Id: roomId})

	resp = roomInfo.ToChatroomDTO()
	actionRoom := &ActionRoom{}
	actionRoom.roomInfo(roomInfo, userId, roomId, &resp)
	return
}

func (l *Room) Update(c *gin.Context, req *request_room.RoomUpdateReq) error {
	if req.Name != "" && !validateChineseCharacters(req.Name, 10) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	if req.Notice != "" && !validateChineseCharacters(req.Notice, 300) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	userId := handle.GetUserId(c)
	rule := &acl.RoomAcl{}
	ok, _ := rule.CheckUserRule(userId, req.RoomId, acl.EditRoom)
	if !ok {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}

	roomDao := dao.RoomDao{}
	roomInfo, err := roomDao.FindOne(&model.Room{Id: req.RoomId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})

	}
	backImg := roomInfo.BackgroundImg
	var isUpdateImg bool
	if req.BackgroundImg != "" {
		backImg = helper.RemovePrefixImgUrl(helper.RemoveEscapedSlashes(req.BackgroundImg))
		if backImg != roomInfo.BackgroundImg {
			isUpdateImg = true
		}
	}
	bgsrDao := dao.RoomBgsRDao{}

	bgsrRes := bgsrDao.GetBgsByImg(backImg)
	if bgsrRes.Backgroud == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}

	tx := coreDb.GetMasterDb().Begin()
	rbgsDao := dao.RoomBgsDao{}
	if len(req.Notice) > 0 {
		roomInfo.Notice = req.Notice
	}
	if len(req.Name) > 0 {
		roomInfo.Name = req.Name
	}
	if len(req.CoverImg) > 0 {
		roomInfo.CoverImg = helper.RemovePrefixImgUrl(req.CoverImg)
	}
	roomInfo.BackgroundImg = backImg
	roomInfo.UpdateTime = time.Now()
	err = rbgsDao.UpdateRoomBgs(tx, roomInfo, &bgsrRes)
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomUpdateErr,
			Msg:  nil,
		})
	}
	tx.Commit()
	if isUpdateImg {
		imComm := new(service_im.ImCommonService)
		imComm.Send(c, handle.GetUserId(c), nil, roomInfo.Id, enum.MsgCustom, helper.FormatImgUrl(bgsrRes.Backgroud), enum.Room_BackGroudImg_Update)
	}
	if len(req.Notice) > 0 {
		imComm := new(service_im.ImCommonService)
		imComm.Send(c, handle.GetUserId(c), []string{handle.GetUserId(c)}, roomInfo.Id, enum.MsgCustom, req.Notice, enum.Room_Notice_Update)
	}
	return nil
}

func (l *Room) ReportingCenter(c *gin.Context, req *request_room.ReportingCenterReq) {
	userId := handle.GetUserId(c)
	ReDao := dao.ReportingCenterDao{}
	if _, ok := enum.ReportingObject[req.Object]; !ok {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	if _, ok := enum.ReportingSence[req.Scene]; !ok {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	for k, v := range req.Pics {
		if len(v) > 0 {
			req.Pics[k] = helper.RemovePrefixImgUrl(v)
		}
	}
	add := ReDao.AddReportingCenter(userId, req.DstId, req.Content, req.Pics, req.Object, req.Scene, req.ReportTypes)
	err := ReDao.Insert(add)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeReportingCenterErr,
			Msg:  nil,
		})
	}
}

// 长度<= n 为true
func validateChineseCharacters(input string, n int) bool {
	lens := utf8.RuneCountInString(input)
	return lens <= n
}

// ApplyRoom 申请房间
func (l *Room) ApplyRoom(c *gin.Context, req *request_room.ApplyAnchorRoomReq) (res response_room.ApplyAnchorRoomRes) {
	userId := helper.GetUserId(c)
	// 参数校验
	switch req.RoomType {
	case enum.RoomTypeAnchorVoice:
		return applyAnchorRoom(userId, req)
	//case enum.RoomTypeAnchorVideo:
	//case enum.RoomTypePersonal:
	default:
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	return
}

// 申请直播间
func applyAnchorRoom(userId string, req *request_room.ApplyAnchorRoomReq) (res response_room.ApplyAnchorRoomRes) {
	// 用户信息
	userInfo, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 用户是否有主播资质
	userCerdDao := &dao.DaoUserPractitionerCerd{
		UserId: userId,
	}
	result, _ := userCerdDao.First(enum.UserPractitionerAnchor)
	if result.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrCodeNotAnchor,
			Msg:  nil,
		})
	}

	// 用户是否有语音直播间 资格取消重考
	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.FindOne(&model.Room{RoomType: enum.RoomTypeAnchorVoice, UserId: userId})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(roomInfo.Id) == 0 {
		// 用户是否有视频直播间 资格取消重考
		roomInfo, err = roomDao.FindOne(&model.Room{RoomType: enum.RoomTypeAnchorVideo, UserId: userId})
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
	}
	bgsrDao := dao.RoomBgsRDao{}
	var bgsrRes model.RoomBgsResource
	if req.BackgroundImg == "" {
		bgsrRes = bgsrDao.GetDefaultBgs()
	} else {
		bgsrRes = bgsrDao.GetBgsByImg(helper.RemovePrefixImgUrl(helper.RemoveEscapedSlashes(req.BackgroundImg)))
	}
	if bgsrRes.Backgroud == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	tx := coreDb.GetMasterDb().Begin()
	if len(roomInfo.Id) == 0 { // 没有房间 新建立
		// 生成房间ID
		roomId := coreSnowflake.GetSnowId()
		roomInfo = &model.Room{
			Id:                 roomId,
			UserId:             userId,
			RoomNo:             userInfo.UserNo,
			RoomType:           req.RoomType,
			LiveType:           2,
			TemplateId:         "2001",
			CoverImg:           helper.RemovePrefixImgUrl(req.CoverImg),
			BackgroundImg:      helper.RemovePrefixImgUrl(req.BackgroundImg),
			Notice:             req.RoomNotice,
			Name:               req.RoomName,
			IsHot:              false,
			Status:             1,
			DaySettleUserId:    userId,
			MonthSettleUserId:  userId,
			CreateTime:         time.Now(),
			UpdateTime:         time.Now(),
			HiddenStatus:       2,
			PublicScreenStatus: 1,
			FreedMicStatus:     2,
			FreedSpeakStatus:   1,
			GuildId:            "0",
		}
		err = tx.Create(roomInfo).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	} else { // 有房间 重新打开
		roomInfo.RoomNo = userInfo.UserNo
		roomInfo.RoomType = req.RoomType
		roomInfo.RoomNo = req.RoomNotice
		roomInfo.Name = req.RoomName
		roomInfo.CoverImg = helper.RemovePrefixImgUrl(req.CoverImg)
		roomInfo.BackgroundImg = helper.RemovePrefixImgUrl(req.BackgroundImg)
		roomInfo.Status = 1
		roomInfo.CreateTime = time.Now() // 重置创建时间 月结补贴使用
		roomInfo.UpdateTime = time.Now()
		err = tx.Save(*roomInfo).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	}

	rbgsDao := dao.RoomBgsDao{}
	err = rbgsDao.UpdateRoomBgs(tx, roomInfo, &bgsrRes)
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 添加权限
	data := model.AuthRoleAccess{
		UserID: userId,
		RoleID: enum.AnchorRoleId,
		RoomID: roomInfo.Id,
	}
	err = tx.Create(data).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 添加为房间从业者
	practitioner := &model.UserPractitioner{
		Id:               0,
		RoomId:           roomInfo.Id,
		UserId:           userId,
		PractitionerType: enum.UserPractitionerAnchor,
		Status:           1,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
	err = tx.Create(practitioner).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	tx.Commit()
	key1 := redisKey.UserRules(userId, roomInfo.Id)
	key2 := redisKey.UserRoles(userId, roomInfo.Id)
	key3 := redisKey.UserCompereRules(userId, roomInfo.Id)
	coreRedis.GetUserRedis().Del(context.Background(), key1, key2, key3)
	// TODO 把房间加入在线房间列表
	res.RoomId = roomInfo.Id
	return
}

// ApplyRoomInfo 申请房间回显信息
func (l *Room) ApplyRoomInfo(c *gin.Context) (res response_room.ApplyRoomInfoRes) {
	userId := helper.GetUserId(c)
	// 用户是否有主播资质
	userCerdDao := &dao.DaoUserPractitionerCerd{
		UserId: userId,
	}
	result, _ := userCerdDao.First(enum.UserPractitionerAnchor)
	if result.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrCodeNotAnchor,
			Msg:  nil,
		})
	}
	// 用户是否有语音直播间
	roomDao := new(dao.RoomDao)
	roomInfo, _ := roomDao.FindOne(&model.Room{RoomType: enum.RoomTypeAnchorVoice, UserId: userId, Status: 1})

	if len(roomInfo.Id) == 0 {
		// 用户是否有视频直播间 资格取消重考
		roomInfo, _ = roomDao.FindOne(&model.Room{RoomType: enum.RoomTypeAnchorVideo, UserId: userId, Status: 1})
	}
	res = response_room.ApplyRoomInfoRes{
		RoomId:        roomInfo.Id,
		RoomNo:        roomInfo.RoomNo,
		RoomName:      roomInfo.Name,
		CoverImg:      helper.FormatImgUrl(roomInfo.CoverImg),
		LiveType:      roomInfo.LiveType,
		RoomType:      roomInfo.RoomType,
		Notice:        roomInfo.Notice,
		BackgroundImg: helper.FormatImgUrl(roomInfo.BackgroundImg),
	}
	return
}

// 房间前端展示标签
func (l *Room) RoomWebStatus(roomInfo model.Room) int {

	// RoomShowStatusClose 未开播（直播）
	// RoomShowStatusAnchoring 直播中(个播房主播开播中)
	// RoomShowStatusInteraction 互动中(情感、听歌、交友，老板位有人)
	// RoomShowStatusInteractWait 等待互动(情感、听歌、交友，老板位没有人)
	// RoomShowStatusSing 演唱中(听歌，演唱位有人)
	seatList := service_room.GetSeatPositionsByRoomId(roomInfo.Id)
	// 房间类型总和
	var RoomTypeList = []int{
		enum.RoomTypeEmoMan, enum.RoomTypeEmoWoman, enum.RoomTypeFriend, enum.RoomTypeSing,
	}
	switch roomInfo.LiveType {
	case enum.LiveTypeChatroom: //1:聊天室
		if easy.InArray(roomInfo.RoomType, RoomTypeList) { //情感，听歌，交友
			if len(seatList) > 1 {
				if seatList[1].UserInfo.UserId != "" {
					return enum.RoomShowStatusInteraction
				}
			}
			return enum.RoomShowStatusInteractWait
		} else {
			return enum.RoomShowStatusInteractWait
		}
	case enum.LiveTypeAnchor: //2:直播
		if seatList[0].UserInfo.UserId != "" {
			return enum.RoomShowStatusSing
		} else {
			return enum.RoomShowStatusClose
		}
	default:
		return enum.RoomShowStatusInteractWait
	}
}

// 需要使用热度值排序 包括取30条
func (l *Room) GetTop(c *gin.Context, req *request_index.TopListReq) (res response.BasePageRes) {
	var (
		ids      []interface{}
		redisArr []redis.Z
	)

	chatRooms, _ := service_room.GetRoomHotsList(c, enum.LiveTypeChatroom, 30)
	redisArr = append(redisArr, chatRooms...)
	anchors, _ := service_room.GetRoomHotsList(c, enum.LiveTypeAnchor, 30)
	redisArr = append(redisArr, anchors...)
	sort.Slice(redisArr, func(i, j int) bool {
		return redisArr[i].Score > redisArr[j].Score
	})
	if len(redisArr) > 30 {
		redisArr = redisArr[:30]
	}
	for _, v := range redisArr {
		ids = append(ids, v.Member)
	}
	roomDao := &dao.RoomDao{}
	modelList, err := roomDao.GetRoomByIdsInterface(ids)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	res.Total = int64(len(modelList))
	res.Data = ForRoomsAction(modelList, false)
	//TODO: banner，取后台配置，无数据，该板块隐藏不展示
	res.Banner = nil
	return
}

// SearchAll 搜索用户、聊天室、直播间
func (l *Room) SearchAll(c *gin.Context, req *request_index.SearchAllReq) (res index.SearchAllRes) {
	userId := helper.GetUserId(c)
	switch req.SearchType {
	case 1:
		res.UserList = searchUser(userId, req.Keyword, helper.GetClientType(c))
	case 2:
		res.ChatroomList = searchRoom(enum.LiveTypeChatroom, req.Keyword)
	case 3:
		res.AnchorRoomList = searchRoom(enum.LiveTypeAnchor, req.Keyword)
	default:
		res.UserList = searchUser(userId, req.Keyword, helper.GetClientType(c))
		res.ChatroomList = searchRoom(enum.LiveTypeChatroom, req.Keyword)
		res.AnchorRoomList = searchRoom(enum.LiveTypeAnchor, req.Keyword)
	}
	return
}

// 查询用户
func searchUser(userId, keyword, clientType string) (res []*index.SearchUserInfo) {
	userList, err := new(dao.UserDao).FindUserByKeyword(keyword, 0, 50)
	if err != nil {
		return
	}
	var userIdList []string
	for _, data := range userList {
		if userId == data.Id {
			continue
		}
		resp := &index.SearchUserInfo{
			UserId:     data.Id,
			UserNo:     data.UserNo,
			Uid32:      cast.ToInt32(data.OriUserNo),
			Nickname:   data.Nickname,
			Avatar:     helper.FormatImgUrl(data.Avatar),
			Sex:        data.Sex,
			Introduce:  data.Introduce,
			IsOnline:   service_user.IsOnline(data.Id),
			InRoom:     service_user.GetUserInRoomInfo(data.Id, clientType),
			UserPlaque: service_user.GetUserLevelPlaque(data.Id, clientType),
		}
		userIdList = append(userIdList, data.Id)
		res = append(res, resp)
	}
	userFollowDao := new(dao.UserFollowDao)
	// 我关注的列表
	followMap, _ := userFollowDao.GetUserFollowMap(userId, userIdList)
	// 关注我的列表
	followerMap := userFollowDao.GetUserFollowerMap(userId, userIdList)
	for _, info := range res {
		follow, isOk := followMap[info.UserId]
		if !isOk {
			// 未关注对方
			info.FollowedType = 0
			// 对方是否关注我
			if _, isExist := followerMap[info.UserId]; isExist {
				info.FollowedType = 3
			}
			continue
		}
		info.FollowedType = 1
		if follow.IsMutualFollow {
			info.FollowedType = 2
		}
	}

	return
}

// 查询房间
func searchRoom(liveType int, keyword string) (res []*response_room.RoomInfo) {
	data, err := new(dao.RoomDao).FindRoomByKeyword(keyword, liveType, 0, 50)
	if err != nil {
		return
	}
	return ForRoomsAction(data, false)
}
