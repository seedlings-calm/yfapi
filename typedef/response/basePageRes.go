package response

type BasePageRes struct {
	// swagger:model
	Data        any   `json:"data"`        //对应数据
	HasNext     bool  `json:"hasNext"`     //是否有下一页
	Total       int64 `json:"total"`       //总条数
	Size        int   `json:"size"`        //分页条数
	CurrentPage int   `json:"currentPage"` //当前页码
	RoomType    int   `json:"roomType"`    //房间类型
	Banner      any   `json:"banner"`      //当前页面的banner信息
}

func (s *BasePageRes) CalcHasNext() {
	s.HasNext = (s.Total - int64((s.CurrentPage+1)*s.Size)) > 0
}
