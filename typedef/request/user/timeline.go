package user

type TimelinePublishReq struct {
	ContentType int     `json:"contentType" validate:"required"` // 动态类型 1图文 2视频
	TextContent string  `json:"textContent"`
	ImgList     string  `json:"imgList"`
	VideoDTO    string  `json:"videoDTO"`
	CityName    string  `json:"cityName"`    // 城市名称
	AddressName string  `json:"addressName"` // 地址名称
	Latitude    float64 `json:"latitude"`    // 纬度
	Longitude   float64 `json:"longitude"`   // 经度
}

type TimelineListReq struct {
	TargetUserId string `json:"targetUserId" form:"targetUserId" validate:"required"`
	Page         int    `json:"page" form:"page"`
	Size         int    `json:"size" form:"size" validate:"required"`
}

type TimelineReplyReq struct {
	TimelineId   int64  `json:"timelineId"`
	ToReplyId    int64  `json:"toReplyId"`
	ReplyContent string `json:"replyContent" validate:"required"`
}

type TimelineReplyListReq struct {
	TimelineId int64 `json:"timelineId" form:"timelineId" validate:"required"`
	Page       int   `json:"page" form:"page"`
	Size       int   `json:"size" form:"size" validate:"required"`
}

type TimelineSubReplyListReq struct {
	ReplyId int64 `json:"replyId" form:"replyId" validate:"required"`
	Page    int   `json:"page" form:"page"`
	Size    int   `json:"size" form:"size" validate:"required"`
}

type TimelineListByTypeReq struct {
	CategoryId int `json:"categoryId" form:"categoryId" validate:"required"`
	Page       int `json:"page" form:"page"`
	Size       int `json:"size" form:"size" validate:"required"`
}
