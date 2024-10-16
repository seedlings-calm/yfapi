package user

import (
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	logic_user "yfapi/internal/logic"
	"yfapi/internal/model"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
)

// PerfectInfo
//
//	@Summary	完善用户信息
//	@Schemes
//	@Description	完善用户信息
//	@Tags			用户相关
//	@Param			req	body	request_user.PerfectInfoReq	true	"完善信息参数"
//	@Accept			json
//	@Produce		json
//	@Router			/v1/user/perfectInfo [post]
func PerfectInfo(context *gin.Context) {
	req := new(request_user.PerfectInfoReq)
	handle.BindBody(context, req)
	service := new(logic_user.UserInfo)
	service.PerfectInfo(req, context)
	response.SuccessResponse(context, "")
}

// ImServer
//
//	@Summary	获取im服务器
//	@Schemes
//	@Description	获取可连接得im服务器
//	@Tags			用户相关
//	@Accept			json
//	@Produce		json
//	@Router			/v1/user/imserver [get]
func ImServer(context *gin.Context) {
	service := new(logic_user.UserInfo)
	response.SuccessResponse(context, service.ImServer(context))
}

// CheckRepeatName
//
//	@Summary	检测重复用户名
//	@Schemes
//	@Description	检测重复用户名
//	@Tags			用户相关
//	@Param nickname	query string	true	"用户昵称"
//	@Accept			json
//	@Produce		json
//	@Router			/v1/user/checkRepeatName [get]
func CheckRepeatName(c *gin.Context) {
	nickname, _ := c.GetQuery("nickname")
	if len(nickname) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	userDao := new(dao.UserDao)
	count := userDao.Count(&model.User{
		Nickname: nickname,
	})
	if count > 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNicknameRepeat,
			Msg:  nil,
		})
	}
	response.SuccessResponse(c, "")
}

// EditUserInfo
//
//	@Description: 修改用户信息
func EditUserInfo(context *gin.Context) {
	req := new(request_user.EditUserInfoReq)
	handle.BindBody(context, req)
	service := new(logic_user.UserInfo)
	resp := service.EditUserInfo(req, context)
	response.SuccessResponse(context, resp)
}

// GetUserInfo
//
//	@Description: 获取用户信息
func GetUserInfo(context *gin.Context) {
	req := new(request_user.UserInfoReq)
	handle.BindQuery(context, req)
	service := new(logic_user.UserInfo)
	resp := service.GetUserInfo(req, context)
	response.SuccessResponse(context, resp)
}

func SearchUserInfo(c *gin.Context) {
	keyword, _ := c.GetQuery("keyword")
	service := new(logic_user.UserInfo)
	response.SuccessResponse(c, service.SearchUserInfo(c, keyword))
}

// 根据userNo获取用户信息
func SearchUserInfoByUserNo(c *gin.Context) {
	userNo, _ := c.GetQuery("userNo")
	service := new(logic_user.UserInfo)
	response.SuccessResponse(c, service.SearchUserInfoByUserNo(c, userNo))
}

// 获取用户实名认证信息
func GetUserRealNameInfo(c *gin.Context) {
	res := new(logic_user.UserInfo).GerUserRealNameInfo(c)
	response.SuccessResponse(c, res)
}
