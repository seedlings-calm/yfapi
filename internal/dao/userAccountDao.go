package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
)

type UserAccountDao struct {
}

// GetUserAccountByUserId 查询用户账户信息
func (u *UserAccountDao) GetUserAccountByUserId(userId string) (res model.UserAccount, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserAccount{}).Where("user_id=?", userId).First(&res).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = coreDb.GetMasterDb().Exec("insert into t_user_account(user_id) values(?)", userId).Error
		if err == nil {
			return u.GetUserAccountByUserId(userId)
		}
	}
	return
}

// UpdateUserAccount 更新用户账户信息
func (u *UserAccountDao) UpdateUserAccount(dst *model.UserAccount) bool {
	dst.Version++
	return coreDb.GetMasterDb().Model(model.UserAccount{}).Where("version<?", dst.Version).Save(dst).RowsAffected > 0
}

// GetUserAccountSubsidyTotalCountByUserId 查询用户补贴账号余额
func (u *UserAccountDao) GetUserAccountSubsidyTotalCountByUserId(userId string) (amount string, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserAccountSubsidy{}).Where("user_id=?", userId).Select("IFNULL(sum(subsidy_amount), 0)").Scan(&amount).Error
	return
}

// GetUserAccountSubsidyListByUserId 查询用户补贴账号列表信息
func (u *UserAccountDao) GetUserAccountSubsidyListByUserId(userId string) (result []*model.UserAccountSubsidy, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserAccountSubsidy{}).Where("user_id=? and subsidy_amount>0", userId).Scan(&result).Error
	return
}

// GetUserAccountRoomSubsidy 查询用户房间补贴账号信息
func (u *UserAccountDao) GetUserAccountRoomSubsidy(userId, roomId string) (result model.UserAccountSubsidy, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserAccountSubsidy{}).Where("user_id=? and room_id=?", userId, roomId).Scan(&result).Error
	return
}

// GetUserAccountGuildSubsidy 查询用户公会补贴账号信息
func (u *UserAccountDao) GetUserAccountGuildSubsidy(userId, guildId string) (result model.UserAccountSubsidy, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserAccountSubsidy{}).Where("user_id=? and guild_id=?", userId, guildId).Scan(&result).Error
	return
}

// UpdateUserAccountSubsidy 更新用户补贴账号信息
func (u *UserAccountDao) UpdateUserAccountSubsidy(dst *model.UserAccountSubsidy) bool {
	return coreDb.GetMasterDb().Model(model.UserAccountSubsidy{}).Save(dst).RowsAffected > 0
}

// GetUserGuildHistorySubsidyAmount 查询用户历史公会账号补贴
func (u *UserAccountDao) GetUserGuildHistorySubsidyAmount(userId string) (amount string, err error) {
	orderTypeList := []int{accountBook.ChangeStarlightRoomFlowDailySettlement, accountBook.ChangeStarlightRoomFlowMonthlySettlement, accountBook.ChangeStarlightGuildLiveRoomSubsidyMonthlySettlement, accountBook.ChangeStarlightGuildFlowSubsidyMonthlySettlement}
	err = coreDb.GetSlaveDb().Model(model.OrderBill{}).Where("user_id=? and currency=? and fund_flow=?", userId, accountBook.CURRENCY_STARLIGHT_SUBSIDY, accountBook.FUND_INFLOW).
		Where("order_type in ?", orderTypeList).Pluck("sum(amount)", &amount).Error
	return
}

// GetUserAccountGuildSubsidy 查询用户房间补贴账号信息
func (u *UserAccountDao) GetRoomAccountGuildSubsidy(userId, roomId string) (result model.UserAccountSubsidy, err error) {
	err = coreDb.GetSlaveDb().Model(model.UserAccountSubsidy{}).Where("user_id=? and room_id=?", userId, roomId).Scan(&result).Error
	return
}

// GetUserGuildHistorySubsidyAmount 查询用户历史房间账号补贴
func (u *UserAccountDao) GetUserRoomHistorySubsidyAmount(userId, roomId string) (amount string, err error) {
	orderTypeList := []int{accountBook.ChangeStarlightRoomFlowDailySettlement, accountBook.ChangeStarlightRoomFlowMonthlySettlement}
	err = coreDb.GetSlaveDb().Model(model.OrderBill{}).Where("user_id=? and room_id=? and currency=? and fund_flow=?", userId, roomId, accountBook.CURRENCY_STARLIGHT_SUBSIDY, accountBook.FUND_INFLOW).
		Where("order_type in ?", orderTypeList).Pluck("sum(amount)", &amount).Error
	return
}
