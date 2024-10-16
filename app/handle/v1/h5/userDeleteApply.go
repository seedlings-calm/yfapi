package h5

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	logic_user "yfapi/internal/logic"
	request_user "yfapi/typedef/request/h5"
	"yfapi/typedef/response"
)

func DeleteUserApply(c *gin.Context) {
	req := new(request_user.UserDeleteApplyReq)
	handle.BindBody(c, req)
	service := new(logic_user.UserDeleteApply)
	service.DeleteUserApply(c, req)
	response.SuccessResponse(c, "")
}
