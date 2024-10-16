package rankList

import (
	"github.com/gin-gonic/gin"
	"yfapi/app/handle"
	"yfapi/internal/logic"
	"yfapi/typedef/request/rankList"
	"yfapi/typedef/response"
)

// 排行榜
func RankList(c *gin.Context) {
	req := &rankList.RankListReq{}
	handle.BindQuery(c, req)
	logicSer := &logic.RankListLogic{}
	response.SuccessResponse(c, logicSer.GetRankList(c, req))
}
