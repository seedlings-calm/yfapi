package logic

import (
	"strings"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreJwtToken"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/riskCheck/shumei"
	service_user "yfapi/internal/service/user"
	typedef_enum "yfapi/typedef/enum"
	common_data "yfapi/typedef/redisKey"
	request_user "yfapi/typedef/request/user"
	response_login "yfapi/typedef/response/login"
	"yfapi/typedef/response/user"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type UserSetting struct {
}

func (u *UserSetting) SetPassword(req *request_user.SetPasswordReq, c *gin.Context) {
	userId := handle.GetUserId(c)
	userDao := new(dao.UserDao)
	user, _ := userDao.FindOne(&model.User{
		Id: userId,
	})
	if len(user.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	req.Password = strings.ToLower(req.Password)
	if helper.CheckPasswordLever(req.Password) != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodePassEasy,
			Msg:  nil,
		})
	}
	oldPass := user.Password
	newPass := helper.EncodeUserPass(req.Password)
	if oldPass == newPass {
		panic(error2.I18nError{
			Code: error2.ErrorCodePassEq,
			Msg:  nil,
		})
	}
	repeat := new(dao.UserDao).RepeatPassword(user.Mobile, newPass)
	if repeat {
		panic(error2.ErrCodePasswordRepeat)
	}
	//检测验证码
	sms := &Sms{
		Mobile:     user.Mobile,
		Code:       req.Code,
		RegionCode: user.RegionCode,
		Type:       typedef_enum.SmsCodeSetPass,
	}
	err := sms.CheckSms(c)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeCaptchaInvalid,
			Msg:  nil,
		})
	}

	userDao.UpdateById(&model.User{
		Id:       user.Id,
		Password: newPass,
	})
	return
}

// 实名认证
func (u *UserSetting) RealName(c *gin.Context, req *request_user.UserRealNameReq) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{
		Id: userId,
	})
	if err != nil || len(userModel.Id) == 0 {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	//添加实名认证是否是待审核或审核通过
	isRealName, err := new(dao.UserDao).FindUserIsRealName(userModel.Id)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeReadDB})
	}
	if len(isRealName) > 0 {
		panic(error2.I18nError{Code: error2.ErrorCodeAlreadyApplyRealName})
	}
	userModel.RealNameStatus = typedef_enum.UserRealNameVerifying
	tx := coreDb.GetMasterDb().Begin()
	err = tx.Model(&model.User{}).Where("id", userId).Updates(userModel).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{Code: error2.ErrCodeUserRealName})
	}
	realNameInfo := &model.UserRealName{
		UserId:     userModel.Id,
		TrueName:   req.RealName,
		IdNo:       req.IdNo,
		FontUrl:    req.FontUrl,
		BackUrl:    req.BackUrl,
		Status:     typedef_enum.UserRealNameReviewStatusWait,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	err = tx.Model(&model.UserRealName{}).Create(realNameInfo).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{Code: error2.ErrCodeUserRealName})
	}
	tx.Commit()
	return
}

// 实名认证状态
func (u *UserSetting) RealNameStatus(c *gin.Context) (res user.RealNameResp) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{
		Id: userId,
	})
	if err != nil || len(userModel.Id) == 0 {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	res.RealNameStatus = userModel.RealNameStatus
	return
}

// 获取当前用户手机号绑定的所有账号
func (u *UserSetting) GetAccountByMobile(c *gin.Context) (res []user.UserAccount) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{
		Id: userId,
	})
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	users := new(dao.UserDao).GetUsersByMobile(userModel.RegionCode, userModel.Mobile)
	for _, u := range users {
		info := user.UserAccount{
			Id:       u.Id,
			UserNo:   u.UserNo,
			Uid32:    cast.ToInt32(u.OriUserNo),
			Avatar:   helper.FormatImgUrl(u.Avatar),
			Nickname: u.Nickname,
		}
		res = append(res, info)
	}
	return
}

// 创建账号检测
func (u *UserSetting) CheckCreateNewAccount(c *gin.Context) (res user.UserAccount) {
	userId := helper.GetUserId(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{
		Id: userId,
	})
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	count := new(dao.UserDao).Count(&model.User{RegionCode: userModel.RegionCode, Mobile: userModel.Mobile})
	if count >= 3 {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserAccountExceed,
			Msg: map[string]interface{}{
				"num": 3,
			},
		})
	}
	res.Avatar = helper.FormatImgUrl(helper.GetUserDefaultAvatar(userModel.Sex))
	res.Nickname = "用户" + helper.GeneUserNickname()
	return
}

// 创建新账号
func (u *UserSetting) CreateNewAccount(c *gin.Context, req *request_user.UserCreateAccountReq) (res user.UserAccount) {
	userId := helper.GetUserId(c)
	headerData := helper.GetHeaderData(c)
	userModel, err := new(dao.UserDao).FindOne(&model.User{
		Id: userId,
	})
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	count := new(dao.UserDao).Count(&model.User{RegionCode: userModel.RegionCode, Mobile: userModel.Mobile})
	if count >= 3 {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserAccountExceed,
			Msg: map[string]interface{}{
				"num": 3,
			},
		})
	}
	if headerData.Platform != "android" && headerData.Platform != "ios" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	if len(req.Avatar) > 0 {
		req.Avatar = helper.RemovePrefixImgUrl(req.Avatar)
		if ok := new(shumei.ShuMei).AvatarSyncCheck(userId, helper.FormatImgUrl(req.Avatar)); !ok {
			panic(error2.I18nError{
				Code: error2.ErrorCodeAvatarCheckReject,
				Msg:  nil,
			})
		}
	}
	if new(dao.UserDao).CheckRepeatNickname(req.Nickname) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNicknameRepeat,
			Msg:  nil,
		})
	}
	newUser := new(model.User)
	newUser.Id = coreSnowflake.GetSnowId()
	if ok := new(shumei.ShuMei).NicknameCheck(newUser.Id, req.Nickname); !ok {
		panic(error2.I18nError{
			Code: error2.ErrorCodeNicknameCheckReject,
			Msg:  nil,
		})
	}

	newUser.Mobile = userModel.Mobile
	newUser.RegionCode = userModel.RegionCode
	loginIp := c.ClientIP()
	userNo := helper.GeneUserNo()
	newUser.UserNo = userNo
	newUser.OriUserNo = userNo
	newUser.Status = typedef_enum.UserStatusNormal
	newUser.Nickname = req.Nickname
	newUser.Sex = typedef_enum.UserSexTypeWoman
	newUser.Avatar = helper.RemovePrefixImgUrl(req.Avatar)
	newUser.RegisterChannel = headerData.Channel
	newUser.RegisterPlatform = headerData.Platform
	newUser.RegisterMachineCode = headerData.MachineCode
	newUser.CreateTime = time.Now()
	newUser.Guide = 2
	if len(req.Password) > 0 {
		req.Password = strings.ToLower(req.Password)
		if helper.CheckPasswordLever(req.Password) != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodePassEasy,
				Msg:  nil,
			})
		}
		repeat := new(dao.UserDao).RepeatPassword(userModel.Mobile, helper.EncodeUserPass(req.Password))
		if repeat {
			panic(error2.I18nError{Code: error2.ErrCodePasswordRepeat})
		}
		newUser.Password = helper.EncodeUserPass(req.Password)
	}
	tx := coreDb.GetMasterDb().Begin()
	//注册用户
	err = tx.Model(newUser).Create(newUser).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 注册用户账户信息
	err = tx.Exec("insert into t_user_account(user_id) values(?)", newUser.Id).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	//记录注册信息
	registerRecordModel := &model.UserRegisterRecord{
		UserId:           newUser.Id,
		RegisterPlatform: headerData.Platform,
		SignType:         2,
		RegisterChannel:  headerData.Channel,
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
	err = tx.Model(&model.UserRegisterRecord{}).Create(registerRecordModel).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	//添加登录信息
	loginRecordModel := &model.UserLoginRecord{
		UserId:        newUser.Id,
		LoginPlatform: headerData.Platform,
		ClientVersion: headerData.AppVersion,
		DeviceID:      headerData.MachineCode,
		LoginIp:       loginIp,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	err = tx.Model(&model.UserLoginRecord{}).Create(loginRecordModel).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	tx.Commit()
	res.Avatar = helper.FormatImgUrl(newUser.Avatar)
	res.Nickname = newUser.Nickname
	res.Id = newUser.Id
	res.UserNo = newUser.UserNo
	res.Uid32 = cast.ToInt32(newUser.OriUserNo)
	token, err := coreJwtToken.Encode(coreJwtToken.Claims{
		UserId:         newUser.Id,
		Mobile:         newUser.Mobile,
		ClientType:     headerData.Platform, //Android Ios Pc Web H5
		DeviceID:       headerData.MachineCode,
		StandardClaims: gojwt.StandardClaims{},
	}, []byte(coreConfig.GetHotConf().JwtSecret), time.Now().Add(typedef_enum.TokenExpireLifeTime).Unix())
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	res.Token = token
	res.Mobile = helper.PrivateMobile(newUser.Mobile)
	return
}

// 选择账号登录
func (u *UserSetting) SwitchAccount(c *gin.Context, req *request_user.SwitchUserAccountReq) (res response_login.LoginCodeRes) {
	userHeaderData := helper.GetHeaderData(c)
	claims, err := coreJwtToken.Decode(req.Token, []byte(coreConfig.GetHotConf().JwtSecret))
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeSwitchUserAccountTokenExpire,
		})
	}
	coreLog.LogInfo("userLoginRecord claims %+v", claims)
	switchUserId := claims.UserId
	switchUserDeviceId := claims.DeviceID
	if userHeaderData.MachineCode != switchUserDeviceId {
		panic(error2.I18nError{
			Code: error2.ErrCodeSwitchUserAccountTokenExpire,
		})
	}
	userLoginRecordModel := new(dao.UserLoginRecordDao).LastOneByUserId(switchUserId)
	coreLog.Info("userLoginRecord %+v", userLoginRecordModel)
	if userLoginRecordModel != nil {
		if userLoginRecordModel.DeviceID != switchUserDeviceId {
			panic(error2.I18nError{
				Code: error2.ErrCodeSwitchUserAccountTokenExpire,
			})
		}
	}
	userModel, err := new(dao.UserDao).FindOne(&model.User{
		Id: switchUserId,
	})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
		})
	}
	if userModel.Status == typedef_enum.UserStatusFreezing || userModel.Status == typedef_enum.UserStatusInvalid {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
		})
	}
	new(dao.UserLoginRecordDao).LoginRecord(c, userModel, userHeaderData)
	token, err := coreJwtToken.Encode(coreJwtToken.Claims{
		UserId:         userModel.Id,
		Mobile:         userModel.Mobile,
		ClientType:     userHeaderData.Platform, //Android Ios Pc Web H5
		DeviceID:       userHeaderData.MachineCode,
		StandardClaims: gojwt.StandardClaims{},
	}, []byte(coreConfig.GetHotConf().JwtSecret), time.Now().Add(typedef_enum.TokenExpireLifeTime).Unix())
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	res.JumpType = typedef_enum.LoginJumpTypeHome
	res.Id = userModel.Id
	res.Token = token
	switch userHeaderData.Platform {
	case typedef_enum.ClientTypeAndroid, typedef_enum.ClientTypeIos: //app
		coreRedis.GetUserRedis().Set(c, common_data.UserLoginInfo("app", userModel.Id), token, typedef_enum.RedisTokenExpireLifeTime)
	case typedef_enum.ClientTypePc: //pc
		coreRedis.GetUserRedis().Set(c, common_data.UserLoginInfo("pc", userModel.Id), token, typedef_enum.RedisTokenExpireLifeTime)
	case typedef_enum.ClientTypeH5: //h5
		coreRedis.GetUserRedis().Set(c, common_data.UserLoginInfo("h5", userModel.Id), token, typedef_enum.RedisTokenExpireLifeTime)
	}
	//如果注销中取消注销
	service_user.UnCancelAccount(userModel.Id)
	return
}

// 获取用户隐私信息
func (u *UserSetting) GetUserPrivateInfo(c *gin.Context) (res user.UserPrivateInfo) {
	userId := helper.GetUserId(c)
	userModel := service_user.GetUserBaseInfo(userId)
	res.Nickname = userModel.Nickname
	res.Id = userModel.Id
	res.Mobile = helper.PrivateMobile(userModel.Mobile)
	one := new(dao.UserRealNameDao).FindOne(&model.UserRealName{
		UserId: userId,
	})
	if one.Id == 0 {
		res.RealNameStatus = typedef_enum.UserRealName(typedef_enum.UserRealNameUnverified).String()
	} else {
		if one.Status == typedef_enum.UserRealNameReviewStatusWait {
			res.RealNameStatus = typedef_enum.UserRealNameReviewStatus(typedef_enum.UserRealNameReviewStatusWait).String()
		} else if one.Status == typedef_enum.UserRealNameReviewStatusReject {
			res.RealNameStatus = typedef_enum.UserRealName(typedef_enum.UserRealNameUnverified).String()
		} else {
			res.RealNameStatus = typedef_enum.UserRealName(typedef_enum.UserRealNameAuthenticated).String()
			res.RealName = helper.PrivateRealName(one.TrueName)
			res.CardNum = helper.PrivateIdNo(one.IdNo)
		}
	}
	return
}

// 发送验证码校验是否原始手机号
func (u *UserSetting) SendVerifyMobileCode(c *gin.Context) {
	userId := helper.GetUserId(c)
	userInfo, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	ser := &Sms{
		Mobile:     userInfo.Mobile,
		Type:       typedef_enum.SmsCodeVerifyMobile,
		RegionCode: userInfo.RegionCode,
	}
	err = ser.SendSms(c)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUnknown})
	}
}

// 验证手机号
func (u *UserSetting) VerifyMobile(c *gin.Context, req *request_user.VerifyMobileReq) {
	userId := helper.GetUserId(c)
	userInfo, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	ser := &Sms{
		Mobile:     userInfo.Mobile,
		Type:       typedef_enum.SmsCodeVerifyMobile,
		Code:       req.Code,
		RegionCode: userInfo.RegionCode,
	}
	err = ser.CheckSms(c)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeCaptchaInvalid})
	}
	key := common_data.ChangeUserMobileBeforeVerify(userId)
	coreRedis.GetUserRedis().Set(c, key, 1, time.Second*300)
}

// 更换手机号
func (u *UserSetting) ChangeMobile(c *gin.Context, req *request_user.ChangeMobileReq) {
	userId := helper.GetUserId(c)
	userInfo, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserNotFound})
	}
	if userInfo.RegionCode+userInfo.Mobile == req.RegionCode+req.Mobile {
		panic(error2.I18nError{Code: error2.ErrorCodeNewMobileOldMobileSame})
	}
	ser := &Sms{
		Mobile:     req.Mobile,
		Type:       typedef_enum.SmsCodeChangeMobile,
		Code:       req.Code,
		RegionCode: req.RegionCode,
	}
	err = ser.CheckSms(c)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeCaptchaInvalid})
	}
	key := common_data.ChangeUserMobileBeforeVerify(userId)
	ok, _ := coreRedis.GetUserRedis().Exists(c, key).Result()
	if ok != 1 {
		panic(error2.I18nError{Code: error2.ErrorCodeSystemBusy})
	}
	err = new(dao.UserDao).UpdateMobile(userInfo.RegionCode, userInfo.Mobile, req.RegionCode, req.Mobile)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeSystemBusy})
	}
	coreRedis.GetUserRedis().Del(c, key)
}

func (u *UserSetting) GetLoginLog(c *gin.Context) (resp []*user.LoginRecordResponse) {
	resp = make([]*user.LoginRecordResponse, 0)
	userId := handle.GetUserId(c)
	recordDao := dao.UserLoginRecordDao{}
	res := recordDao.FindByUserId(userId)
	if len(res) > 0 {
		for _, v := range res {
			item := &user.LoginRecordResponse{
				LoginPlatform: v.LoginPlatform,
				LoginModel:    v.LoginModel,
				ClientVersion: v.ClientVersion,
				DeviceID:      v.DeviceID,
				CreateTime:    v.CreateTime.Format(time.DateTime),
				Address:       v.Address,
			}
			resp = append(resp, item)
		}
	}
	return
}
