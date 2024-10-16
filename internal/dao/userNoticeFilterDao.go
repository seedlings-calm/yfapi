package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type UserNoticeFilterDao struct {
}

func (u *UserNoticeFilterDao) FindOne(params *model.UserNoticeFilter) *model.UserNoticeFilter {
	res := new(model.UserNoticeFilter)
	coreDb.GetMasterDb().Where(params).First(res)
	return res
}

// 删除
func (u *UserNoticeFilterDao) Del(data *model.UserNoticeFilter) (err error) {
	err = coreDb.GetMasterDb().Model(data).Delete(data).Error
	return
}

// 添加
func (u *UserNoticeFilterDao) Add(data *model.UserNoticeFilter) (err error) {
	err = coreDb.GetMasterDb().Create(data).Error
	return
}

func (u *UserNoticeFilterDao) GetSwitchType(userId, toId string, types int) bool {
	res := new(model.UserNoticeFilter)
	coreDb.GetMasterDb().Model(&model.UserNoticeFilter{}).Where("user_id = ? and to_id = ? and types = ?", userId, toId, types).First(res)
	if res.ID == 0 {
		return true
	}
	return false
}

func (u *UserNoticeFilterDao) GetIds(userId string, types int) []string {
	ids := []string{}
	coreDb.GetMasterDb().Model(&model.UserNoticeFilter{}).Where("to_id = ? and types = ?", userId, types).Pluck("user_id", &ids)
	return ids
}
