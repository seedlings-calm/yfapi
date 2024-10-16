package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type GuildDao struct {
}

// Create 添加
func (g *GuildDao) Create(data *model.Guild) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// FindOne 条件查询
func (g *GuildDao) FindOne(param *model.Guild) (data *model.Guild, err error) {
	data = new(model.Guild)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (g *GuildDao) Count(param *model.Guild) (count int64) {
	coreDb.GetMasterDb().Model(param).Where(param).Count(&count)
	return
}

// FindList 查询列表
func (g *GuildDao) FindList(param *model.Guild) (result []model.Guild, err error) {
	err = coreDb.GetMasterDb().Where(param).Find(&result).Error
	return
}

// FindByIds 根据ids查询结果
func (g *GuildDao) FindByIds(ids []string) (result []model.Guild) {
	coreDb.GetMasterDb().Find(&result, ids)
	return
}
