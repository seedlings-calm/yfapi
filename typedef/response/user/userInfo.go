package user

import (
	"yfapi/typedef/response"
	response_goods "yfapi/typedef/response/goods"
)

type UserInfo struct {
	UserId              string                         `json:"userId"`              // 用户id
	Uid32               int32                          `json:"uid32"`               //int32类型uid
	UserNo              string                         `json:"userNo"`              // 用户No
	Nickname            string                         `json:"nickname"`            // 昵称
	Avatar              string                         `json:"avatar"`              // 头像
	LvCurrExp           int                            `json:"lvCurrExp"`           // lv当前经验
	LvMinExp            int                            `json:"lvMinExp"`            // lv最小经验
	LvMaxExp            int                            `json:"lvMaxExp"`            // lv最大经验
	IsOnline            bool                           `json:"isOnline"`            // 是否在线
	IsSetPwd            bool                           `json:"isSetPwd"`            // 是否设置密码
	VoiceUrl            string                         `json:"voiceUrl"`            // 语音介绍
	VoiceLength         int                            `json:"voiceLength"`         // 语音时长
	Introduce           string                         `json:"introduce"`           // 个性签名
	Sex                 int                            `json:"sex"`                 // 性别
	FriendNum           int                            `json:"friendNum"`           // 好友数量
	FansNum             int                            `json:"fansNum"`             // 粉丝数量
	FollowedNum         int                            `json:"followedNum"`         // 关注别人数量
	LikeNum             int                            `json:"likeNum"`             // 获赞数量
	VisitorNum          int                            `json:"visitorNum"`          // 访客数量
	BornDate            string                         `json:"bornDate"`            // 生日
	FollowedType        int                            `json:"followedType"`        //关注状态 0未关注对方 1已关注对方 2互相关注 3对方关注你
	RegionCode          string                         `json:"regionCode"`          // 区号
	Mobile              string                         `json:"mobile"`              // 手机号
	IsAnchor            bool                           `json:"isAnchor"`            // 是否为主播
	GuildName           string                         `json:"guildName"`           // 公会名称
	InRoom              *RoomInfo                      `json:"inRoom,omitempty"`    // 玩家在房信息
	Headwear            *response_goods.SpecialEffects `json:"headwear,omitempty"`  //头饰装扮
	IsBlacklist         bool                           `json:"isBlacklist"`         //是否被拉黑
	UserPlaque          response.UserPlaqueInfo        `json:"userPlaque"`          // 用户铭牌信息
	DontLetHeSeeMoments bool                           `json:"dontLetHeSeeMoments"` //不让他看动态
	DontSeeHeMoments    bool                           `json:"dontSeeHeMoments"`    //不看他的动态
	MomentsNoticeSwitch bool                           `json:"momentsNoticeSwitch"` //动态通知
	LiveNoticeSwitch    bool                           `json:"liveNoticeSwitch"`    //直播通知
	RoleIdList          []int                          `json:"roleIdList"`          //角色信息
	TrueName            string                         `json:"trueName"`            //真实姓名
	IdNo                string                         `json:"idNo"`                //身份证号
}

type RoomInfo struct {
	RoomId   string
	RoomName string
}

// 会话列表用户信息
type SessionListUserInfo struct {
	UserId       string                  `json:"userId"` //用户id
	UserNo       string                  `json:"userNo"`
	Uid32        int32                   `json:"uid32"`
	Nickname     string                  `json:"nickname"` //昵称
	Avatar       string                  `json:"avatar"`   //头像
	IsOnline     bool                    `json:"isOnline"` //是否在线
	Sex          int                     `json:"sex"`
	IsBlacklist  bool                    `json:"isBlacklist"`      //是否拉黑
	InRoom       *RoomInfo               `json:"inRoom,omitempty"` // 玩家在房信息
	UserPlaque   response.UserPlaqueInfo `json:"userPlaque"`       // 用户铭牌信息
	FollowedType int                     `json:"followedType"`     //关注状态 0未关注对方 1已关注对方 2互相关注 3对方关注你
}

type UserH5BasicInfo struct {
	UserId   string                         `json:"userId"` //用户id
	UserNo   string                         `json:"userNo"`
	Uid32    int32                          `json:"uid32"`
	Nickname string                         `json:"nickname"`           //昵称
	Avatar   string                         `json:"avatar"`             //头像
	Headwear *response_goods.SpecialEffects `json:"headwear,omitempty"` //头饰装扮
}

type RealNameResp struct {
	RealNameStatus int `json:"realNameStatus"` //实名状态 1未认证 2已认证
}

type UserAccount struct {
	Id       string `json:"id,omitempty"`
	UserNo   string `json:"userNo,omitempty"`
	Uid32    int32  `json:"uid32"`
	Avatar   string `json:"avatar,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Token    string `json:"token,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
}

// 用户隐私信息
type UserPrivateInfo struct {
	Id             string `json:"id"`
	Nickname       string `json:"nickname"`
	Qrcode         string `json:"qrcode"`
	Mobile         string `json:"mobile"`
	RealNameStatus string `json:"realNameStatus"` //实名状态 1未认证 3审核中 2认证成功
	RealName       string `json:"realName"`
	CardNum        string `json:"cardNum"`
}

type SearchUserByUserNoResp struct {
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
	UserNo   string `json:"userNo"`
	Avatar   string `json:"avatar"` // 头像
}

type UserRealNameInfoResp struct {
	RealName       string `json:"realName"`
	CardNum        string `json:"cardNum"`
	RealNameStatus string `json:"realNameStatus"`
}
