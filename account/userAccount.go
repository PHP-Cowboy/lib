package account

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type UserAccount struct {
	ID        uint64    `gorm:"primarykey;column:id;type:int(11);" json:"id" db:"id"`
	Uid       uint64    `gorm:"column:uid;type:bigint(20);not null;default:0;index:uid_key;" json:"uid" db:"uid"`                    //用户ID
	Ymd       uint32    `gorm:"column:ymd;type:int(11);not null;default:0;index:uid_key;index:ymd_key;" json:"ymd" db:"ymd"`         //年月日
	Uuid      string    `gorm:"column:uuid;type:varchar(32);not null;default:'';" json:"uuid" db:"uuid"`                             //局标识
	WinCash   int       `gorm:"column: win_cash;type:int(11);not null;default:0;" json:" win_cash" db:" win_cash"`                   //可以用金额
	Cash      int       `gorm:"column:cash;type:int(11);not null;default:0;" json:"cash" db:"cash"`                                  //下注金额
	Bonus     int       `gorm:"column:bonus;type:int(11);not null;default:0;" json:"bonus" db:"bonus"`                               //半冻结解金额
	Option    int       `gorm:"column:option;type:int(11);not null;default:0;" json:"option" db:"option"`                            //操作金额
	Tax       int       `gorm:"column:tax;type:int(11);not null;default:0;" json:"tax" db:"tax"`                                     //税
	Type      uint8     `gorm:"column:type;type:tinyint(4);not null;default:0;index:game_key;index:type_key;" json:"type" db:"type"` //1=每日任务，2=VIP,3=邮件,4=签到,5=充值,6=结算,7=储钱罐,8=救济金,9=礼包,10=注册赠送
	GameId    uint32    `gorm:"column:game_id;type:int(11);not null;default:0;index:game_key;" json:"game_id" db:"game_id"`          //游戏ID
	Room      uint32    `gorm:"column:room;type:int(11);not null;default:0;" json:"room" db:"room"`                                  //房间号
	Desk      uint32    `gorm:"column:desk;type:int(11);not null;default:0;" json:"desk" db:"desk"`                                  //桌子
	CreatedAt time.Time `gorm:"column:created_at;" json:"created_at" db:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;" json:"updated_at" db:"updated_at"`
}

func (t *UserAccount) TableName(uid uint64) string {
	return fmt.Sprintf("user_account_0%v", uid%10)
}

func (t *UserAccount) SaveTx(tx *sqlx.Tx, fields, values string, mp map[string]interface{}, uid uint64) (err error) {
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", t.TableName(uid), fields, values)
	_, err = tx.NamedExec(sql, mp)
	return
}
