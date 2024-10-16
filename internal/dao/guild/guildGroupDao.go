package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/request/guild"

	"github.com/gin-gonic/gin"
)

type GuildGroupDao struct {
}

// Create 添加
func (g *GuildGroupDao) Create(data *model.GuildGroup) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

func (g *GuildGroupDao) Save(data *model.GuildGroup) (err error) {
	err = coreDb.GetMasterDb().Model(data).Save(data).Error
	return
}

// FindOne 根据ID条件查询
func (g *GuildGroupDao) FindOne(param *model.GuildGroup) (data *model.GuildGroup, err error) {
	data = new(model.GuildGroup)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	return
}

// Count 通过条件计算数量
func (g *GuildGroupDao) Count(param *model.GuildGroup) (count int64) {
	coreDb.GetMasterDb().Model(param).Where(param).Count(&count)
	return
}

// FindList 查询列表
func (g *GuildGroupDao) FindList(param *model.GuildGroup) (result []model.GuildGroup, err error) {
	err = coreDb.GetMasterDb().Where(param).Find(&result).Error
	return
}

// FindByIds 根据ids查询结果
func (g *GuildGroupDao) FindByIds(ids []string) (result []model.GuildGroup) {
	coreDb.GetMasterDb().Find(&result, ids)
	return
}

// GetGuildGroupListPage 获取公会成员分组列表
func (r *GuildGroupDao) GetGuildGroupListPage(req *guild.GuildGroupListreq, c *gin.Context) (list interface{}, count int64, err error) {
	limit := req.Size
	offset := req.Size * (req.CurrentPage - 1)
	guildID := c.GetString("guildId")
	db := coreDb.GetMasterDb().Table("t_guild_group tgg").
		Where("tgg.guild_id = ? and tgg.status = 1", guildID)
	var dataList []guild.MemberGroupUpdateRes

	err = db.Count(&count).Error

	if err != nil {
		return
	}
	db = db.Joins("left join t_guild_member tgm on tgg.id = tgm.group_id and tgm.status != 3")

	err = db.Select("tgg.*,IFNULL(count(tgm.group_id),0) person_count").
		Group("tgg.id").
		Limit(limit).Offset(offset).Scan(&dataList).Error
	return dataList, count, err
}
