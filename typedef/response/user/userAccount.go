package user

import "yfapi/internal/model"

type UserAccountRes struct {
	UserId              string `json:"userId"`              // 用户id
	Status              int    `json:"status"`              // 用户账号状态
	WithdrawStatus      int    `json:"withdrawStatus"`      // 用户账号提现状态
	DiamondAmount       string `json:"diamondAmount"`       // 钻石余额
	StarlightAmount     string `json:"starlightAmount"`     // 星光总余额
	StarlightUnWithdraw string `json:"starlightUnWithdraw"` // 不可提现星光余额
	StarlightWithdraw   string `json:"starlightWithdraw"`   // 可提现星光余额
	StarlightSubsidy    string `json:"starlightSubsidy"`    // 补贴星光余额
}

type UserAccountDTO struct {
	UserId              string `json:"userId"`              // 用户id
	Status              int    `json:"status"`              // 用户账号状态
	WithdrawStatus      int    `json:"withdrawStatus"`      // 用户账号提现状态
	DiamondAmount       string `json:"diamondAmount"`       // 钻石余额
	StarlightAmount     string `json:"starlightAmount"`     // 星光总余额
	StarlightUnWithdraw string `json:"starlightUnWithdraw"` // 不可提现星光余额
	StarlightWithdraw   string `json:"starlightWithdraw"`   // 可提现星光余额
	StarlightSubsidy    string `json:"starlightSubsidy"`    // 补贴星光余额
	Version             int    `json:"version"`             // 账号版本
	SubsidyList         []*model.UserAccountSubsidy
}

type RechargeDiamondRes struct {
	UserId        string                `json:"userId"`        // 用户id
	UserNickname  string                `json:"nickname"`      //用户昵称
	UserNo        string                `json:"userNo"`        //用户昵称
	UserAvatar    string                `json:"userAvatar"`    //用户头像
	Status        int                   `json:"status"`        // 用户账号状态
	DiamondAmount string                `json:"diamondAmount"` // 钻石余额
	ChannelInfo   model.ConfigChannel   `json:"channelInfo"`   //渠道详情
	ChannelGoods  []NewConfigDiamondRes `json:"channelGoods"`  //渠道充值商品
}

type NewConfigDiamondRes struct {
	Keys     string `json:"keys"       description:"商品ID"` //商品ID
	Nums     int    `json:"nums"       description:"充值金额"` //充值金额
	GotoNums int    `json:"gotoNums"   description:"到账金额"` //到账金额
}
