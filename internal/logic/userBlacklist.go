package logic

import (
	"yfapi/app/handle"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_user "yfapi/internal/service/user"
	response_room "yfapi/typedef/response/room"

	"github.com/gin-gonic/gin"
)

type BlackListLogic struct {
}

func (b BlackListLogic) AddBlacklist(c *gin.Context, toId string) {
	fromId := handle.GetUserId(c)
	blackDao := dao.UserBlackListDao{}
	err := blackDao.Create(&model.UserBlacklist{FromID: fromId, ToID: toId, RoomID: "0", Types: 2, IsEffective: true})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeBlackListAddErr,
			Msg:  nil,
		})
	}
}

func (b BlackListLogic) DelBlacklist(c *gin.Context, toId string) {
	fromId := handle.GetUserId(c)
	blackDao := dao.UserBlackListDao{}
	err := blackDao.Update(&model.UserBlacklist{UnsealID: fromId, FromID: fromId, ToID: toId, RoomID: "0", Types: 2, IsEffective: false})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeBlackListDelErr,
			Msg:  nil,
		})
	}
}

// GetUserBlacklist 获取用户黑名单列表
func (b BlackListLogic) GetUserBlacklist(c *gin.Context) (res []*response_room.BlackListAndUserInfo) {
	userId := handle.GetUserId(c)
	blackDao := dao.UserBlackListDao{}
	res = blackDao.GetUserBlackList(userId)
	for _, info := range res {
		info.Avatar = helper.FormatImgUrl(info.Avatar)
		// 查询用户的铭牌信息
		info.UserPlaque = service_user.GetUserLevelPlaque(info.UserId, helper.GetClientType(c))
	}
	return
}
