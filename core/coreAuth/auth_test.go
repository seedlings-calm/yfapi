package coreAuth

import (
	"fmt"
	"testing"
)

func getGoAuth() *CoreAuth {
	return New(CoreAuthConfig{
		UserName: "root",
		PassWord: "root",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "v_chat",
	})
}

func TestGoAuth_AddRule(t *testing.T) {
	goAuth := getGoAuth()
	if err := goAuth.AddRule("enter_room_hidden", "隐身进厅", 1050, 0, "测试权限"); err != nil {
		fmt.Println(err)
	}
}

func TestGoAuth_EditRule(t *testing.T) {
	goAuth := getGoAuth()
	if err := goAuth.EditRule(1, "enter_room_hidden", "隐身进入房间", 1, 1, ""); err != nil {
		fmt.Println(err)
	}
}

func TestCoreAuth_DeleteRule(t *testing.T) {
	goAuth := getGoAuth()
	goAuth.DeleteRule(2)
}

func TestGoAuth_AddRole(t *testing.T) {
	goAuth := getGoAuth()
	if err := goAuth.AddRole(1001, "主持人", "1,3"); err != nil {
		t.Error(err)
	}
	if err := goAuth.AddRole(1002, "麦未嘉宾", "1"); err != nil {
		t.Error(err)
	}
	if err := goAuth.AddRole(1003, "歌手", "1"); err != nil {
		t.Error(err)
	}
}

func TestGoAuth_GiveUserRole(t *testing.T) {
	goAuth := getGoAuth()
	if err := goAuth.GiveUserRole(1, 1, 1); err != nil {
		t.Error(err)
	}
	if err := goAuth.GiveUserRole(1, 2, 1); err != nil {
		t.Error(err)
	}
}

func TestGoAuth_ShowRoleList(t *testing.T) {
	goAuth := getGoAuth()
	list, err := goAuth.ShowRoleList()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list)
}

func TestGoAuth_GetRoleRules(t *testing.T) {
	goAuth := getGoAuth()
	list, err := goAuth.GetRoleRules(5)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(list)
}

func TestCoreAuth_VerifyAuth(t *testing.T) {
	goAuth := getGoAuth()
	rules := []string{
		"up_hidden_mic",
		"up_compere_mic",
		"apply_normal_mic",
	}
	auth, err := goAuth.VerifyAuth(1816650045333737472, 1810568840619270144, false, rules...)
	fmt.Println(auth, err)
}
