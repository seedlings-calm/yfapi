package guild

import (
	"yfapi/util/easy"
)

type ChatroomTypeResp struct {
	ID       int    `json:"id" gorm:"column:id"`              //主键id
	TypeId   int    `json:"typeId" gorm:"column:type_id"`     // 房间类型ID
	TypeName string `json:"typeName" gorm:"column:type_name"` // 类型名称
	LiveType int8   `json:"liveType" gorm:"column:live_type"` // 房间直播类型 1=聊天室 2=个播 3=个人
	SortNum  int8   `json:"sortNum" gorm:"column:sort_num"`   // 排序
}
type GuildUserBaseInfo struct {
	UserId   string `json:"userId" gorm:"column:id"`
	UserNo   string `json:"userNo" gorm:"column:user_no"`     // 展示的用户id
	Nickname string `json:"nickname" gorm:"column:nickname"`  // 昵称
	Mobile   string `json:"mobile" gorm:"column:mobile"`      // 手机号
	Avatar   string `json:"avatar" gorm:"column:avatar"`      // 头像
	TrueName string `json:"trueName" gorm:"column:true_name"` // 真实姓名
}

// GuildRoomApplyListResp
// @Description: 公会房间申请列表响应
type GuildRoomApplyListResp struct {
	Id                  int64          `json:"id"`
	RoomNo              string         `json:"roomNo"`              // 房间id
	RoomName            string         `json:"roomName"`            // 房间名称
	RoomDesc            string         `json:"roomDesc"`            //房间描述
	RoomType            int            `json:"roomType"`            // 房间类型
	CoverImg            string         `json:"coverImg"`            //厅图
	UserNo              string         `json:"userNo"`              // 房主id
	NickName            string         `json:"nickName"`            // 房主昵称
	DaySettleNickname   string         `json:"daySettleNickname"`   // 日结算人昵称
	MonthSettleNickname string         `json:"monthSettleNickname"` // 月结算人昵称
	CreateTime          easy.LocalTime `json:"createTime"`          // 创建时间
	Status              int8           `json:"status"`              // 状态 1=待审核 2=审核通过 3=审核拒绝
	UpdateTime          easy.LocalTime `json:"updateTime"`          // 审核时间
	Reason              string         `json:"reason"`              // 拒绝原因
}

type GuildRoomListResp struct {
	Id                  string         `json:"id"`                  // 房间id
	RoomNo              string         `json:"roomNo"`              // 房间roomNo
	Name                string         `json:"name"`                // 房间名称
	Notice              string         `json:"notice"`              // 房间公告
	RoomType            int            `json:"roomType"`            // 房间类型
	CoverImg            string         `json:"coverImg"`            // 厅图
	UserNo              string         `json:"userNo"`              // 房主id
	NickName            string         `json:"nickName"`            // 房主昵称
	DaySettleUserNo     string         `json:"daySettleUserNo"`     // 日结算人userNo
	MonthSettleUserNo   string         `json:"monthSettleUserNo"`   // 月结算人userNo
	DaySettleNickname   string         `json:"daySettleNickname"`   // 日结算人昵称
	MonthSettleNickname string         `json:"monthSettleNickname"` // 月结算人昵称
	Status              int8           `json:"status"`              // 状态 1开启  2关闭 3作废
	CreateTime          easy.LocalTime `json:"createTime"`          // 创建时间
}

// RoomTypeInfo
// @Description: 房间类型
type RoomTypeInfo struct {
	TypeList []struct {
		TypeId   int    `json:"typeId"  ` // 房间类型ID
		TypeName string `json:"typeName"` // 类型名称
		LiveType int    `json:"liveType"` // 房间直播类型 1=聊天室 2=个播 3=个人
	} `json:"typeList"` // 房间类型列表
	TemplateList []struct {
		Id           string `json:"id"`           // 模板ID
		TemplateName string `json:"templateName"` // 模板名称
		LiveType     int    `json:"liveType"`     // 房间直播类型 1=聊天室 2=个播 3=个人
	} `json:"templateList"` // 房间模板列表
}
