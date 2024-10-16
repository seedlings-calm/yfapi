package room

import response_goods "yfapi/typedef/response/goods"

// RoomShowBaseRes 展示的房间列表信息
type RoomShowBaseRes struct {
	Id            string `json:"id"`            // 对应id
	UserId        string `json:"userId"`        // 房主id
	RoomNo        string `json:"roomNo"`        // 房间号
	RoomType      int    `json:"roomType"`      // 房间类型
	LiveType      int    `json:"liveType"`      // 房间直播类型 1聊天室 2个播 3个人
	TemplateId    string `json:"templateId"`    // 模板类型
	CoverImg      string `json:"coverImg"`      // 封面图
	BackgroundImg string `json:"backgroundImg"` // 背景图
	Name          string `json:"name"`          // 房间名称
	Status        int    `json:"status"`        // 房间状态
	RoomPwd       string `json:"roomPwd"`       // 房间密码
}

// 麦位结构体
type RoomWheatPosition struct {
	Id       int               `json:"id"`                 //序号
	Identity string            `json:"identity"`           //麦位身份 （主持麦、嘉宾麦、音乐人麦、咨询师麦、普通麦）
	UserInfo RoomWheatUserInfo `json:"userInfo,omitempty"` //麦位上用户信息
	Status   int               `json:"status"`             //麦位状态 1:正常 2：在麦中 3：关闭
	Mute     bool              `json:"mute"`               // true 开启静音 false 关闭静音
}

type RoomWheatUserInfo struct {
	UserId     string                        `json:"userId"` // 在麦用户ID
	UserNo     int                           `json:"userNo"` // 在麦用户userNo
	Uid32      int32                         `json:"uid32"`
	UserName   string                        `json:"userName"`        //在麦用户昵称
	UserAvatar string                        `json:"userAvatar"`      //在麦用户头像
	CharmCount int                           `json:"charmCount"`      // 魅力值
	Frame      response_goods.SpecialEffects `json:"frame,omitempty"` // 麦位框装扮
	Voice      response_goods.SpecialEffects `json:"voice,omitempty"` // 声浪装扮
}
