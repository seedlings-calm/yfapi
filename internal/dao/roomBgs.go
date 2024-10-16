package dao

import (
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"

	"gorm.io/gorm"
)

type RoomBgsDao struct {
}

func (RoomBgsDao) GetRoomBgs(roomId string) (res model.RoomBgs, err error) {
	//清除过期的背景为未使用状态
	coreDb.GetMasterDb().Model(model.RoomBgs{}).Where("expire_time <= ?", time.Now()).Update("is_use", 1)
	err = coreDb.GetMasterDb().Model(model.RoomBgs{}).Where("room_id = ? and expire_time >= ? and is_use = 2", roomId, time.Now()).Find(&res).Error
	return
}

// 更换房间背景
func (RoomBgsDao) UpdateRoomBgs(tx *gorm.DB, roomInfo *model.Room, bgInfo *model.RoomBgsResource) (err error) {

	//先取消之前使用的，然后更改当前使用
	err = tx.Model(model.RoomBgs{}).Where("room_id = ?", roomInfo.Id).Updates(map[string]interface{}{
		"is_use": 1,
		// "is_cron": 1,
	}).Error
	if err != nil {
		tx.Rollback()
		return
	}
	var one model.RoomBgs
	tx.Model(model.RoomBgs{}).
		Where(model.RoomBgs{
			RoomId: roomInfo.Id,
			TrbrId: int(bgInfo.Id),
		}).First(&one)
	if one.Id > 0 {
		err = tx.Model(model.RoomBgs{}).Where("id = ?", one.Id).Updates(model.RoomBgs{IsUse: 2, IsCron: 1}).Error
	} else {
		err = tx.Model(model.RoomBgs{}).Create(&model.RoomBgs{
			RoomId:     roomInfo.Id,
			TrbrId:     int(bgInfo.Id),
			ExpireTime: time.Now().AddDate(100, 0, 0),
			IsUse:      2,
			IsCron:     1,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}).Error
	}
	if err != nil {
		tx.Rollback()
		return
	}
	roomInfo.BackgroundImg = bgInfo.Backgroud
	err = tx.Model(model.Room{}).Where("id = ?", roomInfo.Id).Save(roomInfo).Error
	if err != nil {
		tx.Rollback()
		return
	}
	return
}

type RoomBgsRDao struct {
}

func (RoomBgsRDao) GetBgsById(id int) (res model.RoomBgsResource) {
	coreDb.GetMasterDb().Model(model.RoomBgsResource{}).Where("id = ? and status = 1", id).First(&res)
	return
}

func (RoomBgsRDao) GetBgsByImg(img string) (res model.RoomBgsResource) {
	coreDb.GetMasterDb().Model(model.RoomBgsResource{}).Where("backgroud = ? and status = 1", img).First(&res)
	return
}

func (RoomBgsRDao) GetBgs(types int) (res []model.RoomBgsResource) {
	db := coreDb.GetMasterDb().Model(model.RoomBgsResource{}).Where("status = 1")
	if types > 0 {
		db = db.Where("types = ?", types)
	}
	db.Find(&res)
	return
}

func (RoomBgsRDao) GetDefaultBgs() (res model.RoomBgsResource) {
	coreDb.GetMasterDb().Model(model.RoomBgsResource{}).Where("status = ? and types = ?", 1, 1).Order("create_time asc").First(&res)
	return
}
