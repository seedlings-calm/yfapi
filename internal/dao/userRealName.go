package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type UserRealNameDao struct {
}

func (u *UserRealNameDao) FindOne(params *model.UserRealName) model.UserRealName {
	res := model.UserRealName{}
	coreDb.GetMasterDb().Model(params).Where(params).Order("create_time desc").Limit(1).First(&res)
	return res
}
