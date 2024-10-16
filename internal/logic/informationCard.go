package logic

import (
	"log"
	"strconv"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreConfig"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/acl"
	"yfapi/internal/service/auth"
	service_im "yfapi/internal/service/im"
	service_room "yfapi/internal/service/room"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_room "yfapi/typedef/request/room"
	resquest_user "yfapi/typedef/request/user"
	response_room "yfapi/typedef/response/room"

	ginI18n "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (a *ActionRoom) GetInformationCardByUserId(c *gin.Context, req *resquest_user.UserPractitionerReq) (resp response_room.InformationCardResponse) {
	// res, _ := coreRedis.GetChatroomRedis().Get(c, redisKey.InformationCardKey(req.UserId, req.RoomId)).Result()
	// if res != "" {
	// 	err := json.Unmarshal([]byte(res), &resp)
	// 	if err == nil {
	// 		return
	// 	}
	// }
	userDao := &dao.UserDao{}
	userInfo, err := userDao.FindOne(&model.User{Id: req.UserId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	resp.Avatar = coreConfig.GetHotConf().ImagePrefix + userInfo.Avatar
	resp.Nickname = userInfo.Nickname
	resp.IsTrueName = userInfo.RealNameStatus == enum.UserRealNameAuthenticated
	resp.Sex = userInfo.Sex
	resp.UserNo = userInfo.UserNo
	resp.UserId = req.UserId
	resp.Introduce = userInfo.Introduce
	ownerId := handle.GetUserId(c)
	followDao := &dao.UserFollowDao{}
	resF, _ := followDao.GetUserFollow(ownerId, req.UserId)
	if resF.Id == 0 {
		resp.IsFollow = 0
	} else {
		resp.IsFollow = 1
		if resF.IsMutualFollow {
			resp.IsFollow = 2
		}
	}
	resp.FollowedNum = followDao.GetUserFollowedNum(req.UserId)
	resp.FansNum = followDao.GetUserFansNum(req.UserId)
	PractitionerCerdDao := &dao.DaoUserPractitionerCerd{UserId: req.UserId}
	pras, _ := PractitionerCerdDao.Find()
	if len(pras) >= 1 {
		for _, v := range pras {
			resp.Practitions = append(resp.Practitions, v.PractitionerType)
		}
	}

	//公会
	guildDao := &dao.GuildDao{}
	resp.GuildName = guildDao.GetGuildNameByUserId(req.UserId)

	// 当前房间的身份
	resp.RoleIdList, _, _ = new(auth.Auth).GetRoleListByRoomIdAndUserId(req.RoomId, req.UserId)
	if len(resp.RoleIdList) == 0 {
		resp.RoleIdList = []int{enum.NormalRoleId}
	}

	resp.IsOnline = service_room.IsUserInRoom(req.RoomId, req.UserId)
	// 查询用户的铭牌信息
	resp.UserPlaque = service_user.GetUserLevelPlaque(req.UserId, helper.GetClientType(c))
	// bytes, err := json.Marshal(resp)
	// if err == nil {
	// 	coreRedis.GetChatroomRedis().Set(c, redisKey.InformationCardKey(req.UserId, req.RoomId), bytes, 1*time.Minute)
	// }
	return
}

// 执行拉黑相关操作
func (a *ActionRoom) DoBlackout(c *gin.Context, toId string, roomId string, command string, content string) {
	ownerId := handle.GetUserId(c)
	if command == acl.AddRoomBlacklist { //拉黑，取消拉黑
		if content == "del" {
			a.DelBlackOut(c, ownerId, toId, roomId)
		} else {
			a.AddBlackout(c, ownerId, toId, roomId)
		}
	}
}

// 拉黑用户
func (a *ActionRoom) AddBlackout(c *gin.Context, fromId, toId string, roomId string) {
	blackDao := dao.UserBlackListDao{}
	err := blackDao.Create(&model.UserBlacklist{FromID: fromId, ToID: toId, RoomID: roomId, Types: 1, IsEffective: true})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}

	//退出房间逻辑处理
	serviceUser := service_room.RoomUsersOnlie{RoomId: roomId}
	roomInfo, _ := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	serviceUser.RemoveUserToRoom(c, toId, "", roomInfo)

	tokenData := handle.GetTokenData(c)
	//通知im离开房间
	new(service_im.ImPublicService).SendActionMsg(c, map[string]string{
		"content": ginI18n.MustGetMessage(
			c,
			&i18n.LocalizeConfig{
				MessageID: strconv.Itoa(error2.ErrCodeBlackErr),
				TemplateData: map[string]interface{}{
					"roomName": roomInfo.Name,
				},
			}),
	}, toId, "", roomId, tokenData.ClientType, enum.BLACKOUT_ROOM_MSG)
}
func (a *ActionRoom) BlackList(c *gin.Context, roomId string) (resp response_room.BlackListResponse) {
	blackDao := dao.UserBlackListDao{}
	blackResp := blackDao.GetListByRoomId(roomId)
	resp.Count = len(blackResp)
	resp.List = make([]response_room.BlackListAndUserInfo, 0)
	if resp.Count > 0 {
		for _, v := range blackResp {
			v.Avatar = helper.FormatImgUrl(v.Avatar)
			// 查询用户的铭牌信息
			v.UserPlaque = service_user.GetUserLevelPlaque(v.UserId, helper.GetClientType(c))
			resp.List = append(resp.List, v)
		}
	}
	return
}

// 取消拉黑操作
func (a *ActionRoom) DelBlackOut(c *gin.Context, fromId, toId string, roomId string) {
	blackDao := dao.UserBlackListDao{}
	err := blackDao.Update(&model.UserBlacklist{UnsealID: fromId, ToID: toId, RoomID: roomId, Types: 1, IsEffective: false})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeRoomPermissionDenied,
			Msg:  nil,
		})
	}
}

// 踢出房间
func (a *ActionRoom) KickOut(c *gin.Context, req *request_room.KickOutReq) {
	var timeArr = map[string]time.Duration{
		"1":  1,
		"5":  5,
		"10": 10,
		"30": 30,
	}
	if _, ok := timeArr[req.Times]; !ok {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}

	err := coreRedis.GetChatroomRedis().Set(c, redisKey.RoomKickOutKey(req.UserId, req.RoomId), timeArr[req.Times], timeArr[req.Times]*time.Minute).Err()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeKickOutErr,
			Msg:  nil,
		})
	}
	//退出房间逻辑处理
	serviceUser := service_room.RoomUsersOnlie{RoomId: req.RoomId}
	err = serviceUser.RemoveUserToRoom(c, req.UserId, "", nil)
	log.Println("执行踢出，退出房间，踢出结构体:", req)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeKickOutErr,
			Msg:  nil,
		})
	}
	tokenData := handle.GetTokenData(c)
	new(service_im.ImPublicService).SendActionMsg(c, map[string]string{
		"content": ginI18n.MustGetMessage(
			c,
			&i18n.LocalizeConfig{
				MessageID: strconv.Itoa(error2.ErrCodeBlackOutImErr),
				TemplateData: map[string]interface{}{
					"times": req.Times,
				},
			}),
	}, req.UserId, "", req.RoomId, tokenData.ClientType, enum.KICKOUT_ROOM_MSG)
}
