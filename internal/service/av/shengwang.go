package av

import (
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/rtctokenbuilder2"
	"github.com/spf13/cast"
	"yfapi/core/coreConfig"
	i18n_err "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
)

type ShengWangAv struct {
}

func (av *ShengWangAv) GetToken(userId, channelName string) (string, error) {
	user, err := new(dao.UserDao).FindOne(&model.User{Id: userId})
	if err != nil {
		panic(i18n_err.I18nError{
			Code: i18n_err.ErrorCodeUserNotFound,
		})
	}
	token := ""
	// Token 的有效时间，单位秒
	tokenExpirationInSeconds := uint32(3600 * 24 * 30)
	// 所有的权限的有效时间，单位秒
	privilegeExpirationInSeconds := uint32(3600 * 24 * 30)
	conf := coreConfig.GetHotConf().Av
	//token, err := rtctokenbuilder.BuildTokenWithUserAccount(conf.AppId, conf.Certificate, channelName, userId, rtctokenbuilder.RoleSubscriber, tokenExpirationInSeconds, privilegeExpirationInSeconds)
	// 生成 Token
	token, err = rtctokenbuilder.BuildTokenWithUid(conf.AppId, conf.Certificate, channelName, cast.ToUint32(user.OriUserNo), rtctokenbuilder.RoleSubscriber, tokenExpirationInSeconds, privilegeExpirationInSeconds)
	return token, err
}
