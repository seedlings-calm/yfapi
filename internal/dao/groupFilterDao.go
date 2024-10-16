package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
)

type GroupFilterDao struct {
}

// 查询一条记录
func (u *GroupFilterDao) FindOne(params *model.GroupFilter) *model.GroupFilter {
	res := new(model.GroupFilter)
	coreDb.GetMasterDb().Where(params).First(res)
	return res
}

// 删除规则
func (u *GroupFilterDao) Del(data *model.GroupFilter) (err error) {
	err = coreDb.GetMasterDb().Model(data).Delete(data).Error
	return
}

// 添加
func (u *GroupFilterDao) Add(data *model.GroupFilter) (err error) {
	err = coreDb.GetMasterDb().Create(data).Error
	return
}

func (u *GroupFilterDao) Save(data *model.GroupFilter) {
	coreDb.GetMasterDb().Save(data)
	return
}

// 群禁言状态
func (u *GroupFilterDao) MuteStatus(groupId string) bool {
	res := new(model.GroupFilter)
	coreDb.GetMasterDb().Model(&model.GroupFilter{}).Where("group_id = ? and types = ?", groupId, enum.GroupMuteSwitch).First(res)
	if res.ID > 0 {
		return true
	}
	return false
}
