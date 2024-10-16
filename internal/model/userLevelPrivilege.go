package model

import "time"

// 用户等级特权权益表
type UserLevelPrivilege struct {
	ID            int       `json:"id" gorm:"column:id"`
	Name          string    `json:"name" gorm:"column:name"`                    // 用户lv等级id
	Icon          string    `json:"icon" gorm:"column:icon"`                    // 图标
	LightEffect   string    `json:"lightEffect" gorm:"column:light_effect"`     // 点亮效果
	MinLv         int       `json:"minLv"        gorm:"column:min_lv"`          // lv解锁等级 -1未配置
	MinVip        int       `json:"minVip"       gorm:"column:min_vip"`         // vip解锁等级 -1未配置
	MinStar       int       `json:"minStar"      gorm:"column:min_star"`        // 星光解锁等级 -1未配置
	ColorList     string    `json:"colorList" gorm:"column:color_list"`         // 颜色列表设置
	Explain       string    `json:"explain" gorm:"column:explain"`              // 说明
	FirstOperator string    `json:"firstOperator" gorm:"column:first_operator"` // 创建人
	LastOperator  string    `json:"lastOperator" gorm:"column:last_operator"`   // 最新操作人
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *UserLevelPrivilege) TableName() string {
	return "t_user_level_privilege"
}
