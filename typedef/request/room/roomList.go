package room

type RoomListReq struct {
	Page     int `json:"page" form:"page"`                             //页码
	Size     int `json:"size" form:"size" validate:"required,min=1"`   //条数
	RoomType int `json:"roomType" form:"roomType" validate:"required"` //房间类型
}

type RoomUpdateReq struct {
	RoomId        string `json:"roomId" form:"roomId" validate:"required"`                        //房间ID
	CoverImg      string `json:"coverImg" form:"coverImg" gorm:"column:cover_img"`                // 封面图
	Notice        string `json:"notice" form:"notice" gorm:"column:notice" `                      // 房间公告
	Name          string `json:"name" form:"name" gorm:"column:name" `                            // 房间名称
	BackgroundImg string `json:"backgroundImg" form:"backgroundImg" gorm:"column:background_img"` // 房间背景
}

type RoomLockReq struct {
	RoomId string `json:"roomId" form:"roomId" validate:"required"` //房间ID
	Pwd    string `json:"pwd" form:"pwd"`                           //锁定设置密码
}

// 踢出用户逇传参
type KickOutReq struct {
	RoomId string `json:"roomId" form:"roomId"` //房间ID
	UserId string `json:"userId" form:"userId"` //用户ID
	Times  string `json:"times" form:"times"`   //踢出房间,单位（分钟）
}

type ReportingCenterReq struct {
	DstId       string   `json:"dstId" form:"dstId" validate:"required"`                          //被举报ID
	Object      int      `json:"objectType" form:"objectType" validate:"required,min=1,max=2"`    //1房间，2用户
	Scene       int      `json:"scene" form:"scene" validate:"required,min=1"`                    //举报场景 举报场景 1房间 2消息私聊 3个人主页 4个人动态 5动态评价 6声音派对
	ReportTypes int      `json:"reportTypes" form:"reportTypes" validate:"required,min=1,max=99"` //举报类型 1政治 2诈骗 3侵权 4色情 5辱骂诋毁 6广告拉人 7脱离平台交易 99其他原因
	Content     string   `json:"content" form:"content"`                                          //举报内容
	Pics        []string `json:"pics" form:"pics" validate:"required"`                            //举报图片 值为 json字符串
}

// ApplyAnchorRoomReq 申请直播间请求
type ApplyAnchorRoomReq struct {
	RoomName      string `json:"roomName" validate:"required"` // 房间名称
	CoverImg      string `json:"coverImg" validate:"required"` // 房间封面图
	RoomType      int    `json:"roomType" validate:"required"` // 房间类型 201语音直播 202视频直播 301个人房
	RoomNotice    string `json:"roomNotice"`                   // 房间介绍
	BackgroundImg string `json:"backgroundImg"`                // 背景图
}
