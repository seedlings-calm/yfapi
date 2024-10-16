package h5

import (
	"yfapi/typedef/response"
	"yfapi/util/easy"
)

type JoinGuildRes struct {
	Status int `json:"status"` //入会状态
}
type GuildInfoRes struct {
	GuildList      []*GuildInfo     `json:"guildList"`      //公会列表
	Guild          *GuildMsg        `json:"guild"`          //公会信息
	LiveRoomList   []*RoomInfo      `json:"liveRoomList"`   //直播间列表
	ChatRoomList   []*RoomInfo      `json:"chatRoomList"`   //聊天室列表
	LiveRoomCount  int64            `json:"LiveRoomCount"`  //直播间总数
	ChatRoomCount  int64            `json:"ChatRoomCount"`  //聊天室总数
	IsPractitioner bool             `json:"isPractitioner"` //是否为从业者
	QuitGuildApply *MemberApplyInfo `json:"quitGuildApply"` //申请退出公会信息
}
type GuildInfo struct {
	GuildId        string `json:"guildId"`        //公会id
	GuildNo        string `json:"guildNo"`        //公会编号
	GuildName      string `json:"guildName"`      //公会名称
	GuildLogo      string `json:"guildLogo"`      //公会logo
	GuildBriefDesc string `json:"guildBriefDesc"` //公会介绍
	Number         int    `json:"number"`         //成员数量
}
type RoomInfo struct {
	Id           string `json:"id"`           //房间id
	Name         string `json:"name"`         //房间名称
	RoomNo       string `json:"roomNo"`       //房间编号
	CoverImg     string `json:"coverImg"`     //封面图
	Status       int    `json:"status"`       //房间状态1开启  2关闭 3作废
	RoomType     int    `json:"type"`         //房型
	RoomTypeDesc string `json:"roomTypeDesc"` //房型描述
	Hot          int    `json:"hot"`          //热度
}
type GuildMsg struct {
	GuildId        string `json:"guildId"`        //公会id
	GuildNo        string `json:"guildNo"`        //公会编号
	GuildName      string `json:"guildName"`      //公会名称
	GuildLogo      string `json:"guildLogo"`      //公会logo
	GuildBriefDesc string `json:"guildBriefDesc"` //公会介绍
	NickName       string `json:"nickName"`       //会长昵称
	Avatar         string `json:"avatar"`         //会长头像
	Number         int    `json:"number"`         //成员数量
	ApplyStatus    int    `json:"applyStatus"`    //申请状态 1= 待审核  2= 同意 3=拒绝 4=自动拒绝 5=强制申请自动退出 6=取消申请

}

// GuildMemberListRes
// @Description: 获取工会成员列表
type GuildMemberListRes struct {
	GuildMemberList []*GuildMemberInfo `json:"guildMemberList"` //工会成员列表
}
type GuildMemberInfo struct {
	UserId     string                  `json:"userId"` //用户id
	UserNo     string                  `json:"userNo"` //用户编号
	Uid32      int32                   `json:"uid32"`
	Sex        int                     `json:"sex"`        //性别 0:保密 1=男 2=女
	NickName   string                  `json:"nickName"`   //用户昵称
	Avatar     string                  `json:"avatar"`     //用户头像
	Introduce  string                  `json:"introduce"`  //用户介绍
	Role       int8                    `json:"role"`       //角色 1=会长 2=成员
	UserPlaque response.UserPlaqueInfo `json:"userPlaque"` //用户铭牌
}

// GuildPenaltyDetailRes
// @Description: 违约金详情返回
type GuildPenaltyDetailRes struct {
	PenaltyList    []PenaltyDetail `json:"penaltyList"`
	StarLevel      int             `json:"starLevel"`      // 星光等级
	CurrExp        int             `json:"currExp"`        // 当前星光经验
	PenaltyDiamond string          `json:"penaltyDiamond"` // 违约金钻石
	DeductExp      int             `json:"deductExp"`      // 扣除星光经验
}

type PenaltyDetail struct {
	LevelName      string `json:"levelName"`      // 等级名称
	MinExp         int    `json:"minExp"`         // 最小经验值
	PenaltyRate    string `json:"penaltyRate"`    // 违约金比例
	PenaltyDiamond string `json:"penaltyDiamond"` // 违约金钻石
}

type MemberApplyInfo struct {
	GuildID    string         `json:"guildId" gorm:"column:guild_id"`     //工会id
	UserID     string         `json:"userId" gorm:"column:user_id"`       //用户ID
	ApplyType  int            `json:"applyType" gorm:"column:apply_type"` // 1=入会申请, 2=退会申请
	Force      int            `json:"force" gorm:"column:force"`          // 1=强制申请, 0=非强制申请
	Status     int            `json:"status" gorm:"column:status"`        // 1=待审核, 2=同意, 3=拒绝, 4=自动拒绝, 5=强制申请自动退出, 6=取消申请
	Reason     string         `json:"reason" gorm:"column:reason"`        // 拒绝原因
	CreateTime easy.LocalTime `json:"createTime" gorm:"column:create_time"`
	UpdateTime easy.LocalTime `json:"updateTime" gorm:"column:update_time"`
}
