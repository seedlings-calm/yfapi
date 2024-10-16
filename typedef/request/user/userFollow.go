package user

type AddFollowReq struct {
	FocusUserId string `json:"focusUserId" validate:"required"` // 被关注用户ID，用于关联被关注的用户信息
}
type RemoveFollowReq struct {
	FocusUserId string `json:"focusUserId" validate:"required"` // 取消关注用户ID，用于关联被关注的用户信息
}
type GetUserFollowingListReq struct {
	UserId string `json:"userId"`
	Page   int    `json:"page" form:"page"`                           //页码
	Size   int    `json:"size" form:"size" validate:"required,min=1"` //每页条数
}
type GetUserFriendsListReq struct {
	UserId string `json:"userId"`
	Page   int    `json:"page" form:"page"`                           //页码
	Size   int    `json:"size" form:"size" validate:"required,min=1"` //每页条数
}

type DeleteFansReq struct {
	FocusUserId string `json:"focusUserId" validate:"required"` // 取消关注用户ID，用于关联被关注的用户信息
}
