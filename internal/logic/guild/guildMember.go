package guild

import (
	"errors"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	dao2 "yfapi/internal/dao/guild"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	model2 "yfapi/internal/model/guild"
	"yfapi/internal/service/accountBook"
	"yfapi/typedef/enum"
	typedef_enum "yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_room "yfapi/typedef/request/guild"
	"yfapi/typedef/response"
	response_guild "yfapi/typedef/response/guild"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type GuildMember struct {
}

func (g *GuildMember) SaveMemberGroup(c *gin.Context, req *request_room.AddGuildGroupReq) (resp error) {
	GuildGroupDao := &model.GuildGroup{
		GuildID:    c.GetString("guildId"),
		GroupName:  req.GroupName,
		Desc:       req.Desc,
		Status:     1,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err := new(dao2.GuildGroupDao).Create(GuildGroupDao)
	if err != nil {
		panic(i18n_err.ErrorCodeUpdateDB)
	}

	return
}

func (g *GuildMember) MemberGroupUpdate(c *gin.Context, req *request_room.MemberGroupUpdateReq) (resp error) {
	db := new(dao2.GuildGroupDao)
	res, err := db.FindOne(&model.GuildGroup{ID: req.Id})
	if err != nil {
		panic(i18n_err.ErrorCodeDataNotFound)
	}
	res.GroupName = req.GroupName
	res.Desc = req.Desc
	res.UpdateTime = time.Now()

	err = db.Save(res)
	if err != nil {
		panic(i18n_err.ErrorCodeUpdateDB)
	}
	return
}

func (g *GuildMember) MemberGroupDelete(c *gin.Context, guildId string) (err error) {
	db := new(dao2.GuildGroupDao)
	res, err := db.FindOne(&model.GuildGroup{ID: cast.ToInt(guildId)})
	if err != nil {
		panic(i18n_err.ErrorCodeDataNotFound)
	}
	tx := coreDb.GetMasterDb().Begin()
	res.Status = 2
	err = tx.Model(res).Save(res).Error
	if err != nil {
		tx.Rollback()
		return
	}
	err = tx.Model(model.GuildMember{}).Where("group_id = ?", res.ID).Update("group_id", 0).Error
	if err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

// 分组列表增加成员数量
func (g *GuildMember) MemberGroup(c *gin.Context, req *request_room.GuildGroupListreq) (resp response.AdminPageRes) {
	list, count, err := new(dao2.GuildGroupDao).GetGuildGroupListPage(req, c)
	if err != nil {
		panic(i18n_err.ErrorCodeDataNotFound)
	}
	resp.Total = count
	resp.Data = list
	resp.CurrentPage = req.CurrentPage
	resp.Size = req.Size
	return
}

func (g *GuildMember) SetGroupByMembers(c *gin.Context, req *request_room.SetGroupByMembersReq) (err error) {
	if len(req.Ids) == 0 {
		return
	}
	gdao := dao2.GuildGroupDao{}
	info, err := gdao.FindOne(&model.GuildGroup{ID: req.GroupId})
	if err != nil || info.ID == 0 {
		return
	}
	err = coreDb.GetMasterDb().Model(model.GuildMember{}).Where("id in ?", req.Ids).UpdateColumn("group_id", req.GroupId).Error
	return
}

func (g *GuildMember) getMemberCount(req *request_room.GuildMemberListReq, guildId string) (total int64) {
	// 初始化查询
	tx := coreDb.GetSlaveDb().Table("t_guild_member gm").
		Joins("left join t_order_bill ob on gm.user_id=ob.user_id and ob.guild_id=? and ob.order_type=?", guildId, accountBook.ChangeStarlightRewardIncome).
		Joins("left join t_user u on u.id=gm.user_id").
		Joins("left join t_user_practitioner up on up.user_id=gm.user_id and up.status=1").
		Joins("left join t_guild_group gg on gg.id = gm.group_id").
		Where("gm.guild_id=? and gm.status != 3 ", guildId)
	// 根据用户关键词过滤
	if len(req.UserKeyword) > 0 {
		tx = tx.Where("u.nickname like ? OR u.user_no like ?", easy.GenLikeSql(req.UserKeyword), easy.GenLikeSql(req.UserKeyword))
	}

	// 根据身份证号过滤
	if len(req.IdCard) > 0 {
		tx = tx.Where("up.practitioner_type in ?", req.IdCard)
	}

	if req.GroupID > 0 {
		tx = tx.Where("gg.guild_id = ? and  gg.id = ?", guildId, req.GroupID)
	}
	// 统计总数
	err := tx.Distinct("gm.user_id").Count(&total).Error
	if err != nil {
		total = 0
	}
	return
}
func (g *GuildMember) MemberList(c *gin.Context, req *request_room.GuildMemberListReq) (res response.AdminPageRes) {
	guildId := helper.GetGuildId(c)
	res.Total = g.getMemberCount(req, guildId)
	// 初始化查询
	tx := coreDb.GetSlaveDb().Table("t_guild_member gm").
		Joins("left join t_order_bill ob on gm.user_id=ob.user_id and ob.guild_id=? and ob.order_type=?", guildId, accountBook.ChangeStarlightRewardIncome).
		Joins("left join t_user u on u.id=gm.user_id").
		Joins("left join t_user_practitioner up on up.user_id=gm.user_id and up.status=1").
		Joins("left join t_guild_group gg on gg.id = gm.group_id").
		Where("gm.guild_id=? and gm.status != 3 ", guildId)
	// 根据用户关键词过滤
	if len(req.UserKeyword) > 0 {
		tx = tx.Where("u.nickname like ? OR u.user_no like ?", easy.GenLikeSql(req.UserKeyword), easy.GenLikeSql(req.UserKeyword))
	}

	// 根据身份证号过滤
	if len(req.IdCard) > 0 {
		tx = tx.Where("up.practitioner_type in ?", req.IdCard)
	}

	if req.GroupID > 0 {
		tx = tx.Where("gg.guild_id = ? and  gg.id = ?", guildId, req.GroupID)
	}
	//获取公会的所有房间ID
	roomDao := dao.RoomDao{}
	roomIds, _ := roomDao.GetRoomIdsByGuildId(guildId)
	var userRules []model.AuthRoleAccess
	if len(roomIds) > 0 {
		coreDb.GetMasterDb().Table("t_auth_role_access ").Where("room_id in ? and role_id in ?", roomIds, []int{enum.PresidentRoleId, enum.HouseOwnerRoleId}).Find(&userRules)
	}
	// 公会成员
	var result []*response_guild.GuildMemberListRes
	// 查询具体数据
	err := tx.Select(`
			gm.id,
			gm.user_id,
			gm.create_time as join_time,
			u.user_no,
			u.nickname,
			u.avatar,
			gg.group_name,
			GROUP_CONCAT(distinct up.practitioner_type) AS id_cards,
			SUM(ob.diamond) AS profit_amount,
			COUNT(*) AS reward_count
		`).Group("gm.id,gm.user_id,join_time, u.user_no, u.nickname, u.avatar,gg.group_name").
		Order("profit_amount desc").Limit(req.Size).
		Offset((req.CurrentPage - 1) * req.Size).Scan(&result).Error

	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if len(result) > 0 && len(userRules) > 0 {
		var userRuleMap = map[string]model.AuthRoleAccess{}
		for _, v := range userRules {
			userRuleMap[v.UserID+","+cast.ToString(v.RoleID)] = v
		}
		for _, v := range result {
			if _, ok := userRuleMap[v.UserId+","+cast.ToString(enum.PresidentRoleId)]; ok {
				v.Identity = "会长"
				if _, ok := userRuleMap[v.UserId+","+cast.ToString(enum.HouseOwnerRoleId)]; ok {
					v.Identity += ",房主"
				}
			} else {
				if _, ok := userRuleMap[v.UserId+","+cast.ToString(enum.HouseOwnerRoleId)]; ok {
					v.Identity = "房主"
				} else {
					v.Identity = "成员"
				}
			}
			v.Avatar = coreConfig.GetHotConf().ImagePrefix + v.Avatar
		}
	}
	res.Data = result
	res.Size = req.Size
	res.CurrentPage = req.CurrentPage
	return
}
func (g *GuildMember) MemberIdcards(c *gin.Context, userId string) (res []*response_guild.MemberIdcardsInfoRes) {
	guildId := helper.GetGuildId(c)
	// 获取公会的所有房间ID
	roomDao := dao.RoomDao{}
	roomIds, _ := roomDao.GetRoomIdsByGuildId(guildId)
	coreDb.GetSlaveDb().Table("t_user_practitioner up").
		Joins("left join t_room tr on up.room_id=tr.id ").
		Joins("left join t_user u on u.id=tr.user_id").
		Where("tr.id in ? and up.status = 1 and up.user_id = ?", roomIds, userId).
		Select("tr.room_no,tr.id,tr.user_id,tr.name as room_name,up.practitioner_type as id_cards,up.create_time, u.user_no, u.nickname").Scan(&res)

	return
}

// 入会申请分页
func (g *GuildMember) MemberShipList(c *gin.Context, req *request_room.MemberShipListreq) (resp response.AdminPageRes) {
	guildId := helper.GetGuildId(c)
	var data []*response_guild.MemberJoinApplyInfo
	tx := coreDb.GetSlaveDb().Table("t_guild_member_apply gma").Joins("left join t_user u on u.id=gma.user_id").Where("gma.guild_id=? and gma.apply_type=1 and gma.status!=6", guildId)
	if len(req.UserKeyword) > 0 {
		tx = tx.Where("u.user_no like ? or u.nickname like ?", easy.GenLikeSql(req.UserKeyword), easy.GenLikeSql(req.UserKeyword))
	}
	if req.Status > 0 {
		tx = tx.Where("gma.status", req.Status)
	}
	tx.Count(&resp.Total)
	err := tx.Select("gma.id, gma.user_id, gma.status, gma.reason, gma.create_time, gma.update_time, u.user_no, u.nickname, u.avatar").
		Order("gma.status asc, gma.create_time desc").Limit(req.Size).Offset((req.CurrentPage - 1) * req.Size).Scan(&data).Error
	if err != nil {
		panic(i18n_err.ErrorCodeReadDB)
	}
	for _, info := range data {
		info.Avatar = helper.FormatImgUrl(info.Avatar)
		if info.Status == 1 {
			info.UpdateTime = easy.LocalTime(time.Time{})
		}
	}
	resp.Data = data
	resp.CurrentPage = req.CurrentPage
	resp.Size = req.Size
	return
}

// 退会申请分页
func (g *GuildMember) WithdrawMemberShipList(c *gin.Context, req *request_room.MemberShipListreq) (resp response.AdminPageRes) {
	guildId := helper.GetGuildId(c)
	var data []*response_guild.LeaveMemberShipListRsp
	tx := coreDb.GetSlaveDb().Table("t_guild_member_apply gma").Joins("left join t_user u on u.id=gma.user_id").
		Where("gma.guild_id=? and gma.apply_type=2 and gma.status!=6", guildId)
	if len(req.UserKeyword) > 0 {
		tx = tx.Where("u.user_no like ? or u.nickname like ?", easy.GenLikeSql(req.UserKeyword), easy.GenLikeSql(req.UserKeyword))
	}
	if req.Status > 0 {
		tx = tx.Where("gma.status", req.Status)
	}
	tx.Count(&resp.Total)
	err := tx.Select("gma.id, gma.user_id, gma.status, gma.force, gma.reason, gma.create_time, gma.update_time, u.user_no, u.nickname, u.avatar").
		Order("gma.status asc, gma.create_time desc").Limit(req.Size).Offset((req.CurrentPage - 1) * req.Size).Scan(&data).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	for _, info := range data {
		info.Avatar = helper.FormatImgUrl(info.Avatar)
		var starTime, endTime string
		switch info.Status {
		case 1, 3, 4: // 审核中 已拒绝 自动拒绝 在公会
			if info.Status == 1 {
				info.UpdateTime = easy.LocalTime(time.Time{})
			}
			// 查询公会成员信息
			memberInfo, err := new(dao.GuildDao).GetGuildMemberInfo(info.UserId)
			if err != nil {
				panic(i18n_err.I18nError{
					Code: i18n_err.ErrorCodeReadDB,
					Msg:  nil,
				})
			}
			info.JoinTime = easy.LocalTime(memberInfo.CreateTime)
			starTime = memberInfo.CreateTime.Format(time.DateTime)
			endTime = time.Now().Format(time.DateTime)
		case 2, 5: // 已同意 离开公会
			// 查询公会成员信息
			var memberInfo model.GuildMember
			err = coreDb.GetSlaveDb().Model(memberInfo).Where("create_time<? and status=3", info.CreateTime).Order("create_time desc").First(&memberInfo).Error
			if err != nil {
				panic(i18n_err.I18nError{
					Code: i18n_err.ErrorCodeReadDB,
					Msg:  nil,
				})
			}
			info.JoinTime = easy.LocalTime(memberInfo.CreateTime)
			starTime = memberInfo.CreateTime.Format(time.DateTime)
			endTime = time.Time(info.CreateTime).Format(time.DateTime)
		}
		// 查询用户流水
		_ = coreDb.GetSlaveDb().Model(model.OrderBill{}).Where("user_id=? and guild_id=?", info.UserId, guildId).
			Where("create_time between ? and ?", starTime, endTime).Select("IFNULL(sum(diamond),0) reward_diamonds, count(*) reward_num").Scan(&info)
	}
	resp.Data = data
	resp.CurrentPage = req.CurrentPage
	resp.Size = req.Size
	return
}

// GetPractitionerActionRecord
//
//	@Description: 查询用户从业者行为记录
//	@receiver g
//	@param c *gin.Context -
//	@return res -
func (g *GuildMember) GetPractitionerActionRecord(c *gin.Context) (res []*response_guild.UserPractitionerAction) {
	userId := c.GetString("userId")
	_ = coreDb.GetSlaveDb().Table("t_user_practitioner_action upa").Joins("left join t_user u on u.id=upa.user_id").Where("upa.user_id", userId).
		Select("upa.*, u.user_no, u.nickname, u.avatar").Order("upa.create_time desc").Scan(&res).Error
	for _, info := range res {
		info.Avatar = helper.FormatImgUrl(info.Avatar)
	}
	return
}

// 入会申请审核
func (g *GuildMember) MemberApplyReview(c *gin.Context, req *request_room.GuildMemberApplyReviewReq) (err error) {
	guildMemberApplyDao := new(dao2.GuildMemberApplyDao)
	//查询公会成员入会申请信息
	memberApplyInfo, err := guildMemberApplyDao.FindOne(&model.GuildMemberApply{
		ID:        req.Id,
		ApplyType: 1,
	})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if memberApplyInfo.ID == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeMemberApplyNotExist,
			Msg:  nil,
		})
	}
	memberApplyInfo.Status = req.Status
	if req.Status == typedef_enum.GuildMemberApplyStatusInactive {
		memberApplyInfo.Reason = req.Reason
	}
	//更新成员申请信息
	memberApplyUpdates := model.GuildMemberApply{
		ID:         memberApplyInfo.ID,
		Status:     memberApplyInfo.Status,
		Reason:     memberApplyInfo.Reason,
		UpdateTime: time.Now(),
	}
	tx := coreDb.GetMasterDb().Begin()
	err = tx.Model(model.GuildMemberApply{}).Where("id", memberApplyInfo.ID).Updates(&memberApplyUpdates).Error
	if err != nil {
		tx.Rollback()
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}
	if req.Status == typedef_enum.GuildMemberApplyStatusActive {
		//增加入会时星光等级和经验
		starLevelDao := dao.UserLevelStarDao{}
		starLevel, err := starLevelDao.GetUserStarLevel(memberApplyInfo.UserID)
		var (
			slevel = 1
			sExp   = 0
		)
		if err == nil {
			slevel = starLevel.Level
			sExp = starLevel.CurrExp
		}
		//添加公会成员
		guildMemberUpdates := model.GuildMember{
			GuildID:    memberApplyInfo.GuildID,
			UserID:     memberApplyInfo.UserID,
			Status:     typedef_enum.GuildMemberStatusNormal,
			StarLevel:  slevel,
			StarExp:    sExp,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		err = tx.Model(model.GuildMember{}).Create(&guildMemberUpdates).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
			})
		}
	}
	tx.Commit()
	return
}

// 退会申请审核
func (g *GuildMember) MemberWithdrawApplyReview(c *gin.Context, req *request_room.GuildMemberWithdrawReviewReq) (err error) {
	guildId := helper.GetGuildId(c)
	guildMemberApplyDao := new(dao2.GuildMemberApplyDao)
	//查询公会成员退会申请信息
	memberApplyInfo, err := guildMemberApplyDao.FindOne(&model.GuildMemberApply{
		ID: req.Id,
	})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if memberApplyInfo.ID == 0 || memberApplyInfo.ApplyType == 1 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeMemberApplyNotExist,
			Msg:  nil,
		})
	}
	if memberApplyInfo.Status != typedef_enum.GuildMemberApplyStatusWait {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeMemberApplyNotExist,
			Msg:  nil,
		})
	}
	// 强制退会申请 不能拒绝
	if memberApplyInfo.Force == 1 && req.Status == typedef_enum.GuildMemberApplyStatusInactive {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeParam,
			Msg:  nil,
		})
	}
	//查询成员公会信息
	userId := memberApplyInfo.UserID
	memberInfo, err := new(dao.GuildDao).GetGuildMemberInfo(userId)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if memberInfo.GuildID != guildId {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeParam,
			Msg:  nil,
		})
	}

	memberApplyInfo.Status = req.Status
	if req.Status == typedef_enum.GuildMemberApplyStatusInactive {
		memberApplyInfo.Reason = req.Reason
	}
	//更新工会成员信息
	memberApplyUpdates := model.GuildMemberApply{
		ID:         memberApplyInfo.ID,
		Status:     memberApplyInfo.Status,
		Reason:     memberApplyInfo.Reason,
		UpdateTime: time.Now(),
	}
	tx := coreDb.GetMasterDb().Begin()
	err = tx.Model(model.GuildMemberApply{}).Where("id", memberApplyInfo.ID).Updates(&memberApplyUpdates).Error
	if err != nil {
		tx.Rollback()
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUpdateDB,
			Msg:  nil,
		})
	}

	if req.Status == typedef_enum.GuildMemberApplyStatusActive {
		//如果成员有房主身份，提示"该成员有房主身份，请先取消对应身份后再进行操作！"
		roomInfo, er := new(dao.RoomDao).FindOne(&model.Room{
			UserId:   memberInfo.UserID,
			GuildId:  memberInfo.GuildID,
			LiveType: typedef_enum.LiveTypeChatroom,
		})
		if err != nil && !errors.Is(er, gorm.ErrRecordNotFound) {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeReadDB,
				Msg:  nil,
			})
		}
		if len(roomInfo.Id) > 0 {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUserIsRoomOwner,
				Msg:  nil,
			})
		}
		// 查询公会所有房间列表
		roomList, err := new(dao.RoomDao).GetRoomsByGuildId(guildId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
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
					panic(i18n_err.I18nError{
						Code: i18n_err.ErrorCodeUpdateDB,
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
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		// 删除玩家房间管理员
		err = tx.Model(model.RoomAdmin{}).Where("user_id=? and room_id in ?", userId, roomIdList).Delete(&model.RoomAdmin{}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
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
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
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
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		// 申请中的房间如果是房主自动拒绝
		err = tx.Model(model2.GuildRoomApply{}).Where("room_user_id=? and status=1 and guild_id=?", userId, guildId).Updates(map[string]interface{}{
			"status":      3,
			"reason":      "房主强制退出公会",
			"update_time": time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		// 退出公会
		err = tx.Model(model.GuildMember{}).Where("user_id=? and status != 3 and guild_id=?", userId, guildId).Updates(map[string]interface{}{
			"status":       3,
			"leave_reason": "强制退出公会",
			"update_time":  time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}

		// 扣除星光等级
		// 当前星光等级低于入会等级，保留当前等级不做处理
		var userStarLevel = new(model.UserStarLevel)
		err = tx.Model(model.UserStarLevel{}).Where("user_id = ?", userId).First(userStarLevel).Error
		if err != nil && err == gorm.ErrRecordNotFound { //如果没有星光等级记录，不处理
		} else if userStarLevel.Level >= memberInfo.StarLevel && userStarLevel.CurrExp > memberInfo.StarExp {
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
				panic(i18n_err.I18nError{
					Code: i18n_err.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
		}
	}
	tx.Commit()
	return
}

func (g *GuildMember) GuildKickoutMember(c *gin.Context, userId string, content string) (err error) {
	//踢出公会， 聊天室房主,会长不能踢出， 其余可以踢出， 从业者踢出时，从业者身份取消，主播剔除，房间取消关联公会
	guildId := c.GetString("guildId")

	guildDao := dao2.GuildDao{}
	gFirst, err := guildDao.FindOne(&model.Guild{ID: guildId})
	if err != nil {
		panic(i18n_err.ErrorCodeDataNotFound)
	}
	if gFirst.UserID == userId {
		panic(i18n_err.ErrorCodeRoomGuildOwner)
	}

	guildUserDao := dao2.GuildMemberDao{}
	first := &model.GuildMember{
		GuildID: guildId,
		UserID:  userId,
		Status:  1,
	}
	first, err = guildUserDao.FindOne(first)
	if err != nil {
		panic(i18n_err.ErrorCodeDataNotFound)
	}

	//查询出工会的所有聊天室房间ID
	roomDao := dao.RoomDao{}
	roomList, err := roomDao.GetRoomsByGuildId(guildId)
	if err != nil {
		panic(i18n_err.ErrorCodeRoomNotExist)
	}
	tx := coreDb.GetMasterDb().Begin()
	var ids []string
	for _, v := range roomList {
		if v.UserId == userId && v.LiveType == enum.LiveTypeAnchor { // 是主播且已并入公会
			// 主播房间和公会解约
			err = tx.Model(model.Room{}).Where("user_id=? and room_id=?", userId, v.Id).Updates(map[string]interface{}{
				"guild_id":    "",
				"update_time": time.Now(),
			}).Error
			if err != nil {
				tx.Rollback()
				panic(i18n_err.I18nError{
					Code: i18n_err.ErrorCodeUpdateDB,
					Msg:  nil,
				})
			}
			continue
		}
		ids = append(ids, v.Id)
		if v.UserId == userId && v.LiveType == enum.LiveTypeChatroom {
			tx.Rollback()
			panic(i18n_err.ErrorCodeRoomOwner)
		}
	}
	//下掉公会成员身份
	first.Status = 3
	first.UpdateTime = time.Now()
	first.LeaveReason = content
	first.LeaveTime = &first.UpdateTime
	err = tx.Save(&first).Error
	if err != nil {
		tx.Rollback()
		panic(i18n_err.ErrorCodeUpdateDB)
	}
	if len(ids) > 0 {
		//下掉从业者身份
		err = tx.Model(model.UserPractitioner{}).Where("user_id = ? and room_id in ? and status = 1", userId, ids).UpdateColumns(map[string]interface{}{
			"status":         4,
			"abolish_reason": "踢出公会",
		}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.ErrorCodeUpdateDB)
		}
		//驳回从业者审核记录
		err = tx.Model(model.UserPractitioner{}).Where("user_id = ? and room_id in ? and status = 2", userId, ids).UpdateColumns(map[string]interface{}{
			"status":         3,
			"abolish_reason": "踢出公会",
		}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.ErrorCodeUpdateDB)
		}
		//TODO:  可能会删除掉当前用户对于这些房间的超管，巡查权限，核实需不需要删除
		err = tx.Model(model.AuthRoleAccess{}).Where("user_id = ? and room_id in ?", userId, ids).Delete(&model.AuthRoleAccess{}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.ErrorCodeUpdateDB)
		}
		//删除用户对于其他房间的管理员身份
		err = tx.Model(model.RoomAdmin{}).Where("user_id = ? and room_id in ?", userId, ids).Delete(&model.RoomAdmin{}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.ErrorCodeUpdateDB)
		}
		// 申请中的房间如果是房主自动拒绝
		err = tx.Model(model2.GuildRoomApply{}).Where("room_user_id=? and status=1 and guild_id=?", userId, guildId).Updates(map[string]interface{}{
			"status":      3,
			"reason":      "踢出公会",
			"update_time": time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		// 如果有待审核的申请退会记录，自动拒绝
		err = tx.Model(model.GuildMemberApply{}).Where("user_id=? and guild_id=? and apply_type=2 and status=1", userId, guildId).Updates(map[string]interface{}{
			"status":      4,
			"reason":      "踢出公会",
			"update_time": time.Now(),
		}).Error
		if err != nil {
			tx.Rollback()
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUpdateDB,
				Msg:  nil,
			})
		}
		//执行星光降级
		var userStarLevel = new(model.UserStarLevel)
		err = tx.Model(model.UserStarLevel{}).Where("user_id = ?", userId).First(userStarLevel).Error
		if err != nil && err == gorm.ErrRecordNotFound { //如果没有星光等级记录，不处理
		} else {
			if userStarLevel.Level >= first.StarLevel && userStarLevel.CurrExp > first.StarExp { //降级重新计算有效期，并且保级经验清零
				userStarLevel.ExpireTime = time.Now().AddDate(0, 0, 30)
				userStarLevel.KeepExp = 0
				userStarLevel.Level = first.StarLevel
				userStarLevel.CurrExp = first.StarExp
				userStarLevel.UpdateTime = time.Now()
				err = tx.Model(userStarLevel).Save(userStarLevel).Error
				if err != nil {
					tx.Rollback()
					panic(i18n_err.ErrorCodeUpdateDB)
				}
			} else { // 当前星光等级低于入会等级，保留当前等级不做处理
			}
		}
		//清除权限缓存
		var keys []string
		for _, roomId := range ids {
			key1 := redisKey.UserRules(userId, roomId)
			key2 := redisKey.UserRoles(userId, roomId)
			key3 := redisKey.UserCompereRules(userId, roomId)
			keys = append(keys, key1, key2, key3)
		}
		coreRedis.GetUserRedis().Del(c, keys...)
	}
	tx.Commit()
	return
}
