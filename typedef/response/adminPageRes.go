package response

type AdminPageRes struct {
	Data        interface{} `json:"data"`        //对应数据
	Total       int64       `json:"total"`       //总条数
	CurrentPage int         `json:"currentPage"` //当前页码
	Size        int         `json:"size"`        //分页条数
}
