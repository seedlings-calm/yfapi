package logic

import (
	"context"
	"math"
	"strings"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	"yfapi/internal/service/acl"
	"yfapi/internal/service/auth"
	service_goods "yfapi/internal/service/goods"
	service_im "yfapi/internal/service/im"
	service_level "yfapi/internal/service/level"
	"yfapi/internal/service/rankList"
	service_room "yfapi/internal/service/room"
	"yfapi/internal/service/user"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/message"
	"yfapi/typedef/redisKey"
	request_gift "yfapi/typedef/request/gift"
	response_gift "yfapi/typedef/response/gift"
	"yfapi/util/easy"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

type Gift struct {
}

func (g *Gift) GetGiftSourceList(c *gin.Context) (res *response_gift.GiftSourceListRes) {
	res = new(response_gift.GiftSourceListRes)
	sourceList, _ := new(dao.GiftDao).GetGiftSourceList()
	header := helper.GetHeaderData(c)
	isJson := header.Platform == typedef_enum.ClientTypePc
	for _, info := range sourceList {
		dst := response_gift.GiftSource{GiftCode: info.GiftCode}
		if isJson {
			dst.AnimationJsonUrl = helper.FormatImgUrl(info.AnimationJsonUrl)
		} else {
			dst.AnimationUrl = helper.FormatImgUrl(info.AnimationUrl)
		}
		res.List = append(res.List, dst)
	}
	return
}

func (g *Gift) GetRoomGiftList(c *gin.Context, req *request_gift.GiftListReq) (res *response_gift.GiftListRes) {
	res = new(response_gift.GiftListRes)
	res.CategoryType = req.CategoryType
	giftVersion := coreRedis.GetChatroomRedis().Get(c, redisKey.GetGiftVersionKey(req.CategoryType)).Val()
	if len(req.GiftVersion) > 0 && req.GiftVersion == giftVersion { // 版本未改变
		return
	} else {
		res.GiftVersion = giftVersion
	}

	nowTime := time.Now()
	giftList, err := new(dao.GiftDao).GetRoomGiftList(req.LiveType, req.RoomType, req.CategoryType, nowTime.Format(time.DateTime))
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	giftSendCountDao := new(dao.GiftSendCountDao)
	for _, info := range giftList {
		giftDTO := response_gift.GiftDTO{
			GiftId:           info.ID,
			GiftCode:         info.GiftCode,
			GiftName:         info.GiftName,
			GiftImage:        helper.FormatImgUrl(info.GiftImage),
			GiftAmountType:   info.GiftAmountType,
			GiftDiamond:      info.GiftDiamond,
			GiftGrade:        info.GiftGrade,
			CategoryType:     info.CategoryType,
			AnimationUrl:     helper.FormatImgUrl(info.AnimationUrl),
			AnimationJsonUrl: helper.FormatImgUrl(info.AnimationJsonUrl),
			SendCountList:    []response_gift.GiftSendCount{},
		}
		if !info.SubscriptStartTime.IsZero() {
			if info.SubscriptStartTime.Before(nowTime) || info.SubscriptEndTime.After(nowTime) {
				giftDTO.SubscriptContent = info.SubscriptContent
				giftDTO.SubscriptIcon = helper.FormatImgUrl(info.SubscriptIcon)
			}
		}
		// 处理赠送礼物数量列表
		if len(info.SendCountList) > 0 {
			giftDTO.SendCountList, _ = giftSendCountDao.GetListByIdList(strings.Split(info.SendCountList, ","))
		}
		res.List = append(res.List, giftDTO)
	}

	return
}

func (g *Gift) SendGift(c *gin.Context, req *request_gift.SendGiftReq) (res *response_gift.SendGiftRes) {
	res = new(response_gift.SendGiftRes)
	_ = helper.GetHeaderData(c)
	toUserCount := len(req.ToUserIdList)
	if toUserCount < 1 || toUserCount > 10 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 打赏人信息
	userId := helper.GetUserId(c)
	if len(req.ToUserIdList) == 1 {
		if !new(acl.RoomAcl).IsInRoom(userId, req.RoomId, helper.GetClientType(c)) {
			panic(error2.I18nError{
				Code: error2.ErrCodeUserNotInRoom,
				Msg:  nil,
			})
		}
	}
	userInfo, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if errCode := userInfo.CheckUserStatus(); errCode != error2.SuccessCode {
		panic(errCode)
	}
	// 超管巡查检查 无法打赏
	authService := new(auth.Auth)
	if authService.IsSuperAdminRole(req.RoomId, userId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSendGiftFromSuperAdmin,
			Msg:  nil,
		})
	}
	for _, toUserId := range req.ToUserIdList {
		if userId == toUserId {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		// 超管巡查检查 无法被打赏
		if authService.IsSuperAdminRole(req.RoomId, toUserId) {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSendGiftSuperAdmin,
				Msg:  nil,
			})
		}
	}
	// 房间信息
	roomInfo, err := new(dao.RoomDao).GetRoomById(req.RoomId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if roomInfo.Status != typedef_enum.RoomStatusNormal {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	// 礼物信息
	giftInfo, err := new(dao.GiftDao).GetGiftByCode(req.GiftCode)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if req.GiftDiamond != giftInfo.GiftDiamond {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}

	// 赠送的礼物价格
	sendDiamond := req.GiftCount * giftInfo.GiftDiamond
	if sendDiamond == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeRequestAbnormal,
			Msg:  nil,
		})
	}

	// 分布式锁
	success, unlock, err := coreRedis.UserLock(context.Background(), redisKey.SendGiftLockKey(userId), time.Second*5)
	if err != nil || !success {
		coreLog.Error("SendGift Lock err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFrequent,
			Msg:  nil,
		})
	}
	defer unlock()

	// 总赠送的礼物价格
	totalSendDiamond := sendDiamond * toUserCount
	// 处理打赏人和被打赏人的账户
	res.DiamondAmount = sendGiftAccountProcess(userId, giftInfo.GiftName, req, roomInfo)
	coreLog.Info("总打赏钻石：【%v】", totalSendDiamond)
	// 当前礼物的经验值
	currExp := cast.ToFloat64(giftInfo.GiftDiamond*req.GiftCount) * giftInfo.ExpTimes
	currStarExp := cast.ToInt(float64(sendDiamond*req.GiftCount) * typedef_enum.StarlightRate)
	// 处理玩家经验流程
	addLevelExpProcess(c, userId, userInfo.Nickname, req.RoomId, helper.GetClientType(c), req.ToUserIdList, cast.ToInt(currExp), currStarExp, res)
	// 增加房间热度
	go service_room.UpdateRoomHotBySendGift(req.RoomId, currExp, roomInfo.LiveType)
	//res.LevelConfig = 1
	//res.LvExp = 1
	// 连击
	// comboHitMap := make(map[string]int)
	// for _, toUserId := range req.ToUserIdList {
	comboCount := addGiftSendComboHitCount(req.RoomId, userId, req.ToUserIdList, req.GiftCount, userInfo, giftInfo, totalSendDiamond, req.IsBatch)
	// }
	// 增加贡献值
	go func() {
		pipe := coreRedis.GetChatroomRedis().Pipeline()
		onlineKey := redisKey.RoomUsersOnlineList(req.RoomId)
		consumeKey := redisKey.RoomUsersDayList(req.RoomId)
		seatKey := redisKey.RoomWheatPosition(req.RoomId)
		// 在线列表贡献值
		pipe.ZIncrBy(c, onlineKey, float64(totalSendDiamond), userId)
		// 贡献值榜单
		pipe.ZIncrBy(c, consumeKey, float64(totalSendDiamond), userId)
		charmMap := make(map[int]int)
		// 打赏礼物飞屏动效通知
		noticeData := message.MsgSendGiftSeat{GiftImage: helper.FormatImgUrl(giftInfo.GiftImage), FromSeatId: -1}
		for _, toUserId := range req.ToUserIdList {
			// 玩家是否在麦
			seatInfo := getUserInSeat(req.RoomId, toUserId)
			if seatInfo != nil {
				// 魅力值
				if roomInfo.LiveType == typedef_enum.LiveTypeChatroom {
					pipe.HIncrBy(c, redisKey.ChatroomUserCharmKey(req.RoomId), toUserId, int64(sendDiamond))
				} else {
					pipe.HIncrBy(c, redisKey.AnchorRoomUserCharmKey(req.RoomId), toUserId, int64(sendDiamond))
				}
				seatInfo.UserInfo.CharmCount += sendDiamond
				pipe.HSet(c, seatKey, seatInfo.Id, easy.JSONStringFormObject(seatInfo))
				charmMap[seatInfo.Id] = seatInfo.UserInfo.CharmCount
				// 在麦打赏动效
				noticeData.ToSeatIdList = append(noticeData.ToSeatIdList, seatInfo.Id)
			}
		}
		pipe.Expire(c, consumeKey, 24*time.Hour)
		_, _ = pipe.Exec(c)
		//加入，退出，变成贡献值，刷新榜单前三用户信息
		serviceRoom := &service_room.RoomUsersOnlie{}
		serviceRoom.OnlineChangeToThree(roomInfo.Id)
		// 推送魅力值变化
		if len(charmMap) > 0 {
			new(service_im.ImPublicService).SendCustomMsg(roomInfo.Id, charmMap, typedef_enum.SEAT_CHARM_MSG)
		}
		// 打赏礼物飞屏动效通知
		if len(noticeData.ToSeatIdList) > 0 {
			seatInfo := getUserInSeat(req.RoomId, userId)
			if seatInfo != nil {
				noticeData.FromSeatId = seatInfo.Id
			}
			new(service_im.ImPublicService).SendCustomMsg(roomInfo.Id, noticeData, typedef_enum.ROOM_SEND_GIFT_SEAT_MSG)
		}
	}()

	// 发布打赏消息
	if req.IsBatch && toUserCount > 1 { // 全麦打赏
		new(service_im.ImPublicService).SendGiftMsg(c, userId, "", req.RoomId, helper.GetHeaderData(c).Platform, "", req.GiftCount, comboCount, totalSendDiamond, req.IsBatch, &giftInfo)
	} else {
		req.IsBatch = false
		for _, toUserId := range req.ToUserIdList {
			md5Key := easy.Md5(req.RoomId+userId+req.GiftCode+toUserId, 0, false)
			new(service_im.ImPublicService).SendGiftMsg(c, userId, toUserId, req.RoomId, helper.GetHeaderData(c).Platform, md5Key, req.GiftCount, comboCount, sendDiamond, req.IsBatch, &giftInfo)
		}
	}
	//排行榜计算
	go func() {
		rankListService := rankList.Instance()
		if len(req.ToUserIdList) > 0 {
			if giftInfo.GiftAmountType == 2 {
				//免费礼物
				for _, v := range req.ToUserIdList {
					rankListService.Calculate(rankList.CalculateReq{
						FromUserId: userId,
						ToUserId:   v,
						Types:      "freeGift",
						RoomId:     req.RoomId,
					})
				}
			}
			if giftInfo.GiftAmountType == 1 {
				//收费礼物
				for _, v := range req.ToUserIdList {
					rankListService.Calculate(rankList.CalculateReq{
						FromUserId: userId,
						ToUserId:   v,
						Types:      "gift",
						Diamond:    req.GiftCount * giftInfo.GiftDiamond,
						RoomId:     req.RoomId,
					})
				}
			}
		}
	}()
	// //大于等于3000钻石，增加头条信息
	// if totalSendDiamond >= 3000 {
	// 	go service_gift.StoreTopMsg(userInfo, req.IsBatch, &giftInfo, toUserCount)
	// }
	//直播数据统计
	go service_room.DoRoomWheatTimeOperation(userId, req.RoomId, float64(totalSendDiamond), service_room.RewardCount)
	go service_room.DoRoomWheatTimeOperation(userId, req.RoomId, float64(toUserCount), service_room.RewardTimes)
	go service_room.DoRoomWheatTimeOperation(userId, req.RoomId, 1, service_room.RewardUserCount)
	return
}

// 增加经验
func addLevelExpProcess(c *gin.Context, userId, nickname, roomId, clientType string, toUserIdList []string, addExp, addStarExp int, res *response_gift.SendGiftRes) {
	// 增加打赏人的lv经验
	addLvExp(userId, nickname, roomId, addExp, res, c)
	// 增加打赏人的vip经验
	addVipExp(c, userId, nickname, roomId, addExp)
	// 增加被打赏人的星光经验
	for _, toUserId := range toUserIdList {
		// 是否为本房间的从业者 主持、音乐人、咨询师、主播
		checkRoleIdList := []int{typedef_enum.CompereRoleId, typedef_enum.MusicianRoleId, typedef_enum.CounselorRoleId, typedef_enum.AnchorRoleId}
		isHave := new(auth.Auth).IsHaveCurrRole(roomId, toUserId, checkRoleIdList)
		if isHave {
			addStarlightExp(c, toUserId, roomId, addStarExp)
		}
	}
}

// 增加lv经验
func addLvExp(userId, nickname, roomId string, addExp int, res *response_gift.SendGiftRes, c *gin.Context) {
	// 增加lv经验
	userLv, err := new(dao.UserLevelLvDao).GetUserLvLevel(userId)
	if err != nil {
		coreLog.Error("addLevelExpProcess.addLvExp.GetUserLvLevel userId:%v addExp:%v 增加lv经验失败：%v", userId, addExp, err)
		return
	}
	lvConfig, err := new(dao.LvConfigDao).GetAllLvConfigMap()
	if err != nil {
		coreLog.Error("addLevelExpProcess.addLvExp.GetAllLvConfigMap userId:%v addExp:%v 增加lv经验失败：%v", userId, addExp, err)
		return
	}
	userLv.CurrExp += addExp
	userLv.UpdateTime = time.Now()
	for userLv.CurrExp >= getMinExpByNextLvLevel(userLv.Level, lvConfig) {
		userLv.Level++
		// 获取特权物品配置
		privilegeData, err := service_level.GetLvConfigWithPrivilege(userLv.Level)
		if err != nil {
			coreLog.Error("addLevelExpProcess.addLvExp.GetLvConfigWithPrivilege userId:%v addExp:%v 增加lv经验失败：%v", userId, addExp, err)
			return
		}
		// 发放特权物品
		for _, item := range privilegeData.PrivilegeItemList {
			err = service_goods.UserGoods{}.SendGoodsToUser(c, userId, item.GoodsId, typedef_enum.GoodsGrantSourceLV, item.ExpirationDate, UpdateGoodsCallback())
			if err != nil {
				coreLog.Error("addLevelExpProcess.addLvExp.SendGoodsToUser userId:%v itemId:%v 发放等级特权物品失败：%v", userId, item.GoodsId, err)
			}
		}
		// 推送升级文案
		//if userLv.Level >= 50 && len(roomId) > 0 {
		if len(roomId) > 0 {
			msg := message.MsgLevelUp{
				UserId:   userId,
				Nickname: nickname,
				Icon:     privilegeData.Icon,
			}
			new(service_im.ImPublicService).SendCustomMsg(roomId, msg, typedef_enum.USER_LEVEL_UP_MSG)
		}
	}
	res.LvLevel = userLv.Level
	res.LvCurrExp = userLv.CurrExp
	res.LvMinExp = lvConfig[userLv.Level].MinExperience
	res.LvMaxExp = lvConfig[userLv.Level].MaxExperience
	res.LvIcon = helper.FormatImgUrl(lvConfig[userLv.Level].Icon)
	// 更新lv等级信息
	err = new(dao.UserLevelLvDao).Save(userLv)
	if err != nil {
		coreLog.Error("addLevelExpProcess.addLvExp.Save userId:%v addExp:%v 增加lv经验失败：%v", userId, addExp, err)
	}
}

// 获取下一级的lv经验
func getMinExpByNextLvLevel(currLevel int, configM map[int]*model.UserLvConfig) int {
	if _, isExist := configM[currLevel+1]; !isExist {
		return math.MaxInt
	}
	return configM[currLevel+1].MinExperience
}

// 增加vip经验
func addVipExp(c *gin.Context, userId, nickname, roomId string, addExp int) {
	// 增加vip经验
	userVip, err := new(dao.UserLevelVipDao).GetUserVipLevel(userId)
	if err != nil {
		coreLog.Error("addLevelExpProcess.addVipExp.GetUserLvLevel userId:%v addExp:%v 增加vip经验失败：%v", userId, addExp, err)
		return
	}
	vipConfig, err := new(dao.VipConfigDao).GetAllVipConfigMap()
	if err != nil {
		coreLog.Error("addLevelExpProcess.addVipExp.GetAllLvConfigMap userId:%v addExp:%v 增加vip经验失败：%v", userId, addExp, err)
		return
	}
	userVip.CurrExp += addExp
	userVip.KeepExp += addExp
	userVip.UpdateTime = time.Now()
	for userVip.CurrExp >= getMinExpByNextVipLevel(userVip.Level, vipConfig) {
		userVip.Level++
		userVip.KeepExp = 0
		// 更新过期时间
		userVip.ExpireTime = easy.GetCurrDayEndTime(time.Now()).AddDate(0, 0, 30)
		// 获取特权物品配置
		privilegeData, err := service_level.GetVipConfigWithPrivilege(userVip.Level)
		if err != nil {
			coreLog.Error("addLevelExpProcess.addVipExp.GetVipConfigWithPrivilege userId:%v addExp:%v 增加vip经验失败：%v", userId, addExp, err)
			return
		}
		// 发放特权物品
		for _, item := range privilegeData.PrivilegeItemList {
			err = service_goods.UserGoods{}.SendGoodsToUser(c, userId, item.GoodsId, typedef_enum.GoodsGrantSourceVIP, item.ExpirationDate, UpdateGoodsCallback())
			if err != nil {
				coreLog.Error("addLevelExpProcess.addVipExp.SendGoodsToUser userId:%v itemId:%v 发放等级特权物品失败：%v", userId, item.GoodsId, err)
			}
		}
		// 推送升级文案
		//if userVip.Level >= 6 && len(roomId) > 0 {
		if len(roomId) > 0 {
			msg := message.MsgLevelUp{
				UserId:   userId,
				Nickname: nickname,
				Icon:     privilegeData.Icon,
			}
			new(service_im.ImPublicService).SendCustomMsg(roomId, msg, typedef_enum.USER_LEVEL_UP_MSG)
		}
	}
	// 更新vip等级信息
	err = new(dao.UserLevelVipDao).Save(userVip)
	if err != nil {
		coreLog.Error("addLevelExpProcess.addVipExp.Save userId:%v addExp:%v 增加vip经验失败：%v", userId, addExp, err)
	}
}

// 获取下一级vip经验
func getMinExpByNextVipLevel(currLevel int, configM map[int]*model.UserVipConfig) int {
	if _, isExist := configM[currLevel+1]; !isExist {
		return math.MaxInt
	}
	return configM[currLevel+1].MinExperience
}

// 增加星光经验
func addStarlightExp(c *gin.Context, userId, roomId string, addExp int) {
	// 增加星光经验
	userStar, err := new(dao.UserLevelStarDao).GetUserStarLevel(userId)
	if err != nil {
		coreLog.Error("addLevelExpProcess.addStarlightExp.GetUserLvLevel userId:%v addExp:%v 增加星光经验失败：%v", userId, addExp, err)
		return
	}
	vipConfig, err := new(dao.StarConfigDao).GetAllStarConfigMap()
	if err != nil {
		coreLog.Error("addLevelExpProcess.addStarlightExp.GetAllLvConfigMap userId:%v addExp:%v 增加星光经验失败：%v", userId, addExp, err)
		return
	}
	userStar.CurrExp += addExp
	userStar.KeepExp += addExp
	userStar.UpdateTime = time.Now()
	for userStar.CurrExp >= getMinExpByNextStarLevel(userStar.Level, vipConfig) {
		userStar.Level++
		userStar.KeepExp = 0
		// 更新过期时间
		userStar.ExpireTime = easy.GetCurrDayEndTime(time.Now()).AddDate(0, 0, 30)
		// 获取特权物品配置
		privilegeData, err := service_level.GetStarConfigWithPrivilege(userStar.Level)
		if err != nil {
			coreLog.Error("addLevelExpProcess.addStarlightExp.GetStarConfigWithPrivilege userId:%v addExp:%v 增加星光经验失败：%v", userId, addExp, err)
			return
		}
		// 发放特权物品
		for _, item := range privilegeData.PrivilegeItemList {
			err = service_goods.UserGoods{}.SendGoodsToUser(c, userId, item.GoodsId, typedef_enum.GoodsGrantSourceStar, item.ExpirationDate, UpdateGoodsCallback())
			if err != nil {
				coreLog.Error("addLevelExpProcess.addStarlightExp.SendGoodsToUser userId:%v itemId:%v 发放等级特权物品失败：%v", userId, item.GoodsId, err)
			}
		}
		// 推送升级文案
		//if userStar.Level >= 6 && len(roomId) > 0 {
		if len(roomId) > 0 {
			userInfo := user.GetUserBaseInfo(userId)
			msg := message.MsgLevelUp{
				UserId:   userId,
				Nickname: userInfo.Nickname,
				Icon:     privilegeData.Icon,
			}
			new(service_im.ImPublicService).SendCustomMsg(roomId, msg, typedef_enum.USER_LEVEL_UP_MSG)
		}
	}
	// 更新星光等级信息
	err = new(dao.UserLevelStarDao).Save(userStar)
	if err != nil {
		coreLog.Error("addLevelExpProcess.addStarlightExp.Save userId:%v addExp:%v 增加星光经验失败：%v", userId, addExp, err)
	}
}

// 获取下一级星光经验
func getMinExpByNextStarLevel(currLevel int, configM map[int]*model.UserStarConfig) int {
	if _, isExist := configM[currLevel+1]; !isExist {
		return math.MaxInt
	}
	return configM[currLevel+1].MinExperience
}

// 处理打赏人和被打赏人的账户
func sendGiftAccountProcess(userId, giftName string, req *request_gift.SendGiftReq, roomInfo model.Room) string {
	toUserCount := len(req.ToUserIdList)
	// 赠送的礼物价格
	sendDiamond := req.GiftCount * req.GiftDiamond
	// 总赠送的礼物价格
	totalSendDiamond := int64(sendDiamond * toUserCount)
	// 检查打赏人账户余额是否充足
	accountDao := new(dao.UserAccountDao)
	fromUserAccount, err := accountDao.GetUserAccountByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if easy.StringToDecimal(fromUserAccount.DiamondAmount).LessThan(decimal.NewFromInt(totalSendDiamond)) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDiamondNotEnough,
			Msg:  nil,
		})
	}
	// 生成打赏订单
	tx := coreDb.GetMasterDb().Begin()
	// 生成订单
	orderInfo := &model.Order{
		OrderId:         new(accountBook.Order).OrderNum(accountBook.ORDER_DS),
		UserId:          userId,
		ToUserIdList:    strings.Join(req.ToUserIdList, ","),
		RoomId:          roomInfo.Id,
		GuildId:         roomInfo.GuildId,
		Gid:             req.GiftCode,
		TotalAmount:     cast.ToString(totalSendDiamond),
		PayAmount:       cast.ToString(totalSendDiamond),
		DiscountsAmount: "0",
		Num:             req.GiftCount * toUserCount,
		Currency:        accountBook.CURRENCY_DIAMOND,
		OrderType:       accountBook.ChangeDiamondRewardGift,
		OrderStatus:     1,
		PayStatus:       1,
		Note:            giftName,
		StatDate:        time.Now().Format(time.DateOnly),
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	err = tx.Create(orderInfo).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 处理打赏人和被打赏人的账户
	// 扣除打赏人钻石
	user.UpdateUserAccount(&user.UpdateAccountParam{
		Tx:           tx,
		UserId:       userId,
		ToUserIdList: strings.Join(req.ToUserIdList, ","),
		Gid:          req.GiftCode,
		Num:          req.GiftCount * toUserCount,
		Currency:     accountBook.CURRENCY_DIAMOND,
		FundFlow:     accountBook.FUND_OUTFLOW,
		Amount:       cast.ToString(totalSendDiamond),
		OrderId:      orderInfo.OrderId,
		OrderType:    accountBook.ChangeDiamondRewardGift,
		RoomId:       roomInfo.Id,
		GuildId:      roomInfo.GuildId,
		Note:         giftName,
	})
	// 被打赏人增加星光收益
	for _, toUserId := range req.ToUserIdList {
		currency := accountBook.CURRENCY_STARLIGHT_UNWITHDRAW
		// 是否为本房间从业者
		if user.IsRoomPractitioner(toUserId, roomInfo.Id) {
			currency = accountBook.CURRENCY_STARLIGHT_WITHDRAW
		}
		user.UpdateUserAccount(&user.UpdateAccountParam{
			Tx:         tx,
			UserId:     toUserId,
			FromUserId: userId,
			Gid:        req.GiftCode,
			Num:        req.GiftCount,
			Diamond:    sendDiamond,
			Currency:   currency,
			FundFlow:   accountBook.FUND_INFLOW,
			Amount:     cast.ToString(float64(sendDiamond) * typedef_enum.StarlightRate),
			OrderId:    orderInfo.OrderId,
			OrderType:  accountBook.ChangeStarlightRewardIncome,
			RoomId:     roomInfo.Id,
			GuildId:    roomInfo.GuildId,
			Note:       giftName,
		})
	}
	tx.Commit()
	return easy.StringFixed(easy.StringToDecimal(fromUserAccount.DiamondAmount).Sub(decimal.NewFromInt(totalSendDiamond)))
}

func addGiftSendComboHitCount(roomId, userId string, toUserId []string, giftCount int, userInfo *model.User, giftInfo model.GiftDTO, totalSendDiamond int, batch bool) int {
	toUserIdS := strings.Join(toUserId, ",")
	key := redisKey.GiftComboHitCountKey(roomId, userId, giftInfo.GiftCode, toUserIdS, giftCount)
	count, err := coreRedis.GetChatroomRedis().Incr(context.Background(), key).Result()
	if err != nil {
		coreLog.Error("打赏礼物增加连击次数失败：%v", err)
	}
	_ = coreRedis.GetChatroomRedis().Expire(context.Background(), key, typedef_enum.ComboHitGapTime)
	// 使用 HASH 存储详细数据
	taskID := redisKey.TopMsgCallBack(key)
	if batch && len(toUserId) > 1 {
		toUserIdS = "全麦"
	}
	_, err = coreRedis.GetChatroomRedis().HSet(context.Background(), taskID, map[string]interface{}{
		"user_id":       userInfo.Id,
		"nickname":      userInfo.Nickname,
		"avatar":        userInfo.Avatar,
		"types":         "全服礼物",
		"operate":       "打赏",
		"to_user":       toUserIdS,
		"total_diamond": totalSendDiamond,   //连击的单次总价
		"gift_count":    giftCount,          //购买数量
		"gift_img":      giftInfo.GiftImage, //礼物logo
		"gift_name":     giftInfo.GiftName,  //礼物名称
		"combo_count":   count,              // 连击次数
	}).Result()
	if err != nil {
		coreLog.Error("存储打赏信息失败：%v", err)
		return int(count)
	}

	// 将任务添加到延迟队列，延迟 4 秒
	zsetKey := redisKey.TopMsgEqueueKey()
	timestamp := float64(time.Now().Add(typedef_enum.ComboHitGapTime).Unix())
	coreRedis.GetChatroomRedis().ZAdd(context.Background(), zsetKey, redis.Z{
		Score:  timestamp,
		Member: key,
	})
	return int(count)
}
