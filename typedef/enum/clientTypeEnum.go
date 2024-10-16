package enum

const (
	ClientTypePc      = "pc"
	ClientTypeIos     = "ios"
	ClientTypeAndroid = "android"
	ClientTypeH5      = "h5"
)

var (
	ClientTypeArray = []string{ClientTypeAndroid, ClientTypePc, ClientTypeIos, ClientTypeH5}
)

// 菜单模块类型
const (
	ModuleTypeUserCenter = 1 // 个人主页模块
)

// 菜单类型
const (
	MenuTypeWallet        = 1 // 我的钱包
	MenuTypeOrder         = 2 // 我的订单
	MenuTypeSetting       = 3 // 设置与隐私
	MenuTypeFans          = 4 // 我的粉丝团
	MenuTypeGuild         = 5 // 我的公会
	MenuTypePractitioners = 6 // 从业者考核
	MenuTypeAnchorRoom    = 7 // 我的直播间
	MenuTypeCostume       = 8 // 装扮商城
	MenuTypeLevel         = 9 // 会员体系
)
