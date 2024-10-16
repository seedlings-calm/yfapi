package dao

import (
	"yfapi/core/coreDb"
	"yfapi/internal/model"
	"yfapi/typedef/response/roomOwner"
	response_user "yfapi/typedef/response/user"

	"gorm.io/gorm"
)

type DaoPractitionerExamineQuestion struct {
	ExamineLevel     int
	PractitionerType string
	Limit            int
}

func (u *DaoPractitionerExamineQuestion) Find() (result []response_user.PractitionerExamineQuestion, err error) {

	db := coreDb.GetMasterDb().Table("t_practitioner_examine_question as tp")
	if u.ExamineLevel == 1 { //如果是基础题，可以选择类型1,2,4,获取回答选项，简答题不需要,
		db = db.Preload("Answers").Where("FIND_IN_SET(3,tp.question_type) = 0")
	} else { //简答题，只能选择 3
		db = db.Where("FIND_IN_SET(3,tp.question_type)")
	}
	err = db.Where("tp.status = 1").Where("FIND_IN_SET(?,tp.practitioner_type)", u.PractitionerType).
		Limit(u.Limit).Order("RAND()").Find(&result).Error
	return
}

type DaoPractitionerAnswer struct {
	QuestionID     int
	CorrectAnswers string
}

// 根据question_id 集合查询答案
func (u *DaoPractitionerAnswer) Find(ids []int) (result []DaoPractitionerAnswer, err error) {
	err = coreDb.GetMasterDb().
		Model(model.PractitionerAnswer{}).
		Select("question_id, GROUP_CONCAT(id SEPARATOR ',') AS correct_answers").
		Where("question_id in (?) and is_correct = ?", ids, true).
		Group("question_id").
		Scan(&result).Error
	return
}

type DaoPractitionerExamineLog struct {
}

func (u *DaoPractitionerExamineLog) Creates(tx *gorm.DB, data []*model.PractitionerExamineLog) (err error) {
	for _, v := range data {
		err = tx.Create(&v).Error
		if err != nil {
			return
		}
	}
	return
}

func (u *DaoPractitionerExamineLog) IsLog(userId, status string, types int) (res model.PractitionerExamineLog, err error) {
	err = coreDb.GetMasterDb().Model(model.PractitionerExamineLog{}).Where("user_id = ? and practitioner_type = ? and status in (?)", userId, types, status).First(&res).Error
	return
}

type DaoPractitionerShortAnswerLog struct {
}

func (u *DaoPractitionerShortAnswerLog) Creates(tx *gorm.DB, data []*model.PractitionerShortAnswerLog) (err error) {
	db := coreDb.GetMasterDb().Model(model.PractitionerShortAnswerLog{})
	for _, v := range data {
		err = db.Create(&v).Error
		if err != nil {
			return err
		}
	}
	return
}

type DaoPractitionerSingerExamine struct {
}

func (u *DaoPractitionerSingerExamine) Create(data model.PractitionerSingerExamine) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(&data).Error
	return
}

func (u *DaoPractitionerSingerExamine) IsLog(userId, status string) (res model.PractitionerSingerExamine, err error) {
	err = coreDb.GetMasterDb().Model(model.PractitionerSingerExamine{}).Where("user_id = ? and status in (?)", userId, status).First(&res).Error
	return
}

type DaoUserPractitioner struct {
}

func (du *DaoUserPractitioner) Find(userId string, roomId string) (res []model.UserPractitioner, err error) {
	err = coreDb.GetMasterDb().Model(model.UserPractitioner{}).Where("room_id = ? and user_id = ? and status=1", roomId, userId).Find(&res).Error
	return
}

func (du *DaoUserPractitioner) FindByRoomId(roomId string) (res []*roomOwner.PersonListRes, err error) {
	db := coreDb.GetMasterDb().Table("t_user_practitioner as tup").
		Joins("left join t_user as tu on tup.user_id = tu.id").
		Where("tup.room_id = ? and tup.status = 1", roomId)
	err = db.Select("group_concat(tup.practitioner_type) as id_cards,tu.id as user_id,tu.user_no,tu.nickname,tu.avatar").Group("tup.user_id").Scan(&res).Error
	return
}

// GetGuildPractitionerRecord
//
//	@Description: 查询用户公会从业者记录
//	@receiver du
//	@param userId string -
//	@param guildId string -
//	@return res -
//	@return err -
func (du *DaoUserPractitioner) GetGuildPractitionerRecord(userId, guildId string) (res []model.UserPractitioner, err error) {
	err = coreDb.GetSlaveDb().Table("t_user_practitioner up").Joins("left join t_room r on r.id=up.room_id").
		Where("r.guild_id=? and up.user_id=? and up.status in (1,4)", guildId, userId).Scan(&res).Error
	return
}
