package account

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type StoreBankLog struct {
	ID        uint64    `gorm:"primarykey;column:id;type:int(11);" json:"id" db:"id"`
	Uuid      string    `gorm:"column:uuid;type:varchar(32);not null;default:'';unique:uuid_key;" json:"uuid" db:"uuid"` //局标识
	Uid       uint64    `gorm:"column:uid;type:bigint(20);not null;default:0;index:uid_key;" json:"uid" db:"uid"`        // 用户ID
	Ymd       uint32    `gorm:"column:ymd;type:int(11);not null;default:0;index:uid_key;" json:"ymd" db:"ymd"`           // 年月日
	Total     int       `gorm:"column:total;type:int(11);not null;default:0;" json:"total" db:"total"`                   //总计
	Draw      int       `gorm:"column:draw;type:int(11);not null;default:0;" json:"draw" db:"draw"`                      // 兑换额度
	Amount    int       `gorm:"column:amount;type:int(11);not null;default:0;" json:"amount" db:"amount"`                // 操作金额
	Type      uint8     `gorm:"column:type;type:tinyint(4);not null;default:0;" json:"type" db:"type"`                   // 类别 (1=每日任务，2=VIP,3=邮件,4=签到,5=充值,6=结算,7=储钱罐,8=救济金,9=礼包)
	Category  uint8     `gorm:"column:category;type:tinyint(4);not null;default:0;" json:"category" db:"category"`       // 类型(1=兑换,2=获取)
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at" db:"created_at"`
}

func (t *StoreBankLog) TableName() string {
	return "store_bank_log_" + time.Now().Local().Format("200601")
}

func (t *StoreBankLog) SaveTx(tx *sqlx.Tx, fields, values string, mp map[string]interface{}) (err error) {
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", t.TableName(), fields, values)

	_, err = tx.NamedExec(sql, mp)

	return
}
