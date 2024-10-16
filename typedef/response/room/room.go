package room

type ChatroomDTO struct {
	RoomId             string               `json:"roomId"`             // 房间ID
	RoomNo             string               `json:"roomNo"`             // 房间No
	RoomName           string               `json:"roomName"`           // 房间名称
	CoverImg           string               `json:"coverImg"`           // 房间封面
	Notice             string               `json:"notice"`             // 房间公告
	LiveType           int                  `json:"liveType"`           // 房间直播类型
	RoomType           int                  `json:"roomType"`           // 房间类型
	RoomTypeDesc       string               `json:"roomTypeDesc"`       // 房间类型描述
	TemplateId         string               `json:"templateId"`         // 模板ID
	IsLocked           bool                 `json:"isLocked"`           // 房间是否加锁
	IsCollect          bool                 `json:"isCollect"`          // 是否收藏房间
	BackgroundImg      string               `json:"backgroundImg"`      // 背景图
	WellNoIcon         string               `json:"wellNoIcon"`         // 房间靓号ICON
	SeatList           []*RoomWheatPosition `json:"seatList"`           // 房间座位列表
	IsPublicScreen     bool                 `json:"isPublicScreen"`     // 是否开启公屏
	IsFreedMic         bool                 `json:"isFreedMic"`         //自由上下麦状态 true开启  false关闭
	IsFreedSpeak       bool                 `json:"isFreedSpeak"`       //自由发言状态 true开启  false关闭
	IsMute             bool                 `json:"isMute"`             //关闭声音状态 true开启  false关闭
	AutoWelcomeContent string               `json:"autoWelcomeContent"` // 当前玩家设置的自动欢迎语
	GiftCategoryList   []GiftShowCategory   `json:"giftCategoryList"`   // 礼物展示类目列表
	GuildName          string               `json:"guildName"`          // 公会名称
	UpSeatApplyCount   int64                `json:"upSeatApplyCount"`   // 申请上麦人数
	RoleIdList         []int                `json:"roleIdList"`         // 当前房间的身份列表
	Hot                int                  `json:"hot"`                // 热度值
	HotStr             string               `json:"hotStr"`             // 热度值显示
	UserMute           bool                 `json:"userMute"`           //是否被禁言
	IsRelateWheat      bool                 `json:"isRelateWheat"`      //是否处于连麦中
	OwnerUserInfo      UpSeatApplyInfo      `json:"ownerUserInfo"`      // 房主信息
	IsOnHiddenMic      bool                 `json:"isOnHiddenMic"`      //是否再隐藏麦
}

type GiftShowCategory struct {
	CategoryId   int64  `json:"categoryId"`   // 类目ID
	CategoryName string `json:"categoryName"` // 类目名称
}

type UserRoomSwitchStatus struct {
	RoomMute           int `json:"roomMute"`           //房间静音 1开启 2关闭
	RoomSpecialEffects int `json:"roomSpecialEffects"` //房间特效 1开启 2关闭
}

type ChatroomExtra struct {
	RoomId         string `json:"roomId"`         // 房间ID
	IsPublicScreen bool   `json:"isPublicScreen"` // 是否开启公屏
	IsFreedMic     bool   `json:"isFreedMic"`     //自由上下麦状态 true开启  false关闭
	IsFreedSpeak   bool   `json:"isFreedSpeak"`   //自由发言状态 true开启  false关闭
}
