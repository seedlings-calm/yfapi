package user

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	logic_user "yfapi/internal/logic"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
)

// 添加关注
func AddFollow(context *gin.Context) {
	req := new(request_user.AddFollowReq)
	handle.BindBody(context, req)
	service := new(logic_user.Follow)
	response.SuccessResponse(context, service.AddFollowUser(req, context))
}

// 取消关注
func RemoveFollow(context *gin.Context) {
	req := new(request_user.RemoveFollowReq)
	handle.BindBody(context, req)
	service := new(logic_user.Follow)
	response.SuccessResponse(context, service.RemoveFollowUser(req, context))
}

// 获取关注列表
func GetUserFollowingList(context *gin.Context) {
	req := new(request_user.GetUserFollowingListReq)
	handle.BindBody(context, req)
	service := new(logic_user.Follow)
	response.SuccessResponse(context, service.GetUserFollowingUserList(req, context))
}

// 获取粉丝列表
func GetFollowersList(context *gin.Context) {
	req := new(request_user.GetUserFollowingListReq)
	handle.BindBody(context, req)
	service := new(logic_user.Follow)
	response.SuccessResponse(context, service.GetUserFollowersList(req, context))
}

// 获取好友列表
func GetFriendsList(context *gin.Context) {
	req := new(request_user.GetUserFriendsListReq)
	handle.BindBody(context, req)
	service := new(logic_user.Follow)
	response.SuccessResponse(context, service.GetUserFriendList(req, context))
}

// 删除粉丝
func DeleteFans(context *gin.Context) {
	req := new(request_user.DeleteFansReq)
	handle.BindBody(context, req)
	service := new(logic_user.Follow)
	response.SuccessResponse(context, service.DeleteFansUser(req, context))
}
