package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type RoomTemplateDao struct {
}

// FindOne 条件查询
func (g *RoomTemplateDao) FindOne(param *model.RoomTemplate) (data *model.RoomTemplate, err error) {
	data = new(model.RoomTemplate)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	return
}
