package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type GoodsDao struct {
}

// 根据id 读取物品信息,不涉及物品状态
func (GoodsDao) FirstByGoodsId(goodsId int) (res model.Goods, err error) {
	err = coreDb.GetMasterDb().Model(model.Goods{}).
		Where("id = ?", goodsId).
		First(&res).Error
	return
}
