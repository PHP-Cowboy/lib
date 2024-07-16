package account

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserPrize struct {
	ID        uint64    `gorm:"primarykey;column:id;type:int(11);" json:"id" db:"id"`
	Uuid      string    `gorm:"column:uuid;type:varchar(32);not null;default:0;unique:uuid_key" json:"uuid" db:"uuid"`      // 唯一标识
	Uid       uint64    `gorm:"column:uid;type:bigint(20);not null;default:0;index:uid_key;" json:"uid" db:"uid"`           // 用户ID
	GameId    uint32    `gorm:"column:game_id;type:int(11);not null;default:0;index:game_key;" json:"game_id" db:"game_id"` // 游戏ID
	Ymd       uint32    `gorm:"column:ymd;type:int(11);not null;default:0;index:uid_ymd;" json:"ymd" db:"ymd"`              // 年月日
	Type      uint8     `gorm:"column:type;type:tinyint(4);not null;default:1;index:type_key;" json:"type" db:"type"`       //类别 (1=每日任务，2=VIP,3=邮件,4=签到,5=充值,6=结算,7=储钱罐,8=救济金,9=礼包)
	WinCash   int       `gorm:"column:win_cash;type:int(11);not null;default:0;" json:"win_cash" db:"win_cash"`             // 无限制金额
	Cash      int       `gorm:"column:cash;type:int(11);not null;default:0;" json:"cash" db:"cash"`                         // cash金额
	Bonus     int       `gorm:"column:bonus;type:int(11);not null;default:0;" json:"store" db:"bonus"`                      // 储钱罐金额
	Category  uint8     `gorm:"column:category;type:tinyint(4);not null;default:0;" json:"category" db:"category"`          // 金额操作类型（1=加，2=减）
	Tax       int       `gorm:"column:tax;type:int(11);not null;default:0;" json:"tax" db:"tax"`                            // 税
	MsgTime   time.Time `gorm:"column:msg_time;" json:"msg_time" db:"msg_time"`
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at" db:"created_at"`
}

func (t *UserPrize) TableName() string {
	return "user_prize_" + time.Now().Local().Format("200601")
}

func (t *UserPrize) SaveTx(tx *sqlx.Tx, fields, values string, mp map[string]interface{}) (err error) {
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", t.TableName(), fields, values)

	_, err = tx.NamedExec(sql, mp)
	return
}
