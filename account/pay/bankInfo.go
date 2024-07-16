package pay

import (
	"time"
	"za.game/lib/dbconn"
	"za.game/lib/rds"

	"gorm.io/gorm"
)

type BankInfo struct {
	ID        uint64         `db:"id" gorm:"primarykey;column:id;type:int(11);";json:"id"`
	Uid       uint64         `db:"uid" gorm:"column:uid;type:bigint(20);not null;default:0;";json:"uid"`                        // 用户ID
	BankCode  string         `db:"bank_code" gorm:"column:bank_code;type:varchar(32);not null;default:'';";json:"bank_code"`    // 银行类型
	BankName  string         `db:"bank_name" gorm:"column:bank_name;type:varchar(32);not null;default:'';";json:"bank_name"`    // 银行名称
	AccountNo string         `db:"account_no" gorm:"column:account_no;type:varchar(32);not null;default:'';";json:"account_no"` // 银行账号
	Ifsc      string         `db:"ifsc" gorm:"column:ifsc;type:varchar(32);not null;default:'';";json:"ifsc"`                   // ifsc号
	Name      string         `db:"name" gorm:"column:name;type:varchar(32);not null;default:'';";json:"name"`                   //客户姓名
	Email     string         `db:"email" gorm:"column:email;type:varchar(32);not null;default:'';";json:"email"`                //客户邮箱
	Phone     string         `db:"phone" gorm:"column:phone;type:varchar(32);not null;default:'',";json:"phone"`                //客户手机
	Address   string         `db:"address" gorm:"column:address;type:varchar(64);not null;default:'';";json:"address"`          //客户地址
	Vpa       string         `db:"vpa" gorm:"column:vpa;type:varchar(32);not null;default:'';";json:"vpa"`                      //vpa
	Remark    string         `db:"remark" gorm:"column:remark;type:varchar(64);default:'';";json:"remark"`                      //备注
	CreatedAt time.Time      `db:"created_at" gorm:"column:created_at;";json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" gorm:"column:updated_at;";json:"updated_at"`
	DeletedAt gorm.DeletedAt `db:"deleted_at" gorm:"column:deleted_at;index;"`
}

func (table *BankInfo) TableName() (tableName string) {
	tableName = "bank_info"
	return
}

//	func BaseBankInfo() *gorm.DB {
//		return config.DB.Model(new(BankInfo))
//	}

// 获取账户信息
func GetBankUserInfo(uid uint64) (err error, bank BankInfo) {
	bank = BankInfo{}
	err = rds.SqlxGet(dbconn.PayDB, &bank,
		"SELECT id, uid, bank_code, account_no, ifsc, name, email, phone, address, vpa, remark FROM bank_info WHERE uid = ?",
		uid)
	if err != nil {
		return
	}
	return
}
