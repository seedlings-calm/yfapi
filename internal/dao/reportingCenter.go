package dao

import (
	"encoding/json"
	"time"
	"yfapi/core/coreDb"
	"yfapi/internal/model"
)

type ReportingCenterDao struct{}

func (r ReportingCenterDao) AddReportingCenter(userId, dstId, content string, picUrl []string, object, scene, types int) *model.ReportingCenter {
	resp := &model.ReportingCenter{
		SrcUserID:     userId,
		DstUserID:     dstId,
		Object:        object,
		Scene:         scene,
		State:         1,
		ReportType:    types,
		ReportContent: content,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	if len(picUrl) == 0 {
		resp.ReportPicURL = "[]"
	} else {
		picsS, _ := json.Marshal(&picUrl)
		resp.ReportPicURL = string(picsS)
	}
	return resp
}

func (r *ReportingCenterDao) Insert(param *model.ReportingCenter) error {
	err := coreDb.GetMasterDb().Model(model.ReportingCenter{}).Save(param).Error
	return err
}
