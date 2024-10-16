package logic

import (
	"github.com/gin-gonic/gin"
	"yfapi/internal/service/acl"
	"yfapi/typedef/request/room"
)

type RoomAuth struct {
}

func (r *RoomAuth) GetAuthMenu(c *gin.Context, userId string, req *room.RoomAuthMenuReq) any {
	aclInstance := &acl.RoomAcl{
		UserId:       userId,
		TargetUserId: req.TargetUserId,
		RoomId:       req.RoomId,
		Seat:         req.Seat,
		Scene:        req.Scene,
	}
	return aclInstance.GetAcl()
}
