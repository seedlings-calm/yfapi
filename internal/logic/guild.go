package logic

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	dao2 "yfapi/internal/dao/guild"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	model2 "yfapi/internal/model/guild"
	"yfapi/internal/service/accountBook"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_h5 "yfapi/typedef/request/h5"
	response_h5 "yfapi/typedef/response/h5"
)

type Guild struct {
}

func (g *Guild) JoinGuild(req *request_h5.JoinGuildReq, context *gin.Context) (res response_h5.JoinGuildRes) {
	userId := handle.GetUserId(context)
	//查询工会是否存在
	GuildDao := &dao.GuildDao{}
	data, err := GuildDao.FindById(&model.Guild{
		ID: req.GuildId})
	if err != nil || data == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGuildNotExist,
			Msg:  nil,
		})
	}

	// 检测用户是否已经是工会的成员
	isMember, err := GuildDao.IsGuildMember(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if isMember.ID > 0 || isMember.GuildID == req.GuildId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGuildAlreadyJoin,
			Msg:  nil,
		})
	}
	//查询用户是不是在申请工会
	checkUserGuild := &model.GuildMemberApply{
		GuildID: req.GuildId,
		UserID:  userId,
	}
	datas, err := GuildDao.GetCheckUserApplication(checkUserGuild)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if datas.Status == 1 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGuildAlreadyApplication,
			Msg:  nil,
		})
	}
	GuildMemberDao := &model.GuildMemberApply{
		GuildID:    req.GuildId,
		UserID:     userId,
		ApplyType:  1,
		Force:      0,
		Status:     1,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err = GuildDao.Create(GuildMemberDao)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	res.Status = GuildMemberDao.Status
	return
}

// GetGuildInfo
//
//	@Description: 获取工会详情
//	@receiver g
//	@param req
//	@param context
//	@return res
func (g *Guild) GetGuildInfo(req *request_h5.GuildInfoReq, context *gin.Context) (res response_h5.GuildInfoRes) {
	userId := handle.GetUserId(context)
	GuildDao := &dao.GuildDao{}
	now := time.Now()
	//如果传了工会id，就直接查询工会详情
	if len(req.GuildId) > 0 {
		//查询用户有什么申请加入工会或者已经在工会里
		//查询用户有没有工会
		isMembers, err := GuildDao.IsGuildMember(userId)
		fmt.Println(userId)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		//查询用户有没有在申请
		checkUserGuild := &model.GuildMemberApply{
			UserID: userId,
		}
		CheckUserGuild, err := GuildDao.GetCheckUserApplication(checkUserGuild)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		isOk := CheckUserGuild.CreateTime.Add(time.Hour * 24).Before(now)
		var applyStatus int
		if isMembers.GuildID == req.GuildId {
			applyStatus = 2
		} else if CheckUserGuild.GuildID == req.GuildId && isOk == false {
			applyStatus = 1
		}
		//如果有就返回工会信息,聊天室列表,直播列表信息
		list, err := GuildDao.GetRoomList(req.GuildId)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		//查询工会信息
		data, err := GuildDao.FindById(&model.Guild{
			ID: req.GuildId})
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		userInfoDao := &dao.UserDao{}
		//查询工会会长昵称和头像
		userInfo, err := userInfoDao.FindOne(&model.User{
			Id: data.UserID})
		if err != nil || list == nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserNotFound,
				Msg:  nil,
			})
		}
		guildMember, err := GuildDao.GerGuildMember(req.GuildId)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		var num1, num2 int
		for _, v := range list {
			if v.LiveType == 1 {
				//返回聊天列表
				res.ChatRoomList = append(res.ChatRoomList, &response_h5.RoomInfo{
					Id:           v.Id,
					Name:         v.Name,
					RoomNo:       v.RoomNo,
					CoverImg:     helper.FormatImgUrl(v.CoverImg),
					Status:       v.Status,
					RoomType:     v.RoomType,
					RoomTypeDesc: enum.RoomType(v.RoomType).String(),
				})
				num1++
			} else if v.LiveType == 2 {
				//返回直播列表
				res.LiveRoomList = append(res.LiveRoomList, &response_h5.RoomInfo{
					Id:           v.Id,
					Name:         v.Name,
					RoomNo:       v.RoomNo,
					CoverImg:     helper.FormatImgUrl(v.CoverImg),
					Status:       v.Status,
					RoomType:     v.RoomType,
					RoomTypeDesc: enum.RoomType(v.RoomType).String(),
				})
				num2++
			}
		}
		res.Guild = &response_h5.GuildMsg{
			GuildId:        data.ID,
			GuildNo:        data.GuildNo,
			GuildName:      data.Name,
			GuildLogo:      helper.FormatImgUrl(data.LogoImg),
			GuildBriefDesc: data.BriefDesc,
			NickName:       userInfo.Nickname,
			Avatar:         helper.FormatImgUrl(userInfo.Avatar),
			Number:         int(guildMember),
			ApplyStatus:    applyStatus,
		}
		res.ChatRoomCount = int64(num1)
		res.LiveRoomCount = int64(num2)
	} else {
		//如果没有传递工会id，执行以下逻辑
		//查询用户有没有工会
		isMember, err := GuildDao.IsGuildMember(userId)
		fmt.Println(userId)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		//查询用户有没有在申请
		checkUserGuild := &model.GuildMemberApply{
			UserID: userId,
		}
		datas, err := GuildDao.GetCheckUserApplication(checkUserGuild)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		var guildID string
		var applyStatus int
		if isMember.ID > 0 {
			// 用户已经是工会成员
			guildID = isMember.GuildID
			applyStatus = 2
		} else if datas != nil && datas.Status == 1 {
			// 用户在申请工会
			applyStatus = 1
			guildID = datas.GuildID
		}
		isOk := datas.CreateTime.Add(time.Hour * 24).Before(now)
		if isMember.ID > 0 || datas.Status == 1 && datas.ID > 0 && isOk == false {
			//如果有就返回工会信息,聊天室列表,直播列表信息
			list, err := GuildDao.GetRoomList(guildID)
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeReadDB,
					Msg:  nil,
				})
			}
			//查询工会信息
			data, err := GuildDao.FindById(&model.Guild{
				ID: guildID})
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeReadDB,
					Msg:  nil,
				})
			}
			userInfoDao := &dao.UserDao{}
			//查询工会会长昵称和头像
			userInfo, err := userInfoDao.FindOne(&model.User{
				Id: data.UserID})
			if err != nil || list == nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeUserNotFound,
					Msg:  nil,
				})
			}
			guildMember, err := GuildDao.GerGuildMember(guildID)
			var num1, num2 int
			for _, v := range list {
				if v.LiveType == 1 {
					//返回聊天列表
					res.ChatRoomList = append(res.ChatRoomList, &response_h5.RoomInfo{
						Id:           v.Id,
						Name:         v.Name,
						RoomNo:       v.RoomNo,
						CoverImg:     helper.FormatImgUrl(v.CoverImg),
						Status:       v.Status,
						RoomType:     v.RoomType,
						RoomTypeDesc: enum.RoomType(v.RoomType).String(),
					})
					num1++
				} else if v.LiveType == 2 {
					//返回直播列表
					res.LiveRoomList = append(res.LiveRoomList, &response_h5.RoomInfo{
						Id:           v.Id,
						Name:         v.Name,
						RoomNo:       v.RoomNo,
						CoverImg:     helper.FormatImgUrl(v.CoverImg),
						Status:       v.Status,
						RoomType:     v.RoomType,
						RoomTypeDesc: enum.RoomType(v.RoomType).String(),
					})
					num2++
				}
			}
			res.Guild = &response_h5.GuildMsg{
				GuildId:        data.ID,
				GuildNo:        data.GuildNo,
				GuildName:      data.Name,
				GuildLogo:      helper.FormatImgUrl(data.LogoImg),
				GuildBriefDesc: data.BriefDesc,
				NickName:       userInfo.Nickname,
				Avatar:         helper.FormatImgUrl(userInfo.Avatar),
				Number:         int(guildMember),
				ApplyStatus:    applyStatus,
			}
			res.ChatRoomCount = int64(num1)
			res.LiveRoomCount = int64(num2)
			// 查询当前用户是否有从业者记录
			record, err := new(dao.DaoUserPractitioner).GetGuildPractitionerRecord(userId, data.ID)
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeReadDB,
					Msg:  nil,
				})
			}
			// 用户是否为从业者
			if len(record) > 0 {
				res.IsPractitioner = true
			}
			// 申请退会记录
			err = coreDb.GetSlaveDb().Model(model.GuildMemberApply{}).Where("user_id=? and guild_id=? and apply_type=2 and status<5 and create_time>?", userId, data.ID, isMember.CreateTime).
				Order("create_time desc").Limit(1).Scan(&res.QuitGuildApply).Error
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeReadDB,
					Msg:  nil,
				})
			}
		} else {
			//如果没有则返回工会列表
			//查询工会列表
			guildList, err := GuildDao.GetGuildList()
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeReadDB,
					Msg:  nil,
				})
			}
			for _, guild := range guildList {
				member, err := GuildDao.GerGuildMember(guild.ID)
				if err != nil {
					panic(error2.I18nError{
						Code: error2.ErrorCodeReadDB,
						Msg:  nil,
					})
				}
				res.GuildList = append(res.GuildList, &response_h5.GuildInfo{
					GuildId:        guild.ID,
					GuildNo:        guild.GuildNo,
					GuildName:      guild.Name,
					GuildLogo:      helper.FormatImgUrl(guild.LogoImg),
					GuildBriefDesc: guild.BriefDesc,
					Number:         int(member),
				})
			}
		}
	}
	return
}

// GetGuildMemberList
//
//	@Description:	获取工会成员列表
//	@receiver g
//	@param req
//	@param context
//	@return res
func (g *Guild) GetGuildMemberList(req *request_h5.GuildMemberListReq, context *gin.Context) (res response_h5.GuildMemberListRes) {
	userId := handle.GetUserId(context)
	GuildDao := &dao.GuildDao{}
	//查询用户是不是会长，如果不是不能查看工会成员
	data, err := GuildDao.FindById(&model.Guild{
		ID: req.GuildId})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if data == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGuildNotExist,
			Msg:  nil,
		})
	}
	if data.UserID != userId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGuildNotLeader,
			Msg:  nil,
		})
	}
	list, err := GuildDao.GetGuildMemberList(req.GuildId)
	for _, v := range list {
		var role int8
		if v.UserID == userId {
			role = 1 //会长
		} else {
			role = 2 //成员
		}
		userInfoDao := &dao.UserDao{}
		userInfo, err := userInfoDao.FindOne(&model.User{
			Id: v.UserID})
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		userPlaque := service_user.GetUserLevelPlaque(userId, helper.GetClientType(context))
		res.GuildMemberList = append(res.GuildMemberList, &response_h5.GuildMemberInfo{
			UserId:     v.UserID,
			UserNo:     userInfo.UserNo,
			Uid32:      cast.ToInt32(userInfo.OriUserNo),
			Sex:        userInfo.Sex,
			NickName:   userInfo.Nickname,
			Avatar:     helper.FormatImgUrl(userInfo.Avatar),
			Introduce:  userInfo.Introduce,
			Role:       role,
			UserPlaque: userPlaque,
		})
	}
	return
}

// QuitGuildApply
//
//	@Description: 退出公会申请
//	@receiver g
//	@param c *gin.Context -
//	@param req *request_h5.QuitGuildApplyReq -
func (g *Guild) QuitGuildApply(c *gin.Context, req *request_h5.QuitGuildApplyReq) {
	userId := handle.GetUserId(c)
	guildInfo, err := new(dao.GuildDao).FindById(&model.Guild{ID: req.GuildId})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGuildNotExist,
			Msg:  nil,
		})
	} else if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 拦截会长的请求
	if guildInfo.UserID == userId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 是否为该公会成员
	if !new(dao.GuildDao).GetCheckUserInGuild(req.GuildId, userId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 用户是否有待审核的退会申请
	memberApplyDao := new(dao2.GuildMemberApplyDao)
	applyInfo, err := memberApplyDao.GetGuildMemberApply(userId, req.GuildId, 2, 1)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if applyInfo.ID > 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeQuitGuildApplyExamine,
			Msg:  nil,
		})
	}
	// 强制申请退出
	if req.IsForced {
		// 查询当前用户是否有从业者记录
		record, err := new(dao.DaoUserPractitioner).GetGuildPractitionerRecord(userId, req.GuildId)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		// 从业者无法直接强制退出，需先缴纳违约金
		if len(record) > 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeForcedQuitGuildFee,
				Msg:  nil,
			})
		}
		// 查询用户最新的被拒绝的退会记录
		applyInfo, err = memberApplyDao.GetGuildMemberApply(userId, req.GuildId, 2, 3)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		// 没有被拒绝的退会记录
		if applyInfo.ID == 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		// 修改记录为强制退会
		err = coreDb.GetMasterDb().Model(model.GuildMemberApply{}).Where("id", applyInfo.ID).Updates(map[string]interface{}{
			"force":       1,
			"status":      1,
			"update_time": time.Now(),
		}).Error
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		return
	}

	// 添加申请退会记录
	err = memberApplyDao.Create(&model.GuildMemberApply{
		ID:         0,
		GuildID:    req.GuildId,
		UserID:     userId,
		ApplyType:  2,
		Force:      0,
		Status:     1,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
}

// QuitGuildApplyCancel
//
//	@Description: 取消退出公会申请
//	@param c *gin.Context -
//	@param req *response_h5.GuildInfo -
func (g *Guild) QuitGuildApplyCancel(c *gin.Context, req *request_h5.GuildInfoReq) {
	userId := handle.GetUserId(c)
	// 是否为该公会成员
	if !new(dao.GuildDao).GetCheckUserInGuild(req.GuildId, userId) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 用户是否有待审核的退会申请
	memberApplyDao := new(dao2.GuildMemberApplyDao)
	applyInfo, err := memberApplyDao.GetGuildMemberApply(userId, req.GuildId, 2, 1)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if applyInfo.ID > 0 {
		err = memberApplyDao.MemberApplyUpdate(model.GuildMemberApply{
			ID:         applyInfo.ID,
			Status:     6,
			UpdateTime: time.Now(),
		})
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
	}
}

// GetGuildPenaltyDetail
//
//	@Description:违约金详情
//	@param c *gin.Context -
//	@return res -
func (g *Guild) GetGuildPenaltyDetail(c *gin.Context) (res response_h5.GuildPenaltyDetailRes) {
	userId := helper.GetUserId(c)
	// 查询用户星光等级
	starInfo, _ := new(dao.UserLevelStarDao).FirstById(userId)
	// 查询星光等级配置
	configList, _ := new(dao.StarConfigDao).GetAllStarConfigList()
	for _, info := range configList {
		penaltyDiamond := cast.ToString(info.MinExperience * info.PenaltyRate / 100)
		if len(penaltyDiamond) == 0 {
			penaltyDiamond = "0"
		}
		res.PenaltyList = append(res.PenaltyList, response_h5.PenaltyDetail{
			LevelName:      info.LevelName + "级",
			MinExp:         info.MinExperience,
			PenaltyRate:    cast.ToString(info.PenaltyRate),
			PenaltyDiamond: cast.ToString(info.MinExperience * info.PenaltyRate / 100),
		})
		if starInfo.Level > 0 && starInfo.Level == info.Level {
			res.StarLevel = starInfo.Level
			res.CurrExp = starInfo.CurrExp
			res.PenaltyDiamond = cast.ToString(info.MinExperience * info.PenaltyRate / 100)
			if len(res.PenaltyDiamond) == 0 {
				res.PenaltyDiamond = "0"
			}
			// 计算扣除的星光经验
			// 公会成员信息
			memberInfo, err := new(dao.GuildDao).GetGuildMemberInfo(userId)
			if err != nil {
				panic(error2.I18nError{
					Code: error2.ErrorCodeReadDB,
					Msg:  nil,
				})
			}
			res.DeductExp = starInfo.CurrExp - memberInfo.StarExp
		}
	}
	return
}

// PayGuildPenalty
//
//	@Description: 缴纳违约金(从业者强制退会)
//	@param c *gin.Context -
//	@param req request_h5.GuildInfoReq -
func (g *Guild) PayGuildPenalty(c *gin.Context, req *request_h5.GuildInfoReq) {
	userId := handle.GetUserId(c)
	// 是否为该公会成员
	memberInfo, err := new(dao.GuildDao).GetGuildMemberInfo(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if memberInfo.GuildID != req.GuildId {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 查询当前用户是否有从业者记录
	record, err := new(dao.DaoUserPractitioner).GetGuildPractitionerRecord(userId, req.GuildId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 不是从业者，无需缴纳违约金
	if len(record) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 公会信息
	guildInfo, _ := new(dao.GuildDao).FindById(&model.Guild{ID: req.GuildId})
	if len(guildInfo.ID) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeGuildNotExist,
			Msg:  nil,
		})
	}
	// 查询违约金
	config := g.GetGuildPenaltyDetail(c)
	orderId := new(accountBook.Order).OrderNum(accountBook.ORDER_OR)
	tx := coreDb.GetMasterDb().Begin()
	if len(config.PenaltyDiamond) > 0 && config.PenaltyDiamond != "0" {
		// 扣除违约金
		service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
			Tx:           tx,
			UserId:       userId,
			ToUserIdList: guildInfo.UserID,
			Num:          1,
			Currency:     accountBook.CURRENCY_DIAMOND,
			FundFlow:     2,
			Amount:       config.PenaltyDiamond,
			OrderId:      orderId,
			OrderType:    accountBook.ChangeDiamondPenalty,
			GuildId:      req.GuildId,
			Note:         "强制退出公会违约金",
		})
		// 给会长增加钻石
		service_user.UpdateUserAccount(&service_user.UpdateAccountParam{
			Tx:         tx,
			UserId:     guildInfo.UserID,
			FromUserId: userId,
			Num:        1,
			Currency:   accountBook.CURRENCY_DIAMOND,
			FundFlow:   1,
			Amount:     config.PenaltyDiamond,
			OrderId:    orderId,
			OrderType:  accountBook.ChangeDiamondPenalty,
			GuildId:    req.GuildId,
			Note:       "强制退出公会违约金",
		})
		// 创建缴纳违约金订单
		err = tx.Create(&model.Order{
			OrderId:         orderId,
			UserId:          userId,
			ToUserIdList:    guildInfo.UserID,
			RoomId:          "0",
			GuildId:         req.GuildId,
			TotalAmount:     config.PenaltyDiamond,
			PayAmount:       config.PenaltyDiamond,
			DiscountsAmount: "0",
			Num:             1,
			Currency:        accountBook.CURRENCY_DIAMOND,
			OrderType:       accountBook.ChangeDiamondPenalty,
			OrderStatus:     1,
			PayType:         1,
			PayStatus:       1,
			WithdrawStatus:  1,
			Note:            "强制退出公会违约金",
			StatDate:        time.Now().Format(time.DateOnly),
			CreateTime:      time.Now(),
			UpdateTime:      time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	}
	// 如果是房主把房间退给会长
	err = tx.Model(model.Room{}).Where("user_id=? and guild_id=? and live_type=1", userId, req.GuildId).Updates(map[string]interface{}{
		"user_id":              guildInfo.UserID,
		"day_settle_user_id":   guildInfo.UserID,
		"month_settle_user_id": guildInfo.UserID,
	}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 查询公会所有房间列表
	roomList, err := new(dao.RoomDao).GetRoomsByGuildId(req.GuildId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
	} else if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	var roomIdList []string
	for _, info := range roomList {
		if info.LiveType == enum.LiveTypeAnchor && info.UserId == userId { // 是主播且已并入公会
			// 主播房间和公会解约
			err = tx.Model(model.Room{}).Where("user_id=? and room_id=?", userId, info.Id).Updates(map[string]interface{}{
				"guild_id":    "",
				"update_time": time.Now(),
			}).Error
			if err != nil {
				tx.Rollback()
				panic(error2.I18nError{
					Code: error2.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
			continue
		}
		roomIdList = append(roomIdList, info.Id)
	}
	// 删除玩家房间身份
	err = tx.Model(model.AuthRoleAccess{}).Where("user_id=? and room_id in ?", userId, roomIdList).Delete(&model.AuthRoleAccess{}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 删除玩家房间管理员
	err = tx.Model(model.RoomAdmin{}).Where("user_id=? and room_id in ?", userId, roomIdList).Delete(&model.RoomAdmin{}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 取消从业者身份
	err = tx.Model(model.UserPractitioner{}).Where("user_id=? and status=1 and room_id in ?", userId, roomIdList).Updates(map[string]interface{}{
		"status":         4,
		"abolish_reason": "强制退出公会",
		"update_time":    time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 申请中的从业者身份自动拒绝
	err = tx.Model(model.UserPractitioner{}).Where("user_id=? and status=2 and room_id in ?", userId, roomIdList).Updates(map[string]interface{}{
		"status":         3,
		"abolish_reason": "强制退出公会",
		"update_time":    time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 申请中的房间如果是房主自动拒绝
	err = tx.Model(model2.GuildRoomApply{}).Where("room_user_id=? and status=1 and guild_id=?", userId, req.GuildId).Updates(map[string]interface{}{
		"status":      3,
		"reason":      "房主强制退出公会",
		"update_time": time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 退出公会
	err = tx.Model(model.GuildMember{}).Where("user_id=? and status != 3 and guild_id=?", userId, req.GuildId).Updates(map[string]interface{}{
		"status":       3,
		"leave_reason": "强制退出公会",
		"update_time":  time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 修改强制申请退会记录
	// 查询用户最新的被拒绝的退会记录
	applyInfo, err := new(dao2.GuildMemberApplyDao).GetGuildMemberApply(userId, req.GuildId, 2, 3)
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 没有被拒绝的退会记录
	if applyInfo.ID == 0 {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 修改记录为强制退会
	err = coreDb.GetMasterDb().Model(model.GuildMemberApply{}).Where("id", applyInfo.ID).Updates(map[string]interface{}{
		"force":       1,
		"status":      5,
		"update_time": time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	// 增加用户的从业者行为记录
	err = tx.Create(&model.UserPractitionerAction{
		UserId:     userId,
		GuildId:    req.GuildId,
		Action:     "强制退会",
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	// 扣除星光等级
	// 当前星光等级低于入会等级，保留当前等级不做处理
	if config.StarLevel >= memberInfo.StarLevel && config.CurrExp > memberInfo.StarExp {
		// 降级重新计算有效期，并且保级经验清零
		err = tx.Model(model.UserStarLevel{}).Where("user_id=?", userId).Updates(map[string]interface{}{
			"level":       memberInfo.StarLevel,
			"curr_exp":    memberInfo.StarExp,
			"keep_exp":    "0",
			"expire_time": time.Now().AddDate(0, 0, 30),
			"update_time": time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			panic(error2.I18nError{
				Code: error2.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
	}
	tx.Commit()
	// 删除房间身份缓存
	var keys []string
	for _, roomId := range roomIdList {
		key1 := redisKey.UserRules(userId, roomId)
		key2 := redisKey.UserRoles(userId, roomId)
		key3 := redisKey.UserCompereRules(userId, roomId)
		keys = append(keys, key1, key2, key3)
	}
	coreRedis.GetUserRedis().Del(c, keys...)
}
