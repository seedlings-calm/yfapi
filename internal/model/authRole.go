package model

type AuthRole struct {
	ID     int    `json:"id" gorm:"column:id"`           // 主键ID
	RoleID int    `json:"role_id" gorm:"column:role_id"` // 角色id
	Title  string `json:"title" gorm:"column:title"`     // 角色中文名
	Status int    `json:"status" gorm:"column:status"`   // 状态 1正常 0禁用
	Rules  string `json:"rules" gorm:"column:rules"`     // 角色拥有规则ID,多个则用,隔开
}

func (m *AuthRole) TableName() string {
	return "t_auth_role"
}
