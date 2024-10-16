package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type RoomTypeDao struct {
}

// FindOne 条件查询
func (g *RoomTypeDao) FindOne(param *model.RoomType) (data *model.RoomType, err error) {
	data = new(model.RoomType)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// FindList 查询列表
func (g *RoomTypeDao) FindList(param *model.RoomType) (result []model.RoomType, err error) {
	err = coreDb.GetMasterDb().Where(param).Find(&result).Error
	return
}
