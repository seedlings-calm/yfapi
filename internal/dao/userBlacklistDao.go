package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
	response_room "yfapi/typedef/response/room"
)

type UserBlackListDao struct {
}

// 判断是否拉黑，true 拉黑中
func (u UserBlackListDao) IsLog(param *model.UserBlacklist) bool {
	var count int64
	coreDb.GetMasterDb().Model(&model.UserBlacklist{}).Where(param).Count(&count)

	return count >= 1
}

func (u UserBlackListDao) Create(param *model.UserBlacklist) error {
	models := &model.UserBlacklist{
		ToID:        param.ToID,
		RoomID:      param.RoomID,
		Types:       param.Types,
		IsEffective: false,
	}
	if param.Types == 2 {
		models.FromID = param.FromID
	}
	err := coreDb.GetMasterDb().Model(model.UserBlacklist{}).
		Where(models).Assign(&model.UserBlacklist{
		FromID:      param.FromID,
		IsEffective: param.IsEffective,
		CreateTime:  time.Now(),
		UpdateTime:  time.Now(),
	}).FirstOrCreate(&model.UserBlacklist{}).Error
	return err
}

func (u UserBlackListDao) Update(param *model.UserBlacklist) (err error) {
	models := &model.UserBlacklist{ToID: param.ToID, RoomID: param.RoomID, Types: param.Types, IsEffective: true}
	if param.Types == 2 {
		models.FromID = param.FromID
	}
	err = coreDb.GetMasterDb().Model(model.UserBlacklist{}).
		Where(models).
		Updates(map[string]interface{}{
			"unseal_id":    param.UnsealID,
			"is_effective": false,
			"update_time":  time.Now(),
		}).Error
	return
}

func (u UserBlackListDao) GetListByRoomId(roomId string) (res []response_room.BlackListAndUserInfo) {
	coreDb.GetMasterDb().Table("t_user_blacklist as a").
		Joins("inner join t_user as b on a.to_id = b.id").
		Where("a.room_id = ? and a.types = 1 and a.is_effective = 1", roomId).
		Select("a.to_id as user_id,b.user_no,b.avatar,b.nickname, b.sex").Order("a.update_time").Scan(&res)
	return
}

// GetUserBlackList 获取用户黑名单列表
func (u UserBlackListDao) GetUserBlackList(userId string) (res []*response_room.BlackListAndUserInfo) {
	coreDb.GetMasterDb().Table("t_user_blacklist as a").
		Joins("inner join t_user as b on a.to_id = b.id").
		Where("from_id=? and a.types =? and a.is_effective = 1", userId, enum.BlacklistTypeUser).
		Select("a.to_id user_id,b.user_no,b.avatar,b.nickname, b.sex, b.introduce").Order("a.update_time").Scan(&res)
	return
}

// 获取黑名单数量
func (u UserBlackListDao) GetCount(userId string, types int) (count int64) {
	coreDb.GetMasterDb().Model(&model.UserBlacklist{}).Where("from_id = ?  and types = ? and is_effective = 1", userId, types).Count(&count)
	return
}
