package h5

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	logic_user "yfapi/internal/logic"
	request_h5 "yfapi/typedef/request/h5"
	"yfapi/typedef/response"
	response_h5 "yfapi/typedef/response/h5"
)

// @Summary 等级基本信息
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  req body request_h5.LevelBaseInfoReq   true "请求参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_h5.LevelBaseInfoRes{}
// @Router /v1/h5/level/base [post]
func LevelBaseInfo(c *gin.Context) {
	req := new(request_h5.LevelBaseInfoReq)
	handle.BindBody(c, req)
	res := new(response_h5.LevelBaseInfoRes)
	res = new(logic_user.LevelConfig).GetLevelBaseInfo(c, req)
	response.SuccessResponse(c, res)
}

// @Summary 等级配置列表
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  req body request_h5.LevelBaseInfoReq   true "请求参数"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_h5.LevelConfigListRes{}
// @Router /v1/h5/level/config [post]
func GetLevelConfigList(c *gin.Context) {
	req := new(request_h5.LevelBaseInfoReq)
	handle.BindBody(c, req)
	res := response_h5.LevelConfigListRes{}
	res = new(logic_user.LevelConfig).GetLevelConfigList(c, req)
	response.SuccessResponse(c, res)
}
