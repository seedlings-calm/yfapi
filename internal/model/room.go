package model

import (
	"time"
	"yfapi/core/coreConfig"
	"yfapi/typedef/response/room"
)

type Room struct {
	Id                  string    `json:"id" gorm:"column:id"`
	UserId              string    `json:"userId" gorm:"column:user_id"` // 房主id
	UserInfo            User      `gorm:"foreignKey:UserId"`
	RoomNo              string    `json:"roomNo" gorm:"column:room_no"`                     // 房间号
	FancyNoIcon         string    `json:"fancyNoIcon" gorm:"column:fancy_no_icon"`          // 靓号icon
	RoomType            int       `json:"roomType" gorm:"column:room_type"`                 // 房间类型
	LiveType            int       `json:"liveType" gorm:"column:live_type"`                 // 直播类型
	TemplateId          string    `json:"templateId" gorm:"column:template_id"`             // 模板类型
	CoverImg            string    `json:"coverImg" gorm:"column:cover_img"`                 // 封面图
	Notice              string    `json:"notice" gorm:"column:notice"`                      // 房间公告
	Name                string    `json:"name" gorm:"column:name"`                          // 房间名称
	IsHot               bool      `json:"isHot" gorm:"column:is_hot"`                       // 是否热门
	Status              int       `json:"status" gorm:"column:status"`                      // 房间状态 1开启  2关闭 3作废
	RoomPwd             string    `json:"roomPwd" gorm:"column:room_pwd"`                   // 房间密码
	BackgroundImg       string    `json:"BackgroundImg"`                                    // 房间背景
	UpSeatType          int       `json:"upSeatType" gorm:"column:up_seat_type"`            // 上麦模式 1=排麦  2=自由麦
	DaySettleUserId     string    `json:"daySettleUserId" gorm:"column:day_settle_user_id"` // 日结算人id
	DaySettleUserInfo   User      `gorm:"foreignKey:DaySettleUserId"`
	MonthSettleUserId   string    `json:"monthSettleUserId" gorm:"column:month_settle_user_id"` // 月结算人id
	MonthSettleUserInfo User      `gorm:"foreignKey:MonthSettleUserId"`
	GuildId             string    `json:"guildId" gorm:"column:guild_id"`                        // 公会Id
	CreateTime          time.Time `json:"createTime" gorm:"column:create_time"`                  // 创建时间
	UpdateTime          time.Time `json:"updateTime" gorm:"column:update_time"`                  // 更新时间
	HiddenStatus        int       `json:"hiddenStatus" gorm:"column:hidden_status"`              // 隐藏房间 2关闭 1开启
	PublicScreenStatus  int       `json:"publicScreenStatus" gorm:"column:public_screen_status"` // 公屏显示 2关闭 1开启
	FreedMicStatus      int       `json:"freedMicStatus" gorm:"column:freed_mic_status"`         // 自由上下麦 2关闭 1 开启
	FreedSpeakStatus    int       `json:"freedSpeakStatus" gorm:"column:freed_speak_status"`     // 自由发言 2关闭 1开启
}

func (m *Room) TableName() string {
	return "t_room"
}

func (m *Room) ToShowBase() room.RoomShowBaseRes {
	return room.RoomShowBaseRes{
		Id:         m.Id,
		UserId:     m.UserId,
		RoomNo:     m.RoomNo,
		RoomType:   m.RoomType,
		LiveType:   m.LiveType,
		TemplateId: m.TemplateId,
		CoverImg:   formatUrl(m.CoverImg),
		Name:       m.Name,
		Status:     m.Status,
		RoomPwd:    m.RoomPwd,
	}
}

func (m *Room) ToChatroomDTO() room.ChatroomDTO {
	return room.ChatroomDTO{
		RoomId:        m.Id,
		RoomNo:        m.RoomNo,
		RoomName:      m.Name,
		CoverImg:      formatUrl(m.CoverImg),
		Notice:        m.Notice,
		LiveType:      m.LiveType,
		RoomType:      m.RoomType,
		TemplateId:    m.TemplateId,
		WellNoIcon:    m.FancyNoIcon,
		BackgroundImg: formatUrl(m.BackgroundImg),
	}
}

func formatUrl(url string) string {
	if len(url) == 0 {
		return ""
	}
	return coreConfig.GetHotConf().ImagePrefix + url
}
