package h5

import (
	"yfapi/app/handle"
	logic_user "yfapi/internal/logic"
	request_h5 "yfapi/typedef/request/h5"
	"yfapi/typedef/response"
	response_h5 "yfapi/typedef/response/h5"

	"github.com/gin-gonic/gin"
)

// JoinGuild
//
//	@Description: 加入公会
func JoinGuild(c *gin.Context) {
	req := new(request_h5.JoinGuildReq)
	handle.BindBody(c, req)
	service := new(logic_user.Guild)
	res := service.JoinGuild(req, c)
	response.SuccessResponse(c, res)
}

// GuildInfo
// @Summary 获取公会信息
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  req body request_h5.GuildInfoReq   true "请求参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_h5.GuildInfoRes{}
// @Router /v1/h5/guild/info [post]
func GuildInfo(c *gin.Context) {
	req := new(request_h5.GuildInfoReq)
	handle.BindBody(c, req)
	service := new(logic_user.Guild)
	res := service.GetGuildInfo(req, c)
	response.SuccessResponse(c, res)

}

// GetGuildMemberList
// @Summary 获取公会成员列表
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  req body request_h5.GuildMemberListReq   true "请求参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_h5.GuildMemberListRes{}
// @Router /v1/h5/guild/memberList [post]
func GetGuildMemberList(c *gin.Context) {
	req := new(request_h5.GuildMemberListReq)
	handle.BindBody(c, req)
	service := new(logic_user.Guild)
	res := response_h5.GuildMemberListRes{}
	res = service.GetGuildMemberList(req, c)
	response.SuccessResponse(c, res)
}

// QuitGuildApply
// @Summary 退出公会申请
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  req body request_h5.QuitGuildApplyReq   true "请求参数"
// @Success 0 {object} response.Response{}
// @Router /v1/h5/guild/quitApply [post]
func QuitGuildApply(c *gin.Context) {
	req := new(request_h5.QuitGuildApplyReq)
	handle.BindBody(c, req)
	new(logic_user.Guild).QuitGuildApply(c, req)
	response.SuccessResponse(c, "")
}

// QuitGuildApplyCancel
// @Summary 取消退出公会申请
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  req body request_h5.GuildInfoReq   true "请求参数"
// @Success 0 {object} response.Response{}
// @Router /v1/h5/guild/quitApplyCancel [post]
func QuitGuildApplyCancel(c *gin.Context) {
	req := new(request_h5.GuildInfoReq)
	handle.BindBody(c, req)
	new(logic_user.Guild).QuitGuildApplyCancel(c, req)
	response.SuccessResponse(c, "")
}

// GetGuildPenaltyDetail
// @Summary 违约金详情
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Success 0 {object} response_h5.GuildPenaltyDetailRes{}
// @Router /v1/h5/guild/penaltyDetail [get]
func GetGuildPenaltyDetail(c *gin.Context) {
	response.SuccessResponse(c, new(logic_user.Guild).GetGuildPenaltyDetail(c))
}

// PayGuildPenalty
// @Summary 缴纳违约金
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  req body request_h5.GuildInfoReq   true "请求参数"
// @Success 0 {object} response.Response{}
// @Router /v1/h5/guild/payPenalty [post]
func PayGuildPenalty(c *gin.Context) {
	req := new(request_h5.GuildInfoReq)
	handle.BindBody(c, req)
	new(logic_user.Guild).PayGuildPenalty(c, req)
	response.SuccessResponse(c, "")
}
