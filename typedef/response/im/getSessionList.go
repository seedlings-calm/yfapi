package response_im

import "yfapi/typedef/response/user"

type GetSessionListRes struct {
	UserInfo    user.SessionListUserInfo `json:"userInfo"`    //用户信息
	Timestamp   int64                    `json:"timestamp"`   //时间戳
	TextColor   string                   `json:"textColor"`   //文本颜色
	ShowContent string                   `json:"showContent"` //展示的文本信息
	NotReadNum  int                      `json:"notReadNum"`  //未读消息数
	IsTop       bool                     `json:"isTop"`
	Types       int                      `json:"types"` //0私聊 1系统通知 2官方公告 3互动消息
}
