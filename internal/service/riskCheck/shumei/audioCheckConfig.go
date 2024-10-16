package shumei

// 签名音频检测
type SignAudioAsyncCheckConfig struct {
	EventId     string
	Type        string
	UserId      string
	AccessKey   string
	AppId       string
	Callback    string
	BtId        string
	ContentType string
	Content     string
}

func (n *SignAudioAsyncCheckConfig) getPayLoadData() any {
	payload := AudioCheckReq{
		AccessKey:   n.AccessKey,
		AppId:       n.AppId,
		EventId:     "sign",
		Type:        "POLITICAL_PORN_AD_ABUSE_MOAN",
		BtId:        n.BtId,
		ContentType: "URL",
		Content:     n.Content,
		Callback:    n.Callback,
		Data: AudioCheckReqData{
			TokenId: n.UserId,
		},
	}
	return payload
}

// 私聊音频检测
type PrivateChatAudioCheckConfig struct {
	EventId     string
	Type        string
	UserId      string
	AccessKey   string
	AppId       string
	BtId        string
	ContentType string
	Content     string
}

func (n *PrivateChatAudioCheckConfig) getPayLoadData() any {
	payload := AudioCheckReq{
		AccessKey:   n.AccessKey,
		AppId:       n.AppId,
		EventId:     "message",
		Type:        "POLITICAL_PORN_AD_ABUSE_MOAN",
		BtId:        n.BtId,
		ContentType: "URL",
		Content:     n.Content,
		Data: AudioCheckReqData{
			TokenId: n.UserId,
		},
	}
	return payload
}
