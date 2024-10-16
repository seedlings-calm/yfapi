package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef"

	"github.com/gin-gonic/gin"
)

type UserLoginRecordDao struct {
}

func (u *UserLoginRecordDao) LastOneByUserId(userId string) (res *model.UserLoginRecord) {
	coreDb.GetMasterDb().Model(&model.UserLoginRecord{}).Where("user_id = ?", userId).Order("create_time desc").First(&res)
	return
}

func (l *UserLoginRecordDao) LoginRecord(c *gin.Context, user *model.User, headData typedef.HeaderData) {
	//添加登录信息
	loginRecordModel := &model.UserLoginRecord{
		UserId:        user.Id,
		LoginPlatform: headData.Platform,
		LoginModel:    headData.Models,
		ClientVersion: headData.AppVersion,
		DeviceID:      headData.MachineCode,
		LoginIp:       headData.Ip,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	coreDb.GetMasterDb().Model(loginRecordModel).Create(loginRecordModel)
}

func (u *UserLoginRecordDao) FindByUserId(userId string) (res []*model.UserLoginRecord) {
	coreDb.GetMasterDb().Table("t_user_login_record tulr").Where("tulr.user_id = ?", userId).
		Where("tulr.create_time >= ?", time.Now().AddDate(0, -3, 0)).
		Where("tulr.create_time = (SELECT MAX(t2.create_time) FROM t_user_login_record AS t2 WHERE t2.device_id = tulr.device_id AND t2.user_id = ?)", userId).
		Order("tulr.create_time desc").
		Find(&res)
	return
}
