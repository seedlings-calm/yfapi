package dao

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type TimelineReplyDao struct {
}

// Create
//
//	@Description: 新增评论
//	@receiver t
//	@param param *model.TimelineReply -
//	@return pkId -
//	@return err -
func (t *TimelineReplyDao) Create(param *model.TimelineReply) (pkId int64, err error) {
	pkId = coreDb.GetMasterDb().Model(param).Create(param).RowsAffected
	return
}

// Update
//
//	@Description: 更新评论
//	@receiver t
//	@param param *model.TimelineReply -
//	@return err -
func (t *TimelineReplyDao) Update(param *model.TimelineReply) (err error) {
	return coreDb.GetMasterDb().Model(param).Updates(param).Error
}

// Delete
//
//	@Description: 删除评论
//	@receiver t
//	@param id int64 -
//	@return err -
func (t *TimelineReplyDao) Delete(id int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.TimelineReply{}).Where(model.TimelineReply{Id: id}).Updates(model.TimelineReply{Status: 3}).Error
}

// GetTimelineReplyList
//
//	@Description: 获取评论列表（不包含子评论）
//	@receiver t
//	@param timelineId int64 -
//	@param page int -
//	@param size int -
//	@return result -
//	@return count -
//	@return err -
func (t *TimelineReplyDao) GetTimelineReplyList(timelineId int64, page, size int) (result []model.TimelineReply, count int64, err error) {
	tx := coreDb.GetMasterDb().Model(model.TimelineReply{}).Where("timeline_id=? and status=1 and to_reply_id=0", timelineId).Count(&count)
	err = tx.Order("praised_count desc, id desc").Limit(size).Offset(page * size).Scan(&result).Error
	return
}

// GetAllTimelineReplyList
//
//	@Description: 根据动态Id获取所有的评论、子评论列表
//	@receiver t
//	@param timelineId int64 -
//	@return result -
//	@return err -
func (t *TimelineReplyDao) GetAllTimelineReplyList(timelineId int64) (result []model.TimelineReply, err error) {
	err = coreDb.GetMasterDb().Model(model.TimelineReply{}).Where("timeline_id=? and status!=3", timelineId).Scan(&result).Error
	return
}

// GetTimelineSubReplyList
//
//	@Description: 获取子评论列表
//	@receiver t
//	@param timelineId int64 -
//	@param toReplyId int64 -
//	@param page int -
//	@param size int -
//	@return result -
//	@return count -
//	@return err -
func (t *TimelineReplyDao) GetTimelineSubReplyList(timelineId, toReplyId int64, page, size int) (result []model.TimelineReply, count int64, err error) {
	tx := coreDb.GetMasterDb().Model(model.TimelineReply{}).Where(&model.TimelineReply{TimelineId: timelineId, Status: 1, ToReplyId: toReplyId}).Count(&count)
	err = tx.Order("id").Limit(size).Offset(page * size).Scan(&result).Error
	return
}

// IncrSubReplyCount
//
//	@Description: 增加评论的回复数量
//	@receiver t
//	@param id int64 -
//	@return err -
func (t *TimelineReplyDao) IncrSubReplyCount(id int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.TimelineReply{Id: id}).Update("sub_reply_count", gorm.Expr("sub_reply_count+1")).Error
}

// DecrSubReplyCount
//
//	@Description: 减少评论的回复数量
//	@receiver t
//	@param id int64 -
//	@return err -
func (t *TimelineReplyDao) DecrSubReplyCount(id int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.TimelineReply{Id: id}).Where("sub_reply_count>=1").Update("sub_reply_count", gorm.Expr("sub_reply_count-1")).Error
}

// IncrSubReplyPraisedCount
//
//	@Description: 增加评论的点赞数量
//	@receiver t
//	@param id int64 -
//	@return err -
func (t *TimelineReplyDao) IncrSubReplyPraisedCount(id int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.TimelineReply{Id: id}).Update("praised_count", gorm.Expr("praised_count+1")).Error
}

// DecrSubReplyPraisedCount
//
//	@Description: 减少评论的点赞数量
//	@receiver t
//	@param id int64 -
//	@return err -
func (t *TimelineReplyDao) DecrSubReplyPraisedCount(id int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.TimelineReply{Id: id}).Where("praised_count>=1").Update("praised_count", gorm.Expr("praised_count-1")).Error
}

// GetTimelineReplyById
//
//	@Description: 查询评论信息
//	@receiver t
//	@param id int64 -
//	@return result -
//	@return err -
func (t *TimelineReplyDao) GetTimelineReplyById(id int64) (result model.TimelineReply, err error) {
	err = coreDb.GetMasterDb().Model(model.TimelineReply{}).Where(&model.TimelineReply{Id: id, Status: 1}).Scan(&result).Error
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		return
	}
	return
}

// DeleteReplyByTimelineId
//
//	@Description: 删除动态下的所有评论、子评论
//	@receiver t
//	@param timelineId int64 -
//	@return err -
func (t *TimelineReplyDao) DeleteReplyByTimelineId(timelineId int64) (err error) {
	return coreDb.GetMasterDb().Model(&model.TimelineReply{}).Where(model.TimelineReply{TimelineId: timelineId}).Updates(model.TimelineReply{Status: 3}).Error
}

// GetUserTimelineReplyPraisedCount
//
//	@Description: 查询用户评论的总点赞数量
//	@receiver t
//	@param userId string -
//	@return count -
//	@return err -
func (t *TimelineReplyDao) GetUserTimelineReplyPraisedCount(userId string) (count int, err error) {
	err = coreDb.GetSlaveDb().Model(model.TimelineReply{}).Select("IFNULL(sum(praised_count),0)").Where("replier_id=? and status!=3", userId).Scan(&count).Error
	return
}
