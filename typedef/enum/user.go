package enum

import "time"

// 用户状态
const (
	// UserStatusNormal 用户状态 正常
	UserStatusNormal = iota + 1
	// UserStatusFreezing 用户状态 冻结
	UserStatusFreezing
	// UserStatusApplyInvalid 用户状态 申请注销
	UserStatusApplyInvalid
	// UserStatusInvalid 用户状态 已注销
	UserStatusInvalid
)

const (
	// OfficialUserId 官方用户ID(系统通知....)
	OperationUserId   = "1000" //运营资金账号
	SystematicUserId  = "1001" //系统通知账号
	OfficialUserId    = "1002" //官方公告账号
	InteractiveUserId = "1003" //互动消息账号
)

var OfficialUserIdList = []string{OperationUserId, SystematicUserId, OfficialUserId, InteractiveUserId}

// 账号注销申请状态
const (
	// UserDeleteStatusApplying 账号注销状态 已申请
	UserDeleteStatusApplying = iota + 1
	// UserDeleteStatusCancel 账号注销状态 已取消
	UserDeleteStatusCancel
	// UserDeleteStatusDelete  账号注销状态 已注销
	UserDeleteStatusDelete
)

// 用户性别
const (
	// UserSexTypeMan 用户性别 男
	UserSexTypeMan = 1
	// UserSexTypeWoman 用户性别 女
	UserSexTypeWoman = 2
)

// 从业者身份标识
const (
	UserPractitionerCompere   = 1 //主持人身份
	UserPractitionerMusician  = 2 //音乐人
	UserPractitionerCounselor = 3 //咨询师身份
	UserPractitionerAnchor    = 4 //主播身份
)

type PractitionerType int

func (p PractitionerType) String() string {
	switch p {
	case UserPractitionerCompere:
		return "主持人"
	case UserPractitionerMusician:
		return "音乐人"
	case UserPractitionerCounselor:
		return "咨询师"
	case UserPractitionerAnchor:
		return "主播"
	default:
		return "未知类型"
	}
}

// 实名认证状态
const (
	UserRealNameUnverified    = 1 //未认证
	UserRealNameAuthenticated = 2 //已认证
	UserRealNameVerifying     = 3 //认证中
)

type UserRealName int

func (u UserRealName) String() string {
	switch u {
	case UserRealNameUnverified:
		return "未认证"
	case UserRealNameAuthenticated:
		return "已认证"
	case UserRealNameVerifying:
		return "认证中"
	default:
		return "未知类型"
	}
}

const (
	RedisTokenExpireLifeTime = time.Hour * 24 * 7
	TokenExpireLifeTime      = time.Hour * 24 * 365
	TokenExpireChooseUser    = time.Hour * 1
)

const (
	BlacklistTypeRoom = 1 // 房间拉黑
	BlacklistTypeUser = 2 // 用户拉黑
)

// 动态控制开关
const (
	DontLetHeSeeMoments = 1 //不让他看动态
	DontSeeHeMoments    = 2 //不看他动态
)

// 实名热证审核状态
// 审核状态 1:待审核 2:审核通过 3:审核拒绝
const (
	UserRealNameReviewStatusWait   = iota + 1 //待审核
	UserRealNameReviewStatusPass              //审核通过
	UserRealNameReviewStatusReject            //审核拒绝
)

type UserRealNameReviewStatus int

func (UserRealNameReviewStatus) String() string {
	switch UserRealNameReviewStatusWait {
	case UserRealNameReviewStatusWait:
		return "待审核"
	case UserRealNameReviewStatusPass:
		return "审核通过"
	case UserRealNameReviewStatusReject:
		return "审核拒绝"
	default:
		return "未知类型"
	}
}

// 通知开关
const (
	MomentsNoticeSwitch = 1 //动态通知
	LiveNoticeSwitch    = 2 //直播通知
)
