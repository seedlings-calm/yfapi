package shumei

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"yfapi/core/coreConfig"
	"yfapi/core/coreLog"
)

// 昵称检测 true通过 false不通过
// 同步检测
func (s *ShuMei) NicknameCheck(userId, nickname string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId":   userId,
		"nickname": nickname,
	}
	config := coreConfig.GetHotConf().ShuMei
	check, err := s.textCheck(&NicknameTextCheckConfig{
		EventId:   "nickname",
		UserId:    userId,
		Text:      nickname,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
	})
	coreLog.LogInfo("昵称数美检测结果:%+v", check)
	if err != nil {
		coreLog.Error("NicknameCheck err:%+v req:%+v", err, req)
		return false
	}
	if res, ok := check.(TextCheckResp); ok {
		if res.Code != 1100 {
			coreLog.LogError("NicknameCheck err:%s", res.Message)
			return false
		} else {
			if res.RiskLevel == RiskLevelReject { //发生违规
				return false
			}
		}
	}
	return true
}

// 私聊检测 true通过 false不通过
// 同步检测
func (s *ShuMei) PrivateChatCheck(userId, receiveUserId, text string) (string, bool) {
	if !coreConfig.GetHotConf().RiskSwitch {
		return "{}", true
	}
	req := map[string]any{
		"userId": userId,
		"text":   text,
	}
	config := coreConfig.GetHotConf().ShuMei
	check, err := s.textCheck(&PrivateChatTextCheckConfig{
		EventId:       "message",
		UserId:        userId,
		Text:          text,
		AccessKey:     config.AccessKey,
		AppId:         config.AppId,
		ReceiveUserId: receiveUserId,
		Topic:         userId + receiveUserId,
	})
	coreLog.LogInfo("私聊数美检测结果:%+v", check)
	var riskRes string = ""
	if err != nil {
		coreLog.Error("PrivateChatCheck err:%+v req:%+v", err, req)
		return riskRes, false
	}
	if res, ok := check.(TextCheckResp); ok {
		marshal, _ := json.Marshal(res)
		riskRes = string(marshal)
		if res.Code != 1100 {
			coreLog.LogError("PrivateChatCheck err:%s", res.Message)
			return riskRes, false
		} else {
			if res.RiskLevel == RiskLevelReject { //发生违规
				return riskRes, false
			}
		}
		return riskRes, true
	} else {
		coreLog.Error("PrivateChatCheck check.(TextCheckResp) err %+v", check)
	}
	return riskRes, false
}

// 朋友圈文本检测
// 同步检测
func (s *ShuMei) MomentsCheck(userId, text string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId": userId,
		"text":   text,
	}
	config := coreConfig.GetHotConf().ShuMei
	check, err := s.textCheck(&MomentsTextCheckConfig{
		EventId:   "dynamic",
		UserId:    userId,
		Text:      text,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
	})
	coreLog.LogInfo("朋友圈数美检测结果:%+v", check)
	if err != nil {
		coreLog.Error("MomentsCheck err:%+v req:%+v", err, req)
		return false
	}
	if res, ok := check.(TextCheckResp); ok {
		if res.Code != 1100 {
			coreLog.LogError("MomentsCheck err:%s", res.Message)
			return false
		} else {
			if res.RiskLevel == RiskLevelReject { //发生违规
				return false
			}
		}
	}
	return true
}

// 评论检测
// 同步检测
func (s *ShuMei) CommentCheck(userId, text string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId": userId,
		"text":   text,
	}
	config := coreConfig.GetHotConf().ShuMei
	check, err := s.textCheck(&CommentTextCheckConfig{
		EventId:   "comment",
		UserId:    userId,
		Text:      text,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
	})
	coreLog.LogInfo("评论数美检测结果:%+v", check)
	if err != nil {
		coreLog.Error("CommentChatCheck err:%+v req:%+v", err, req)
		return false
	}
	if res, ok := check.(TextCheckResp); ok {
		if res.Code != 1100 {
			coreLog.LogError("CommentChatCheck err:%s", res.Message)
			return false
		} else {
			if res.RiskLevel == RiskLevelReject { //发生违规
				return false
			}
		}
	}
	return true
}

// 签名检测
// 同步检测
func (s *ShuMei) SignCheck(userId, text string) bool {
	if !coreConfig.GetHotConf().RiskSwitch {
		return true
	}
	req := map[string]any{
		"userId": userId,
		"text":   text,
	}
	config := coreConfig.GetHotConf().ShuMei
	check, err := s.textCheck(&SignTextCheckConfig{
		EventId:   "sign",
		UserId:    userId,
		Text:      text,
		AccessKey: config.AccessKey,
		AppId:     config.AppId,
	})
	coreLog.LogInfo("签名数美检测结果:%+v", check)
	if err != nil {
		coreLog.Error("SignChatCheck err:%+v req:%+v", err, req)
		return false
	}
	if res, ok := check.(TextCheckResp); ok {
		if res.Code != 1100 {
			coreLog.LogError("SignChatCheck err:%s", res.Message)
			return false
		} else {
			if res.RiskLevel == RiskLevelReject { //发生违规
				return false
			}
		}
	}
	return true
}

// 文本检测
func (s *ShuMei) textCheck(config CheckConfig) (any, error) {
	payload := config.getPayLoadData()
	coreLog.LogInfo("textCheck payload :%+v", payload)
	var data TextCheckResp
	url := GetTextUrl()
	b, _ := json.Marshal(payload)
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(b))
	if resp != nil {
		respBytes, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(respBytes, &data)
	}
	return data, nil
}
