package room

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_user "yfapi/internal/service/user"
	response_roomowner "yfapi/typedef/response/roomOwner"
	"yfapi/util/easy"
)

type RoomAccountInfo struct{}

func (h *RoomAccountInfo) RoomAccountInfoLogic(c *gin.Context) (res response_roomowner.RoomAccountInfoRes) {
	userId := helper.GetUserId(c)
	roomId := helper.GetRoomId(c)
	// 查询玩家房间账户信息
	accountDao := new(dao.UserAccountDao)
	subsidyInfo, err := accountDao.GetRoomAccountGuildSubsidy(userId, roomId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据不存在 初始化
		subsidyInfo = model.UserAccountSubsidy{
			Id:            0,
			UserId:        userId,
			AccountType:   1,
			RoomId:        roomId,
			GuildId:       "0",
			Status:        1,
			SubsidyAmount: "0",
		}
		err = coreDb.GetMasterDb().Create(&subsidyInfo).Error
		if err != nil {
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
	} else if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 查询补贴账户总收益
	totalIncome, _ := accountDao.GetUserRoomHistorySubsidyAmount(userId, roomId)
	res = response_roomowner.RoomAccountInfoRes{
		UserId:          userId,
		Status:          subsidyInfo.Status,
		CashAmount:      easy.StringFixed(easy.StringToDecimal(subsidyInfo.SubsidyAmount).Div(decimal.NewFromInt(10))),
		TotalCashIncome: easy.StringFixed(easy.StringToDecimal(totalIncome).Div(decimal.NewFromInt(10))),
	}
	//查询手续费和提现说明
	userWithdrawDao := dao.UserWithdrawDao{}
	withdrawInfo, err := userWithdrawDao.UserWithdrawInfoById(1)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if withdrawInfo.ID == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeWithdrawConfigNotExist,
			Msg:  nil,
		})
	}
	res.Desc = withdrawInfo.Desc
	res.SettlementRate = withdrawInfo.SettlementRate
	//查询用户信息
	userInfo := service_user.GetUserBaseInfo(userId)
	res.Mobile = userInfo.Mobile
	res.TrueName = userInfo.TrueName
	userBankDao := new(dao.UserBankDao)
	bankInfo, err := userBankDao.GetUserBankListByUserId(userId)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	res.BankList = make([]response_roomowner.BankInfo, 0)
	if len(bankInfo) > 0 {
		for _, v := range bankInfo {
			bank := response_roomowner.BankInfo{
				Id:         v.Id,
				BankNo:     v.BankNo,
				BankName:   v.BankName,
				BankHolder: v.BankHolder,
				BankBranch: v.BankBranch,
				IsDefault:  v.IsDefault,
			}
			res.BankList = append(res.BankList, bank)
		}
	}
	return res
}
