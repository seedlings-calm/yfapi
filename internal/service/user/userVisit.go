package user

import (
	"yfapi/internal/dao"
	"yfapi/typedef/enum"
	"yfapi/util/easy"
)

// RecordUserVisit 记录用户访问
func RecordUserVisit(userId, targetUserId string) {
	// 过滤系统用户
	if easy.InArray(targetUserId, enum.OfficialUserIdList) {
		return
	}
	_ = new(dao.UserVisitDao).RecordUserVisit(userId, targetUserId, false)
}
