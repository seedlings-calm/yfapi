package enum

const (
	SuperAdminRoleId = 1001 //超管角色ID
	PatrolRoleId     = 1002 //巡查角色ID
	PresidentRoleId  = 1003 //公会会长角色ID
	HouseOwnerRoleId = 1004 //房主角色ID
	CompereRoleId    = 1005 //主持人角色Id
	RoomAdminRoleId  = 1006 //房间管理角色ID
	MusicianRoleId   = 1007 //音乐人角色Id
	CounselorRoleId  = 1008 //咨询师角色ID
	AnchorRoleId     = 1009 //主播角色ID
	NormalRoleId     = 1010 //普通用户
)

// 麦位定义
const (
	HiddenMicSeat    = "hiddenMicSeat"    //隐藏麦
	CompereMicSeat   = "compereMicSeat"   //主持麦
	GuestMicSeat     = "guestMicSeat"     //嘉宾麦
	MusicianMicSeat  = "musicianMicSeat"  //音乐人麦
	CounselorMicSeat = "counselorMicSeat" //咨询师麦
	NormalMicSeat    = "normalMicSeat"    //普通麦
)

type RoleType int

// PlaqueTypeName
//
//	@Description: 身份铭牌类型名称
//	@receiver r
//	@return string -
func (r RoleType) PlaqueTypeName() string {
	switch r {
	case SuperAdminRoleId:
		return "chao"
	case PatrolRoleId:
		return "xun"
	case PresidentRoleId:
		return "hui"
	case HouseOwnerRoleId:
		return "fang"
	case CompereRoleId:
		return "zhu"
	case RoomAdminRoleId:
		return "guan"
	case MusicianRoleId:
		return ""
	case CounselorRoleId:
		return ""
	case AnchorRoleId:
		return ""
	case NormalRoleId:
		return ""
	}
	return ""
}
