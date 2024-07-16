package account

import (
	"fmt"
	"strings"
	"time"

	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/alitto/pond"
	"github.com/jmoiron/sqlx"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
	"za.game/lib/tool"
)

var log_pool *pond.WorkerPool

const (
	LogType_Default                   int = 0  //默认
	LogType_day_task                  int = 1  //每日任务
	LogType_vip                       int = 2  //VIP
	LogType_email                     int = 3  //邮件
	LogType_sign                      int = 4  //签到
	LogType_recharge                  int = 5  //充值
	LogType_frozen                    int = 6  //冻结
	LogType_bank                      int = 7  //储钱罐
	LogType_benefit                   int = 8  //救济金
	LogType_giftbag                   int = 9  //礼包
	LogType_register_send             int = 10 //注册赠送
	LogType_luck_spin                 int = 11 //luckspin
	LogType_withdraw                  int = 12 //赠送
	LogType_admin                     int = 13 //管理后台操作
	LogType_RechargeGift              int = 14 //充值赠送礼包
	LogType_EventGift                 int = 15 //发货二选一礼包
	LogType_OnlyOneGift               int = 16 //发货OnlyOne三选一礼包
	LogType_BenefitGift               int = 17 //救济金礼包
	LogType_RechargeRoomGift          int = 18 //房间特惠礼包
	LogType_Activity_NewPlayer        int = 19 //新手嘉年华
	LogType_Bonus_Draw                int = 20 //bonus 提取
	LogType_Tyro_Cash                 int = 21 //新用户首次
	LogType_luck_spin_recharge        int = 22 //luckspin充值
	LogType_WithdrawSend              int = 23 //withdraw 添加首次添加账户
	LogType_SurpriseGift              int = 24 //Surprise礼包充值
	LogType_Activity_NewPlayerActCoin int = 25 //
	LogType_SurpriseGiftFree          int = 26 //Surprise礼包充值(免费)
	LogType_FxlbYqYj                  int = 27 //分享裂变邀请佣金
	LogType_FxlbCzYj                  int = 28 //分享裂变充值佣金
	LogType_FxlbDlYj                  int = 28 //分享裂变代理佣金
	LogType_FxlbYjWithdraw            int = 29 //分享裂变提现
	LogType_FxlbYjZhuan               int = 30 //分享裂变佣金转化
)

func init() {
	//创建一个线程池，工作线程（1），最多任务（100）。
	log_pool = pond.New(1, 100)
}

/**记录金币流水
*@param left: 账户剩余信息
*@param cashNum: 变动的 cash 金额
*@param winCashNum: 变动 winCash 的金额
*@param bonus: 存储罐金额
*@param can_exchange_bonus: 可提现存储罐金额
*@param fileds: 日志数据库额外字段(用"," 分割)
*@param args: 可变参数，对应额外字段的值*/
func LogMoney(logdb *sqlx.DB, left *LeftUserMoneyInfo, logtype int, cashNum int64, winCashNum int64, bonus int64, can_exchange_bonus int64, fileds string, args ...interface{}) {
	//用户离开同步数据
	insert_columns := map[string]interface{}{
		"created_at":              time.Now(),
		"ymd":                     tool.TimeGetYmd(),
		"uid":                     left.Uid,
		"type":                    logtype,
		"win_cash":                winCashNum,
		"left_win_cash":           left.LeftWinCash,
		"cash":                    cashNum,
		"left_cash":               left.LeftCash,
		"bonus":                   bonus,
		"left_bonus":              left.LeftBonus,
		"can_exchange_bonus":      can_exchange_bonus,
		"left_can_exchange_bonus": left.LeftCanExchangeBonus,
	}
	arr_filed := strings.Split(fileds, ",")
	fileds_value_name := ""
	if len(arr_filed) != len(args) {
		logger.Logger.Errorf("LogMoney: the len of fileds and args are not same! fileds:[%v],args:[%v]", fileds, args)
		return
	}
	for i := 0; i < len(arr_filed); i++ {
		if fileds_value_name != "" {
			fileds_value_name += ","
		}
		fileds_value_name += ":" + arr_filed[i]
		key := arr_filed[i]
		value := args[i]
		insert_columns[key] = value
	}

	log_pool.Submit(func() {
		tableYm := time.Now().Local().Format("200601")
		tablename := "user_funds_flow_log_" + tableYm
		sql_str := fmt.Sprintf("insert into %s(created_at,ymd,uid,type,win_cash,left_win_cash,cash,left_cash,bonus,left_bonus,can_exchange_bonus,left_can_exchange_bonus,%s) values(:created_at,:ymd,:uid,:type,:win_cash,:left_win_cash,:cash,:left_cash,:bonus,:left_bonus,:can_exchange_bonus,:left_can_exchange_bonus,%s)", tablename, fileds, fileds_value_name)
		_, err := rds.SqlxNamedExec(logdb, sql_str, insert_columns)
		if err != nil {
			logger.Logger.Errorf("LogMoney: insert log failed! err:[%v]", err)
			return
		}
		//用户流水属性
		if cashNum+winCashNum != 0 || bonus != 0 {
			err = UpdateRecentRecord(left.Uid, UserRecordItem{
				Cash:          int(cashNum + winCashNum),
				Bonus:         int(bonus),
				Type:          logtype,
				RecordTimeStr: time.Now().Format("2006-01-02 15:04:05"),
			})
			if err != nil {
				logger.Logger.Errorf("LogMoney: UpdateRecentRecord failed! err:[%v]", err)
				return
			}
		}
	})
}

/**记录物品流水
*@param left: 账户剩余信息
*@param cashNum: 变动的 cash 金额
*@param winCashNum: 变动 winCash 的金额
*@param bonus: 存储罐金额
*@param can_exchange_bonus: 可提现存储罐金额
*@param fileds: 日志数据库额外字段(用"," 分割)
*@param args: 可变参数，对应额外字段的值*/
func LogItems(uid uint64, left *LeftItemInfo, logtype int, itemid int64, addnum int64, fileds string, args ...interface{}) {
	//用户离开同步数据
	insert_columns := map[string]interface{}{
		"created_at": time.Now(),
		"uid":        uid,
		"type":       logtype,
		"itemid":     itemid,
		"addnum":     addnum,
		"leftnum":    left.LeftNum,
	}
	arr_filed := strings.Split(fileds, ",")
	fileds_value_name := ""
	if len(arr_filed) != len(args) {
		logger.Logger.Errorf("LogItems: the len of fileds and args are not same! fileds:[%v],args:[%v]", fileds, args)
		return
	}
	for i := 0; i < len(arr_filed); i++ {
		if fileds_value_name != "" {
			fileds_value_name += ","
		}
		fileds_value_name += ":" + arr_filed[i]
		key := arr_filed[i]
		value := args[i]
		insert_columns[key] = value
	}

	log_pool.Submit(func() {
		tableYm := time.Now().Local().Format("200601")
		tablename := "items_log_" + tableYm
		sql_str := fmt.Sprintf("insert into %s(created_at,uid,type,itemid,addnum,leftnum,%s) values(:created_at,:uid,:type,:itemid,:addnum,:leftnum,%s)", tablename, fileds, fileds_value_name)
		_, err := rds.SqlxNamedExec(dbconn.LogDB, sql_str, insert_columns)
		if err != nil {
			logger.Logger.Errorf("LogItems: insert log failed! err:[%v]", err)
			return
		}
	})
}

type LogFxlbParams struct {
	Uid        uint64 //用户id
	Keti       int    //可提现佣金变动的值
	LeftKeti   int64  //可提现佣金剩余的值
	WeiZh      int64  //未转变动的值
	LeftWeiZh  int64  //未转剩余的值
	Type       int    //类型
	ToUid      uint64 //对方id
	ToYj       int64  //	对方佣金变动的值
	LeftToYj   int64  //	对方账户佣金剩余值
	ToKeti     int    //	对方可提现佣金变动的值
	LeftToKeti int64  //	对方可提现佣金剩余值
	Remark     string //备注
	RoomId     int64  //房间id
}

/**记录分享裂变佣金流水，和流水统计*/
func LogFxlbYj(params *LogFxlbParams) {
	insert_columns := map[string]interface{}{
		"created_at": time.Now(),
		"uid":        params.Uid,
		"type":       params.Type,
		"to_uid":     params.ToUid,
	}
	filed_list := ""
	value_list := ""
	if params.Keti != 0 {
		insert_columns["keti"] = params.Keti
		insert_columns["leftketi"] = params.LeftKeti
		if filed_list != "" {
			filed_list += ","
			value_list += ","
		}
		filed_list += "keti,leftketi"
		value_list += ":keti,:leftketi"
	}
	if params.ToYj != 0 {
		insert_columns["to_yj"] = params.ToYj
		insert_columns["left_to_yj"] = params.LeftToYj
		if filed_list != "" {
			filed_list += ","
			value_list += ","
		}
		filed_list += "to_yj,left_to_yj"
		value_list += ":to_yj,:left_to_yj"
	}
	if params.ToKeti != 0 {
		insert_columns["to_keti"] = params.ToKeti
		insert_columns["left_to_keti"] = params.LeftToKeti
		if filed_list != "" {
			filed_list += ","
			value_list += ","
		}
		filed_list += "to_keti,left_to_keti"
		value_list += ":to_keti,:left_to_keti"
	}
	if params.WeiZh != 0 {
		insert_columns["weizh"] = params.WeiZh
		insert_columns["left_weizh"] = params.LeftWeiZh
		if filed_list != "" {
			filed_list += ","
			value_list += ","
		}
		filed_list += "weizh,left_weizh"
		value_list += ":weizh,:left_weizh"
	}
	if params.Remark != "" {
		insert_columns["remark"] = params.Remark
		if filed_list != "" {
			filed_list += ","
			value_list += ","
		}
		filed_list += "remark"
		value_list += ":remark"
	}
	if params.RoomId != 0 {
		insert_columns["roomid"] = params.RoomId
		if filed_list != "" {
			filed_list += ","
			value_list += ","
		}
		filed_list += "roomid"
		value_list += ":roomid"
	}

	log_pool.Submit(func() {
		tableYm := time.Now().Local().Format("200601")
		tablename := "fxlb_log_" + tableYm
		sql_str := fmt.Sprintf("insert into %s(created_at,uid,type,to_uid,%s) values(:created_at,:uid,:type,:to_uid,%s)", tablename, filed_list, value_list)
		_, err := rds.SqlxNamedExec(dbconn.LogDB, sql_str, insert_columns)
		if err != nil {
			logger.Logger.Errorf("LogFxlbYj: 分享裂变-记录佣金流水日志失败! err:[%v]", err)
			return
		}
	})
}

/**佣金流水统计*/
func LogFxlbState(params *AddFxlbParams) {
	update_columns := map[string]interface{}{
		"day":   time.Now().Format("20060102"),
		"uid":   params.FromUid,
		"yq_yj": params.YqYj,
		"cz_yj": params.CzYj,
		"dl_yj": params.DlYj,
	}

	log_pool.Submit(func() {
		tableYm := time.Now().Local().Format("200601")
		tablename := "fxlb_stat_log_" + tableYm
		sql_str := fmt.Sprintf("update %s set yq_yj=yq_yj+:yq_yj,cz_yj=cz_yj+:cz_yj,cz_yj=cz_yj+:dl_yj where day=:day and uid=:uid", tablename)
		set, err := rds.SqlxNamedExec(dbconn.LogDB, sql_str, update_columns)
		if err != nil {
			logger.Logger.Errorf("LogFxlbState: 分享裂变-更新用户佣金统计失败! err:[%v]", err)
			return
		}
		if rows, _ := set.RowsAffected(); rows == 0 { //没有数据，则插入
			sql_str := fmt.Sprintf("insert into %s(day,uid,yq_yj,cz_yj,dl_yj) values(:day,:uid,:yq_yj,:cz_yj,:dl_yj)", tablename)
			_, err = rds.SqlxNamedExec(dbconn.LogDB, sql_str, update_columns)
			if err != nil {
				logger.Logger.Errorf("LogFxlbState: 分享裂变-插入用户佣金统计失败! err:[%v]", err)
				return
			}
		}
	})
}
