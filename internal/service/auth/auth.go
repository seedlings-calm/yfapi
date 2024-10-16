package auth

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	dao2 "yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"

	"github.com/spf13/cast"
)

type Auth struct {
}

// 添加权限
func (a *Auth) AddRule(name, title string, ruleId int) error {
	data := model.AuthRule{
		RuleID: ruleId,
		Name:   name,
		Title:  title,
		Status: 1,
	}
	dao := new(dao2.UserAuthDao)
	err := dao.AddRule(&data)
	if err != nil {
		return err
	}
	return nil
}

// GiveUserRole
//
//	@Description: 给用户分配角色
//	@receiver a
func (a *Auth) GiveUserRole(userId, roomId string, roleId int) error {
	data := model.AuthRoleAccess{
		UserID: userId,
		RoleID: roleId,
		RoomID: roomId,
	}
	dao := new(dao2.UserAuthDao)
	err := dao.AddAuthRoleAccess(&data)
	if err != nil {
		return err
	}
	a.ClearRedisCache(userId, roomId)
	return nil
}

// DelUserRole
//
//	@Description: 删除用户角色
//	@receiver a
//	@param userId
//	@param roomId
//	@param roleId
//	@return error
func (a *Auth) DelUserRole(userId, roomId string, roleId int) error {
	data := model.AuthRoleAccess{
		UserID: userId,
		RoleID: roleId,
		RoomID: roomId,
	}
	dao := new(dao2.UserAuthDao)
	err := dao.DelAuthRoleAccess(&data)
	if err != nil {
		return err
	}
	a.ClearRedisCache(userId, roomId)
	return nil
}

// VerifyAuth
//
//	@Description:验证权限
//	@receiver a
//	@param userId
//	@param roomId
//	@param isCompere
//	@param ruleName
func (a *Auth) VerifyAuth(userId, roomId string, isCompere bool, ruleName ...string) (map[string]bool, error) {
	key := redisKey.UserRules(userId, roomId)
	if isCompere {
		key = redisKey.UserCompereRules(userId, roomId)
	}
	cacheRules := a.getRedisCache(key)
	res := map[string]bool{}
	var err error
	rules := []string{}
	if len(cacheRules) > 0 {
		err = json.Unmarshal([]byte(cacheRules), &rules)
		if err != nil {
			coreLog.LogError("%+v", err)
		}
	} else {
		rules, err = a.getUserRules(userId, roomId, isCompere)
		if err != nil {
			coreLog.LogError("%+v", err)
			return res, err
		}
		if len(rules) > 0 {
			marshal, _ := json.Marshal(rules)
			a.setRedisCache(key, string(marshal), time.Hour)
		}
	}
	ruleMap := map[string]struct{}{}
	for _, v := range rules {
		ruleMap[v] = struct{}{}
	}
	for _, name := range ruleName {
		if _, ok := ruleMap[name]; ok {
			res[name] = true
		} else {
			res[name] = false
		}
	}
	return res, nil
}

func (a *Auth) getUserRules(userId string, roomId string, isCompere bool) ([]string, error) {
	sql := `select role.rules from ` + new(model.AuthRoleAccess).TableName() + ` access inner join ` + new(model.AuthRole).TableName() + ` role on access.role_id=role.role_id where role.status=1 and access.user_id=? and (access.room_id=? or access.room_id = 0)`
	if !isCompere {
		sql = `select role.rules from ` + new(model.AuthRoleAccess).TableName() + ` access inner join ` + new(model.AuthRole).TableName() + ` role on access.role_id=role.role_id where role.status=1 and access.user_id=? and access.role_id!=1005 and (access.room_id=? or access.room_id = 0)`
	}
	ruleSlice := []string{}
	err := coreDb.GetMasterDb().Raw(sql, userId, roomId).Scan(&ruleSlice).Error
	if err != nil {
		return nil, err
	}
	if len(ruleSlice) == 0 { //没有特殊身份则是普通身份
		sql = `select rules from ` + new(model.AuthRole).TableName() + ` where role_id=` + cast.ToString(enum.NormalRoleId)
		ruleStr := ""
		err = coreDb.GetMasterDb().Raw(sql).Scan(&ruleStr).Error
		if err != nil {
			return nil, err
		}
		ruleSlice = append(ruleSlice, strings.Split(ruleStr, ",")...)
	}
	sql = `select name from ` + new(model.AuthRule).TableName() + ` where rule_id in(` + strings.Join(ruleSlice, ",") + `)`
	res := []string{}
	err = coreDb.GetMasterDb().Raw(sql).Scan(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

type RoleInfo struct {
	RoleId int    `json:"roleId"`
	Icon   string `json:"icon"`
}

// 获取用户在当前房间所有角色
func (a *Auth) GetRoleListByRoomIdAndUserId(roomId, userId string) ([]int, []RoleInfo, error) {
	key := redisKey.UserRoles(userId, roomId)
	var roles []RoleInfo
	var roleIds []int
	cache := a.getRedisCache(key)
	if len(cache) > 0 {
		err := json.Unmarshal([]byte(cache), &roles)
		if err != nil {
			coreLog.LogError("%+v", err)
		}
	} else {
		sql := "select ra.role_id, ar.icon from t_auth_role_access ra left join t_auth_role ar on ra.role_id=ar.role_id where (ra.room_id = ? or ra.room_id = 0) and ra.user_id = ? order by ra.role_id"
		err := coreDb.GetMasterDb().Raw(sql, roomId, userId).Scan(&roles).Error
		if err != nil {
			return roleIds, roles, err
		}
		if len(roles) > 0 {
			marshal, _ := json.Marshal(roles)
			a.setRedisCache(key, string(marshal), time.Hour)
		}
	}
	for _, info := range roles {
		roleIds = append(roleIds, info.RoleId)
	}

	return roleIds, roles, nil
}

// GetRoomRoleListByRoleId 查询房间某个角色的所有用户
func (a *Auth) GetRoomRoleListByRoleId(roomId string, roleId int) (userIdList []string) {
	_ = coreDb.GetSlaveDb().Model(model.AuthRoleAccess{}).Where("room_id=? and role_id=?", roomId, roleId).Select("user_id").Scan(&userIdList)
	return
}

func (a *Auth) IsHasRoomRole(roomId, userId string, roleId int) bool {
	roleIdList, _, _ := a.GetRoleListByRoomIdAndUserId(roomId, userId)
	for _, id := range roleIdList {
		if id == roleId {
			return true
		}
	}
	return false
}

// 是否超管
func (a *Auth) IsSuperAdminRole(roomId, userId string) bool {
	roleIdList, _, _ := a.GetRoleListByRoomIdAndUserId(roomId, userId)
	for _, id := range roleIdList {
		if id == enum.SuperAdminRoleId || id == enum.PatrolRoleId {
			return true
		}
	}
	return false
}

// IsHaveCurrRole 是否拥有传入的身份，满足任意一个即可
func (a *Auth) IsHaveCurrRole(roomId, userId string, checkRoleIdList []int) bool {
	roleIdList, _, _ := a.GetRoleListByRoomIdAndUserId(roomId, userId)
	for _, id := range roleIdList {
		for _, roleId := range checkRoleIdList {
			if id == roleId {
				return true
			}
		}
	}
	return false
}

func (a *Auth) ClearRedisCache(userId, roomId string) {
	key1 := redisKey.UserRules(userId, roomId)
	key2 := redisKey.UserRoles(userId, roomId)
	key3 := redisKey.UserCompereRules(userId, roomId)
	coreRedis.GetUserRedis().Del(context.Background(), key1, key2, key3)
}

func (a *Auth) setRedisCache(key string, data any, expire time.Duration) {
	coreRedis.GetUserRedis().Set(context.Background(), key, data, expire)
}

func (a *Auth) getRedisCache(key string) string {
	result, err := coreRedis.GetUserRedis().Get(context.Background(), key).Result()
	if err != nil {
		coreLog.LogError("getRedisCache err:%+v", err)
		return ""
	}
	return result
}

// 根据权限获取用户所在房间有此权限的房间列表
func (a *Auth) GetUsersRules(userId string, roleId ...int) (res []string) {
	_ = coreDb.GetSlaveDb().Model(model.AuthRoleAccess{}).Where("user_id=? and role_id in ?", userId, roleId).Select("room_id").Scan(&res)
	return
}
