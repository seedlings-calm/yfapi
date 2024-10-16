package enum

import "github.com/spf13/cast"

const (
	LiveTypeChatroom = 1 //聊天室
	LiveTypeAnchor   = 2 //个播
	LiveTypePersonal = 3 //个人
)

// 个播的模板定义
const (
	RoomTemplateOne   = 2001
	RoomTemplateTwo   = 2002
	RoomTemplateThree = 2003
	RoomTemplateFour  = 2004
	RoomTemplateSix   = 2005
	RoomTemplateNine  = 2006
)

// 连麦模板
var RoomTemplates = map[string]int{
	cast.ToString(RoomTemplateTwo):   RoomTemplateTwo,
	cast.ToString(RoomTemplateThree): RoomTemplateThree,
	cast.ToString(RoomTemplateFour):  RoomTemplateFour,
	cast.ToString(RoomTemplateSix):   RoomTemplateSix,
	cast.ToString(RoomTemplateNine):  RoomTemplateNine,
}
