package room

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	dao2 "yfapi/internal/dao/roomowner"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	resquest_room "yfapi/typedef/request/room"
	request_login "yfapi/typedef/request/roomOwner"
	resquest_Practitioner "yfapi/typedef/request/roomOwner"
)

type RoomAdmin struct {
}

func (r *RoomAdmin) RoomAdminAdd(c *gin.Context, req *resquest_room.RoomAdminAddReq) {
	userId := helper.GetUserId(c)
	req.RoomId = helper.GetRoomId(c)
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: req.RoomId})
	if err != nil || room.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	//判断房间状态
	if room.Status != enum.RoomStatusNormal {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	// 判断操作人是否是房主
	if room.UserId != userId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotRoomOwner,
			Msg:  nil,
		})
	}
	//判断要添加的管理员和操作人即房主是否是一个工会
	isExist := new(dao.GuildDao).GetCheckUserInGuild(room.GuildId, req.UserId)
	if isExist == false {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotInGuild,
			Msg:  nil,
		})
	}
	//判断用户是否存在
	user, err := new(dao.UserDao).FindOne(&model.User{Id: req.UserId})
	if err != nil || user.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	//判断用户是否已是管理员
	isAdmin := new(dao.RoomAdminDao).IsAdmin(req.UserId, req.RoomId)
	if isAdmin == true {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserIsAdmin,
			Msg:  nil,
		})
	}
	//查询操作人昵称
	StaffInfo, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	roomAdmin := &model.RoomAdmin{
		UserId:     req.UserId,
		RoomId:     req.RoomId,
		RoomNo:     room.RoomNo,
		StaffName:  StaffInfo.Nickname,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err = new(dao.RoomAdminDao).Create(roomAdmin)
	return
}

func (r *RoomAdmin) RoomAdminRemove(c *gin.Context, req *request_login.RoomCommonReq) {
	userId := helper.GetUserId(c)
	roomId := helper.GetRoomId(c)
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	if err != nil || room.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	//判断房间状态
	if room.Status != enum.RoomStatusNormal {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	// 判断操作人是否是房主
	if room.UserId != userId {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeNotRoomOwner,
			Msg:  nil,
		})
	}
	err = new(dao.RoomAdminDao).Delete(req.UserId, roomId)
	return
}

func (r *RoomAdmin) RoomPractitionerAdd(c *gin.Context, req *resquest_Practitioner.RoomPractitionerAddReq) {
	roomId := helper.GetRoomId(c)
	opuserId := helper.GetUserId(c)
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	if err != nil || room.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	//判断房间状态
	if room.Status != enum.RoomStatusNormal {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	//判断用户是否存在
	user, err := new(dao.UserDao).FindOne(&model.User{Id: req.UserId})
	if err != nil || user.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeCheckUserID,
			Msg:  nil,
		})
	}
	// 判断用户是否拥有资质
	_, err = (&dao.DaoUserPractitionerCerd{UserId: req.UserId}).First(req.PractitionerType)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有从业者资质
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotPractitionerCredExist,
			Msg:  nil,
		})
	} else if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	//判断要添加的从业者是否是一个公会
	isExist := new(dao.GuildDao).GetCheckUserInGuild(room.GuildId, req.UserId)
	if isExist == false {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotGuildMember,
			Msg:  nil,
		})
	}
	//判断用户是否已是房间从业者
	var data model.UserPractitioner
	err = coreDb.GetSlaveDb().Model(data).Where("user_id=? and room_id=? and practitioner_type=? and status in ?", req.UserId, roomId, req.PractitionerType, []int{1, 2}).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if data.Id > 0 {
		if data.Status == 1 {
			// 已是从业者
			panic(error2.I18nError{
				Code: error2.ErrorCodeIsPractitionerExist,
				Msg:  nil,
			})
		} else {
			// 已申请从业者，审核中
			panic(error2.I18nError{
				Code: error2.ErrorCodePractitionerExamine,
				Msg:  nil,
			})
		}
	}
	//查询操作人昵称
	StaffInfo, err := new(dao.UserDao).FindOne(&model.User{Id: opuserId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	err = new(dao2.DaoRoomPractitioner).Create(&model.UserPractitioner{
		RoomId:            roomId,
		UserId:            req.UserId,
		PractitionerType:  req.PractitionerType,
		PractitionerBrief: req.PractitionerBrief,
		Status:            2,
		StaffName:         StaffInfo.Nickname,
		AbolishReason:     "",
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
	})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	return
}

func (r *RoomAdmin) RoomPractitionerRemove(c *gin.Context, req *resquest_Practitioner.RoomPractitionerAddReq) {
	roomId := helper.GetRoomId(c)
	opuserId := helper.GetUserId(c)
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	if err != nil || room.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	//判断房间状态
	if room.Status != enum.RoomStatusNormal {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	//判断用户是否存在
	user, err := new(dao.UserDao).FindOne(&model.User{Id: req.UserId})
	if err != nil || user.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeCheckUserID,
			Msg:  nil,
		})
	}
	//判断要移除的从业者是否是一个公会
	isExist := new(dao.GuildDao).GetCheckUserInGuild(room.GuildId, req.UserId)
	if isExist == false {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotGuildMember,
			Msg:  nil,
		})
	}
	//判断用户是否已是房间从业者
	var data model.UserPractitioner
	err = coreDb.GetSlaveDb().Model(data).Where("user_id=? and room_id=? and practitioner_type=? and status=1", req.UserId, roomId, req.PractitionerType).First(&data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有从业者资质
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotRoomPractitioner,
			Msg:  nil,
		})
	} else if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}

	//查询操作人昵称
	StaffInfo, err := new(dao.UserDao).FindOne(&model.User{Id: opuserId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	tx := coreDb.GetMasterDb().Begin()
	// 取消本房间从业者身份
	err = tx.Model(model.UserPractitioner{}).Where("user_id=? and room_id=? and practitioner_type=? and status=1", req.UserId, roomId, req.PractitionerType).Updates(map[string]interface{}{"status": "4", "staff_name": StaffInfo.Nickname}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 取消本房间从业者权限
	err = tx.Model(model.AuthRoleAccess{}).Where("user_id = ? and room_id = ?", req.UserId, roomId).Delete(&model.AuthRoleAccess{}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	tx.Commit()
	key1 := redisKey.UserRules(req.UserId, roomId)
	key2 := redisKey.UserRoles(req.UserId, roomId)
	key3 := redisKey.UserCompereRules(req.UserId, roomId)
	coreRedis.GetUserRedis().Del(c, key1, key2, key3)
	return
}

// 从业者后台再次提交
func (r *RoomAdmin) RoomPractitionerReSave(c *gin.Context, req *resquest_Practitioner.RoomPractitionerUpdateReq) {
	//userId := req.UserId
	opUserId := helper.GetUserId(c)
	roomId := helper.GetRoomId(c)
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	if err != nil || room.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	//判断房间状态
	if room.Status != enum.RoomStatusNormal {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	// 数据状态
	var info model.UserPractitioner
	err = coreDb.GetSlaveDb().Model(model.UserPractitioner{}).Where("id", req.Id).Scan(&info).Error
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if info.RoomId != roomId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotGuildMember,
			Msg:  nil,
		})
	}
	if info.Status != 3 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}

	//判断用户是否存在
	user, err := new(dao.UserDao).FindOne(&model.User{Id: info.UserId})
	if err != nil || user.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeCheckUserID,
			Msg:  nil,
		})
	}
	//判断是否是一个公会会
	isExist := new(dao.GuildDao).GetCheckUserInGuild(room.GuildId, info.UserId)
	if isExist == false {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotGuildMember,
			Msg:  nil,
		})
	}
	//判断用户是否已是房间从业者
	var data model.UserPractitioner
	err = coreDb.GetSlaveDb().Model(data).Where("user_id=? and room_id=? and practitioner_type=? and status in ?", info.UserId, roomId, info.PractitionerType, []int{1, 2}).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if data.Id > 0 {
		if data.Status == 1 {
			// 已是从业者
			panic(error2.I18nError{
				Code: error2.ErrorCodeIsPractitionerExist,
				Msg:  nil,
			})
		} else {
			// 已申请从业者，审核中
			panic(error2.I18nError{
				Code: error2.ErrorCodePractitionerExamine,
				Msg:  nil,
			})
		}
	}
	//查询操作人昵称
	StaffInfo, err := new(dao.UserDao).FindOne(&model.User{Id: opUserId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	info.Status = 2
	info.StaffName = StaffInfo.Nickname
	info.UpdateTime = time.Now()
	info.PractitionerBrief = req.PractitionerBrief
	err = new(dao2.DaoRoomPractitioner).Update(&info)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	return
}

// RoomPractitionerInvalid
//
//	@Description: 从业者申请作废
//	@receiver r
//	@param c *gin.Context -
//	@param req *request_login.RoomCommonReq -
func (r *RoomAdmin) RoomPractitionerInvalid(c *gin.Context, req *request_login.RoomPractitionerUpdateReq) {
	roomId := helper.GetRoomId(c)
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	if err != nil || room.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	// 数据状态
	var info model.UserPractitioner
	err = coreDb.GetSlaveDb().Model(model.UserPractitioner{}).Where("id", req.Id).Scan(&info).Error
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if info.RoomId != roomId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotGuildMember,
			Msg:  nil,
		})
	}
	if info.Status != 3 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	err = coreDb.GetMasterDb().Where("id", req.Id).Delete(model.UserPractitioner{}).Error
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	return
}

// 审核和拒绝
func (r *RoomAdmin) RoomPractitionerApply(c *gin.Context, req *request_login.RoomPractitionerApply) {
	userId := helper.GetUserId(c)
	roomid := c.GetHeader("roomId")
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomid})
	if err != nil || room.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	// 判断操作人是否是房主
	if room.UserId != userId {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeNotRoomOwner,
			Msg:  nil,
		})
	}
	//判断用户是否存在
	user, err := new(dao.UserDao).FindOne(&model.User{Id: req.UserId})
	if err != nil || user.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	data := &model.UserPractitioner{
		RoomId:        roomid,
		UserId:        req.UserId,
		Status:        req.Status,
		StaffName:     user.Nickname,
		AbolishReason: req.AbolishReason,
		UpdateTime:    time.Now(),
	}
	err = new(dao2.DaoRoomPractitioner).Update(data)
	return
}
