package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type UserFriendDao struct {
}

// follow 添加好友
func (u *UserFriendDao) FollowUser(data *model.UserFriends) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}
