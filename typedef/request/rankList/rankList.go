package rankList

type RankListReq struct {
	Types  string `json:"types" form:"types" validate:"required,oneof=contributor popularity"` //contributor贡献榜 popularity人气榜
	Range  string `json:"range" form:"range" validate:"required,oneof=day week month"`         //day 日榜 week 周榜单 month 月榜
	RoomId string `json:"roomId" form:"roomId"`
}
