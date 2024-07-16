package sqlInfo

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"za.game/lib/rds"
)

type UserCharge struct {
	Uid        int64        `db:"uid" gorm:"primarykey;column:uid;type:bigint(20);not null;index:user_charge_idx2;" json:"uid" `
	Recharge   int          `db:"recharge" gorm:"primarykey;column:recharge;type:int(11);not null;unique:user_charge_idx1;" json:"recharge" `
	Num        int          `db:"num" gorm:"column:num;type:int(11);default:0;" json:"num" `
	UpdateTime sql.NullTime `db:"update_time" gorm:"column:update_time;type:datetime;not null;default:CURRENT_TIMESTAMP;" json:"update_time" `
}

func (table *UserCharge) TableName() (tableName string) {
	tableName = "user_charge"
	return
}
func (table *UserCharge) GetUserChargeCount(db *sqlx.DB, uid uint64) (int, error) {

	count := 0
	rows, err := rds.Query(db, "SELECT count(uid)'num' from %s where uid = %d", table.TableName(), uid)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}

	return count, nil
}
