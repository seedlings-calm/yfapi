package model

import "time"

// 黑名单
type UserBlacklist struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;comment:'主键编码'" json:"id"`
	FromID      string    `gorm:"type:varchar(100);default:NULL;comment:'操作用户ID'" json:"from_id"`
	ToID        string    `gorm:"type:varchar(100);default:NULL;comment:'被封用户ID'" json:"to_id"`
	RoomID      string    `gorm:"type:varchar(100);default:NULL;comment:'房间ID'" json:"room_id"`
	UnsealID    string    `gorm:"type:varchar(100);default:NULL;comment:'解封用户操作ID'" json:"unseal_id"`
	Types       int       `gorm:"default:1;comment:'1直播间黑名单,2用户黑名单'" json:"types"`
	IsEffective bool      `gorm:"type:tinyint(1);default:0;comment:'是否有效 =1 拉黑中'" json:"is_effective"`
	CreateTime  time.Time `gorm:"type:datetime(3);default:NULL;comment:'拉黑时间'" json:"create_time"`
	UpdateTime  time.Time `gorm:"type:datetime(3);default:NULL;comment:'更新时间'" json:"update_time"`
}

func (u UserBlacklist) TableName() string {
	return "t_user_blacklist"
}
