package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	response_gift "yfapi/typedef/response/gift"
)

type GiftSendCountDao struct {
}

func (g *GiftSendCountDao) GetListByIdList(idList []string) (result []response_gift.GiftSendCount, err error) {
	err = coreDb.GetSlaveDb().Model(model.GiftSendCount{}).Where("id in ?", idList).Scan(&result).Error
	return
}
