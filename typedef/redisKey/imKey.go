package redisKey

import "fmt"

// 用户不同端登录状态
var UserClientLoginStatus = func(userId, clientType string) string {
	return fmt.Sprintf("imOnline:user:%s:%s", clientType, userId)
}

// 用户是否在线
var ImOnlineUser = func(userId string) string {
	return fmt.Sprintf("imOnline:user:%s", userId)
}
