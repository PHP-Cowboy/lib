package sqlInfo

import (
	"encoding/json"
)

type SpecialDataItem struct {
	PoolId int    `json:"id"`
	YMD    string ` json:"ymd"`   //日期
	Count  int    ` json:"count"` //当前触发次数
}
type UserActivity struct {
	ID                   uint64 `db:"id" gorm:"primarykey;column:id;type:int(11);";json:"id"`
	Uid                  uint64 `db:"uid" gorm:"column:uid;type:bigint(20);not null;default:0;";json:"uid"` //用户ID
	UserTag              int    `db:"user_tag" json:"user_tag"`                                             //用户标识
	TpNewControl         int    `db:"tp_new" json:"tp_new"`                                                 //tp新手次数 调控
	TyroCash             int    `db:"tyro_cash" json:"tyro_cash"`                                           //新手赠送cash
	SlotSaveTime         int    `db:"slot_save_time" json:"slot_save_time"`                                 //slot拉回次数
	DynamicRechargeValue int64  `db:"dynamic_recharge_value" json:"dynamic_recharge_value"`                 //充值动态值（数据库）
	FirstSpecialPoker    int    `db:"FirstSpecialPoker" json:"FirstSpecialPoker"`                           //新手ss
	SpecialExtData       string `db:"SpecialExtData"  json:"SpecialExtData,omitempty"`                      //强体验

	SpecialData map[int]*SpecialDataItem
}

func (ua *UserActivity) ParseSpecialExtData() error {
	var items []*SpecialDataItem
	ua.SpecialData = make(map[int]*SpecialDataItem, len(items))
	err := json.Unmarshal([]byte(ua.SpecialExtData), &items)
	if err != nil {
		return err
	}
	for _, item := range items {
		ua.SpecialData[item.PoolId] = item
	}

	return nil
}
func (ua *UserActivity) ParseSpecialData() (string, error) {
	items := make([]SpecialDataItem, 0, len(ua.SpecialData))
	for _, item := range ua.SpecialData {
		items = append(items, *item)
	}
	jsonBytes, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
func (table *UserActivity) TableName() (tableName string) {
	tableName = "user_activity"
	return
}
