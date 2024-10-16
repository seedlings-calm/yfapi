package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"time"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_user "yfapi/internal/service/user"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
	response_user "yfapi/typedef/response/user"
)

type UserFilter struct {
}

// 动态规则设置
func (u *UserFilter) TimelineFilter(c *gin.Context, req *request_user.TimelineFilterReq) {
	userId := helper.GetUserId(c)
	if userId == req.UserId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
		})
	}
	data, err := new(dao.UserDao).FindOne(&model.User{Id: req.UserId})
	if err != nil || len(data.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
		})
	}
	serDao := new(dao.UserTimelineFilterDao)
	userTimelineFilterModel := serDao.FindOne(&model.UserTimelineFilter{
		UserID: userId,
		ToID:   req.UserId,
		Types:  req.Types,
	})
	if userTimelineFilterModel.ID > 0 {
		err := serDao.Del(userTimelineFilterModel)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
			})
		}
	} else {
		userTimelineFilterModel.UserID = userId
		userTimelineFilterModel.ToID = req.UserId
		userTimelineFilterModel.Types = req.Types
		userTimelineFilterModel.CreateTime = time.Now()
		err := serDao.Add(userTimelineFilterModel)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
			})
		}
	}
	return
}

// 查询用户动态规则设置
func (ser *UserFilter) GetTimelineFilterList(c *gin.Context, req *request_user.GetTimelineFilterListReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	UserTimelineFilterDao := &dao.UserTimelineFilterDao{}
	rows, count, err := UserTimelineFilterDao.GetList(userId, req.Page, req.Size, req.Types)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	var userInfoList []*response_user.UserFilterList
	// 遍历结果集并收集好友的用户ID及信息
	var userIdList []string
	// 遍历结果集并收集关注的用户ID
	for _, v := range rows {
		userIdList = append(userIdList, v.ToID)
		userInfoList = append(userInfoList, &response_user.UserFilterList{
			UserId: v.ToID,
		})
	}
	userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
	for _, info := range userInfoList {
		info.UserNo = userInfoMap[info.UserId].UserNo
		info.Nickname = userInfoMap[info.UserId].Nickname
		info.Avatar = userInfoMap[info.UserId].Avatar
		info.Uid32 = cast.ToInt32(userInfoMap[info.UserId].UserNo)
		info.Sex = userInfoMap[info.UserId].Sex
		info.UserPlaque = service_user.GetUserLevelPlaque(info.UserId, helper.GetClientType(c))
	}
	// 返回好友列表
	res.Data = userInfoList
	res.Total = count
	res.Size = req.Size
	res.CurrentPage = req.Page
	res.CalcHasNext()
	return res
}

// 通知规则开关
func (u *UserFilter) NoticeFilter(c *gin.Context, req *request_user.NoticeFilterReq) {
	userId := helper.GetUserId(c)
	if userId == req.UserId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
		})
	}
	data, err := new(dao.UserDao).FindOne(&model.User{Id: req.UserId})
	if err != nil || len(data.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
		})
	}
	serDao := new(dao.UserNoticeFilterDao)
	userNoticeFilterModel := serDao.FindOne(&model.UserNoticeFilter{
		UserID: userId,
		ToID:   req.UserId,
		Types:  req.Types,
	})
	if userNoticeFilterModel.ID > 0 {
		err := serDao.Del(userNoticeFilterModel)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
			})
		}
	} else {
		userNoticeFilterModel.UserID = userId
		userNoticeFilterModel.ToID = req.UserId
		userNoticeFilterModel.Types = req.Types
		userNoticeFilterModel.CreateTime = time.Now()
		err := serDao.Add(userNoticeFilterModel)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
			})
		}
	}
	return
}
