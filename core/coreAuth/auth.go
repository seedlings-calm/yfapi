package coreAuth

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var coreAuth *CoreAuth
var coreAuthOnce sync.Once

type CoreAuth struct {
	dsn   string
	cache map[string]cacheData
	l     sync.Mutex
	p     Pool
}

type CoreAuthConfig struct {
	UserName string
	PassWord string
	Host     string
	Port     int
	Database string
}

func (c CoreAuthConfig) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4", c.UserName, c.PassWord, c.Host, c.Port, c.Database)
}

func Instance() *CoreAuth {
	return coreAuth
}

func New(config CoreAuthConfig) *CoreAuth {
	coreAuthOnce.Do(func() {
		auth := &CoreAuth{
			dsn:   config.GetDsn(),
			cache: map[string]cacheData{},
		}
		init, err := auth.poolInit()
		if err != nil {
			panic(errors.New("初始化连接池失败:" + err.Error()))
		}
		auth.p = init
		go auth.checkCache()
		fmt.Printf("初始化auth成功。连接池数:%d\r\n", auth.p.Len())
		coreAuth = auth
	})
	return coreAuth
}

func (this *CoreAuth) poolInit() (Pool, error) {
	//factory 创建连接的方法
	factory := func() (any, error) { return sql.Open("mysql", this.dsn) }
	//close 关闭连接的方法
	closed := func(v any) error { return v.(*sql.DB).Close() }
	//ping 检测连接的方法
	ping := func(v any) error { return v.(*sql.DB).Ping() }
	//创建一个连接池： 初始化5，最大空闲连接是20，最大并发连接30
	poolConfig := &Config{
		InitialCap: 1,  //资源池初始连接数
		MaxIdle:    1,  //最大空闲连接数
		MaxCap:     50, //最大并发连接数
		Factory:    factory,
		Close:      closed,
		Ping:       ping,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}
	return NewChannelPool(poolConfig)
}

// 添加权限
func (this *CoreAuth) AddRule(name, title string, ruleId, category int, categoryName string) error {
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `insert into ` + auth_rule_name + ` (rule_id,name,title,show_name,category,categoryName) values(?,?,?,?,?,?)`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ruleId, name, title, "", category, categoryName)
	if err != nil {
		return err
	}
	return nil
}

// 修改权限
func (this *CoreAuth) EditRule(id int, name, title string, category, ruleId int, categoryName string) error {
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `update ` + auth_rule_name + ` set rule_id=?, name=?,title=?,category=?,categoryName=? where id=?`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(ruleId, name, title, category, categoryName, id)
	if err != nil {
		return err
	}
	this.clear()
	return nil
}

// 删除权限
func (this *CoreAuth) DeleteRule(id int) error {
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `delete from ` + auth_rule_name + ` where id=?`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	this.clear()
	return nil
}

// 添加新角色
func (this *CoreAuth) AddRole(roleId int, title string, rules string) error {
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `insert into ` + auth_role_name + ` (role_id,title,rules) values(?,?,?)`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(roleId, title, rules)
	if err != nil {
		return err
	}
	return nil
}

// 删除角色
func (this *CoreAuth) DeleteRole(id int) error {
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `delete from ` + auth_role_name + ` where id=?`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	this.clear()
	return nil
}

// 修改角色
func (this *CoreAuth) EditRole(id, roleId int, title, rules string) error {
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `update ` + auth_role_name + ` set title=?,rules=?,role_id=? where id=?`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(title, rules, roleId, id)
	if err != nil {
		return err
	}
	this.clear()
	return nil
}

// 分配角色
func (this *CoreAuth) GiveUserRole(userId, roleId, roomId int) error {
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `insert into ` + auth_role_access_name + ` (user_id,role_id,roomId) values(?,?,?)`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userId, roleId, roomId)
	if err != nil {
		return err
	}
	this.clear()
	return nil
}

// 删除用户角色
func (this *CoreAuth) DeleteUserRole(userId, roleId, roomId int) error {
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `delete from ` + auth_role_access_name + ` where user_id=? and role_id=? and room_id=?`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userId, roleId, roomId)
	if err != nil {
		return err
	}
	this.clear()
	return nil
}

// 获取角色列表
func (this *CoreAuth) ShowRoleList() ([]map[string]any, error) {
	key := "cache_role_list"
	c := this.get(key)
	if c != nil {
		//return c.([]map[string]any), nil
	}
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `select * from ` + auth_role_name + ` where status=1`
	rows, err := db.(*sql.DB).Query(query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	list := []map[string]any{}
	var id int
	var title string
	var status int
	var rules string
	var roleId int
	var icon string
	for rows.Next() {
		if err := rows.Scan(&id, &roleId, &title, &status, &rules, &icon); err != nil {
			return nil, err
		}
		list = append(list, map[string]any{"id": id, "roleId": roleId, "title": title, "status": status, "rules": rules, "icon": icon})
	}
	this.set(key, list)
	return list, nil
}

// 获取角色下的用户
func (this *CoreAuth) ShowRoleUserList(roleId int) ([]map[string]any, error) {
	key := "cache_role_user_list_" + strconv.Itoa(roleId)
	c := this.get(key)
	if c != nil {
		//return c.([]map[string]any), nil
	}
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `select user.id as userId from ` + auth_role_name + ` role inner join ` + auth_role_access_name +
		` access on role.id=access.role_id inner join ` + user_table_name + ` user on access.user_id=user.id where role.id=?`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(roleId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	list := []map[string]any{}
	var userId, name any
	for rows.Next() {
		if err := rows.Scan(&userId, &name); err != nil {
			return nil, err
		}
		list = append(list, map[string]any{"userId": userId, "name": name})
	}
	this.set(key, list)
	return list, nil
}

// 获取角色权限
func (this *CoreAuth) GetRoleRules(roleId int) (map[int][]any, error) {
	key := "cache_role_rule_list_" + strconv.Itoa(roleId)
	c := this.get(key)
	if c != nil {
		//return c.(map[int][]any), nil
	}
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `select rules from ` + auth_role_name + ` where status=1 and role_id=?`
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return nil, err
	}
	roleRows := stmt.QueryRow(roleId)
	var rules string
	err = roleRows.Scan(&rules)
	if err != nil {
		return nil, err
	}
	ruleIds := strings.Split(rules, ",")
	query = `select id,name,title,category,categoryName from ` + auth_rule_name
	ruleRows, err := db.(*sql.DB).Query(query)
	defer ruleRows.Close()
	if err != nil {
		return nil, err
	}
	allRules := []map[string]any{}
	var id, category, ruleId int
	var name, title, categoryName string
	for ruleRows.Next() {
		if err := ruleRows.Scan(&id, &ruleId, &name, &title, &category, &categoryName); err != nil {
			return nil, err
		}
		allRules = append(allRules, map[string]any{"id": id, "ruleId": ruleId, "name": name, "title": title, "category": category, "categoryName": categoryName, "select": 0})
	}
	for k, rule := range allRules {
		if id, ok := rule["id"].(int); ok {
			for _, ruleId := range ruleIds {
				idid, _ := strconv.Atoi(ruleId)
				if idid == id {
					allRules[k]["select"] = 1
				}
			}
		}
	}
	var sortList = map[int][]any{}
	for _, v := range allRules {
		if category, ok := v["category"].(int); ok {
			sortList[category] = append(sortList[category], v)
		}
	}
	this.set(key, sortList)
	return sortList, nil
}

// 验证用户是否拥有权限
func (this *CoreAuth) VerifyAuth(userId, roomId int, isCompere bool, ruleName ...string) (map[string]bool, error) {
	rules, err := this.getUserRules(userId, roomId, isCompere)
	res := map[string]bool{}
	if err != nil {
		return res, err
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

// 获取用户当前房间所有权限
func (this *CoreAuth) getUserRules(userId int, roomId int, isCompere bool) ([]string, error) {
	key := "cache_user_rule_list_" + strconv.Itoa(roomId) + "_" + strconv.Itoa(userId)
	c := this.get(key)
	if c != nil {
		//return c.([]string), nil
	}
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := `select role.rules from ` + auth_role_access_name + ` access inner join ` + auth_role_name + ` role on access.role_id=role.role_id where role.status=1 and access.user_id=? and (access.room_id=? or access.room_id = 0)`
	if !isCompere {
		query = `select role.rules from ` + auth_role_access_name + ` access inner join ` + auth_role_name + ` role on access.role_id=role.role_id where role.status=1 and access.user_id=? and access.role_id!=1005 and (access.room_id=? or access.room_id = 0)`
	}
	stms, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stms.Query(userId, roomId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var ruleStr string
	ruleSlice := []string{}
	for rows.Next() {
		if err := rows.Scan(&ruleStr); err != nil {
			return nil, err
		}
		ruleSlice = append(ruleSlice, strings.Split(ruleStr, ",")...)
	}
	if len(ruleSlice) == 0 { //没有特殊身份则是普通身份
		query = `select rules from ` + auth_role_name + ` where role_id=?`
		stms, err = db.(*sql.DB).Prepare(query)
		if err != nil {
			return nil, err
		}
		rows, err = stms.Query(1010)
		defer rows.Close()
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			if err = rows.Scan(&ruleStr); err != nil {
				return nil, err
			}
			ruleSlice = append(ruleSlice, strings.Split(ruleStr, ",")...)
		}
	}
	var ruleRows *sql.Rows
	query = `select name from ` + auth_rule_name + ` where rule_id in(` + strings.Join(ruleSlice, ",") + `)`
	ruleRows, err = db.(*sql.DB).Query(query)
	if err != nil {
		return nil, err
	}
	defer ruleRows.Close()
	list := []string{}
	for ruleRows.Next() {
		ruleRows.Scan(&ruleStr)
		list = append(list, ruleStr)
	}
	this.set(key, list)
	return list, err
}

func (this *CoreAuth) GetRoleListByRoomIdAndUserId(roomId, userId string) ([]int, error) {
	key := fmt.Sprintf("cache_user_roles_%s_%s", roomId, userId)
	c := this.get(key)
	if c != nil {
		//return c.([]int), nil
	}
	db, err := this.p.Get()
	defer this.p.Put(db)
	query := fmt.Sprintf("select role_id from %s where room_id = %s and user_id = %s", auth_role_access_name, "?", "?")
	stmt, err := db.(*sql.DB).Prepare(query)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(roomId, userId)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	list := []int{}
	var roleId int
	for rows.Next() {
		if err := rows.Scan(&roleId); err != nil {
			return nil, err
		}
		list = append(list, roleId)
	}
	this.set(key, list)
	return list, nil
}

type cacheData struct {
	expires int64
	data    interface{}
}

func (this *CoreAuth) get(key string) any {
	this.l.Lock()
	defer this.l.Unlock()
	data := this.cache[key]
	if time.Now().Unix() > data.expires {
		delete(this.cache, key)
		return nil
	}
	return data.data
}

func (this *CoreAuth) set(key string, value any) {
	this.l.Lock()
	defer this.l.Unlock()
	data := cacheData{
		expires: time.Now().Unix() + 120,
		data:    value,
	}
	this.cache[key] = data
}

func (this *CoreAuth) del(key string) {
	this.l.Lock()
	defer this.l.Unlock()
	delete(this.cache, key)
}

func (this *CoreAuth) checkCache() {
	for {
		this.l.Lock()
		for key, value := range this.cache {
			if time.Now().Unix() > value.expires {
				delete(this.cache, key)
			}
		}
		this.l.Unlock()
		time.Sleep(time.Second * 300)
	}
}

func (this *CoreAuth) clear() {
	this.l.Lock()
	defer this.l.Unlock()
	this.cache = make(map[string]cacheData)
}
