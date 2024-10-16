package logic

import (
	"encoding/json"
	"strconv"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreDb"
	"yfapi/core/coreRedis"
	"yfapi/core/coreSnowflake"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/model"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_user "yfapi/typedef/request/user"
	response_user "yfapi/typedef/response/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Practitioner struct {
}

/**
 * @description  判断用户的从业者身份
 * @param userId 用户id  - types 从业者身份
 * @returns bool 已有此身份为true
 */
func (p Practitioner) IsCerdStatus(userId string, types int) bool {
	cerdDao := dao.DaoUserPractitionerCerd{
		UserId: userId,
	}
	cerdres, _ := cerdDao.First(types)
	if cerdres.UserId == userId {
		return true
	}
	if types == enum.UserPractitionerMusician {
		singerDao := dao.DaoPractitionerSingerExamine{}
		singerres, _ := singerDao.IsLog(userId, "1,2")
		if singerres.Id > 0 {
			return true
		}
	} else {
		answerDao := dao.DaoPractitionerExamineLog{}
		answerres, _ := answerDao.IsLog(userId, "1,2", types)
		if answerres.Id != "" {
			return true
		}
	}
	return false
}

func (p Practitioner) GetQuestion(types string, c *gin.Context) (res []response_user.PractitionerExamineQuestion) {
	userId := handle.GetUserId(c)
	isTypes, _ := strconv.Atoi(types)
	//身份状态检查
	ok := p.IsCerdStatus(userId, isTypes)
	if ok {
		panic(error2.I18nError{Code: error2.ErrorCodeIsTrue})
	}

	examineLevel := 1 //基础答题（20题），简答答题（2题） 1:基础 ，3:简答
	dao := &dao.DaoPractitionerExamineQuestion{
		ExamineLevel:     examineLevel,
		PractitionerType: types,
		Limit:            20,
	}
	val := coreRedis.GetUserRedis().Get(c, redisKey.UserPractitionerQuestion(userId, types)).Val()
	if val != "" {
		examineLevel = 3
		dao.ExamineLevel = examineLevel
		dao.Limit = 2
	}
	//只有基础考题 需要限制每日三次 TODO: 测试期间注释
	if examineLevel == 1 {
		//考试次数检查
		key := redisKey.UserPractitionerQuestionNums(userId, types)
		//   增加计数器,如果不存在则会初始化为0，然后加1
		count, err := coreRedis.GetUserRedis().Incr(c, key).Result()
		if err != nil {
			panic(error2.I18nError{Code: error2.ErrorCodePractitionerAnswer, Msg: map[string]interface{}{
				"times": 3,
			}})
		}
		//表示初始化，增加过期时间
		if count == 1 {
			_, err = coreRedis.GetUserRedis().Expire(c, key, 24*60*time.Minute).Result()
			if err != nil {
				panic(&error2.I18nError{Code: error2.ErrorCodePractitionerAnswer, Msg: map[string]interface{}{
					"times": 3,
				}})
			}
		}
		// 检查是否超过最大请求次数
		if count > int64(3) {
			panic(&error2.I18nError{Code: error2.ErrorCodePractitionerAnswer, Msg: map[string]interface{}{
				"times": 3,
			}})
		}
	}

	for i := 0; i < 5; i++ {
		if dao.Limit != 0 {
			result, err := dao.Find()
			if err != nil {
				panic(error2.I18nError{Code: error2.ErrorCodeUserPractitioner})
			}
			res = append(res, result...)
			dao.Limit -= len(result)
		} else {
			break
		}
	}
	if examineLevel == 1 { //基础考核，需要考虑是否有20题
		if len(res) != 20 {
			panic(error2.I18nError{Code: error2.ErrorCodeUserPractitioner})
		}
	}

	return
}

func (p Practitioner) PullAnswer(req *request_user.PractitionerAwnserReq, c *gin.Context) int {
	if len(req.Answers) != 20 {
		panic(error2.I18nError{Code: error2.ErrorCodeUserPractitionerAwnser})
	}

	getIdMap := func(items *request_user.PractitionerAwnserReq) []int {
		ids := make([]int, 0)
		for _, item := range items.Answers {
			ids = append(ids, item.Id)
		}
		return ids
	}
	ids := getIdMap(req)
	dao := dao.DaoPractitionerAnswer{}
	if len(ids) == 0 {
		panic(error2.I18nError{Code: error2.ErrorCodeUserPractitionerAwnser})
	}

	result, err := dao.Find(ids)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserPractitionerAwnser})
	}
	var Num = 100 //分数
	//遍历所有题目的正确答案，
	for _, v := range result {
		for _, av := range req.Answers {
			if v.QuestionID == av.Id && v.CorrectAnswers != av.Awnsers {
				Num -= 5
			}
		}
	}
	if Num >= 80 {
		userId := handle.GetUserId(c)
		ok := p.IsCerdStatus(userId, req.Types)
		if ok {
			panic(error2.I18nError{Code: error2.ErrorCodeIsTrue})
		}

		value := model.PractitionerExamineLog{
			Id:               coreSnowflake.GetSnowId(),
			UserId:           userId,
			Score:            Num,
			PractitionerType: req.Types,
			Status:           1,
			ExamTime:         time.Now(),
			CreateTime:       time.Now(),
		}
		val, _ := json.Marshal(value)
		coreRedis.GetUserRedis().Set(c, redisKey.UserPractitionerQuestion(userId, strconv.Itoa(req.Types)), val, 0)
	}
	return Num
}

func (p Practitioner) PullShortAnswer(req request_user.PractitionerShortAwnserReq, c *gin.Context) error {
	if len(req.Answers) != 2 {
		panic(error2.I18nError{Code: error2.ErrorCodeUserPractitionerAwnser})
	}
	userId := handle.GetUserId(c)
	times := time.Now()
	res, err := coreRedis.GetUserRedis().Get(c, redisKey.UserPractitionerQuestion(userId, strconv.Itoa(req.Types))).Result()
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserPractitionerAwnser})
	}
	answer := &model.PractitionerExamineLog{}
	err = json.Unmarshal([]byte(res), answer)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserPractitionerAwnser})
	}
	answer.ExamTime = times
	answer.CreateTime = times
	var creates = make([]*model.PractitionerShortAnswerLog, 0)
	for _, v := range req.Answers {
		item := model.PractitionerShortAnswerLog{
			ExamineLogId:    answer.Id,
			UserId:          userId,
			QuestionContent: v.QuestionName,
			AnswerContent:   v.Awnsers,
			CreateTime:      times,
		}
		creates = append(creates, &item)
	}
	err = coreDb.GetMasterDb().Transaction(func(tx *gorm.DB) error {
		logDao := dao.DaoPractitionerExamineLog{}
		logDaos := make([]*model.PractitionerExamineLog, 0)
		logDaos = append(logDaos, answer)
		err := logDao.Creates(tx, logDaos)
		if err != nil {
			return err
		}
		ShortDao := dao.DaoPractitionerShortAnswerLog{}
		err = ShortDao.Creates(tx, creates)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUserPractitionerAwnser})
	}
	coreRedis.GetUserRedis().Del(c, redisKey.UserPractitionerQuestion(userId, strconv.Itoa(req.Types)))
	return nil
}

func (p Practitioner) PullMusic(req *request_user.PractitionerMusicianReq, c *gin.Context) error {
	userId := handle.GetUserId(c)
	ok := p.IsCerdStatus(userId, enum.UserPractitionerMusician)
	if ok {
		panic(error2.I18nError{Code: error2.ErrorCodeIsTrue})
	}
	dao := dao.DaoPractitionerSingerExamine{}
	data := model.PractitionerSingerExamine{
		UserId:      userId,
		Img:         req.Image,
		Audio:       req.Audio,
		Description: req.Description,
		Status:      1,
		ExamTime:    time.Now(),
		CreateTime:  time.Now(),
	}
	err := dao.Create(data)
	if err != nil {
		panic(error2.I18nError{Code: error2.ErrorCodeUpdateDB})
	}
	return err
}
