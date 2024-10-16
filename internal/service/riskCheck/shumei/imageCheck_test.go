package shumei

import (
	"fmt"
	"testing"
)

func TestShuMei_NicknameCheck(t *testing.T) {
	check := new(ShuMei).NicknameCheck("10001", "我操你妈的")
	fmt.Println(check)
}

func TestShuMei_PrivateChatCheck(t *testing.T) {
	check, ok := new(ShuMei).PrivateChatCheck("10001", "10002", "我是毛泽东")
	fmt.Println(check, ok)
}

func TestShuMei_MomentsCheck(t *testing.T) {
	check := new(ShuMei).MomentsCheck("10001", "习近平习近平毛泽东东")
	fmt.Println(check)
}

func TestShuMei_CommentCheck(t *testing.T) {
	check := new(ShuMei).CommentCheck("10001", "习近平习近平毛泽东东")
	fmt.Println(check)
}

func TestShuMei_SignCheck(t *testing.T) {
	check := new(ShuMei).SignCheck("10001", "习近平习近平毛泽东东")
	fmt.Println(check)
}
