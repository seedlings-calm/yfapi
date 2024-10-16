package user

import (
	"yfapi/typedef/response"
)

type TimelineInfo struct {
	TimelineId      int64                   `json:"timelineId"`  // 动态表主键
	ContentType     int                     `json:"contentType"` // 1=图片动态 2=视频动态
	UserId          string                  `json:"userId"`      // 用户Id
	TextContent     string                  `json:"textContent"` // 动态内容
	Status          int                     `json:"status"`      // 状态 1=正常 0=删除
	LoveCount       int                     `json:"loveCount"`   // 点赞量
	Latitude        float64                 `json:"latitude"`    // 纬度
	Longitude       float64                 `json:"longitude"`   // 经度
	CityName        string                  `json:"cityName"`    // 城市名
	AddressName     string                  `json:"addressName"` // 地址名称
	ReplyCount      int                     `json:"replyCount"`  // 评论数
	ImgDTOList      []TimelineImgDTO        `json:"imgDTOList,omitempty"`
	VideoDTO        *TimelineVideoDTO       `json:"videoDTO,omitempty"`
	IsTop           bool                    `json:"isTop"`      // 是否置顶
	IsPraised       bool                    `json:"isPraised"`  //是否点赞过
	IsFollow        bool                    `json:"isFollow"`   //是否关注
	CreateTime      string                  `json:"createTime"` // 创建时间
	TimelineTimeStr string                  `json:"createTimeStr"`
	UpdateTime      string                  `json:"updateTime"`          // 更新时间
	ReplyList       []*TimelineReplyInfo    `json:"replyList,omitempty"` //评论列表
	UserPlaque      response.UserPlaqueInfo `json:"userPlaque"`          // 用户铭牌信息

	UserNo             string `json:"userNo"`
	Uid32              int32  `json:"uid32"`
	Nickname           string `json:"nickname"`
	Avatar             string `json:"avatar"`
	Sex                int    `json:"sex"`
	UserLastActiveTime string `json:"userLastActiveTime"` //用户在线时间
}

// 动态的图片对象
type TimelineImgDTO struct {
	ImgPhotoKey string `json:"imgPhotoKey,omitempty"`
	ImgUrl      string `json:"imgUrl"`
	Width       string `json:"width"`
	Height      string `json:"height"`
}

// 动态的视频对象
type TimelineVideoDTO struct {
	VideoKey         string `json:"videoKey,omitempty"`
	Width            string `json:"width"`
	Height           string `json:"height"`
	Duration         string `json:"duration"`
	VideoUrl         string `json:"videoUrl"`         // 视频原始地址
	VideoWaterUrl    string `json:"videoWaterUrl"`    // 水印视频地址
	VideoCoverImgUrl string `json:"videoCoverImgUrl"` // 视频封面图
}

type TimelineReplyInfo struct {
	ReplyId       int64                   //评论ID
	TimelineId    int64                   //动态ID
	ReplierId     string                  //评论人ID
	ReplierName   string                  //评论人昵称
	ReplierAvatar string                  //评论人头像
	ReplyContent  string                  //评论内容
	ToReplyId     int64                   // 被评论的帖子ID
	ToReplierId   string                  // 被评论的用户Id
	ToReplierName string                  //被评论的用户昵称
	SubReplyCount int                     // 子评论数
	ToSubReplyId  int64                   //实际评论得子评论ID
	IsPraised     bool                    // 是否赞过
	PraisedCount  int                     //赞的次数
	UserPlaque    response.UserPlaqueInfo `json:"userPlaque"` // 用户铭牌信息

	CreateTime   string
	SubReplyList []*TimelineReplyInfo `json:"subReplyList,omitempty"` //子评论列表

	CreateTimeStr string
}

// PraisedUserInfo 点赞用户信息
type PraisedUserInfo struct {
	UserId     string                  `json:"userId"`     // 用户id
	UserNo     string                  `json:"userNo"`     // 用户no
	Nickname   string                  `json:"nickname"`   // 昵称
	Avatar     string                  `json:"avatar"`     // 头像
	Sex        int                     `json:"sex"`        // 性别
	UserPlaque response.UserPlaqueInfo `json:"userPlaque"` // 用户铭牌信息
}
