package sqlInfo

import (
	"database/sql"
	"time"
)

type ItemsLogDB struct {
	Id        int          `db:"id" gorm:"primarykey;column:id;type:int(11) unsigned;not null;" json:"id" `
	CreatedAt sql.NullTime `db:"created_at" gorm:"column:created_at;type:datetime;not null;" json:"created_at" `
	Uid       int          `db:"uid" gorm:"column:uid;type:int(11);not null;" json:"uid" `
	Type      int8         `db:"type" gorm:"column:type;type:tinyint(4);not null;default:0;" json:"type" `
	Itemid    int64        `db:"itemid" gorm:"column:itemid;type:bigint(11);not null;" json:"itemid" `
	Addnum    int64        `db:"addnum" gorm:"column:addnum;type:bigint(11);not null;" json:"addnum" `
	Leftnum   int64        `db:"leftnum" gorm:"column:leftnum;type:bigint(11);not null;" json:"leftnum" `
	Remark    string       `db:"remark" gorm:"column:remark;type:varchar(64);not null;" json:"remark" `
	OrderNo   string       `db:"order_no" gorm:"column:order_no;type:varchar(32);" json:"order_no" `
	MOrderNo  string       `db:"m_order_no" gorm:"column:m_order_no;type:varchar(32);" json:"m_order_no" `
	Roomid    int64        `db:"roomid" gorm:"column:roomid;type:bigint(20);" json:"roomid" `
}

func (table *ItemsLogDB) TableName() (tableName string) {
	tableYm := time.Now().Local().Format("200601")
	tableName = "items_log_" + tableYm
	return
}
