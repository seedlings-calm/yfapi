package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
)

type UserTimelineFilterDao struct {
}

func (u *UserTimelineFilterDao) FindOne(params *model.UserTimelineFilter) *model.UserTimelineFilter {
	res := new(model.UserTimelineFilter)
	coreDb.GetMasterDb().Where(params).First(res)
	return res
}

// 删除
func (u *UserTimelineFilterDao) Del(data *model.UserTimelineFilter) (err error) {
	err = coreDb.GetMasterDb().Model(data).Delete(data).Error
	return
}

// 添加
func (u *UserTimelineFilterDao) Add(data *model.UserTimelineFilter) (err error) {
	err = coreDb.GetMasterDb().Create(data).Error
	return
}

func (u *UserTimelineFilterDao) GetSwitchType(userId, toId string, types int) bool {
	res := new(model.UserTimelineFilter)
	coreDb.GetMasterDb().Model(&model.UserTimelineFilter{}).Where("user_id = ? and to_id = ? and types = ?", userId, toId, types).First(res)
	if res.ID == 0 {
		return false
	}
	return true
}

func (u *UserTimelineFilterDao) GetFilterUserIds(userId string) (ids []string) {
	res1 := []string{}
	coreDb.GetMasterDb().Model(&model.UserTimelineFilter{}).Where("user_id = ? and types = ?", userId, enum.DontSeeHeMoments).Pluck("to_id", &res1)
	res2 := []string{}
	coreDb.GetMasterDb().Model(&model.UserTimelineFilter{}).Where("to_id = ? and types = ?", userId, enum.DontLetHeSeeMoments).Pluck("user_id", &res2)
	ids = append(ids, res1...)
	ids = append(ids, res2...)
	return
}

// 获取分页列表
func (u *UserTimelineFilterDao) GetList(userid string, page, size, types int) (res []*model.UserTimelineFilter, count int64, err error) {
	tx := coreDb.GetMasterDb().Model(&model.UserTimelineFilter{}).Where("user_id = ? AND types = ? ", userid, types).Count(&count)
	err = tx.Offset(page * size).Limit(size).Find(&res).Error
	if err != nil {
		return
	}
	return
}

// 获取数量
func (u *UserTimelineFilterDao) GetCount(userId string, types int) (count int64) {
	coreDb.GetMasterDb().Model(&model.UserTimelineFilter{}).Where("user_id = ? AND types = ? ", userId, types).Count(&count)
	return
}
