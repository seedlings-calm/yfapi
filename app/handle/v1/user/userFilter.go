package user

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	logic_user "yfapi/internal/logic"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
)

// 动态规则设置
func TimelineFilter(c *gin.Context) {
	req := new(request_user.TimelineFilterReq)
	handle.BindBody(c, req)
	service := new(logic_user.UserFilter)
	service.TimelineFilter(c, req)
	response.SuccessResponse(c, true)
}

// 动态规则列表
func TimelineFilterList(c *gin.Context) {
	req := new(request_user.GetTimelineFilterListReq)
	handle.BindQuery(c, req)
	service := new(logic_user.UserFilter)
	res := service.GetTimelineFilterList(c, req)
	response.SuccessResponse(c, res)
}

// 通知规则开关
func NoticeFilter(c *gin.Context) {
	req := new(request_user.NoticeFilterReq)
	handle.BindBody(c, req)
	service := new(logic_user.UserFilter)
	service.NoticeFilter(c, req)
	response.SuccessResponse(c, true)
}
