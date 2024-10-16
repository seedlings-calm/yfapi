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

// 私聊语音
// 同步检测
func (s *ShuMei) PrivateChatAudioCheck(userId, voiceUrl, btId string) (string, bool) {
	if !coreConfig.GetHotConf().RiskSwitch {
		return "{}", true
	}
	req := map[string]any{
		"userId":   userId,
		"voiceUrl": voiceUrl,
	}
	config := coreConfig.GetHotConf().ShuMei
	res := &AudioCheckSyncResp{}
	err := s.audioCheck(&PrivateChatAudioCheckConfig{
		UserId:    userId,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
		BtId:      btId,
		Content:   voiceUrl,
	}, AudioCheckUrl, res)
	marshal, _ := json.Marshal(res)
	resStr := string(marshal)
	coreLog.LogInfo("数美私聊音频检测结果:%+v", res)
	if err != nil {
		coreLog.Error("PrivateChatAudioCheck err:%+v req:%+v", err, req)
		return resStr, false
	}
	if res.Code != 1100 {
		coreLog.LogError("PrivateChatAudioCheck err:%s", res.Message)
		return resStr, false
	} else {
		if res.Detail.RiskLevel == RiskLevelReject { //发生违规
			return resStr, false
		}
	}
	return resStr, true
}

// 签名语音检测 异步检测
func (s *ShuMei) SignAudioAsyncCheck(userId, voiceUrl, uniqueId, btId string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId":   userId,
		"voiceUrl": voiceUrl,
	}
	config := coreConfig.GetHotConf().ShuMei
	encrypt := easy.AesEncrypt(cast.ToString(uniqueId), coreConfig.GetHotConf().InnerSecret.AESEncryptKey1)
	toString := base64.StdEncoding.EncodeToString([]byte(encrypt))
	res := &AudioCheckResp{}
	err := s.audioCheck(&SignAudioAsyncCheckConfig{
		UserId:    userId,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
		Callback:  config.AudioSignCallback + toString,
		BtId:      btId,
		Content:   voiceUrl,
	}, GetAudioUrl(), res)
	coreLog.LogInfo("数美声音签名检测结果:%+v", res)
	if err != nil {
		coreLog.Error("SignAudioAsyncCheck err:%+v req:%+v", err, req)
		return false
	}
	if res.Code != 1100 {
		coreLog.LogError("SignAudioAsyncCheck err:%s", res.Message)
		return false
	}
	return true
}

// 音频检测
func (s *ShuMei) audioCheck(config CheckConfig, url string, data any) error {
	payload := config.getPayLoadData()
	coreLog.Info("数美音频 req:%+v", payload)
	b, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	respBytes, _ := io.ReadAll(resp.Body)
	coreLog.Info("数美音频 resp:%s", string(respBytes))
	err = json.Unmarshal(respBytes, &data)
	if err != nil {
		return err
	}
	return nil
}
