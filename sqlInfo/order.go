package sqlInfo

import (
	"time"

	"za.game/lib/rds"

	"github.com/jmoiron/sqlx"
)

type Order struct {
	ID           uint64    `db:"id" insert_db:"-" `
	Uid          uint64    `db:"uid"`           // 用户ID
	Ymd          int       `db:"ymd"`           // 年月日
	OrderNo      string    `db:"order_no"`      // 订单ID
	MOrderNo     string    `db:"m_order_no"`    // 商户订单ID
	PayId        int       `db:"pay_id"`        // 支付配置id
	Account      int       `db:"account"`       // 支付金额
	Cash         int       `db:"cash"`          //cash金额
	GiftCash     int       `db:"gift_cash"`     //赠送的cash
	Bonus        int       `db:"bonus"`         //储钱罐
	RequestTime  time.Time `db:"request_time"`  //下单时间
	Email        string    `db:"email"`         // 邮箱地址
	Name         string    `db:"name"`          // 用户名
	Phone        string    `db:"phone"`         //手机
	RedirectTime time.Time `db:"redirect_time"` //下单拿h5地址时间
	Status       int8      `db:"status"`        // 状态 0=等待支付，1=支付完成，2=下单失败
	CompleteTime int       `db:"complete_time"` //订单完成时间
	H5Url        string    `db:"h5_url"`        // h5地址
	Type         int       `db:"type"`          // 类型 1:充值;2:礼包
	GiftId       int       `db:"gift_id"`       // 类型 充值项
	Remark       string    `db:"remark"`        //备注
	RoomId       string    `db:"room_id"`       //房间id，房间内充值数值记录
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	Channel      int       `db:"channel"` //渠道
}

func (table *Order) TableName() (tableName string) {
	tableYm := time.Now().Local().Format("200601")
	tableName = "order_" + tableYm
	return
}

/**查询充值礼包的购买数量*/
func GetBuyNumRechargeGift(db *sqlx.DB, uid int64, giftid int) (int, error) {
	// order := Order{}

	// buy_num := 0
	// rows, err := rds.Query(db, "select count(*)'num' from %s where uid = %d and gift_id=%d and (status = 3 or status = 1)",
	// 	order.TableName(), uid, giftid)
	// if err != nil {
	// 	return buy_num, err
	// }
	// defer rows.Close()

	// if rows.Next() {
	// 	if err := rows.Scan(&buy_num); err != nil {
	// 		return buy_num, err
	// 	}
	// }

	buy_num := 0
	rows, err := rds.Query(db, "select buynum from recharge_gift_user where uid=%v and giftid=%v", uid, giftid)
	if err != nil {
		return buy_num, err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&buy_num); err != nil {
			return buy_num, err
		}
	}

	return buy_num, nil
}

/**查询二选一的购买数量（一天购买一次）
*@param db: 数据库连接
*@param uid: 用户id
*@param giftidmin: 属于该礼包的最小值
*@param giftidmax: 属于该礼包的最大值*/
func GetBuyNumEventGift(db *sqlx.DB, uid int64, giftidmin int, giftidmax int) (int, error) {
	order := Order{}

	buy_num := 0
	rows, err := rds.Query(db, "select count(*)'num' from %s where uid = %d and gift_id > %d and gift_id <= %d and date(updated_at)=CURDATE() and (status = 3 or status = 1)",
		order.TableName(), uid, giftidmin, giftidmax)
	if err != nil {
		return buy_num, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&buy_num); err != nil {
			return buy_num, err
		}
	}

	return buy_num, nil
}

/**查询OnlyOne三选一的购买数量
*@param db: 数据库连接
*@param uid: 用户id
*@param giftidmin: 属于该礼包的最小值
*@param giftidmax: 属于该礼包的最大值*/
func GetBuyNumOnlyOneGift(db *sqlx.DB, uid int64, giftidmin int, giftidmax int) (int, error) {
	order := Order{}

	buy_num := 0
	rows, err := rds.Query(db, "select count(*)'num' from %s where uid = %d and gift_id > %d and gift_id <= %d and (status = 3 or status = 1)",
		order.TableName(), uid, giftidmin, giftidmax)
	if err != nil {
		return buy_num, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&buy_num); err != nil {
			return buy_num, err
		}
	}

	return buy_num, nil
}

/**查询救济礼包购买数量（一天购买一次）
*@param db: 数据库连接
*@param uid: 用户id
*@param giftidmin: 属于该礼包的最小值
*@param giftidmax: 属于该礼包的最大值*/
func GetBuyNumBenefitGift(db *sqlx.DB, uid int64, giftidmin int, giftidmax int) (int, error) {
	order := Order{}

	buy_num := 0
	rows, err := rds.Query(db, "select count(*)'num' from %s where uid = %d and gift_id > %d and gift_id <= %d and date(updated_at)=CURDATE() and (status = 3 or status = 1)",
		order.TableName(), uid, giftidmin, giftidmax)
	if err != nil {
		return buy_num, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&buy_num); err != nil {
			return buy_num, err
		}
	}

	return buy_num, nil
}
