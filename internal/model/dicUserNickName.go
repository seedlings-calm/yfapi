package model

type DicUserNickName struct {
	Id           int64  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT" ` //
	UserNickName string `json:"userNickName" gorm:"column:user_nick_name"`       //生成的用户昵称
	Status       int    `json:"status" gorm:"column:status"`                     //是否使用 0:未使用 1：已使用
}

func (data *DicUserNickName) TableName() string {
	return "t_dic_user_nick_name"
}
