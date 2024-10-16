package enum

// 动态类别
const (
	// TimelineListTypeChosen 动态列表类型 精选
	TimelineListTypeChosen = iota + 1
	// TimelineListTypeLatest 动态列表类型 最新
	TimelineListTypeLatest
	// TimelineListTypeFollow 动态列表类型 关注
	TimelineListTypeFollow
)

// 动态状态
const (
	// TimelineStatusHidden 动态状态 隐藏
	TimelineStatusHidden = 0
	// TimelineStatusNormal 动态状态 正常
	TimelineStatusNormal = 1
	// TimelineStatusWaitExamine 动态状态 待审核
	TimelineStatusWaitExamine = 2
	// TimelineStatusRefuseExamine 动态状态 审核拒绝
	TimelineStatusRefuseExamine = 3
	// TimelineStatusDelete 动态状态 删除
	TimelineStatusDelete = 4
)

const (
	TimelineImgType   = 1 //图文动态
	TimelineVideoType = 2 //视频动态
)

const (
	// VideoCoverImgSuffix 视频url封面图后缀
	VideoCoverImgSuffix = "?x-oss-process=video/snapshot,t_0,f_jpg,w_200,h_200,m_fast"
)
