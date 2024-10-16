package bank

import (
	"time"
	"yfapi/core/coreDb"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/util/bankCard"
)

type Bank struct {
}

func (b Bank) BankBind(userId, bankNo, bankName, bankHolder, bankBranch string) (code error2.ErrCode) {
	//查询有无添加过此银行卡
	bankDao := dao.UserBankDao{}
	bankInfo, err := bankDao.GetUserBankByUserIdAndBankNo(userId, bankNo)
	if err != nil {
		return error2.ErrorCodeReadDB
	}
	if bankInfo.Id > 0 {
		return error2.ErrorCodeBankCardExist
	}
	//添加银行卡信息
	info, err := bankCard.GetBankInfo(bankNo)
	if err != nil {
		return error2.ErrorCodeBankCardInfoErr
	}
	bank := &model.UserBank{
		UserId:     userId,
		BankName:   bankName,
		BankNo:     bankNo,
		BankHolder: bankHolder,
		BankBranch: bankBranch,
		IsDefault:  1,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		BankCode:   info.Bank,
	}
	tx := coreDb.GetMasterDb().Begin()
	_ = tx.Model(&model.UserBank{}).Where("user_id=?", userId).Update("is_default", 0).Error
	err = tx.Create(&bank).Error
	if err != nil {
		tx.Rollback()
		return error2.ErrorCodeReadDB
	}
	tx.Commit()
	return
}
