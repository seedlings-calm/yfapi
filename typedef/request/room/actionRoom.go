package room

type ActionRoomReq struct {
	RoomId       string `json:"roomId" validate:"required"` //房间Id
	Pwd          string `json:"pwd"`                        //房间密码
	FollowUserId string `json:"followUserId"`               //跟随进房用户ID
}

type RoomAuthMenuReq struct {
	RoomId       string `json:"roomId" form:"roomId" validate:"required"`                              //房间ID
	TargetUserId string `json:"targetUserId" form:"targetUserId" validate:"omitempty"`                 //目标用户
	Scene        string `json:"scene" form:"scene" validate:"required,oneof=mic card more hidden_mic"` //触发场景位置 mic麦位，card资料卡,more更多
	Seat         int    `json:"seat" form:"seat" validate:"omitempty"`                                 //座位
	SeatTag      string `json:"seatTag" form:"seatTag" validate:"omitempty"`                           //标记
}

type ExecCommandReq struct {
	RoomId       string `json:"roomId" validate:"required"`  // 房间Id
	Seat         int    `json:"seat"`                        // 座位号
	Command      string `json:"command" validate:"required"` // 执行的命令
	TargetUserId string `json:"targetUserId"`                // 目标用户ID
	Content      string `json:"content"`                     // 内容信息
	Count        int64  `json:"count"`                       // 数量
}

type UpSeatApplyListReq struct {
	RoomId string `json:"roomId" form:"roomId" validate:"required"` // 房间Id
}

type FreeUpSeatReq struct {
	RoomId string `json:"roomId" form:"roomId" validate:"required"` // 房间Id
	Seat   *int   `json:"seat" form:"seat" validate:"required"`     // 座位号
}

type MutLocalSeatReq struct {
	RoomId string `json:"roomId" form:"roomId" validate:"required"` // 房间Id
}

type UpSeatUserListReq struct {
	RoomId   string `json:"roomId" form:"roomId" validate:"required"`     // 房间Id
	SeatType int    `json:"seatType" form:"seatType" validate:"required"` // 麦位类型 1主持 2普通
}
