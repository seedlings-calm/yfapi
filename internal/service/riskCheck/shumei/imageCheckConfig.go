package shumei

// 头像图片检测
type AvatarImageAsyncCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	AccessKey string
	AppId     string
	Callback  string
	Avatar    string
	BtId      string
}

func (n *AvatarImageAsyncCheckConfig) getPayLoadData() any {
	imgs := []ImagesAsyncCheckReqImgs{
		{
			BtId: n.BtId,
			Img:  n.Avatar,
		},
	}
	payload := ImagesAsyncCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		Data: ImagesAsyncCheckReqData{
			Imgs:    imgs,
			TokenId: n.UserId,
		},
		EventId:  n.EventId,
		Type:     "POLITY_EROTIC_VIOLENT_ADVERT_QRCODE_IMGTEXTRISK",
		Callback: n.Callback,
	}
	return payload
}

// 头像同步检测
type AvatarImageSyncCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	AccessKey string
	AppId     string
	Avatar    string
}

func (n *AvatarImageSyncCheckConfig) getPayLoadData() any {
	payload := OneImageSyncCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		Data: OneImageSyncCheckReqData{
			Img:     n.Avatar,
			TokenId: n.UserId,
		},
		EventId: n.EventId,
		Type:    "POLITY_EROTIC_VIOLENT_ADVERT_QRCODE_IMGTEXTRISK",
	}
	return payload
}

// 私聊图片同步检测
type OnePrivateChatImageSyncCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	AccessKey string
	AppId     string
	Image     string
}

func (n *OnePrivateChatImageSyncCheckConfig) getPayLoadData() any {
	payload := OneImageSyncCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		Data: OneImageSyncCheckReqData{
			Img:     n.Image,
			TokenId: n.UserId,
		},
		EventId: n.EventId,
		Type:    "POLITY_EROTIC_VIOLENT_ADVERT_QRCODE_IMGTEXTRISK",
	}
	return payload
}

// 圈子图片检测
type MomentsImageAsyncCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	AccessKey string
	AppId     string
	Callback  string
	Images    []ImagesAsyncCheckReqImgs
}

func (n *MomentsImageAsyncCheckConfig) getPayLoadData() any {
	payload := ImagesAsyncCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		Data: ImagesAsyncCheckReqData{
			Imgs:    n.Images,
			TokenId: n.UserId,
		},
		EventId:  n.EventId,
		Type:     "POLITY_EROTIC_VIOLENT_ADVERT_QRCODE_IMGTEXTRISK",
		Callback: n.Callback,
	}
	return payload
}
