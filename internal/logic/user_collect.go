package logic

import (
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"

	"github.com/gin-gonic/gin"
)

type UserCollect struct {
}

// 添加
func (d *UserCollect) AddCollect(roomId string, c *gin.Context) error {
	userId := handle.GetUserId(c)
	userDao := dao.DaoUserCollect{
		UserId: userId,
	}
	roomDao := dao.RoomDao{}
	roomInfo, err := roomDao.FindOne(&model.Room{
		Id: roomId,
	})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeCollect,
			Msg:  nil,
		})
	}
	if roomInfo.Status == 3 {
		panic(error2.I18nError{
			Code: error2.ErrCodeCollect,
			Msg:  nil,
		})
	}
	err = userDao.Create(roomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeCollect,
			Msg:  nil,
		})
	}
	return nil
}

// 取消
func (d *UserCollect) DelCollect(roomId string, c *gin.Context) error {
	userId := handle.GetUserId(c)
	dao := dao.DaoUserCollect{
		UserId: userId,
	}

	err := dao.Update(roomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeDelCollect,
			Msg:  nil,
		})
	}
	return nil
}
