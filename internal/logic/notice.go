package logic

import (
	"github.com/gin-gonic/gin"
	"time"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_im "yfapi/internal/service/im"
	service_user "yfapi/internal/service/user"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/message"
	"yfapi/util/easy"
)

type Notice struct {
}

// 动态发布通知
func (n *Notice) MomentsPublishNotice(c *gin.Context, userId string, timelineId int64) {
	//获取用户粉丝
	ids := new(dao.UserFollowDao).GetUserFollowersIds(userId)
	if len(ids) == 0 {
		return
	}
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	nowTime := time.Now().Format(time.DateTime)
	timelineModel, _ := new(dao.TimelineDao).GetTimelineById(timelineId)
	timelineInfo := genTimelineInfo(c, timelineModel, userId, helper.GetClientType(c))
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
	userIds := []string{}
	getIds := new(dao.UserNoticeFilterDao).GetIds(userId, typedef_enum.MomentsNoticeSwitch)
	for _, v := range ids {
		if !easy.InArray(v, getIds) {
			userIds = append(userIds, v)
		}
	}
	if len(userIds) == 0 {
		return
	}
	go new(service_im.ImNoticeService).SendInteractiveMsg(c, message.SendInteractiveMsg{
		Avatar:          helper.FormatImgUrl(userModel.Avatar),
		NickName:        userModel.Nickname,
		CreateTime:      nowTime,
		ImageUrl:        imageUrl,
		VideoUrl:        videoUrl,
		Msg:             userModel.Nickname + i18n_msg.GetI18nMsg(c, i18n_msg.NewPostPublishedKey),
		MsgType:         typedef_enum.UserPublishNewMoments,
		TimelineId:      timelineInfo.TimelineId,
		TimelineContent: timelineInfo.TextContent,
		UserPlaque:      service_user.GetUserLevelPlaque(userModel.Id, helper.GetClientType(c)),
	}, userIds)
}

// 直播发布通知
func (n *Notice) LivePublishNotice(c *gin.Context, userId string, roomId string) {

	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	//获取用户粉丝
	ids := new(dao.UserFollowDao).GetUserFollowersIds(userId)
	if len(ids) == 0 {
		return
	}
	nowTime := time.Now().Format(time.DateTime)
	roomInfo, _ := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	userIds := []string{}
	getIds := new(dao.UserNoticeFilterDao).GetIds(userId, typedef_enum.LiveNoticeSwitch)
	for _, v := range ids {
		if !easy.InArray(v, getIds) {
			userIds = append(userIds, v)
		}
	}
	if len(userIds) == 0 {
		return
	}
	go new(service_im.ImNoticeService).SendInteractiveMsg(c, message.SendInteractiveMsg{
		Avatar:     helper.FormatImgUrl(userModel.Avatar),
		NickName:   userModel.Nickname,
		CreateTime: nowTime,
		ImageUrl:   helper.FormatImgUrl(roomInfo.CoverImg),
		Msg:        userModel.Nickname + i18n_msg.GetI18nMsg(c, i18n_msg.LiveStreamCurrentlyKey),
		MsgType:    typedef_enum.StartLivingStream,
		RoomId:     roomId,
		UserPlaque: service_user.GetUserLevelPlaque(userModel.Id, helper.GetClientType(c)),
		UserId:     userId,
	}, userIds)
}
