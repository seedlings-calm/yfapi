package enum

import "time"

const (
	UserConnectImServiceLife = time.Hour * 24
)

const (
	//互动消息点赞
	InteractiveMsgLikeTypes = 1
	//互动评论消息
	InteractiveCommentTypes = 2
	//回复通知
	InteractiveReplyTypes = 3
)
