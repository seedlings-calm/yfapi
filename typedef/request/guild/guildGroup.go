package guild

import (
	"yfapi/typedef/request"
	"yfapi/util/easy"
)

type GuildGroupListreq struct {
	request.PageInfo
}
type AddGuildGroupReq struct {
	GroupName string `json:"groupName"` //分组标题
	Desc      string `json:"desc"`      //描述
}

// GetRoomRankListReq
// @Description: 查询房间排行榜列表请求
type GetRoomRankListReq struct {
	StartTime string `json:"startTime" form:"startTime"` // 开始日期 2024-09-25
	EndTime   string `json:"endTime" form:"endTime"`     // 结束日期 2024-09-25
}

type MemberGroupUpdateReq struct {
	Id        int    `json:"id" validate:"required"`
	GroupName string `json:"groupName" validate:"required"` //分组标题
	Desc      string `json:"desc" validate:"required"`      //描述
}

type MemberGroupUpdateRes struct {
	ID          int             `json:"id" gorm:"column:id;primaryKey"`     //主键id
	GuildID     string          `json:"guildId" gorm:"column:guild_id"`     //工会id
	GroupName   string          `json:"groupName" gorm:"column:group_name"` //分组名称
	Desc        string          `json:"desc" gorm:"column:desc"`            //分组描述
	Status      int8            `json:"status" gorm:"column:status"`        // 1=正常, 2=作废
	PersonCount int             `json:"personCount" validate:"required"`    //分组人数
	CreateTime  *easy.LocalTime `json:"createTime" gorm:"column:create_time"`
	UpdateTime  *easy.LocalTime `json:"updateTime" gorm:"column:update_time"`
}

type SetGroupByMembersReq struct {
	GroupId int   `json:"groupId" validate:"required"` //分组ID
	Ids     []int `json:"ids"`                         //成员列表ID
}
