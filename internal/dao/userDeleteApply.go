package dao

import (
	"database/sql"
	"errors"
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	typedef_enum "yfapi/typedef/enum"
)

type UserDeleteApplyDao struct {
}

func (u *UserDeleteApplyDao) GetUserDeleteApply(userId string) (result model.UserDeleteApply, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserDeleteApply{}).Where("user_id=? and status=?", userId, typedef_enum.UserDeleteStatusApplying).Scan(&result).Error
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}

func (u *UserDeleteApplyDao) UpdateUserDeleteApply(id int64, status int) (err error) {
	return coreDb.GetMasterDb().Model(model.UserDeleteApply{}).Where("id=?", id).Updates(&model.UserDeleteApply{
		Status:     status,
		UpdateTime: time.Now(),
	}).Error
}

func (u *UserDeleteApplyDao) GetCanDeleteUserList() (result []model.UserDeleteApply, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserDeleteApply{}).Where("status=? and expire_time<=?", typedef_enum.UserDeleteStatusApplying, time.Now().Format(time.DateTime)).Scan(&result).Error
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}
