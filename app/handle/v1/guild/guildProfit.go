package guild

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	logic_guild "yfapi/internal/logic/guild"
	request_login "yfapi/typedef/request/guild"
	"yfapi/typedef/response"
	_ "yfapi/typedef/response/guild"
)

// GetGuildMemberProfitList
// @Summary 公会成员流水列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body	request_login.GetGuildMemberProfitReq	true	"参数"
// @Success 0 {object} []guild.MemberProfit{}
// @Router /v1/guild/memberProfit [post]
func GetGuildMemberProfitList(c *gin.Context) {
	req := new(request_login.GetGuildMemberProfitReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic_guild.GuildProfit).GetGuildMemberProfitList(c, req))
}

// GetGuildRoomProfitList
// @Summary 公会房间流水列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body	request_login.GetGuildRoomProfitReq	true	"参数"
// @Success 0 {object} []guild.RoomProfit{}
// @Router /v1/guild/roomProfit [post]
func GetGuildRoomProfitList(c *gin.Context) {
	req := new(request_login.GetGuildRoomProfitReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic_guild.GuildProfit).GetGuildRoomProfitList(c, req))
}

// GetGuildRewardList
// @Summary 公会礼物打赏详情列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body	request_login.GetGuildRewardListReq	true	"参数"
// @Success 0 {object} []guild.RewardDetail{}
// @Router /v1/guild/guildReward [post]
func GetGuildRewardList(c *gin.Context) {
	req := new(request_login.GetGuildRewardListReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic_guild.GuildProfit).GetGuildRewardList(c, req))
}

// GetAccountBillList
// @Summary 账户交易明细列表
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body	request_login.GetAccountBillReq	true	"参数"
// @Success 0 {object} []guild.AccountBill{}
// @Router /v1/guild/accountBill [post]
func GetAccountBillList(c *gin.Context) {
	req := new(request_login.GetAccountBillReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(logic_guild.GuildProfit).GetAccountBillList(c, req))
}
