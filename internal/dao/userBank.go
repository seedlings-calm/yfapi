package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type UserBankDao struct {
}

// GetUserBankByUserIdAndBankNo
//
//	@Description: 查询用户银行卡信息ByUserIdAndBankNo
//	@receiver u
//	@param userId
//	@param bankNo
//	@return res
//	@return err
func (u *UserBankDao) GetUserBankByUserIdAndBankNo(userId string, bankNo string) (res model.UserBank, err error) {
	//err = coreDb.GetMasterDb().Model(model.UserBank{}).Where("user_id =?", userId).Where("bank_no=?", bankNo).First(&res).Error
	err = coreDb.GetMasterDb().Model(model.UserBank{}).Where("bank_no=?", bankNo).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}
	return
}

// GetUserBankByUserId
//
//	@Description: 查询用户银行卡信息列表
//	@receiver u
//	@param userId
//	@return res
//	@return err
func (u *UserBankDao) GetUserBankListByUserId(userId string) (res []*model.UserBank, err error) {
	err = coreDb.GetMasterDb().Model(model.UserBank{}).Where("user_id =?", userId).Order("create_time desc").Find(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}
	return
}

// BankAdd
//
//	@Description: 	添加银行卡
//	@receiver u
//	@param bank
//	@return err
func (u *UserBankDao) BankAdd(bank *model.UserBank) (err error) {
	err = coreDb.GetMasterDb().Model(model.UserBank{}).Create(bank).Error
	return
}

// 通过银行卡表id获取银行卡信息
func (u *UserBankDao) GetUserBankById(id int) (res model.UserBank, err error) {
	err = coreDb.GetMasterDb().Model(model.UserBank{}).Where("id =?", id).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}
	return
}
