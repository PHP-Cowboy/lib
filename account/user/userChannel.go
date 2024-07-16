package user

import (
	"database/sql"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
)

type UserChannel struct {
	Id          int          `db:"id" gorm:"primaryKey;column:id;type:int(11);not null;" json:"id" `
	ChannelName string       `db:"channel_name" gorm:"column:channel_name;type:varchar(32);not null;default:'';" json:"channel_name" `
	Code        string       `db:"code" gorm:"column:code;type:varchar(255);not null;default:'';unique:code;" json:"code" `
	Remark      string       `db:"remark" gorm:"column:remark;type:varchar(32);not null;default:'';" json:"remark" `
	CreatedAt   sql.NullTime `db:"created_at" gorm:"column:created_at;type:timestamp;" json:"created_at" `
	UpdatedAt   sql.NullTime `db:"updated_at" gorm:"column:updated_at;type:timestamp;" json:"updated_at" `
	DeletedAt   sql.NullTime `db:"deleted_at" gorm:"column:deleted_at;type:timestamp;" json:"deleted_at" `
}

// 获取账户信息
func GetUserChannelList() (dataList []UserChannel, err error) {

	err = rds.SqlxSelect(
		dbconn.NDB,
		&dataList,
		"SELECT id, channel_name, code FROM user_channel",
	)

	return
}
