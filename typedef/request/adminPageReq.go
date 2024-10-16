package request

import "gorm.io/gorm"

// PageInfo Paging common input parameter structure
type PageInfo struct {
	CurrentPage int `json:"pageNum" form:"pageNum"`   //当前页码
	Size        int `json:"pageSize" form:"pageSize"` //分页条数
}

func (r *PageInfo) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.CurrentPage <= 0 {
			r.CurrentPage = 1
		}
		switch {
		case r.Size > 100:
			r.Size = 100
		case r.Size <= 0:
			r.Size = 10
		}
		offset := (r.CurrentPage - 1) * r.Size
		return db.Offset(offset).Limit(r.Size)
	}
}
