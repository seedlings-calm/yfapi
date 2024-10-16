package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	response_goods "yfapi/typedef/response/goods"
)

type GoodsUseDao struct {
}

func (GoodsUseDao) FindByGoodsType(typeId int) (res []*response_goods.GoodsListByTypesRes, err error) {

	db := coreDb.GetMasterDb().Table("t_goods_use tgu").
		Joins("left join t_goods ts on tgu.goods_id = ts.id").
		Where("tgu.status = ? and  tgu.is_del = ?", 1, 1)
	if typeId != 0 {
		db = db.Where("tgu.goods_type_id = ?", typeId)
	}
	err = db.Select("tgu.goods_id", "ts.name as goods_name", "tgu.goods_type_id", "ts.icon", "ts.animation_url", "ts.animation_json_url", "tgu.money", "tgu.moneys", "tgu.create_time").
		Order("create_time desc").Scan(&res).Error
	return
}

func (GoodsUseDao) FindGoodsTypeIds(isStatus bool) (res []int64, err error) {
	db := coreDb.GetMasterDb().Model(model.GoodsUse{})
	if isStatus {
		db = db.Where("status = 1")
	}
	err = db.Where(" is_del = 1").Group("goods_type_id").Pluck("goods_type_id", &res).
		Error
	return
}

func (GoodsUseDao) FirstByGoodsId(goodsId int) (res model.GoodsUse, err error) {
	err = coreDb.GetMasterDb().Model(model.GoodsUse{}).
		Where("goods_id = ?", goodsId).
		Where("status = 1 and is_del = 1").First(&res).Error
	return
}
