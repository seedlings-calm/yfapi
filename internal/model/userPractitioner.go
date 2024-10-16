package model

import "time"

type PractitionerExamineQuestion struct {
	Id               int    `json:"id"               description:""`
	QuestionType     int    `json:"questionType"     description:"类型(1单选,2多选,3简答,4判断题)"`
	Content          string `json:"content"          description:"问题内容"`
	PractitionerType string `json:"practitionerType" description:"从业者类型，已逗号区分"`
	Status           bool   `json:"status"           description:"是否启用"`
	// Answers          []PractitionerAnswer `json:"answers" gorm:"foreignKey:QuestionId" description:"题目的回答"`
	CreateTime time.Time `json:"createTime"       description:"创建时间"`
}

func (m *PractitionerExamineQuestion) TableName() string {
	return "t_practitioner_examine_question"
}

type PractitionerAnswer struct {
	Id         int64     `json:"id"         description:""`
	QuestionId int       `json:"questionId" description:"问题"`
	Content    string    `json:"content"    description:"答案内容"`
	IsCorrect  bool      `json:"isCorrect"  description:"是否正确"`
	SortNum    int       `json:"sortNum"    description:"答案排序值"`
	CreateTime time.Time `json:"-" description:"创建时间"`
}

func (m *PractitionerAnswer) TableName() string {
	return "t_practitioner_answer"
}

type PractitionerExamineLog struct {
	Id               string    `json:"id"               description:""`
	UserId           string    `json:"userId"           description:"用户id"`
	PractitionerType int       `json:"practitionerType" description:"从业者类型"`
	Score            int       `json:"score"            description:"分数"`
	Status           int       `json:"status"           description:"状态"`
	ExamBy           string    `json:"examBy"           description:"审核人"`
	ExamTime         time.Time `json:"examTime"         description:"审核时间"`
	ExamFeedback     string    `json:"examFeedback"     description:"审核反馈"`
	CreateTime       time.Time `json:"createTime"       description:""`
}

func (m *PractitionerExamineLog) TableName() string {
	return "t_practitioner_examine_log"
}

type PractitionerShortAnswerLog struct {
	Id              int64     `json:"id"              description:""`
	ExamineLogId    string    `json:"examineLogId"    description:"考试记录id"`
	UserId          string    `json:"userId"          description:"用户id"`
	QuestionContent string    `json:"questionContent" description:"问题"`
	AnswerContent   string    `json:"answerContent"   description:"答案"`
	CreateTime      time.Time `json:"createTime"      description:""`
}

func (m *PractitionerShortAnswerLog) TableName() string {
	return "t_practitioner_short_answer_log"
}

type PractitionerSingerExamine struct {
	Id           int64     `json:"id"           description:""`
	UserId       string    `json:"userId"       description:"用户id"`
	Img          string    `json:"img"          description:"认证资料图"`
	Description  string    `json:"description"  description:"技能介绍"`
	Audio        string    `json:"audio"        description:"语音介绍"`
	Status       int       `json:"status"       description:"状态 1待审核,2审核通过,3未通过,4取消身份"`
	ExamBy       string    `json:"examBy"       description:"审核人"`
	ExamTime     time.Time `json:"examTime"     description:"审核时间"`
	ExamFeedback string    `json:"examFeedback" description:"审核反馈"`
	CreateTime   time.Time `json:"createTime"   description:""`
}

func (m *PractitionerSingerExamine) TableName() string {
	return "t_practitioner_singer_examine"
}

type UserPractitioner struct {
	Id                int64     `json:"id"               description:""`
	RoomId            string    `json:"roomId"           description:"房间id"`
	UserId            string    `json:"userId"           description:"用户id"`
	PractitionerType  int       `json:"practitionerType" description:"从业者类型 1主持 2音乐 3咨询 4主播"`
	PractitionerBrief string    `json:"practitionerBrief"` // 从业者简介
	Status            int       `json:"status"           description:"状态 1正常,2审核中,3审核拒绝，4取消"`
	StaffName         string    `json:"staffName"        description:"操作人"`
	AbolishReason     string    `json:"abolishReason"    description:"取消原因"`
	CreateTime        time.Time `json:"createTime"       description:""`
	UpdateTime        time.Time `json:"updateTime"       description:""`
}

func (m *UserPractitioner) TableName() string {
	return "t_user_practitioner"
}
