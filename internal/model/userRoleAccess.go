package model

// AuthRoleAccess 表示 `t_auth_role_access` 表
type AuthRoleAccess struct {
	UserID string `gorm:"column:user_id;type:bigint unsigned;not null;default:0;comment:'用户ID'"`
	RoleID int    `gorm:"column:role_id;type:int unsigned;not null;default:0;comment:'角色ID'"`
	RoomID string `gorm:"column:room_id;type:bigint;not null;default:0;comment:'房间ID'"`
}

// TableName returns the name of the table in the database
func (AuthRoleAccess) TableName() string {
	return "t_auth_role_access"
}
