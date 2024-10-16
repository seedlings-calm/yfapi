package v1_cron

import (
	"fmt"
	"net/http"
	"time"
	"yfapi/core/coreLog"

	"github.com/gin-gonic/gin"
)

func TestTask(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			coreLog.Error("cron run err name:%s err:%+v", "TestTask", err)
			c.JSON(http.StatusInternalServerError, err)
		}
	}()
	fmt.Println("执行定时任务成功", time.Now().Format(time.DateTime))
	c.JSON(http.StatusOK, true)
}
