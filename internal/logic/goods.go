package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"yfapi/app/handle"
	"yfapi/core/coreConfig"
	"yfapi/core/coreDb"
	"yfapi/core/coreLog"
	"yfapi/core/coreRedis"
	error2 "yfapi/i18n/error"
	"yfapi/internal/dao"
	"yfapi/internal/helper"
	"yfapi/internal/model"
	"yfapi/internal/service/accountBook"
	service_goods "yfapi/internal/service/goods"
	service_im "yfapi/internal/service/im"
	"yfapi/internal/service/user"
	"yfapi/typedef/enum"
	"yfapi/typedef/redisKey"
	request_goods "yfapi/typedef/request/goods"
	response_goods "yfapi/typedef/response/goods"
	"yfapi/typedef/response/h5"
	"yfapi/util/easy"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
)

type GoodsUseLogic struct {
}

/**
 * @description 用户更换装饰之后，判断是否在房间推送更换装饰消息
 * @param userId string true 用户ID
 * @param goodsId string strue 操作的商品ID
 */
func UpdateGoodsCallback() func(c *gin.Context, userId string, goodsId string, isUse int) {
	return func(c *gin.Context, userId string, goodsId string, isUse int) {
		key := redisKey.UserInWhichRoom(userId, helper.GetClientType(c))
		roomId, err := coreRedis.GetChatroomRedis().Get(context.Background(), key).Result()
		coreLog.Info("进入更换装扮回调中<>userId:%s goodsId:%s roomId:%s isUse:%d", userId, goodsId, roomId, isUse)
		if err == nil && roomId != "" {
			resData := new(response_goods.SpecialEffects)
			var goodsRes model.Goods
			goodsDao := dao.GoodsDao{}
			goodsRes, err = goodsDao.FirstByGoodsId(cast.ToInt(goodsId))
			if err != nil {
				return
			}
			resData.GoodsId = cast.ToString(goodsRes.Id)
			resData.GoodsTypeKey = goodsRes.TypeKey
			resData.GoodsName = goodsRes.Name
			if isUse == 2 { //装扮使用，资源放上去
				resData.GoodsIcon = helper.FormatImgUrl(goodsRes.Icon)
				resData.AnimationUrl = helper.FormatImgUrl(goodsRes.AnimationUrl)
				resData.AnimationJsonUrl = helper.FormatImgUrl(goodsRes.AnimationJsonUrl)
			}
			//推送更新装扮信息
			new(service_im.ImCommonService).Send(c, userId, nil, roomId, enum.MsgCustom, resData, enum.User_Goods_Change)
		}
	}

}

func (GoodsUseLogic) GoodsTypesList(c *gin.Context) (res []*response_goods.GoodsTypesListRes) {
	GoodsUseLogic := dao.GoodsUseDao{}
	ids, err := GoodsUseLogic.FindGoodsTypeIds(true)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	if len(ids) > 0 {
		GoodsTypeDao := dao.GoodsType{}
		res, err = GoodsTypeDao.FindsByIds(ids, true)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		if len(res) > 0 {
			for _, v := range res {
				v.Icon = coreConfig.GetHotConf().ImagePrefix + v.Icon
			}
		}
	}
	return
}

func (GoodsUseLogic) GoodsListByTypes(c *gin.Context, typeId int) (res []*response_goods.GoodsListByTypesRes) {
	goodsUseDao := dao.GoodsUseDao{}
	var err error
	res, err = goodsUseDao.FindByGoodsType(typeId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	if len(res) > 0 {
		for _, v := range res {
			v.AnimationJsonUrl = coreConfig.GetHotConf().ImagePrefix + v.AnimationJsonUrl
			v.AnimationUrl = coreConfig.GetHotConf().ImagePrefix + v.AnimationUrl
			v.Icon = coreConfig.GetHotConf().ImagePrefix + v.Icon
		}
	}
	return
}

func (GoodsUseLogic) GoodsListToUser(c *gin.Context) (res []*response_goods.GoodsListToUserRes) {
	userId := handle.GetUserId(c)
	userGoodsDao := dao.UserGoodsDao{}
	useGoodsDao := dao.GoodsUseDao{}
	useGoodsTypeIds, err := useGoodsDao.FindGoodsTypeIds(false)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	//根据用户装扮获取类型ID，防止类型删除，装扮展示不出来
	userTypeIds, _ := userGoodsDao.GetGoodsTypeByUser(userId)
	if len(useGoodsTypeIds) == 0 && len(userTypeIds) == 0 {
		return
	}
	useGoodsTypeIds = append(useGoodsTypeIds, userTypeIds...)
	typeList, err := dao.GoodsType{}.FindsByIds(useGoodsTypeIds, false)
	if err != nil {
		return
	}
	goodsList, err := userGoodsDao.FindListByUser(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}

	for _, v := range typeList {
		item := new(response_goods.GoodsListToUserRes)
		item.Id = v.Id
		item.Icon = coreConfig.GetHotConf().ImagePrefix + v.Icon
		item.Keys = v.Keys
		item.Name = v.Name
		item.Sort = v.Sort
		item.Nums = 1000
		item.Status = v.Status
		item.GoodsList = make([]*response_goods.UserGoods, 0)
		for _, v1 := range goodsList {
			if v1.GoodsTypeId == v.Id {
				v1.Icon = coreConfig.GetHotConf().ImagePrefix + v1.Icon
				v1.AnimationUrl = coreConfig.GetHotConf().ImagePrefix + v1.AnimationUrl
				if v1.AnimationJsonUrl != "" {
					v1.AnimationJsonUrl = coreConfig.GetHotConf().ImagePrefix + v1.AnimationJsonUrl
				}
				item.GoodsList = append(item.GoodsList, &v1)
				if v1.Nums <= 100 { //商品展示红点， 分类赋予红点效果
					item.Nums = v1.Nums
				}
			}
		}
		res = append(res, item)
	}
	return
}

func (GoodsUseLogic) UseGoodsToUser(c *gin.Context, goodId string) {
	userId := handle.GetUserId(c)
	userGoodsDao := dao.UserGoodsDao{}
	goodsInfo, err := userGoodsDao.GetUserGoodsOne(userId, goodId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	if goodsInfo.IsUse == 1 {
		if userGoodsDao.IsExpireTime(goodsInfo.ExpireTime) {
			panic(error2.I18nError{
				Code: error2.ErrCodeUserGoodsExpireErr,
				Msg:  nil,
			})
		}
	}
	userGoodsService := service_goods.UserGoods{}
	err = userGoodsService.UpdateUserGoods(c, goodsInfo, helper.GetClientType(c), UpdateGoodsCallback())
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserGoodsUseErr,
		})
	}
}
func (GoodsUseLogic) UserGoodsDel(c *gin.Context, goodId string) {
	userId := handle.GetUserId(c)
	userGoodsDao := dao.UserGoodsDao{}
	goodsInfo, err := userGoodsDao.GetUserGoodsOne(userId, goodId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	//如果未过期，不处理
	if !userGoodsDao.IsExpireTime(goodsInfo.ExpireTime) {
		return
	}
	err = userGoodsDao.DelExpireTimeGoods(goodsInfo.Id)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrCodeUserGoodsDelErr,
			Msg:  nil,
		})
	}
}

func (GoodsUseLogic) BuyGoodsToUser(c *gin.Context, req *request_goods.BuyGoodsToUserReq) {
	userId := handle.GetUserId(c)
	goodsUseDao := dao.GoodsUseDao{}
	goodsInfo, err := goodsUseDao.FirstByGoodsId(req.GoodsId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	var (
		//计算支付总价
		totalAmount int
	)
	if req.Days != 1 {
		moneySlice := strings.Split(goodsInfo.Moneys, ",")
		var money int
		if len(moneySlice) >= 1 {
			switch req.Days {
			case 7:
				money = cast.ToInt(moneySlice[0])
			case 15:
				money = cast.ToInt(moneySlice[1])
			case 30:
				money = cast.ToInt(moneySlice[2])
			default:
				panic(error2.I18nError{
					Code: error2.ErrorCodeParam,
					Msg:  nil,
				})
			}
			totalAmount = req.Num * money
		}

	} else {
		totalAmount = req.Num * goodsInfo.Money
	}
	if totalAmount <= 0 {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	// 检查账户余额是否充足
	accountDao := new(dao.UserAccountDao)
	fromUserAccount, err := accountDao.GetUserAccountByUserId(userId)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeReadDB,
			Msg:  nil,
		})
	}
	if easy.StringToDecimal(fromUserAccount.DiamondAmount).LessThan(decimal.NewFromInt(int64(totalAmount))) {
		panic(error2.I18nError{
			Code: error2.ErrorCodeDiamondNotEnough,
			Msg:  nil,
		})
	}
	tx := coreDb.GetMasterDb().Begin()

	stringTotalAmount := cast.ToString(totalAmount)
	// 生成订单
	orderInfo := &model.Order{
		OrderId:         new(accountBook.Order).OrderNum(accountBook.ORDER_SC),
		UserId:          userId,
		RoomId:          "0",
		GuildId:         "0",
		Gid:             cast.ToString(req.GoodsId),
		TotalAmount:     stringTotalAmount,
		PayAmount:       stringTotalAmount,
		DiscountsAmount: "0",
		Num:             req.Num,
		Currency:        accountBook.CURRENCY_DIAMOND,
		OrderType:       accountBook.ChangeDiamondMallConsumption,
		OrderStatus:     1,
		PayStatus:       1,
		Note:            fmt.Sprintf("%d*%d", req.Days, req.Num),
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
		StatDate:        time.Now().Format(time.DateOnly),
	}
	err = tx.Create(orderInfo).Error
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	//添加商品到背包
	userGoodsService := service_goods.UserGoods{}
	var isUse int
	isUse, err = userGoodsService.AddGoods(tx, userId, cast.ToString(goodsInfo.GoodsId), cast.ToString(goodsInfo.GoodsTypeId), goodsInfo.GoodsTypeKey, req.Days*req.Num, orderInfo.OrderId)
	if err != nil {
		tx.Rollback()
		panic(error2.I18nError{
			Code: error2.ErrorCodeSystemBusy,
			Msg:  nil,
		})
	}
	// 扣除购买人钻石
	user.UpdateUserAccount(&user.UpdateAccountParam{
		Tx:        tx,
		UserId:    userId,
		Gid:       cast.ToString(req.GoodsId),
		Num:       req.Num,
		Currency:  accountBook.CURRENCY_DIAMOND,
		FundFlow:  accountBook.FUND_OUTFLOW,
		Amount:    stringTotalAmount,
		OrderId:   orderInfo.OrderId,
		OrderType: accountBook.ChangeDiamondMallConsumption,
		RoomId:    "0",
		GuildId:   "0",
		Note:      fmt.Sprintf("%d*%d", req.Days, req.Num),
	})
	tx.Commit()
	UpdateGoodsCallback()(c, userId, cast.ToString(goodsInfo.Id), isUse)
}

func (GoodsUseLogic) GoodsAll(c *gin.Context) (resp []*response_goods.GoodsAllRes) {

	bytes, _ := coreRedis.GetUserRedis().Get(c, redisKey.GoodsAllCacheKey()).Result()
	if bytes != "" {
		err := json.Unmarshal([]byte(bytes), &resp)
		if err == nil {
			return
		}
	}
	GoodsUseLogic := dao.GoodsUseDao{}
	ids, err := GoodsUseLogic.FindGoodsTypeIds(false)
	if err != nil {
		panic(error2.I18nError{
			Code: error2.ErrorCodeParam,
			Msg:  nil,
		})
	}
	if len(ids) > 0 {
		GoodsTypeDao := dao.GoodsType{}
		res, err := GoodsTypeDao.FindsByIds(ids, true)
		if err != nil {
			panic(error2.I18nError{
				Code: error2.ErrorCodeParam,
				Msg:  nil,
			})
		}
		if len(res) > 0 {
			goodsList, _ := GoodsUseLogic.FindByGoodsType(0)

			for _, v := range res {
				item := new(response_goods.GoodsAllRes)
				item.Id = v.Id
				item.Icon = coreConfig.GetHotConf().ImagePrefix + v.Icon
				item.Keys = v.Keys
				item.Name = v.Name
				item.Sort = v.Sort
				item.GoodsList = make([]*response_goods.GoodsListByTypesRes, 0)
				for _, v1 := range goodsList {
					if v1.GoodsTypeId == v.Id {
						v1.Icon = coreConfig.GetHotConf().ImagePrefix + v1.Icon
						v1.AnimationUrl = coreConfig.GetHotConf().ImagePrefix + v1.AnimationUrl
						if v1.AnimationJsonUrl != "" {
							v1.AnimationJsonUrl = coreConfig.GetHotConf().ImagePrefix + v1.AnimationJsonUrl
						}
						item.GoodsList = append(item.GoodsList, v1)
					}
				}
				resp = append(resp, item)
			}
		}
	}
	jsonString, err := json.Marshal(resp)
	if err != nil {
		return
	}
	coreRedis.GetUserRedis().Set(c, redisKey.GoodsAllCacheKey(), string(jsonString), time.Second*30)
	return
}

func (GoodsUseLogic) DelRedDot(c *gin.Context, goodsId string) {
	userId := handle.GetUserId(c)
	userGoodsDao := dao.UserGoodsDao{}
	err := userGoodsDao.DelRedHot(userId, goodsId)
	if err != nil {
		coreLog.Error(err)
	}
}

func (GoodsUseLogic) GetUserGoodsAccount(c *gin.Context) (res h5.UserAccountRes) {
	userId := handle.GetUserId(c)
	account := dao.UserAccountDao{}
	resAccount, _ := account.GetUserAccountByUserId(userId)
	res.UserId = resAccount.UserId
	res.DiamondAmount = resAccount.DiamondAmount
	userInfo := user.GetUserBaseInfo(userId)
	res.Avatar = userInfo.Avatar
	res.Nickname = userInfo.Nickname
	res.UserNo = userInfo.UserNo
	res.Uid32 = cast.ToInt32(userInfo.OriUserNo)
	goodsRes := dao.UserGoodsDao{}.GetUserGoods(userId)
	if len(goodsRes) > 0 {
		res.MyGoodsIsHave = true
		for _, v := range goodsRes {
			if v.Nums <= 100 {
				res.MyGoodsIsRed = true
			}
		}
	}
	return
}
