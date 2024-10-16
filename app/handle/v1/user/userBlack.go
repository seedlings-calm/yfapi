package user

import (
	error2 "yfapi/i18n/error"
	"yfapi/internal/logic"
	"yfapi/typedef/response"
	response_room "yfapi/typedef/response/room"

	"github.com/gin-gonic/gin"
)

// @Summary 执行拉黑
// @Description
// @Tags 拉黑管理
// @Accept json
// @Produce json
// @Param  toId query int  true "被拉黑用户ID"
// @Success 0 {object} response.Response{}
// @Router /v1/blacklist/add [get]
func AddBlacklist(c *gin.Context) {
	toId := c.Query("toId")
	if toId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	bLogic := new(logic.BlackListLogic)
	bLogic.AddBlacklist(c, toId)
	response.SuccessResponse(c, "")
}

// @Summary 取消拉黑
// @Description
// @Tags 拉黑管理
// @Accept json
// @Produce json
// @Param toId query int  true "被取消拉黑的ID"
// @Success 0 {object} response.Response{}
// @Router /v1/blacklist/del [get]
func DelBlacklist(c *gin.Context) {
	toId := c.Query("toId")
	if toId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	bLogic := new(logic.BlackListLogic)
	bLogic.DelBlacklist(c, toId)
	response.SuccessResponse(c, "")
	// tx := coreDb.GetMasterDb().Begin()
	// for i := 0; i < 10; i++ {

	// 	orderInfo := &model.Order{
	// 		OrderId:         new(accountBook.Order).OrderNum(accountBook.ORDER_SC),
	// 		UserId:          "1",
	// 		RoomId:          "0",
	// 		GuildId:         "0",
	// 		TotalAmount:     strconv.Itoa((i + 1) * 100),
	// 		PayAmount:       "100",
	// 		DiscountsAmount: "0",
	// 		Num:             i + 1,
	// 		Currency:        accountBook.CURRENCY_CNY,
	// 		OrderType:       accountBook.ChangeDiamondMallConsumption,
	// 		OrderStatus:     1,
	// 		PayType:         1,
	// 		PayStatus:       1,
	// 		Note:            "测试商城订单",
	// 		CreateTime:      time.Now(),
	// 		UpdateTime:      time.Now(),
	// 	}
	// 	tx.Create(orderInfo)
	// }

	// tx.Commit()
}

// @Summary 用户黑名单列表
// @Description
// @Tags 用户相关
// @Accept json
// @Produce json
// @Success 0 {object} []response_room.BlackListAndUserInfo{}
// @Router /v1/blacklist/userList [get]
func GetUserBlacklist(c *gin.Context) {
	var res []*response_room.BlackListAndUserInfo
	res = new(logic.BlackListLogic).GetUserBlacklist(c)
	response.SuccessResponse(c, res)
}
