package roomOwner

import (
	"yfapi/internal/helper"
	"yfapi/internal/logic/room"
	"yfapi/typedef/response"
	response_roomowner "yfapi/typedef/response/roomOwner"

	"github.com/gin-gonic/gin"
)

// PersonList
// @Summary 成员列表
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param startTime  query string  false "开始时间"
// @Param endTime  query string  false "结束时间"
// @Param sort  query string  false "排序"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_roomowner.PersonListRes{}
// @Router /v1/roomOwner/member/list [get]
func PersonList(c *gin.Context) {
	pd := &room.PersonDao{
		RoomId: c.GetHeader("roomId"),
		UserId: helper.GetUserId(c),
	}
	var res = make([]*response_roomowner.PersonListRes, 0)
	res = append(res, pd.PersonList(c, c.Query("startTime"), c.Query("endTime"), c.Query("sort"))...)
	response.SuccessResponse(c, res)
}

// PersonListDetail
// @Summary 成员受赏详情
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param userId  query string  true "成员长ID"
// @Param startTime  query string  false "开始时间"
// @Param endTime  query string  false "结束时间"
// @Param sort  query string  false "排序"
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_roomowner.PersonListDetailRes{}
// @Router /v1/roomOwner/member/detail [get]
func PersonListDetail(c *gin.Context) {
	pd := &room.PersonDao{
		RoomId: c.GetHeader("roomId"),
		UserId: helper.GetUserId(c),
	}
	var res = make([]*response_roomowner.PersonListDetailRes, 0)
	res = append(res, pd.PersonListDetail(c, c.Query("userId"), c.Query("startTime"), c.Query("endTime"), c.Query("sort"))...)
	response.SuccessResponse(c, res)

}

// RoomDashBoard
// @Summary 房间数据接口
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} response_roomowner.RoomDashBoardRes{}
// @Router /v1/roomOwner/room/dashboard [get]
func RoomDashBoard(c *gin.Context) {
	pd := &room.PersonDao{
		RoomId: c.GetHeader("roomId"),
		UserId: helper.GetUserId(c),
	}
	response.SuccessResponse(c, pd.RoomDashBoard(c))
}

// RoomDashBoardMoneysChart
// @Summary 房间流水chart
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param timeType query string  true "参数:day,month"
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_roomowner.RoomDashBoardChartRes{}
// @Router /v1/roomOwner/room/dashboardMoney [get]
func RoomDashBoardMoneysChart(c *gin.Context) {
	pd := &room.PersonDao{
		RoomId: c.GetHeader("roomId"),
		UserId: helper.GetUserId(c),
	}
	response.SuccessResponse(c, pd.RoomDashBoardMoneysChart(c))
}

// RoomDashBoardTimesChart
// @Summary 有效观看数据
// @Description
// @Tags 房主后台
// @Accept json
// @Produce json
// @Param timeType query string  true "参数:day,month"
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_roomowner.RoomDashBoardChartRes{}
// @Router /v1/roomOwner/room/dashboardTimes [get]
func RoomDashBoardTimesChart(c *gin.Context) {
	pd := &room.PersonDao{
		RoomId: c.GetHeader("roomId"),
		UserId: helper.GetUserId(c),
	}
	response.SuccessResponse(c, pd.RoomDashBoardTimesChart(c))
}
