package dao

import (
	"errors"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	model2 "yfapi/internal/model/guild"
	"yfapi/typedef/request/guild"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GuildMemberDao struct {
}

// Create 添加
func (g *GuildMemberDao) Create(data *model.GuildMember) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// FindOne 根据ID条件查询
func (g *GuildMemberDao) FindOne(param *model.GuildMember) (data *model.GuildMember, err error) {
	data = new(model.GuildMember)
	err = coreDb.GetMasterDb().Where(param).Order("id desc").First(data).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// Count 通过条件计算数量
func (g *GuildMemberDao) Count(param *model.GuildMember) (count int64) {
	coreDb.GetMasterDb().Model(param).Where(param).Count(&count)
	return
}

// FindList 查询列表
func (g *GuildMemberDao) FindList(param *model.GuildMember) (result []model.GuildMember, err error) {
	err = coreDb.GetMasterDb().Where(param).Find(&result).Error
	return
}

// FindByIds 根据ids查询结果
func (g *GuildMemberDao) FindByIds(ids []string) (result []model.GuildMember) {
	coreDb.GetMasterDb().Find(&result, ids)
	return
}

// GetRoomListPage 获取公会后台房间列表
func (r *GuildMemberDao) GetMemberListPage(req *guild.GuildMemberListreq, c *gin.Context) (list interface{}, count int64, err error) {
	limit := req.Size
	offset := req.Size * (req.CurrentPage - 1)
	guildID := c.GetString("guildId")
	db := coreDb.GetSlaveDb().Model(&model2.GuildMember{GuildID: guildID})
	var dataList []model2.GuildMember

	// 如果有条件搜索 下方会自动创建搜索语句
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.UserKeyword != "" {
		db = db.Where("user_id LIKE ?", "%"+req.UserKeyword+"%")

	}
	//if req.RoomType != 0 {
	//	db = db.Where("room_type = ?", req.RoomType)
	//}
	err = db.Where(&model2.GuildMember{GuildID: guildID}).Count(&count).Error
	if err != nil {
		return
	}
	err = db.Preload("Users").Limit(limit).Offset(offset).Find(&dataList).Error

	return dataList, count, err
}
