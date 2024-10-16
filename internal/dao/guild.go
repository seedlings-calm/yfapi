package dao

import (
	"errors"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
	response_guild "yfapi/typedef/response/guild"

	"gorm.io/gorm"
)

type GuildDao struct {
}

// FindById
//
//	@Description: 查询工会是否存在
//	@receiver g
//	@param param
//	@return data
//	@return err
func (g *GuildDao) FindById(param *model.Guild) (data *model.Guild, err error) {
	data = new(model.Guild)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	return
}

// 查询用户是否已经已在工会中
func (g *GuildDao) GetCheckUserInGuild(guildId, userId string) bool {
	guildMember := new(model.GuildMember)
	err := coreDb.GetMasterDb().Model(&model.GuildMember{}).Where("guild_id = ? AND user_id = ? AND status!=3", guildId, userId).First(&guildMember).Error
	if err != nil {
		return false
	}
	return true
}

// 查询公会成员信息
func (g *GuildDao) GetGuildMemberInfo(userId string) (res model.GuildMember, err error) {
	err = coreDb.GetSlaveDb().Model(model.GuildMember{}).Where("user_id=? and status!=3", userId).Order("create_time desc").First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// IsUserMemberOfOtherGuild
// 查询用户是否已经是其他工会成员
func (g *GuildDao) IsGuildMember(userId string) (res *model.GuildMember, err error) {
	err = coreDb.GetMasterDb().Model(&model.GuildMember{}).Where("user_id = ? and status<?", userId, enum.GuildMemberStatusLeave).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// 查询用户是否已经申请
func (g *GuildDao) GetCheckUserApplication(data *model.GuildMemberApply) (res *model.GuildMemberApply, err error) {
	err = coreDb.GetMasterDb().Model(&model.GuildMemberApply{}).Where(data).Where("status = ? and apply_type=1", 1).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// 申请加入工会
func (g *GuildDao) Create(data *model.GuildMemberApply) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// 查询工会列表
func (g *GuildDao) GetGuildList() (res []*model.Guild, err error) {
	err = coreDb.GetMasterDb().Model(&model.Guild{}).Find(&res).Error
	return

}

// 查询房间列表
func (g *GuildDao) GetRoomList(guildId string) (res []*model.Room, err error) {
	err = coreDb.GetMasterDb().Model(&model.Room{}).Where("guild_id", guildId).Find(&res).Error
	return
}

// 通过工会id查询工会成员的数量
func (g *GuildDao) GerGuildMember(guildId string) (count int64, err error) {
	coreDb.GetMasterDb().Model(&model.GuildMember{}).Where("guild_id=? and status<?", guildId, enum.GuildMemberStatusLeave).Count(&count)
	return
}

func (g *GuildDao) GetGuildNameByUserId(userId string) (name string) {
	coreDb.GetMasterDb().Table("t_guild as a").
		Joins("left join t_guild_member b on b.guild_id = a.id").
		Where("b.user_id = ? and b.status<?", userId, enum.GuildMemberStatusLeave).
		Select("a.name").Scan(&name)
	return
}
func (g *GuildDao) GetGuildMemberList(guildId string) (res []*model.GuildMember, err error) {
	err = coreDb.GetMasterDb().Model(&model.GuildMember{}).Where("guild_id=? and status<?", guildId, enum.GuildMemberStatusLeave).Find(&res).Error
	return
}

// GetGuildInfo
//
//	@Description: 根据公会id查询公会信息
//	@receiver g
//	@param guildId string -
//	@return res -
//	@return err -
func (g *GuildDao) GetGuildInfo(guildId string) (res response_guild.GuildInfo, err error) {
	err = coreDb.GetSlaveDb().Table("t_guild g").Joins("left join t_user u on u.id=g.user_id").Where("g.id=? and g.status=?", guildId, enum.GuildStatusNormal).
		Select("g.id guild_id, g.logo_img, g.guild_no, u.avatar user_avatar, u.user_no, u.nickname user_name").Scan(&res).Error
	return
}

// GetGuildByIdList
//
//	@Description: 根据公会ID列表查询公会信息list
//	@receiver g
//	@param guildIdList []string -
//	@return res -
//	@return err -
func (g *GuildDao) GetGuildByIdList(guildIdList []string) (res []model.Guild, err error) {
	err = coreDb.GetSlaveDb().Model(model.Guild{}).Where("id in ? and status=1", guildIdList).Scan(&res).Error
	return
}

// GetGuildMapByIdList
//
//	@Description: 根据公会ID列表查询公会信息map
//	@receiver g
//	@param guildIdList []string -
//	@return res -
//	@return err -
func (g *GuildDao) GetGuildMapByIdList(guildIdList []string) (res map[string]model.Guild, err error) {
	res = make(map[string]model.Guild)
	result, e := g.GetGuildByIdList(guildIdList)
	if e != nil {
		err = e
		return
	}
	for _, info := range result {
		res[info.ID] = info
	}
	return
}
