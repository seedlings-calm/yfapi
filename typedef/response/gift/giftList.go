package gift

type GiftListRes struct {
	CategoryType int       `json:"categoryType"` // 礼物类目
	List         []GiftDTO `json:"list"`         // 礼物列表
	GiftVersion  string    `json:"giftVersion"`  // 礼物版本号
}

type GiftDTO struct {
	GiftId           int             `json:"giftId"`           // 礼物主键ID
	GiftCode         string          `json:"giftCode"`         // 礼物ID
	GiftName         string          `json:"giftName"`         // 礼物名称
	GiftImage        string          `json:"giftImage"`        // 礼物Icon
	GiftAmountType   int             `json:"giftAmountType"`   // 礼物币种 1钻石 2红钻
	GiftDiamond      int             `json:"giftDiamond"`      // 礼物价格
	GiftGrade        int             `json:"giftGrade"`        // 礼物等级
	CategoryType     int             `json:"categoryType"`     // 礼物类目
	AnimationUrl     string          `json:"animationUrl"`     // 礼物动效url
	AnimationJsonUrl string          `json:"animationJsonUrl"` // 礼物动效json url
	SubscriptContent string          `json:"subscriptContent"` //角标内容
	SubscriptIcon    string          `json:"subscriptIcon"`    //角标样式
	SendCountList    []GiftSendCount `json:"sendCountList"`    // 礼物赠送数量列表
}

type GiftSendCount struct {
	SendCount int    `json:"sendCount" gorm:"column:send_count"` // 赠送数量
	Desc      string `json:"desc" gorm:"column:desc"`            // 赠送数量描述
}

type SendGiftRes struct {
	DiamondAmount string `json:"diamondAmount"` // 钻石余额
	LvLevel       int    `json:"lvLevel"`       // lv等级
	LvCurrExp     int    `json:"lvCurrExp"`     // lv经验
	LvMinExp      int    `json:"lvMinExp"`      // lv最小经验
	LvMaxExp      int    `json:"lvMaxExp"`      // lv最大经验
	LvIcon        string `json:"lvIcon"`        // lv等级图标
}

type GiftSource struct {
	GiftCode         string `json:"giftCode"`                   // 礼物ID
	AnimationUrl     string `json:"animationUrl,omitempty"`     // 礼物动效url
	AnimationJsonUrl string `json:"animationJsonUrl,omitempty"` // 礼物动效json url
}

type GiftSourceListRes struct {
	List []GiftSource `json:"list"` // 礼物资源地址
}
