package model

type UserAccount struct {
	UserId            string `json:"userId" gorm:"column:user_id"`                        // 用户ID
	DiamondAmount     string `json:"diamondAmount" gorm:"column:diamond_amount"`          // 钻石余额
	StarlightAmount   string `json:"starlightAmount" gorm:"column:starlight_amount"`      // 星光余额
	CanWithdrawAmount string `json:"canWithdrawAmount" gorm:"column:can_withdraw_amount"` // 可提现星光余额
	Status            int    `json:"status" gorm:"column:status"`                         // 账户状态 1=正常 2=冻结 只出不进
	WithdrawStatus    int    `json:"withdrawStatus" gorm:"column:withdraw_status"`        // 提现状态 1允许,2不允许
	Version           int    `json:"version" gorm:"column:version"`                       // 账号版本
}

func (m *UserAccount) TableName() string {
	return "t_user_account"
}

type UserAccountSubsidy struct {
	Id            int    `json:"id" gorm:"column:id"`
	UserId        string `json:"userId" gorm:"column:user_id"`               // 用户ID
	AccountType   int    `json:"accountType" gorm:"column:account_type"`     // 账户类型 1房间 2公会
	RoomId        string `json:"roomId" gorm:"column:room_id"`               // 房间ID
	GuildId       string `json:"guildId" gorm:"column:guild_id"`             // 公会ID
	Status        int    `json:"status" gorm:"column:status"`                // 账户状态 1=正常 2=冻结
	SubsidyAmount string `json:"subsidyAmount" gorm:"column:subsidy_amount"` // 补贴星光余额
}

func (m *UserAccountSubsidy) TableName() string {
	return "t_user_account_subsidy"
}
