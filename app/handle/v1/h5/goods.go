package h5

import (
	"strconv"
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/logic"
	request_goods "yfapi/typedef/request/goods"
	"yfapi/typedef/response"
	response_goods "yfapi/typedef/response/goods"

	"github.com/gin-gonic/gin"
)

// @Summary 物品分类
// @Description
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_goods.GoodsTypesListRes{}
// @Router /v1/goods/types [get]
func GoodsTypesList(c *gin.Context) {
	GoodsUseLogic := new(logic.GoodsUseLogic)
	resp := GoodsUseLogic.GoodsTypesList(c)
	response.SuccessResponse(c, resp)
}

// @Summary 根据分类读取商品列表
// @Description
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Param  typeId query int  true "物品类型id"
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_goods.GoodsListByTypesRes{}
// @Router /v1/goods/type_list [get]
func GoodsListByTypes(c *gin.Context) {
	typeId := c.Query("typeId")
	if typeId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	newTypeId, err := strconv.Atoi(typeId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	GoodsUseLogic := new(logic.GoodsUseLogic)
	resp := GoodsUseLogic.GoodsListByTypes(c, newTypeId)
	response.SuccessResponse(c, resp)
}

// @Summary 我的装扮
// @Description
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_goods.GoodsListToUserRes{}
// @Router /v1/goods/goods_center [get]
func GoodsListToUser(c *gin.Context) {
	GoodsUseLogic := new(logic.GoodsUseLogic)
	resp := make([]*response_goods.GoodsListToUserRes, 0)
	resp = append(resp, GoodsUseLogic.GoodsListToUser(c)...)
	response.SuccessResponse(c, resp)
}

// @Summary 使用取消装扮
// @Description
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Param  goodId query int  true "商品ID"
// @Success 0 {object} response.Response{}
// @Router /v1/goods/use_goods [get]
func UseGoodsToUser(c *gin.Context) {
	goodId := c.Query("goodId")
	if goodId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	GoodsUseLogic := new(logic.GoodsUseLogic)
	GoodsUseLogic.UseGoodsToUser(c, goodId)
	response.SuccessResponse(c, nil)
}

// @Summary 删除过期装扮
// @Description
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Param goodId query int true "商品ID"
// @Success 0 {object} response.Response{}
// @Router /v1/goods/del_goods [get]
func UserGoodsDel(c *gin.Context) {
	goodId := c.Query("goodId")
	if goodId == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	GoodsUseLogic := new(logic.GoodsUseLogic)
	GoodsUseLogic.UserGoodsDel(c, goodId)
	response.SuccessResponse(c, nil)
}

// @Summary 执行购买装扮商品
// @Description
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Param  req body  request_goods.BuyGoodsToUserReq  true "购买参数"
// @Success 0 {object} response.Response{}
// @Router /v1/goods/buy [post]
func ByGoodsToUser(c *gin.Context) {
	req := new(request_goods.BuyGoodsToUserReq)
	handle.BindBody(c, req)
	GoodsUseLogic := new(logic.GoodsUseLogic)
	GoodsUseLogic.BuyGoodsToUser(c, req)
	response.SuccessResponse(c, nil)
}

// @Summary 所有分类带商品
// @Description
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_goods.GoodsAllRes{}
// @Router /v1/goods/list [get]
func GoodsAll(c *gin.Context) {
	GoodsUseLogic := new(logic.GoodsUseLogic)
	resp := make([]*response_goods.GoodsAllRes, 0)
	resp = append(resp, GoodsUseLogic.GoodsAll(c)...)
	response.SuccessResponse(c, resp)
}

// @Summary 清除红点
// @Description
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Param goodsId query int true "商品ID"
// @Success 0 {object} response.Response{}
// @Router /v1/goods/del_reddot [get]
func DelRedDot(c *gin.Context) {
	goodsId := c.Query("goodsId")
	if goodsId == "" {
		panic(error2.I18nError{
			Code: i18n_err.ErrorCodeParam,
			Msg:  nil,
		})
	}
	GoodsUseLogic := new(logic.GoodsUseLogic)
	GoodsUseLogic.DelRedDot(c, goodsId)
	response.SuccessResponse(c, nil)
}

// @Summary 装扮中心的用户信息
// @Description 获取用户的信息
// @Tags 装扮中心
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} h5.UserAccountRes{}
// @Router /v1/goods/goods_center_user [get]
func GetGoodsUserAccounts(c *gin.Context) {
	GoodsUseLogic := new(logic.GoodsUseLogic)

	res := GoodsUseLogic.GetUserGoodsAccount(c)
	response.SuccessResponse(c, res)
}
