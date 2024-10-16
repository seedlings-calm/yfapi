package index

import (
	"yfapi/typedef/response"
	"yfapi/typedef/response/room"
	"yfapi/typedef/response/user"
)

type SearchUserInfo struct {
	UserId       string                  `json:"userId"` //用户id
	UserNo       string                  `json:"userNo"` //用户序号
	Uid32        int32                   `json:"uid32"`
	Nickname     string                  `json:"nickname"`         //昵称
	Avatar       string                  `json:"avatar"`           //头像
	Sex          int                     `json:"sex"`              // 性别
	Introduce    string                  `json:"introduce"`        // 个性签名
	IsOnline     bool                    `json:"isOnline"`         // 是否在线
	FollowedType int                     `json:"followedType"`     //关注状态 0未关注对方 1已关注对方 2互相关注 3对方已关注
	InRoom       *user.RoomInfo          `json:"inRoom,omitempty"` // 玩家在房信息
	UserPlaque   response.UserPlaqueInfo `json:"userPlaque"`       // 用户铭牌信息
}

// SearchAllRes 搜索用户、聊天室、直播间 返回
type SearchAllRes struct {
	UserList       []*SearchUserInfo `json:"userList"`       // 用户
	ChatroomList   []*room.RoomInfo  `json:"chatroomList"`   // 聊天室
	AnchorRoomList []*room.RoomInfo  `json:"anchorRoomList"` // 直播间
}

// AppMenuSetting
// @Description: app菜单配置
type AppMenuSetting struct {
	MenuName string `json:"menu_name"` // 菜单名称
	Icon     string `json:"icon"`      // 菜单icon
	LinkUrl  string `json:"link_url"`  // 跳转地址
}

type TopMsgRes struct {
	UserId    string      `json:"userId"`    //打赏人id
	Nickname  string      `json:"nickname"`  //打赏人昵称
	Avatar    string      `json:"avatar"`    //打赏人头像
	Types     string      `json:"types"`     //类型 全服礼物
	Operate   string      `json:"operate"`   //操作 打赏
	ToUser    interface{} `json:"toUser"`    //收礼人
	GiftImg   string      `json:"giftImg"`   //礼物图片
	GiftName  string      `json:"giftName"`  //礼物名称
	GiftCount int         `json:"giftCount"` //礼物数量
}
