package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	response_goods "yfapi/typedef/response/goods"

	"gorm.io/gorm"
)

type UserGoodsDao struct {
}

// 判断时间是否过期,如果过期，返回true
func (UserGoodsDao) IsExpireTime(expireTime time.Time) bool {
	return expireTime.Before(time.Now())
}

func (UserGoodsDao) DelExpireTimeGoods(id uint64) error {
	err := coreDb.GetMasterDb().Model(model.UserGoods{}).Delete(model.UserGoods{Id: id}).Error
	return err
}

// 根据条件获取用户的装扮信息
func (UserGoodsDao) GetUserGoodsOne(userId, goodsId string) (res model.UserGoods, err error) {
	err = coreDb.GetMasterDb().Model(model.UserGoods{}).
		Where("user_id = ? and goods_id = ?", userId, goodsId).
		First(&res).
		Error
	return
}

// 根据条件获取用户的装扮信息-带有效期筛选
func (UserGoodsDao) GetUserGoodsOneByKey(userId, key string) (res model.UserGoods, err error) {
	err = coreDb.GetMasterDb().Model(model.UserGoods{}).
		Where("user_id = ? and goods_type_key = ? and expire_time >= ? and is_use = 2", userId, key, time.Now()).
		First(&res).
		Error
	return
}

// 根据商品类型key读取用户的商品
func (UserGoodsDao) GetUserGoodsByTypeKey(userId string, keys ...string) (res []*model.UserGoods, err error) {
	err = coreDb.GetMasterDb().Model(model.UserGoods{}).Where("user_id = ? and goods_type_key in ? and is_use = 2", userId, keys).Find(&res).Error
	return
}

// 获取用户持有的商品的分类ID
func (UserGoodsDao) GetGoodsTypeByUser(userId string) (res []int64, err error) {
	err = coreDb.GetMasterDb().Model(model.UserGoods{}).
		Where("user_id = ?", userId).Group("goods_type_id").Pluck("goods_type_id", &res).Error
	return
}

func (UserGoodsDao) FindListByUser(userId string) (res []response_goods.UserGoods, err error) {
	err = coreDb.GetMasterDb().Table("t_user_goods tug").
		Joins("left join t_goods tg on tg.id = tug.goods_id").
		Where("tug.user_id = ?", userId).
		Select("tug.goods_id", "tug.goods_type_id", "tug.expire_time", "tug.is_use", "tug.create_time", "tug.nums", "tg.name", "tg.icon", "tg.animation_url", "tg.animation_json_url").
		Order("tug.create_time desc").
		Scan(&res).Error
	return
}

// 清除装扮的红点, nums标记：触发清除增加100，其余次数为获取当前装扮的次数
func (UserGoodsDao) DelRedHot(userId string, goodsId string) (err error) {
	err = coreDb.GetMasterDb().Model(model.UserGoods{}).Where("user_id = ? and goods_id = ? and nums <= 100", userId, goodsId).Update("nums", gorm.Expr("nums+100")).Error
	return
}

// 获取用户是否有红点
func (UserGoodsDao) GetUserGoods(userId string) (res []model.UserGoods) {
	coreDb.GetMasterDb().Model(model.UserGoods{}).Where("user_id = ?", userId).Find(&res)
	return
}
