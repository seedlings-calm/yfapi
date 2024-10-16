package dao

import (
	"errors"
	"gorm.io/gorm"
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
)

// TimelineDao
// @Description: 用户动态
type TimelineDao struct {
}

// Create
//
//	@Description: 发布动态
//	@receiver t
//	@param param *model.Timeline -
//	@return err -
func (t *TimelineDao) Create(param *model.Timeline) (err error) {
	return coreDb.GetMasterDb().Model(param).Create(param).Error
}

// Update
//
//	@Description: 更新动态
//	@receiver t
//	@param param *model.Timeline -
//	@return err -
func (t *TimelineDao) Update(param *model.Timeline) (err error) {
	return coreDb.GetMasterDb().Model(param).Updates(param).Error
}

// Delete
//
//	@Description: 删除动态
//	@receiver t
//	@param param *model.Timeline -
//	@return err -
func (t *TimelineDao) Delete(param *model.Timeline) (err error) {
	return coreDb.GetMasterDb().Model(param).Updates(&model.Timeline{Status: 4, UpdateTime: time.Now()}).Where("status!=4").Error
}

// IncrLoveCount
//
//	@Description: 增加点赞数
//	@receiver t
//	@param timelineId int64 -
//	@return err -
func (t *TimelineDao) IncrLoveCount(timelineId int64) (err error) {
	//return coreDb.GetMasterDb().Model(&model.Timeline{Id: timelineId}).Update("love_count", "love_count+1").Error
	return coreDb.GetMasterDb().Model(&model.Timeline{Id: timelineId}).Update("love_count", gorm.Expr("love_count+1")).Error
}

// DecrLoveCount
//
//	@Description: 减少点赞数
//	@receiver t
//	@param timelineId int64 -
//	@return err -
func (t *TimelineDao) DecrLoveCount(timelineId int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.Timeline{Id: timelineId}).Where("love_count>=1").Update("love_count", gorm.Expr("love_count-1")).Error
}

// IncrReplyCount
//
//	@Description: 增加评论数量
//	@receiver t
//	@param timelineId int64 -
//	@return err -
func (t *TimelineDao) IncrReplyCount(timelineId int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.Timeline{Id: timelineId}).Update("reply_count", gorm.Expr("reply_count+1")).Error
}

// DecrReplyCount
//
//	@Description: 减少评论数量
//	@receiver t
//	@param timelineId int64 -
//	@return err -
func (t *TimelineDao) DecrReplyCount(timelineId int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.Timeline{Id: timelineId}).Where("reply_count>=1").Update("reply_count", gorm.Expr("reply_count-1")).Error
}

// GetUserTimelineList
//
//	@Description: 获取用户动态列表
//	@receiver t
//	@param userId string -
//	@param page int -
//	@param size int -
//	@return result -
//	@return err -
func (t *TimelineDao) GetUserTimelineList(userId, selfUserId string, page, size int, isSelf bool) (result []model.Timeline, count int64, err error) {
	tx := coreDb.GetMasterDb().Model(&model.Timeline{})
	switchType := false
	if isSelf {
		tx = tx.Where("user_id=? and status!=4", userId).Count(&count)
	} else {
		switchType = new(UserTimelineFilterDao).GetSwitchType(userId, selfUserId, enum.DontLetHeSeeMoments)
		if !switchType {
			tx = tx.Where("user_id=? and status=1", userId).Count(&count)
		}
	}
	if !switchType {
		err = tx.Order("id desc").Limit(size).Offset(page * size).Scan(&result).Error
	}
	return
}

func (t *TimelineDao) GetTimelineById(id int64) (result model.Timeline, err error) {
	err = coreDb.GetMasterDb().Model(model.Timeline{}).Where("id = ? and status!=4", id).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// GetTimelineListLatest
//
//	@Description: 查询最新动态列表
//	@receiver t
//	@param page int -
//	@param size int -
//	@return result -
//	@return count -
//	@return err -
func (t *TimelineDao) GetTimelineListLatest(userId string, page, size int) (result []model.Timeline, count int64, err error) {
	ids := new(UserTimelineFilterDao).GetFilterUserIds(userId)
	tx := coreDb.GetMasterDb()
	if len(ids) > 0 {
		tx = tx.Model(model.Timeline{}).Where("status=1 and user_id not in (?)", ids)
	} else {
		tx = tx.Model(model.Timeline{}).Where("status=1")
	}
	tx.Count(&count)
	err = tx.Order("id desc").Limit(size).Offset(page * size).Scan(&result).Error
	return
}

// GetTimelineListFollow
//
//	@Description: 查询关注好友动态列表
//	@receiver t
//	@param userId string -
//	@param page int -
//	@param size int -
//	@return result -
//	@return count -
//	@return err -
func (t *TimelineDao) GetTimelineListFollow(userId string, page, size int) (result []model.Timeline, count int64, err error) {
	ids := new(UserTimelineFilterDao).GetFilterUserIds(userId)
	tx := coreDb.GetMasterDb().Table("t_timeline t")
	if len(ids) > 0 {
		tx.Joins("inner join t_user_follow uf on uf.focus_user_id=t.user_id").Where("uf.user_id=? and t.status=1 and t.user_id not in (?)", userId, ids).Count(&count)
	} else {
		tx.Joins("inner join t_user_follow uf on uf.focus_user_id=t.user_id").Where("uf.user_id=? and t.status=1", userId).Count(&count)
	}

	err = tx.Select("t.*").Order("t.id desc").Limit(size).Offset(page * size).Scan(&result).Error
	return
}

// GetUserTimelineLoveCount
//
//	@Description: 查询用户动态总点赞数量
//	@receiver t
//	@param userId string -
//	@return count -
//	@return err -
func (t *TimelineDao) GetUserTimelineLoveCount(userId string) (count int, err error) {
	err = coreDb.GetSlaveDb().Model(model.Timeline{}).Select("IFNULL(sum(love_count), 0)").Where("user_id=? and status!=4", userId).Scan(&count).Error
	return
}

// GetTimelineDetail
//
//	@Description: 查询正常的动态详情
//	@receiver t
//	@param timelineId int64 -
//	@return result -
//	@return err -
func (t *TimelineDao) GetTimelineDetail(timelineId int64) (result model.Timeline, err error) {
	err = coreDb.GetSlaveDb().Model(model.Timeline{}).Where("id=? and status=1", timelineId).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}
