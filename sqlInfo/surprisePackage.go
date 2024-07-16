package sqlInfo

import "time"

type SurprisePackage struct {
	Uid      uint64    `db:"uid"  json:"uid"`           //用户ID
	StartDt  time.Time `db:"start_dt" json:"start_dt"`  //开始时间
	State    int       `db:"state" json:"state"`        //状态 0未开始 1已生成 2已完成
	PackData string    `db:"packData"  json:"packData"` //礼包数据
	EndDt    time.Time `db:"end_dt"  json:"end_dt"`     //结束时间
	NextDt   time.Time `db:"next_dt"  json:"next_dt"`   //下一次刷新时间
}

func (table *SurprisePackage) TableName() (tableName string) {
	tableName = "surprise_pack_userinfo"
	return
}

type SurpriseITem struct {
	GiftId    int `json:"giftid"`     //
	Pay       int `json:"pay"`        //
	Cash      int `json:"cash"`       //
	ExtraCash int `json:"extra_cash"` //
	Bonus     int `json:"bonus"`      //
	Total     int `json:"total"`      //
}
