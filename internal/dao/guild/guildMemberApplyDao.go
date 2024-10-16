package dao

import (
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	model2 "yfapi/internal/model/guild"
)

type GuildMemberApplyDao struct {
}

// Create 添加
func (g *GuildMemberApplyDao) Create(data *model.GuildMemberApply) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// FindOne 根据ID条件查询
func (g *GuildMemberApplyDao) FindOne(param *model.GuildMemberApply) (data *model.GuildMemberApply, err error) {
	data = new(model.GuildMemberApply)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// Count 通过条件计算数量
func (g *GuildMemberApplyDao) Count(param *model.GuildMemberApply) (count int64) {
	coreDb.GetMasterDb().Model(param).Where(param).Count(&count)
	return
}

// FindList 查询列表
func (g *GuildMemberApplyDao) FindList(param *model.GuildMemberApply) (result []model.GuildMemberApply, err error) {
	err = coreDb.GetMasterDb().Where(param).Find(&result).Error
	return
}

// FindByIds 根据ids查询结果
func (g *GuildMemberApplyDao) FindByIds(ids []string) (result []model.GuildMemberApply) {
	coreDb.GetMasterDb().Find(&result, ids)
	return
}

// 公会成员申请修改
func (g *GuildMemberApplyDao) MemberApplyUpdate(data model.GuildMemberApply) (err error) {
	err = coreDb.GetMasterDb().Model(&model.GuildMemberApply{}).Where("id =?", data.ID).Updates(&data).Error
	return
}

// GetGuildMemberApply
//
//	@Description: 查询用户的公会申请记录
//	@receiver g
//	@param userId string -
//	@param applyType int -
//	@param status int -
//	@return result -
//	@return err -
func (g *GuildMemberApplyDao) GetGuildMemberApply(userId, guildId string, applyType, status int) (result model.GuildMemberApply, err error) {
	tx := coreDb.GetSlaveDb().Model(model2.GuildMemberApply{}).Where("user_id=? and guild_id=? and apply_type=?", userId, guildId, applyType)
	if status > 0 {
		tx = tx.Where("status=?", status)
	}
	err = tx.Order("create_time desc").Limit(1).Scan(&result).Error
	return
}
