package user

import (
	response_room "yfapi/typedef/response/room"
)

// 收藏返回结构体
type CollectResponse struct {
	ChatRoom      []*response_room.RoomInfo `json:"chatRoom"`
	LiveBroadcast []*response_room.RoomInfo `json:"liveBroadcast"`
	Recommend     []*response_room.RoomInfo `json:"recommend"`
}

// 为你推荐 返回值
type RecommendResponse struct {
	ChatRoom      []*response_room.RoomInfo `json:"chatRoom"`
	LiveBroadcast []*response_room.RoomInfo `json:"liveBroadcast"`
}
