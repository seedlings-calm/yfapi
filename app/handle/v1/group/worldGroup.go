package group

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	"yfapi/internal/logic"
	"yfapi/typedef/request/group"
	"yfapi/typedef/response"
)

// 世界频道设置
func WorldGroupSetting(c *gin.Context) {
	logicSer := new(logic.Group)
	response.SuccessResponse(c, logicSer.GetWorldGroupSetting(c))
}

// 世界频道禁言
func WorldGroupMute(c *gin.Context) {
	logicSer := new(logic.Group)
	logicSer.WorldGroupMute(c)
	response.SuccessResponse(c, true)
}

// 消息通知设置
func WorldGroupNoticeSetting(c *gin.Context) {
	logicSer := new(logic.Group)
	req := new(group.WorldGroupNoticeSettingReq)
	handle.BindBody(c, req)
	logicSer.WorldGroupNoticeSetting(c, req)
	response.SuccessResponse(c, true)
}
