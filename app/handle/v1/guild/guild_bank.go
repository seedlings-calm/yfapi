package guild

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/logic/guild"
	request_guild "yfapi/typedef/request/guild"
	request_users "yfapi/typedef/request/user"
	"yfapi/typedef/response"
)

// GuildBindBank
// @Summary 会长绑定银行卡
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  req body request_guild.BankBindReq  true "会长绑定银行卡参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/guildBindBank [post]
func GuildBindBank(c *gin.Context) {
	req := new(request_guild.BankBindReq)
	handle.BindBody(c, req)
	code := new(guild.GuildBank).GuildBindBank(c, req)
	if code != i18n_err.SuccessCode {
		panic(i18n_err.I18nError{
			Code: code,
			Msg:  nil,
		})
	}
	response.SuccessResponse(c, "")
}

// GuildWithdrawApply
// @Summary 会长提现申请
// @Description
// @Tags 公会后台
// @Accept json
// @Produce json
// @Param  req body request_users.UserWithdrawApplyReq  true "会长提现参数"
// @Success 0 {object} response.Response{}
// @Router /v1/guild/guildWithdrawApply [post]
func GuildWithdrawApply(c *gin.Context) {
	req := new(request_users.UserWithdrawApplyReq)
	handle.BindBody(c, req)
	new(guild.GuildBank).GuildWithdrawApply(c, req)
	response.SuccessResponse(c, "")
}
