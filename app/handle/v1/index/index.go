package index

import (
	"yfapi/app/handle"
	i18n_msg "yfapi/i18n/msg"
	"yfapi/internal/logic"
	request_index "yfapi/typedef/request/index"
	"yfapi/typedef/response"
	response_index "yfapi/typedef/response/index"
	response_room "yfapi/typedef/response/room"
	_ "yfapi/typedef/response/user"

	"github.com/gin-gonic/gin"
)

type Nav struct {
	Id       string
	Name     string
	ParentId string
	Sort     int
	Children []Nav
}

func GetNav(c *gin.Context) []Nav {
	return []Nav{
		{
			Name:     i18n_msg.GetI18nMsg(c, i18n_msg.CollectKey),
			Id:       "1",
			Sort:     1,
			ParentId: "0",
			Children: nil,
		},
		{
			Name:     i18n_msg.GetI18nMsg(c, i18n_msg.RecommendedKey),
			Id:       "2",
			Sort:     2,
			ParentId: "0",
			Children: nil,
		},
		{
			Name:     i18n_msg.GetI18nMsg(c, i18n_msg.EmotionKey),
			Id:       "3",
			Sort:     3,
			ParentId: "0",
			Children: []Nav{
				{Name: i18n_msg.GetI18nMsg(c, i18n_msg.EmotionManKey), Id: "101", Sort: 1, ParentId: "3", Children: nil},
				{Name: i18n_msg.GetI18nMsg(c, i18n_msg.EmotionWomanKey), Id: "102", Sort: 2, ParentId: "3", Children: nil},
			},
		},
		{
			Name:     i18n_msg.GetI18nMsg(c, i18n_msg.LiveKey),
			Id:       "4",
			Sort:     4,
			ParentId: "0",
			Children: []Nav{
				{Name: i18n_msg.GetI18nMsg(c, i18n_msg.LiveVoiceKey), Id: "201", Sort: 1, ParentId: "4", Children: nil},
				{Name: i18n_msg.GetI18nMsg(c, i18n_msg.LiveVideoKey), Id: "202", Sort: 2, ParentId: "4", Children: nil},
			},
		},
		{
			Name:     i18n_msg.GetI18nMsg(c, i18n_msg.SingKey),
			Id:       "5",
			Sort:     5,
			ParentId: "0",
			Children: nil,
		},
	}
}

// @Summary 获取首页导航栏目
// @Description
// @Tags 首页
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Router /v1/index/getNavigation [get]
func GetNavigation(c *gin.Context) {
	response.SuccessResponse(c, GetNav(c))
}

// @Summary 收藏列表
// @Description
// @Tags 首页
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} user.CollectResponse{}
// @Router /v1/index/getCollect [get]
func GetCollect(c *gin.Context) {
	service := new(logic.Indexs)
	resp := service.GetCollect(c)
	response.SuccessResponse(c, resp)
}

// @Summary 导航栏 推荐 使用
// @Description 取后台配置数据
// @Tags 首页
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} user.RecommendResponse{}
// @Router /v1/index/getRecommend [get]
func GetRecommend(c *gin.Context) {
	service := new(logic.Indexs)
	resp := service.GetRecommend(c)
	response.SuccessResponse(c, resp)
}

// @Summary App元氛TOP榜（pc-热门列表）
// @Description
// @Tags 首页
// @Accept json
// @Produce json
// @Param	req	query	request_index.TopListReq	true	"房间列表参数"
// @Success		200 {object} response.BasePageRes
// @Router /v1/index/top [get]
func GetTopRooms(c *gin.Context) {
	req := new(request_index.TopListReq)
	handle.BindQuery(c, req)
	service := new(logic.Room)
	resp := service.GetTop(c, req)
	response.SuccessResponse(c, resp)
}

// @Summary APP菜单配置列表
// @Description
// @Tags 首页
// @Accept json
// @Produce json
// @Param	req	query	request_index.AppMenuSettingReq	true	"房间列表参数"
// @Success		200 {object} []index.AppMenuSetting
// @Router /v1/index/menuList [get]
func GetAppMenuSettingList(c *gin.Context) {
	req := new(request_index.AppMenuSettingReq)
	handle.BindQuery(c, req)
	service := new(logic.AppMenuSetting)
	resp := service.GetAppMenuSettingList(c, req)
	response.SuccessResponse(c, resp)
}

// @Summary pc端 我的房间列表
// @Description
// @Tags 首页
// @Accept json
// @Produce json
// @Success 0 {object} response.Response{}
// @Success 0 {object} []response_room.RoomInfo{}
// @Router /v1/index/getRoomsByPC [get]
func GetRoomsByPC(c *gin.Context) {
	service := new(logic.Indexs)
	resp := make([]*response_room.RoomInfo, 0)
	resp = append(resp, service.GetRoomsByPC(c)...)

	response.SuccessResponse(c, resp)
}

// @Summary 头条中心
// @Description
// @Tags 首页
// @Accept json
// @Produce json
// @Param req query  request_index.TopMsgReq  true "请求参数"
// @Success 0 {object} response.Response{}
// @Router /v1/index/topMsg [get]
func TopMsg(c *gin.Context) {
	req := new(request_index.TopMsgReq)
	handle.BindQuery(c, req)
	service := new(logic.Indexs)
	resp := make([]response_index.TopMsgRes, 0)
	resp = append(resp, service.GetTopMsg(c, req)...)

	response.SuccessResponse(c, resp)
}
