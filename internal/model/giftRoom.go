package model

import "time"

// GiftRoom
// @Description: 房间礼物
type GiftRoom struct {
	ID                 int       `json:"id" gorm:"column:id"`
	GiftCode           string    `json:"gift_code" gorm:"column:gift_code"`                       // 礼物编码ID
	CategoryType       int       `json:"category_type" gorm:"column:category_type"`               // 礼物类目
	RoomLiveType       int       `json:"room_live_type" gorm:"column:room_live_type"`             // 分类 1聊天室，2个播
	RoomType           string    `json:"room_type" gorm:"column:room_type"`                       // 展示厅类型，逗号分割
	Status             int       `json:"status" gorm:"column:status"`                             // 上下架，1上架, 2下架
	DelStatus          int       `json:"del_status" gorm:"column:del_status"`                     // 删除状态，留存记录，1删除，0正常
	SortNo             int       `json:"sort_no" gorm:"column:sort_no"`                           // 礼物序号,越小越靠前
	StartTime          time.Time `json:"start_time" gorm:"column:start_time"`                     // 展示开始时间
	EndTime            time.Time `json:"end_time" gorm:"column:end_time"`                         // 展示结时间
	SubscriptContent   string    `json:"subscript_content" gorm:"column:subscript_content"`       // 角标内容
	SubscriptIcon      string    `json:"subscript_icon" gorm:"column:subscript_icon"`             // 角标样式
	SubscriptStartTime time.Time `json:"subscript_start_time" gorm:"column:subscript_start_time"` // 角标展示开始时间
	SubscriptEndTime   time.Time `json:"subscript_end_time" gorm:"column:subscript_end_time"`     // 角标展示结束时间
	SubscriptStatus    int       `json:"subscript_status" gorm:"column:subscript_status"`         // 角标状态，1开启 2关闭
	IsFreeGift         int       `json:"is_free_gift" gorm:"column:is_free_gift"`                 // 是否为免费礼物 0否 1是
	FreeTimeList       string    `json:"free_time_list" gorm:"column:free_time_list"`             // 免费礼物倒计时列表
	LimitLevelType     int       `json:"limit_level_type" gorm:"column:limit_level_type"`         // 使用等级类型 0无限制 1vip 2lv 3star
	LimitLevel         int       `json:"limit_level" gorm:"column:limit_level"`                   // 使用等级
	StaffName          string    `json:"staff_name" gorm:"column:staff_name"`                     // 操作人昵称
	CreateTime         time.Time `json:"create_time" gorm:"column:create_time"`                   // 创建时间
	UpdateTime         time.Time `json:"update_time" gorm:"column:update_time"`                   // 更新时间
}

type GiftDTO struct {
	ID                 int       `json:"id" gorm:"column:id"`
	GiftCode           string    `json:"gift_code" gorm:"column:gift_code"`                       // 礼物编码ID
	CategoryType       int       `json:"category_type" gorm:"column:category_type"`               // 礼物类目
	RoomLiveType       int       `json:"room_live_type" gorm:"column:room_live_type"`             // 分类 1聊天室，2个播
	RoomType           string    `json:"room_type" gorm:"column:room_type"`                       // 展示厅类型，逗号分割
	Status             int       `json:"status" gorm:"column:status"`                             // 上下架，1上架, 2下架
	DelStatus          int       `json:"del_status" gorm:"column:del_status"`                     // 删除状态，留存记录，1删除，0正常
	SortNo             int       `json:"sort_no" gorm:"column:sort_no"`                           // 礼物序号,越小越靠前
	StartTime          time.Time `json:"start_time" gorm:"column:start_time"`                     // 展示开始时间
	EndTime            time.Time `json:"end_time" gorm:"column:end_time"`                         // 展示结时间
	SubscriptContent   string    `json:"subscript_content" gorm:"column:subscript_content"`       // 角标内容
	SubscriptIcon      string    `json:"subscript_icon" gorm:"column:subscript_icon"`             // 角标样式
	SubscriptStartTime time.Time `json:"subscript_start_time" gorm:"column:subscript_start_time"` // 角标展示开始时间
	SubscriptEndTime   time.Time `json:"subscript_end_time" gorm:"column:subscript_end_time"`     // 角标展示结束时间
	SubscriptStatus    int       `json:"subscript_status" gorm:"column:subscript_status"`         // 角标状态，1开启 2关闭
	IsFreeGift         int       `json:"is_free_gift" gorm:"column:is_free_gift"`                 // 是否为免费礼物 0否 1是
	FreeTimeList       string    `json:"free_time_list" gorm:"column:free_time_list"`             // 免费礼物倒计时列表
	LimitLevelType     int       `json:"limit_level_type" gorm:"column:limit_level_type"`         // 使用等级类型 0无限制 1vip 2lv 3star
	LimitLevel         int       `json:"limit_level" gorm:"column:limit_level"`                   // 使用等级
	StaffName          string    `json:"staff_name" gorm:"column:staff_name"`                     // 操作人昵称
	CreateTime         time.Time `json:"create_time" gorm:"column:create_time"`                   // 创建时间
	UpdateTime         time.Time `json:"update_time" gorm:"column:update_time"`                   // 更新时间

	GiftName         string  `json:"gift_name" gorm:"column:gift_name"`                   // 礼物名称
	GiftImage        string  `json:"gift_image" gorm:"column:gift_image"`                 // 礼物图片
	GiftGrade        int     `json:"gift_grade" gorm:"column:gift_grade"`                 // 礼物等级
	AnimationUrl     string  `json:"animation_url" gorm:"column:animation_url"`           // VAP配置地址
	AnimationJsonUrl string  `json:"animation_json_url" gorm:"column:animation_json_url"` // VAP JSON配置地址
	GiftAmountType   int     `json:"gift_amount_type" gorm:"column:gift_amount_type"`     // 礼物币种 1钻石 2红钻
	GiftDiamond      int     `json:"gift_diamond" gorm:"column:gift_diamond"`             // 礼物价格
	GiftRevenueType  int     `json:"gift_revenue_type" gorm:"column:gift_revenue_type"`   // 礼物收益类型 2红钻 3星光
	ExpTimes         float64 `json:"exp_times" gorm:"column:exp_times"`                   // 经验倍数 默认1倍钻石
	SendCountList    string  `json:"send_count_list" gorm:"column:send_count_list"`       // 赠送数量列表 关联gift_send_count id
}

func (m *GiftRoom) TableName() string {
	return "t_gift_room"
}
