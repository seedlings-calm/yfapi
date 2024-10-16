package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

// UserDefaultAvatarDao
// @Description: 用户默认头像
type UserDefaultAvatarDao struct {
}

// GetUserDefaultAvatar
//
//	@Description: 根据性别获取用户默认头像
//	@receiver u
//	@param sex int -
//	@return result -
//	@return err -
func (u *UserDefaultAvatarDao) GetUserDefaultAvatar(sex int) (result model.UserDefaultAvatar, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserDefaultAvatar{}).Where("sex=? and status=1", sex).Scan(&result).Error
	return
}
