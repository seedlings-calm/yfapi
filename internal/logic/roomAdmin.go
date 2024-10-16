package logic

import (
	"github.com/gin-gonic/gin"
	"time"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/acl"
	"yfapi/typedef/enum"
	resquest_room "yfapi/typedef/request/room"
	response_room "yfapi/typedef/response/room"
)

type RoomAdmin struct {
}

// RoomAdminAdd
//
//	@Description: 房间管理员添加
//	@receiver o
//	@param c
//	@param req
//	@return res
func (r *RoomAdmin) RoomAdminAdd(c *gin.Context, req *resquest_room.RoomAdminAddReq) {
	userId := helper.GetUserId(c)
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
	// 判断操作人是否在房间内
	if !new(acl.RoomAcl).IsInRoom(userId, req.RoomId, helper.GetClientType(c)) {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserNotInRoom,
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
	isExist := new(dao.GuildDao).GetCheckUserInGuild(req.UserId, room.GuildId)
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
	// 判断用户是否在房间内
	if !new(acl.RoomAcl).IsInRoom(req.UserId, req.RoomId, "") {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserNotInRoom,
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

// RoomAdminDelete
//
//	@Description: 房间管理员删除
//	@receiver r
//	@param c
//	@param req
func (r *RoomAdmin) RoomAdminDelete(c *gin.Context, req *resquest_room.RoomAdminDeleteReq) {
	userId := helper.GetUserId(c)
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
	// 判断操作人是否在房间内
	if !new(acl.RoomAcl).IsInRoom(userId, req.RoomId, helper.GetClientType(c)) {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserNotInRoom,
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
	//判断用户是否存在
	user, err := new(dao.UserDao).FindOne(&model.User{Id: req.UserId})
	if err != nil || user.Id == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	err = new(dao.RoomAdminDao).Delete(req.UserId, req.RoomId)
	return
}

// RoomAdminList
//
//	@Description: 房间管理员列表
//	@receiver r
//	@param c
//	@param req
//	@return res
func (r *RoomAdmin) RoomAdminList(c *gin.Context, req *resquest_room.RoomAdminListReq) (res *response_room.RoomAdminRes) {
	userId := helper.GetUserId(c)
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
	// 判断操作人是否在房间内
	if !new(acl.RoomAcl).IsInRoom(userId, req.RoomId, helper.GetClientType(c)) {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserNotInRoom,
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
	//查询房间管理员列表
	roomAdminList, err := new(dao.RoomAdminDao).List(req.RoomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	for _, v := range roomAdminList {
		res.RoomAdmin = append(res.RoomAdmin, response_room.RoomAdminInfo{
			UserId:     v.UserId,
			RoomId:     v.RoomId,
			RoomNo:     v.RoomNo,
			StaffName:  v.StaffName,
			CreateTime: v.CreateTime,
			UpdateTime: v.UpdateTime,
		})
	}
	return
}
