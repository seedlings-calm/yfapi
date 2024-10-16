package user

import (
	"yfapi/app/handle"
	logic_user "yfapi/internal/logic"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetTimelineDetail(c *gin.Context) {
	timelineId, _ := c.GetQuery("timelineId")
	service := new(logic_user.Timeline)
	response.SuccessResponse(c, service.GetTimelineDetail(c, cast.ToInt64(timelineId)))
}

// GetTimelineListByType
//
//	@Description: 按分类获取动态列表
func GetTimelineListByType(c *gin.Context) {
	req := new(request_user.TimelineListByTypeReq)
	handle.BindQuery(c, req)
	service := new(logic_user.Timeline)
	response.SuccessResponse(c, service.GetTimelineListByType(c, req))
}

// PublishTimeline
//
//	@Description: 发布动态
func PublishTimeline(c *gin.Context) {
	req := new(request_user.TimelinePublishReq)
	handle.BindBody(c, req)
	service := new(logic_user.Timeline)
	service.PublishTimeline(c, req)
	response.SuccessResponse(c, "")
}

// GetUserTimelineList
//
//	@Description: 获取用户动态列表
func GetUserTimelineList(c *gin.Context) {
	req := new(request_user.TimelineListReq)
	handle.BindQuery(c, req)
	service := new(logic_user.Timeline)
	response.SuccessResponse(c, service.GetUserTimelineList(c, req))
}

// DeleteTimeline
//
//	@Description: 删除动态
func DeleteTimeline(c *gin.Context) {
	timelineId, _ := c.GetQuery("timelineId")
	service := new(logic_user.Timeline)
	service.DeleteTimeline(c, cast.ToInt64(timelineId))
	response.SuccessResponse(c, "")
}

// PraiseTimeline
//
//	@Description: 动态点赞
func PraiseTimeline(c *gin.Context) {
	timelineId, _ := c.GetQuery("timelineId")
	service := new(logic_user.Timeline)
	service.PraiseTimeline(c, cast.ToInt64(timelineId))
	response.SuccessResponse(c, "")
}

// CancelPraiseTimeline
//
//	@Description: 动态取消点赞
func CancelPraiseTimeline(c *gin.Context) {
	timelineId, _ := c.GetQuery("timelineId")
	service := new(logic_user.Timeline)
	service.CancelPraiseTimeline(c, cast.ToInt64(timelineId))
	response.SuccessResponse(c, "")
}

// GetPraiseUserList
//
//	@Description: 获取点赞列表
func GetPraiseUserList(c *gin.Context) {
	req := new(request_user.TimelineReplyListReq)
	handle.BindQuery(c, req)
	service := new(logic_user.Timeline)
	response.SuccessResponse(c, service.GetPraiseUserList(c, req))
}

// ReplyTimeline
//
//	@Description: 动态评论
func ReplyTimeline(c *gin.Context) {
	req := new(request_user.TimelineReplyReq)
	handle.BindBody(c, req)
	service := new(logic_user.Timeline)
	service.ReplyTimeline(c, req)
	response.SuccessResponse(c, "")
}

// DeleteReplyTimeline
//
//	@Description: 动态删除评论
func DeleteReplyTimeline(c *gin.Context) {
	replyId, _ := c.GetQuery("replyId")
	service := new(logic_user.Timeline)
	service.DeleteReplyTimeline(c, cast.ToInt64(replyId))
	response.SuccessResponse(c, "")
}

// GetTimelineReplyList
//
//	@Description: 获取评论列表
func GetTimelineReplyList(c *gin.Context) {
	req := new(request_user.TimelineReplyListReq)
	handle.BindQuery(c, req)
	service := new(logic_user.Timeline)
	response.SuccessResponse(c, service.GetTimelineReplyList(c, req))
}

// GetTimelineSubReplyList
//
//	@Description: 获取子评论列表
func GetTimelineSubReplyList(c *gin.Context) {
	req := new(request_user.TimelineSubReplyListReq)
	handle.BindQuery(c, req)
	service := new(logic_user.Timeline)
	response.SuccessResponse(c, service.GetTimelineSubReplyList(c, req))
}

// PraiseTimelineReply
//
//	@Description: 评论点赞
func PraiseTimelineReply(c *gin.Context) {
	replyId, _ := c.GetQuery("replyId")
	service := new(logic_user.Timeline)
	service.PraiseTimelineReply(c, cast.ToInt64(replyId))
	response.SuccessResponse(c, "")
}

// CancelPraiseTimelineReply
//
//	@Description: 评论取消点赞
func CancelPraiseTimelineReply(c *gin.Context) {
	replyId, _ := c.GetQuery("replyId")
	service := new(logic_user.Timeline)
	service.CancelPraiseTimelineReply(c, cast.ToInt64(replyId))
	response.SuccessResponse(c, "")
}
