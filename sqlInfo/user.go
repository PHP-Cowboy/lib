package sqlInfo

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         uint64         `db:"id" gorm:"primarykey;column:id;type:int(11);";json:"id"`
	Uid        uint64         `db:"uid" gorm:"column:uid;type:bigint(20);not null;default:0;index:unique;";json:"uid"`        //用户ID
	IsGust     uint8          `db:"is_guest" gorm:"column:is_guest;type:tinyint(4);not null;default:0;";json:"is_guest"`      //是否游客
	IsSend     uint8          `db:"is_send" gorm:"column:is_send;type:tinyint(4);not null;default:0;";json:"is_send"`         //是否赠送
	Device     string         `db:"device" gorm:"column:device;type:varchar(64);not null;default:'';";json:"device"`          //设备码
	UserName   string         `db:"user_name" gorm:"column:user_name;type:varchar(32);not null;default:'';";json:"user_name"` //用户名
	Icon       int8           `db:"icon" gorm:"column:icon;type:tinyint(4);not null;default:0;";json:"icon"`                  //头像
	Phone      string         `db:"phone" gorm:"column:phone;type:varchar(16);not null;default:'';";json:"phone"`             //电话
	Email      string         `db:"email" gorm:"column:email;type:varchar(32);not null;default:''1;";json:"email"`            //邮箱
	Pwd        string         `db:"pwd" gorm:"column:pwd;type:varchar(32);not null;default:'';";json:"pwd"`                   //密码
	Token      string         `db:"token" gorm:"column:token;type:varchar(512);not null;default:'';";json:"token"`            //token
	ChannelId  int            `gorm:"column:channel_id;type:int(11);not null;default:0;" json:"channel_id" db:"channel_id"`   //用户ID
	TpNew      int            `json:"tp_new" db:"tp_new"`
	Expire     time.Time      `db:"expire" gorm:"column:expire;" json:"expire"`
	RegIp      string         `json:"reg_ip" db:"reg_ip"`           //注册ip
	RegVersion string         `json:"reg_version" db:"reg_version"` //注册版本号
	Gpcadid    string         `json:"gpcadid" db:"gpcadid"`         //gpcadid
	CreatedAt  time.Time      `db:"created_at" gorm:"column:created_at;";json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" gorm:"column:updated_at;";json:"updated_at"`
	DeletedAt  gorm.DeletedAt `db:"deleted_at" gorm:"column:deleted_at;index;"`
}

func (table *User) TableName() (tableName string) {
	tableName = "user"
	return
}
