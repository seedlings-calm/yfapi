package room

type RoomInfo struct {
	RoomShowBaseRes
	SeatList  []*RoomWheatPosition `json:"seatList"`  // 房间麦位用户头像列表(不展示老板位),主持人为第一个列
	Hot       int                  `json:"hot"`       // 热度值
	HotStr    string               `json:"hotStr"`    // 热度值显示
	ShowLabel int                  `json:"showLabel"` // 房间标签
	//TODO:房间标签  表示房间的当时状态，pk中，直播中，等
}

// ApplyAnchorRoomRes 申请直播间返回
type ApplyAnchorRoomRes struct {
	RoomId string `json:"roomId"` // 房间ID
}

// 申请房间房间信息回显
type ApplyRoomInfoRes struct {
	RoomId        string `json:"roomId"`        // 房间ID
	RoomNo        string `json:"roomNo"`        // 房间No
	RoomName      string `json:"roomName"`      // 房间名称
	CoverImg      string `json:"coverImg"`      // 房间封面
	LiveType      int    `json:"liveType"`      // 房间直播类型
	RoomType      int    `json:"roomType"`      // 房间类型
	Notice        string `json:"notice"`        // 房间公告
	BackgroundImg string `json:"backgroundImg"` // 背景图
}
