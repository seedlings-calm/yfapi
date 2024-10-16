package h5

import (
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	logic_user "yfapi/internal/logic"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
	_ "yfapi/typedef/response/user"

	"github.com/gin-gonic/gin"
)

// @Summary 获取考核题目接口
// @Description (20题基础题)或者 （2到简答题）
// @Tags H5
// @Accept json
// @Produce json
// @Param  types path int  true  "身份 1：主持人，3：咨询师,4:主播"
// @Success 0 {object} response.Response{}
// @Success 0 {object} []user.PractitionerExamineQuestion{}
// @Router /v1/h5/getQuestion/{types} [get]
func GetQuestion(c *gin.Context) {
	types := c.Param("types")
	if len(types) == 0 || types == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := new(logic_user.Practitioner)
	data := service.GetQuestion(types, c)
	response.SuccessResponse(c, data)

}

// @Summary 上传客户的2O题考核答案
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param 	data body request_user.PractitionerAwnserReq  true "考核答案"
// @Success 0 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Router /v1/h5/pullAnswer [post]
func PullAnswer(c *gin.Context) {

	var req request_user.PractitionerAwnserReq
	handle.BindJson(c, &req)
	service := new(logic_user.Practitioner)
	resp := service.PullAnswer(&req, c)
	response.SuccessResponse(c, resp)
}

// @Summary 上传客户的简答 答案
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  data body  request_user.PractitionerShortAwnserReq true "简答答案参数"
// @Success 0 {object} response.Response{}
// @Failure 400 {object} response.Response{}
// @Router   /v1/h5/pullShortAnswer [post]
func PullShortAnswer(c *gin.Context) {

	var req request_user.PractitionerShortAwnserReq
	handle.BindJson(c, &req)
	service := new(logic_user.Practitioner)
	resp := service.PullShortAnswer(req, c)
	response.SuccessResponse(c, resp)

}

// @Summary 考核音乐人信息提交
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  data body request_user.PractitionerMusicianReq  true "考核信息"
// @Success 0 {object} response.Response{}
// @Router  /v1/h5/pullMusic [post]
func PullMusic(c *gin.Context) {
	req := &request_user.PractitionerMusicianReq{}
	handle.BindBody(c, req)
	service := new(logic_user.Practitioner)
	resp := service.PullMusic(req, c)
	response.SuccessResponse(c, resp)
}

// @Summary 申请从业者资格
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Param  types path int  true  "身份 1：主持人，2：音乐人，3：咨询师,4:主播"
// @Success 0 {object} response.Response{}
// @Success 0 {object} user.CerdAuthResponse{}
// @Router /v1/h5/cerdAuth/{types} [get]
func CerdAuth(c *gin.Context) {
	types := c.Param("types")
	if len(types) == 0 || types == "" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := new(logic_user.PractitionerCerd)
	resp := service.CerdAuth(types, c)
	response.SuccessResponse(c, resp)

}

// @Summary 申请条件是否满足情况
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} user.ApplyJoinResultResponse{}
// @Router /v1/h5/applyJoinResult [get]
func ApplyJoinResult(c *gin.Context) {
	service := new(logic_user.PractitionerCerd)
	resp := service.ApplyJoinResult(c)
	response.SuccessResponse(c, resp)

}

// @Summary 获取用户的基本信息
// @Description
// @Tags H5
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object}  user.UserH5BasicInfo{}
// @Router /v1/h5/getUserInfo [get]
func GetUserInfo(c *gin.Context) {
	service := new(logic_user.UserInfo)
	resp := service.GetUserBasicInfo(handle.GetUserId(c))
	response.SuccessResponse(c, resp)
}
