package model

import "time"

type Timeline struct {
	Id               int64     `json:"id" gorm:"column:id"`                                // 动态表主键
	ContentType      int       `json:"contentType" gorm:"column:content_type"`             // 1=图片动态 2=视频动态
	UserId           string    `json:"userId" gorm:"column:user_id"`                       // 用户Id
	TextContent      string    `json:"textContent" gorm:"column:text_content"`             // 动态内容
	VideoOriginUrl   string    `json:"videoOriginUrl" gorm:"column:video_origin_url"`      // 原始视频地址
	VideoWaterUrl    string    `json:"videoWaterUrl" gorm:"column:video_water_url"`        // 水印视频地址
	VideoCoverImgUrl string    `json:"videoCoverImgUrl" gorm:"column:video_cover_img_url"` // 视频封面图
	VideoSize        string    `json:"videoSize" gorm:"column:video_size"`                 // 视频尺寸 720x1280
	Status           int       `json:"status" gorm:"column:status"`                        // 状态 1=正常 2=待审核 3审核未通过 4删除
	ShumeiExamine    string    `json:"shumeiExamine" gorm:"column:shumei_examine"`         // 数美审核结果
	LoveCount        int       `json:"loveCount" gorm:"column:love_count"`                 // 点赞量
	Latitude         float64   `json:"latitude" gorm:"column:latitude"`                    // 纬度
	Longitude        float64   `json:"longitude" gorm:"column:longitude"`                  // 经度
	CityName         string    `json:"cityName" gorm:"column:city_name"`                   // 城市名
	AddressName      string    `json:"addressName" gorm:"column:address_name"`             // 地址名称
	StaffName        string    `json:"staff_name" gorm:"column:staff_name"`                // 审核、删除操作人
	Reason           string    `json:"reason" gorm:"column:reason"`                        // 原因
	IsRecommend      int       `json:"isRecommend" gorm:"column:is_recommend"`             // 是否推荐 1= 推荐 0=不推荐
	ReplyCount       int       `json:"replyCount" gorm:"column:reply_count"`               // 评论数
	ImgList          string    `json:"imgList" gorm:"column:img_list"`                     // 图片列表
	VideoData        string    `json:"videoData" gorm:"column:video_data"`                 // 视频对象
	IsTop            bool      `json:"isTop" gorm:"column:is_top"`                         // 是否置顶 1=是 0=否
	CreateTime       time.Time `json:"createTime" gorm:"column:create_time"`               // 创建时间
	UpdateTime       time.Time `json:"updateTime" gorm:"column:update_time"`               // 更新时间
}

func (m *Timeline) TableName() string {
	return "t_timeline"
}
