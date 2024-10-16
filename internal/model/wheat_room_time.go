package model

type RoomWheatTime struct {
	ID              uint    `json:"id" redis:"-" gorm:"primaryKey;autoIncrement"`                        // 自增ID
	RoomID          string  `json:"roomId" redis:"roomId" gorm:"size:100" comment:"房间ID"`                // 房间ID
	GuildID         string  `json:"guildId" redis:"guildId" gorm:"size:100" comment:"公会ID"`              // 公会ID
	UserID          string  `json:"userId" redis:"userId" gorm:"size:100" comment:"开播用户"`                // 开播用户
	RoomType        int     `json:"roomType" redis:"roomType" gorm:"comment:'房间类型'"`                     // 房间类型
	OnwheatTime     string  `json:"onwheatTime" redis:"onwheatTime" gorm:"type:datetime" comment:"开播时间"` // 开播时间
	StatDate        string  `json:"statDate" redis:"statDate" gorm:"type:date" comment:"开播时间"`           // 开播时间天
	OnTime          int     `json:"onTime" redis:"-" comment:"开播总时间(秒)"`                                 // 开播总时间(秒)
	UpwheatTime     string  `json:"upwheatTime" redis:"-" gorm:"type:datetime" comment:"下播时间"`           // 下播时间
	EnterCount      int     `json:"enterCount" redis:"enterCount" comment:"进房总人数"`                       // 进房总人数
	EnterTimes      int     `json:"enterTimes" redis:"enterTimes" comment:"进房总人次"`                       // 进房总人次
	RewardCount     float64 `json:"rewardCount" redis:"rewardCount" comment:"打赏总金额"`                     // 打赏总金额
	RewardTimes     int     `json:"rewardTimes" redis:"rewardTimes" comment:"打赏次数"`                      // 打赏次数
	RewardUserCount int     `json:"rewardUserCount" redis:"rewardUserCount" comment:"打赏人数"`              // 打赏人数
}

// redis 存储关联接口提
type DoRoomWheatTime struct {
	ID              uint        `json:"id" redis:"-" `                                             // 自增ID
	RoomID          interface{} `json:"roomId" redis:"roomId"  `                                   // 房间ID
	GuildID         interface{} `json:"guildId" redis:"guildId"  `                                 // 公会ID
	UserID          interface{} `json:"userId" redis:"userId"  `                                   // 开播用户
	RoomType        interface{} `json:"roomType" redis:"roomType"`                                 // 房间类型
	OnwheatTime     interface{} `json:"onwheatTime" redis:"onwheatTime"  `                         // 开播时间
	StatDate        interface{} `json:"statDate" redis:"statDate" gorm:"type:date" comment:"开播时间"` // 开播时间天
	OnTime          interface{} `json:"onTime" redis:"-" `                                         // 开播总时间(秒)
	UpwheatTime     interface{} `json:"upwheatTime" redis:"-"  `                                   // 下播时间
	EnterCount      interface{} `json:"enterCount" redis:"enterCount" `                            // 进房总人数
	EnterTimes      interface{} `json:"enterTimes" redis:"enterTimes" `                            // 进房总人次
	RewardCount     interface{} `json:"rewardCount" redis:"rewardCount" `                          // 打赏总金额
	RewardTimes     interface{} `json:"rewardTimes" redis:"rewardTimes" `                          // 打赏次数
	RewardUserCount interface{} `json:"rewardUserCount" redis:"rewardUserCount" `                  // 打赏人数
}

// 房间开播统计
func (RoomWheatTime) TableName() string {
	return "t_room_wheat_time"
}
