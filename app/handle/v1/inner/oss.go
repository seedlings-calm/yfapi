package inner

import (
	"github.com/gin-gonic/gin"
	"yfapi/internal/logic"
	"yfapi/typedef/response"
)

// GetUploadToken 获取上传oss token
func GetUploadToken(context *gin.Context) {
	service := new(logic.OssStsToken)
	response.SuccessResponse(context, service.GetOssStsToken(context))
}
