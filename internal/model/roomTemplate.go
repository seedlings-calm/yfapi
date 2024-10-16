package model

import "time"

type RoomTemplate struct {
	Id            string    `json:"id"            description:"模板ID"`
	TemplateName  string    `json:"templateName"  description:"模板名称"`
	LiveType      int       `json:"liveType"      description:"房间直播类型 1=聊天室 2=个播 3=个人"`
	BgImg         string    `json:"bgImg"         description:"默认房间背景图"`
	SeatListCount int       `json:"seatListCount" description:"嘉宾位数量"`
	Status        int       `json:"status"        description:"状态 1=启用 2=禁用"`
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time"` // 更新时间
}

func (m *RoomTemplate) TableName() string {
	return "t_room_template"
}
