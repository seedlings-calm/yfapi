package enum

var ReportingSence = map[int]string{
	1: "房间",
	2: "消息私聊",
	3: "个人主页",
	4: "个人动态",
	5: "动态评价",
	6: "声音派对",
}
var ReportingSenceKey = []int{
	1, 2, 3, 4, 5, 6,
}
var ReportingObject = map[int]string{
	1: "房间",
	2: "用户",
}

// 举报类型
var ReportingType = map[int]string{
	1:  "政治",
	2:  "诈骗",
	3:  "侵权",
	4:  "色情",
	5:  "辱骂诋毁",
	6:  "广告拉人",
	7:  "脱离平台交易",
	99: "其他原因",
}

// 举报类型key
var ReportingTypeKey = []int{
	1, 2, 3, 4, 5, 6, 7, 99,
}

var ReportingState = map[int]string{
	0: "待审核",
	1: "举报成功",
	2: "举报失败",
}
