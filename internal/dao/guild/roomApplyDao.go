package dao

import (
	"github.com/gin-gonic/gin"
	"yfapi/core/coreDb"
	"yfapi/internal/helper"
	model2 "yfapi/internal/model/guild"
	service_user "yfapi/internal/service/user"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/request/guild"
	response_guild "yfapi/typedef/response/guild"
	"yfapi/util/easy"
)

type RoomApplyDao struct {
}

// GetRoomApplyList 获取公会后台房间申请列表
func (r *RoomApplyDao) GetRoomApplyList(req *guild.GuildRoomApplyListReq, c *gin.Context) (list interface{}, count int64, err error) {
	limit := req.Size
	offset := req.Size * (req.CurrentPage - 1)
	guildId := helper.GetGuildId(c)
	db := coreDb.GetMasterDb().Table("t_guild_room_apply as gra").
		Joins("left join t_room as tr on tr.id = gra.room_id").
		Joins("left join t_user as u on u.id = gra.room_user_id").
		Where("gra.guild_id =?", guildId)
	if req.RoomKeyword != "" {
		db = db.Where("(gra.room_name LIKE ? OR tr.room_no LIKE ?)", "%"+req.RoomKeyword+"%", "%"+req.RoomKeyword+"%")
	}
	if req.UserKeyWord != "" {
		db = db.Where("u.user_no LIKE ? Or u.nickname LIKE ?", "%"+req.UserKeyWord+"%", "%"+req.UserKeyWord+"%")
	}
	if req.Status != 0 {
		db = db.Where("gra.status = ?", req.Status)
	}

	err = db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	List := []*response_guild.GuildRoomApplyListResp{}
	var roomApply []RoomApply
	err = db.Select("gra.*,tr.room_no,u.user_no,u.nickname").Order("gra.create_time").Limit(limit).Offset(offset).Scan(&roomApply).Error
	if err != nil {
		return nil, 0, err
	}
	// 查询日月结算人信息
	userIdMap := make(map[string]struct{})
	for _, info := range roomApply {
		userIdMap[info.DaySettleUserID] = struct{}{}
		userIdMap[info.MonthSettleUserID] = struct{}{}
	}
	var userIdList []string
	for key := range userIdMap {
		userIdList = append(userIdList, key)
	}
	userMap := service_user.GetUserBaseInfoMap(userIdList)
	for _, v := range roomApply {
		List = append(List, &response_guild.GuildRoomApplyListResp{
			Id:                  v.ID,
			RoomNo:              v.RoomNo,
			RoomName:            v.RoomName,
			RoomDesc:            v.RoomDesc,
			RoomType:            v.RoomType,
			CoverImg:            helper.FormatImgUrl(v.RoomAvatar),
			UserNo:              v.UserNo,
			NickName:            v.Nickname,
			CreateTime:          easy.LocalTime(v.CreateTime),
			Status:              v.Status,
			DaySettleNickname:   userMap[v.DaySettleUserID].Nickname,
			MonthSettleNickname: userMap[v.MonthSettleUserID].Nickname,
		})
		if v.Status == typedef_enum.GuildRoomApplyStatusPass {
			List[len(List)-1].UpdateTime = easy.LocalTime(v.UpdateTime)
		}
		if v.Status == typedef_enum.GuildRoomApplyStatusRefuse {
			List[len(List)-1].Reason = v.Reason
		}
	}
	return List, count, err
}

type RoomApply struct {
	model2.GuildRoomApply
	RoomNo   string `json:"roomNo"`
	UserNo   string `json:"userNo"`
	Nickname string `json:"nickname"`
}

// 查询房间申请信息
func (g *RoomApplyDao) FindOne(param *model2.GuildRoomApply) (data *model2.GuildRoomApply, err error) {
	data = new(model2.GuildRoomApply)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	return
}

func (g *RoomApplyDao) UpdateApply(param *model2.GuildRoomApply) (err error) {
	err = coreDb.GetMasterDb().Where(param).Updates(param).Error
	return
}
