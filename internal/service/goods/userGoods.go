package service_goods

import (
	"errors"
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	response_goods "yfapi/typedef/response/goods"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type UserGoods struct {
}

// 获取用户的单个装扮信息
func (ug UserGoods) GetGoodsByKey(userId string, isImg bool, key string) (res *response_goods.SpecialEffects) {
	res = new(response_goods.SpecialEffects)
	userGoodsDao := dao.UserGoodsDao{}
	daoRes, err := userGoodsDao.GetUserGoodsOneByKey(userId, key)
	if err != nil {
		return
	}
	var goodsRes model.Goods
	goodsDao := dao.GoodsDao{}
	goodsRes, err = goodsDao.FirstByGoodsId(cast.ToInt(res.GoodsId))
	if err != nil {
		return
	}
	res.GoodsId = daoRes.GoodsId
	res.GoodsTypeKey = daoRes.GoodsTypeKey
	res.GoodsName = goodsRes.Name
	if isImg {
		res.GoodsIcon = helper.FormatImgUrl(goodsRes.Icon)
		res.AnimationUrl = helper.FormatImgUrl(goodsRes.AnimationUrl)
		res.AnimationJsonUrl = helper.FormatImgUrl(goodsRes.AnimationJsonUrl)
	} else {
		res.GoodsIcon = goodsRes.Icon
		res.AnimationUrl = goodsRes.AnimationUrl
		res.AnimationJsonUrl = goodsRes.AnimationJsonUrl
	}
	return
}

// 根据商品key ，查询用户的使用商品,isImg 是否需要图片全路径，需要true
func (ug UserGoods) GetGoodsByKeys(userId string, isImg bool, keys ...string) (res []*response_goods.SpecialEffects, err error) {
	userGoodsDao := dao.UserGoodsDao{}
	daoRes, err := userGoodsDao.GetUserGoodsByTypeKey(userId, keys...)
	if err != nil {
		return
	}
	if len(daoRes) > 0 {
		goodsDao := dao.GoodsDao{}
		for _, v := range daoRes {
			item := new(response_goods.SpecialEffects)
			item.GoodsId = v.GoodsId
			item.GoodsTypeKey = v.GoodsTypeKey
			var goodsRes model.Goods
			goodsRes, err = goodsDao.FirstByGoodsId(cast.ToInt(item.GoodsId))
			if err != nil {
				continue
			}
			item.GoodsName = goodsRes.Name
			if isImg {
				item.GoodsIcon = helper.FormatImgUrl(goodsRes.Icon)
				item.AnimationUrl = helper.FormatImgUrl(goodsRes.AnimationUrl)
				item.AnimationJsonUrl = helper.FormatImgUrl(goodsRes.AnimationJsonUrl)
			} else {
				item.GoodsIcon = goodsRes.Icon
				item.AnimationUrl = goodsRes.AnimationUrl
				item.AnimationJsonUrl = goodsRes.AnimationJsonUrl

			}
			res = append(res, item)
		}

	}

	return
}

// 物品发放给用户 统一接口
// sendSource 发放来源，枚举：GoodsGrantSourceAdmin
// days 发放天数
func (ug UserGoods) SendGoodsToUser(c *gin.Context, userId string, goodsId, sendSource, days int, callback func(c *gin.Context, userId, goodsId string, isUse int)) (err error) {
	if userId == "" || goodsId == 0 || days == 0 {
		return errors.New("error params")
	}
	goodsInfo, err := new(dao.GoodsDao).FirstByGoodsId(goodsId)
	if err != nil {
		return
	}
	var nowTime = time.Now()
	tx := coreDb.GetMasterDb().Begin()
	var grantM model.GoodsGrant
	grantM.GoodsId = int64(goodsInfo.Id)
	grantM.Source = sendSource
	grantM.Day = days
	grantM.UserId = userId
	grantM.CreateTime = nowTime
	err = tx.Model(model.GoodsGrant{}).Create(&grantM).Error
	if err != nil {
		tx.Rollback()
		return
	}
	var isUse int
	isUse, err = ug.AddGoods(tx, userId, cast.ToString(goodsInfo.Id), cast.ToString(goodsInfo.TypeId), goodsInfo.TypeKey, days, "")
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	callback(c, userId, cast.ToString(goodsId), isUse)
	return
}

// 购买的物品存储到装扮背包里，涉及逻辑使用装扮
func (ug UserGoods) AddGoods(tx *gorm.DB, userId, goodsId string, goodsTypeId, goodsTypeKey string, days int, orderId string) (isUse int, err error) {
	var res []model.UserGoods
	tx.Model(model.UserGoods{}).Where("user_id = ?", userId).Where("goods_type_id = ? and is_use = 2", goodsTypeId).Find(&res)
	isUse = 2
	if len(res) > 0 { //如果当前商品的类型已经存在使用的装扮，这个不使用
		isUse = 1
	}
	where := &model.UserGoods{
		UserId:  userId,
		GoodsId: goodsId,
	}
	var goods model.UserGoods
	err = tx.Model(model.UserGoods{}).
		Where(where).First(&goods).Error
	if err != nil {

		if err != gorm.ErrRecordNotFound {
			return
		}
		do := &model.UserGoods{
			UserId:       userId,
			GoodsId:      goodsId,
			GoodsTypeId:  goodsTypeId,
			GoodsTypeKey: goodsTypeKey,
			IsUse:        isUse,
			ExpireTime:   time.Now().AddDate(0, 0, days),
			Nums:         1,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		if orderId != "" {
			do.OrderId = orderId
		}
		err = tx.Model(model.UserGoods{}).Create(do).Error
	} else {
		if goods.Id > 0 {

			goods.ExpireTime = goods.ExpireTime.AddDate(0, 0, days)
			goods.Nums = goods.Nums + 1
			if orderId != "" {
				goods.OrderId = goods.OrderId + "," + orderId
			}
			err = tx.Save(&goods).Error
		}
	}

	return
}

// 此方法操作，执行使用装扮时，下掉之前使用装扮，如果是下掉使用装扮，不额外处理
func (ug UserGoods) UpdateUserGoods(c *gin.Context, goodsInfo model.UserGoods, clientType string, callBack func(c *gin.Context, userId, goodsId string, isUse int)) (err error) {
	tx := coreDb.GetMasterDb().Begin()
	var isUse int
	//根据当前装扮的状态执行反状态
	if goodsInfo.IsUse == 1 {
		isUse = 2
	} else {
		isUse = 1
	}

	if isUse == 2 { //如果是使用装扮，同类型的需要取消装扮
		err = tx.Model(model.UserGoods{}).Where("user_id = ? and goods_type_id = ? and expire_time >= ?", goodsInfo.UserId, goodsInfo.GoodsTypeId, time.Now()).Updates(map[string]interface{}{
			"is_use":      1,
			"update_time": time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Model(model.UserGoods{}).Where("user_id = ? and goods_id = ?", goodsInfo.UserId, goodsInfo.GoodsId).Update("is_use", isUse).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	callBack(c, goodsInfo.UserId, goodsInfo.GoodsId, isUse)
	return
}
