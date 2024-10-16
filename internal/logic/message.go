package logic

import (
	"encoding/json"
	"yfapi/core/coreLog"
	"yfapi/internal/dao"
	"yfapi/typedef/enum"
	"yfapi/typedef/request/message"
)

type Message struct {
}

// 获取会话历史聊天记录
func (m *Message) GetMessageList(selfUserId string, req message.MessageListReq) any {
	list := new(dao.ChatStoreDao).GetPrivateChatList(selfUserId, req.OtherUserId, req.Limit, req.Timestamp)
	data := []any{}
	if len(list) == 0 {
		return data
	}
	if req.OtherUserId != enum.InteractiveUserId {
		l := len(list) - 1
		for i := 0; i <= l; i++ {
			baseMsg := map[string]any{}
			err := json.Unmarshal([]byte(list[l-i].Message), &baseMsg)
			if err != nil {
				coreLog.Error("GetMessageList Unmarshal err:%+v", err)
				continue
			}
			data = append(data, baseMsg)
		}
	} else {
		for _, v := range list {
			baseMsg := map[string]any{}
			err := json.Unmarshal([]byte(v.Message), &baseMsg)
			if err != nil {
				coreLog.Error("GetMessageList Unmarshal err:%+v", err)
				continue
			}
			data = append(data, baseMsg)
		}
	}
	return data
}
