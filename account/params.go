package account

import "time"

type Account struct {
	Uid        uint64
	UUID       string
	Cate       uint8 //类别 (1=每日任务，2=VIP,3=邮件,4=签到,5=充值,6=结算,7=储钱罐,8=救济金,9=礼包)
	OptionType uint8 //1=加操作，2=减操作
	Account    int
	Money      int64
	Cash       int64
	Store      int64
	Recharge   int64
	GameId     uint32
	Ymd        uint32
	Tax        int
	IsBank     uint8 //1=加储钱罐，2=加兑换额度
	MsgTime    time.Time
}
