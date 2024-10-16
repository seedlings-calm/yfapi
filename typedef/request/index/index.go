package request_index

type TopListReq struct {
	Page int `json:"page" form:"page"`                           //页码
	Size int `json:"size" form:"size" validate:"required,min=1"` //条数
}

// SearchAllReq 搜索用户、聊天室、直播间 请求
type SearchAllReq struct {
	SearchType int    `json:"searchType" form:"searchType"`               // 1用户 2聊天室 3直播间
	Keyword    string `json:"keyword" form:"keyword" validate:"required"` // id或昵称
}

type AppMenuSettingReq struct {
	ModuleType int `json:"moduleType" form:"moduleType"` // 模块类型 1个人主页菜单
}

// 头条中心请求参数
type TopMsgReq struct {
	Page int `json:"page" form:"page"`                           //页码
	Size int `json:"size" form:"size" validate:"required,min=1"` //条数
}
