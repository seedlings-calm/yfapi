package model

import "time"

type UserLevelPrivilegeItems struct {
	ID         int       `json:"id" gorm:"column:id"`
	Level      int       `json:"level" gorm:"column:level"`            // 等级id
	GoodsId    int       `json:"goodsId" gorm:"column:goods_id"`       // 物品id
	LevelType  int       `json:"levelType" gorm:"column:level_type"`   // 1:lv 2:vip 3:star
	ExpireDate int       `json:"expireDate" gorm:"column:expire_date"` // 有效期(天)
	Explain    string    `json:"explain" gorm:"column:explain"`        // 说明
	StaffName  string    `json:"staffName" gorm:"column:staff_name"`   // 操作人
	CreateTime time.Time `json:"createTime" gorm:"column:create_time"`
	UpdateTime time.Time `json:"updateTime" gorm:"column:update_time"`
}

func (m *UserLevelPrivilegeItems) TableName() string {
	return "t_user_level_privilege_items"
}

// LevelPrivilegeItemDTO 等级权益物品详情
type LevelPrivilegeItemDTO struct {
	ID                    int    `json:"id" gorm:"column:id"`
	Level                 int    `json:"level" gorm:"column:level"`                                    // 等级id
	GoodsId               int    `json:"goodsId" gorm:"column:goods_id"`                               // 物品id
	LevelType             int    `json:"levelType" gorm:"column:level_type"`                           // 1:lv 2:vip 3:star
	ExpireDate            int    `json:"expireDate" gorm:"column:expire_date"`                         // 有效期(天)
	Explain               string `json:"explain" gorm:"column:explain"`                                // 说明
	GoodsName             string `json:"goodsName" gorm:"column:goods_name"`                           // 物品名称
	GoodsType             int    `json:"goodsType" gorm:"column:goods_type"`                           // 物品类型
	GoodsTypeName         string `json:"goodsTypeName" gorm:"column:goods_type_name"`                  // 物品类型名称
	GoodsTypeKey          string `json:"goodsTypeKey" gorm:"goods_type_key"`                           // 物品类型key
	GoodsIcon             string `json:"goodsIcon" gorm:"column:goods_icon"`                           // 图标
	GoodsAnimationUrl     string `json:"goodsAnimationUrl" gorm:"column:goods_animation_url"`          // 图片动效
	GoodsAnimationJsonUrl string `json:"goodsAnimationJsonUrl" gorm:"column:goods_animation_json_url"` // json文件动效
}
