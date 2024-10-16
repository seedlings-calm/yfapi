package logic

import (
	"context"
	"database/sql"
	"time"
	"unicode/utf8"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/auth"
	service_goods "yfapi/internal/service/goods"
	service_im "yfapi/internal/service/im"
	"yfapi/internal/service/riskCheck/shumei"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	common_data "yfapi/typedef/redisKey"
	request_user "yfapi/typedef/request/user"
	response_user "yfapi/typedef/response/user"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

type UserInfo struct {
}

// PerfectInfo 完善用户信息
func (ser *UserInfo) PerfectInfo(req *request_user.PerfectInfoReq, c *gin.Context) {
	userId := helper.GetUserId(c)
	success, unlock, err := coreRedis.UserLock(context.Background(), common_data.UserModifyInfoLock(userId), time.Second*5)
	if err != nil || !success {
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFrequent,
			Msg:  nil,
		})
	}
	defer unlock()
	if req.Sex == 0 {
		req.Sex = enum.UserSexTypeWoman
	}
	if len(req.Nickname) > 0 {
		if ok := new(shumei.ShuMei).NicknameCheck(userId, req.Nickname); !ok {
			req.Nickname = "用户" + helper.GeneUserNickname()
			new(service_im.ImNoticeService).SendSystematicMsg(c, i18n_msg.GetI18nMsg(c, i18n_msg.NicknameRejectKey), "", i18n_msg.GetI18nMsg(c, i18n_msg.NicknameRejectContextKey), "", "", []string{userId})
		}
	}
	if len(req.Avatar) > 0 {
		if ok := new(shumei.ShuMei).AvatarSyncCheck(userId, helper.FormatImgUrl(req.Avatar)); !ok {
			req.Avatar = helper.GetUserDefaultAvatar(req.Sex)
			new(service_im.ImNoticeService).SendSystematicMsg(c, i18n_msg.GetI18nMsg(c, i18n_msg.AvatarRejectKey), "", i18n_msg.GetI18nMsg(c, i18n_msg.AvatarRejectContextKey), "", "", []string{userId})
		}
	}
	user, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if len(req.Nickname) > 0 {
		if !new(dao.UserDao).CheckRepeatNickname(req.Nickname) {
			user.Nickname = req.Nickname
		} else {
			user.Nickname = "用户" + helper.GeneUserNickname()
		}
	}
	user.Sex = req.Sex
	if len(req.BornDate) != 0 {
		user.BornDate = sql.NullString{
			String: req.BornDate,
			Valid:  true,
		}
	}
	user.Avatar = req.Avatar
	userDao := new(dao.UserDao)
	user.Guide = 2
	err = userDao.UpdateById(user)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
		})
	}
}

// 获取用户分配得im服务
func (ser *UserInfo) ImServer(c *gin.Context) (res response_user.UserImRes) {
	userId := helper.GetUserId(c)
	ok, unlock, err := coreRedis.ImLock(context.Background(), common_data.UserImServerLock(userId), time.Second*5)
	if err != nil || !ok {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserImServer,
			Msg:  nil,
		})
	}
	defer unlock()
	imServer, _ := coreRedis.GetImRedis().Get(context.Background(), common_data.UserImServer(userId)).Result()
	_, err = coreRedis.GetImRedis().ZScore(context.Background(), common_data.OnlineImServer(), imServer).Result()
	if len(imServer) != 0 && err != redis.Nil {
		//检测im地址是否可访问
		connection, err := helper.CheckWebSocketConnection(helper.WsFull(imServer))
		if err != nil || !connection {
			coreLog.Error("im地址不可用 address:%s err:%+v", imServer, err)
			failCount, _ := coreRedis.GetImRedis().Get(context.Background(), common_data.ImServerConnectFailCount(imServer)).Result()
			if cast.ToInt(failCount) > 5 {
				coreRedis.GetImRedis().ZRem(context.Background(), common_data.OnlineImServer(), imServer)
			}
			coreRedis.GetImRedis().Incr(context.Background(), common_data.ImServerConnectFailCount(imServer))
			coreRedis.GetImRedis().Expire(context.Background(), common_data.ImServerConnectFailCount(imServer), time.Second*60)
		}
		res.Service = helper.WsFull(imServer)
	} else {
		key := common_data.OnlineImServer()
		result, err := coreRedis.GetImRedis().ZRangeByScoreWithScores(context.Background(), key, &redis.ZRangeBy{
			Min:    "-inf",
			Max:    "inf",
			Offset: 0,
			Count:  1,
		}).Result()
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserImServer,
				Msg:  nil,
			})
		}
		_, err = coreRedis.GetImRedis().Set(context.Background(), common_data.UserImServer(userId), cast.ToString(result[0].Member), enum.UserConnectImServiceLife).Result()
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserImServer,
				Msg:  nil,
			})
		}
		if len(result) > 0 {
			res.Service = helper.WsFull(cast.ToString(result[0].Member))
		}
	}
	return
}

// EditUserInfo
//
//	@Description: 更新用户信息
//	@receiver ser
//	@param req
//	@param context
func (ser *UserInfo) EditUserInfo(req *request_user.EditUserInfoReq, c *gin.Context) (resp response_user.UserInfo) {
	if utf8.RuneCountInString(req.Nickname) > 16 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	userId := helper.GetUserId(c)
	success, unlock, err := coreRedis.UserLock(context.Background(), common_data.UserModifyInfoLock(userId), time.Second*5)
	if err != nil || !success {
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFrequent,
			Msg:  nil,
		})
	}
	defer unlock()
	user, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	rds := coreRedis.GetUserRedis()
	data := map[string]any{}
	if len(req.Avatar) != 0 && req.Avatar != helper.FormatImgUrl(user.Avatar) {
		req.Avatar = helper.RemovePrefixImgUrl(req.Avatar)
		if ok := new(shumei.ShuMei).AvatarSyncCheck(userId, helper.FormatImgUrl(req.Avatar)); !ok {
			panic(error2.I18nError{
				Code: error2.ErrorCodeAvatarCheckReject,
				Msg:  nil,
			})
		}
		data["avatar"] = req.Avatar
	}
	if len(req.Nickname) != 0 && req.Nickname != user.Nickname {
		nicknameEditNum, _ := rds.Get(context.Background(), common_data.UserNicknameEditNum(userId)).Result()
		if cast.ToInt(nicknameEditNum) >= 3 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeNicknameEditNum,
				Msg:  nil,
			})
		}
		if ok := new(shumei.ShuMei).NicknameCheck(userId, req.Nickname); !ok {
			panic(error2.I18nError{
				Code: error2.ErrorCodeNicknameCheckReject,
				Msg:  nil,
			})
		}
		if new(dao.UserDao).CheckRepeatNickname(req.Nickname) {
			panic(error2.I18nError{
				Code: error2.ErrorCodeNicknameRepeat,
				Msg:  nil,
			})
		}
		data["nickname"] = req.Nickname
	}
	if req.Sex != user.Sex {
		if user.SexEditNum > 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSexEditNum,
				Msg:  nil,
			})
		}
		data["sex"] = req.Sex
		data["sex_edit_num"] = user.SexEditNum + 1
	}
	if len(req.BornDate) != 0 && req.BornDate != user.BornDate.String {
		data["born_date"] = req.BornDate
	}
	if len(req.Introduce) != 0 && req.Introduce != user.Introduce {
		if ok := new(shumei.ShuMei).SignCheck(userId, req.Introduce); !ok {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSignCheckReject,
				Msg:  nil,
			})
		}
		data["introduce"] = req.Introduce
	}
	if req.Voice != nil && req.Voice.Url != "" && req.Voice.Url != helper.FormatImgUrl(user.VoiceUrl) {
		req.Voice.Url = helper.RemovePrefixImgUrl(req.Voice.Url)
		voiceEditNum, _ := rds.Get(context.Background(), common_data.UserVoiceEditNum(userId)).Result()
		if cast.ToInt(voiceEditNum) >= 3 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeVoiceEditNum,
				Msg:  nil,
			})
		}
		data["voice_url"] = req.Voice.Url
		data["voice_length"] = req.Voice.Length
		data["voice_status"] = 1
	} else {
		data["voice_url"] = ""
		data["voice_length"] = 0
		data["voice_status"] = 1
	}
	userDao := new(dao.UserDao)
	err = userDao.UpdateUserFieldsByUserId(userId, data)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	if len(req.Nickname) != 0 {
		rds.Incr(context.Background(), common_data.UserNicknameEditNum(userId))
		rds.Expire(context.Background(), common_data.UserNicknameEditNum(userId), time.Hour*24*31)
	}
	if req.Voice != nil && req.Voice.Url != "" && req.Voice.Url != helper.FormatImgUrl(user.VoiceUrl) {
		rds.Incr(context.Background(), common_data.UserVoiceEditNum(userId))
		rds.Expire(context.Background(), common_data.UserVoiceEditNum(userId), time.Hour*24*31)
		//异步检测声音签名
		new(shumei.ShuMei).SignAudioAsyncCheck(userId, helper.FormatImgUrl(req.Voice.Url), userId, coreSnowflake.GetSnowId())
	}
	resp = ser.GetUserInfo(&request_user.UserInfoReq{UserId: userId}, c)
	return
}

// 获取用户信息
func (ser *UserInfo) GetUserInfo(req *request_user.UserInfoReq, c *gin.Context) (resp response_user.UserInfo) {
	userId := helper.GetUserId(c)
	data, err := new(dao.UserDao).FindOne(&model.User{Id: req.UserId, Status: 1})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	//查询身份证号
	realNameInfo, err := new(dao.UserDao).FindUserIdNo(req.UserId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
		})
	}
	if realNameInfo.IdNo != "" {
		resp.IdNo = realNameInfo.IdNo
	}
	resp.UserId = data.Id
	resp.UserNo = data.UserNo
	resp.Uid32 = cast.ToInt32(data.OriUserNo)
	resp.Sex = data.Sex
	resp.Avatar = helper.FormatImgUrl(data.Avatar)
	resp.Nickname = data.Nickname
	resp.Introduce = data.Introduce
	resp.TrueName = data.TrueName
	resp.VoiceUrl = helper.FormatImgUrl(data.VoiceUrl)
	if data.Password != "" {
		resp.IsSetPwd = true
	}
	resp.VoiceLength = data.VoiceLength
	if len(data.BornDate.String) != 0 {
		parse, _ := time.Parse(time.RFC3339, data.BornDate.String)
		resp.BornDate = parse.Format(time.DateOnly)
	}
	//增加用户逇头饰装扮
	resp.Headwear = service_goods.UserGoods{}.GetGoodsByKey(req.UserId, true, enum.GoodsTypeTS)

	if userId == req.UserId {
		friendNum := new(dao.UserFollowDao).GetUserFriendsNum(userId)
		//查看自己信息
		resp.FriendNum = int(friendNum)
		followedNum := new(dao.UserFollowDao).GetUserFollowedNum(userId)
		resp.FollowedNum = int(followedNum)
		fansNum := new(dao.UserFollowDao).GetUserFansNum(userId)
		resp.FansNum = int(fansNum)
		resp.IsOnline = service_user.IsOnline(userId)
		resp.RegionCode = data.RegionCode
		resp.Mobile = data.Mobile
		// 是否有主播资质
		cerdInfo, _ := (&dao.DaoUserPractitionerCerd{UserId: userId}).First(enum.UserPractitionerAnchor)
		if cerdInfo.Id > 0 {
			resp.IsAnchor = true
		}
		// 公会信息
		resp.GuildName = new(dao.GuildDao).GetGuildNameByUserId(userId)
		// lv经验信息
		userLv, _ := new(dao.UserLevelLvDao).GetUserLvLevelDTO(userId)
		if userLv.ID > 0 {
			resp.LvCurrExp = userLv.CurrExp
			resp.LvMinExp = userLv.MinExperience
			resp.LvMaxExp = userLv.MaxExperience
		}
		// 访客数量
		resp.VisitorNum = new(dao.UserVisitDao).GetVisitUserCount(userId)
		resp.RoleIdList, _, _ = new(auth.Auth).GetRoleListByRoomIdAndUserId(enum.WorldGroupId, userId)
		if len(resp.RoleIdList) == 0 {
			resp.RoleIdList = []int{enum.NormalRoleId}
		}
	} else { //查看他人信息
		followedNum := new(dao.UserFollowDao).GetUserFollowedNum(req.UserId)
		resp.FollowedNum = int(followedNum)
		fansNum := new(dao.UserFollowDao).GetUserFansNum(req.UserId)
		resp.FansNum = int(fansNum)
		resp.LikeNum = service_user.GetUserTotalPraisedCount(req.UserId)
		follow, _ := new(dao.UserFollowDao).GetUserFollow(userId, req.UserId)
		if follow.Id == 0 {
			resp.FollowedType = 0
			otherFollow, _ := new(dao.UserFollowDao).GetUserFollow(req.UserId, userId)
			if otherFollow.Id != 0 {
				resp.FollowedType = 3
			}
		} else {
			resp.FollowedType = 1
			if follow.IsMutualFollow {
				resp.FollowedType = 2
			}
		}
		resp.IsOnline = service_user.IsOnline(req.UserId)
		if resp.IsOnline {
			// 当前玩家是否在房
			resp.InRoom = service_user.GetUserInRoomInfo(req.UserId, helper.GetClientType(c))
		}
		//判断是否拉黑这个用户
		resp.IsBlacklist = service_user.IsBlacklist(userId, req.UserId, "0", enum.BlacklistTypeUser)
		resp.DontLetHeSeeMoments = new(dao.UserTimelineFilterDao).GetSwitchType(userId, req.UserId, enum.DontLetHeSeeMoments)
		resp.DontSeeHeMoments = new(dao.UserTimelineFilterDao).GetSwitchType(userId, req.UserId, enum.DontSeeHeMoments)
		resp.MomentsNoticeSwitch = new(dao.UserNoticeFilterDao).GetSwitchType(userId, req.UserId, enum.MomentsNoticeSwitch)
		resp.LiveNoticeSwitch = new(dao.UserNoticeFilterDao).GetSwitchType(userId, req.UserId, enum.LiveNoticeSwitch)
		// 记录用户访问
		go service_user.RecordUserVisit(userId, req.UserId)
	}
	// 查询用户的铭牌信息
	resp.UserPlaque = service_user.GetUserLevelPlaque(req.UserId, helper.GetClientType(c))
	return
}

// SearchUserInfo
//
//	@Description: 模糊查询符合的user_no或nickname
//	@receiver u
//	@param c *gin.Context -
//	@param keyword string -
//	@return res -
func (u *UserInfo) SearchUserInfo(c *gin.Context, keyword string) (res []response_user.UserInfo) {
	userId := helper.GetUserId(c)
	userList, err := new(dao.UserDao).FindUserByKeyword(keyword, 0, 50)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	for _, data := range userList {
		resp := response_user.UserInfo{}
		resp.UserId = data.Id
		resp.UserNo = data.UserNo
		resp.Uid32 = cast.ToInt32(data.OriUserNo)
		resp.Sex = data.Sex
		resp.Avatar = helper.FormatImgUrl(data.Avatar)
		resp.Nickname = data.Nickname
		resp.Introduce = data.Introduce
		resp.VoiceUrl = helper.FormatImgUrl(data.VoiceUrl)
		resp.VoiceLength = data.VoiceLength
		resp.Headwear = service_goods.UserGoods{}.GetGoodsByKey(data.Id, true, enum.GoodsTypeTS)
		if userId == data.Id {
			//查看自己信息
			friendNum := new(dao.UserFollowDao).GetUserFriendsNum(userId)
			//查看自己信息
			resp.FriendNum = int(friendNum)
			followedNum := new(dao.UserFollowDao).GetUserFollowedNum(userId)
			resp.FollowedNum = int(followedNum)
			fansNum := new(dao.UserFollowDao).GetUserFansNum(userId)
			resp.FansNum = int(fansNum)
			resp.IsOnline = service_user.IsOnline(userId)
		} else { //查看他人信息
			followedNum := new(dao.UserFollowDao).GetUserFollowedNum(data.Id)
			resp.FollowedNum = int(followedNum)
			fansNum := new(dao.UserFollowDao).GetUserFansNum(data.Id)
			resp.FansNum = int(fansNum)
			resp.LikeNum = service_user.GetUserTotalPraisedCount(data.Id)
			follow, _ := new(dao.UserFollowDao).GetUserFollow(userId, data.Id)
			if follow.Id == 0 {
				resp.FollowedType = 0
			} else {
				resp.FollowedType = 1
				if follow.IsMutualFollow {
					resp.FollowedType = 2
				}
			}
			resp.IsOnline = service_user.IsOnline(data.Id)
		}
		res = append(res, resp)
	}
	return
}

func (u *UserInfo) GetUserBasicInfo(userId string) response_user.UserH5BasicInfo {
	users, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	resp := response_user.UserH5BasicInfo{
		UserId:   users.Id,
		Nickname: users.Nickname,
		UserNo:   users.UserNo,
		Uid32:    cast.ToInt32(users.OriUserNo),
		Avatar:   coreConfig.GetHotConf().ImagePrefix + users.Avatar,
		Headwear: service_goods.UserGoods{}.GetGoodsByKey(users.Id, true, enum.GoodsTypeTS),
	}
	return resp
}

// 根据用户userNo查询用户信息
func (u *UserInfo) SearchUserInfoByUserNo(c *gin.Context, userNo string) (res response_user.SearchUserByUserNoResp) {
	if len(userNo) == 0 {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	userModel, err := new(dao.UserDao).FindUserByUserNo(userNo)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	res.Id = userModel.Id
	res.UserNo = userModel.UserNo
	res.Nickname = userModel.Nickname
	res.Avatar = helper.FormatImgUrl(userModel.Avatar)
	return
}

func (u *UserInfo) GerUserRealNameInfo(c *gin.Context) (res response_user.UserRealNameInfoResp) {
	userId := helper.GetUserId(c)
	one := new(dao.UserRealNameDao).FindOne(&model.UserRealName{
		UserId: userId,
	})
	if one.Id == 0 {
		res.RealNameStatus = enum.UserRealName(enum.UserRealNameUnverified).String()
	} else {
		if one.Status == enum.UserRealNameReviewStatusWait {
			res.RealNameStatus = enum.UserRealNameReviewStatus(enum.UserRealNameReviewStatusWait).String()
		} else if one.Status == enum.UserRealNameReviewStatusReject {
			res.RealNameStatus = enum.UserRealName(enum.UserRealNameUnverified).String()
		} else {
			res.RealNameStatus = enum.UserRealName(enum.UserRealNameAuthenticated).String()
			res.RealName = helper.PrivateRealName(one.TrueName)
			res.CardNum = helper.PrivateIdNo(one.IdNo)
		}
	}
	return
}
