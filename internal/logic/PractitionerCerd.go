package logic

import (
	"log"
	"strconv"
	"yfapi/app/handle"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	response_user "yfapi/typedef/response/user"

	"github.com/gin-gonic/gin"
)

type PractitionerCerd struct{}

func (p PractitionerCerd) CerdAuth(types string, c *gin.Context) response_user.CerdAuthResponse {
	resp := response_user.CerdAuthResponse{
		IsTrueName:   1,
		IsMobile:     false,
		IsCerd:       1,
		ExamineLevel: 1,
		AnswerNum:    0,
		QuestionNum:  0,
	}

	userId := handle.GetUserId(c)
	val := coreRedis.GetUserRedis().Get(c, redisKey.UserPractitionerQuestion(userId, types)).Val()
	if val != "" {
		resp.ExamineLevel = 3
	}

	num := coreRedis.GetUserRedis().Get(c, redisKey.UserPractitionerQuestionNums(userId, types)).Val()
	if num != "" {
		newNum, _ := strconv.Atoi(num)
		resp.AnswerNum = newNum
	}

	user := &model.User{
		Id: userId,
	}
	userDao := &dao.UserDao{}
	userInfo, _ := userDao.FindOne(user)
	resp.IsTrueName = userInfo.RealNameStatus
	if userInfo.Mobile != "" {
		resp.IsMobile = true
	}
	//根据身份处理状态
	tv, _ := strconv.Atoi(types)
	switch tv {
	case enum.UserPractitionerCompere, enum.UserPractitionerAnchor, enum.UserPractitionerCounselor: //考核一套流程
		dao := dao.DaoPractitionerExamineLog{}
		res, _ := dao.IsLog(user.Id, "1", tv)
		if res.Id != "" {
			resp.IsCerd = 2
		}
	case enum.UserPractitionerMusician: //音乐人特殊考核
		dao := dao.DaoPractitionerSingerExamine{}
		res, _ := dao.IsLog(user.Id, "1")
		if res.Status == 1 {
			resp.IsCerd = 2
		}
	default:
		log.Printf("从业者身份存储的类型不匹配:%#v", tv)
		return resp
	}

	userCerddao := &dao.DaoUserPractitionerCerd{
		UserId: user.Id,
	}
	res, _ := userCerddao.First(tv)
	//从业者身份存储中如果有此身份并且是可用状态
	if res.Id > 0 {
		resp.IsCerd = 3
	}

	questionDao := &dao.DaoPractitionerExamineQuestion{
		ExamineLevel:     resp.ExamineLevel,
		PractitionerType: types,
		Limit:            20,
	}
	quesRes, _ := questionDao.Find()
	resp.QuestionNum = len(quesRes)
	return resp
}
func (p PractitionerCerd) ApplyJoinResult(c *gin.Context) (resp response_user.ApplyJoinResultResponse) {

	dao := &dao.DaoUserPractitionerCerd{
		UserId: handle.GetUserId(c),
	}
	res, err := dao.Find()
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodePractitionerResult,
			Msg:  nil,
		})
	}
	for _, v := range res {
		switch v.PractitionerType {
		case enum.UserPractitionerCompere:
			resp.IsCompere = true
		case enum.UserPractitionerCounselor:
			resp.IsCounselor = true
		case enum.UserPractitionerAnchor:
			resp.IsAnchor = true
		case enum.UserPractitionerMusician:
			resp.IsMusician = true
		default:
			log.Printf("从业者身份存储的类型不匹配:%#v", v)
		}
	}
	return resp
}
