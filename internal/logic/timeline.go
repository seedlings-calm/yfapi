package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_im "yfapi/internal/service/im"
	"yfapi/internal/service/riskCheck/shumei"
	service_user "yfapi/internal/service/user"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/message"
	"yfapi/typedef/redisKey"
	"yfapi/typedef/request/user"
	"yfapi/typedef/response"
	response_user "yfapi/typedef/response/user"
)

type Timeline struct {
}

// GetTimelineDetail
//
//	@Description: 查询正常的动态详情
//	@receiver t
//	@param c *gin.Context -
//	@param timelineId int64 -
//	@return res -
func (t *Timeline) GetTimelineDetail(c *gin.Context, timelineId int64) (res *response_user.TimelineInfo) {
	userId := helper.GetUserId(c)
	timelineModel, err := new(dao.TimelineDao).GetTimelineDetail(timelineId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	if timelineModel.Id == 0 {
		return nil
	}
	res = genTimelineInfo(c, timelineModel, userId, helper.GetClientType(c))
	return
}

// GetTimelineListByType
//
//	@Description: 按分类获取动态列表
//	@receiver t
//	@param c *gin.Context -
//	@param req *user.TimelineListByTypeReq -
//	@return res -
func (t *Timeline) GetTimelineListByType(c *gin.Context, req *user.TimelineListByTypeReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}

	var modelList []model.Timeline
	var count int64
	timelineDao := &dao.TimelineDao{}
	switch req.CategoryId {
	case typedef_enum.TimelineListTypeLatest: // 最新
		modelList, count, err = timelineDao.GetTimelineListLatest(userId, req.Page, req.Size)
	case typedef_enum.TimelineListTypeFollow: // 关注
		modelList, count, err = timelineDao.GetTimelineListFollow(userId, req.Page, req.Size)
		// TODO
	default: // 默认精选
		// TODO
	}
	var result []*response_user.TimelineInfo
	for _, data := range modelList {
		info := genTimelineInfo(c, data, userId, helper.GetClientType(c))
		info.UserPlaque.HeadList = []response.PlaqueInfo{}
		result = append(result, info)
	}
	res.Total = count
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Data = result
	res.CalcHasNext()
	return
}

// GetUserTimelineList
//
//	@Description: 获取用户动态列表
//	@receiver t
//	@param c *gin.Context -
//	@param req *user.TimelineListReq -
//	@return res -
func (t *Timeline) GetUserTimelineList(c *gin.Context, req *user.TimelineListReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}

	modelList, count, err := new(dao.TimelineDao).GetUserTimelineList(req.TargetUserId, userId, req.Page, req.Size, userId == req.TargetUserId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	var result []*response_user.TimelineInfo
	for _, data := range modelList {
		result = append(result, genTimelineInfo(c, data, userId, helper.GetClientType(c)))
	}
	res.Total = count
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Data = result
	res.CalcHasNext()
	return
}

// PublishTimeline
//
//	@Description: 发布动态
//	@receiver t
//	@param c *gin.Context -
//	@param req *user.TimelinePublishReq -
func (t *Timeline) PublishTimeline(c *gin.Context, req *user.TimelinePublishReq) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}

	if len(req.TextContent) != 0 {
		if ok := new(shumei.ShuMei).MomentsCheck(userId, req.TextContent); !ok {
			panic(error2.I18nError{
				Code: error2.ErrorCodeTextCheckReject,
				Msg:  nil,
			})
		}
	}
	status := typedef_enum.TimelineStatusNormal
	var imgList []response_user.TimelineImgDTO
	if req.ContentType == typedef_enum.TimelineImgType {
		if len(req.ImgList) > 0 {
			err := json.Unmarshal([]byte(req.ImgList), &imgList)
			if err != nil {
				coreLog.Error("PublishTimeline ImgList unmarshal err:%+v", err)
				panic(error2.I18nError{
					Code: error2.ErrorCodeParam,
					Msg:  nil,
				})
			}
			for i := range imgList {
				if len(imgList[i].Height) == 0 {
					imgList[i].Height = "0"
				}
				if len(imgList[i].Width) == 0 {
					imgList[i].Width = "0"
				}
			}
			if coreConfig.GetHotConf().RiskSwitch {
				status = typedef_enum.TimelineStatusHidden
			}
		}
	}
	var videoDto = &response_user.TimelineVideoDTO{}
	if req.ContentType == typedef_enum.TimelineVideoType { //视频
		if len(req.VideoDTO) > 0 {
			err = json.Unmarshal([]byte(req.VideoDTO), videoDto)
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeParam,
					Msg:  nil,
				})
			}
			if len(videoDto.Width) == 0 {
				videoDto.Width = "0"
			}
			if len(videoDto.Height) == 0 {
				videoDto.Height = "0"
			}
			if len(videoDto.Duration) == 0 {
				videoDto.Duration = "0"
			}
			if coreConfig.GetHotConf().RiskSwitch {
				status = typedef_enum.TimelineStatusHidden
			}
		}
	}

	param := &model.Timeline{
		ContentType: req.ContentType,
		UserId:      userId,
		TextContent: req.TextContent,
		Status:      status,
		LoveCount:   0,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		CityName:    req.CityName,
		AddressName: req.AddressName,
		ImgList:     req.ImgList,
		VideoData:   req.VideoDTO,
		IsTop:       false,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}
	timelineDao := &dao.TimelineDao{}
	err = timelineDao.Create(param)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	if len(req.ImgList) > 0 && req.ContentType == typedef_enum.TimelineImgType && coreConfig.GetHotConf().RiskSwitch {
		images := []shumei.ImagesAsyncCheckReqImgs{}
		for _, v := range imgList {
			images = append(images, shumei.ImagesAsyncCheckReqImgs{
				BtId: coreSnowflake.GetSnowId(),
				Img:  helper.FormatImgUrl(v.ImgPhotoKey),
			})
		}
		new(shumei.ShuMei).MomentsImageAsyncCheck(userId, images, cast.ToString(param.Id))
	}
	if videoDto != nil && req.ContentType == typedef_enum.TimelineVideoType && coreConfig.GetHotConf().RiskSwitch {
		new(shumei.ShuMei).MomentsVideoAsyncCheck(userId, helper.FormatImgUrl(videoDto.VideoUrl), cast.ToString(param.Id), coreSnowflake.GetSnowId())
	}
	if status == typedef_enum.TimelineStatusNormal {
		//正常发布 通知用户
		go new(Notice).MomentsPublishNotice(c, userId, param.Id)
	}
	return
}

// DeleteTimeline
//
//	@Description: 删除动态
//	@receiver t
//	@param c *gin.Context -
//	@param timelineId int64 -
func (t *Timeline) DeleteTimeline(c *gin.Context, timelineId int64) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	timelineDao := new(dao.TimelineDao)
	_ = timelineDao
	timelineModel, err := timelineDao.GetTimelineById(timelineId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if timelineModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	if timelineModel.UserId != userId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	err = timelineDao.Delete(&model.Timeline{Id: timelineId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}

	// 删除动态的点赞列表
	coreRedis.GetUserRedis().Del(c, redisKey.TimelinePraisedUserList(timelineId))

	// 删除所有评论、子评论点赞列表
	timelineReplyDao := &dao.TimelineReplyDao{}
	replyModelList, _ := timelineReplyDao.GetAllTimelineReplyList(timelineId)
	var replyPraisedKey []string
	for _, replyModel := range replyModelList {
		// 删除点赞列表
		replyPraisedKey = append(replyPraisedKey, redisKey.TimelineReplyPraisedUserList(replyModel.Id))
	}
	if len(replyPraisedKey) > 0 {
		coreRedis.GetUserRedis().Del(c, replyPraisedKey...)
	}

	// 删除所有评论和子评论
	_ = timelineReplyDao.DeleteReplyByTimelineId(timelineId)
}

// PraiseTimeline
//
//	@Description: 动态点赞
//	@receiver t
//	@param c *gin.Context -
//	@param timelineId int64 -
func (t *Timeline) PraiseTimeline(c *gin.Context, timelineId int64) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	timelineDao := new(dao.TimelineDao)
	_ = timelineDao
	timelineModel, err := timelineDao.GetTimelineById(timelineId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if timelineModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	// 是否点赞过
	if checkUserIsPraisedTimeline(userId, timelineId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeIsPraiseTimeline,
			Msg:  nil,
		})
	}
	// 增加点赞数
	err = timelineDao.IncrLoveCount(timelineId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 加入到点赞列表
	member := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: userId,
	}
	_ = coreRedis.GetUserRedis().ZAdd(c, redisKey.TimelinePraisedUserList(timelineId), member).Err()
	// 给动态作者发互动消息
	if userId != timelineModel.UserId {
		timelineInfo := genTimelineInfo(c, timelineModel, userId, helper.GetClientType(c))
		msg := i18n_msg.GetI18nMsg(c, i18n_msg.LikedYourPostMsgKey, map[string]any{"nickname": userModel.Nickname})
		nowTime := time.Now().Format(time.DateTime)
		//发送点赞通知
		imageUrl := ""
		videoUrl := ""
		if timelineInfo.ContentType == typedef_enum.TimelineImgType {
			if len(timelineInfo.ImgDTOList) > 0 {
				imageUrl = timelineInfo.ImgDTOList[0].ImgUrl
			}
		}
		if timelineInfo.ContentType == typedef_enum.TimelineVideoType {
			if timelineInfo.VideoDTO != nil {
				videoUrl = timelineInfo.VideoDTO.VideoUrl + typedef_enum.VideoCoverImgSuffix
			}
		}
		go new(service_im.ImNoticeService).SendInteractiveMsg(c, message.SendInteractiveMsg{
			Avatar:          helper.FormatImgUrl(userModel.Avatar),
			NickName:        userModel.Nickname,
			PraiseIcon:      "",
			CreateTime:      nowTime,
			ImageUrl:        imageUrl,
			VideoUrl:        videoUrl,
			Msg:             msg,
			MsgType:         typedef_enum.DynamicMsgTypeLike,
			TimelineId:      timelineInfo.TimelineId,
			TimelineContent: timelineInfo.TextContent,
			UserPlaque:      service_user.GetUserLevelPlaque(userModel.Id, helper.GetClientType(c)),
		}, []string{timelineModel.UserId})
	}
	return
}

// CancelPraiseTimeline
//
//	@Description: 动态取消点赞
//	@receiver t
//	@param c *gin.Context -
//	@param timelineId int64 -
func (t *Timeline) CancelPraiseTimeline(c *gin.Context, timelineId int64) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	timelineDao := new(dao.TimelineDao)
	_ = timelineDao
	timelineModel, err := timelineDao.GetTimelineById(timelineId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if timelineModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	// 是否点赞过
	if !checkUserIsPraisedTimeline(userId, timelineId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotPraiseTimeline,
			Msg:  nil,
		})
	}
	// 减少点赞数
	err = timelineDao.DecrLoveCount(timelineId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 移除点赞列表
	_ = coreRedis.GetUserRedis().ZRem(c, redisKey.TimelinePraisedUserList(timelineId), userId).Err()
	return
}

// ReplyTimeline
//
//	@Description: 动态评论
//	@receiver t
//	@param c *gin.Context -
//	@param req *user.TimelineReplyReq -
//	@return res -
func (t *Timeline) ReplyTimeline(c *gin.Context, req *user.TimelineReplyReq) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	timelineDao := new(dao.TimelineDao)
	timelineModel, err := timelineDao.GetTimelineById(req.TimelineId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if timelineModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	// TODO 是否被拉黑
	// code.ErrorCodeIsBlacklist

	if len(req.ReplyContent) != 0 {
		if ok := new(shumei.ShuMei).CommentCheck(userId, req.ReplyContent); !ok {
			panic(error2.I18nError{
				Code: error2.ErrorCodeTextCheckReject,
				Msg:  nil,
			})
		}
	}

	// 如果评论的是二级评论,ToReplyId设置为一级的评论ID
	toReplyId, toReplierId := int64(0), ""
	if req.ToReplyId > 0 {
		replyModel, err := new(dao.TimelineReplyDao).GetTimelineReplyById(req.ToReplyId)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		if replyModel.ToReplyId > 0 {
			toReplyId = replyModel.ToReplyId
		} else {
			toReplyId = replyModel.Id
		}
		toReplierId = replyModel.ReplierId
	}

	param := &model.TimelineReply{
		ReplierId:    userId,
		TimelineId:   req.TimelineId,
		ToReplyId:    toReplyId,
		ReplyContent: req.ReplyContent,
		Status:       1,
		ToSubReplyId: req.ToReplyId,
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}

	if len(toReplierId) > 0 {
		// TODO 如果被评论人拉黑了评论人则不能评论
		// code.ErrorCodeIsBlacklist
		param.ToReplierId = sql.NullString{String: toReplierId, Valid: true}
	}

	tx := coreDb.GetMasterDb().Begin()
	err = tx.Model(param).Create(param).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}

	// 如果是子评论 增加被评论对象的子评论数
	if toReplyId > 0 {
		err = tx.Model(&model.TimelineReply{Id: toReplyId}).Update("sub_reply_count", gorm.Expr("sub_reply_count+1")).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	} else {
		// 增加动态评论数量
		err = tx.Model(&model.Timeline{Id: req.TimelineId}).Update("reply_count", gorm.Expr("reply_count+1")).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	}
	tx.Commit()

	// 如果是回复评论 给评论人发一条
	if len(toReplierId) > 0 && toReplierId != userId {
		msg := i18n_msg.GetI18nMsg(c, i18n_msg.RepliedYourCommentMsgKey, map[string]any{"nickname": userModel.Nickname})
		_ = msg
		//_ = sendReplyOrPraiseTimelineMsg(timelineModel.UserId, &userModel, &timelineModel, msgTitle, requestModel.ReplyContent, replyModel.CreateTime)
		// TODO　发送通知
		timelineInfo := genTimelineInfo(c, timelineModel, userId, helper.GetClientType(c))
		nowTime := time.Now().Format(time.DateTime)
		imageUrl := ""
		videoUrl := ""
		if timelineInfo.ContentType == typedef_enum.TimelineImgType {
			if len(timelineInfo.ImgDTOList) > 0 {
				imageUrl = timelineInfo.ImgDTOList[0].ImgUrl
			}
		}
		if timelineInfo.ContentType == typedef_enum.TimelineVideoType {
			if timelineInfo.VideoDTO != nil {
				videoUrl = timelineInfo.VideoDTO.VideoUrl + typedef_enum.VideoCoverImgSuffix
			}
		}
		// 给评论人发一条
		go new(service_im.ImNoticeService).SendInteractiveMsg(c, message.SendInteractiveMsg{
			Avatar:          helper.FormatImgUrl(userModel.Avatar),
			NickName:        userModel.Nickname,
			PraiseIcon:      "",
			CreateTime:      nowTime,
			ImageUrl:        imageUrl,
			VideoUrl:        videoUrl,
			Msg:             msg,
			MsgType:         typedef_enum.DynamicMsgTypeMyCommentReply,
			CommentContent:  req.ReplyContent,
			TimelineId:      timelineInfo.TimelineId,
			TimelineContent: timelineInfo.TextContent,
			UserPlaque:      service_user.GetUserLevelPlaque(userModel.Id, helper.GetClientType(c)),
		}, []string{toReplierId})

		//_ = sendReplyOrPraiseTimelineMsg(toReplierId, &userModel, &timelineModel, msgTitle2, requestModel.ReplyContent, replyModel.CreateTime)
		toReplyUserInfo := service_user.GetUserBaseInfo(toReplierId)
		// 给动态作者发互动消息
		go new(service_im.ImNoticeService).SendInteractiveMsg(c, message.SendInteractiveMsg{
			Avatar:          helper.FormatImgUrl(userModel.Avatar),
			NickName:        userModel.Nickname,
			PraiseIcon:      "",
			CreateTime:      nowTime,
			ImageUrl:        imageUrl,
			VideoUrl:        videoUrl,
			Msg:             msg,
			MsgType:         typedef_enum.DynamicMsgTypeOtherCommentReply,
			CommentContent:  req.ReplyContent,
			ReplyUserId:     toReplyUserInfo.Id,
			ReplyNickName:   toReplyUserInfo.Nickname,
			TimelineId:      timelineInfo.TimelineId,
			TimelineContent: timelineInfo.TextContent,
			UserPlaque:      service_user.GetUserLevelPlaque(userModel.Id, helper.GetClientType(c)),
		}, []string{timelineModel.UserId})
		// userId 回复 toReplierId 消息内容
	} else {
		// 给动态作者发互动消息
		if userId != timelineModel.UserId {
			msg := i18n_msg.GetI18nMsg(c, i18n_msg.CommentedYourPostMsgKey, map[string]any{"nickname": userModel.Nickname})
			_ = msg
			//_ = sendReplyOrPraiseTimelineMsg(timelineModel.UserId, &userModel, &timelineModel, msgTitle, requestModel.ReplyContent, replyModel.CreateTime)
			// TODO　发送通知
			timelineInfo := genTimelineInfo(c, timelineModel, userId, helper.GetClientType(c))
			nowTime := time.Now().Format(time.DateTime)
			imageUrl := ""
			videoUrl := ""
			if timelineInfo.ContentType == typedef_enum.TimelineImgType {
				if len(timelineInfo.ImgDTOList) > 0 {
					imageUrl = timelineInfo.ImgDTOList[0].ImgUrl
				}
			}
			if timelineInfo.ContentType == typedef_enum.TimelineVideoType {
				if timelineInfo.VideoDTO != nil {
					videoUrl = timelineInfo.VideoDTO.VideoUrl + typedef_enum.VideoCoverImgSuffix
				}
			}
			go new(service_im.ImNoticeService).SendInteractiveMsg(c, message.SendInteractiveMsg{
				Avatar:          helper.FormatImgUrl(userModel.Avatar),
				NickName:        userModel.Nickname,
				PraiseIcon:      "",
				CreateTime:      nowTime,
				ImageUrl:        imageUrl,
				VideoUrl:        videoUrl,
				Msg:             msg,
				MsgType:         typedef_enum.DynamicMsgTypeMyDynamicComment,
				CommentContent:  req.ReplyContent,
				TimelineId:      timelineInfo.TimelineId,
				TimelineContent: timelineInfo.TextContent,
				UserPlaque:      service_user.GetUserLevelPlaque(userModel.Id, helper.GetClientType(c)),
			}, []string{timelineModel.UserId})
		}
	}

	// TODO　评论红点提示
	return
}

// DeleteReplyTimeline
//
//	@Description: 动态删除评论
//	@receiver t
//	@param c *gin.Context -
//	@param replyId int64 -
func (t *Timeline) DeleteReplyTimeline(c *gin.Context, replyId int64) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}

	timelineReplyDao := &dao.TimelineReplyDao{}
	replyModel, _ := timelineReplyDao.GetTimelineReplyById(replyId)
	if replyModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	if replyModel.ReplierId != userId { // 只能删除自己的评论
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	// 修改状态为删除
	tx := coreDb.GetMasterDb().Begin()
	err = tx.Model(&model.TimelineReply{Id: replyId}).Updates(model.TimelineReply{Status: 3}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}

	if replyModel.ToReplyId > 0 {
		// 减少评论的评论数
		err = tx.Model(&model.TimelineReply{Id: replyModel.Id}).Where("sub_reply_count>=1").Update("sub_reply_count", gorm.Expr("sub_reply_count-1")).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	} else {
		// 减少动态评论数
		err = tx.Model(&model.Timeline{Id: replyModel.TimelineId}).Where("reply_count>=1").Update("reply_count", gorm.Expr("reply_count-1")).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	}
	tx.Commit()
	return
}

// GetTimelineReplyList
//
//	@Description: 获取评论列表
//	@receiver t
//	@param c *gin.Context -
//	@param req *user.TimelineReplyListReq -
//	@return res -
func (t *Timeline) GetTimelineReplyList(c *gin.Context, req *user.TimelineReplyListReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	timelineReplyDao := &dao.TimelineReplyDao{}
	modelList, count, err := timelineReplyDao.GetTimelineReplyList(req.TimelineId, req.Page, req.Size)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	var result []*response_user.TimelineReplyInfo
	for _, data := range modelList {
		dst := genTimelineReplyInfo(data, helper.GetClientType(c))
		dst.IsPraised = checkUserIsPraisedTimelineReply(userId, data.Id)
		if data.ToReplyId == 0 {
			// 获取第一条子评论
			subModelList, _, _ := timelineReplyDao.GetTimelineSubReplyList(req.TimelineId, data.Id, 0, 1)
			if len(subModelList) > 0 {
				var subReplyList []*response_user.TimelineReplyInfo
				for _, sub := range subModelList {
					sDst := genTimelineReplyInfo(sub, helper.GetClientType(c))
					// 是否点赞过
					sDst.IsPraised = checkUserIsPraisedTimelineReply(userId, sub.Id)
					subReplyList = append(subReplyList, sDst)
				}
				dst.SubReplyList = subReplyList
			}

		}
		result = append(result, dst)
	}
	res.Total = count
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Data = result
	res.CalcHasNext()
	return
}

// GetTimelineSubReplyList
//
//	@Description: 获取子评论列表
//	@receiver t
//	@param c *gin.Context -
//	@param req *user.TimelineSubReplyListReq -
//	@return res -
func (t *Timeline) GetTimelineSubReplyList(c *gin.Context, req *user.TimelineSubReplyListReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	timelineReplyDao := &dao.TimelineReplyDao{}
	replyModel, _ := timelineReplyDao.GetTimelineReplyById(req.ReplyId)
	if replyModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	modelList, count, err := timelineReplyDao.GetTimelineSubReplyList(replyModel.TimelineId, req.ReplyId, req.Page, req.Size)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	var result []*response_user.TimelineReplyInfo
	for _, data := range modelList {
		dst := genTimelineReplyInfo(data, helper.GetClientType(c))
		dst.IsPraised = checkUserIsPraisedTimelineReply(userId, data.Id)
		result = append(result, dst)
	}
	res.Total = count
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Data = result
	res.CalcHasNext()
	return
}

// GetPraiseUserList
//
//	@Description: 获取点赞列表
//	@receiver t
//	@param c *gin.Context -
//	@param req *user.TimelineReplyListReq -
//	@return res -
func (t *Timeline) GetPraiseUserList(c *gin.Context, req *user.TimelineReplyListReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}

	start := int64(req.Page * req.Size)
	end := start + int64(req.Size) - 1

	totalCount := coreRedis.GetUserRedis().ZCard(c, redisKey.TimelinePraisedUserList(req.TimelineId)).Val()
	if start >= totalCount {
		return
	}
	var result []response_user.PraisedUserInfo
	if end > totalCount {
		end = totalCount
	}

	zList := coreRedis.GetUserRedis().ZRevRangeWithScores(c, redisKey.TimelinePraisedUserList(req.TimelineId), start, end).Val()

	var userIdList []string
	for _, z := range zList {
		currUserId, _ := z.Member.(string)
		result = append(result, response_user.PraisedUserInfo{UserId: currUserId})
		userIdList = append(userIdList, currUserId)
	}
	userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
	for i, v := range result {
		result[i].UserNo = userInfoMap[v.UserId].UserNo
		result[i].Nickname = userInfoMap[v.UserId].Nickname
		result[i].Avatar = userInfoMap[v.UserId].Avatar
		result[i].Sex = userInfoMap[v.UserId].Sex
		result[i].UserPlaque = service_user.GetUserLevelPlaque(v.UserId, helper.GetClientType(c))
	}

	res.Total = totalCount
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Data = result
	res.CalcHasNext()
	return
}

// PraiseTimelineReply
//
//	@Description: 评论点赞
//	@receiver t
//	@param c *gin.Context -
//	@param replyId int64 -
func (t *Timeline) PraiseTimelineReply(c *gin.Context, replyId int64) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}

	timelineReplyDao := &dao.TimelineReplyDao{}
	replyModel, _ := timelineReplyDao.GetTimelineReplyById(replyId)
	if replyModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	if replyModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	// 是否点赞过
	if checkUserIsPraisedTimelineReply(userId, replyId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeIsPraiseTimeline,
			Msg:  nil,
		})
	}
	// 增加点赞数
	err = timelineReplyDao.IncrSubReplyPraisedCount(replyId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}

	// 加入到点赞列表
	member := redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: userId,
	}
	_ = coreRedis.GetUserRedis().ZAdd(c, redisKey.TimelineReplyPraisedUserList(replyId), member).Err()
	return
}

// CancelPraiseTimelineReply
//
//	@Description: 评论取消点赞
//	@receiver t
//	@param c *gin.Context -
//	@param replyId int64 -
func (t *Timeline) CancelPraiseTimelineReply(c *gin.Context, replyId int64) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}

	timelineReplyDao := &dao.TimelineReplyDao{}
	replyModel, _ := timelineReplyDao.GetTimelineReplyById(replyId)
	if replyModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	if replyModel.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDataNotFound,
			Msg:  nil,
		})
	}
	// 是否点赞过
	if !checkUserIsPraisedTimelineReply(userId, replyId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNotPraiseTimeline,
			Msg:  nil,
		})
	}
	// 减少点赞数
	err = timelineReplyDao.DecrSubReplyPraisedCount(replyId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 移除点赞列表
	_ = coreRedis.GetUserRedis().ZRem(c, redisKey.TimelineReplyPraisedUserList(replyId), userId).Err()
	return
}

// genTimelineInfo
//
//	@Description: 包装动态数据
//	@param data model.Timeline -
//	@return dst -
func genTimelineInfo(c *gin.Context, data model.Timeline, userId string, clientType string) (dst *response_user.TimelineInfo) {
	dst = &response_user.TimelineInfo{
		TimelineId:      data.Id,
		ContentType:     data.ContentType,
		UserId:          data.UserId,
		TextContent:     data.TextContent,
		Status:          data.Status,
		LoveCount:       data.LoveCount,
		Latitude:        data.Latitude,
		Longitude:       data.Longitude,
		CityName:        data.CityName,
		AddressName:     data.AddressName,
		ReplyCount:      data.ReplyCount,
		ImgDTOList:      nil,
		VideoDTO:        nil,
		IsTop:           data.IsTop,
		IsPraised:       false,
		CreateTime:      data.CreateTime.Format(time.DateTime),
		UpdateTime:      data.UpdateTime.Format(time.DateTime),
		TimelineTimeStr: getTimelineTimeStr(c, data.CreateTime),
		UserPlaque:      service_user.GetUserLevelPlaque(data.UserId, clientType),
	}

	dst.IsPraised = checkUserIsPraisedTimeline(userId, data.Id)

	if data.UserId != userId {
		dst.IsFollow = new(dao.UserFollowDao).IsUserFollowed(userId, data.UserId)
	}

	// 处理视频相关
	if data.ContentType == 2 {
		var videoDTO response_user.TimelineVideoDTO
		_ = json.Unmarshal([]byte(data.VideoData), &videoDTO)
		videoDTO.VideoUrl = helper.FormatImgUrl(videoDTO.VideoUrl)
		videoDTO.VideoCoverImgUrl = helper.FormatImgUrl(videoDTO.VideoCoverImgUrl)
		dst.VideoDTO = &videoDTO
	} else if data.ContentType == 1 {
		if len(data.ImgList) > 0 {
			var tempList []response_user.TimelineImgDTO
			_ = json.Unmarshal([]byte(data.ImgList), &tempList)
			var imgDTOList []response_user.TimelineImgDTO
			for _, img := range tempList {
				img.ImgUrl = helper.FormatImgUrl(img.ImgPhotoKey)
				img.ImgPhotoKey = ""
				imgDTOList = append(imgDTOList, img)
			}
			dst.ImgDTOList = imgDTOList
		}
	}

	// 补充用户信息
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: data.UserId})
	if err != nil {
		return dst
	}
	dst.UserNo = userModel.UserNo
	dst.Uid32 = cast.ToInt32(userModel.OriUserNo)
	dst.Nickname = userModel.Nickname
	dst.Avatar = helper.FormatImgUrl(userModel.Avatar)
	dst.Sex = userModel.Sex
	dst.UserLastActiveTime = "" // TODO
	return
}

// getTimelineTimeStr
//
//	@Description: 动态时间格式化
//	@param createTime time.Time -
//	@return dst -
func getTimelineTimeStr(c *gin.Context, createTime time.Time) (dst string) {
	year := createTime.Year()
	now := time.Now()
	currYear := now.Year()
	if year != currYear {
		return createTime.Format("2006-01-02 15:04")
	}

	minutes := time.Since(createTime).Minutes()
	switch {
	case minutes > 60*24 || createTime.Day() != now.Day():
		dst = fmt.Sprintf("%v", createTime.Format("01-02 15:04"))
	case minutes > 60:
		dst = fmt.Sprintf("%v", createTime.Format("15:04"))
	case minutes > 5:
		//dst = fmt.Sprintf("%v分钟前", cast.ToInt(minutes))
		dst = i18n_msg.GetI18nMsg(c, i18n_msg.MinutesAgoMsgKey, map[string]any{"minute": cast.ToInt(minutes)})
	default:
		//dst = "刚刚"
		dst = i18n_msg.GetI18nMsg(c, i18n_msg.JustNowMsgKey)
	}
	return
}

func getTimelineReplyTimeStr(createTime time.Time) (dst string) {
	year := createTime.Year()
	now := time.Now()
	currYear := now.Year()
	if year != currYear {
		return createTime.Format("2006")
	}
	if createTime.Day() == now.Day() {
		dst = createTime.Format("15:04")
	} else {
		dst = createTime.Format("01-02")
	}
	return
}

// checkUserIsPraisedTimeline
//
//	@Description: 检查玩家是否点赞过动态
//	@param userId string -
//	@param timelineId int64 -
//	@return bool -
func checkUserIsPraisedTimeline(userId string, timelineId int64) bool {
	_, err := coreRedis.GetUserRedis().ZRank(context.Background(), redisKey.TimelinePraisedUserList(timelineId), userId).Result()
	if err != nil {
		return false
	}
	return true
}

// checkUserIsPraisedTimelineReply
//
//	@Description: 检查玩家是否点赞过动态评论
//	@param userId string -
//	@param replyId int64 -
//	@return bool -
func checkUserIsPraisedTimelineReply(userId string, replyId int64) bool {
	_, err := coreRedis.GetUserRedis().ZRank(context.Background(), redisKey.TimelineReplyPraisedUserList(replyId), userId).Result()
	if err != nil {
		return false
	}
	return true
}

// genTimelineReplyInfo
//
//	@Description: 包装动态评论数据
//	@param data model.TimelineReply -
//	@return dst -
func genTimelineReplyInfo(data model.TimelineReply, clientType string) (dst *response_user.TimelineReplyInfo) {
	dst = &response_user.TimelineReplyInfo{
		ReplyId:       data.Id,
		TimelineId:    data.TimelineId,
		ReplierId:     data.ReplierId,
		ReplyContent:  data.ReplyContent,
		ToReplyId:     data.ToReplyId,
		ToReplierId:   data.ToReplierId.String,
		SubReplyCount: data.SubReplyCount,
		ToSubReplyId:  data.ToSubReplyId,
		IsPraised:     false,
		PraisedCount:  data.PraisedCount,
		CreateTime:    data.CreateTime.Format(time.DateTime),
		SubReplyList:  nil,
		CreateTimeStr: getTimelineReplyTimeStr(data.CreateTime),
		UserPlaque:    service_user.GetUserLevelPlaque(data.ReplierId, clientType),
	}

	// 评论人
	userDao := &dao.UserDao{}
	replierModel, _ := userDao.FindOne(&model.User{Id: data.ReplierId})
	if len(replierModel.Id) > 0 {
		dst.ReplierName = replierModel.Nickname
		dst.ReplierAvatar = helper.FormatImgUrl(replierModel.Avatar)
	}

	// 被评论人
	if len(data.ToReplierId.String) > 0 {
		toReplierModel, _ := userDao.FindOne(&model.User{Id: data.ToReplierId.String})
		if len(toReplierModel.Id) > 0 {
			dst.ToReplierName = toReplierModel.Nickname
		}
	}
	return
}
