package logic

import (
	"github.com/gin-gonic/gin"
	"time"
	error2 "yfapi/i18n/error"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/request"
	"yfapi/typedef/response"
	response_user "yfapi/typedef/response/user"
)

// UserVisit
// @Description: 用户访问足迹
type UserVisit struct {
}

// GetUserVisitRecordList
//
//	@Description: 查询用户足迹记录列表
//	@receiver u
//	@param c *gin.Context -
//	@param req *request.BasePageReq -
//	@return res -
func (u *UserVisit) GetUserVisitRecordList(c *gin.Context, req *request.BasePageReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	// TODO 只能查两百条记录
	var result []*response_user.UserVisitInfo
	dataList, count, _ := new(dao.UserVisitDao).GetUserVisitRecordList(userId, req.Page, req.Size)
	for _, info := range dataList {
		dst := &response_user.UserVisitInfo{
			UserId:    info.TargetUserId,
			Nickname:  info.Nickname,
			Avatar:    helper.FormatImgUrl(info.Avatar),
			Sex:       info.Sex,
			Introduce: info.Introduce,
		}
		// 时间描述
		dst.TimeDesc = genVisitTimeDesc(c, info.UpdateTime, true)
		// 额外描述
		dst.Extra = genVisitExtra(c, info.IsVisit, false)
		// 查询用户的铭牌信息
		dst.UserPlaque = service_user.GetUserLevelPlaque(info.TargetUserId, helper.GetClientType(c))
		result = append(result, dst)
	}
	// 处理返回结果
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Total = count
	res.Data = result
	res.CalcHasNext()
	return
}

// GetVisitUserRecordList
//
//	@Description: 查询用户访客记录列表
//	@receiver u
//	@param c *gin.Context -
//	@param req *request.BasePageReq -
//	@return res -
func (u *UserVisit) GetVisitUserRecordList(c *gin.Context, req *request.BasePageReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	// TODO 只能查两百条记录
	var result []*response_user.UserVisitInfo
	dataList, count, _ := new(dao.UserVisitDao).GetVisitUserRecordList(userId, req.Page, req.Size)
	for _, info := range dataList {
		dst := &response_user.UserVisitInfo{
			UserId:   info.UserId,
			Nickname: info.Nickname,
			Avatar:   helper.FormatImgUrl(info.Avatar),
			Sex:      info.Sex,
		}
		// 时间描述
		dst.TimeDesc = genVisitTimeDesc(c, info.UpdateTime, false)
		// 额外描述
		dst.Extra = genVisitExtra(c, info.IsVisit, true)
		// 查询用户的铭牌信息
		dst.UserPlaque = service_user.GetUserLevelPlaque(info.UserId, helper.GetClientType(c))
		result = append(result, dst)
	}
	// 处理返回结果
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Total = count
	res.Data = result
	res.CalcHasNext()
	return
}

// ClearUserVisitRecord
//
//	@Description: 清除用户足迹
//	@receiver u
//	@param c *gin.Context -
func (u *UserVisit) ClearUserVisitRecord(c *gin.Context) {
	userId := helper.GetUserId(c)
	err := new(dao.UserVisitDao).ClearUserVisitRecord(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	return
}

func genVisitTimeDesc(c *gin.Context, createTime time.Time, isSelf bool) (dst string) {
	year := createTime.Year()
	now := time.Now()
	currYear := now.Year()
	if year != currYear {
		msyKey := i18n_msg.VisitYearHeMsyKey
		if isSelf {
			msyKey = i18n_msg.VisitYearMeMsyKey
		}
		dst = i18n_msg.GetI18nMsg(c, msyKey, map[string]any{"year": createTime.Year(), "month": int(createTime.Month())})
		return
	}

	minutes := time.Since(createTime).Minutes()
	switch {
	case minutes > 7*24*60: // 大于七天
		msyKey := i18n_msg.VisitMonthHeMsyKey
		if isSelf {
			msyKey = i18n_msg.VisitMonthMeMsyKey
		}
		dst = i18n_msg.GetI18nMsg(c, msyKey, map[string]any{"month": int(createTime.Month()), "day": createTime.Day()})
	case minutes > 24*60: // 大于1天 小于七天
		msgKey := i18n_msg.VisitDayHeMsgKey
		if isSelf {
			msgKey = i18n_msg.VisitDayMeMsgKey
		}
		dst = i18n_msg.GetI18nMsg(c, msgKey, map[string]any{"num": int(minutes / (24 * 60))})
	case minutes > 60:
		msgKey := i18n_msg.VisitHourHeMsgKey
		if isSelf {
			msgKey = i18n_msg.VisitHourMeMsgKey
		}
		dst = i18n_msg.GetI18nMsg(c, msgKey, map[string]any{"num": int(minutes / 60)})
	default:
		if minutes < 1 {
			minutes = 1
		}
		msgKey := i18n_msg.VisitMinuteHeMsgKey
		if isSelf {
			msgKey = i18n_msg.VisitMinuteMeMsgKey
		}
		dst = i18n_msg.GetI18nMsg(c, msgKey, map[string]any{"num": int(minutes)})
	}
	return
}

func genVisitExtra(c *gin.Context, isVisit, isSelf bool) (dst string) {
	if isVisit {
		msgKey := i18n_msg.VisitExtraHeMsgKey
		if isSelf {
			msgKey = i18n_msg.VisitExtraMeMsgKey
		}
		dst = i18n_msg.GetI18nMsg(c, msgKey)
	}
	return
}
