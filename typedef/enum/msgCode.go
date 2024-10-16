package enum

const (
	MsgListTextColorNormal = "#999999" //正常
)
const (
	MsgAudio  = "audio"  //语音
	MsgText   = "text"   //文本
	MsgImg    = "img"    //图片
	MsgCustom = "custom" //自定义
	MsgAction = "action" //动作消息
	MsgGift   = "gift"
)

const (
	ALL_SERVICE_MESSAGE_TYPE  = 0 //全服消息
	ROOM_MESSAGE_TYPE         = 1 //房间消息
	USER_MESSAGE_TYPE         = 2 //用户消息
	USER_MESSAGE_CLIENT_TYYPE = 3 //指定端用户
)

const (
	ChatMessagePass          = 0 //私聊消息通过
	ChatMessageReview        = 1 //私聊消息可疑
	ChatMessageReject        = 2 //私聊消息拦截
	PrivateChatMessageRead   = 1 //私聊消息已读
	PrivateChatMessageUnRead = 2 //私聊消息未读
)

// 10001--10999 用户相私聊相关操作
const (
	USER_LOGIN_OTHER_CLIENT_MSG = 10001 //用户在其他端登录 退出登录
	USER_TEXT_MSG               = 10002 //用户单聊文本信息
	USER_IMG_MSG                = 10003 //用户单聊图片信息
	USER_AUDIO_MSG              = 10004 //用户单聊语音信息
	USER_CUSTOM_MSG             = 10005 //用户单聊自定义信息
	USER_SYSTEM_MSG             = 10006 //系统消息
	USER_OFFICIAL_MSG           = 10007 //官方公告
	USER_INTERACTIVE_MSG        = 10008 //互动消息
)

// 20001--20999  房间消息
const (
	JOIN_ROOM_MSG                    = 20001 //加入房间
	LEAVE_ROOM_MSG                   = 20002 //离开房间
	ROOM_TEXT_MSG                    = 20003 //房间公屏文本信息
	ROOM_IMG_MSG                     = 20004 //房间公屏图片信息
	ROOM_CUSTOM_MSG                  = 20005 //房间自定义信息
	ROOM_GIFT_MSG                    = 20006 //房间礼物消息
	ROOM_ACTION_MSG                  = 20007 //房间动作消息
	BLACKOUT_ROOM_MSG                = 20009 //操作拉黑，踢出房间
	KICKOUT_ROOM_MSG                 = 20010 //操作踢出，踢出房间
	ONLINE_CHANGE_THREE_MSG          = 20011 //在线列表前三榜单变化
	AUTO_WELCOME_MSG                 = 20012 //自动欢迎语设置
	ROOM_SEND_GIFT_SEAT_MSG          = 20013 //打赏礼物飞屏动效麦位列表
	ROOM_HOT_UPDATE_MSG              = 20014 //房间热度值更新
	USER_LEVEL_UP_MSG                = 20015 //用户等级升级
	USER_JOIN_ROOM_FROM_OTHER_CLIENT = 20016 //用户在其他端进入房间
	ROOM_FORCE_OFF_MSG               = 20017 //房间强制关播通知
	// 20100 ~ 20150 麦位相关通知

	UP_SEAT_MSG           = 20100 //用户上麦
	DOWN_SEAT_MSG         = 20101 //用户下麦
	MUTE_SEAT_MSG         = 20102 //麦位静音
	CLOSE_SEAT_MSG        = 20103 //麦位关闭
	APPLY_SEAT_MSG        = 20104 //麦位申请列表更新
	HOLD_UP_SEAT_MSG      = 20105 //用户被抱上麦
	HOLD_DOWN_SEAT_MSG    = 20106 //用户被抱下麦
	FREED_MIC_MSG         = 20107 //自由上下麦状态更新
	FREED_SPEAK_MSG       = 20108 //自由发言状态更新
	CLOSE_PUBLIC_CHAT_MSG = 20109 //关闭公屏状态更新
	CLEAR_PUBLIC_CHAT_MSG = 20110 //清除公屏通知
	HIDDEN_ROOM_MSG       = 20111 //隐藏房间状态更新
	ROOM_WARNING_MSG      = 20112 //房间警告信息
	SEAT_CHARM_MSG        = 20113 //座位魅力值更新
	TIMER_OPEN_MSG        = 20114 //开启倒计时
	MUTE_USER_MSG         = 20115 //用户禁言更新
	SEAT_ACTION_MSG       = 20116 //公屏麦位动作消息
	RESET_CHARM_MSG       = 20117 //重置魅力值
	APPLY_SEAT_RESULT_MSG = 20118 //申请上麦结果通知
	Room_Relate_Wheat_MSG = 20119 //连麦操作
)

// 房间外部刷新，变更等通知
// 30000-31999
const (
	Room_BackGroudImg_Update = 30001 //房间背景更换推送
	User_Goods_Change        = 30002 //用户的装扮信息更换
	Room_Notice_Update       = 30003 //房间公告更新
)

// 用户个人消息
// 40000-41999
const (
	//用户踢下线
	USER_KICK_OUT_MSG = 40001 //用户踢下线
)

// 群聊相关
// 50000-50999
const (
	GROUP_TEXT_MSG       = 50001 //普通群聊文本消息
	GROUP_IMG_MSG        = 50002 //普通群聊图片消息
	GROUP_ADMIN_TEXT_MSG = 50003 //超管文本消息
	GROUP_ADMIN_IMG_MSG  = 50004 //超管图片消息
	GROUP_SUBSCRIBE      = 50005 //订阅群聊
	GROUP_UN_SUBSCRIBE   = 50006 //取消群聊订阅
	GROUP_MUTE_MSG       = 50007 //群聊禁言通知
	GROUP_UNMUTE_MSG     = 50008 //取消禁言
)
