package pay

import (
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	"za.game/lib/tool"
)

var (
	Rand *rand.Rand
)

func init() {
	// 创建自定义的随机数生成器源
	source := rand.NewSource(time.Now().UnixNano())

	// 创建自定义的随机数生成器实例
	Rand = rand.New(source)
}

// 生成订单号
func GetOrderNo() (orderNo string) {
	now := time.Now().Local()
	ymdHis := now.Format("20060102150405")

	num := Rand.Intn(1000) + 1000
	orderNo = tool.StringJoin([]string{"NO", ymdHis, strconv.Itoa(num)})
	return orderNo
}

type Order struct {
	Id           int          `db:"id" gorm:"primarykey;column:id;type:int(11);" json:"id" `                         //
	Uid          int64        `db:"uid" gorm:"column:uid;type:bigint(20);" json:"uid" `                              // 用户ID
	Ymd          int          `db:"ymd" gorm:"column:ymd;type:int(11);" json:"ymd" `                                 // 年月日
	OrderNo      string       `db:"order_no" gorm:"column:order_no;type:varchar(32);" json:"order_no" `              // 订单ID
	MOrderNo     string       `db:"m_order_no" gorm:"column:m_order_no;type:varchar(32);" json:"m_order_no" `        // 商户订单号
	PayId        int          `db:"pay_id" gorm:"column:pay_id;type:int(11);" json:"pay_id" `                        // 支付配置ID
	Account      int          `db:"account" gorm:"column:account;type:int(11);" json:"account" `                     // 支付钱
	Cash         int          `db:"cash" gorm:"column:cash;type:int(11);" json:"cash" `                              // cash
	GiftCash     int          `db:"gift_cash" gorm:"column:gift_cash;type:int(11);" json:"gift_cash" `               // 赠送的cash
	Bonus        int          `db:"bonus" gorm:"column:bonus;type:int(11);" json:"bonus" `                           // bonus
	RequestTime  sql.NullTime `db:"request_time" gorm:"column:request_time;type:timestamp;" json:"request_time" `    // 下单时间
	Email        string       `db:"email" gorm:"column:email;type:varchar(32);" json:"email" `                       // 邮件
	Name         string       `db:"name" gorm:"column:name;type:varchar(32);" json:"name" `                          // 姓名
	Phone        string       `db:"phone" gorm:"column:phone;type:varchar(32);" json:"phone" `                       // 手机
	RedirectTime sql.NullTime `db:"redirect_time" gorm:"column:redirect_time;type:timestamp;" json:"redirect_time" ` // 下单返回时间
	Status       int8         `db:"status" gorm:"column:status;type:tinyint(4);" json:"status" `                     // 状态(0=等待支付,1=完成支付,2=下单失败,3=已发货)
	CompleteTime int          `db:"complete_time" gorm:"column:complete_time;type:int(11);" json:"complete_time" `   // 订单完成时间
	H5Url        string       `db:"h5_url" gorm:"column:h5_url;type:varchar(128);" json:"h5_url" `                   // 支付地址
	Remark       string       `db:"remark" gorm:"column:remark;type:varchar(64);" json:"remark" `                    // 备注
	CreatedAt    sql.NullTime `db:"created_at" gorm:"column:created_at;type:timestamp;" json:"created_at" `          //
	UpdatedAt    sql.NullTime `db:"updated_at" gorm:"column:updated_at;type:timestamp;" json:"updated_at" `          //
	Type         int8         `db:"type" gorm:"column:type;type:tinyint(4);" json:"type" `                           // 类型 1:充值;2:礼包
	GiftId       int          `db:"gift_id" gorm:"column:gift_id;type:int(11);" json:"gift_id" `                     // 充值项
	RoomId       string       `db:"room_id" gorm:"column:room_id;type:varchar(255);" json:"room_id" `                // 房间id，房间内充值数值记录
	Channel      int          `db:"channel" gorm:"column:channel;type:int(11);not null;default:1;" json:"channel"`   //渠道
}

func (t *Order) TableName() string {
	return "order_" + time.Now().Format("200601")
}

func (t *Order) GetTableName(mTime time.Time) string {
	return "order_" + mTime.Format("200601")
}
