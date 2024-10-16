package gift

type GiftListReq struct {
	LiveType     int    `json:"liveType" form:"liveType" validate:"required"`         // 房间直播类型
	RoomType     int    `json:"roomType" form:"roomType" validate:"required"`         // 房间类型
	CategoryType int    `json:"categoryType" form:"categoryType" validate:"required"` // 展示类目类型
	GiftVersion  string `json:"giftVersion" form:"giftVersion"`                       // 选中展示类目的礼物列表版本号
}

type SendGiftReq struct {
	RoomId       string   `json:"roomId" validate:"required"`       // 房间ID
	ToUserIdList []string `json:"toUserIdList" validate:"required"` // 被打赏人列表
	GiftCode     string   `json:"giftCode" validate:"required"`     // 礼物ID
	GiftCount    int      `json:"giftCount" validate:"required"`    // 赠送的礼物数量
	GiftDiamond  int      `json:"giftDiamond" validate:"required"`  // 礼物价格
	IsBatch      bool     `json:"isBatch"`                          // 是否全麦
}
