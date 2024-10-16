package gift

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	"yfapi/internal/logic"
	request_gift "yfapi/typedef/request/gift"
	"yfapi/typedef/response"
	response_gift "yfapi/typedef/response/gift"
)

// @Summary 礼物列表
// @Description
// @Tags 礼物相关
// @Accept json
// @Produce json
// @Param  categoryType query int true "展示类目类型"
// @Param  liveType query int true "房间直播类型"
// @Param  roomType query int true "房间类型"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_gift.GiftListRes{}
// @Router /v1/giftList [get]
func GiftList(c *gin.Context) {
	req := new(request_gift.GiftListReq)
	handle.BindQuery(c, req)
	var res *response_gift.GiftListRes
	res = new(logic.Gift).GetRoomGiftList(c, req)
	response.SuccessResponse(c, res)
}

// @Summary 礼物资源列表
// @Description
// @Tags 礼物相关
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_gift.GiftSourceListRes{}
// @Router /v1/giftSourceList [get]
func GiftSourceList(c *gin.Context) {
	response.SuccessResponse(c, new(logic.Gift).GetGiftSourceList(c))
}

// @Summary 打赏礼物
// @Description
// @Tags 礼物相关
// @Accept json
// @Produce json
// @Param  data body  request_gift.SendGiftReq true "礼物列表参数"
// @Success 0 {object} response.Response{}
// @Router /v1/sendGift [post]
func SendGift(c *gin.Context) {
	req := new(request_gift.SendGiftReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(logic.Gift).SendGift(c, req))
}
