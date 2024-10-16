package user

import (
	"context"
	"yfapi/core/coreRedis"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_goods "yfapi/internal/service/goods"
	typedef_enum "yfapi/typedef/enum"
	common_userInfo "yfapi/typedef/redisKey"
	"yfapi/typedef/response/user"

	"github.com/spf13/cast"
)

func GetUserInfo(targetUserId, clientType string, userId ...string) (result user.UserInfo) {
	if len(targetUserId) == 0 {
		return
	}
	userDao := new(dao.UserDao)
	userInfo, _ := userDao.FindOne(&model.User{
		Id: targetUserId,
	})
	result.UserId = targetUserId
	result.Nickname = userInfo.Nickname
	result.Avatar = helper.FormatImgUrl(userInfo.Avatar)
	result.UserNo = userInfo.UserNo
	result.Sex = userInfo.Sex
	result.LikeNum = GetUserTotalPraisedCount(targetUserId)
	result.IsOnline = IsOnline(targetUserId)
	if len(userId) > 0 {
		follow, _ := new(dao.UserFollowDao).GetUserFollow(userId[0], targetUserId)
		if follow.Id == 0 {
			result.FollowedType = 0
		} else {
			result.FollowedType = 1
			if follow.IsMutualFollow {
				result.FollowedType = 2
			}
		}
	}
	result.Headwear = service_goods.UserGoods{}.GetGoodsByKey(targetUserId, true, typedef_enum.GoodsTypeTS)
	// 查询用户的铭牌信息
	// true彩色昵称 true房间身份
	result.UserPlaque = GetUserLevelPlaque(targetUserId, clientType, true, true)
	return
}

// 判断用户是否在线
func IsOnline(userId string) bool {
	key := common_userInfo.ImOnlineUser(userId)
	result, err := coreRedis.GetImRedis().Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if result == 1 {
		return true
	}
	return false
}

// IsBlacklist 是否拉黑 用户拉黑，roomId传字符串0， 如果处于拉黑中，返回true
func IsBlacklist(fromId, toId, roomId string, typeId int) bool {
	blacklistDao := dao.UserBlackListDao{}
	models := &model.UserBlacklist{
		RoomID:      roomId,
		ToID:        toId,
		IsEffective: true,
		Types:       typeId,
	}
	if typeId == typedef_enum.BlacklistTypeUser { //如果是拉黑用户，需要指定拉黑操作人
		models.FromID = fromId
	}
	return blacklistDao.IsLog(models)
}

// GetUserTotalPraisedCount
//
//	@Description: 查询用户获得的总点赞数量
//	@param userId string -
//	@return count -
func GetUserTotalPraisedCount(userId string) (count int) {
	timelineCount, _ := new(dao.TimelineDao).GetUserTimelineLoveCount(userId)
	replyCount, _ := new(dao.TimelineReplyDao).GetUserTimelineReplyPraisedCount(userId)
	return timelineCount + replyCount
}

// 获取会话列表用户信息
func GetSessionListUserInfo(userId, targetUserId, clientType string) (result user.SessionListUserInfo) {
	if len(targetUserId) == 0 {
		return
	}
	userDao := new(dao.UserDao)
	userInfo, _ := userDao.FindOne(&model.User{
		Id: targetUserId,
	})
	result.UserId = targetUserId
	result.Nickname = userInfo.Nickname
	result.Avatar = helper.FormatImgUrl(userInfo.Avatar)
	result.IsOnline = IsOnline(targetUserId)
	result.UserNo = userInfo.UserNo
	result.Uid32 = cast.ToInt32(userInfo.OriUserNo)
	result.Sex = userInfo.Sex
	result.IsBlacklist = IsBlacklist(userId, targetUserId, "0", typedef_enum.BlacklistTypeUser)
	if result.IsOnline {
		// 当前玩家是否在房
		result.InRoom = GetUserInRoomInfo(targetUserId, clientType)
	}
	result.UserPlaque = GetUserLevelPlaque(targetUserId, clientType)
	follow, _ := new(dao.UserFollowDao).GetUserFollow(userId, targetUserId)
	if follow.Id == 0 {
		result.FollowedType = 0
		otherFollow, _ := new(dao.UserFollowDao).GetUserFollow(targetUserId, userId)
		if otherFollow.Id != 0 {
			result.FollowedType = 3
		}
	} else {
		result.FollowedType = 1
		if follow.IsMutualFollow {
			result.FollowedType = 2
		}
	}
	return
}

func GetUserBaseInfo(userId string) (result *model.User) {
	if len(userId) == 0 {
		return
	}
	userDao := new(dao.UserDao)
	result, _ = userDao.FindOne(&model.User{
		Id: userId,
	})
	result.Avatar = helper.FormatImgUrl(result.Avatar)
	result.VoiceUrl = helper.FormatImgUrl(result.VoiceUrl)
	return
}

func GetUserBaseInfoList(userIdList []string) (result []model.User) {
	if len(userIdList) == 0 {
		return
	}
	result = new(dao.UserDao).FindByIds(userIdList)
	for i := range result {
		result[i].Avatar = helper.FormatImgUrl(result[i].Avatar)
		result[i].VoiceUrl = helper.FormatImgUrl(result[i].VoiceUrl)
	}
	return
}

func GetUserBaseInfoMap(userIdList []string) (result map[string]model.User) {
	result = make(map[string]model.User)
	data := GetUserBaseInfoList(userIdList)
	for _, info := range data {
		result[info.Id] = info
	}
	return
}

// IsRoomPractitioner 是否为房间从业者
func IsRoomPractitioner(userId, roomId string) bool {
	dataList, err := new(dao.DaoUserPractitioner).Find(userId, roomId)
	if err != nil {
		return false
	}
	if len(dataList) > 0 {
		return true
	}
	return false
}

// GetUserInRoomInfo 查询玩家在房信息
func GetUserInRoomInfo(userId, clientType string) (result *user.RoomInfo) {
	// 是否为在房从业者[有从业者资格，加入公会的从业者（主播有从业者资格即可）]
	dataList, err := (&dao.DaoUserPractitionerCerd{UserId: userId}).Find()
	if err != nil {
		return nil
	}
	isOk := false
	for _, data := range dataList {
		// 主播资质
		if data.PractitionerType == typedef_enum.UserPractitionerAnchor {
			isOk = true
			break
		}
	}
	// 拥有其他资质
	if !isOk && len(dataList) > 0 {
		// 是否加入了公会
		guildMember, _ := new(dao.GuildDao).IsGuildMember(userId)
		if len(guildMember.GuildID) > 0 {
			isOk = true
		}
	}
	// 是否符合查询条件
	if isOk {
		roomId, _ := coreRedis.GetChatroomRedis().Get(context.Background(), common_userInfo.UserInWhichRoom(userId, clientType)).Result()
		if len(roomId) > 0 {
			roomInfo, _ := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
			if len(roomInfo.Id) > 0 {
				result = &user.RoomInfo{
					RoomId:   roomInfo.Id,
					RoomName: roomInfo.Name,
				}
				return result
			}
		}
	}
	// 无从业者资格或有资格未加入公会
	return nil
}

// 取消注销账号
func UnCancelAccount(userId string) {
	userDeleteApply := &dao.UserDeleteApplyDao{}
	applyModel, _ := userDeleteApply.GetUserDeleteApply(userId)
	if applyModel.Id > 0 {
		_ = userDeleteApply.UpdateUserDeleteApply(applyModel.Id, typedef_enum.UserDeleteStatusCancel)
		_ = new(dao.UserDao).UpdateById(&model.User{Id: userId, Status: typedef_enum.UserStatusNormal})
	}
}

// 获取用户在线得端
func GetUserLoginClientType(userId string) []string {
	res := []string{}
	for _, client := range typedef_enum.ClientTypeArray {
		key := common_userInfo.UserClientLoginStatus(userId, client)
		result, _ := coreRedis.GetImRedis().Exists(context.Background(), key).Result()
		if result == 1 {
			res = append(res, client)
		}
	}
	return res
}
