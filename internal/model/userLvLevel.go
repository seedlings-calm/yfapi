package model

import "time"

// 用户lv等级表
type UserLvLevel struct {
	ID         int       `json:"id" gorm:"column:id"`
	UserId     string    `json:"userId" gorm:"column:user_id"`
	Level      int       `json:"level" gorm:"column:level"`
	CurrExp    int       `json:"currExp" gorm:"column:curr_exp"`
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *UserLvLevel) TableName() string {
	return "t_user_lv_level"
}

//1-10级 #47DDD3 透明度80%
//11-20级 #5D9CFF 透明度80%
//21-30级 #875DF8 透明度80%
//31-40级 #F2814A 透明度80%
//41-50级 #E1C261 透明度80%
//51-60级 #37C157 透明度80%
//61-70级 #5C5AFF 透明度80%
//71-80级 #E47AFF-#8F52F0 透明度80%
//81-90级 #F85959-#F69057 透明度80%
//91-100级 #413225-#B5863A ；描边#95775C-#FFF5B4-#D29736 透明度80%

// GetColor lv等级进场横幅颜色
func (m *UserLvLevel) GetColor() ([]string, []string) {
	switch {
	case m.Level <= 1:
		return []string{"#47DDD3"}, nil
	case m.Level <= 2:
		return []string{"#5D9CFF"}, nil
	case m.Level <= 3:
		return []string{"#875DF8"}, nil
	case m.Level <= 4:
		return []string{"#F2814A"}, nil
	case m.Level <= 5:
		return []string{"#E1C261"}, nil
	case m.Level <= 6:
		return []string{"#37C157"}, nil
	case m.Level <= 7:
		return []string{"#5C5AFF"}, nil
	case m.Level <= 8:
		return []string{"#E47AFF", "#8F52F0"}, nil
	case m.Level <= 9:
		return []string{"#F85959", "#F69057"}, nil
	case m.Level <= 10:
		return []string{"#413225", "#B5863A"}, []string{"#95775C", "#FFF5B4", "#D29736"}
	}
	return nil, nil
}

type UserLvLevelDTO struct {
	ID            int       `json:"id" gorm:"column:id"`
	UserId        string    `json:"userId" gorm:"column:user_id"`
	Level         int       `json:"level" gorm:"column:level"`
	CurrExp       int       `json:"currExp" gorm:"column:curr_exp"`
	MinExperience int       `json:"minExperience" gorm:"column:min_experience"` // 最小经验
	MaxExperience int       `json:"maxExperience" gorm:"column:max_experience"` // 最高经验
	Icon          string    `json:"icon" gorm:"icon"`                           // 图标
	CreateTime    time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime    time.Time `json:"updateTime" gorm:"column:update_time"`
}
