package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type UserAuthDao struct {
}

// TODO: 换成统一查询权限
func (u *UserAuthDao) IsRoles(userId string, roomId string, roleIds []string) (bool, error) {
	var res []model.AuthRoleAccess
	err := coreDb.GetMasterDb().Model(model.AuthRoleAccess{}).Where("user_id = ? and room_id = ? and role_id in (?)", userId, roomId, roleIds).Find(&res).Error
	if len(res) >= 1 {
		return true, nil
	}
	return false, err
}

// 添加权限
func (u *UserAuthDao) AddRule(data *model.AuthRule) error {
	err := coreDb.GetMasterDb().Create(data).Error
	return err
}

func (u *UserAuthDao) AddAuthRoleAccess(data *model.AuthRoleAccess) error {
	err := coreDb.GetMasterDb().Create(data).Error
	return err
}

func (u *UserAuthDao) DelAuthRoleAccess(data *model.AuthRoleAccess) error {
	err := coreDb.GetMasterDb().Where(data).Delete(&model.AuthRoleAccess{}).Error
	return err
}
