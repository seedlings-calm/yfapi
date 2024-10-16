package enum

// 聊天室房间类型
const (
	// RoomTypeEmoMan 情感男
	RoomTypeEmoMan = 101
	// RoomTypeEmoWoman 情感女
	RoomTypeEmoWoman = 102
	// RoomTypeFriend 交友
	RoomTypeFriend = 103
	// RoomTypeDating 相亲
	RoomTypeDating = 104
	// RoomTypeSing 听歌
	RoomTypeSing = 105
	// RoomTypePodcast 播客
	RoomTypePodcast = 106
	// RoomTypeOrder 派单
	RoomTypeOrder = 107
)

// 个播房间类型
const (
	// RoomTypeAnchorVoice 个播 语音
	RoomTypeAnchorVoice = 201
	// RoomTypeAnchorVideo 个播 视频
	RoomTypeAnchorVideo = 202
)

// 个人房间类型
const (
	// RoomTypePersonal 个人
	RoomTypePersonal = 301
)

// 房间状态
const (
	// RoomStatusNormal 正常
	RoomStatusNormal = iota + 1
	// RoomStatusClose 关闭
	RoomStatusClose
	// RoomStatusInvalid 作废
	RoomStatusInvalid
)

// 房间前端显示状态
const (
	// RoomShowStatusClose  未开播（直播）
	RoomShowStatusClose = 0
	// RoomShowStatusAnchoring 直播中(个播房主播开播中)
	RoomShowStatusAnchoring = 1
	// RoomShowStatusInteraction 互动中(情感、听歌、交友，老板位有人)
	RoomShowStatusInteraction = 2
	// RoomShowStatusInteractWait 等待互动(情感、听歌、交友，老板位没有人)
	RoomShowStatusInteractWait = 3
	// RoomShowStatusSing 演唱中(听歌，演唱位有人)
	RoomShowStatusSing = 4
)

// 麦位状态
const (
	MicStatusNormal = 1 // 正常
	MicStatusUsed   = 2 // 在麦中
	MicStatusClose  = 3 // 关闭
)

const (
	SwitchOpen = 1 //开启
	SwitchOff  = 2 //关闭状态
)

type RoomType int

func (r RoomType) String() string {
	switch r {
	case RoomTypeEmoMan:
		return "情感男"
	case RoomTypeEmoWoman:
		return "情感女"
	case RoomTypeFriend:
		return "交友"
	case RoomTypeDating:
		return "相亲"
	case RoomTypeSing:
		return "听歌"
	case RoomTypePodcast:
		return "播客"
	case RoomTypeOrder:
		return "派单"
	case RoomTypeAnchorVoice:
		return "语音直播"
	case RoomTypeAnchorVideo:
		return "视频直播"
	case RoomTypePersonal:
		return "个人"
	default:
		return "未知"
	}
}
