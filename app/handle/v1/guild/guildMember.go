package guild

import (
	"yfapi/app/handle"
	i18n_err "yfapi/i18n/error"
	logic_guild "yfapi/internal/logic/guild"
	request_room "yfapi/typedef/request/guild"
	"yfapi/typedef/response"
	response_guild "yfapi/typedef/response/guild"

	"github.com/gin-gonic/gin"
)

// MemberGroup
// @Summary 公会分组列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param	req	body	request_room.GuildGroupListreq	true	"参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response.AdminPageRes{}
// @Success 0 {object} []request_room.MemberGroupUpdateRes{}
// @Router /v1/guild/memberGroup [post]
func MemberGroup(context *gin.Context) {
	req := new(request_room.GuildGroupListreq)
	handle.BindBody(context, req)
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(context, service.MemberGroup(context, req))
}

// SaveMemberGroup
// @Summary 公会分组创建
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param	req	body	request_room.AddGuildGroupReq	true	"参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/savememberGroup [post]
func SaveMemberGroup(context *gin.Context) {
	req := new(request_room.AddGuildGroupReq)
	handle.BindBody(context, req)
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(context, service.SaveMemberGroup(context, req))
}

// MemberGroupUpdate
// @Summary 公会分组修改
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  req body request_room.MemberGroupUpdateReq  true "参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/memberGroupUpdate [post]
func MemberGroupUpdate(c *gin.Context) {
	req := new(request_room.MemberGroupUpdateReq)
	handle.BindBody(c, req)
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(c, service.MemberGroupUpdate(c, req))
}

// MemberGroupDelete
// @Summary 公会分组删除
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  id  query  int  true "分组ID"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/groupDelete [get]
func MemberGroupDelete(c *gin.Context) {
	groupId := c.Query("id")
	if groupId == "" {
		panic(i18n_err.ErrorCodeParam)
	}
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(c, service.MemberGroupDelete(c, groupId))

}

// SetGroupByMembers
// @Summary 设置分组
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  req body request_room.SetGroupByMembersReq{}   true "请求参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/setGroup [post]
func SetGroupByMembers(c *gin.Context) {
	req := new(request_room.SetGroupByMembersReq)
	handle.BindBody(c, req)
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(c, service.SetGroupByMembers(c, req))
}

// MemberList
// @Summary 公会成员列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param req body request_room.GuildMemberListReq false "参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object}  []response_guild.GuildMemberListRes{}
// @Router /v1/guild/memberList [post]
func MemberList(context *gin.Context) {
	req := new(request_room.GuildMemberListReq)
	handle.BindBody(context, req)
	service := new(logic_guild.GuildMember)
	_ = make([]response_guild.GuildMemberListRes, 0)
	response.SuccessResponse(context, service.MemberList(context, req))
}

// MemberIdcards
// @Summary 从业者身份列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param userId query string  true "用户ID"
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_guild.MemberIdcardsInfoRes{}
// @Router /v1/guild/memberIdcards [get]
func MemberIdcards(c *gin.Context) {
	service := new(logic_guild.GuildMember)
	_ = make([]response_guild.MemberIdcardsInfoRes, 0)
	response.SuccessResponse(c, service.MemberIdcards(c, c.Query("userId")))
}

// MembershipList
// @Summary 入会申请列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param req body request_room.MemberShipListreq false "参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object}  []response_guild.MemberJoinApplyInfo{}
// @Router /v1/guild/membershipList [post]
func MembershipList(context *gin.Context) {
	req := new(request_room.MemberShipListreq)
	handle.BindBody(context, req)
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(context, service.MemberShipList(context, req))
}

// WithdrawMembershipList
// @Summary 退会申请列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param req body request_room.MemberShipListreq false "参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object}  []response_guild.LeaveMemberShipListRsp{}
// @Router /v1/guild/withdrawList [post]
func WithdrawMembershipList(context *gin.Context) {
	req := new(request_room.MemberShipListreq)
	handle.BindBody(context, req)
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(context, service.WithdrawMemberShipList(context, req))
}

// GetPractitionerActionRecord
// @Summary 从业者行为记录列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param userId query string true "用户长ID"
// @Success 0 {object} response.Response{}
// @Success 0 {object}  []response_guild.UserPractitionerAction{}
// @Router /v1/guild/practitionerAction [get]
func GetPractitionerActionRecord(c *gin.Context) {
	response.SuccessResponse(c, new(logic_guild.GuildMember).GetPractitionerActionRecord(c))
}

// MemberApplyReview
// @Summary 入会申请审核
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  req body request_room.GuildMemberApplyReviewReq  true "入会申请审核参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/guildMemberApplyReview [post]
func MemberApplyReview(context *gin.Context) {
	req := new(request_room.GuildMemberApplyReviewReq)
	handle.BindBody(context, req)
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(context, service.MemberApplyReview(context, req))
}

// MemberWithdrawApplyReview
// @Summary 退会申请审核
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  req body request_room.GuildMemberWithdrawReviewReq  true "退会申请审核参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/guildMemberWithdrawApplyReview [post]
func MemberWithdrawApplyReview(context *gin.Context) {
	req := new(request_room.GuildMemberWithdrawReviewReq)
	handle.BindBody(context, req)
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(context, service.MemberWithdrawApplyReview(context, req))
}

// GuildKickoutMember
// @Summary 踢出公会成员
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param userId query string  true "被踢出用户ID"
// @Param reason query string  true "踢出原因"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/kickOutMember [get]
func GuildKickoutMember(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		panic(i18n_err.ErrorCodeParam)
	}
	service := new(logic_guild.GuildMember)
	response.SuccessResponse(c, service.GuildKickoutMember(c, userId, c.Query("reason")))
}
