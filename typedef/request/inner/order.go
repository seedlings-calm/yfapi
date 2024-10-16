package request_inner

type AccountChangeReq struct {
	UserNo    string  `json:"UserNo"           description:"用户ID" validate:"required"`
	FundFlow  int     `json:"fundFlow"     description:"变动类型：1入账 2出账" validate:"min=1,max=2"`
	Note      string  `json:"note" dc:"变动原因" validate:"required"`
	Money     float64 `json:"money"            description:"变动数量" validate:"required"`
	AdminName string  `json:"adminName" dc:"后台操作人" validate:"required"`
}

// OperationChangeAccountReq
// @Description: 后台变动用户资产请求
type OperationChangeAccountReq struct {
	UserId            string // 用户ID
	Currency          string // 币种
	FundFlow          int    // 资金方向 1入 2出
	Amount            string // 变动数量
	OrderId           string // 关联订单ID
	OrderType         int    // 订单类型
	RoomId            string // 关联房间ID
	GuildId           string // 关联公会ID
	Note              string // 备注信息
	SubsidyType       int    // 补贴类型 1房间日补贴 2房间月补贴 4公会月补贴
	SubsidyAmountType int    // 补贴账户类型 1房间 2公会
}
