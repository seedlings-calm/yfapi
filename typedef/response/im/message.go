package response_im

import response_goods "yfapi/typedef/response/goods"

type ResponseMsg struct {
	Sn    int64  `json:"sn"`
	Msg   any    `json:"msg"`
	Extra string `json:"extra"`
}

// JoinRoomImResponse 加入房间推送信息
type JoinRoomImResponse struct {
	Content     string                        `json:"content"`          //描述信息
	JoinSE      response_goods.SpecialEffects `json:"joinSE,omitempty"` // 进场特效信息
	Color       []string                      `json:"color"`            //进场横幅颜色
	StrokeColor []string                      `json:"strokeColor"`      // 进场横幅描边颜色
}
