package account

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type StoreBankUser struct {
	ID         uint64    `gorm:"primarykey;column:id;type:int(11);" json:"id" db:"id"`
	Uid        uint64    `gorm:"column:uid;type:bigint(20);not null;default:0;" json:"uid" db:"uid"`            //用户ID
	Bank       int       `gorm:"column:bank;type:int(11);not null;default:0;" json:"bank" db:"bank"`            //储钱罐金额
	Withdrawal int       `gorm:"column:withdrawal;type:int(11);default:0;" json:"withdrawal" db:"withdrawal"`   //可提现金额
	Remark     string    `gorm:"column:remark;type:varchar(32);not null;default:'';" json:"remark" db:"remark"` //备注
	CreatedAt  time.Time `gorm:"column:created_at;" json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;" json:"updated_at" db:"updated_at"`
}

func (t *StoreBankUser) TableName(uid uint64) string {
	return fmt.Sprintf("store_bank_user_%02d", uid%5)
}

func (t *StoreBankUser) GetOneByUid(tx *sqlx.Tx, fields string, uid uint64) (store StoreBankUser, err error) {
	sql := fmt.Sprintf("select %s from %s where uid = ? limit 1", fields, t.TableName(uid))

	err = tx.Get(&store, sql, uid)
	return
}

func (t *StoreBankUser) Save(tx *sqlx.Tx, fields, values string, mp map[string]interface{}, uid uint64) (err error) {
	sql := fmt.Sprintf("insert into %s(%s) values (%s)", t.TableName(uid), fields, values)

	_, err = tx.NamedExec(sql, mp)

	return
}

func (t *StoreBankUser) UpdateTxBankByUid(tx *sqlx.Tx, bank int, uid uint64) (err error) {
	sql := fmt.Sprintf("update %s set bank = ? where uid = ?", t.TableName(uid))

	_, err = tx.Exec(sql, bank, uid)

	return err
}

func (t *StoreBankUser) UpdateTxWithdrawalByUid(tx *sqlx.Tx, withdrawal int, uid uint64) (err error) {
	sql := fmt.Sprintf("update %s set withdrawal = ? where uid = ?", t.TableName(uid))

	_, err = tx.Exec(sql, withdrawal, uid)

	return err
}

func (t *StoreBankUser) UpdateTxBankWithdrawalByUid(tx *sqlx.Tx, bank, withdrawal int, uid uint64) (err error) {
	sql := fmt.Sprintf("update %s set bank = ?, withdrawal = ? where uid = ?", t.TableName(uid))

	_, err = tx.Exec(sql, bank, withdrawal, uid)

	return err
}
