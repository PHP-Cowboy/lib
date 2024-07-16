package sqlInfo

import (
	"time"
)

/**分享裂变，流水统计*/
type FxlbStatLog struct {
	Id   int   `db:"id" gorm:"primarykey;column:id;type:int(11);not null;" json:"id" `
	Day  int   `db:"day" gorm:"column:day;type:int(11);not null;default:0;comment:日期：20240701;" json:"day" `
	Uid  int64 `db:"uid" gorm:"column:uid;type:bigint(20);not null;comment:用户id;" json:"uid" `
	YqYj int64 `db:"yq_yj" gorm:"column:yq_yj;type:bigint(20);not null;default:0;comment:邀请佣金;" json:"yq_yj" `
	CzYj int64 `db:"cz_yj" gorm:"column:cz_yj;type:bigint(20);not null;default:0;comment:充值佣金;" json:"cz_yj" `
	DlYj int64 `db:"dl_yj" gorm:"column:dl_yj;type:bigint(20);not null;default:0;comment:代理佣金;" json:"dl_yj" `
}

func GetFxlbStatTable(t int64) string {
	//todo 检验表是否存在，不存在则创建
	if t > 0 {
		return "fxlb_stat_log_" + time.Unix(t, 0).Format("200601")
	} else {
		return "fxlb_stat_log_" + time.Now().Format("200601")
	}
}
