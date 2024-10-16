package resource

import (
	error2 "yfapi/i18n/error"
	"yfapi/internal/logic"
	"yfapi/typedef/response"
	response_resource "yfapi/typedef/response/resource"

	"github.com/gin-gonic/gin"
)

// ResourceList
//
// @Summary 资源列表
// @Schemes
// @Description 资源列表
// @Tags 资源相关
// @Produce json
// @Success 200 {object} response_resource.ResourceListRes
// @Router /v1/resource/list [get]
func ResourceList(c *gin.Context) {
	data := response_resource.ResourceListRes{
		FileUrl:  "https://photo.storage.swdws.com/vchatline.com/assets/img/op_1697709818.mp4",
		FileName: "测试vap",
		FileType: "vap",
	}
	var result []response_resource.ResourceListRes
	result = append(result, data)
	response.SuccessResponse(c, result)
}

// GetUploadToken
//
// @Summary 获取上传token
// @Schemes
// @Description 获取上传token
// @Tags 资源相关
// @Produce json
// @Param type	query string	true	"上传类型 img图片,video视频,audio音频"
// @Router /v1/resource/getUploadToken [get]
func GetUploadToken(c *gin.Context) {
	uploadType, _ := c.GetQuery("type")
	if len(uploadType) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := new(logic.OssStsToken)
	switch uploadType {
	case "img":
		response.SuccessResponse(c, service.GetOssStsToken(c))
		break
	case "video":
		response.SuccessResponse(c, service.GetOssStsToken(c))
		break
	case "audio":
		response.SuccessResponse(c, service.GetOssStsToken(c))
		break
	default:
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}

}

// GetUploadTokenToGuild
// @Summary 获取上传token
// @Schemes
// @Description 获取上传token
// @Tags 公会后台
// @Produce json
// @Param type	query string	true	"上传类型 img图片,video视频,audio音频"
// @Router /v1/guild/getUploadToken [get]
func GetUploadTokenToGuild(c *gin.Context) {
	uploadType, _ := c.GetQuery("type")
	if len(uploadType) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	service := new(logic.OssStsToken)
	switch uploadType {
	case "img":
		response.SuccessResponse(c, service.GetOssStsToken(c))
		break
	case "video":
		response.SuccessResponse(c, service.GetOssStsToken(c))
		break
	case "audio":
		response.SuccessResponse(c, service.GetOssStsToken(c))
		break
	default:
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}

}
