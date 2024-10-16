package inner

import (
	"yfapi/app/handle"
	"yfapi/internal/logic"
	request_login "yfapi/typedef/request/login"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
)

// @Summary 发送短信
// @Description 内部调用
// @Tags 内部调用
// @Accept json
// @Produce json
// @Param data body request_login.SendSmsReq  true "请求参数"
// @Success 0 {object} response.Response{}
// @Router /v1/inner/sendSms  [post]
func SendSms(context *gin.Context) {
	req := new(request_login.SendSmsReq)
	handle.BindBody(context, req)
	service := logic.Sms{
		Mobile:     req.Mobile,
		Type:       req.Type,
		RegionCode: req.RegionCode,
	}
	response.SuccessResponse(context, service.SendSms(context))
}
