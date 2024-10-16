package inner

import (
	"yfapi/app/handle"
	"yfapi/internal/logic"
	request_recharge "yfapi/typedef/request/recharge"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
)

// 代付
func AnotherPay(c *gin.Context) {
	req := new(request_recharge.AnotherPayReq)
	handle.BindBody(c, req)
	ok, err := new(logic.Recharge).AnotherPay(c, req)
	response.SuccessResponse(c, map[string]interface{}{
		"status": ok,
		"msg":    err,
	})
}
