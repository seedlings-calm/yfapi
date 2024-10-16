package model

import (
	"database/sql"
	"time"
)

type TimelineReply struct {
	Id            int64          `json:"id" gorm:"column:id"`                         // 动态评论表主键
	ReplierId     string         `json:"replierId" gorm:"column:replier_id"`          // 评论人Id
	TimelineId    int64          `json:"timelineId" gorm:"column:timeline_id"`        // 动态Id
	ToReplyId     int64          `json:"toReplyId" gorm:"column:to_reply_id"`         // 对哪个评论的回复
	ToReplierId   sql.NullString `json:"toReplierId" gorm:"column:to_replier_id"`     // 对哪个人的回复
	ReplyContent  string         `json:"replyContent" gorm:"column:reply_content"`    // 评论内容
	Status        int            `json:"status" gorm:"column:status"`                 // 评论状态 0=隐藏 1=正常 2=审核未通过 3=删除
	SubReplyCount int            `json:"subReplyCount" gorm:"column:sub_reply_count"` // 子评论数
	ToSubReplyId  int64          `json:"toSubReplyId" gorm:"column:to_sub_reply_id"`  // 对哪个子评论的回复
	PraisedCount  int            `json:"praisedCount" gorm:"column:praised_count"`    // 被赞次数
	CreateTime    time.Time      `json:"createTime" gorm:"column:create_time"`        // 创建时间
	UpdateTime    time.Time      `json:"updateTime" gorm:"column:update_time"`        // 更新时间
}

func (m *TimelineReply) TableName() string {
	return "t_timeline_reply"
}
