package logic

import (
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/core/coreDb"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	typedef_enum "yfapi/typedef/enum"
	request_h5 "yfapi/typedef/request/h5"
)

type UserDeleteApply struct {
}

func (u *UserDeleteApply) DeleteUserApply(c *gin.Context, req *request_h5.UserDeleteApplyReq) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}

	//检测验证码
	sms := &Sms{
		Mobile:     userModel.Mobile,
		Code:       req.Captcha,
		RegionCode: userModel.RegionCode,
		Type:       typedef_enum.SmsCodeUserDeleteApply,
	}
	err = sms.CheckSms(c)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeCaptchaInvalid,
			Msg:  nil,
		})
	}

	// 生成申请记录
	param := &model.UserDeleteApply{
		UserId:     userId,
		Status:     typedef_enum.UserDeleteStatusApplying,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		ExpireTime: time.Now().Add(30 * 24 * time.Hour),
	}
	tx := coreDb.GetMasterDb().Begin()
	err = tx.Model(model.UserDeleteApply{}).Create(param).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 更改用户状态
	err = tx.Model(model.User{}).Where("id=?", userId).Updates(&model.User{Status: typedef_enum.UserStatusApplyInvalid}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	tx.Commit()
	return
}
