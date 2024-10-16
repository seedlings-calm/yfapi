package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/internal/model"
	"yfapi/util/easy"
)

type UserFollowDao struct {
}

// 查询是否已经关注过该用户
func (u *UserFollowDao) GetUserFollow(userid, followUserid string) (res *model.UserFollow, err error) {
	res = new(model.UserFollow)
	err = coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("user_id = ? AND focus_user_id = ?", userid, followUserid).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// GetUserFollowList 查询用户目标用户列表的关注信息list
func (u *UserFollowDao) GetUserFollowList(userId string, followUserIdList []string) (res []model.UserFollow, err error) {
	err = coreDb.GetSlaveDb().Model(&model.UserFollow{}).Where("user_id=? and focus_user_id in "+easy.GetInSql(followUserIdList), userId).Scan(&res).Error
	return
}

// GetUserFollowMap 查询用户目标用户列表的关注信息map
func (u *UserFollowDao) GetUserFollowMap(userId string, followUserIdList []string) (res map[string]model.UserFollow, err error) {
	res = make(map[string]model.UserFollow)
	followData, err := u.GetUserFollowList(userId, followUserIdList)
	if err != nil {
		return res, err
	}
	for _, info := range followData {
		res[info.FocusUserID] = info
	}
	return
}

// follow //修改互相关注状态
func (u *UserFollowDao) UpdateFollowUserStatus(id int, status bool) (err error) {
	err = coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("id = ?", id).Update("is_mutual_follow", status).Error
	return
}

// follow 获取关注列表
func (u *UserFollowDao) GetUserFollowingUserLists(userid string, page, size int) (res []*model.UserFollow, count int64, err error) {
	tx := coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("user_id = ?", userid).Count(&count)
	err = tx.Offset(page * size).Limit(size).Find(&res).Error
	if err != nil {
		return
	}
	return
}

// follow 获取粉丝列表
func (u *UserFollowDao) GetUserFollowersList(userid string, page, size int) (res []*model.UserFollow, count int64, err error) {
	tx := coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("focus_user_id = ?", userid).Count(&count)
	err = tx.Offset(page * size).Limit(size).Find(&res).Error
	if err != nil {
		return
	}
	return
}

// GetUserFollowerMap 根据传入用户列表获取用户粉丝map
func (u *UserFollowDao) GetUserFollowerMap(userId string, userIdList []string) (res map[string]struct{}) {
	res = make(map[string]struct{})
	var result []*model.UserFollow
	err := coreDb.GetSlaveDb().Model(model.UserFollow{}).Where("focus_user_id = ? and user_id in ?", userId, userIdList).Scan(&result).Error
	if err != nil {
		return
	}
	for _, info := range result {
		res[info.UserID] = struct{}{}
	}
	return
}

// follow 获取好友列表
func (u *UserFollowDao) GetUserFriendsList(userid string, page, size int) (res []*model.UserFollow, count int64, err error) {
	tx := coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("user_id = ? AND is_mutual_follow = ? ", userid, 1).Count(&count)
	err = tx.Offset(page * size).Limit(size).Find(&res).Error
	if err != nil {
		return
	}
	return
}

// 获取用户好友数量
func (u *UserFollowDao) GetUserFriendsNum(userId string) int64 {
	var num int64 = 0
	err := coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("user_id = ? and is_mutual_follow = 1", userId).Count(&num).Error
	if err != nil {
		coreLog.Error("GetUserFriendsNum err:%+v", err)
	}
	return num
}

// 获取关注数量
func (u *UserFollowDao) GetUserFollowedNum(userId string) int64 {
	var num int64 = 0
	err := coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("user_id = ?", userId).Count(&num).Error
	if err != nil {
		coreLog.Error("GetUserFollowedNum err:%+v", err)
	}
	return num
}

// 获取粉丝数量
func (u *UserFollowDao) GetUserFansNum(userId string) int64 {
	var num int64 = 0
	err := coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("focus_user_id = ?", userId).Count(&num).Error
	if err != nil {
		coreLog.Error("GetUserFansNum err:%+v", err)
	}
	return num
}

// IsUserFollowed
//
//	@Description: 查询用户是否关注目标用户
//	@receiver u
//	@param userId string -
//	@param targetUserId string -
//	@return bool -
func (u *UserFollowDao) IsUserFollowed(userId, targetUserId string) bool {
	var count int
	err := coreDb.GetSlaveDb().Model(model.UserFollow{}).Select("count(1)").Where("user_id=? and focus_user_id=?", userId, targetUserId).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// 获取粉丝ids
func (u *UserFollowDao) GetUserFollowersIds(userid string) (res []string) {
	coreDb.GetMasterDb().Model(&model.UserFollow{}).Where("focus_user_id = ?", userid).Pluck("user_id", &res)
	return
}
