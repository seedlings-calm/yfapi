package h5

type LevelBaseInfoRes struct {
	Nickname string     `json:"nickname"` //用户昵称
	Avatar   string     `json:"avatar"`   //用户头像
	LvBase   LevelBase  `json:"lvBase"`   // lv等级信息
	VipBase  LevelBase  `json:"vipBase"`  // vip等级信息
	StarBase *LevelBase `json:"starBase"` // star等级信息
}

type LevelBase struct {
	LevelType int         `json:"levelType"` // 等级类型 1:lv 2:vip 3:星光
	CurrLevel LevelConfig // 当前等级信息
	NextLevel LevelConfig // 下一级信息
}

// LevelConfig
// @Description: LV等级配置信息
type LevelConfig struct {
	LevelName         string                 `json:"levelName"` //等级名称
	Level             int                    `json:"level"`     //等级
	Icon              string                 `json:"icon"`      //等级图标
	LogoIcon          string                 `json:"logoIcon"`  //等级logo图标
	CurrExp           int                    `json:"currExp"`   //当前经验值
	MinExp            int                    `json:"minExp"`    //最小经验值
	MaxExp            int                    `json:"maxExp"`    //最大经验值
	MaxLevel          int                    `json:"maxLevel"`  //最大等级
	PrivilegeCount    int                    // 当前权益数量
	PrivilegeList     []*PrivilegeConfig     // 等级权益
	PrivilegeItemList []*PrivilegeItemConfig // 等级权益物品
}

// PrivilegeConfig
// @Description: 特权权益配置信息
type PrivilegeConfig struct {
	Name        string `json:"name"`        //权益名称
	Icon        string `json:"icon"`        //权益图标
	LightEffect string `json:"lightEffect"` //点亮效果
	MinLv       int    `json:"minLv"`       // lv解锁等级
	MinVip      int    `json:"minVip"`      // vip解锁等级
	MinStar     int    `json:"minStar"`     // 星光解锁等级
	Explain     string `json:"explain"`     // 说明
}

// PrivilegeItemConfig
// @Description: 特权物品配置信息
type PrivilegeItemConfig struct {
	GoodsId          int    `json:"goodsId"`          // 物品ID
	Name             string `json:"name"`             //物品名称
	Icon             string `json:"icon"`             //物品图标
	AnimationUrl     string `json:"animationUrl"`     //物品动画
	AnimationJsonUrl string `json:"animationJsonUrl"` //json文件动效
	TypeName         string `json:"typeName"`         //特权物品类型名称
	TypeKey          string `json:"typeKey"`          //物品类型key
	Explain          string `json:"explain"`          //特权物品说明
	ExpirationDate   int    `json:"expirationDate"`   //特权物品有效期
}

type LevelConfigListRes struct {
	LvList   []LevelInfo `json:"lvList"`   // lv等级配置列表
	VipList  []LevelInfo `json:"vipList"`  // vip等级配置列表
	StarList []LevelInfo `json:"starList"` // star等级配置列表
}

type LevelInfo struct {
	LevelName string `json:"levelName"` //等级名称
	Icon      string `json:"icon"`      //等级图标
	MinExp    int    `json:"minExp"`    //最小经验值
	KeepExp   int    `json:"keepExp"`   //保级经验值
}
