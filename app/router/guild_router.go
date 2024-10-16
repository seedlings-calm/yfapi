package router

import (
	"sync"
	v1_guild "yfapi/app/handle/v1/guild"
	v1_resource "yfapi/app/handle/v1/resource"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

// 公会后台的专属路由
// SetGuildRouter 公会后台路由
func SetGuildRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/guild")
	//手机验证码登录
	group.POST("/loginAction", BaseRouter(oriEngine, wg, v1_guild.SmsCodeLogin))
	//发送短信
	group.POST("/sendSms", BaseRouter(oriEngine, wg, v1_guild.GetLoginMobileSMS))
	//国家区号
	//group.GET("/regionCode", BaseRouter(oriEngine, wg, v1_login.RegionCodeList))
	group.Use(middle.GuildAuth())
	{
		group.GET("/getUploadToken", BaseRouter(oriEngine, wg, v1_resource.GetUploadTokenToGuild)) // 获取上传token
		// 首页相关
		group.POST("/guildInfo", BaseRouter(oriEngine, wg, v1_guild.GuildInfo))          // 公会信息
		group.GET("/statInfo", BaseRouter(oriEngine, wg, v1_guild.GetGuildStatInfo))     // 首页公会统计信息
		group.GET("/profitInfo", BaseRouter(oriEngine, wg, v1_guild.GetGuildProfitInfo)) // 首页公会流水信息
		group.GET("/roomRank", BaseRouter(oriEngine, wg, v1_guild.GetRoomRankList))      // 首页公会房间排行榜
		// 账户信息
		group.GET("/accountInfo", BaseRouter(oriEngine, wg, v1_guild.AccountInfo))                // 账户信息
		group.POST("/accountBill", BaseRouter(oriEngine, wg, v1_guild.GetAccountBillList))        // 账户交易明细列表
		group.POST("/exchange", BaseRouter(oriEngine, wg, v1_guild.ExchangeDiamond))              // 兑换钻石
		group.POST("/guildBindBank", BaseRouter(oriEngine, wg, v1_guild.GuildBindBank))           // 会长绑定银行卡
		group.POST("/guildWithdrawApply", BaseRouter(oriEngine, wg, v1_guild.GuildWithdrawApply)) // 公会提现申请

		// 成员相关
		group.POST("/memberList", BaseRouter(oriEngine, wg, v1_guild.MemberList))                                    // 获取公会成员列表
		group.GET("/memberIdcards", BaseRouter(oriEngine, wg, v1_guild.MemberIdcards))                               // 公会成员身份列表
		group.POST("/membershipList", BaseRouter(oriEngine, wg, v1_guild.MembershipList))                            // 入会申请列表
		group.GET("/practitionerAction", BaseRouter(oriEngine, wg, v1_guild.GetPractitionerActionRecord))            // 从业者行为记录
		group.POST("/guildMemberApplyReview", BaseRouter(oriEngine, wg, v1_guild.MemberApplyReview))                 // 入会申请审核
		group.POST("/withdrawList", BaseRouter(oriEngine, wg, v1_guild.WithdrawMembershipList))                      // 退会申请列表
		group.POST("/guildMemberWithdrawApplyReview", BaseRouter(oriEngine, wg, v1_guild.MemberWithdrawApplyReview)) // 退会申请审核

		// 获取公会成员分组
		group.POST("/memberGroup", BaseRouter(oriEngine, wg, v1_guild.MemberGroup))             // 公会分组列表
		group.POST("/savememberGroup", BaseRouter(oriEngine, wg, v1_guild.SaveMemberGroup))     // 创建分组
		group.POST("/memberGroupUpdate", BaseRouter(oriEngine, wg, v1_guild.MemberGroupUpdate)) // 更新分组信息
		group.GET("/groupDelete", BaseRouter(oriEngine, wg, v1_guild.MemberGroupDelete))        // 删除分组
		group.GET("/kickOutMember", BaseRouter(oriEngine, wg, v1_guild.GuildKickoutMember))     // 踢出公会成员
		group.POST("/setGroup", BaseRouter(oriEngine, wg, v1_guild.SetGroupByMembers))          // 设置分组

		// 公会房间
		group.POST("/chatRoomType", BaseRouter(oriEngine, wg, v1_guild.RoomTypeList))               // 房间类型
		group.POST("/guildRoomApply", BaseRouter(oriEngine, wg, v1_guild.RoomApply))                // 公会房间申请
		group.POST("/guildRoomApplyList", BaseRouter(oriEngine, wg, v1_guild.RoomApplyList))        // 公会房间申请列表
		group.POST("/roomList", BaseRouter(oriEngine, wg, v1_guild.RoomList))                       // 房间列表
		group.POST("/updateRoom", BaseRouter(oriEngine, wg, v1_guild.UpdateRoom))                   // 更换房主ID\日结算\月结算ID
		group.POST("/getUserBaseInfo", BaseRouter(oriEngine, wg, v1_guild.GetUserBaseInfoByUserNo)) // 通过userNo返回用户的基本信息
		group.POST("/closeRoom", BaseRouter(oriEngine, wg, v1_guild.CloseRoom))                     // 关闭房间

		//group.POST("/roomBase", BaseRouter(oriEngine, wg, v1_guild.MemberList))   //厅下拉
		//group.POST("/roomOnline", BaseRouter(oriEngine, wg, v1_guild.MemberList)) //房间在线时长统计列表

		// 公会流水
		group.POST("/roomProfit", BaseRouter(oriEngine, wg, v1_guild.GetGuildRoomProfitList))     // 公会房间流水
		group.POST("/memberProfit", BaseRouter(oriEngine, wg, v1_guild.GetGuildMemberProfitList)) // 公会成员礼物收入
		group.POST("/guildReward", BaseRouter(oriEngine, wg, v1_guild.GetGuildRewardList))        // 公会礼物打赏详情列表
		//group.POST("/totalReward", BaseRouter(oriEngine, wg, v1_guild.MemberList))                // 公会流水汇总
	}
}
