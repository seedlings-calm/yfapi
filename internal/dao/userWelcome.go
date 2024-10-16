package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type AutoWelcome struct{}

func (a *AutoWelcome) Save(param *model.UserAutoWelcome) error {
	err := coreDb.GetMasterDb().Model(model.UserAutoWelcome{}).
		Where(&model.UserAutoWelcome{
			UserID: param.UserID,
			State:  1,
		}).Assign(&model.UserAutoWelcome{
		WelcomeContent: param.WelcomeContent,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}).FirstOrCreate(&model.UserAutoWelcome{}).Error
	return err
}

func (a *AutoWelcome) Del(userId string) error {
	return coreDb.GetMasterDb().Model(model.UserAutoWelcome{}).Where("user_id = ? and state = 1", userId).Updates(map[string]interface{}{
		"state":       2,
		"staff_name":  userId,
		"update_time": time.Now(),
	}).Error
}

func (a *AutoWelcome) FindToUserIds(ids []string) (res []string) {
	coreDb.GetMasterDb().Model(model.UserAutoWelcome{}).Where("user_id in (?) and state = 1", ids).Select("user_id").Scan(&res)
	return
}

func (a *AutoWelcome) FirstContent(userId string) (res string) {
	coreDb.GetMasterDb().Model(model.UserAutoWelcome{}).Where("user_id = ? and state = 1", userId).Select("welcome_content").Scan(&res)
	return
}
