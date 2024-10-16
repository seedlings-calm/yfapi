package logic

import (
	"github.com/gin-gonic/gin"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/typedef/response/user"
)

type UserPrivacy struct {
}

// 获取隐私设置
func (u *UserPrivacy) GetPrivacySetting(c *gin.Context) (res user.PrivacySettingResp) {
	userId := helper.GetUserId(c)
	res.BlacklistNum = new(dao.UserBlackListDao).GetCount(userId, 2)
	res.DontSeeHeMomentsNum = new(dao.UserTimelineFilterDao).GetCount(userId, 2)
	res.DontLetHeSeeMeNum = new(dao.UserTimelineFilterDao).GetCount(userId, 1)
	return
}
