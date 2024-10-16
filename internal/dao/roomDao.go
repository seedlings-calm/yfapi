package dao

import (
	"errors"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	model2 "yfapi/internal/model/guild"
	typedef_enum "yfapi/typedef/enum"

	"gorm.io/gorm"
)

type RoomDao struct {
}

// Create 添加
func (u *RoomDao) Create(data *model.Room) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// Update 修改
func (u *RoomDao) Update(data model.Room) (err error) {
	err = coreDb.GetMasterDb().Model(&model.Room{Id: data.Id}).Updates(data).Error
	return
}

// Save 修改  会保存空值
func (u *RoomDao) Save(data model.Room) (err error) {
	err = coreDb.GetMasterDb().Model(&model.Room{Id: data.Id}).Save(data).Error
	return
}

// FindOne 条件查询
func (u *RoomDao) FindOne(param *model.Room) (data *model.Room, err error) {
	data = new(model.Room)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	return
}

// FindList 查询列表
func (u *RoomDao) FindList(param *model.Room) (result []model.Room, err error) {
	err = coreDb.GetMasterDb().Where(param).Where(model.Room{Status: typedef_enum.RoomStatusNormal}).Find(&result).Error
	return
}
func (u *RoomDao) FindListByLiveType(userId string, liveType int) (result []model.Room, err error) {
	err = coreDb.GetMasterDb().Model(&model.Room{}).Where("user_id=? and live_type = ? and status<3", userId, liveType).Find(&result).Error
	return
}

// FindByIds 根据ids查询结果
func (u *RoomDao) FindByIds(ids []string, sorts string) (result []model.Room) {
	coreDb.GetMasterDb().Where(model.Room{Status: typedef_enum.RoomStatusNormal}).Find(&result, ids).Order(sorts)
	return
}

// FindRoomByKeyword 房间类型 关键字 查询房间列表
func (u *RoomDao) FindRoomByKeyword(keyword string, liveType, page, size int) (result []model.Room, err error) {
	keyword = "%" + keyword + "%"
	tx := coreDb.GetSlaveDb().Model(model.Room{}).Where("(name like ? or room_no like ?) and status=?", keyword, keyword, typedef_enum.UserStatusNormal)
	if liveType > 0 {
		tx = tx.Where("live_type", liveType)
	}
	err = tx.Order("name").Limit(size).Offset(page * size).Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (u *RoomDao) GetRoomList(page, size, roomType int) (result []model.Room, count int64, err error) {
	tx := coreDb.GetSlaveDb().Model(&model.Room{}).
		Where(model.Room{RoomType: roomType, Status: typedef_enum.RoomStatusNormal, HiddenStatus: 2}).
		Count(&count)
	err = tx.Limit(size).Offset((page) * size).Find(&result).Error
	if err != nil {
		return
	}
	return
}

func (u *RoomDao) GetTopsRoomList(page, size int) (result []model.Room, count int64, err error) {
	tx := coreDb.GetSlaveDb().Model(&model.Room{}).
		Where(model.Room{Status: typedef_enum.RoomStatusNormal, HiddenStatus: 2}).
		Count(&count)
	err = tx.Limit(size).Offset((page) * size).Find(&result).Error
	if err != nil {
		return
	}
	return
}

func (u *RoomDao) GetRoomById(roomId string) (result model.Room, err error) {
	err = coreDb.GetSlaveDb().Model(model.Room{}).Where("id=? and status!=3", roomId).Scan(&result).Error
	return
}

func (u *RoomDao) GetRoomByIds(Ids []string) (result []model.Room, err error) {
	err = coreDb.GetSlaveDb().Model(model.Room{}).Where("id in ? and status!=3", Ids).Find(&result).Error
	return
}

func (u *RoomDao) GetRoomByIdsInterface(Ids []interface{}) (result []model.Room, err error) {
	err = coreDb.GetSlaveDb().Model(model.Room{}).Where("id in ? and status!=3", Ids).Find(&result).Error
	return
}

// GetRoomMapByIdList 根据房间id列表查询房间信息map
func (u *RoomDao) GetRoomMapByIdList(idList []string) (result map[string]model.Room, err error) {
	result = map[string]model.Room{}
	res, e := u.GetRoomByIds(idList)
	if e != nil {
		err = e
		return
	}
	for _, info := range res {
		result[info.Id] = info
	}
	return
}

type PositionRes struct {
	Num      int `json:"num"`
	IsBoss   int `json:"is_boss"`
	RoomType int `json:"room_type"`
}

// 获取当前房间使用的模板麦位  TODO:  后台操作关闭房间，需要删除redis 房间的麦位信息
func (u *RoomDao) GetRoomPositions(roomId string) (res PositionRes, err error) {

	err = coreDb.GetMasterDb().Table("t_room_template as a").
		Joins("inner join t_room as b on b.template_Id = a.id").
		Select("a.seat_list_count as num,a.is_boss,b.room_type").
		Where("b.id = ?", roomId).Scan(&res).Error
	return
}
func (u *RoomDao) GetRoomTypeList() (res []*model.RoomType, err error) {
	err = coreDb.GetMasterDb().Model(model.RoomType{}).Scan(&res).Error
	return
}

func (u *RoomDao) GetGuildRoomCount(guildId string) (count int64, err error) {
	err = coreDb.GetMasterDb().Model(model.Room{}).Where("guild_id=?", guildId).Count(&count).Error
	return
}

// 添加房间申请
func (u *RoomDao) AddRoomApply(data *model2.GuildRoomApply) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

func (u *RoomDao) GetRoomsByGuildId(guildId string) (res []model.Room, err error) {
	err = coreDb.GetMasterDb().Model(model.Room{}).Where("guild_id=?  and status!=3", guildId).Find(&res).Error
	return
}
func (u *RoomDao) GetRoomsByGuildIdAndLiveType(guildId string, liveType string) (res []model.Room, err error) {
	err = coreDb.GetMasterDb().Model(model.Room{}).Where("guild_id=? and live_type = ?  and status!=3", guildId, liveType).Find(&res).Error
	return
}

func (u *RoomDao) GetRoomIdsByGuildId(guildId string) (res []string, err error) {
	err = coreDb.GetMasterDb().Model(model.Room{}).Where("guild_id=?  and status!=3", guildId).Pluck("id", &res).Error
	return
}

// 随机获取房间号
func (u *RoomDao) GetRoomNo() (res *model.DicRoomNo, err error) {
	res = new(model.DicRoomNo)
	err = coreDb.GetMasterDb().Model(model.DicRoomNo{}).Select("room_no").Where("status=?", false).Order("rand()").First(res).Error
	if err != nil {
		return nil, err
	}
	// 更新获取到的房号
	err = coreDb.GetMasterDb().Model(model.DicRoomNo{}).Where("id=?", res.Id).Update("status", true).Error
	if err != nil {
		return nil, err
	}
	return
}

// 查询是否已经申请过
func (u *RoomDao) FindRoomApply(userId string, status int) (res *model2.GuildRoomApply, err error) {
	res = new(model2.GuildRoomApply)
	err = coreDb.GetMasterDb().Model(model2.GuildRoomApply{}).Where("room_user_id=? and status=?", userId, status).First(res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}
