package redisKey

import (
	"fmt"
	"time"
)

var (
	// LockUserRegister 用户注册锁
	LockUserRegister = func(mobile, regionCode string) string {
		return fmt.Sprintf("userRegisterLock:%v:%v", regionCode, mobile)
	}
	// UserLoginInfo 用户登录信息
	UserLoginInfo = func(platform, userId string) string {
		return fmt.Sprintf("userLoginInfo:%v:%v", platform, userId)
	}
	// GuildAdminLoginTime 公会后台登录时间记录
	GuildAdminLoginTime = func(guildId string) string {
		return fmt.Sprintf("GuildAdminLoginTime:%v", guildId)
	}
	// RoomAdminLoginTime 房主后台登录时间记录
	RoomAdminLoginTime = func(userId string) string {
		return fmt.Sprintf("RoomAdminLoginTime:%v", userId)
	}
	// RepeatSubmit 重复提交
	RepeatSubmit = func(key string) string {
		return fmt.Sprintf("repeatSubmit:%v", key)
	}
	// UserImServer 用户im服务
	UserImServer = func(userId string) string {
		return fmt.Sprintf("user:connectImService:%s", userId)
	}
	UserImServerLock = func(userId string) string {
		return fmt.Sprintf("lock:user:connectImService:%s", userId)
	}
	OnlineImServer = func() string {
		return "im:onlineServices"
	}
	// ImUserInfo im用户信息
	ImUserInfo = func(userId string) string {
		return fmt.Sprintf("im:userinfo:%v", userId)
	}
	// ImOneSessionSortId 用户会话列表排序
	ImOneSessionSortId = func(userId string) string {
		return fmt.Sprintf("im:one:session:id:%v", userId)
	}
	// ImOneSessionList im单聊会话列表
	ImOneSessionList = func(userId string) string {
		return fmt.Sprintf("im:one:session:list:%v", userId)
	}
	// ImOneMsgNotReadNum im单聊消息未读数
	ImOneMsgNotReadNum = func(userId string) string {
		return fmt.Sprintf("im:one:msg:num:%v", userId)
	}
	OssUploadPhotoStsToken = func() string {
		return fmt.Sprintf("OssUpload:StsToken")
	}
	// TimelinePraisedUserList 动态点赞列表
	TimelinePraisedUserList = func(timelineId int64) string {
		return fmt.Sprintf("TimelinePraisedUserList:%v", timelineId)
	}
	// TimelineReplyPraisedUserList 动态评论点赞列表
	TimelineReplyPraisedUserList = func(replyId int64) string {
		return fmt.Sprintf("TimelineReplyPraisedUserList:%v", replyId)
	}

	//用户修改信息锁
	UserModifyInfoLock = func(userId string) string {
		return fmt.Sprintf("userModifyInfoLock:%v", userId)
	}

	//验证码key
	UserSmsCode = func(area, mobile string, types int) string {
		return fmt.Sprintf("SmsCode:%s:%d", area+mobile, types)
	}
	//用户昵称修改次数
	UserNicknameEditNum = func(userId string) string {
		return fmt.Sprintf("userNicknameEditNum:%s:%v", time.Now().Month().String(), userId)
	}

	//用户声音修改次数
	UserVoiceEditNum = func(userId string) string {
		return fmt.Sprintf("userVoiceEditNum:%s:%v", time.Now().Month().String(), userId)
	}

	//创建用户No锁
	CreateUserNoLock = func() string {
		return fmt.Sprintf("createUserNoLock")
	}

	//随机生成userNickName锁
	CreateUserNickName = func() string {
		return fmt.Sprintf("createUserNickNameLock")
	}
	//从业者考核数据存储
	UserPractitionerQuestion = func(userId, types string) string {
		return fmt.Sprintf("questionAnswer:%s:%s", userId, types)
	}

	//从业者基础考试次数存储 一天上限3次
	UserPractitionerQuestionNums = func(userId, types string) string {
		return fmt.Sprintf("questionAnswerNums:%s:%s", userId, types)
	}

	//用户所在房间
	UserInWhichRoom = func(userId, clientType string) string {
		return fmt.Sprintf("userInWhithRoom:%s:%s", clientType, userId)
	}

	//用户权限缓存
	UserRules = func(userId, roomId string) string {
		return fmt.Sprintf("userAuth:rules:%s:%s", roomId, userId)
	}

	UserCompereRules = func(userId, roomId string) string {
		return fmt.Sprintf("userAuth:rules:compere:%s:%s", roomId, userId)
	}

	//用户角色缓存
	UserRoles = func(userId, roomId string) string {
		return fmt.Sprintf("userAuth:roles:%s:%s", roomId, userId)
	}

	//静音请求频次
	MuteLocalSeatReqRate = func(userId, roomId string) string {
		return fmt.Sprintf("rate:muteLocal:%s:%s", roomId, userId)
	}

	// UserFollowJoinRoom 用户跟随进房缓存
	UserFollowJoinRoom = func(userId string) string {
		return fmt.Sprintf("UserFollowJoinRoom:%v:%v", userId, time.Now().Format("20060102"))
	}
	//装扮中心商品信息缓存
	GoodsAllCacheKey = func() string {
		return "goodsAll"
	}

	//更换手机号前置验证
	ChangeUserMobileBeforeVerify = func(userId string) string {
		return fmt.Sprintf("ChangeUserMobileBeforeVerify:%s", userId)
	}

	//用户充值支付锁
	UserPayLock = func(userId string) string {
		return fmt.Sprintf("UserPayLock:%s", userId)
	}

	//im服务器错误次数
	ImServerConnectFailCount = func(imServer string) string {
		return fmt.Sprintf("imServerConnectFail:%s", imServer)
	}
)
