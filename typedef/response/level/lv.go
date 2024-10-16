package level

// LvListResp
// @Description: lv等级列表返回
type LvListResp struct {
	UserInfo                      LvListUserInfo              `json:"userInfo"`
	NextLvUnlockGoodsAndPrivilege NextLvGoodsAndPrivilegeInfo `json:"nextLvUnlockGoodsAndPrivilege"`
	Goods                         []LvListGoodsInfo           `json:"goods"`
	Privileges                    []LvListPrivilegeInfo       `json:"privileges"`
}

type LvListUserInfo struct {
	UserId            string `json:"userId"`            // 用户ID
	Nickname          string `json:"nickname"`          // 昵称
	Avatar            string `json:"avatar"`            // 头像
	LevelName         string `json:"levelName"`         // 等级名称
	CurrentExperience int    `json:"currentExperience"` // 当前经验值
	MaxExperience     int    `json:"maxExperience"`     // 最大经验值
	Icon              string `json:"icon"`              // 等级图标
}
type NextLvGoodsAndPrivilegeInfo struct {
	Name string `json:"name"` // 名称
	Icon string `json:"icon"` // 图标
}
type LvListGoodsInfo struct {
	Name      string `json:"name"`      // 名称
	Icon      string `json:"icon"`      // 图标
	LevelName int    `json:"levelName"` // 等级名称
}
type LvListPrivilegeInfo struct {
	Name      string `json:"name"`      // 名称
	Icon      string `json:"icon"`      // 图标
	LevelName int    `json:"levelName"` // 等级名称
}
