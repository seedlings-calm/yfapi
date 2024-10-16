package guild

import (
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	dao2 "yfapi/internal/dao/guild"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	model2 "yfapi/internal/model/guild"
	service_im "yfapi/internal/service/im"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_room "yfapi/typedef/request/guild"
	"yfapi/typedef/response"
	response_guild "yfapi/typedef/response/guild"
)

type GuildRoom struct {
}

func (g *GuildRoom) RoomList(c *gin.Context, req *request_room.GuildRoomListreq) (resp response.AdminPageRes) {
	list, count, err := new(dao2.RoomListDao).GetRoomListPage(req, c)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	resp.Total = count
	resp.Data = list
	resp.CurrentPage = req.CurrentPage
	resp.Size = req.Size
	return
}

// RoomTypeList
//
//	@Description: 房间类型列表
//	@receiver g
//	@param req *request_room.RoomTypeReq -
//	@return resp -
func (g *GuildRoom) RoomTypeList(req *request_room.RoomTypeReq) (resp response_guild.RoomTypeInfo) {
	tx := coreDb.GetSlaveDb().Model(model.RoomType{})
	if req.LiveType > 0 {
		tx = tx.Where("live_type", req.LiveType)
	}
	err := tx.Select("type_id, type_name, live_type").Order("sort_num").Scan(&resp.TypeList).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	tx = coreDb.GetSlaveDb().Model(model.RoomTemplate{}).Where("status", 1)
	if req.LiveType > 0 {
		tx = tx.Where("live_type", req.LiveType)
	}
	err = tx.Select("id, template_name, live_type").Scan(&resp.TemplateList).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	return
}
func (g *GuildRoom) UpdateRoom(c *gin.Context, req *request_room.ChangeRoomParamReq) any {
	guildId := helper.GetGuildId(c)
	//公会信息校验
	guildDao := new(dao2.GuildDao)
	guildInfo, err := guildDao.FindOne(&model.Guild{ID: guildId})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(guildInfo.ID) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeGuildNotExist,
			Msg:  nil,
		})
	}
	// 房间信息
	roomInfo, err := new(dao.RoomDao).GetRoomById(req.RoomID)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(roomInfo.Id) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	if roomInfo.LiveType == typedef_enum.LiveTypeAnchor { // 直播间不可改动
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 用户信息
	userDao := new(dao.UserDao)
	userInfo, err := userDao.FindUserByUserNo(req.UserNo)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if userInfo.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	//实名认证
	if userInfo.RealNameStatus != typedef_enum.UserRealNameAuthenticated {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeIDCardAuth,
			Msg:  nil,
		})
	}
	//检测是否是公会成员
	guildMemberDao := new(dao2.GuildMemberDao)
	memberInfo, err := guildMemberDao.FindOne(&model.GuildMember{
		GuildID: guildId,
		UserID:  userInfo.Id,
		Status:  1,
	})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if memberInfo.ID == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotInGuild,
		})
	}

	//执行更换
	roomDao := new(dao.RoomDao)
	switch req.Op {
	case 1: // 更换房主
		if userInfo.Id == roomInfo.UserId { // 无变动
			return nil
		}
		// 变更房主
		tx := coreDb.GetMasterDb().Begin()
		err = tx.Model(model.Room{}).Where("id", req.RoomID).Updates(map[string]interface{}{
			"user_id":              userInfo.Id,
			"day_settle_user_id":   userInfo.Id,
			"month_settle_user_id": userInfo.Id,
		}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		// 删除旧房主的房主权限
		err = tx.Model(model.AuthRoleAccess{}).Where("user_id=? and room_id=? and role_id=?", roomInfo.UserId, roomInfo.Id, typedef_enum.HouseOwnerRoleId).Delete(model.AuthRoleAccess{}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		// 新增新房主的房主权限
		err = tx.Create(&model.AuthRoleAccess{
			UserID: userInfo.Id,
			RoleID: typedef_enum.HouseOwnerRoleId,
			RoomID: roomInfo.Id,
		}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		tx.Commit()
		// 删除身份权限缓存
		key1 := redisKey.UserRules(userInfo.Id, req.RoomID)
		key2 := redisKey.UserRoles(userInfo.Id, req.RoomID)
		key3 := redisKey.UserRoles(roomInfo.UserId, req.RoomID)
		key4 := redisKey.UserRoles(roomInfo.UserId, req.RoomID)
		coreRedis.GetUserRedis().Del(c, key1, key2, key3, key4)
		return err
	case 2: // 更换日结算人
		if userInfo.Id != roomInfo.UserId && userInfo.Id != guildInfo.UserID {
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeDaySettleInvalidRole,
				Msg:  nil,
			})
		}
		err = roomDao.Update(model.Room{DaySettleUserId: userInfo.Id, Id: req.RoomID})
		return err
	case 3: // 更换月结算人
		if userInfo.Id != roomInfo.UserId && userInfo.Id != guildInfo.UserID {
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeMonthSettleInvalidRole,
				Msg:  nil,
			})
		}
		err = roomDao.Update(model.Room{MonthSettleUserId: userInfo.Id, Id: req.RoomID})
		return err
	}
	return nil
}

// 通过userno返回基本用户信息
func (g *GuildRoom) GetUserInfoByUserNo(c *gin.Context, req *request_room.UserNoParamReq) any {
	//返回查询
	userDao := new(dao.UserDao)
	data, err := userDao.FindUserByUserNo(req.UserNo)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if data.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	checkGuild := new(dao2.GuildMemberDao)
	_, err = checkGuild.FindOne(&model.GuildMember{GuildID: c.GetString("guildId"), UserID: data.Id})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotInGuild,
			Msg:  nil,
		})
	}
	var guildUser response_guild.GuildUserBaseInfo
	guildUser.Nickname = data.Nickname
	guildUser.TrueName = data.TrueName
	guildUser.Avatar = coreConfig.GetHotConf().ImagePrefix + data.Avatar
	guildUser.UserNo = data.UserNo
	guildUser.UserId = data.Id
	guildUser.Mobile = data.Mobile
	return guildUser
}
func (g *GuildRoom) CloseRoom(c *gin.Context, req *request_room.CloseRoomReq) (err error) {
	switch req.Status {
	case typedef_enum.RoomStatusNormal, typedef_enum.RoomStatusClose:
	default:
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeParam,
			Msg:  nil,
		})
	}
	roomDao := new(dao.RoomDao)
	roomInfo, err := roomDao.GetRoomById(req.RoomID)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	if len(roomInfo.Id) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	if roomInfo.Status == req.Status { // 状态无变动
		return
	}
	if roomInfo.LiveType == typedef_enum.LiveTypeAnchor { // 直播间不可改动
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 开启房间 检查是否运营后台封禁房间
	if req.Status == typedef_enum.RoomStatusNormal {
		var punishRecord model.UserReportPunishRecord
		_ = coreDb.GetSlaveDb().Model(punishRecord).Where("dst_user_id=? and object=1 and expire_time>?", roomInfo.Id, time.Now().Format(time.DateTime)).
			Order("expire_time desc").Limit(1).Scan(&punishRecord).Error
		if punishRecord.ID > 0 {
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeRoomFreezing,
				Msg:  nil,
			})
		}
	}
	err = roomDao.Update(model.Room{Status: req.Status, Id: req.RoomID})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	if req.Status == typedef_enum.RoomStatusClose {
		new(service_im.ImPublicService).SendCustomMsg(req.RoomID, "房间已关闭", typedef_enum.ROOM_FORCE_OFF_MSG)
	}
	return err
}

// 申请房间
func (g *GuildRoom) RoomApply(c *gin.Context, req *request_room.GuildRoomApplyReq) (err error) {
	roomDao := new(dao.RoomDao)
	guildId := helper.GetGuildId(c)
	//房主信息验证
	userDao := new(dao.UserDao)
	roomUser, err := userDao.FindUserByUserNo(req.UserNo)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(roomUser.Id) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	//公会信息校验
	guildDao := new(dao2.GuildDao)
	guildInfo, err := guildDao.FindOne(&model.Guild{ID: guildId})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(guildInfo.ID) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeGuildNotExist,
			Msg:  nil,
		})
	}
	//房间上限验证
	roomCount, err := roomDao.GetGuildRoomCount(guildId)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if int(roomCount) >= guildInfo.RoomMax {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeGuildRoomMaxNum,
			Msg:  nil,
		})
	}

	//实名认证
	if roomUser.RealNameStatus != typedef_enum.UserRealNameAuthenticated {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeIDCardAuth,
			Msg:  nil,
		})
	}
	//检测是否是公会成员
	guildMemberDao := new(dao2.GuildMemberDao)
	data, err := guildMemberDao.FindOne(&model.GuildMember{
		GuildID: guildId,
		UserID:  roomUser.Id,
		Status:  1,
	})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if data.ID == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotInGuild,
		})
	}

	//日结算人信息验证
	dayUserInfo, err := userDao.FindUserByUserNo(req.DaySettleUserNo)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(dayUserInfo.Id) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if dayUserInfo.UserNo != roomUser.UserNo && dayUserInfo.Id != guildInfo.UserID {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeDaySettleInvalidRole,
			Msg:  nil,
		})
	}
	//月结算人信息验证
	MonthUserInfo, err := userDao.FindUserByUserNo(req.MonthSettleUserNo)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(MonthUserInfo.Id) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if MonthUserInfo.UserNo != roomUser.UserNo && MonthUserInfo.Id != guildInfo.UserID {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeMonthSettleInvalidRole,
			Msg:  nil,
		})
	}
	//房间类型验证
	roomTypeDao := new(dao2.RoomTypeDao)
	roomType, err := roomTypeDao.FindOne(&model.RoomType{TypeId: req.RoomType})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if roomType.ID == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomTypeNotExist,
			Msg:  nil,
		})
	}
	//房间模版验证
	roomTemplateDao := new(dao2.RoomTemplateDao)
	roomTemplate, err := roomTemplateDao.FindOne(&model.RoomTemplate{
		Id:     req.TemplateId,
		Status: 1,
	})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(roomTemplate.Id) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomTemplateNotExist,
			Msg:  nil,
		})
	}
	//添加房间申请
	roomApplyInfo := &model2.GuildRoomApply{
		GuildID:           guildId,
		RoomName:          req.RoomName,
		RoomID:            "0",
		RoomDesc:          req.RoomDesc,
		RoomUserID:        roomUser.Id,
		RoomAvatar:        req.CoverImg,
		Status:            typedef_enum.GuildRoomApplyStatusWaitReview,
		RoomType:          req.RoomType,
		TemplateID:        req.TemplateId,
		DaySettleUserID:   dayUserInfo.Id,
		MonthSettleUserID: MonthUserInfo.Id,
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
	}
	err = roomDao.AddRoomApply(roomApplyInfo)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	return err
}

func (g *GuildRoom) RoomApplyList(c *gin.Context, req *request_room.GuildRoomApplyListReq) (resp response.AdminPageRes) {
	list, count, err := new(dao2.RoomApplyDao).GetRoomApplyList(req, c)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	resp.Data = list
	resp.Total = count
	resp.CurrentPage = req.CurrentPage
	resp.Size = req.Size
	return
}
