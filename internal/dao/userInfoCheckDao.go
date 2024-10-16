package dao

import (
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	error2 "yfapi/i18n/error"
	"yfapi/internal/model"
)

type UserInfoCheckDao struct {
}

// 查询是否有待审核内容  types 1头像审核 2语音审核
func (u *UserInfoCheckDao) FindOnePending(userId string, types int) bool {
	data := &model.UserInfoCheck{}
	err := coreDb.GetMasterDb().Model(&model.UserInfoCheck{}).Where("user_id = ? and type = ?", userId, types).
		Where("auto_status = ?", 0).First(data).Error
	if err == gorm.ErrRecordNotFound {
		return false
	}
	if err != nil {
		coreLog.Error("FindOnePending err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	return true
}

// 创建审核记录
func (u *UserInfoCheckDao) Create(userId string, types int, content string) int64 {
	ok := u.FindOnePending(userId, 1)
	if ok {
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFail,
			Msg:  nil,
		})
	}
	now := time.Now()
	data := &model.UserInfoCheck{
		UserID:     userId,
		Type:       types,
		Content:    content,
		CreateTime: &now,
	}
	err := coreDb.GetMasterDb().Model(&model.UserInfoCheck{}).Create(data).Error
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUnknown,
			Msg:  nil,
		})
	}
	return data.ID
}

// 审核通过
func (u *UserInfoCheckDao) PendingSuccess(id int64, result string) {
	data := &model.UserInfoCheck{}
	err := coreDb.GetMasterDb().Model(&model.UserInfoCheck{}).Where("id = ?", id).
		Where("auto_status = ?", 0).First(data).Error
	if err == gorm.ErrRecordNotFound {
		coreLog.Error("PendingSuccess ErrRecordNotFound id:%d", id)
	}
	upData := map[string]any{}
	if data.Type == 1 {
		upData["avatar"] = data.Content
	}
	if data.Type == 2 {
		upData["voice_url"] = data.Content
		upData["voice_length"] = data.VoiceLength
	}
	err = new(UserDao).UpdateUserFieldsByUserId(data.UserID, upData)
	if err != nil {
		coreLog.Error("PendingSuccess err:%+v", err)
		return
	}
	upData2 := map[string]any{
		"auto_status": 1,
		"update_time": time.Now(),
		"auto_result": result,
	}
	coreDb.GetMasterDb().Model(&model.UserInfoCheck{}).Where("id = ?", id).Updates(&upData2)
}

// 审核失败
func (u *UserInfoCheckDao) PendingFail(id int64, result string) {
	data := &model.UserInfoCheck{}
	err := coreDb.GetMasterDb().Model(&model.UserInfoCheck{}).Where("id = ?", id).
		Where("auto_status = ?", 0).First(data).Error
	if err == gorm.ErrRecordNotFound {
		coreLog.Error("PendingFail ErrRecordNotFound id:%d", id)
	}
	upData := map[string]any{
		"auto_status": 2,
		"update_time": time.Now(),
		"auto_result": result,
	}
	coreDb.GetMasterDb().Model(&model.UserInfoCheck{}).Where("id = ?", id).Updates(&upData)
}
