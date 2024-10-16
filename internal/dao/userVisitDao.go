package dao

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type UserVisitDao struct {
}

// RecordUserVisit 记录用户访问
func (u *UserVisitDao) RecordUserVisit(userId, targetUserId string, isHidden bool) error {
	// 是否有访问记录
	var record model.UserVisit
	err := coreDb.GetSlaveDb().Model(record).Where("user_id=? and target_user_id=?", userId, targetUserId).First(&record).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	isVisit := false
	if !record.IsVisit { // 没有互相访问
		// 查询是否互相访问过
		var result model.UserVisit
		err = coreDb.GetSlaveDb().Model(record).Where("user_id=? and target_user_id=?", targetUserId, userId).First(&result).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if result.ID > 0 { // 互相访问
			isVisit = true
			// 更新互相访问标识
			result.IsVisit = true
			_ = coreDb.GetMasterDb().Model(result).Save(&result).Error
		}

	} else {
		isVisit = record.IsVisit
	}
	if record.ID == 0 { // 新增记录
		record = model.UserVisit{
			UserId:       userId,
			TargetUserId: targetUserId,
			IsVisit:      isVisit,
			VisitCount:   1,
			VisitHidden:  isHidden,
			ClearVisit:   false,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err = coreDb.GetMasterDb().Model(record).Create(&record).Error
		if err != nil {
			return err
		}
	} else { // 更新记录
		record.IsVisit = isVisit
		record.VisitCount++
		record.VisitHidden = isHidden
		record.ClearVisit = false
		record.UpdateTime = time.Now()
		err = coreDb.GetMasterDb().Model(record).Save(&record).Error
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUserVisitRecordList 查询用户足迹列表
func (u *UserVisitDao) GetUserVisitRecordList(userId string, page, size int) (result []model.UserVisitDTO, count int64, err error) {
	tx := coreDb.GetSlaveDb().Table("t_user_visit uv").Joins("left join t_user u on u.id=uv.target_user_id").Where("uv.user_id=? and uv.clear_visit=0", userId).Count(&count)
	err = tx.Select("uv.*, u.nickname, u.sex, u.avatar, u.introduce").Order("uv.update_time desc").Limit(size).Offset(page * size).Scan(&result).Error
	return
}

// GetVisitUserRecordList 查询用户访客列表
func (u *UserVisitDao) GetVisitUserRecordList(userId string, page, size int) (result []model.UserVisitDTO, count int64, err error) {
	tx := coreDb.GetSlaveDb().Table("t_user_visit uv").Joins("left join t_user u on u.id=uv.user_id").Where("uv.target_user_id=? and uv.visit_hidden=0", userId).Count(&count)
	err = tx.Select("uv.*, u.nickname, u.sex, u.avatar, u.introduce").Order("uv.update_time desc").Limit(size).Offset(page * size).Scan(&result).Error
	return
}

// GetVisitUserCount
//
//	@Description: 查询用户访客数量
//	@receiver u
//	@param userId string -
//	@return count -
func (u *UserVisitDao) GetVisitUserCount(userId string) int {
	var count int64
	_ = coreDb.GetSlaveDb().Model(model.UserVisit{}).Where("target_user_id=? and visit_hidden=0", userId).Count(&count).Error
	return int(count)
}

// ClearUserVisitRecord 清空我的足迹
func (u *UserVisitDao) ClearUserVisitRecord(userId string) (err error) {
	err = coreDb.GetMasterDb().Model(model.UserVisit{}).Where("user_id=? and clear_visit=0", userId).Update("clear_visit", 1).Error
	return
}
