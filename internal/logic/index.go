package logic

import (
	"log"
	"slices"
	"sort"
	"time"
	"yfapi/app/handle"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/auth"
	service_gift "yfapi/internal/service/gift"
	service_room "yfapi/internal/service/room"
	"yfapi/typedef/enum"
	request_index "yfapi/typedef/request/index"
	response_index "yfapi/typedef/response/index"
	response_room "yfapi/typedef/response/room"
	response_user "yfapi/typedef/response/user"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

type Indexs struct{}

func (i Indexs) GetCollect(c *gin.Context) (resp response_user.CollectResponse) {
	userId := handle.GetUserId(c)
	collectDao := &dao.DaoUserCollect{
		UserId: userId,
	}
	roomIds, _ := collectDao.GetRoomIds()
	if len(roomIds) == 0 {
		resp.LiveBroadcast = nil
		resp.ChatRoom = nil
	}
	res := []model.Room{}
	roomDao := &dao.RoomDao{}
	if len(roomIds) > 0 {
		// 直播间收藏排序：优先展示开播房间、再展示未开播房间；
		rooms := roomDao.FindByIds(roomIds, "status asc")
		var chatRooms []*response_room.RoomInfo
		var liveBroadcast []*response_room.RoomInfo
		//TODO:根据热度排序，未处理,热度值为赋值
		roomLists := ForRoomsAction(rooms, false)
		for _, v := range roomLists {
			switch v.LiveType {
			case enum.LiveTypeChatroom:
				chatRooms = append(chatRooms, v)
			case enum.LiveTypeAnchor: // 展示主播头像
				liveBroadcast = append(liveBroadcast, v)
			default:
				log.Println("个人目前收藏不展示")
			}
		}

		// 聊天室收藏，按照房间热度值依次排序  TODO:
		resp.ChatRoom = chatRooms
		resp.LiveBroadcast = liveBroadcast
		// 为你推荐，取后台配置数据；数据为空  TODO: 展示平台热度前四名的聊天室（情感、交友、唱歌、播客）
		roomRecDao := &dao.RoomRecommendDao{}
		nowDate := time.Now().Format(time.DateOnly)
		nowTime := time.Now().Format(time.TimeOnly)
		res = roomRecDao.GetRooms(nowDate, nowTime, 0)
	}
	if len(res) > 0 {
		roomLists := ForRoomsAction(res, false)
		resp.Recommend = append(resp.Recommend, roomLists...)
	} else {
		redisArr := []redis.Z{}
		recommendRoom, _ := service_room.GetRoomHotsList(c, enum.LiveTypeChatroom, 3)
		redisArr = append(redisArr, recommendRoom...)
		anchors, _ := service_room.GetRoomHotsList(c, enum.LiveTypeAnchor, 3)
		redisArr = append(redisArr, anchors...)
		sort.Slice(redisArr, func(i, j int) bool {
			return redisArr[i].Score > redisArr[j].Score
		})
		if len(redisArr) > 4 {
			redisArr = redisArr[:4]
		}
		recommendRoomIds := []string{}
		for _, v := range redisArr {
			recommendRoomIds = append(recommendRoomIds, cast.ToString(v.Member))
		}
		recommendRooms := roomDao.FindByIds(recommendRoomIds, "status asc")
		recommendRoomLists := ForRoomsAction(recommendRooms, false)
		resp.Recommend = recommendRoomLists
	}
	//对收藏房间排序，直播中在前
	if len(resp.LiveBroadcast) > 0 {
		resp.LiveBroadcast = helper.MoveRoomToEnd(resp.LiveBroadcast)
	}
	return
}

func (i Indexs) GetRecommend(c *gin.Context) (resp response_user.RecommendResponse) {

	nowDate := time.Now().Format(time.DateOnly)
	nowTime := time.Now().Format(time.TimeOnly)

	roomDao := &dao.RoomRecommendDao{}
	chatroom := roomDao.GetRooms(nowDate, nowTime, enum.LiveTypeChatroom)
	if chatroom != nil {
		resp.ChatRoom = ForRoomsAction(chatroom, false)
		//过滤主持位无人的聊天室
		if len(resp.ChatRoom) > 0 {
			// 使用 slices.DeleteFunc 移除 主持位没有人的元素
			resp.ChatRoom = slices.DeleteFunc(resp.ChatRoom, func(p *response_room.RoomInfo) bool {
				return p.SeatList[0].UserInfo.UserId == ""
			})
		}
		if len(resp.ChatRoom) != 3 {
			resp.ChatRoom = nil
		}
	}
	anchor := roomDao.GetRooms(nowDate, nowTime, enum.LiveTypeAnchor)
	if anchor != nil {
		resp.LiveBroadcast = ForRoomsAction(anchor, true)
	}
	return
}

func (i Indexs) GetRoomsByPC(c *gin.Context) (resp []*response_room.RoomInfo) {
	userId := handle.GetUserId(c)

	auth := auth.Auth{}
	//本房间从业者，会长，房主
	roomIds := auth.GetUsersRules(userId, enum.PresidentRoleId, enum.HouseOwnerRoleId, enum.MusicianRoleId, enum.CounselorRoleId, enum.AnchorRoleId)
	if len(roomIds) == 0 {
		return
	}
	roomDao := &dao.RoomDao{}
	roomList, err := roomDao.GetRoomByIds(roomIds)
	if err != nil {
		return
	}
	resp = ForRoomsAction(roomList, false)
	return
}

func (i Indexs) GetTopMsg(c *gin.Context, req *request_index.TopMsgReq) (res []response_index.TopMsgRes) {
	return service_gift.GetTopMsg(int64(req.Page), int64(req.Size))
}
