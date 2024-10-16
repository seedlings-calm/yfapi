package request

type BasePageReq struct {
	Page int `json:"page" form:"page"`                           //页码
	Size int `json:"size" form:"size" validate:"required,min=1"` //每页条数
}
