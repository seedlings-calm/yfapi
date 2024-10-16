package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type ConfigDiamondDao struct {
}

func (ConfigDiamondDao) Find(platform string) (res []*model.ConfigDiamond, err error) {
	err = coreDb.GetMasterDb().Model(model.ConfigDiamond{}).Where("platform = ? and status = 1", platform).Order("nums asc").Find(&res).Error
	return
}

// 根据keys查询一条商品
func (ConfigDiamondDao) FindOneByKeys(keys string) *model.ConfigDiamond {
	res := &model.ConfigDiamond{}
	coreDb.GetMasterDb().Model(&model.ConfigDiamond{}).Where("`keys` = ?", keys).First(res)
	return res
}
