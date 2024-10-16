package message

import (
	"yfapi/typedef/response"
	response_goods "yfapi/typedef/response/goods"
)

// 基础消息
type BaseMsg struct {
	MessageId    string `json:"messageId"`              //消息ID
	Timestamp    int64  `json:"timestamp"`              //时间戳毫秒
	MsgType      string `json:"msgType"`                //消息类型 语音audio 文本text 图片img 自定义custom
	Code         int    `json:"code"`                   //消息编码  用户消息 10000-19999  房间消息 20000-29999 系统消息 30000-39999
	MsgData      any    `json:"msgData,omitempty"`      //消息内容
	FromUserInfo any    `json:"fromUserInfo,omitempty"` //来源用户信息
	ToUserInfo   any    `json:"toUserInfo,omitempty"`   //目标用户信息
	RiskLevel    string `json:"riskLevel,omitempty"`    //拦截等级
	RoomId       string `json:"roomId"`                 //房间ID
	ExtraInfo    any    `json:"extraInfo,omitempty"`    // 扩展信息
}

type MsgText struct {
	Content string                        `json:"content"`          //文本内容
	Bubble  response_goods.SpecialEffects `json:"bubble,omitempty"` // 气泡
}

type MsgImg struct {
	Content string `json:"content"` //图片地址
	Width   int    `json:"width"`   //宽度
	Height  int    `json:"height"`  //长度
}

type MsgAudio struct {
	Content string `json:"content"` //音频文件地址
	Length  int    `json:"length"`  //音频秒数
}

type MsgAction struct {
	Content string `json:"content"`
	Data    any    `json:"data"`
}

// 互动消息
type InteractiveMsg struct {
	Content string `json:"content"`
	Types   int    `json:"types"`
	ID      int    `json:"id"`
}

// 系统消息
type SystematicMsg struct {
	Title     string `json:"title"`     // 标题
	Img       string `json:"img"`       // 图片
	Content   string `json:"content"`   // 内容
	Link      string `json:"link"`      // 链接
	H5Content string `json:"h5Content"` //富文本内容
}

// 官方公告消息
type OfficialMsg struct {
	Title     string `json:"title"`     // 标题
	Img       string `json:"img"`       // 图片
	Content   string `json:"content"`   // 内容
	Link      string `json:"link"`      // 链接
	H5Content string `json:"h5Content"` //富文本内容
}

// OneMsgListModel 会话列表数据
type OneMsgListModel struct {
	ShowContent string `json:"showContent"` //展示的文本信息
	Timestamp   int64  `json:"timestamp"`   //时间戳毫秒级
	TextColor   string `json:"textColor"`   //文本颜色
	ToUserId    string `json:"toUserId"`    //对应的用户id
	IsTop       bool   `json:"isTop"`       //是否置顶
}

// MsgGift 礼物打赏通知
type MsgGift struct {
	GiftName         string `json:"giftName"`         // 礼物名称
	GiftImage        string `json:"giftImage"`        // 礼物图片
	AnimationUrl     string `json:"animationUrl"`     // VAP配置地址
	AnimationJsonUrl string `json:"animationJsonUrl"` // VAP JSON配置地址
	GiftCount        int    `json:"giftCount"`        // 赠送数量
	ComboCount       int    `json:"comboCount"`       // 连击数
	ComboKey         string `json:"comboKey"`         // 连击key
	TotalGiftDiamond int    `json:"totalGiftDiamond"` // 打赏总钻石数
	IsBatch          bool   `json:"isBatch"`          // 是否全麦打赏
}

// MsgCustom 自定义消息通知
type MsgCustom struct {
	MessageId string `json:"messageId"` //消息ID
	Timestamp int64  `json:"timestamp"` //时间戳毫秒
	MsgType   string `json:"msgType"`   //消息类型 语音audio 文本text 图片img 自定义custom
	Code      int    `json:"code"`      //消息编码  用户消息 10000-19999  房间消息 20000-29999 系统消息 30000-39999
	MsgData   any    `json:"msgData"`   //消息内容
	RoomId    string `json:"roomId"`    //房间ID
	ExtraInfo any    `json:"extraInfo"` // 扩展信息
}

// MsgSendGiftSeat 被打赏麦位动效通知
type MsgSendGiftSeat struct {
	GiftImage    string `json:"giftImage"`    // 礼物图片
	FromSeatId   int    `json:"fromSeatId"`   // 打赏人麦位ID -1不在麦
	ToSeatIdList []int  `json:"toSeatIdList"` // 被打赏的麦位ID列表
}

// MsgJoinRoom 加入房间信息通知
type MsgJoinRoom struct {
	Content string `json:"content"` // 进房公屏信息 王道长加入房间  王道长踩着冯宝宝的小尾巴进入房间
	Url     string `json:"url"`     // 进场特效
}

// MsgLevelUp 等级升级信息
type MsgLevelUp struct {
	UserId   string `json:"userId"`
	Nickname string `json:"nickname"`
	Icon     string `json:"icon"`
}

// SendInteractiveMsg 互动消息
type SendInteractiveMsg struct {
	Avatar         string `json:"avatar"`         //头像
	NickName       string `json:"nickName"`       //昵称
	PraiseIcon     string `json:"praiseIcon"`     //点赞图标
	CreateTime     string `json:"createTime"`     //创建时间
	ImageUrl       string `json:"imageUrl"`       //动态缩略图
	VideoUrl       string `json:"videoUrl"`       //视频地址
	Msg            string `json:"content"`        //消息
	MsgType        int    `json:"msgType"`        //消息类型 1点赞 2我的评论回复 3我的动态评论 4我的动态其他用户评论回复 5入群申请
	CommentContent string `json:"commentContent"` //评论内容
	//ToUserId        string                  `json:"toUserId"`        //接收方用户id
	ReplyUserId     string                  `json:"replyUserId"`     //回复用户id
	ReplyNickName   string                  `json:"replyNickName"`   //回复用户昵称
	TimelineId      int64                   `json:"timelineId"`      //动态ID
	TimelineContent string                  `json:"timelineContent"` //动态内容
	UserPlaque      response.UserPlaqueInfo `json:"userPlaque"`      // 用户铭牌信息
	RoomId          string                  `json:"roomId"`
	UserId          string                  `json:"userId"`
}
