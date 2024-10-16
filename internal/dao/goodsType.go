package dao

import (
	"errors"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	response_goods "yfapi/typedef/response/goods"

	"gorm.io/gorm"
)

type GoodsType struct {
}

// isStatus 判断是否加上分类有效的判断
func (GoodsType) FindsByIds(ids []int64, isStatus bool) (res []*response_goods.GoodsTypesListRes, err error) {

	db := coreDb.GetMasterDb().Model(model.GoodsType{})
	if len(ids) > 0 {
		db = db.Where("id in (?)", ids)
	}
	if isStatus {
		db = db.Where("status = 1")
	}
	err = db.Where("is_del = 1").Select("id", "icon", "keys", "name", "status", "sort", "create_time").
		Order("sort desc").Order("create_time desc").
		Scan(&res).Error
	return
}

// FindByTypeId
//
//	@Description: 根据id获取分类信息
//	@receiver GoodsType
//	@param id
//	@return res
//	@return err
func (GoodsType) FindByTypeId(id int) (res *model.GoodsType, err error) {
	err = coreDb.GetMasterDb().Model(model.GoodsType{}).Where("id = ?", id).Scan(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}
