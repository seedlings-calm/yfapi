package user

type PractitionerExamineQuestion struct {
	Id           int                  `json:"id"               description:""`
	QuestionType int                  `json:"questionType"     description:"类型(1单选,2多选,3简答,4判断题)"`
	Content      string               `json:"content"          description:"问题内容"`
	Answers      []PractitionerAnswer `json:"answers" gorm:"foreignKey:QuestionId" description:"题目的回答"`
}

func (m *PractitionerExamineQuestion) TableName() string {
	return "t_practitioner_examine_question"
}

type PractitionerAnswer struct {
	Id         int64  `json:"id"         description:""`
	QuestionId int    `json:"questionId" description:"问题"`
	Content    string `json:"content"    description:"答案内容"`
	SortNum    int    `json:"sortNum"    description:"答案排序值"`
}

func (m *PractitionerAnswer) TableName() string {
	return "t_practitioner_answer"
}

type ApplyJoinResultResponse struct {
	IsMusician  bool `json:"isMusician"`  //是否有音乐人身份 false：未考核   true: 考核通过
	IsCompere   bool `json:"isCompere"`   //是否有主持人身份 false：未考核   true: 考核通过
	IsCounselor bool `json:"isCounselor"` //是否有咨询师身份 false：未考核   true: 考核通过
	IsAnchor    bool `json:"isAnchor"`    //是否有主播身份 false：未考核   true: 考核通过
}

type CerdAuthResponse struct {
	IsTrueName   int  `json:"isTrueName"`   //是否实名  1未认证 2已认证 3认证中
	IsMobile     bool `json:"isMobile"`     //是否绑定手机号	false 未绑定， true 已绑定
	IsCerd       int  `json:"isCerd"`       // 当前身份的状态  1:未考核  2:考核中  3: 考核通过
	ExamineLevel int  `json:"examineLevel"` //考核等级 1:20题基础，3:2题简答
	AnswerNum    int  `json:"answerNum"`    //已考试次数
	QuestionNum  int  `json:"questionNum"`  //考试题数量

}
