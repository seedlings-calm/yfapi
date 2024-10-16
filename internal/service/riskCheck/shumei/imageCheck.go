package shumei

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/spf13/cast"
	"io"
	"net/http"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
	"yfapi/util/easy"
)

// 头像检测
// 同步
func (s *ShuMei) AvatarSyncCheck(userId, avatar string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId": userId,
		"avatar": avatar,
	}
	config := coreConfig.GetHotConf().ShuMei
	res := &OneImageSyncCheckResp{}
	err := s.imageCheck(&AvatarImageSyncCheckConfig{
		EventId:   "headImage",
		UserId:    userId,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
		Avatar:    avatar,
	}, GetOneImageUrl(), res)
	coreLog.LogInfo("头像同步数美检测结果:%+v", res)
	if err != nil {
		coreLog.Error("AvatarCheck err:%+v req:%+v", err, req)
		return false
	}
	if res.Code != 1100 {
		coreLog.LogError("AvatarCheck err:%s", res.Message)
		return false
	} else {
		if res.RiskLevel == RiskLevelReject {
			return false
		}
	}
	return true
}

// 单聊图片单张检测
func (s *ShuMei) OneChatImageSyncCheck(userId, image string) (string, bool) {
	if !coreConfig.GetHotConf().RiskSwitch {
		return "{}", true
	}
	req := map[string]any{
		"userId": userId,
		"image":  image,
	}
	config := coreConfig.GetHotConf().ShuMei
	res := &OneImageSyncCheckResp{}
	err := s.imageCheck(&OnePrivateChatImageSyncCheckConfig{
		EventId:   "message",
		UserId:    userId,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
		Image:     image,
	}, GetOneImageUrl(), res)
	coreLog.LogInfo("私聊图片同步数美检测结果:%+v", res)
	coreLog.LogInfo("OnePrivateChatImageSyncCheck res:%+v", res)
	marshal, _ := json.Marshal(res)
	riskRes := string(marshal)
	if err != nil {
		coreLog.Error("OnePrivateChatImageSyncCheck err:%+v req:%+v", err, req)
		return riskRes, false
	}
	if res.Code != 1100 {
		coreLog.LogError("OnePrivateChatImageSyncCheck err:%s", res.Message)
		return riskRes, false
	} else {
		if res.RiskLevel == RiskLevelReject {
			return riskRes, false
		}
	}
	return riskRes, true
}

// 朋友圈图片检测
// 异步
func (s *ShuMei) MomentsImageAsyncCheck(userId string, images []ImagesAsyncCheckReqImgs, uniqueId string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId": userId,
		"images": images,
	}
	config := coreConfig.GetHotConf().ShuMei
	encrypt := easy.AesEncrypt(cast.ToString(uniqueId), coreConfig.GetHotConf().InnerSecret.AESEncryptKey1)
	toString := base64.StdEncoding.EncodeToString([]byte(encrypt))
	res := &ImagesAsyncCheckResp{}
	err := s.imageCheck(&MomentsImageAsyncCheckConfig{
		EventId:   "dynamic",
		UserId:    userId,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
		Callback:  config.MomentsImageCallback + toString,
		Images:    images,
	}, GetImageUrl(), res)
	coreLog.LogInfo("朋友圈图片异步数美检测结果:%+v", res)
	if err != nil {
		coreLog.Error("MomentsImageCheck err:%+v req:%+v", err, req)
		return false
	}
	if res.Code != 1100 {
		coreLog.LogError("MomentsImageCheck err:%s", res.Message)
		return false
	}
	return true
}

func (s *ShuMei) imageCheck(config CheckConfig, url string, data any) error {
	payload := config.getPayLoadData()
	b, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	respBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &data)
	if err != nil {
		return err
	}
	return nil
}
