package model

import "time"

// Gift
// @Description: 礼物
type Gift struct {
	ID               int       `json:"id" gorm:"column:id"`
	GiftCode         string    `json:"gift_code" gorm:"column:gift_code"`                   // 礼物编码ID
	GiftName         string    `json:"gift_name" gorm:"column:gift_name"`                   // 礼物名称
	GiftImage        string    `json:"gift_image" gorm:"column:gift_image"`                 // 礼物图片
	GiftGrade        int       `json:"gift_grade" gorm:"column:gift_grade"`                 // 礼物等级
	AnimationUrl     string    `json:"animation_url" gorm:"column:animation_url"`           // VAP配置地址
	AnimationJsonUrl string    `json:"animation_json_url" gorm:"column:animation_json_url"` // VAP JSON配置地址
	GiftAmountType   int       `json:"gift_amount_type" gorm:"column:gift_amount_type"`     // 礼物币种 1钻石 2红钻
	GiftDiamond      int       `json:"gift_diamond" gorm:"column:gift_diamond"`             // 礼物价格
	GiftRevenueType  int       `json:"gift_revenue_type" gorm:"column:gift_revenue_type"`   // 礼物收益类型 2红钻 3星光
	ExpTimes         float64   `json:"exp_times" gorm:"column:exp_times"`                   // 经验倍数 默认1倍钻石
	UsageType        int       `json:"usage_type" gorm:"column:usage_type"`                 // 礼物使用方式
	Status           int       `json:"status" gorm:"column:status"`                         // 礼物状态 1=正常  2=作废
	SendCountList    string    `json:"send_count_list" gorm:"column:send_count_list"`       // 赠送数量列表 关联gift_send_count id
	StaffName        string    `json:"staff_name" gorm:"column:staff_name"`                 // 操作人昵称
	CreateTime       time.Time `json:"create_time" gorm:"column:create_time"`               // 创建时间
	UpdateTime       time.Time `json:"update_time" gorm:"column:update_time"`               // 更新时间
}

func (m *Gift) TableName() string {
	return "t_gift"
}
