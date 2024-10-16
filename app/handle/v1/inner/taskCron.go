package inner

import (
	"github.com/gin-gonic/gin"
	"yfapi/internal/logic"
	"yfapi/typedef/response"
)

func AutoDeleteUser(context *gin.Context) {
	service := new(logic.TaskCron)
	service.AutoDeleteUser(context)
	response.SuccessResponse(context, "ok")
}
