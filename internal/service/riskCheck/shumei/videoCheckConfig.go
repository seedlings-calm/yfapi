package shumei

// 朋友圈视频检测
type MomentsVideoAsyncCheckConfig struct {
	EventId   string
	UserId    string
	AccessKey string
	AppId     string
	Callback  string
	BtId      string
	VideoUrl  string
}

func (n *MomentsVideoAsyncCheckConfig) getPayLoadData() any {
	payload := VideoCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		AudioType: "POLITICAL_PORN_AD_ABUSE_MOAN",
		Callback:  n.Callback,
		Data: VideoCheckReqData{
			BtId:    n.BtId,
			TokenId: n.UserId,
			Url:     n.VideoUrl,
		},
		EventId: "dynamic",
		ImgType: "POLITY_EROTIC_VIOLENT_ADVERT_QRCODE_IMGTEXTRISK",
	}
	return payload
}

// 私聊视频检测
type PrivateChatVideoAsyncCheckConfig struct {
	EventId   string
	UserId    string
	AccessKey string
	AppId     string
	Callback  string
	BtId      string
	VideoUrl  string
}

func (n *PrivateChatVideoAsyncCheckConfig) getPayLoadData() any {
	payload := VideoCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		AudioType: "POLITICAL_PORN_AD_ABUSE_MOAN",
		Callback:  n.Callback,
		Data: VideoCheckReqData{
			BtId:    n.BtId,
			TokenId: n.UserId,
			Url:     n.VideoUrl,
		},
		EventId: "message",
		ImgType: "POLITY_EROTIC_VIOLENT_ADVERT_QRCODE_IMGTEXTRISK",
	}
	return payload
}
