package room

import "yfapi/typedef/response"

// 资料卡
type InformationCardResponse struct {
	UserId             string                  `json:"userId"`             //用户id
	UserNo             string                  `json:"userNo"`             //用户序号
	Nickname           string                  `json:"nickname"`           //昵称
	Avatar             string                  `json:"avatar"`             //头像
	Sex                int                     `json:"sex"`                //性别
	Introduce          string                  `json:"introduce"`          //个人简介
	IsTrueName         bool                    `json:"isTrueName"`         //是否实名
	IsFollow           int                     `json:"isFollow"`           //是否关注此用户
	IsOnline           bool                    `json:"isOnline"`           //是否在房间
	FansNum            int64                   `json:"fansNum"`            //粉丝数量
	FollowedNum        int64                   `json:"followedNum"`        //关注别人数量
	GuildName          string                  `json:"guildName"`          //公会名称 最多六个字符
	Practitions        []int                   `json:"practitions"`        //从业者身份
	RoleIdList         []int                   `json:"roleIdList"`         // 当前房间的身份列表
	GiftExhibitionHall any                     `json:"giftExhibitionHall"` //礼物展馆 TODO:
	GiftWall           any                     `json:"giftWall"`           //礼物墙 TODO:
	RechargeNameplate  bool                    `json:"rechargeNameplate"`  //充值铭牌，后台设置，TODO:
	UserPlaque         response.UserPlaqueInfo `json:"userPlaque"`         // 用户铭牌信息
}

// 用户的铭牌信息   TODO:  调用的都为实现，
type UserIdCards struct {
	VIP       string `json:"vip"`
	LV        string `json:"lv"`
	Starlight string `json:"starlight"` //是否展示星光  该房间从业者展示星光铭牌，非从业者不展示星光铭牌
}

type BlackListResponse struct {
	Count int                    `json:"count"`
	List  []BlackListAndUserInfo `json:"list"`
}

// 黑名单和用户关联结构体
type BlackListAndUserInfo struct {
	UserId     string                  `json:"userId"`     //用户id
	UserNo     string                  `json:"userNo"`     //用户序号
	Nickname   string                  `json:"nickname"`   //昵称
	Avatar     string                  `json:"avatar"`     //头像
	Sex        int                     `json:"sex"`        //性别;0:保密,1:男,2:女
	Introduce  string                  `json:"introduce"`  // 个性签名
	UserPlaque response.UserPlaqueInfo `json:"userPlaque"` // 用户铭牌信息
}

// 在线用户列表返回结构体
type OnlineUsersResponse struct {
	IsShowNums       bool             `json:"isShowNums"`       //是否展示贡献值
	DayUsersCount    int              `json:"dayUsersCount"`    // 1000贡献榜人数
	OnlineLists      []*RoomUsersBase `json:"onlineLists"`      //在线用户列表
	OnlineUsersCount int              `json:"onlineUsersCount"` //在线用户数量，最高200人
	OwnerInfo        OwnerInfo        `json:"ownerInfo"`        //自己的信息
}

type OwnerInfo struct {
	RoomUsersBase
	UpgradeNum           string `json:"upgradeNum"`           //距离上榜 差多少贡献值
	OriginalContribution string `json:"originalContribution"` //未处理贡献值
}

type RoomUsersBase struct {
	UserId       string                  `json:"userId"` //用户id
	UserNo       string                  `json:"userNo"` //用户序号
	Uid32        int32                   `json:"uid32"`
	Nickname     string                  `json:"nickname"`     //昵称
	Avatar       string                  `json:"avatar"`       //头像
	Sex          int                     `json:"sex"`          //性别;0:保密,1:男,2:女
	Contribution string                  `json:"contribution"` //贡献值
	UserPlaque   response.UserPlaqueInfo `json:"userPlaque"`   // 用户铭牌信息
}

// 1000贡献榜返回结构体
type DayUsersResponse struct {
	OnlineUsersCount   int              `json:"onlineUsersCount"`   //在线用户数量
	NoOnlineUsersCount int              `json:"noOnlineUsersCount"` //不在线用户数量
	OnlineUsers        []*RoomUsersBase `json:"onlineUsers"`        //在线列表
	NoOnlineUsers      []*RoomUsersBase `json:"NoOnlineUsers"`      //不在线列表
	OwnerInfo          OwnerInfo        `json:"onwerInfo"`          //自己的信息
}

// 高等级用户接口返回体
type HighGradeUsersResponse struct {
	FirstCount  int              `json:"firstCount"`  //51级以上
	SecondCount int              `json:"secondCount"` //41-50级
	ThreeCount  int              `json:"threeCount"`  // 31-40级
	ThreeList   []*RoomUsersBase `json:"threeList"`
	FirstList   []*RoomUsersBase `json:"firstList"`
	SecondList  []*RoomUsersBase `json:"secondList"`
	OwnerInfo   RoomUsersBase    `json:"onwerInfo"` //自己的信息
}

type HighGradeUsersCountResponse struct {
	LevelList []HighGradeUsersCountItem
}

type HighGradeUsersCountItem struct {
	Rule  string `json:"rule"`  //规则  >=51 high  >40 <=50  mid   <40 low
	Icon  string `json:"icon"`  //图标路径
	Count int    `json:"count"` //统计人数
}
