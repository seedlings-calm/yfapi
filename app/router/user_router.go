package router

import (
	"sync"
	v1_user "yfapi/app/handle/v1/user"
	"yfapi/app/middle"
	"yfapi/internal/engine"

	"github.com/gin-gonic/gin"
)

func SetLUserRouter(oriEngine *engine.OriEngine, router *gin.RouterGroup, wg *sync.WaitGroup) {
	group := router.Group("/v1/user")
	group.Use(middle.CheckHeaderData())
	group.Use(middle.Auth())
	{
		group.POST("/perfectInfo", BaseRouter(oriEngine, wg, v1_user.PerfectInfo))               // 完善用户信息
		group.GET("/imserver", BaseRouter(oriEngine, wg, v1_user.ImServer))                      // 获取im服务器
		group.GET("/checkRepeatName", BaseRouter(oriEngine, wg, v1_user.CheckRepeatName))        // 检测重复用户名
		group.POST("/editUserInfo", BaseRouter(oriEngine, wg, v1_user.EditUserInfo))             // 更新用户信息
		group.GET("/userInfo", BaseRouter(oriEngine, wg, v1_user.GetUserInfo))                   // 查询用户主页信息
		group.GET("/search", BaseRouter(oriEngine, wg, v1_user.SearchUserInfo))                  // 搜索用户信息
		group.POST("/setPassword", BaseRouter(oriEngine, wg, v1_user.SetPassword))               // 设置密码
		group.POST("/realName", BaseRouter(oriEngine, wg, v1_user.RealName))                     // 实名认证
		group.GET("/realNameStatus", BaseRouter(oriEngine, wg, v1_user.RealNameStatus))          // 实名认证状态
		group.GET("/userAccounts", BaseRouter(oriEngine, wg, v1_user.GetUserAccounts))           // 用户手机号关联账号
		group.POST("/checkCreateAccount", BaseRouter(oriEngine, wg, v1_user.CheckCreateAccount)) // 创建账号检测
		group.POST("/createAccount", BaseRouter(oriEngine, wg, v1_user.CreateNewAccount))        // 创建账号
		group.POST("/switchAccount", BaseRouter(oriEngine, wg, v1_user.SwitchAccount))           // 选择账号
		group.GET("/privateInfo", BaseRouter(oriEngine, wg, v1_user.GetUserPrivateInfo))         // 隐私资料

		group.POST("/sendCode", BaseRouter(oriEngine, wg, v1_user.SendVerifyCode))   // 发送手机验证码
		group.POST("/verifyCode", BaseRouter(oriEngine, wg, v1_user.VerifyMobile))   // 校验手机
		group.POST("/changeMobile", BaseRouter(oriEngine, wg, v1_user.ChangeMobile)) // 变更手机号
		group.POST("/noticeFilter", BaseRouter(oriEngine, wg, v1_user.NoticeFilter)) //通知规则
		group.GET("/loginRecord", BaseRouter(oriEngine, wg, v1_user.GetLoginLog))

		group.GET("/privacySetting", BaseRouter(oriEngine, wg, v1_user.GetPrivacySetting)) //用户隐私设置信息

		group.GET("/userInfoByUserNo", BaseRouter(oriEngine, wg, v1_user.SearchUserInfoByUserNo)) // 查询用户信息
		group.GET("/userRealNameInfo", BaseRouter(oriEngine, wg, v1_user.GetUserRealNameInfo))    // 查询用户实名认证信息
	}
	{
		group.GET("/account", BaseRouter(oriEngine, wg, v1_user.GetUserAccountInfo))                 // 用户账户信息
		group.GET("/bill/diamond", BaseRouter(oriEngine, wg, v1_user.GetUserDiamondBill))            // 用户钻石流水
		group.GET("/bill/starlight", BaseRouter(oriEngine, wg, v1_user.GetUserStarlightBill))        // 用户星光流水
		group.POST("/exchangeDiamond", BaseRouter(oriEngine, wg, v1_user.StarlightExchangeDiamond))  // 用户星光兑换钻石
		group.GET("/rechargeDiamond", BaseRouter(oriEngine, wg, v1_user.RechargeDiamond))            // 用户充值钻石页面
		group.GET("/bill/rechargeDiamondLog", BaseRouter(oriEngine, wg, v1_user.RechargeDiamondLog)) // 充值钻石日志

	}
	{
		//关注用户
		group.POST("/addFollow", BaseRouter(oriEngine, wg, v1_user.AddFollow))
		//取消关注用户
		group.POST("/removeFollow", BaseRouter(oriEngine, wg, v1_user.RemoveFollow))
		//获取关注列表
		group.POST("/getUserFollowingList", BaseRouter(oriEngine, wg, v1_user.GetUserFollowingList))
		//获取粉丝列表
		group.POST("/getFollowersList", BaseRouter(oriEngine, wg, v1_user.GetFollowersList))
		//获取好友列表
		group.POST("/getFriendsList", BaseRouter(oriEngine, wg, v1_user.GetFriendsList))
		//删除粉丝
		group.POST("/deleteFans", BaseRouter(oriEngine, wg, v1_user.DeleteFans))
	}
	{
		// 动态相关
		group.GET("/timeline/detail", BaseRouter(oriEngine, wg, v1_user.GetTimelineDetail))                      // 动态详情信息
		group.GET("/timeline/category/list", BaseRouter(oriEngine, wg, v1_user.GetTimelineListByType))           // 动态分类列表
		group.POST("/timeline/publish", BaseRouter(oriEngine, wg, v1_user.PublishTimeline))                      // 发布动态
		group.GET("/timeline/list", BaseRouter(oriEngine, wg, v1_user.GetUserTimelineList))                      // 动态列表
		group.GET("/timeline/delete", BaseRouter(oriEngine, wg, v1_user.DeleteTimeline))                         // 删除动态
		group.GET("/timeline/praise", BaseRouter(oriEngine, wg, v1_user.PraiseTimeline))                         // 动态点赞
		group.GET("/timeline/praise/cancel", BaseRouter(oriEngine, wg, v1_user.CancelPraiseTimeline))            // 动态取消点赞
		group.GET("/timeline/praise/list", BaseRouter(oriEngine, wg, v1_user.GetPraiseUserList))                 // 动态点赞列表
		group.POST("/timeline/reply", BaseRouter(oriEngine, wg, v1_user.ReplyTimeline))                          // 动态评论
		group.GET("/timeline/reply/delete", BaseRouter(oriEngine, wg, v1_user.DeleteReplyTimeline))              // 删除动态评论
		group.GET("/timeline/reply/list", BaseRouter(oriEngine, wg, v1_user.GetTimelineReplyList))               // 动态评论列表
		group.GET("/timeline/reply/praise", BaseRouter(oriEngine, wg, v1_user.PraiseTimelineReply))              // 动态评论点赞
		group.GET("/timeline/reply/praise/cancel", BaseRouter(oriEngine, wg, v1_user.CancelPraiseTimelineReply)) // 动态评论取消点赞
		group.GET("/timeline/reply/sub/list", BaseRouter(oriEngine, wg, v1_user.GetTimelineSubReplyList))        // 动态评论回复列表
		group.POST("/timeline/filter", BaseRouter(oriEngine, wg, v1_user.TimelineFilter))                        //动态规则
		group.GET("/timeline/filterList", BaseRouter(oriEngine, wg, v1_user.TimelineFilterList))                 //动态规则列表
	}
	{
		//收藏相关
		group.GET("/collect", middle.RateLimitMiddleware(), BaseRouter(oriEngine, wg, v1_user.AddCollect)) //添加收藏
		group.PUT("/collect", middle.RateLimitMiddleware(), BaseRouter(oriEngine, wg, v1_user.DelCollect)) //取消收藏
	}
	{
		// 访客足迹
		group.GET("/visitClear", middle.RateLimitMiddleware(), BaseRouter(oriEngine, wg, v1_user.ClearUserVisitRecord)) //清除我的足迹
		group.GET("/myVisit", BaseRouter(oriEngine, wg, v1_user.GetUserVisitRecordList))                                // 我的足迹
		group.GET("/visitMe", BaseRouter(oriEngine, wg, v1_user.GetVisitUserRecordList))                                // 访客记录
	}

	{
		//提现
		group.GET("/appSettingInfo", BaseRouter(oriEngine, wg, v1_user.AppSettingInfo))    //获取提现详情
		group.POST("/appSettingApply", BaseRouter(oriEngine, wg, v1_user.AppSettingApply)) //提现申请
		group.POST("/bankAdd", BaseRouter(oriEngine, wg, v1_user.BankAdd))                 //绑定银行卡
		group.POST("/bankUnBind", BaseRouter(oriEngine, wg, v1_user.BankUnBind))           //解绑银行卡
		group.GET("/bankList", BaseRouter(oriEngine, wg, v1_user.BankList))                //获取银行卡列表

	}
}
