package room

import "time"

type GetRoomBgsRes struct {
	TrbrId     int       `json:"trbrId"     description:"房间背景ID"`           //房间背景ID
	Name       string    `json:"name"       description:"背景名称"`             //背景名称
	Icon       string    `json:"icon"       description:"房间背景缩略图"`          //房间背景缩略图
	Backgroud  string    `json:"backgroud"  description:"背景图"`              //背景图
	Types      int       `json:"types"      description:"1:默认，2：活动"`        //1:默认，2：活动"
	IsUse      int       `json:"isUse"      description:"是否使用 1：未使用 2：使用中"` //是否使用 1：未使用 2：使用中
	CreateTime time.Time `json:"createTime" description:""`
}
