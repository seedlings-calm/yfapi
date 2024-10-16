package model

import (
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/typedef/response/room"
)

type GuildRoom struct {
	Id                  string    `json:"id" gorm:"column:id"`
	UserId              string    `json:"userId" gorm:"column:user_id"` // 房主id
	UserNo              string    `json:"userNo"`
	UserNickName        string    `json:"userNickName"`                                         // 房主昵称
	RoomNo              string    `json:"roomNo" gorm:"column:room_no"`                         // 房间号
	FancyNoIcon         string    `json:"fancyNoIcon" gorm:"column:fancy_no_icon"`              // 靓号icon
	RoomType            int       `json:"roomType" gorm:"column:room_type"`                     // 房间类型
	LiveType            int       `json:"liveType" gorm:"column:live_type"`                     // 直播类型
	TemplateId          string    `json:"templateId" gorm:"column:template_id"`                 // 模板类型
	CoverImg            string    `json:"coverImg" gorm:"column:cover_img"`                     // 封面图
	Notice              string    `json:"notice" gorm:"column:notice"`                          // 房间公告
	Name                string    `json:"name" gorm:"column:name"`                              // 房间名称
	Status              int       `json:"status" gorm:"column:status"`                          // 房间状态 1开启  2关闭 3作废
	DaySettleUserId     string    `json:"daySettleUserId" gorm:"column:day_settle_user_id"`     // 日结算人id
	DaySettleUserNo     string    `json:"daySettleUserNo"`                                      // 日结算人用户编号
	DaySettleNickName   string    `json:"daySettleNickName"`                                    // 日结算人昵称
	MonthSettleUserId   string    `json:"monthSettleUserId" gorm:"column:month_settle_user_id"` // 月结算人id
	MonthSettleUserNo   string    `json:"monthSettleUserNo"`                                    // 月结算人用户No
	MonthSettleNickName string    `json:"monthSettleNickName"`
	GuildId             string    `json:"guildId" gorm:"column:guild_id"`       // 公会Id
	CreateTime          time.Time `json:"createTime" gorm:"column:create_time"` // 创建时间
	CreateTimeStr       string    `json:"createTimeStr"`                        // 创建时间
	UpdateTime          time.Time `json:"updateTime" gorm:"column:update_time"` // 更新时间
	UpdateTimeStr       string    `json:"updateTimeStr"`                        // 更新时间
}

func (m *GuildRoom) TableName() string {
	return "t_room"
}
func (m *GuildRoom) AfterFind(db *gorm.DB) (err error) {
	result := new(*GuildUser)
	m.UserNickName = (*result).GetUserInfo(db, m.UserId).Nickname
	m.DaySettleNickName = (*result).GetUserInfo(db, m.DaySettleUserId).Nickname
	m.MonthSettleNickName = (*result).GetUserInfo(db, m.MonthSettleUserId).Nickname
	m.UserNo = (*result).GetUserInfo(db, m.UserId).UserNo
	m.DaySettleUserNo = (*result).GetUserInfo(db, m.DaySettleUserId).UserNo
	m.MonthSettleUserNo = (*result).GetUserInfo(db, m.MonthSettleUserId).UserNo
	m.CoverImg = coreConfig.GetHotConf().ImagePrefix + m.CoverImg
	m.CreateTimeStr = m.CreateTime.Format("2006-01-02 15:04:05")
	m.UpdateTimeStr = m.UpdateTime.Format("2006-01-02 15:04:05")
	return
}
func (m *GuildRoom) ToChatroomDTO() room.ChatroomDTO {
	return room.ChatroomDTO{
		RoomId:     m.Id,
		RoomNo:     m.RoomNo,
		RoomName:   m.Name,
		CoverImg:   coreConfig.GetHotConf().ImagePrefix + m.CoverImg,
		Notice:     m.Notice,
		LiveType:   m.LiveType,
		RoomType:   m.RoomType,
		TemplateId: m.TemplateId,
		WellNoIcon: m.FancyNoIcon,
	}
}
