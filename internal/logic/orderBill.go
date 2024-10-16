package logic

import (
	"fmt"
	"strings"
	"time"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/service/accountBook"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/request/orderBill"
	"yfapi/typedef/response"
	response_orderBill "yfapi/typedef/response/orderBill"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
)

type OrderBill struct {
}

// GetUserDiamondBillList 查询玩家钻石流水
func (o *OrderBill) GetUserDiamondBillList(c *gin.Context, req *orderBill.DiamondBillReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	data, count := new(dao.OrderBillDao).GetUserDiamondBillList(userId, req.Page, req.Size, req.FundFlow, req.TimeKey)
	var result []response_orderBill.DiamondBill
	userIdMap := make(map[string]struct{})
	for _, info := range data {
		dst := response_orderBill.DiamondBill{
			FundFlow:   info.FundFlow,
			Amount:     easy.StringToDecimalFixed(info.Amount),
			RoomName:   info.RoomName,
			CreateTime: info.CreateTime.Format(time.DateTime),
			TimeKey:    info.CreateTime.Format("2006-01"),
		}
		if info.OrderType == accountBook.ChangeDiamondRewardGift { // 礼物打赏
			dst.Title = fmt.Sprintf("%v-%vx%v", accountBook.EnumOrderType(info.OrderType).String(), info.Note, info.Num)
		} else {
			dst.Title = accountBook.EnumOrderType(info.OrderType).String()
		}
		if len(info.ToUserIdList) > 0 {
			toUserList := strings.Split(info.ToUserIdList, ",")
			for _, toUserId := range toUserList {
				dst.ToUserList = append(dst.ToUserList, &response_orderBill.UserInfo{UserId: toUserId})
				userIdMap[toUserId] = struct{}{}
			}
		}
		result = append(result, dst)
	}
	var userIdList []string
	for toUserId, _ := range userIdMap {
		userIdList = append(userIdList, toUserId)
	}
	if len(userIdList) > 0 {
		userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
		for _, info := range result {
			for _, uInfo := range info.ToUserList {
				if uInfo != nil {
					uInfo.Nickname = userInfoMap[uInfo.UserId].Nickname
					uInfo.Avatar = userInfoMap[uInfo.UserId].Avatar
				}
			}
		}
	}
	res.Data = result
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Total = count
	res.CalcHasNext()
	return
}

// GetUserStarlightBillList 查询玩家星光流水
func (o *OrderBill) GetUserStarlightBillList(c *gin.Context, req *orderBill.DiamondBillReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	data, count := new(dao.OrderBillDao).GetUserStarlightBillList(userId, req.Page, req.Size, req.FundFlow, req.TimeKey)
	var result []response_orderBill.StarlightBill
	userIdMap := make(map[string]struct{})
	for _, info := range data {
		dst := response_orderBill.StarlightBill{
			FundFlow:   info.FundFlow,
			Amount:     easy.StringToDecimalFixed(info.Amount),
			RoomName:   info.RoomName,
			CreateTime: info.CreateTime.Format(time.DateTime),
			TimeKey:    info.CreateTime.Format("2006-01"),
		}
		if info.OrderType == accountBook.ChangeStarlightRewardIncome { // 打赏收益
			dst.Title = fmt.Sprintf("%v-%vx%v", accountBook.EnumOrderType(info.OrderType).String(), info.Note, info.Num)
		} else {
			dst.Title = accountBook.EnumOrderType(info.OrderType).String()
		}
		if len(info.FromUserId) > 0 && info.FromUserId != "0" {
			dst.FormUser = &response_orderBill.UserInfo{UserId: info.FromUserId}
			userIdMap[info.FromUserId] = struct{}{}
		}
		result = append(result, dst)
	}
	var userIdList []string
	for toUserId, _ := range userIdMap {
		userIdList = append(userIdList, toUserId)
	}
	if len(userIdList) > 0 {
		userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
		for _, info := range result {
			if info.FormUser != nil {
				info.FormUser.Nickname = userInfoMap[info.FormUser.UserId].Nickname
				info.FormUser.Avatar = userInfoMap[info.FormUser.UserId].Avatar
			}
		}
	}
	res.Data = result
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Total = count
	res.CalcHasNext()
	return
}

func (o *OrderBill) GetUserRechargeDiamonLogList(c *gin.Context, req *orderBill.RechargeDiamondReq) (res response.BasePageRes) {
	userId := helper.GetUserId(c)
	data, count := new(dao.OrderBillDao).GetUserDiamondBillList(userId, req.Page, req.Size, 1, req.TimeKey, accountBook.ChangeDiamondRecharge)
	var result []response_orderBill.DiamondBill
	userIdMap := make(map[string]struct{})
	for _, info := range data {
		dst := response_orderBill.DiamondBill{
			FundFlow:   info.FundFlow,
			Amount:     easy.StringToDecimalFixed(info.Amount),
			RoomName:   info.RoomName,
			CreateTime: info.CreateTime.Format(time.DateTime),
			TimeKey:    info.CreateTime.Format("2006-01"),
		}
		if info.OrderType == accountBook.ChangeDiamondRewardGift { // 礼物打赏
			dst.Title = fmt.Sprintf("%v-%vx%v", accountBook.EnumOrderType(info.OrderType).String(), info.Note, info.Num)
		} else {
			dst.Title = accountBook.EnumOrderType(info.OrderType).String()
		}
		if len(info.ToUserIdList) > 0 {
			toUserList := strings.Split(info.ToUserIdList, ",")
			for _, toUserId := range toUserList {
				dst.ToUserList = append(dst.ToUserList, &response_orderBill.UserInfo{UserId: toUserId})
				userIdMap[toUserId] = struct{}{}
			}
		}
		result = append(result, dst)
	}
	var userIdList []string
	for toUserId, _ := range userIdMap {
		userIdList = append(userIdList, toUserId)
	}
	if len(userIdList) > 0 {
		userInfoMap := service_user.GetUserBaseInfoMap(userIdList)
		for _, info := range result {
			for _, uInfo := range info.ToUserList {
				if uInfo != nil {
					uInfo.Nickname = userInfoMap[uInfo.UserId].Nickname
					uInfo.Avatar = userInfoMap[uInfo.UserId].Avatar
				}
			}
		}
	}
	res.Data = result
	res.CurrentPage = req.Page
	res.Size = req.Size
	res.Total = count
	res.CalcHasNext()
	return
}
