package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type AgentDao struct {
}

func (a *AgentDao) GetSecretByAppId(appId string) (secret string) {
	coreDb.GetMasterDb().Model(&model.Agent{}).Where("app_id = ? and status = 1", appId).Select("secret").Find(&secret)
	return
}
