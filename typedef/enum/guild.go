package enum

const (
	// GuildMemberStatusNormal 公会成员状态 正常
	GuildMemberStatusNormal = iota + 1
	// GuildMemberStatusFreezing 公会成员状态 冻结
	GuildMemberStatusFreezing
	// GuildMemberStatusLeave 公会成员状态 脱离
	GuildMemberStatusLeave
)

const (

	// GuildStatusNormal 公会状态 正常
	GuildStatusNormal = iota + 1
	// GuildStatusInvalid 公会状态 作废
	GuildStatusInvalid
)

// GuildRoomApply
// 待审核 1 审核通过 2  审核拒绝 3
const (

	// GuildRoomApplyStatusWaitReview 公会房间申请状态 待审核
	GuildRoomApplyStatusWaitReview = iota + 1
	// GuildRoomApplyStatusPass 公会房间申请状态 审核通过
	GuildRoomApplyStatusPass
	// GuildRoomApplyStatusRefuse 公会房间申请状态 审核拒绝
	GuildRoomApplyStatusRefuse
)

// 工会成员申请状态
const (
	//公会成员申请状态 1= 待审核  2= 同意 3=拒绝 4=自动拒绝 5=强制申请自动退出 6=取消申请
	GuildMemberApplyStatusWait     = iota + 1 //待审核
	GuildMemberApplyStatusActive              //同意
	GuildMemberApplyStatusInactive            //拒绝
	GuildMemberApplyStatusDeleting            //自动拒绝
	GuildMemberApplyStatusDeleted             //强制申请自动退出
	GuildMemberApplyStatusCancel              //取消申请
)
