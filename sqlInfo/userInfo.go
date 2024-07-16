package sqlInfo

import (
	"database/sql"
	"fmt"
	"time"
)

type UserInfo struct {
	ID               uint64       `db:"id" gorm:"primarykey;column:id;type:int(11);not null;" json:"id" `
	Uid              uint64       `db:"uid" gorm:"column:uid;type:bigint(20);not null;default:0;unique:uid_unique;comment:用户ID;" json:"uid" `
	Channel          int8         `db:"channel" gorm:"column:channel;type:tinyint(4);not null;default:1;comment:渠道;" json:"channel" `
	Sex              int8         `db:"sex" gorm:"column:sex;type:tinyint(4);not null;default:1;comment:性别（1=男，2=女）;" json:"sex" `
	Age              int8         `db:"age" gorm:"column:age;type:tinyint(4);not null;default:0;comment:年龄;" json:"age" `
	Vip              int64        `db:"vip" gorm:"column:vip;type:bigint(4);not null;default:0;index:vip_key;comment:vip经验;" json:"vip" `
	WinCash          int64        `db:"win_cash" gorm:"column:win_cash;type:bigint(20);not null;default:0;comment:可用金额（提现下注）;" json:"win_cash" `
	Cash             int64        `db:"cash" gorm:"column:cash;type:bigint(20);not null;default:0;comment:半冻结金额（不可提现可下注）;" json:"cash" `
	Bonus            int64        `db:"bonus" gorm:"column:bonus;type:bigint(20);not null;default:0;comment:冻结金额（不可提现下注）;" json:"bonus" `
	CanExchangeBonus int64        `db:"can_exchange_bonus" gorm:"column:can_exchange_bonus;type:bigint(20);default:0;comment:可兑换bonus数量;" json:"can_exchange_bonus" `
	Recharge         int64        `db:"recharge" gorm:"column:recharge;type:bigint(20);not null;default:0;index:recharge_key;comment:充值总金额;" json:"recharge" `
	RechargeCount    int          `db:"recharge_count" gorm:"column:recharge_count;type:int(11);not null;default:0;comment:充值总次数;" json:"recharge_count" `
	Bet              int64        `db:"bet" gorm:"column:bet;type:bigint(20);not null;default:0;comment:总投注;" json:"bet" `
	Transport        int64        `db:"transport" gorm:"column:transport;type:bigint(20);not null;default:0;comment:总输;" json:"transport" `
	Win              int64        `db:"win" gorm:"column:win;type:bigint(20);not null;default:0;comment:总赢;" json:"win" `
	Tax              int64        `db:"tax" gorm:"column:tax;type:bigint(20);not null;default:0;comment:税;" json:"tax" `
	Remark           string       `db:"remark" gorm:"column:remark;type:varchar(32);not null;default:'';comment:备注;" json:"remark" `
	CreatedAt        time.Time    `db:"created_at" gorm:"column:created_at;type:timestamp;" json:"created_at" `
	UpdatedAt        time.Time    `db:"updated_at" gorm:"column:updated_at;type:timestamp;" json:"updated_at" `
	DeletedAt        sql.NullTime `db:"deleted_at" gorm:"column:deleted_at;type:timestamp;" json:"deleted_at" `
	VipCashBank1     int          `db:"vip_cash_bank1" gorm:"column:vip_cash_bank1;type:int(20);not null;default:0;comment:可以领取的vip存钱罐奖励;" json:"vip_cash_bank1" `
	VipCashBank2     int          `db:"vip_cash_bank2" gorm:"column:vip_cash_bank2;type:int(20);not null;default:0;comment:当天累计的vip存钱罐奖励;" json:"vip_cash_bank2" `
	VipCashBankTime  time.Time    `db:"vip_cash_bank_time" gorm:"column:vip_cash_bank_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:上一次累计的时间;" json:"vip_cash_bank_time" `
	WithdrawMoney    int          `db:"withdraw_money" gorm:"column:withdraw_money;type:int(11);default:0;comment:体现审核中;" json:"withdraw_money" `
	WithdrawedMoney  int          `db:"withdrawed_money" gorm:"column:withdrawed_money;type:int(11);default:0;comment:体现到账;" json:"withdrawed_money" `
	ExtraGift        int64        `db:"extra_gift" gorm:"column:extra_gift;type:bigint(20);default:0;comment:额外福利;" json:"extra_gift" `
	//PoolTax          int64        `db:"pool_tax" json:"pool_tax"`
}

func GetUserInfoTableName(uid uint64) (tableName string) {
	tableName = fmt.Sprintf("user_info_%02d", uid%5)
	return
}

// redis 不重要数据，redis保存
type UserInRoomInfo struct {
	CurRechargeCount    int   `json:"CurRechargeCount"`    //当前房间充值次数
	DailyControlWinCash int64 `json:"DailyControlWinCash"` //日控赢钱
}

// vip配置信息
type VipConfig struct {
	CashBankRate int //vip存储罐比例
	CashBankTop  int //vip存储罐上限
}

/**分享裂变信息*/
type FxlbInfo struct {
	FxlbYjWei   int64   `db:"fxlb_yj_wei"` //分享裂变给邀请用户提供的可转换的佣金
	FxlbYjZhuan int64   //分享裂变转给邀请者可获得佣金
	Yj_zh_rate  float64 //分享裂变佣金转化比例
	FromUid     uint64  //邀请者的id
}

type GameUserSqlItem struct {
	User
	UserInfo
	UserActivity
	UserInRoomInfo
	VipConfig
	FxlbInfo
}
