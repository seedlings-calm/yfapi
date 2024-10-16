package logic

import (
	"github.com/gin-gonic/gin"
	"yfapi/core/coreDb"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_bank "yfapi/internal/service/bank"
	service_user "yfapi/internal/service/user"
	typedef_enum "yfapi/typedef/enum"
	request_user "yfapi/typedef/request/user"
	response_user "yfapi/typedef/response/user"
)

type UserBank struct{}

// 添加银行卡
func (u *UserBank) Add(c *gin.Context, req *request_user.UserBankAddReq) (code error2.ErrCode) {
	userId := helper.GetUserId(c)
	userInfo := service_user.GetUserBaseInfo(userId)
	if req.BankHolder != userInfo.TrueName || userInfo.RealNameStatus != typedef_enum.UserRealNameAuthenticated {
		return error2.ErrorCodeIDCardAuth
	}
	if req.Mobile != userInfo.Mobile {
		return error2.ErrCodeMobileOrRegionCode
	}
	//	校验验证码
	sms := &Sms{
		Mobile:     req.Mobile,
		Code:       req.Code,
		RegionCode: req.RegionCode,
		Type:       typedef_enum.SmsCodeBindBankCard,
	}
	err := sms.CheckSms(c)
	if err != nil {
		return error2.ErrorCodeCaptchaInvalid
	}
	code = service_bank.Bank{}.BankBind(userId, req.BankNo, req.BankName, req.BankHolder, req.BankBranch)
	return code
}

// 解绑银行卡
func (u *UserBank) UnBind(c *gin.Context, req *request_user.UserBankUnBindReq) {
	userId := helper.GetUserId(c)
	//检测身份证号
	userDao := dao.UserDao{}
	IdNoInfo, err := userDao.CheckUserIdNo(req.IdNo, userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if IdNoInfo.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeIdCardNotExist,
			Msg:  nil,
		})
	}
	//查询有无添加过此银行卡
	bankDao := dao.UserBankDao{}
	bankInfo, err := bankDao.GetUserBankById(req.Id)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if bankInfo.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserBankNotBind,
			Msg:  nil,
		})
	}
	tx := coreDb.GetMasterDb().Begin()
	err = tx.Delete(&model.UserBank{}, req.Id).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 如果该银行卡是默认银行卡，将最近添加的一张银行卡默认状态改为1
	if bankInfo.IsDefault == 1 {
		bankInfos, err := bankDao.GetUserBankListByUserId(userId)
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		// 找到最近添加的银行卡，且不是要删除的银行卡
		var newDefaultBankInfo *model.UserBank
		for _, Info := range bankInfos {
			if Info.Id != bankInfo.Id {
				newDefaultBankInfo = Info
				break
			}
		}
		// 如果找到了新的默认银行卡，更新其IsDefault状态为1
		if newDefaultBankInfo != nil {
			e := tx.Model(&model.UserBank{}).Where("id =?", newDefaultBankInfo.Id).Update("is_default", 1).Error
			if e != nil {
				tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
		}
	}
	tx.Commit()
}

// GetBankList
//
//	@Description: 获取用户银行卡列表
//	@receiver u
//	@param c
//	@return res
func (u *UserBank) GetBankList(c *gin.Context) (res []*response_user.UserBankInfo) {
	userId := helper.GetUserId(c)
	//查询有无添加过此银行卡
	bankDao := dao.UserBankDao{}
	bankInfos, err := bankDao.GetUserBankListByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	res = make([]*response_user.UserBankInfo, 0)
	for _, Info := range bankInfos {
		res = append(res, &response_user.UserBankInfo{
			Id:         Info.Id,
			UserId:     Info.UserId,
			BankName:   Info.BankName,
			BankNo:     Info.BankNo,
			BankHolder: Info.BankHolder,
			BankBranch: Info.BankBranch,
			BankCode:   Info.BankCode,
			IsDefault:  Info.IsDefault,
			CreateTime: Info.CreateTime,
			UpdateTime: Info.UpdateTime,
		})
	}
	return
}
