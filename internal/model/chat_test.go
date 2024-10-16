package model

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
	"yfapi/core/coreSnowflake"
)

func newPg() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", "192.168.77.107", "postgres", "weisheng@123#pg", "v_chat", "5432")
	open, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return open, err
}

func TestWritePublicChat(t *testing.T) {
	pg, err := newPg()
	if err != nil {
		t.Error(err)
		return
	}
	node, _ := coreSnowflake.New(1)
	fromUserId := node.Generate().String()
	toUserId := node.Generate().String()
	roomId := node.Generate().String()
	message := `{"code":1001,"msg":"公屏测试消息"}`
	for i := 0; i < 20; i++ {
		publicRecord := PublicChat{
			MessageId:  node.Generate().String(),
			FromUserId: fromUserId,
			ToUserId:   toUserId,
			RoomId:     roomId,
			Message:    message,
			Timestamp:  time.Now().UnixMicro(),
		}
		pg.Create(&publicRecord)
	}
}

func TestWritePrivateChat(t *testing.T) {
	pg, err := newPg()
	if err != nil {
		t.Error(err)
		return
	}
	node, _ := coreSnowflake.New(1)
	fromUserId := node.Generate().String()
	toUserId := node.Generate().String()
	message := `{"code":1001,"msg":"公屏测试消息"}`
	for i := 0; i < 20; i++ {
		privateRecord := PrivateChat{
			MessageId:  node.Generate().String(),
			UniteId:    fromUserId + toUserId,
			FromUserId: fromUserId,
			ToUserId:   toUserId,
			Message:    message,
			Read:       2,
			Timestamp:  time.Now().UnixMicro(),
			Type:       1,
		}
		pg.Create(&privateRecord)
	}
}

func TestSelectPrivateChat(t *testing.T) {
	pg, err := newPg()
	if err != nil {
		t.Error(err)
		return
	}
	list := []PrivateChat{}
	err = pg.Model(&PrivateChat{}).Where("unite_id = ? and read = ?", "18076701318835118081807670131883511809", 2).Order("timestamp desc").Find(&list).Error
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(list), list)
}
