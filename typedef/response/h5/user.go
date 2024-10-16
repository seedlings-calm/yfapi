package h5

type UserAccountRes struct {
	UserId        string `json:"userId"` // 用户id
	UserNo        string `json:"userNo"`
	Uid32         int32  `json:"uid32"`
	Nickname      string `json:"nickname"`      //昵称
	Avatar        string `json:"avatar"`        //头像
	DiamondAmount string `json:"diamondAmount"` // 钻石余额
	MyGoodsIsRed  bool   `json:"myGoodsIsRed"`  //我的装扮中心， 是否有红点
	MyGoodsIsHave bool   `json:"myGoodsIsHave"` //是否有装扮
}
