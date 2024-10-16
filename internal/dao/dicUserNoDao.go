package dao

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/model"
	"yfapi/typedef/redisKey"
)

var (
	userNoHead     = 10 //用户ID头
	userNoBodyBits = 99999
	createCount    = 10000 //每次生成数量
)

type DicUserNoDao struct {
}

// Create 添加
func (u *DicUserNoDao) Create(data *model.DicUserNo) (err error) {
	err = coreDb.GetMasterDb().Model(data).Create(data).Error
	return
}

// Update 修改
func (u *DicUserNoDao) Update(data model.DicUserNo) (err error) {
	err = coreDb.GetMasterDb().Model(&model.DicUserNo{Id: data.Id}).Updates(data).Error
	return
}

// Update 修改昵称id使用状态
func (u *DicUserNoDao) UpdateNickNameStatus(data model.DicUserNickName) (err error) {
	err = coreDb.GetMasterDb().Model(&model.DicUserNickName{Id: data.Id}).Updates(data).Error
	return
}

// FindOne 条件查询
func (u *DicUserNoDao) FindOne(param *model.DicUserNo) (data *model.DicUserNo, err error) {
	data = new(model.DicUserNo)
	err = coreDb.GetMasterDb().Where(param).First(data).Error
	return
}

// FindList 查询列表
func (u *DicUserNoDao) FindList(param *model.DicUserNo) (result []model.DicUserNo, err error) {
	err = coreDb.GetMasterDb().Where(param).Find(&result).Error
	return
}

// FindByIds 根据ids查询结果
func (u *DicUserNoDao) FindByIds(ids []string) (result []model.DicUserNo) {
	coreDb.GetMasterDb().Find(&result, ids)
	return
}

// GetRandData 获取随机一条数据
func (u *DicUserNoDao) GetRandData() (result model.DicUserNo) {
	coreDb.GetMasterDb().Where("status = ?", 0).Order("RAND()").First(&result)
	//进行自动生成
	if result.Id == 0 {
		success, unlock, err := coreRedis.UserLock(context.Background(), redisKey.CreateUserNoLock(), time.Second*10)
		if err != nil || !success {
			coreLog.Error("GetRandData err:%+v", err)
			panic(error2.I18nError{
				Code: error2.ErrorCodeOperationFrequent,
				Msg:  nil,
			})
		}
		defer unlock()
		Gene2(0)
		coreDb.GetMasterDb().Where(&model.DicUserNo{
			Status: 0,
		}).Order("RAND()").First(&result)
	}
	return
}

// 获取随机一条数据(昵称随机数)
func (u *DicUserNoDao) GetRandUserNickName() (result model.DicUserNickName) {
	coreDb.GetMasterDb().Where("status = ?", 0).Order("RAND()").First(&result)
	//进行自动生成
	success, unlock, err := coreRedis.UserLock(context.Background(), redisKey.CreateUserNoLock(), time.Second*10)
	if err != nil || !success {
		coreLog.Error("GetRandData err:%+v", err)
		panic(error2.I18nError{
			Code: error2.ErrorCodeOperationFrequent,
			Msg:  nil,
		})
	}
	defer unlock()
	if result.Id == 0 {
		GeneNickName()
		coreDb.GetMasterDb().Where(&model.DicUserNickName{
			Status: 0,
		}).Order("RAND()").First(&result)
	}
	return
}

// 判断存在四位重复数字
func hasFourIdenticalDigits(numStr string) bool {
	count := make(map[rune]int)
	for _, char := range numStr {
		if char >= '0' && char <= '9' {
			count[char]++
			if count[char] >= 4 {
				return true
			}
		}
	}
	return false
}

func Gene2(count int) {
	var lastData model.DicUserNo
	coreDb.GetMasterDb().Last(&lastData)
	lastUserNo := lastData.UserNo
	var dataList []model.DicUserNo
	headNum := userNoHead
	bodyNum := 0
	if len(lastUserNo) > 0 {
		last := strings.Split(lastUserNo, "")
		headNum = cast.ToInt(last[0] + last[1])
		bodyNumStr := ""
		for i := 2; i < len(last); i++ {
			bodyNumStr += last[i]
		}
		atoi, _ := strconv.Atoi(bodyNumStr)
		bodyNum = atoi + 1
	}
	if count == 0 {
		count = createCount
	}
	for i := 0; i < count; i++ {
		if bodyNum >= userNoBodyBits {
			bodyNum = 0
			headNum++
		}
		bodyNumStr := fmt.Sprintf("%05d", bodyNum)
		full := cast.ToString(headNum) + bodyNumStr
		if !hasFourIdenticalDigits(full) {
			dataList = append(dataList, model.DicUserNo{
				UserNo: full,
				Status: 0,
			})
		}
		bodyNum++
	}
	coreDb.GetMasterDb().Model(&model.DicUserNo{}).CreateInBatches(dataList, len(dataList))
}

// 随机生成4位及以上的昵称id
func GeneNickName() {
	result := model.DicUserNickName{}
	coreDb.GetMasterDb().Order("id DESC").Limit(1).Find(&result)
	// 生成昵称编号
	num, _ := strconv.Atoi(result.UserNickName)
	num += 1
	var lastData []model.DicUserNickName
	for i := num; i < num+1000; i++ {
		nickname := strconv.Itoa(i)
		// 如果昵称不足4位，用0补足到4位
		if len(nickname) < 4 {
			nickname = fmt.Sprintf("%04s", nickname)
		}
		lastData = append(lastData, model.DicUserNickName{
			UserNickName: nickname,
			Status:       0,
		})
	}
	coreDb.GetMasterDb().Model(&model.DicUserNickName{}).CreateInBatches(lastData, len(lastData))
}

// 旧方法，有缺陷
func Gene() {
	var lastData model.DicUserNo
	coreDb.GetMasterDb().Last(&lastData)
	lastUserNo := lastData.UserNo
	var dataList []model.DicUserNo
	prefix := "10"

	if len(lastUserNo) > 0 {
		last := strings.Split(lastUserNo, "")
		prefix = last[0] + last[1]
	}
	prefixInt := cast.ToInt(prefix)
	first := 1
	for i := 0; i < 10000; i++ {
		numStr := ""
		num := first + i
		if num < 10000 {
			numStr = fmt.Sprintf("%04d", num)
			numStr = prefix + numStr
		} else {
			numStr = strconv.Itoa(prefixInt + num)
		}

		numChars := strings.Split(numStr, "")
		eqData := ""
		eqNum := 0
		for _, item := range numChars {
			if eqNum == 4 {
				break
			}
			if eqData == item {
				eqNum++
			} else {
				eqNum = 1
			}
			eqData = item
		}
		if eqNum >= 4 {
			continue
		}
		dataList = append(dataList, model.DicUserNo{
			UserNo: numStr,
			Status: 0,
		})
	}
	coreDb.GetMasterDb().Model(&model.DicUserNo{}).CreateInBatches(dataList, len(dataList))
}
