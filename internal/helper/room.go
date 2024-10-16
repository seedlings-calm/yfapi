package helper

import (
	"yfapi/typedef/enum"
	"yfapi/typedef/response/room"
)

// 将未开播房间移动到切片最后
func MoveRoomToEnd(slice []*room.RoomInfo) []*room.RoomInfo {
	newSlice := []*room.RoomInfo{}
	toEnd := map[int]*room.RoomInfo{}
	for i, v := range slice {
		if v.ShowLabel == enum.RoomShowStatusClose {
			toEnd[i] = v
		} else {
			newSlice = append(newSlice, v)
		}
	}

	if len(toEnd) > 0 {
		for _, v := range toEnd {
			newSlice = append(newSlice, v)
		}
	}
	return newSlice
}
