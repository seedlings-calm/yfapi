package dao

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	typedef_enum "yfapi/typedef/enum"
)

type UserDao struct {
}

// Create 添加
func (u *UserDao) Create(data *model.User) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// UpdateById Update 修改
func (u *UserDao) UpdateById(data *model.User) (err error) {
	err = coreDb.GetMasterDb().Model(model.User{Id: data.Id}).Updates(data).Error
	return
}

// FindOne 条件查询
func (u *UserDao) FindOne(param *model.User) (data *model.User, err error) {
	data = new(model.User)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	return
}
func (u *UserDao) FindOneAndRows(param *model.User) (data *model.User, rows int64, err error) {
	data = new(model.User)
	result := coreDb.GetMasterDb().Where(param).First(data)
	rows = result.RowsAffected //返回找到的记录数
	err = result.Error
	return
}
func (u *UserDao) Count(param *model.User) (count int64) {
	coreDb.GetMasterDb().Model(param).Where(param).Count(&count)
	return
}

// FindList 查询列表
func (u *UserDao) FindList(param *model.User) (result []model.User, err error) {
	err = coreDb.GetMasterDb().Where(param).Find(&result).Error
	return
}

// FindByIds 根据ids查询结果
func (u *UserDao) FindByIds(ids []string) (result []model.User) {
	coreDb.GetMasterDb().Find(&result, ids)
	return
}

func (u *UserDao) UpdateUserFieldsByUserId(userId string, data map[string]any) error {
	err := coreDb.GetMasterDb().Model(&model.User{}).Where("id = ?", userId).Updates(data).Error
	return err
}

// FindUserByKeyword
//
//	@Description: 模糊查询符合的user_no或nickname
//	@receiver u
//	@param keyword string -
//	@return result -
//	@return err -
func (u *UserDao) FindUserByKeyword(keyword string, page, size int) (result []model.User, err error) {
	keyword = "%" + keyword + "%"
	err = coreDb.GetSlaveDb().Model(model.User{}).Where("(nickname like ? or user_no like ?) and status=?", keyword, keyword, typedef_enum.UserStatusNormal).Order("nickname").Limit(size).Offset(page * size).Find(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// FindUserByMobile
//
//	@Description: 通过手机号查询用户
func (u *UserDao) FindUserByMobile(mobile string) (result model.User, err error) {
	err = coreDb.GetSlaveDb().Model(model.User{}).Where("mobile = ? and status=?", mobile, typedef_enum.UserStatusNormal).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// FindUserByUserNo
//
//	@Description: 通过userNo查询用户
func (u *UserDao) FindUserByUserNo(userno string) (result model.User, err error) {
	err = coreDb.GetSlaveDb().Model(model.User{}).Where("user_no = ? and status=?", userno, typedef_enum.UserStatusNormal).Scan(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

// 检测是否有重复昵称存在
func (u *UserDao) CheckRepeatNickname(nickname string) bool {
	count := u.Count(&model.User{
		Nickname: nickname,
	})
	if count > 0 {
		return true
	}
	return false
}

// 用户登录信息记录
func (u *UserDao) UserLoginRecord(data *model.UserLoginRecord) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// 根据偏移量获取用户iD
func (u *UserDao) GetUserIdsOffsetId(id int64, limit int) []string {
	userIds := []string{}
	coreDb.GetMasterDb().Model(&model.User{}).Where("status = ? and source in ? and id > ?", 1, []int{0, 1}, id).Order("id asc").Limit(limit).Pluck("id", &userIds)
	return userIds
}

// 根据手机号获取对应用户
func (u *UserDao) GetUsersByMobile(regionCode, mobile string) []model.User {
	res := []model.User{}
	if len(regionCode) == 0 || len(mobile) == 0 {
		return res
	}
	coreDb.GetMasterDb().Where("region_code = ? and mobile = ? and status in ?", regionCode, mobile, []int{typedef_enum.UserStatusNormal, typedef_enum.UserStatusFreezing, typedef_enum.UserStatusApplyInvalid}).Find(&res)
	return res
}

// 判断密码是否重复
func (u *UserDao) RepeatPassword(mobile string, password string) bool {
	var num int64
	coreDb.GetMasterDb().Model(&model.User{}).Where("mobile = ? and password = ?", mobile, password).Count(&num)
	if num > 0 {
		return true
	}
	return false
}

func (u *UserDao) UpdateMobile(oldRegionCode, oldMobile, NewRegionCode, newMobile string) (err error) {
	err = coreDb.GetMasterDb().Model(model.User{}).Where("region_code = ? and mobile = ?", oldRegionCode, oldMobile).Updates(map[string]any{"region_code": NewRegionCode, "mobile": newMobile}).Error
	return
}

// 检查用户身份证号是否正确
func (u *UserDao) CheckUserIdNo(idNo, userId string) (res model.UserRealName, err error) {
	err = coreDb.GetMasterDb().Model(model.UserRealName{}).Where("id_no =?", idNo).Where("user_id=?", userId).Scan(&res).Error
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	return
}

// 查询用户身份证号
func (u *UserDao) FindUserIdNo(userId string) (res model.UserRealName, err error) {
	err = coreDb.GetMasterDb().Model(model.UserRealName{}).Where("user_id=?", userId).Scan(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (u *UserDao) GetUserIdsByMobile(regionCode, mobile string) (res []string) {
	coreDb.GetMasterDb().Model(&model.User{}).Where("region_code = ? and mobile = ?", regionCode, mobile).Pluck("id", &res)
	return
}

// 判断用户是否实名是否是待审核或审核通过
func (u *UserDao) FindUserIsRealName(userId string) (res []string, err error) {
	userIds := []string{}
	coreDb.GetMasterDb().Model(&model.UserRealName{}).Where("user_id =? and status in?", userId, []int{1, 2}).Pluck("id", &userIds)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return userIds, err
}
