package guild

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	"yfapi/internal/logic/guild"
	request_login "yfapi/typedef/request/guild"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
	response_guild "yfapi/typedef/response/guild"
)

// AccountInfo
// @Summary 首页账户信息
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Success 0 {object} response_guild.AccountInfoRes{}
// @Router /v1/guild/accountInfo [get]
func AccountInfo(context *gin.Context) {
	service := new(guild.HomeData)
	response.SuccessResponse(context, service.AccountInfoLogic(context))
}

// ExchangeDiamond
// @Summary 兑换钻石
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  data body  request_user.ExchangeDiamondReq true "礼物列表参数"
// @Success 0 {object} response_guild.AccountInfoRes{}
// @Router /v1/guild/exchange [post]
func ExchangeDiamond(c *gin.Context) {
	req := new(request_user.ExchangeDiamondReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(guild.HomeData).ExchangeDiamond(c, req))
}

// GetGuildStatInfo
// @Summary 首页公会统计信息
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Success 0 {object} response_guild.StatGuildInfo{}
// @Router /v1/guild/statInfo [get]
func GetGuildStatInfo(c *gin.Context) {
	var resp response_guild.StatGuildInfo
	resp = new(guild.HomeData).GetGuildStatInfo(c)
	response.SuccessResponse(c, resp)
}

// GetGuildProfitInfo
// @Summary 首页公会流水信息
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Success 0 {object} response_guild.ProfitGuildInfo{}
// @Router /v1/guild/profitInfo [get]
func GetGuildProfitInfo(c *gin.Context) {
	var resp response_guild.ProfitGuildInfo
	resp = new(guild.HomeData).GetGuildProfitInfo(c)
	response.SuccessResponse(c, resp)
}

// GetRoomRankList
// @Summary 首页公会房间排行榜
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  req	query	request_login.GetRoomRankListReq	true	"参数"
// @Success 0 {object} []response_guild.RoomRank{}
// @Router /v1/guild/roomRank [get]
func GetRoomRankList(c *gin.Context) {
	req := new(request_login.GetRoomRankListReq)
	handle.BindQuery(c, req)
	var resp []response_guild.RoomRank
	resp = new(guild.HomeData).GetRoomRankList(c, req)
	response.SuccessResponse(c, resp)
}
