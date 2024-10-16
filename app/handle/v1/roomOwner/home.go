package roomOwner

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	"yfapi/internal/logic/room"
	request_login "yfapi/typedef/request/guild"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
	_ "yfapi/typedef/response/guild"
	_ "yfapi/typedef/response/roomOwner"
)

// ExchangeDiamond
// @Summary 兑换钻石
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body  request_user.ExchangeDiamondReq true "礼物列表参数"
// @Success 0 {object} guild.AccountInfoRes{}
// @Router /v1/roomOwner/exchange [post]
func ExchangeDiamond(c *gin.Context) {
	req := new(request_user.ExchangeDiamondReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(room.RoomHome).ExchangeDiamond(c, req))
}

// GetAccountBillList
// @Summary 账户交易明细列表
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param  data body	request_login.GetAccountBillReq	true	"参数"
// @Success 0 {object} []guild.AccountBill{}
// @Router /v1/roomOwner/accountBill [post]
func GetAccountBillList(c *gin.Context) {
	req := new(request_login.GetAccountBillReq)
	handle.BindQuery(c, req)
	response.SuccessResponse(c, new(room.RoomHome).GetAccountBillList(c, req))
}

// RoomSearchUser
// @Summary 搜索用户信息
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param userNO query string  true "参数:用户昵称/ID"
// @Success 0 {object} roomOwner.SearchUser{}
// @Router /v1/roomOwner/user/search [get]
func RoomSearchUser(c *gin.Context) {
	response.SuccessResponse(c, new(room.RoomHome).RoomSearchUser(c))
}
