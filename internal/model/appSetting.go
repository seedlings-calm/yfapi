package model

import "time"

type AppSetting struct {
	ID             int       `json:"id" gorm:"column:id"`
	Status         int       `json:"status" gorm:"column:status"`                   // 1:开启 2:关闭
	WithdrawDays   string    `json:"withdrawDays" gorm:"column:withdraw_days"`      // 0:周日 1:周一 2:周二 3:周三 4:周四 5:周五 6:周六
	StartTime      string    `json:"startTime" gorm:"column:start_time"`            // 可提现开始时间
	EndTime        string    `json:"endTime" gorm:"column:end_time"`                // 可提现结束时间
	RewardRate     int       `json:"rewardRate" gorm:"column:reward_rate"`          // 打赏收入提现手续费(百分比)
	SettlementRate int       `json:"settlementRate" gorm:"column:settlement_rate"`  // 结算收入提现手续费(百分比)
	WithdrawDesc   string    `json:"withdrawDesc" gorm:"column:withdraw_desc"`      // 可提现金额说明
	UnWithdrawDesc string    `json:"unWithdrawDesc" gorm:"column:un_withdraw_desc"` // 不可提现金额说明
	Desc           string    `json:"desc" gorm:"column:desc"`                       // 提现说明
	StaffName      string    `json:"staffName" gorm:"column:staff_name"`            // 操作人
	UpdateTime     time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *AppSetting) TableName() string {
	return "t_app_setting"
}
