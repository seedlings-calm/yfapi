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

// 私聊视频
// 异步检测
func (s *ShuMei) PrivateChatVideoAsyncCheck(userId, btId, videoUrl string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId":   userId,
		"btId":     btId,
		"videoUrl": videoUrl,
	}
	config := coreConfig.GetHotConf().ShuMei
	res := &VideoCheckResp{}
	err := s.videoCheck(&PrivateChatVideoAsyncCheckConfig{
		UserId:    userId,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
		Callback:  "",
		BtId:      btId,
		VideoUrl:  videoUrl,
	}, GetVideoUrl(), res)
	if err != nil {
		coreLog.Error("PrivateChatVideoAsyncCheck err:%+v req:%+v", err, req)
		return false
	}
	if res.Code != 1100 {
		coreLog.LogError("PrivateChatVideoAsyncCheck err:%s", res.Message)
		return false
	}
	return true
}

// 朋友圈视频检测
// 异步检测
func (s *ShuMei) MomentsVideoAsyncCheck(userId, videoUrl, uniqueId, btId string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId":   userId,
		"btId":     btId,
		"videoUrl": videoUrl,
	}
	config := coreConfig.GetHotConf().ShuMei
	encrypt := easy.AesEncrypt(cast.ToString(uniqueId), coreConfig.GetHotConf().InnerSecret.AESEncryptKey1)
	toString := base64.StdEncoding.EncodeToString([]byte(encrypt))
	res := &VideoCheckResp{}
	err := s.videoCheck(&MomentsVideoAsyncCheckConfig{
		UserId:    userId,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
		Callback:  config.MomentsVideoCallback + toString,
		BtId:      btId,
		VideoUrl:  videoUrl,
	}, GetVideoUrl(), res)
	if err != nil {
		coreLog.Error("MomentsVideoAsyncCheck err:%+v req:%+v", err, req)
		return false
	}
	if res.Code != 1100 {
		coreLog.LogError("MomentsVideoAsyncCheck err:%s", res.Message)
		return false
	}
	return true
}

// 视频检测
func (s *ShuMei) videoCheck(config CheckConfig, url string, data any) error {
	payload := config.getPayLoadData()
	b, _ := json.Marshal(payload)
	resp, err := http.Post(GetVideoUrl(), "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	respBytes, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(respBytes, &data)
	if err != nil {
		return err
	}
	return err
}
