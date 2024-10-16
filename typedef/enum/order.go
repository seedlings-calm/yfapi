package enum

//提现订单状态 0待审核 1审核拒绝 2审核通过 3打款中 4打款成功 5打款失败 6退还成功 7退还失败

// WithdrawStatusWaitReview = 0 //待审核
// WithdrawStatusReject     = 1 //审核拒绝
// WithdrawStatusPass       = 2 //审核通过
// WithdrawStatusPaying     = 3 //打款中
// WithdrawStatusSuccess    = 4 //打款成功
// WithdrawStatusFail       = 5 //打款失败
// WithdrawStatusReturn     = 6 //退还成功
// WithdrawStatusReturnFail = 7 //退还失败
const (
	//WithdrawStatusWaitReview = 0 //待审核
	WithdrawStatusWaitReview = 0
	//WithdrawStatusReject     = 1 //审核拒绝
	WithdrawStatusReject = 1
	//WithdrawStatusPass       = 2 //审核通过
	WithdrawStatusPass = 2
	//WithdrawStatusPaying     = 3 //打款中
	WithdrawStatusPaying = 3
	//WithdrawStatusSuccess    = 4 //打款成功
	WithdrawStatusSuccess = 4
	//WithdrawStatusFail       = 5 //打款失败
	WithdrawStatusFail = 5
	//WithdrawStatusReturn     = 6 //退还成功
	WithdrawStatusReturn = 6
	//WithdrawStatusReturnFail = 7 //退还失败
	WithdrawStatusReturnFail = 7
)

type WithdrawOrderStatus int

func (s WithdrawOrderStatus) String() string {
	switch s {
	case WithdrawStatusWaitReview:
		return "待审核"
	case WithdrawStatusReject:
		return "审核拒绝"
	case WithdrawStatusPass:
		return "审核通过"
	case WithdrawStatusPaying:
		return "打款中"
	case WithdrawStatusSuccess:
		return "打款成功"
	case WithdrawStatusFail:
		return "打款失败"
	case WithdrawStatusReturn:
		return "退还成功"
	case WithdrawStatusReturnFail:
		return "退还失败"
	}
	return ""
}
