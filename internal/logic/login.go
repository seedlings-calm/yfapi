package logic

import (
	"github.com/spf13/cast"
	"strings"
	"time"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreJwtToken"
	"yfapi/core/coreRedis"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	service_im "yfapi/internal/service/im"
	service_user "yfapi/internal/service/user"
	"yfapi/typedef"
	typedef_enum "yfapi/typedef/enum"
	common_data "yfapi/typedef/redisKey"
	request_login "yfapi/typedef/request/login"
	response_login "yfapi/typedef/response/login"
	"yfapi/util/easy"

	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Login struct {
}

// 检测用户登陆前是否在其他端登录
func (l *Login) CheckLoginType(c *gin.Context, req *request_login.LoginCheckReq) {
	otherClientLogin := false
	headerData := helper.GetHeaderData(c)
	userDao := new(dao.UserDao)
	//手机号密码登录 直接检测
	if len(req.Mobile) != 0 && len(req.RegionCode) != 0 && len(req.Password) != 0 {
		user, err := userDao.FindOne(&model.User{
			Mobile:     req.Mobile,
			RegionCode: req.RegionCode,
			Password:   helper.EncodeUserPass(strings.ToLower(req.Password)),
		})
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserPassError,
				Msg:  nil,
			})
		}
		if len(user.Id) == 0 {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserPassError,
				Msg:  nil,
			})
		}
		clientType := service_user.GetUserLoginClientType(user.Id)
		switch headerData.Platform {
		case typedef_enum.ClientTypeAndroid:
			if easy.InArray(typedef_enum.ClientTypeIos, clientType) {
				otherClientLogin = true
			}
			if easy.InArray(typedef_enum.ClientTypeAndroid, clientType) {
				otherClientLogin = true
			}
		case typedef_enum.ClientTypeIos:
			if easy.InArray(typedef_enum.ClientTypeAndroid, clientType) {
				otherClientLogin = true
			}
			if easy.InArray(typedef_enum.ClientTypeIos, clientType) {
				otherClientLogin = true
			}
		}
		if otherClientLogin {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserLoginOtherClient,
			})
		}
		return
	}

	if len(req.UserId) > 0 {
		clientType := service_user.GetUserLoginClientType(req.UserId)
		switch headerData.Platform {
		case typedef_enum.ClientTypeAndroid:
			if easy.InArray(typedef_enum.ClientTypeIos, clientType) {
				otherClientLogin = true
			}
			if easy.InArray(typedef_enum.ClientTypeAndroid, clientType) {
				otherClientLogin = true
			}
		case typedef_enum.ClientTypeIos:
			if easy.InArray(typedef_enum.ClientTypeAndroid, clientType) {
				otherClientLogin = true
			}
			if easy.InArray(typedef_enum.ClientTypeIos, clientType) {
				otherClientLogin = true
			}
		}
		if otherClientLogin {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserLoginOtherClient,
			})
		}
		return
	}
	//如果是验证码登录判断是否多账号
	if len(req.Mobile) != 0 && len(req.RegionCode) != 0 {
		count := userDao.Count(&model.User{
			Mobile:     req.Mobile,
			RegionCode: req.RegionCode,
		})
		if count > 1 {
			//多账号
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserLoginCheckMoreAccount,
			})
		}
		user, _ := userDao.FindOne(&model.User{Mobile: req.Mobile, RegionCode: req.RegionCode})
		if len(user.Id) == 0 {
			return
		}
		clientType := service_user.GetUserLoginClientType(user.Id)
		switch headerData.Platform {
		case typedef_enum.ClientTypeAndroid:
			if easy.InArray(typedef_enum.ClientTypeIos, clientType) {
				otherClientLogin = true
			}
			if easy.InArray(typedef_enum.ClientTypeAndroid, clientType) {
				otherClientLogin = true
			}
		case typedef_enum.ClientTypeIos:
			if easy.InArray(typedef_enum.ClientTypeAndroid, clientType) {
				otherClientLogin = true
			}
			if easy.InArray(typedef_enum.ClientTypeIos, clientType) {
				otherClientLogin = true
			}
		}
		if otherClientLogin {
			panic(error2.I18nError{
				Code: error2.ErrorCodeUserLoginOtherClient,
			})
		}
	}
	return
}

// LoginByCode
//
// @Description 根据验证码登录
func (l *Login) LoginByCode(req *request_login.LoginCodeReq, context *gin.Context) (res response_login.LoginCodeRes) {
	//检测验证码
	sms := &Sms{
		Mobile:     req.Mobile,
		Code:       req.Code,
		RegionCode: req.RegionCode,
		Type:       typedef_enum.SmsCodeLogin,
	}
	data := helper.GetHeaderData(context)
	err := sms.CheckSms(context)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeCaptchaInvalid,
			Msg:  nil,
		})
	}
	//获取用户并且注册
	user, jumpType := l.loginOrRegisterByMobile(req.Mobile, req.RegionCode, data, context)
	return l.buildLoginRes(context, user, jumpType)
}

func (l *Login) loginOrRegisterByMobile(mobile, regionCode string, data typedef.HeaderData, context *gin.Context) (user *model.User, jumpType int) {
	success, unlock, _ := coreRedis.UserLock(context, common_data.LockUserRegister(mobile, regionCode), 20*time.Second)
	//锁不成功，提示异常
	if !success {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	defer unlock()
	//根据手机号查询用户
	userDao := new(dao.UserDao)
	user, _ = userDao.FindOne(&model.User{
		Mobile:     mobile,
		RegionCode: regionCode,
	})

	jumpType = typedef_enum.LoginJumpTypeHome
	//代表该手机号未注册
	if len(user.Id) == 0 {
		user = new(model.User)
		jumpType = typedef_enum.LoginJumpTypeRegister
		//注册用户
		user.Mobile = mobile
		user.RegionCode = regionCode
		loginIp := context.ClientIP()
		l.registerUser(user, data, loginIp)
	} else {
		//用户被封禁
		if user.Status == typedef_enum.UserStatusFreezing {
			panic(error2.I18nError{
				Code: error2.ErrorAccountIsBanned,
				Msg:  nil,
			})
		}
		if user.Guide == 1 {
			jumpType = typedef_enum.LoginJumpTypeRegister
		} else {
			count := userDao.Count(&model.User{
				Mobile:     mobile,
				RegionCode: regionCode,
			})
			if count > 1 {
				jumpType = typedef_enum.LoginJumpTypeChooseUser
			}
		}
	}
	return
}

func (l *Login) buildLoginRes(c *gin.Context, user *model.User, jumpType int) (res response_login.LoginCodeRes) {
	userId := user.Id
	if jumpType == typedef_enum.LoginJumpTypeChooseUser {
		userId = ""
	}
	headerData := helper.GetHeaderData(c)
	//生成token
	token, err := coreJwtToken.Encode(coreJwtToken.Claims{
		UserId:         userId,
		Mobile:         user.Mobile,
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
	res.JumpType = jumpType
	if jumpType != typedef_enum.LoginJumpTypeChooseUser {
		res.Token = token
		//将token放入登录缓存
		//需要根据客户端类型进行判定
		switch headerData.Platform {
		case typedef_enum.ClientTypeAndroid, typedef_enum.ClientTypeIos: //app
			coreRedis.GetUserRedis().Set(c, common_data.UserLoginInfo("app", userId), token, typedef_enum.RedisTokenExpireLifeTime)
		case typedef_enum.ClientTypePc: //pc
			coreRedis.GetUserRedis().Set(c, common_data.UserLoginInfo("pc", userId), token, typedef_enum.RedisTokenExpireLifeTime)
		case typedef_enum.ClientTypeH5: //h5
			coreRedis.GetUserRedis().Set(c, common_data.UserLoginInfo("h5", userId), token, typedef_enum.RedisTokenExpireLifeTime)
		}
		res.Id = userId
		res.Uid32 = cast.ToInt32(user.OriUserNo)
	} else {
		token, err = coreJwtToken.Encode(coreJwtToken.Claims{
			UserId:         userId,
			Mobile:         user.Mobile,
			ClientType:     headerData.Platform, //Android Ios Pc Web H5
			DeviceID:       headerData.MachineCode,
			StandardClaims: gojwt.StandardClaims{},
		}, []byte(coreConfig.GetHotConf().JwtSecret), time.Now().Add(typedef_enum.TokenExpireChooseUser).Unix())
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeSystemBusy,
				Msg:  nil,
			})
		}
		res.ChooseUserToken = token
	}
	new(dao.UserLoginRecordDao).LoginRecord(c, user, headerData)
	// 如果已申请注销账号，自动取消注销账号
	service_user.UnCancelAccount(userId)
	if jumpType != typedef_enum.LoginJumpTypeChooseUser {
		//如果在其他端登录则其他端退出通知
		clientType := service_user.GetUserLoginClientType(user.Id)
		otherClientLogin := false
		loginOnClient := headerData.Platform
		switch headerData.Platform {
		case typedef_enum.ClientTypeAndroid:
			if easy.InArray(typedef_enum.ClientTypeIos, clientType) {
				otherClientLogin = true
				loginOnClient = typedef_enum.ClientTypeIos
			}
			if easy.InArray(typedef_enum.ClientTypeAndroid, clientType) {
				otherClientLogin = true
				loginOnClient = typedef_enum.ClientTypeAndroid
			}
		case typedef_enum.ClientTypeIos:
			if easy.InArray(typedef_enum.ClientTypeAndroid, clientType) {
				otherClientLogin = true
				loginOnClient = typedef_enum.ClientTypeAndroid
			}
			if easy.InArray(typedef_enum.ClientTypeIos, clientType) {
				otherClientLogin = true
				loginOnClient = typedef_enum.ClientTypeIos
			}
		}
		if otherClientLogin {
			toUser := []service_im.ToClientUserInfo{
				{
					UserId: user.Id,
					Client: loginOnClient,
				},
			}
			new(service_im.ImCommonService).SendClientUser(c, typedef_enum.SystematicUserId, toUser, typedef_enum.MsgCustom, nil, typedef_enum.USER_LOGIN_OTHER_CLIENT_MSG)
		}
	}
	return
}

func (l *Login) registerUser(user *model.User, data typedef.HeaderData, loginIp string) {
	if data.Platform != "android" && data.Platform != "ios" {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	userNo := helper.GeneUserNo()
	user.Id = coreSnowflake.GetSnowId()
	user.UserNo = userNo
	user.OriUserNo = userNo
	user.Status = typedef_enum.UserStatusNormal
	user.Nickname = "用户" + helper.GeneUserNickname()
	user.Sex = typedef_enum.UserSexTypeWoman
	user.Avatar = helper.GetUserDefaultAvatar(user.Sex)
	user.RegisterChannel = data.Channel
	user.RegisterPlatform = data.Platform
	user.RegisterMachineCode = data.MachineCode
	user.CreateTime = time.Now()
	user.Guide = 1
	user.RealNameStatus = typedef_enum.UserRealNameUnverified
	tx := coreDb.GetMasterDb().Begin()
	//注册用户
	err := tx.Model(user).Create(user).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 注册用户账户信息
	err = tx.Exec("insert into t_user_account(user_id) values(?)", user.Id).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 注册用户lv等级信息
	err = tx.Exec("insert into t_user_lv_level(user_id, create_time, update_time) values(?,?,?)", user.Id, user.CreateTime, user.CreateTime).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 注册用户vip等级信息
	err = tx.Exec("insert into t_user_vip_level(user_id, create_time, update_time) values(?,?,?)", user.Id, user.CreateTime, user.CreateTime).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// TODO 注册用户星光等级信息

	//记录注册信息
	registerRecordModel := &model.UserRegisterRecord{
		UserId:           user.Id,
		RegisterPlatform: data.Platform,
		SignType:         2,
		RegisterChannel:  data.Channel,
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
	tx.Commit()
}

func (l *Login) LoginByPass(req *request_login.LoginPassReq, context *gin.Context) (res response_login.LoginCodeRes) {
	userDao := new(dao.UserDao)
	user, err := userDao.FindOne(&model.User{
		Mobile:     req.Mobile,
		RegionCode: req.RegionCode,
	})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	if len(user.Password) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserPasswordNotSet,
			Msg:  nil,
		})
	}
	user, err = userDao.FindOne(&model.User{
		Mobile:     req.Mobile,
		RegionCode: req.RegionCode,
		Password:   helper.EncodeUserPass(strings.ToLower(req.Password)),
	})
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserPassError,
			Msg:  nil,
		})
	}
	if len(user.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserPassError,
			Msg:  nil,
		})
	}
	jumpType := typedef_enum.LoginJumpTypeHome
	return l.buildLoginRes(context, user, jumpType)
}

func (l *Login) GetChooseUser(chooseUserToken string, context *gin.Context) (result []response_login.GetChooseUserRes) {
	claims, err := coreJwtToken.Decode(chooseUserToken, []byte(coreConfig.GetHotConf().JwtSecret))
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	mobile := claims.Mobile
	userDao := new(dao.UserDao)
	userList, _ := userDao.FindList(&model.User{
		Mobile: mobile,
	})
	for _, item := range userList {
		result = append(result, response_login.GetChooseUserRes{
			Avatar:   helper.FormatImgUrl(item.Avatar),
			Nickname: item.Nickname,
			Id:       item.Id,
		})
	}
	return
}

func (l *Login) ChooseUserLogin(req *request_login.ChooseUserLoginReq, c *gin.Context) (result response_login.LoginCodeRes) {

	claims, err := coreJwtToken.Decode(req.ChooseUserToken, []byte(coreConfig.GetHotConf().JwtSecret))
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	mobile := claims.Mobile
	userDao := new(dao.UserDao)
	user, _ := userDao.FindOne(&model.User{
		Id: req.Id,
	})
	if mobile != user.Mobile {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}

	return l.buildLoginRes(c, user, typedef_enum.LoginJumpTypeHome)
}

func (l *Login) ForgetPassCheck(req *request_login.LoginCodeReq, c *gin.Context) response_login.ForgetPassCheckRes {
	//检测验证码
	sms := &Sms{
		Mobile:     req.Mobile,
		Code:       req.Code,
		RegionCode: req.RegionCode,
		Type:       typedef_enum.SmsCodeForgetPass,
	}
	err := sms.CheckSms(c)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeCaptchaInvalid,
			Msg:  nil,
		})
	}
	//查询当前手机号有几个用户
	userDao := new(dao.UserDao)
	count := userDao.Count(&model.User{
		Mobile: req.Mobile,
	})
	if count == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	jumpType := 1
	id := ""
	token := ""
	if count > 1 {
		jumpType = 2
	} else {
		user, _ := userDao.FindOne(&model.User{
			Mobile: req.Mobile,
		})
		id = user.Id
	}
	token, _ = coreJwtToken.Encode(coreJwtToken.Claims{
		Mobile:         req.Mobile,
		StandardClaims: gojwt.StandardClaims{},
	}, []byte(coreConfig.GetHotConf().JwtSecret), time.Now().Add(typedef_enum.RedisTokenExpireLifeTime).Unix())
	return response_login.ForgetPassCheckRes{
		Id:              id,
		JumpType:        jumpType,
		ChooseUserToken: token,
	}
}

func (l *Login) ForgetPass(req *request_login.ForgetPassReq, c *gin.Context) {
	userDao := new(dao.UserDao)
	user, _ := userDao.FindOne(&model.User{
		Id: req.Id,
	})
	if len(user.Id) == 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeUserNotFound,
			Msg:  nil,
		})
	}
	claims, err := coreJwtToken.Decode(req.ChooseUserToken, []byte(coreConfig.GetHotConf().JwtSecret))
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	//判断是否是同一个手机号
	if user.Mobile != claims.Mobile {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
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
	userDao.UpdateById(&model.User{
		Id:       user.Id,
		Password: newPass,
	})
}
