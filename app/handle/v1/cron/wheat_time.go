package v1_cron

import (
	"net/http"
	"yfapi/internal/logic"

	"github.com/gin-gonic/gin"
)

// @Summary
// @Description 房间定时任务处理跨天直播数据
// @Accept json
// @Produce json
// @Router /v1/cron/wheatTime [post]
func WheatTimeCron(c *gin.Context) {
	taskLogic := new(logic.TaskCron)
	taskLogic.WheatTimeCron()
	c.JSON(http.StatusOK, true)
}
