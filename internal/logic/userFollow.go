package logic

import (
	"github.com/gin-gonic/gin"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreDb"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_user "yfapi/internal/service/user"
	request_user "yfapi/typedef/request/user"
	"yfapi/typedef/response"
	response_user "yfapi/typedef/response/user"
)

type Follow struct {
}

// 添加关注
func (ser *Follow) AddFollowUser(req *request_user.AddFollowReq, c *gin.Context) (res response_user.AddFollowRes) {
	userId := handle.GetUserId(c)
	followUserId := req.FocusUserId
	if userId == followUserId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	//判断关注的用户是否存在
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: followUserId})
	if err != nil || userModel == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	UserFollowDao := &dao.UserFollowDao{}
	//判断用户是否关注了该用户
	followModel, err := UserFollowDao.GetUserFollow(userId, followUserId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if followModel.Id > 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserAlreadyFollow,
			Msg:  nil,
		})
	}
	//判断用户是否被该用户关注
	rsp, err := UserFollowDao.GetUserFollow(followUserId, userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	isOk := rsp.Id > 0
	//添加关注
	tx := coreDb.GetMasterDb().Begin()
	data := &model.UserFollow{
		UserID:         userId,
		FocusUserID:    followUserId,
		IsMutualFollow: isOk,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
	}
	err = tx.Create(data).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	res.FollowedType = 1 // 已关注
	if isOk {
		//修改互相关注状态
		err = tx.Model(&model.UserFollow{}).Where("id = ?", rsp.Id).Update("is_mutual_follow", isOk).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		res.FollowedType = 2 // 互相关注
	}
	tx.Commit()

	followedNum := new(dao.UserFollowDao).GetUserFollowedNum(followUserId)
	res.FollowedNum = int(followedNum)
	fansNum := new(dao.UserFollowDao).GetUserFansNum(followUserId)
	res.FansNum = int(fansNum)
	return
}

// 取消关注
func (ser *Follow) RemoveFollowUser(req *request_user.RemoveFollowReq, c *gin.Context) (res response_user.AddFollowRes) {
	userId := handle.GetUserId(c)
	followUserId := req.FocusUserId
	//判断取消关注的用户是否存在
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil || userModel == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	//判断用户是否关注了该用户
	UserFollowDao := &dao.UserFollowDao{}
	rsp, err := UserFollowDao.GetUserFollow(userId, followUserId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if rsp.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFollow,
			Msg:  nil,
		})
	}
	//查询用户是否被该用户关注
	rsps, err := UserFollowDao.GetUserFollow(followUserId, userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	isOk := rsps.Id > 0
	//取消关注
	data := &model.UserFollow{
		Id: rsp.Id,
	}
	tx := coreDb.GetMasterDb().Begin()
	err = tx.Model(&model.UserFollow{}).Delete(data).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	//如果被关注修改互相关注状态
	if isOk {
		err = tx.Model(&model.UserFollow{}).Where("id = ?", rsps.Id).Update("is_mutual_follow", !isOk).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	}
	tx.Commit()

	followedNum := new(dao.UserFollowDao).GetUserFollowedNum(followUserId)
	res.FollowedNum = int(followedNum)
	fansNum := new(dao.UserFollowDao).GetUserFansNum(followUserId)
	res.FansNum = int(fansNum)
	return
}

// 查询用户关注列表的方法
func (ser *Follow) GetUserFollowingUserList(req *request_user.GetUserFollowingListReq, c *gin.Context) (res response.BasePageRes) {
	var userId string
	if len(req.UserId) == 0 {
		userId = handle.GetUserId(c)
	} else {
		userId = req.UserId
	}
	UserFollowDao := &dao.UserFollowDao{}
	data, count, err := UserFollowDao.GetUserFollowingUserLists(userId, req.Page, req.Size)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	var userInfoList []*response_user.GetUserFollowingList
	var userIdList []string
	// 遍历结果集并收集关注的用户ID
	for _, v := range data {
		userIdList = append(userIdList, v.FocusUserID)
		userInfoList = append(userInfoList, &response_user.GetUserFollowingList{
			UserId:         v.FocusUserID,
			IsMutualFollow: v.IsMutualFollow,
			FollowedType:   1,
		})
	}
	userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
	for _, info := range userInfoList {
		info.UserNo = userInfoMap[info.UserId].UserNo
		info.Nickname = userInfoMap[info.UserId].Nickname
		info.Avatar = userInfoMap[info.UserId].Avatar
		info.Sex = userInfoMap[info.UserId].Sex
		info.Introduce = userInfoMap[info.UserId].Introduce
		if info.IsMutualFollow {
			info.FollowedType = 2
		}
		info.UserPlaque = service_user.GetUserLevelPlaque(info.UserId, helper.GetClientType(c))
	}
	// 返回关注列表
	res.Data = userInfoList
	res.Total = count
	res.Size = req.Size
	res.CurrentPage = req.Page
	res.CalcHasNext()
	return res
}

// 查询用户粉丝列表的方法
func (ser *Follow) GetUserFollowersList(req *request_user.GetUserFollowingListReq, c *gin.Context) (res response.BasePageRes) {
	var userId string
	if len(req.UserId) == 0 {
		userId = handle.GetUserId(c)
	} else {
		userId = req.UserId
	}
	UserFollowDao := &dao.UserFollowDao{}
	rows, count, err := UserFollowDao.GetUserFollowersList(userId, req.Page, req.Size)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	var userInfoList []*response_user.GetUserFollowingList
	var userIdList []string
	// 遍历结果集并收集粉丝的用户ID及信息
	for _, v := range rows {
		userIdList = append(userIdList, v.UserID)
		userInfoList = append(userInfoList, &response_user.GetUserFollowingList{
			UserId:         v.UserID,
			IsMutualFollow: v.IsMutualFollow,
			FollowedType:   3,
		})
	}
	userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
	for _, info := range userInfoList {
		info.UserNo = userInfoMap[info.UserId].UserNo
		info.Nickname = userInfoMap[info.UserId].Nickname
		info.Avatar = userInfoMap[info.UserId].Avatar
		info.Sex = userInfoMap[info.UserId].Sex
		info.Introduce = userInfoMap[info.UserId].Introduce
		if info.IsMutualFollow {
			info.FollowedType = 2
		}
		info.UserPlaque = service_user.GetUserLevelPlaque(info.UserId, helper.GetClientType(c))
	}
	// 返回粉丝列表
	res.Data = userInfoList
	res.Total = count
	res.Size = req.Size
	res.CurrentPage = req.Page
	res.CalcHasNext()
	return res
}

// 查询用户好友列表的方法
func (ser *Follow) GetUserFriendList(req *request_user.GetUserFriendsListReq, c *gin.Context) (res response.BasePageRes) {
	var userId string
	if len(req.UserId) == 0 {
		userId = handle.GetUserId(c)
	} else {
		userId = req.UserId
	}
	UserFollowDao := &dao.UserFollowDao{}
	rows, count, err := UserFollowDao.GetUserFriendsList(userId, req.Page, req.Size)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	var userInfoList []*response_user.GetUserFollowingList
	// 遍历结果集并收集好友的用户ID及信息
	var userIdList []string
	// 遍历结果集并收集关注的用户ID
	for _, v := range rows {
		userIdList = append(userIdList, v.FocusUserID)
		userInfoList = append(userInfoList, &response_user.GetUserFollowingList{
			UserId:         v.FocusUserID,
			IsMutualFollow: v.IsMutualFollow,
			FollowedType:   2,
		})
	}
	userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
	for _, info := range userInfoList {
		info.UserNo = userInfoMap[info.UserId].UserNo
		info.Nickname = userInfoMap[info.UserId].Nickname
		info.Avatar = userInfoMap[info.UserId].Avatar
		info.Sex = userInfoMap[info.UserId].Sex
		info.Introduce = userInfoMap[info.UserId].Introduce
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

// 删除粉丝
func (ser *Follow) DeleteFansUser(req *request_user.DeleteFansReq, c *gin.Context) (res response_user.AddFollowRes) {
	userId := handle.GetUserId(c)
	followUserId := req.FocusUserId
	//判断取消关注的用户是否存在
	userModel, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil || userModel == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userModel.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	UserFollowDao := &dao.UserFollowDao{}
	//查询用户是否关注自己
	otherRsp, err := UserFollowDao.GetUserFollow(followUserId, userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if otherRsp.Id == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFollow,
			Msg:  nil,
		})
	}
	//判断自己是否关注了该用户
	selfRsp, err := UserFollowDao.GetUserFollow(userId, followUserId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	isOk := selfRsp.Id > 0
	//取消关注
	data := &model.UserFollow{
		Id: otherRsp.Id,
	}
	tx := coreDb.GetMasterDb().Begin()
	err = tx.Model(&model.UserFollow{}).Delete(data).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	//如果被关注修改互相关注状态
	if isOk {
		err = tx.Model(&model.UserFollow{}).Where("id = ?", selfRsp.Id).Update("is_mutual_follow", !isOk).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	}
	tx.Commit()
	//对方的关注数量
	followedNum := new(dao.UserFollowDao).GetUserFollowedNum(followUserId)
	res.FollowedNum = int(followedNum)
	if isOk {
		res.FollowedType = 1
	}
	return
}
