package shumei

// 昵称检测
type NicknameTextCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	Lang      string
	Text      string
	AccessKey string
	AppId     string
}

func (n *NicknameTextCheckConfig) getPayLoadData() any {
	payload := TextCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		EventId:   n.EventId,
		Type:      "TEXTRISK",
		Data: TextCheckReqData{
			Text:    n.Text,
			TokenId: n.UserId,
		},
	}
	return payload
}

// 私聊文本见检测
type PrivateChatTextCheckConfig struct {
	EventId       string
	Type          string
	UserId        string
	ReceiveUserId string
	Lang          string
	Text          string
	AccessKey     string
	AppId         string
	Topic         string
}

func (n *PrivateChatTextCheckConfig) getPayLoadData() any {
	payload := TextCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		EventId:   n.EventId,
		Type:      "TEXTRISK",
		Data: TextCheckReqData{
			Text:    n.Text,
			TokenId: n.UserId,
			Extra: TextCheckReqExtra{
				ReceiveTokenId: n.ReceiveUserId,
				Topic:          n.Topic,
			},
		},
	}
	return payload
}

// 公屏聊天文本配置
type PublicChatTextCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	Lang      string
	Text      string
	AccessKey string
	AppId     string
}

func (n *PublicChatTextCheckConfig) getPayLoadData() any {
	payload := TextCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		EventId:   n.EventId,
		Type:      "TEXTRISK",
		Data: TextCheckReqData{
			Text:    n.Text,
			TokenId: n.UserId,
		},
	}
	return payload
}

// 朋友圈文本检测
type MomentsTextCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	Lang      string
	Text      string
	AccessKey string
	AppId     string
}

func (n *MomentsTextCheckConfig) getPayLoadData() any {
	payload := TextCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		EventId:   n.EventId,
		Type:      "TEXTRISK",
		Data: TextCheckReqData{
			Text:    n.Text,
			TokenId: n.UserId,
		},
	}
	return payload
}

// 评论文本检测
type CommentTextCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	Lang      string
	Text      string
	AccessKey string
	AppId     string
}

func (n *CommentTextCheckConfig) getPayLoadData() any {
	payload := TextCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		EventId:   n.EventId,
		Type:      "TEXTRISK",
		Data: TextCheckReqData{
			Text:    n.Text,
			TokenId: n.UserId,
		},
	}
	return payload
}

// 签名
type SignTextCheckConfig struct {
	EventId   string
	Type      string
	UserId    string
	Lang      string
	Text      string
	AccessKey string
	AppId     string
}

func (n *SignTextCheckConfig) getPayLoadData() any {
	payload := TextCheckReq{
		AccessKey: n.AccessKey,
		AppId:     n.AppId,
		EventId:   n.EventId,
		Type:      "TEXTRISK",
		Data: TextCheckReqData{
			Text:    n.Text,
			TokenId: n.UserId,
		},
	}
	return payload
}
