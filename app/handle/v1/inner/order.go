package inner

import (
	"log"
	"yfapi/app/handle"
	"yfapi/internal/logic"
	request_inner "yfapi/typedef/request/inner"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
)

// @Summary accountChange
// @Description
// @Tags 内部调用
// @Accept json
// @Produce json
// @Param  req body request_inner.AccountChangeReq  true "参数"
// @Success 0 {object} response.Response{}
// @Router /v1/inner/accountChange [post]
func AccountChange(c *gin.Context) {
	req := new(request_inner.AccountChangeReq)
	handle.BindBody(c, req)
	userAccountLogic := new(logic.UserAccount)
	log.Println(req)
	response.SuccessResponse(c, userAccountLogic.AccountChangeToAdmin(c, req))
}

// @Summary 运营后台变动用户账户
// @Description
// @Tags 内部调用
// @Accept json
// @Produce json
// @Param  req body request_inner.OperationChangeAccountReq  true "参数"
// @Success 0 {object} response.Response{}
// @Router /v1/inner/operationChangeAccount [post]
func OperationChangeAccount(c *gin.Context) {
	req := new(request_inner.OperationChangeAccountReq)
	handle.BindBody(c, req)
	response.SuccessResponse(c, new(logic.UserAccount).OperationChangeAccount(c, req))
}
