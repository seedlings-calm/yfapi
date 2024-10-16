package room

type RoomAdminAddReq struct {
	RoomId string `json:"roomId" form:"roomId"` //房间ID
	UserId string `json:"userId" form:"userId"` //用户ID
}
type RoomAdminDeleteReq struct {
	RoomId string `json:"roomId" form:"roomId"`
	UserId string `json:"userId" form:"userId"`
}
type RoomAdminListReq struct {
	RoomId string `json:"roomId" form:"roomId"`
}
