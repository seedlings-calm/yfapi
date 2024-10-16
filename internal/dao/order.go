package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type OrderDao struct {
}

func (o *OrderDao) FindOne(params *model.Order) *model.Order {
	orderModel := &model.Order{}
	coreDb.GetMasterDb().Model(&model.Order{}).Where(params).First(orderModel)
	return orderModel
}
func (o *OrderDao) FindList(params *model.Order) (res []*model.Order) {
	// 执行查询，将结果存储在 res 变量中
	err := coreDb.GetMasterDb().Model(&model.Order{}).Where(params).Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	// 如果找到了记录，返回结果切片
	return res
}

// 检测今日是否提现
func (o *OrderDao) IsUserWithdraw(userId, todayStart, todayEnd string) (res model.Order, err error) {
	err = coreDb.GetMasterDb().Model(&model.Order{}).Where("user_id = ? And create_time between ? and ?",
		userId, todayStart, todayEnd).Order("create_time desc").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}
func (o *OrderDao) Save(params *model.Order) error {
	return coreDb.GetMasterDb().Save(params).Error
}

func (o *OrderDao) FindOneMap(params map[string]any) *model.Order {
	orderModel := &model.Order{}
	coreDb.GetMasterDb().Model(&model.Order{}).Where(params).First(orderModel)
	return orderModel
}

// 添加记录
func (o *OrderDao) Add(params *model.SubsidyActionRecord) error {
	return coreDb.GetMasterDb().Create(params).Error
}
