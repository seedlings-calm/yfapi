package user

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	"yfapi/internal/logic"
	"yfapi/typedef/request"
	"yfapi/typedef/response"
	_ "yfapi/typedef/response/user"
)

// @Summary 我的足迹
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param  req	query	request.BasePageReq	true	"我的足迹列表参数"
// @Success 0 {object} response.BasePageRes{}
// @Success 0 {object} []user.UserVisitInfo{}
// @Router /v1/user/myVisit [get]
func GetUserVisitRecordList(c *gin.Context) {
	req := new(request.BasePageReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic.UserVisit).GetUserVisitRecordList(c, req))
}

// @Summary 访客记录
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Param  req	query	request.BasePageReq	true	"访客记录列表参数"
// @Success 0 {object} response.BasePageRes{}
// @Success 0 {object} []user.UserVisitInfo{}
// @Router /v1/user/visitMe [get]
func GetVisitUserRecordList(c *gin.Context) {
	req := new(request.BasePageReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic.UserVisit).GetVisitUserRecordList(c, req))
}

// @Summary 清除我的足迹
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Router /v1/user/visitClear [get]
func ClearUserVisitRecord(c *gin.Context) {
	new(logic.UserVisit).ClearUserVisitRecord(c)
	response.SuccessResponse(c, "")
}
