package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type DaoUserPractitionerCerd struct {
	UserId string
}

func (u *DaoUserPractitionerCerd) Find() (result []model.UserPractitionerCred, err error) {
	err = coreDb.GetMasterDb().Where("user_id = ? and status  = 1", u.UserId).Find(&result).Error
	return
}

func (u *DaoUserPractitionerCerd) First(cerd int) (res model.UserPractitionerCred, err error) {
	err = coreDb.GetMasterDb().Model(model.UserPractitionerCred{}).Where("user_id = ? and practitioner_type = ? and status = 1", u.UserId, cerd).First(&res).Error
	return
}
