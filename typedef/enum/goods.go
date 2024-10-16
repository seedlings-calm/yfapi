package enum

// 物品固定分类，key不可修改
var GoodsType = map[string]string{
	"DJ":   "道具",
	"JF":   "积分",
	"JY":   "经验",
	"SL":   "声浪",
	"TS":   "头饰",
	"ZJ":   "座驾",
	"ERS":  "进场特效",
	"MWK":  "麦位框",
	"RCMB": "聊天室气泡",
}
var (
	GoodsTypeDJ   = "DJ"
	GoodsTypeJF   = "JF"
	GoodsTypeJY   = "JY"
	GoodsTypeTS   = "TS"
	GoodsTypeSL   = "SL"
	GoodsTypeZJ   = "ZJ"
	GoodsTypeERS  = "ERS"
	GoodsTypeMWK  = "MWK"
	GoodsTypeRCMB = "RCMB"
)

// goods_grant 发放表  方法来源枚举
var (
	// GoodsGrantSourceAdmin  后台发放
	GoodsGrantSourceAdmin = 1
	// GoodsGrantSourceVIP vip等级发放
	GoodsGrantSourceVIP = 2
	// GoodsGrantSourceLV  lv等级发放
	GoodsGrantSourceLV = 3
	// GoodsGrantSourceStar  星光等级发放
	GoodsGrantSourceStar = 4
)
