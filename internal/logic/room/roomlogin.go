package room

import (
	"encoding/json"
	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"sort"
	"strconv"
	"strings"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreJwtToken"
	"yfapi/core/coreRedis"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/logic"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	typedef_enum "yfapi/typedef/enum"
	common_data "yfapi/typedef/redisKey"
	request_login "yfapi/typedef/request/roomOwner"
	"yfapi/typedef/response"
	response_login "yfapi/typedef/response/roomOwner"
	"yfapi/util/easy"
)

type RoomLogin struct {
}

// 手机验证码登录
func (g *RoomLogin) LoginByCode(req *request_login.LoginMobileReq, context *gin.Context) (res response_login.LoginMobileCodeRes) {
	//检测验证码
	sms := &logic.Sms{
		Mobile:     req.Mobile,
		Code:       req.Code,
		RegionCode: req.RegionCode,
		Type:       typedef_enum.SmsCodeRoomAdminLogin,
	}
	//data := helper.GetHeaderData(context)
	//校验手机验证码
	err := sms.CheckSms(context)
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeCaptchaInvalid,
			Msg:  nil,
		})
	}

	//判断用户合法性
	one, err := new(dao.UserDao).FindOne(&model.User{Mobile: req.Mobile})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if one.Status == 2 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserFreezing,
			Msg:  nil,
		})
	}
	//生成token
	token, err := coreJwtToken.RoomEncode(coreJwtToken.RoomClaims{
		UserId:         one.Id,
		Mobile:         one.Mobile,
		StandardClaims: gojwt.StandardClaims{},
	}, []byte(coreConfig.GetHotConf().JwtSecret), time.Now().Add(time.Hour*24*7).Unix())
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	//将token放入登录缓存
	//需要根据客户端类型进行判定
	coreRedis.GetUserRedis().Set(context, common_data.UserLoginInfo("roomPc", one.Id), token, 7*24*time.Hour)

	// 记录本次登录时间
	key := common_data.RoomAdminLoginTime(one.Id)
	var recordList []string
	now := time.Now().Format(time.DateTime)
	recordStr := coreRedis.GetUserRedis().Get(context, key).Val()
	if len(recordStr) == 0 {
		recordList = append(recordList, now)
	} else {
		_ = json.Unmarshal([]byte(recordStr), &recordList)
		recordList = append(recordList, now)
		if len(recordList) > 2 {
			recordList = recordList[1:]
		}
	}
	recordByte, _ := json.Marshal(recordList)
	coreRedis.GetUserRedis().Set(context, key, string(recordByte), 0)

	return response_login.LoginMobileCodeRes{
		Token:  token,
		UserID: one.Id,
	}

}

// 获取手机验证码
func (g *RoomLogin) SearchByMobile(mobile string, context *gin.Context) (err error) {
	one, err := new(dao.UserDao).FindOne(&model.User{Mobile: mobile})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if one.Status == 2 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserFreezing,
			Msg:  nil,
		})
	}

	roomlist, err := new(dao.RoomDao).FindList(&model.Room{UserId: one.Id})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if len(roomlist) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	return err
}

// 获取单个房间信息
func (g *RoomLogin) RoomInfo(c *gin.Context) (resp response_login.RoomHomeInfo) {
	rid := c.GetHeader("roomId")
	if rid == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeNoRoomID,
			Msg:  nil,
		})
	}

	data, _ := new(dao.RoomDao).FindOne(&model.Room{Id: rid})
	one, err := new(dao.UserDao).FindOne(&model.User{Id: data.UserId})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}

	// 上次登录时间
	lastLoginTime := ""
	key := common_data.RoomAdminLoginTime(data.UserId)
	var recordList []string
	recordStr := coreRedis.GetUserRedis().Get(c, key).Val()
	if len(recordStr) > 0 {
		_ = json.Unmarshal([]byte(recordStr), &recordList)
		if len(recordList) > 1 {
			lastLoginTime = recordList[0]
		}
	}

	return response_login.RoomHomeInfo{
		RoomName:      data.Name,
		UserAvatar:    coreConfig.GetHotConf().ImagePrefix + one.Avatar,
		UserNo:        one.UserNo,
		UserName:      one.Nickname,
		LogoImg:       coreConfig.GetHotConf().ImagePrefix + data.CoverImg,
		RoomType:      strconv.Itoa(data.RoomType),
		CreateTime:    data.CreateTime.Format("2006-01-02 15:04:05"),
		Notice:        data.Notice,
		RoomNo:        data.RoomNo,
		Welcome:       "",
		RoomID:        data.Id,
		UserID:        data.UserId,
		LastLoginTime: lastLoginTime,
	}

}

// 聊天室概况
func (g *RoomLogin) RoomBase(c *gin.Context) (resp response_login.RoomHomeBaseInfo) {
	roomId := helper.GetRoomId(c)
	if len(roomId) == 0 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeNoRoomID,
			Msg:  nil,
		})
	}

	// 今日开播时长 本日进房人次
	coreDb.GetSlaveDb().Table("t_room_wheat_time").Where("stat_date=? and room_id=?", time.Now().Format(time.DateOnly), roomId).Select("IFNULL(sum(on_time),0) online_second, IFNULL(sum(enter_times),0) enter_times").Scan(&resp)
	resp.OnlineSecond = easy.SecondFormatString(cast.ToInt64(resp.OnlineSecond))
	// 今日流水
	coreDb.GetSlaveDb().Table("t_order").Where("room_id=?", roomId).
		Where("stat_date=? and order_type=?", time.Now().Format(time.DateOnly), accountBook.ChangeDiamondRewardGift).Select("IFNULL(sum(total_amount),0)").Scan(&resp.TodayProfit)
	// 本月流水
	coreDb.GetSlaveDb().Table("t_order").Where("room_id=? and order_type=?", roomId, accountBook.ChangeDiamondRewardGift).
		Where("stat_date between ? and ?", easy.GetCurrMonthStartTime(time.Now()).Format(time.DateOnly), time.Now().Format(time.DateOnly)).
		Select("IFNULL(sum(total_amount),0)").Scan(&resp.MonthProfit)
	// 从业者 主持 音乐人 咨询师
	// 从业者数量
	var data []struct {
		UserId   string
		TypeList string
	}
	coreDb.GetSlaveDb().Table("t_user_practitioner").Where("room_id=? and status=1", roomId).
		Select("user_id, GROUP_CONCAT(DISTINCT practitioner_type SEPARATOR ',') type_list").Group("user_id").Scan(&data)
	resp.Practitioner = len(data)
	for _, info := range data {
		dst := strings.Split(info.TypeList, ",")
		for _, _type := range dst {
			switch cast.ToInt(_type) {
			case typedef_enum.UserPractitionerCompere:
				resp.Host++
			case typedef_enum.UserPractitionerMusician:
				resp.Musician++
			case typedef_enum.UserPractitionerCounselor:
				resp.Counselor++
			}
		}
	}

	// 最近七天的流水信息
	startTime := time.Now().AddDate(0, 0, -6)
	coreDb.GetSlaveDb().Model(model.Order{}).Where("stat_date between ? and ?", startTime.Format(time.DateOnly), time.Now().Format(time.DateOnly)).
		Where("room_id=? and order_type=?", roomId, accountBook.ChangeDiamondRewardGift).
		Select("stat_date, IFNULL(sum(total_amount),0) profit_amount").Group("stat_date").Order("stat_date").Scan(&resp.LatestWeek)
	weekMap := make(map[string]struct{})
	for i, info := range resp.LatestWeek {
		currDate, _ := time.ParseInLocation(time.RFC3339, info.StatDate, time.Local)
		info.StatDate = currDate.Format(time.DateOnly)
		resp.LatestWeek[i].StatDate = info.StatDate
		weekMap[info.StatDate] = struct{}{}
	}
	for i := 0; i < 7; i++ {
		currDate := startTime.AddDate(0, 0, i).Format(time.DateOnly)
		if _, isExist := weekMap[currDate]; !isExist {
			resp.LatestWeek = append(resp.LatestWeek, response_login.ProfitInfo{
				StatDate:     currDate,
				ProfitAmount: "0",
			})
		}
	}
	// 排序
	sort.Slice(resp.LatestWeek, func(i, j int) bool {
		iDate, _ := time.ParseInLocation(time.DateOnly, resp.LatestWeek[i].StatDate, time.Local)
		jDate, _ := time.ParseInLocation(time.DateOnly, resp.LatestWeek[j].StatDate, time.Local)
		return iDate.Before(jDate)
	})
	return
}

// 房间管理员列表
func (g *RoomLogin) RoomAdminList(req *request_login.RoomAdminListReq, c *gin.Context) (resp response.AdminPageRes) {
	roomId := helper.GetRoomId(c)
	userId := helper.GetUserId(c)
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	if err != nil || room.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	//判断房间状态
	if room.Status != 1 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	// 判断操作人是否是房主
	if room.UserId != userId {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeNotRoomOwner,
			Msg:  nil,
		})
	}
	//查询房间管理员列表
	var data []*response_login.RoomAdminInfo
	limit := req.Size
	offset := req.Size * (req.CurrentPage - 1)
	db := coreDb.GetSlaveDb().Table("t_room_admin ra").Joins("left join t_user u on u.id=ra.user_id").Where("ra.room_id=?", roomId)
	_ = db.Count(&resp.Total)
	err = db.Select("ra.user_id, ra.create_time, ra.staff_name, u.user_no, u.nickname, u.avatar").Order("ra.create_time desc").Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	for _, info := range data {
		info.Avatar = helper.FormatImgUrl(info.Avatar)
	}
	resp.Data = data
	resp.CurrentPage = req.CurrentPage
	resp.Size = req.Size
	return

}

// 从业者列表
func (g *RoomLogin) RoomPractitionerList(req *request_login.RoomPractitionerListReq, c *gin.Context) (resp response.AdminPageRes) {
	roomId := helper.GetRoomId(c)
	// 判断房间是否存在
	room, err := new(dao.RoomDao).FindOne(&model.Room{Id: roomId})
	if err != nil || room.Id == "" {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomNotExist,
			Msg:  nil,
		})
	}
	//判断房间状态
	if room.Status != 1 {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeRoomStatus,
			Msg:  nil,
		})
	}
	//查询房间从业者列表
	var data []*response_login.RoomPractitionerInfo
	tx := coreDb.GetSlaveDb().Table("t_user_practitioner up").Joins("left join t_user u on u.id=up.user_id").Where("up.room_id=? and up.status!=4", roomId)
	if len(req.UserKeyword) > 0 {
		tx = tx.Where("u.user_no like ? or u.nickname like ?", easy.GenLikeSql(req.UserKeyword), easy.GenLikeSql(req.UserKeyword))
	}
	if req.PractitionerType > 0 {
		tx = tx.Where("up.practitioner_type", req.PractitionerType)
	}
	if req.Status > 0 {
		tx = tx.Where("up.status", req.Status)
	}
	tx.Count(&resp.Total)
	err = tx.Select("up.id, up.user_id, up.practitioner_type, up.status, up.abolish_reason, up.create_time, up.update_time, u.user_no, u.nickname, u.avatar").
		Order("up.create_time desc").Limit(req.Size).Offset((req.CurrentPage - 1) * req.Size).Find(&data).Error
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	for _, info := range data {
		info.Avatar = helper.FormatImgUrl(info.Avatar)
		info.PractitionerTypeDesc = typedef_enum.PractitionerType(info.PractitionerType).String()
	}

	resp.Data = data
	resp.CurrentPage = req.CurrentPage
	resp.Size = req.Size
	return

}

// 获取房间列表
func (g *RoomLogin) RoomList(c *gin.Context) (resp []*response_login.RoomInfo) {
	userId := helper.GetUserId(c)
	data, _ := new(dao.RoomDao).FindListByLiveType(userId, typedef_enum.LiveTypeChatroom)
	for _, v := range data {
		one, err := new(dao.UserDao).FindOne(&model.User{Id: v.UserId})
		if err != nil {
			panic(i18n_err.I18nError{
				Code: i18n_err.ErrorCodeUserNotFound,
				Msg:  nil,
			})
		}
		resp = append(resp, &response_login.RoomInfo{
			RoomName:   v.Name,
			UserAvatar: coreConfig.GetHotConf().ImagePrefix + one.Avatar,
			UserNo:     one.UserNo,
			UserName:   one.Nickname,
			LogoImg:    coreConfig.GetHotConf().ImagePrefix + v.CoverImg,
			RoomNo:     v.RoomNo,
			RoomID:     v.Id,
		})
	}
	return resp

}
