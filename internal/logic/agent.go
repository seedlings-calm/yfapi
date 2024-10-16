package logic

import (
	"context"
	gojwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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
	"yfapi/internal/service/agent"
	typedef_enum "yfapi/typedef/enum"
	common_data "yfapi/typedef/redisKey"
)

type Agent struct {
	ClientType string //客户端类型
}

func (a *Agent) OauthToken(c *gin.Context, data map[string]any) string {
	if data["appId"] == nil || data["uid"] == nil || data["nickname"] == nil || data["avatar"] == nil || data["gender"] == nil || data["timestamp"] == nil || data["signature"] == nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
		})
	}
	appId := cast.ToString(data["appId"])
	uid := cast.ToString(data["uid"])
	nickname := cast.ToString(data["nickname"])
	avatar := cast.ToString(data["avatar"])
	sex := cast.ToInt(data["gender"])
	oriSignatue := cast.ToString(data["signature"])
	secret := new(dao.AgentDao).GetSecretByAppId(cast.ToString(data["appId"]))
	service := &agent.Oauth{
		AppId:  cast.ToString(data["appId"]),
		Secret: secret,
	}
	sign := service.Sign(data)
	if oriSignatue != sign {
		panic(error2.I18nError{Code: error2.ErrCodeAgentSignErr})
	}
	user := a.loginOrRegister(appId, uid, nickname, avatar, sex)
	token := a.buildLoginRes(c, user)
	return token
}

func (a *Agent) loginOrRegister(appId, uid, nickname, avatar string, sex int) (user *model.User) {
	success, unlock, _ := coreRedis.UserLock(context.Background(), common_data.LockUserRegister(appId, uid), 20*time.Second)
	//锁不成功，提示异常
	if !success {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	defer unlock()
	//根据uid appId
	userDao := new(dao.UserDao)
	user, _ = userDao.FindOne(&model.User{
		AppId: appId,
		Uid:   uid,
	})
	//未注册
	if len(user.Id) == 0 {
		user = new(model.User)
		user.AppId = appId
		user.Uid = uid
		user.Nickname = nickname
		user.Avatar = avatar
		user.Sex = sex
		user.Source = 1
		a.registerUser(user)
	}
	return
}

func (a *Agent) registerUser(user *model.User) {
	userNo := helper.GeneUserNo()
	user.Id = coreSnowflake.GetSnowId()
	user.UserNo = userNo
	user.OriUserNo = userNo
	user.Status = typedef_enum.UserStatusNormal
	user.RegisterChannel = user.AppId + "_" + user.Uid
	user.RegisterPlatform = a.ClientType
	user.CreateTime = time.Now()
	user.Guide = 2
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
	err = tx.Exec("insert into t_user_account(user_id) values(?)", user.Id).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	tx.Commit()
}

func (a *Agent) buildLoginRes(c *gin.Context, user *model.User) string {
	//生成token
	token, err := coreJwtToken.Encode(coreJwtToken.Claims{
		UserId:         user.Id,
		Mobile:         user.Mobile,
		ClientType:     a.ClientType,
		StandardClaims: gojwt.StandardClaims{},
	}, []byte(coreConfig.GetHotConf().JwtSecret), time.Now().Add(typedef_enum.TokenExpireLifeTime).Unix())
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	coreRedis.GetUserRedis().Set(context.Background(), common_data.UserLoginInfo("h5", user.Id), token, typedef_enum.RedisTokenExpireLifeTime)
	new(dao.UserLoginRecordDao).LoginRecord(c, user, helper.GetHeaderData(c))
	return token
}
