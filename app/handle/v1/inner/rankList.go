package inner

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"yfapi/internal/service/rankList"
)

// 人气榜，贡献榜日结
func RankListDay(c *gin.Context) {
	ser := rankList.Instance()
	//获取贡献榜前三名
	contributorUsers := ser.GetRankListSettle(3, "", rankList.CycleDay, rankList.RankListContributor)
	fmt.Println(contributorUsers)
	//TODO 发放礼物

	//人气榜前三
	popularityUsers := ser.GetRankListSettle(3, "", rankList.CycleDay, rankList.RankListPopularity)
	fmt.Println(popularityUsers)
	//TODO 发放礼物
}

// 人气榜，贡献榜周结
func RankListWeek(c *gin.Context) {
	ser := rankList.Instance()
	//获取贡献榜前三名
	contributorUsers := ser.GetRankListSettle(3, "", rankList.CycleWeek, rankList.RankListContributor)
	fmt.Println(contributorUsers)
	//TODO 发放礼物

	//人气榜前三
	popularityUsers := ser.GetRankListSettle(3, "", rankList.CycleWeek, rankList.RankListPopularity)
	fmt.Println(popularityUsers)
	//TODO 发放礼物
}

// 人气榜，贡献榜月结
func RankListMonth(c *gin.Context) {
	ser := rankList.Instance()
	//获取贡献榜前三名
	contributorUsers := ser.GetRankListSettle(3, "", rankList.CycleMonth, rankList.RankListContributor)
	fmt.Println(contributorUsers)
	//TODO 发放礼物

	//人气榜前三
	popularityUsers := ser.GetRankListSettle(3, "", rankList.CycleMonth, rankList.RankListPopularity)
	fmt.Println(popularityUsers)
	//TODO 发放礼物
}
