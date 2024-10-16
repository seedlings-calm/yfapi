package user

type AddFriendReq struct {
	FriendId string `json:"friendId" validate:"required"` // 好友ID
}
