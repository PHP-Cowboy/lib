package pay

import (
	"time"

	"za.game/lib/account"

	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"za.game/lib/consts"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
	libRes "za.game/lib/response"
)

/**提现参数*/
type GiveMoneyParams struct {
	Uid          uint64 //用户id
	Amount       int    //提现金额
	PayStatus    int
	LeftUserInfo *account.LeftUserMoneyInfo //用户剩余金额
	OrderNo      string                     //订单号
	Name         string                     //提现人姓名
	Type         consts.GiveMoneyType       //提现类型
	GiveId       int                        //提现id
}

func AddRecord(params *GiveMoneyParams) (err *libRes.Response) {

	now := time.Now()
	leftUserInfo := params.LeftUserInfo
	giveMoneyItem := map[string]any{
		"uid":          params.Uid,
		"order_no":     params.OrderNo,
		"amount":       params.Amount,
		"amount_total": leftUserInfo.WithdrawMoney + leftUserInfo.WithdrawedMoney, //赠送总币 = 已成功的赠送额度（数据库字段） + 已提交但是未到账的赠送额度（数据库字段）
		"recharge":     leftUserInfo.Recharge,
		"pay_status":   params.PayStatus,
		"nick_name":    params.Name,
		"channel_id":   leftUserInfo.Channel,
		"created_at":   now,
		"updated_at":   now,
		"type":         params.Type,
		"give_id":      params.GiveId,
	}

	//give_rate 实际赠送金币比 = 已成功的赠送额度（数据库字段） /  订单提交的时候的玩家的充值总币
	//commit_give_rate 提交赠送金币比 = （已成功的赠送额度（数据库字段） + 已提交但是未到账的赠送额度（数据库字段）） /  订单提交的时候的玩家的充值总币
	if leftUserInfo.Recharge > 0 {
		giveMoneyItem["give_rate"] = float64(leftUserInfo.WithdrawedMoney) / float64(leftUserInfo.Recharge)
		giveMoneyItem["commit_give_rate"] = float64(leftUserInfo.WithdrawMoney+leftUserInfo.WithdrawedMoney) / float64(leftUserInfo.Recharge)
	} else {
		giveMoneyItem["give_rate"] = 0
		giveMoneyItem["commit_give_rate"] = 0
	}

	_, err2 := rds.SqlxNamedExec(
		dbconn.PayDB,
		"insert into give_money(created_at, updated_at, order_no,  uid, nick_name,channel_id, amount,amount_total,recharge,give_rate,commit_give_rate,pay_status,`type`,give_id) "+
			"VALUES(:created_at,:updated_at,:order_no,:uid,:nick_name,:channel_id,:amount,:amount_total,:recharge,:give_rate,:commit_give_rate,:pay_status,:type,:give_id)",

		giveMoneyItem,
	)

	if err2 != nil {
		logger.Logger.Errorf("giveMoney: insert recorde failed! err:[%v]", err2)
		return &consts.DataSaveError
	}

	return
}
