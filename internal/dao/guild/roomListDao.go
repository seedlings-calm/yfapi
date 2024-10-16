package dao

import (
	"github.com/gin-gonic/gin"
	"yfapi/core/coreDb"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/request/guild"
	response_guild "yfapi/typedef/response/guild"
	"yfapi/util/easy"
)

type RoomListDao struct {
}

// GetRoomListPage 获取公会后台房间列表
func (r *RoomListDao) GetRoomListPage(req *guild.GuildRoomListreq, c *gin.Context) (list interface{}, count int64, err error) {
	limit := req.Size
	offset := req.Size * (req.CurrentPage - 1)
	guildId := c.GetString("guildId")
	db := coreDb.GetMasterDb().Table("t_room as tr").Joins("left join t_user as u on u.id = tr.user_id").
		Where("tr.guild_id =?", guildId).Where("tr.live_type = ?", req.LiveType)
	if req.RoomKeyword != "" { //房间名称或房间号
		db = db.Where("tr.name LIKE ? Or tr.room_no LIKE ?", "%"+req.RoomKeyword+"%", "%"+req.RoomKeyword+"%")
	}
	if req.UserKeyword != "" {
		db = db.Where("u.user_no LIKE ? Or u.nickname LIKE ?", "%"+req.UserKeyword+"%", "%"+req.UserKeyword+"%")
	}
	if req.RoomType != 0 {
		db = db.Where("tr.room_type = ?", req.RoomType)
	}
	if req.Status != 0 {
		db = db.Where("tr.status = ?", req.Status)
	}
	err = db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	roomList := []*response_guild.GuildRoomListResp{}
	var roomInfo []RoomInfo
	err = db.Select("tr.*,u.user_no,u.nickname").Order("tr.create_time desc").Limit(limit).Offset(offset).Scan(&roomInfo).Error
	if err != nil {
		return nil, 0, err
	}
	for _, v := range roomInfo {
		roomListInfo := &response_guild.GuildRoomListResp{
			Id:         v.Id,
			RoomNo:     v.RoomNo,
			Name:       v.Name,
			Notice:     v.Notice,
			RoomType:   v.RoomType,
			CoverImg:   helper.FormatImgUrl(v.CoverImg),
			UserNo:     v.UserNo,
			NickName:   v.Nickname,
			Status:     int8(v.Status),
			CreateTime: easy.LocalTime(v.CreateTime),
		}
		var daySettleUserInfo = new(model.User)
		var monthSettleUserInfo = new(model.User)
		if v.DaySettleUserId != "" {
			//查询日结算人信息
			daySettleUserInfo = service_user.GetUserBaseInfo(v.DaySettleUserId)
			roomListInfo.DaySettleUserNo = daySettleUserInfo.UserNo
			roomListInfo.DaySettleNickname = daySettleUserInfo.Nickname
		}
		if v.MonthSettleUserId != "" {
			monthSettleUserInfo = service_user.GetUserBaseInfo(v.MonthSettleUserId)
			roomListInfo.MonthSettleUserNo = monthSettleUserInfo.UserNo
			roomListInfo.MonthSettleNickname = monthSettleUserInfo.Nickname
		}

		roomList = append(roomList, roomListInfo)
	}
	return roomList, count, err
}

type RoomInfo struct {
	model.Room
	UserNo   string `json:"user_no"`
	Nickname string `json:"nickname"`
}
