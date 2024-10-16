package room

// 广告位
type AdvertisingResp struct {
	Id       int    `json:"id"`
	OpenType int    `json:"openType"` //打开方式 1房间内 2房间外
	Ratio    int    `json:"ratio"`    //打开比例 1-100
	Image    string `json:"image"`    //展示图片
	Url      string `json:"url"`      //打开链接
	SiteType int    `json:"siteType"` //位置类型 1大挂件 2小挂件
}
