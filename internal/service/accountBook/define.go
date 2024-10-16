package accountBook

import (
	"fmt"
	"time"
)

var (
	//订单号自增记录
	lastOrderAutoIncrNum = func() string {
		return fmt.Sprintf("order:autoIncrNum:%s", time.Now().Format("20060102"))
	}
)

const (
	PAY_TYPE_ALIPAY      = 1 //支付宝
	PAY_TYPE_WEIXIN      = 2 //微信
	PAY_TYPE_APPLE_IAP   = 3 //苹果内购支付
	PAY_TYPE_GOOGLE_IAP  = 4 //google内购
	PAY_TYPE_AGGREGATION = 5 //第三方聚合支付
)

// 订单类型
const (
	ORDER_CZ = 1 //充值订单
	ORDER_SC = 2 //商城订单
	ORDER_TX = 3 //提现订单
	ORDER_PW = 4 //陪玩订单
	ORDER_DS = 5 //礼物打赏订单
	ORDER_OR = 6 //其他订单
)

const (
	ORDER_STATUS_UNCOMPLETION = 0 //未完成订单
	ORDER_STATUS_COMPLETION   = 1 //已完成订单
)

// 支付状态
const (
	PAY_STATUS_UNPAID      = 0 //待支付
	PAY_STATUS_COMPLETION  = 1 //支付完成
	PAY_STATUS_REFUNDING   = 2 //退款中
	PAY_STATUS_REFUNDED    = 4 //退款完成
	PAY_STATUS_REFUND_FAIL = 5 //退款失败
)

// 资金类型
const (
	CURRENCY_DIAMOND              = "diamond"      //钻石
	CURRENCY_STARLIGHT            = "starlight"    //星光 扣除星光时使用
	CURRENCY_STARLIGHT_WITHDRAW   = "starlightW"   //可提现星光
	CURRENCY_STARLIGHT_UNWITHDRAW = "starlightUW"  //不可提现星光
	CURRENCY_STARLIGHT_SUBSIDY    = "starlightSUB" //补贴星光
	CURRENCY_CNY                  = "CNY"          //人民币
	CURRENCY_USD                  = "USD"          //美元
)

// 资金方向
const (
	FUND_INFLOW  = 1 //资金流入
	FUND_OUTFLOW = 2 //资金流出
)

// 钻石账变类型
const (
	ChangeDiamondRoomDailyFundFlowSettlement           = 101 //房间流水日结
	ChangeDiamondRoomMonthlyFundFlowSettlement         = 102 //房间流水月结
	ChangeDiamondGuildLiveRoomSubsidyMonthlySettlement = 103 //公会直播间补贴月结
	ChangeDiamondGuildFlowSubsidyMonthlySettlement     = 104 //公会流水补贴月结
	ChangeDiamondStarlightExchange                     = 105 //兑换钻石
	ChangeDiamondRecharge                              = 106 //充值
	ChangeDiamondOperationGift                         = 107 //运营赠送
	ChangeDiamondInviteFriends                         = 108 //邀请好友
	ChangeDiamondRedPacketIncome                       = 109 //红包收益
	ChangeDiamondViolationDeduction                    = 111 //违规扣除
	ChangeDiamondRewardGift                            = 112 //打赏礼物
	ChangeDiamondMallConsumption                       = 113 //商城消耗
	ChangeDiamondNicknameModification                  = 114 //昵称修改
	ChangeDiamondAvatarModification                    = 115 //头像修改
	ChangeDiamondRedPacketConsumption                  = 116 //红包消耗
	ChangeDiamondCPAccountChange                       = 117 //CP账变
	ChangeDiamondCloseFriendAccountChange              = 118 //挚友账变
	ChangeDiamondPenalty                               = 119 //违约金
	ChangeDiamondPurchaseMysteryPerson                 = 120 //购买神秘人
	ChangeDiamondRenewMysteryPerson                    = 121 //续费神秘人
	ChangeDiamondAppointmentMysteryPerson              = 122 //预约神秘人
	ChangeDiamondTransfer                              = 123 //转账
)

// 星光账变类型
const (
	ChangeStarlightRewardIncome                          = 201 //打赏收益
	ChangeStarlightRoomFlowDailySettlement               = 202 //房间流水日结
	ChangeStarlightRoomFlowMonthlySettlement             = 203 //房间流水月结
	ChangeStarlightGuildLiveRoomSubsidyMonthlySettlement = 204 //公会直播间补贴月结
	ChangeStarlightGuildFlowSubsidyMonthlySettlement     = 205 //公会流水补贴月结
	ChangeStarlightOperationGift                         = 207 //运营赠送
	ChangeStarlightStarlightWithdrawal                   = 208 //提现
	ChangeStarlightViolationDeduction                    = 209 //违规扣除
	ChangeStarlightStarlightExchange                     = 210 //兑换钻石
)

// 额外订单账变类型
const (
	SubsidyRoomProfitDaily  = 301 // 房间流水日结补贴
	SubsidyRoomProfitMonth  = 302 // 房间流水月结补贴
	SubsidyGuildProfitMonth = 303 // 公会流水月结补贴
	SubsidyGuildAnchorMonth = 304 // 公会直播间月结补贴
)

type EnumOrderType int

func (e EnumOrderType) String() string {
	switch e {
	case ChangeDiamondRoomDailyFundFlowSettlement:
		return "房间流水日结"
	case ChangeDiamondRoomMonthlyFundFlowSettlement:
		return "房间流水月结"
	case ChangeDiamondGuildLiveRoomSubsidyMonthlySettlement:
		return "公会直播间补贴月结"
	case ChangeDiamondGuildFlowSubsidyMonthlySettlement:
		return "公会流水补贴月结"
	case ChangeDiamondStarlightExchange:
		return "兑换钻石"
	case ChangeDiamondRecharge:
		return "充值"
	case ChangeDiamondOperationGift:
		return "运营赠送"
	case ChangeDiamondInviteFriends:
		return "邀请好友"
	case ChangeDiamondRedPacketIncome:
		return "红包收益"
	case ChangeDiamondViolationDeduction:
		return "违规扣除"
	case ChangeDiamondRewardGift:
		return "打赏礼物"
	case ChangeDiamondMallConsumption:
		return "商城消耗"
	case ChangeDiamondNicknameModification:
		return "昵称修改"
	case ChangeDiamondAvatarModification:
		return "头像修改"
	case ChangeDiamondRedPacketConsumption:
		return "红包消耗"
	case ChangeDiamondCPAccountChange:
		return "CP账变"
	case ChangeDiamondCloseFriendAccountChange:
		return "挚友账变"
	case ChangeDiamondPenalty:
		return "违约金"
	case ChangeDiamondPurchaseMysteryPerson:
		return "购买神秘人"
	case ChangeDiamondRenewMysteryPerson:
		return "续费神秘人"
	case ChangeDiamondAppointmentMysteryPerson:
		return "预约神秘人"
	case ChangeDiamondTransfer:
		return "转账"
	case ChangeStarlightRewardIncome:
		return "打赏收益"
	case ChangeStarlightRoomFlowDailySettlement:
		return "房间流水日结"
	case ChangeStarlightRoomFlowMonthlySettlement:
		return "房间流水月结"
	case ChangeStarlightGuildLiveRoomSubsidyMonthlySettlement:
		return "公会直播间补贴月结"
	case ChangeStarlightGuildFlowSubsidyMonthlySettlement:
		return "公会流水补贴月结"
	case ChangeStarlightOperationGift:
		return "运营赠送"
	case ChangeStarlightStarlightWithdrawal:
		return "提现"
	case ChangeStarlightViolationDeduction:
		return "违规扣除"
	case ChangeStarlightStarlightExchange:
		return "兑换钻石"
	}
	return ""
}
