package sqlInfo

import (
	"database/sql"
	"time"
)

type RoomFundsFlowLog struct {
	Id           int          `db:"id" gorm:"primarykey;column:id;type:int(11) unsigned;not null;comment:id;" json:"id" `
	CreatedAt    sql.NullTime `db:"created_at" gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间;" json:"created_at" `
	Uid          int          `db:"uid" gorm:"column:uid;type:int(11);not null;comment:用户ID;" json:"uid" `
	GameId       int          `db:"game_id" gorm:"column:game_id;type:int(11);not null;comment:游戏id;" json:"game_id" `
	GameIdX      int64        `db:"game_id_x" gorm:"column:game_id_x;type:bigint(20);not null;default:0;comment:对局标识;" json:"game_id_x" `
	RoomId       int          `db:"room_id" gorm:"column:room_id;type:int(11);not null;comment:房间id;" json:"room_id" `
	DeskId       int          `db:"desk_id" gorm:"column:desk_id;type:int(11);not null;comment:桌子id;" json:"desk_id" `
	Pos          int8         `db:"pos" gorm:"column:pos;type:tinyint(4);not null;comment:座位号;" json:"pos" `
	Cash         int64        `db:"cash" gorm:"column:cash;type:bigint(20);not null;default:0;comment:当前cash余额;" json:"cash" `
	CashNums     int64        `db:"cash_nums" gorm:"column:cash_nums;type:bigint(20);not null;default:0;comment:cash变动数额;" json:"cash_nums" `
	WinCash      int64        `db:"win_cash" gorm:"column:win_cash;type:bigint(20);not null;default:0;comment:充值金币账户;" json:"win_cash" `
	WinCashNums  int64        `db:"win_cash_nums" gorm:"column:win_cash_nums;type:bigint(20);not null;default:0;comment:winCash变动数额;" json:"win_cash_nums" `
	Withdraw     int64        `db:"withdraw" gorm:"column:withdraw;type:bigint(20);not null;default:0;comment:赠送储钱罐提现额度;" json:"withdraw" `
	WithdrawNums int64        `db:"withdraw_nums" gorm:"column:withdraw_nums;type:bigint(20);not null;default:0;comment:变动的储钱罐提现额度;" json:"withdraw_nums" `
	ExtraGift    int64        `db:"extra_gift" gorm:"column:extra_gift;type:bigint(20);not null;default:0;comment:额外赠送;" json:"extra_gift" `
	Tax          int64        `db:"tax" gorm:"column:tax;type:bigint(20);not null;default:0;comment:税收;" json:"tax" `
	Remark       string       `db:"remark" gorm:"column:remark;type:varchar(255);comment:备注;" json:"remark" `
	Bonus        int64        `db:"bonus" gorm:"column:bonus;type:bigint(20);not null;default:0;comment:当前bonus余额;" json:"bonus" `
	BonusNums    int64        `db:"bonus_nums" gorm:"column:bonus_nums;type:bigint(20);not null;default:0;" json:"bonus_nums" `
}

func GetGameFundsFlowLogTable() string {
	//todo 检验表是否存在，不存在则创建
	return "room_funds_flow_log_" + time.Now().Format("200601")
}

type GameRecordLog struct {
	Id        int          `db:"id" gorm:"primarykey;column:id;type:int(11) unsigned;not null;comment:id;" json:"id" `
	GameIdX   int64        `db:"game_idx" gorm:"column:game_idx;type:bigint(20);not null;default:0;index:game_idx;comment:对局标识;" json:"game_idx" `
	CreatedAt sql.NullTime `db:"created_at" gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:创建时间;" json:"created_at" `
	GameId    int          `db:"game_id" gorm:"column:game_id;type:int(11);not null;default:0;comment:游戏id;" json:"game_id" `
	RoomId    int          `db:"room_id" gorm:"column:room_id;type:int(11);not null;default:0;comment:房间id;" json:"room_id" `
	DeskId    int          `db:"desk_id" gorm:"column:desk_id;type:int(11);not null;default:0;comment:桌子id;" json:"desk_id" `
	Details   string       `db:"details" gorm:"column:details;type:varchar(8000);not null;comment:详情;" json:"details" `
}

func GetGameRecordLogTable() string {
	//todo 检验表是否存在，不存在则创建
	return "GameRecord_log_" + time.Now().Format("200601")
}
