package helper

import (
	"fmt"
	"testing"
)

func TestRemovePrefixImgUrl(t *testing.T) {
	url := `http://test.yungoubuy.com/uploads/20240717/54f3c88659dc458b9ed0ffd2b86c9d42.png`
	imgUrl := RemovePrefixImgUrl(url)
	fmt.Println(imgUrl)
}

func TestPrivateMobile(t *testing.T) {
	mobile := PrivateMobile("136823010377")
	fmt.Println(mobile)
}

func TestCheckPasswordLever(t *testing.T) {
	err := CheckPasswordLever("wdne4312234321")
	t.Log(err)
	err = CheckPasswordLever("1234567")
	t.Log(err)
	err = CheckPasswordLever("aaa1234567812154122112541")
	t.Log(err)
	err = CheckPasswordLever("abcdefghijklmnop")
	t.Log(err)
	err = CheckPasswordLever("abc123&^%$#%&")
	t.Log(err)
}

func TestPrivateRealName(t *testing.T) {
	fmt.Println(PrivateRealName("六小龄童"))
	fmt.Println(PrivateRealName("斯琴高娃"))
	fmt.Println(PrivateRealName("刘亦菲"))
	fmt.Println(PrivateRealName("王宝强"))
	fmt.Println(PrivateRealName("汪峰"))
	fmt.Println(PrivateRealName("伍佰"))
	fmt.Println(PrivateRealName("jeck"))
	fmt.Println(PrivateRealName("jeck lis"))
}

func TestPrivateIdNo(t *testing.T) {
	fmt.Println(PrivateIdNo("411122199005100098"))
	fmt.Println(PrivateIdNo("123456789"))
	fmt.Println(PrivateIdNo("1234567"))
	fmt.Println(PrivateIdNo("123456"))
	fmt.Println(PrivateIdNo("12345"))
	fmt.Println(PrivateIdNo("1234"))
	fmt.Println(PrivateIdNo("123"))
	fmt.Println(PrivateIdNo("12"))
	fmt.Println(PrivateIdNo("1"))
	fmt.Println(PrivateIdNo(""))
}
