package model

type DicUserNo struct {
	Id     int64  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT" ` //
	UserNo string `json:"userNo" gorm:"column:user_no"`                    //生成的用户id
	Status int    `json:"status" gorm:"column:status"`                     //是否使用
}

func (data *DicUserNo) TableName() string {
	return "t_dic_user_no"
}
