package acl

import "yfapi/typedef/enum"

// 隐藏麦
const (
	UpHiddenMic   = "up_hidden_mic"   //上隐藏麦
	DownHiddenMic = "down_hidden_mic" //下隐藏麦
)

// 主持麦
const (
	UpCompereMic              = "up_compere_mic"                // 自己上麦
	HoldCompereUpCompereMic   = "hold_compere_up_compere_mic"   // 抱主持人上麦
	SwitchCompereMic          = "switch_compere_mic"            // 关闭麦位
	DownCompereMic            = "down_compere_mic"              // 下麦
	MuteCompereMic            = "mute_compere_mic"              // 静音麦位
	HoldCompereDownCompereMic = "hold_compere_down_compere_mic" // 抱用户下麦
	OutRoom                   = "out_room"                      // 踢出房间
	ReportUser                = "report_user"                   // 举报
	AddRoomBlacklist          = "add_room_blacklist"            // 拉黑
	ShutUp                    = "shut_up"                       // 禁言
)

// 嘉宾麦
const (
	UpGuestMic           = "up_guest_mic"             // 自己上麦
	HoldUserUpGuestMic   = "hold_user_up_guest_mic"   // 抱嘉宾上麦
	SwitchGuestMic       = "switch_guest_mic"         // 关闭麦位
	ApplyGuestMic        = "apply_guest_mic"          // 申请上麦
	DownGuestMic         = "down_guest_mic"           // 下麦
	MuteGuestMic         = "mute_guest_mic"           // 静音麦位
	TimerGuestMic        = "timer_guest_mic"          // 开启倒计时
	HoldUserDownGuestMic = "hold_user_down_guest_mic" // 抱用户下麦
)

// 音乐人麦
const (
	UpMusicianMic           = "up_musician_mic"             // 自己上麦
	HoldUserUpMusicianMic   = "hold_user_up_musician_mic"   // 抱用户上麦
	SwitchMusicianMic       = "switch_musician_mic"         // 关闭麦位
	ApplyMusicianMic        = "apply_musician_mic"          // 申请上麦
	DownMusicianMic         = "down_musician_mic"           // 下麦
	HoldUserDownMusicianMic = "hold_user_down_musician_mic" // 抱用户下麦
	MuteMusicianMic         = "mute_musician_mic"           // 静音麦位
)

// 咨询师麦
const (
	UpCounselorMic           = "up_counselor_mic"             // 自己上麦
	HoldUserUpCounselorMic   = "hold_user_up_counselor_mic"   // 抱用户上麦
	SwitchCounselorMic       = "switch_counselor_mic"         // 关闭麦位
	ApplyCounselorMic        = "apply_counselor_mic"          // 申请上麦
	DownCounselorMic         = "down_counselor_mic"           // 下麦
	HoldUserDownCounselorMic = "hold_user_down_counselor_mic" // 抱用户下麦
	MuteCounselorMic         = "mute_counselor_mic"           // 静音麦位
	TimerCounselorMic        = "timer_counselor_mic"          //开启倒计时
)

// 普通麦
const (
	UpNormalMic           = "up_normal_mic"             // 自己上麦
	HoldUserUpNormalMic   = "hold_user_up_normal_mic"   // 抱人上麦
	SwitchNormalMic       = "switch_normal_mic"         // 关闭麦位
	ApplyNormalMic        = "apply_normal_mic"          // 申请上麦
	MuteNormalMic         = "mute_normal_mic"           // 静音麦位
	DownNormalMic         = "down_normal_mic"           // 下麦
	TimerNormalMic        = "timer_normal_mic"          // 开启倒计时
	HoldUserDownNormalMic = "hold_user_down_normal_mic" // 抱用户下麦
)

// 权限
const (
	RedEnvelope             = "red_envelope"               // 红包
	MyCostume               = "my_costume"                 // 我的装扮
	WishList                = "wish_list"                  // 心愿单
	Greeting                = "greeting"                   // 自动欢迎语
	EditRoom                = "edit_room"                  // 编辑房间
	RoomSettingManager      = "room_setting_manager"       // 设置管理员
	RoomWarningMessage      = "room_warning_message"       // 警告信息
	HiddenRoom              = "hidden_room"                // 隐藏房间
	RoomOutAllMic           = "room_out_all_mic"           // 踢出全麦
	LockRoom                = "lock_room"                  // 锁定房间
	FreedMic                = "freed_mic"                  // 自由上下麦
	FreedSpeak              = "freed_speak"                // 自由发言
	RoomMute                = "room_mute"                  // 关闭声音
	RoomCloseSpecialEffects = "room_close_special_effects" // 关闭动效
	RoomBackgroundMusic     = "room_background_music"      // 背景音乐
	RoomClosePublicChat     = "room_close_public_chat"     // 关闭公屏
	RoomClearPublicChat     = "room_clear_public_chat"     // 清除公屏
	RoomClearMic            = "room_clear_mic"             // 清空麦位
	RoomBlackList           = "room_black_list"            // 黑名单
	RoomShutUpList          = "room_shut_up_list"          // 禁言单
	ResetGlamour            = "reset_glamour"              // 重置魅力值
	ChatRoomBackground      = "chat_room_background"       //聊天室背景
	RoomRelateWheat         = "room_relate_wheat"          //连麦操作
	// 自定义命令
	ClearUpSeatApply  = "clear_up_seat_apply"  // 清空申请上麦列表
	RefuseUpSeatApply = "refuse_up_seat_apply" // 拒绝上麦申请
	AcceptUpSeatApply = "accept_up_seat_apply" // 同意上麦申请
	CancelUpSeatApply = "cancel_up_seat_apply" // 取消上麦申请
)

type Menu struct {
	Name   string `json:"name"`                //权限名
	Title  string `json:"title"`               //展示标题
	Show   bool   `json:"-"`                   //是否隐藏
	Switch int    `json:"switchOff,omitempty"` //开关 1开启 2 关闭
}

var HiddenMic = []Menu{
	{
		Name:  UpHiddenMic,
		Title: "上隐藏麦",
	},
	{
		Name:  DownHiddenMic,
		Title: "下隐藏麦",
	},
}

// 互动包含权限
var HuDong = []Menu{
	//TODO 红包功能没做，暂时隐藏
	//{
	//	Name:  RedEnvelope,
	//	Title: "红包",
	//},
	{
		Name:  MyCostume,
		Title: "我的装扮",
	},
	{
		Name:  WishList,
		Title: "心愿单",
	},
	{
		Name:  Greeting,
		Title: "自动欢迎语",
	},
}

// 管理包含权限
var Manage = []Menu{
	{
		Name:  EditRoom,
		Title: "编辑房间",
	},
	{
		Name:  RoomSettingManager,
		Title: "设置管理员",
	},
	{
		Name:  RoomWarningMessage,
		Title: "警告信息",
	},
	{
		Name:   HiddenRoom,
		Title:  "隐藏房间",
		Switch: enum.SwitchOff,
	},
	{
		Name:  RoomOutAllMic,
		Title: "踢出全麦",
	},
	{
		Name:   LockRoom,
		Title:  "锁定房间",
		Switch: enum.SwitchOff,
	},
	{
		Name:   FreedMic,
		Title:  "自由上下麦",
		Switch: enum.SwitchOff,
	},
	{
		Name:   FreedSpeak,
		Title:  "自由发言",
		Switch: enum.SwitchOff,
	},
}

// 其他权限
var Other = []Menu{
	{
		Name:  RoomBackgroundMusic,
		Title: "背景音乐",
	},
	{
		Name:   RoomClosePublicChat,
		Title:  "开启公屏",
		Switch: enum.SwitchOff,
	},
	{
		Name:  RoomClearPublicChat,
		Title: "清除公屏",
	},
	{
		Name:  RoomClearMic,
		Title: "清空麦位",
	},
	{
		Name:  RoomBlackList,
		Title: "黑名单",
	},
	{
		Name:  RoomShutUpList,
		Title: "禁言单",
	},
	{
		Name:  ResetGlamour,
		Title: "重置魅力值",
	},
	{
		Name:  ChatRoomBackground,
		Title: "房间背景",
	},
	{
		Name:  RoomRelateWheat,
		Title: "连麦操作",
	},
}

// 更多按钮权限结构
var MoreMenu = map[string][]Menu{
	"interact": HuDong,
	"manage":   Manage,
	"more":     Other,
}

// 用户资料卡权限
var UserCard = []Menu{
	{
		Name:  OutRoom,
		Title: "踢出房间",
	},
	{
		Name:  ReportUser,
		Title: "举报",
	},
	{
		Name:  AddRoomBlacklist,
		Title: "拉黑",
	},
	{
		Name:  ShutUp,
		Title: "禁言",
	},
}

// 主持麦位用户资料卡权限
var ZhuChiUserMic = []Menu{
	{
		Name:  DownCompereMic,
		Title: "下麦",
	},
	{
		Name:   MuteCompereMic,
		Title:  "静音麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  HoldCompereDownCompereMic,
		Title: "抱主持下麦",
	},
	{
		Name:  OutRoom,
		Title: "踢出房间",
	},
	{
		Name:  ReportUser,
		Title: "举报",
	},
	{
		Name:  AddRoomBlacklist,
		Title: "拉黑",
	},
	{
		Name:  ShutUp,
		Title: "禁言",
	},
}

// 嘉宾麦位资料卡权限
var JiaBinUserMic = []Menu{
	{
		Name:  DownGuestMic,
		Title: "下麦",
	},
	{
		Name:   MuteGuestMic,
		Title:  "静音麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  TimerGuestMic,
		Title: "开启倒计时",
	},
	{
		Name:  HoldUserDownGuestMic,
		Title: "抱用户下麦",
	},
	{
		Name:  OutRoom,
		Title: "踢出房间",
	},
	{
		Name:  ReportUser,
		Title: "举报",
	},
	{
		Name:  AddRoomBlacklist,
		Title: "拉黑",
	},
	{
		Name:  ShutUp,
		Title: "禁言",
	},
}

// 音乐人资料卡权限
var YinYueRenUserMic = []Menu{
	{
		Name:  DownMusicianMic,
		Title: "下麦",
	},
	{
		Name:   MuteMusicianMic,
		Title:  "静音麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  HoldUserDownMusicianMic,
		Title: "抱用户下麦",
	},
	{
		Name:  OutRoom,
		Title: "踢出房间",
	},
	{
		Name:  ReportUser,
		Title: "举报",
	},
	{
		Name:  AddRoomBlacklist,
		Title: "拉黑",
	},
	{
		Name:  ShutUp,
		Title: "禁言",
	},
}

// 咨询师麦位权限
var ZiXunShiUserMic = []Menu{
	{
		Name:  DownCounselorMic,
		Title: "下麦",
	},
	{
		Name:   MuteCounselorMic,
		Title:  "静音麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  TimerCounselorMic,
		Title: "开启倒计时",
	},
	{
		Name:  HoldUserDownCounselorMic,
		Title: "抱用户下麦",
	},
	{
		Name:  OutRoom,
		Title: "踢出房间",
	},
	{
		Name:  ReportUser,
		Title: "举报",
	},
	{
		Name:  AddRoomBlacklist,
		Title: "拉黑",
	},
	{
		Name:  ShutUp,
		Title: "禁言",
	},
}

// 普通麦权限
var NormalUserMic = []Menu{
	{
		Name:   MuteNormalMic,
		Title:  "静音麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  TimerNormalMic,
		Title: "开启倒计时",
	},
	{
		Name:   DownNormalMic,
		Title:  "下麦",
		Switch: enum.SwitchOff,
	},
	{
		Name:  HoldUserDownNormalMic,
		Title: "抱用户下麦",
	},
	{
		Name:  OutRoom,
		Title: "踢出房间",
	},
	{
		Name:  ReportUser,
		Title: "举报",
	},
	{
		Name:  AddRoomBlacklist,
		Title: "拉黑",
	},
	{
		Name:  ShutUp,
		Title: "禁言",
	},
}

// 主持麦空麦位权限
var ZhuChiEmptyMic = []Menu{
	{
		Name:  UpCompereMic,
		Title: "自己上麦",
	},
	{
		Name:  HoldCompereUpCompereMic,
		Title: "抱主持人上麦",
	},
	{
		Name:   SwitchCompereMic,
		Title:  "关闭麦位",
		Switch: enum.SwitchOff,
	},
}

// 嘉宾麦空麦权限
var JiaBinEmptyMic = []Menu{
	{
		Name:  UpGuestMic,
		Title: "自己上麦",
	},
	{
		Name:  HoldUserUpGuestMic,
		Title: "抱嘉宾上麦",
	},
	{
		Name:   SwitchGuestMic,
		Title:  "关闭麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  ApplyGuestMic,
		Title: "申请上麦",
	},
}

// 音乐人空麦权限
var YinYueRenEmptyMic = []Menu{
	{
		Name:  UpMusicianMic,
		Title: "自己上麦",
	},
	{
		Name:  HoldUserUpMusicianMic,
		Title: "抱用户上麦",
	},
	{
		Name:   SwitchMusicianMic,
		Title:  "关闭麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  ApplyMusicianMic,
		Title: "申请上麦",
	},
}

// 咨询师空麦权限
var ZiXunShiEmptyMic = []Menu{
	{
		Name:  UpCounselorMic,
		Title: "自己上麦",
	},
	{
		Name:  HoldUserUpCounselorMic,
		Title: "抱用户上麦",
	},
	{
		Name:   SwitchCounselorMic,
		Title:  "关闭麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  ApplyCounselorMic,
		Title: "申请上麦",
	},
}

// 普通麦空麦权限
var NormalEmptyMic = []Menu{
	{
		Name:  UpNormalMic,
		Title: "自己上麦",
	},
	{
		Name:  HoldUserUpNormalMic,
		Title: "抱人上麦",
	},
	{
		Name:   SwitchNormalMic,
		Title:  "关闭麦位",
		Switch: enum.SwitchOff,
	},
	{
		Name:  ApplyNormalMic,
		Title: "申请上麦",
	},
}
