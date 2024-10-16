package model

import "time"

// 举报中心
type ReportingCenter struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`                                                      // 主键，自动递增
	SrcUserID     string    `gorm:"size:100;not null;default:'';comment:'举报用户ID'"`                                 // 举报用户ID
	DstUserID     string    `gorm:"size:100;not null;default:'';comment:'被举报人/房ID'"`                               // 被举报人/房ID
	Object        int       `gorm:"not null;default:1;comment:'举报对象 1房间 2用户'"`                                     // 举报对象
	Scene         int       `gorm:"not null;default:1;comment:'举报场景 1房间 2消息私聊 3个人主页 4个人动态 5动态评价 6声音派对'"`           // 举报场景
	ReportType    int       `gorm:"not null;default:99;comment:'举报类型 1政治 2诈骗 3侵权 4色情 5辱骂诋毁 6广告拉人 7脱离平台交易 99其他原因'"` // 举报类型
	ReportContent string    `gorm:"size:255;comment:'举报内容描述'"`                                                     // 举报内容描述
	ReportPicURL  string    `gorm:"type:json;comment:'举报图片'"`                                                      // 举报图片 (假设 JSON 字符串形式存储)
	State         int8      `gorm:"not null;default:0;comment:'状态 1待审核 2已审核'"`                                     // 状态
	SrcReply      string    `gorm:"size:255;comment:'举报人反馈'"`                                                      // 举报人反馈
	DstReply      string    `gorm:"size:255;comment:'被举报人反馈'"`                                                     // 被举报人反馈
	CreateTime    time.Time `gorm:"default:null;comment:'举报时间'"`                                                   // 举报时间
	StaffName     string    `gorm:"size:60;comment:'处理人'"`                                                         // 处理人
	UpdateTime    time.Time `gorm:"default:null;comment:'处理时间'"`                                                   // 处理时间
}

func (rc ReportingCenter) TableName() string {
	return "t_reporting_center"
}
