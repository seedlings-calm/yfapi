package logic

import (
	"time"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_user "yfapi/internal/service/user"
	response_room "yfapi/typedef/response/room"

	"github.com/gin-gonic/gin"
)

type UserMute struct {
}

func (u *UserMute) Mute(c *gin.Context, targetUserId, roomId string, minute int) error {
	userId := helper.GetUserId(c)
	err := new(dao.UserMuteListDao).Create(&model.UserMuteList{
		FromID:    userId,
		ToID:      targetUserId,
		RoomID:    roomId,
		UnsealID:  "0",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Minute * time.Duration(minute)),
	})
	return err
}

func (u *UserMute) UnMute(c *gin.Context, userId, roomId string) error {
	err := new(dao.UserMuteListDao).Delete(userId, roomId)
	return err
}

func (u *UserMute) IsMute(userId, roomId string) model.UserMuteList {
	return new(dao.UserMuteListDao).FindOne(userId, roomId)
}

func (u *UserMute) GetMuteList(c *gin.Context, roomId string) []response_room.UserMuteListRes {
	res := new(dao.UserMuteListDao).List(roomId)
	if len(res) > 0 {
		for k, v := range res {
			res[k].Avatar = helper.FormatImgUrl(v.Avatar)
			// 查询用户的铭牌信息
			res[k].UserPlaque = service_user.GetUserLevelPlaque(v.UserId, helper.GetClientType(c))
		}
	}
	if len(res) == 0 {
		res = make([]response_room.UserMuteListRes, 0)
	}
	return res
}

func (u *UserMute) GetMuteCount(userID int64) int64 {
	return 0
}
