package helper

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"net/url"
	"yfapi/core/coreConfig"
)

func WsFull(server string) string {
	if coreConfig.GetHotConf().ENV == "pro" {
		return fmt.Sprintf("wss://%s", server)
	}
	if coreConfig.GetHotConf().ENV == "test" {
		return fmt.Sprintf("wss://%s", server)
	}
	return fmt.Sprintf("ws://%s:%d", server, 9002)
}

// 获取联合的用户ID
func GetUniteId(userId1, userId2 string) string {
	u1 := cast.ToInt64(userId1)
	u2 := cast.ToInt64(userId2)
	if u2 > u1 {
		return cast.ToString(u1) + cast.ToString(u2)
	}
	return cast.ToString(u2) + cast.ToString(u1)
}

// 检测im服务器是否可用
func CheckWebSocketConnection(urlStr string) (bool, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false, err
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return true, nil
}
