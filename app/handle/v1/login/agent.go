package login

import (
	"github.com/gin-gonic/gin"
	error2 "yfapi/i18n/error"
	"yfapi/internal/logic"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/response"
)

// 授权token
func H5Token(c *gin.Context) {
	req := map[string]any{}
	if err := c.ShouldBind(&req); err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := &logic.Agent{ClientType: typedef_enum.ClientTypeH5}
	response.SuccessResponse(c, service.OauthToken(c, req))
}
