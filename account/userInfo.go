package account

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/jmoiron/sqlx"
	"za.game/lib/consts"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
	"za.game/lib/sqlInfo"
)

type UserInfo struct {
	sqlInfo.UserInfo
}

func (t *UserInfo) TableName(uid uint64) string {
	return fmt.Sprintf("user_info_%02d", uid%5)
}

func (t *UserInfo) GetOneByUid(db *sqlx.DB, fields string, uid uint64) (u UserInfo, err error) {
	sql := fmt.Sprintf("select %s from %s where uid = ? limit 1", fields, t.TableName(uid))
	err = db.Get(&u, sql, uid)
	return
}

func (t *UserInfo) GetOneByUidTx(tx *sqlx.Tx, fields string, uid uint64) (u UserInfo, err error) {
	sql := fmt.Sprintf("select %s from %s where uid = ? limit 1", fields, t.TableName(uid))
	err = tx.Get(&u, sql, uid)
	return
}

func (t *UserInfo) Save(tx *sqlx.Tx, fields, values string, mp map[string]interface{}, uid uint64) (err error) {
	sql := fmt.Sprintf("insert into %s(%s) values(%s)", t.TableName(uid), fields, values)

	_, err = tx.NamedExec(sql, mp)
	return
}

func (t *UserInfo) UpdateTxCashByUid(tx *sqlx.Tx, cash int64, uid uint64) (err error) {
	sql := fmt.Sprintf("update %s set cash = ? where uid = ?", t.TableName(uid))

	_, err = tx.Exec(sql, cash, uid)
	return
}

func (t *UserInfo) UpdateTxWinCashByUid(tx *sqlx.Tx, winCash int64, uid uint64) (err error) {
	sql := fmt.Sprintf("update %s set win_cash = ? where uid = ?", t.TableName(uid))

	_, err = tx.Exec(sql, winCash, uid)
	return
}

func (t *UserInfo) UpdateTxMoneyByUid(tx *sqlx.Tx, winCash, cash, recharge int64, uid uint64) (err error) {
	sql := fmt.Sprintf("update %s set win_cash = ?, cash = ?, recharge = ? where uid = ?", t.TableName(uid))

	_, err = tx.Exec(sql, winCash, cash, recharge, uid)
	return
}

func (t *UserInfo) UpdateTxAccountByUid(tx *sqlx.Tx, winCash, cash, bonus, recharge int64, uid uint64) (err error) {
	sql := fmt.Sprintf("update %s set win_cash = ?, cash = ?,bonus = ?, recharge = ? where uid = ?", t.TableName(uid))

	_, err = tx.Exec(sql, winCash, cash, bonus, recharge, uid)
	return
}

func (t *UserInfo) UpdateTxBonusById(tx *sqlx.Tx, bonus int64, uid, id uint64) (err error) {
	sql := fmt.Sprintf("update %s set bonus = ? where id = ?", t.TableName(uid))

	_, err = tx.Exec(sql, bonus, id)
	return
}

func (t *UserInfo) UpdateTxBonusAndCanExchangeBonusById(tx *sqlx.Tx, bonus, canExchangeBonus int64, uid, id uint64) (err error) {
	sql := fmt.Sprintf("update %s set bonus = ?,can_exchange_bonus = ? where id = ?", t.TableName(uid))

	_, err = tx.Exec(sql, bonus, canExchangeBonus, id)
	return
}

func (t *UserInfo) UpdateTxCanExchangeBonusById(tx *sqlx.Tx, exchangeBonus int64, uid, id uint64) (err error) {
	sql := fmt.Sprintf("update %s set can_exchange_bonus = ? where id = ?", t.TableName(uid))

	_, err = tx.Exec(sql, exchangeBonus, id)
	return
}

func (t *UserInfo) UpdateWithdrawByUid(db *sqlx.DB, withdraw, withdrawed int, uid uint64) (err error) {
	sql := fmt.Sprintf("update %s set withdraw_money = withdraw_money + ?, withdrawed_money = withdrawed_money + ? where uid = ?", t.TableName(uid))

	_, err = db.Exec(sql, withdraw, withdrawed, uid)
	return
}

// 需要更新的用户金额
type UpdateUserMoneyInfo struct {
	Uid              uint64 `db:"uid" json:"uid"`                                //用户ID
	WinCash          int64  `db:"win_cash" json:"win_cash"`                      //可用金额
	Cash             int64  `db:"cash" json:"cash"`                              //金额
	Bonus            int64  `db:"bonus" json:"bonus"`                            //bonus
	CanExchangeBonus int64  `db:"can_exchange_bonus" json:"can_exchange_bonus" ` //可兑换bonus数量
	Recharge         int64  `db:"recharge" json:"recharge"`                      //充值的金额
	WithdrawMoney    int64  `db:"withdraw_money"`                                //提现审核中
	WithdrawedMoney  int64  `db:"withdrawed_money"`                              //提现到账
	ExtraGift        int64  `db:"extra_gift"`                                    //额外福利

	Sync bool //是否同步游戏服务器
}

// 更新后的用户金额
type LeftUserMoneyInfo struct {
	Uid                  uint64 `db:"uid" json:"uid"`                                     //用户ID
	Channel              int    `db:"channel" json:"channel"`                             //渠道id
	Vip                  int64  `db:"vip" json:"vip"`                                     //账户剩余vip经验
	LeftWinCash          int64  `db:"win_cash" json:"left_win_cash"`                      //更新后WinCash
	LeftCash             int64  `db:"cash" json:"left_cash"`                              //更新后Cash
	LeftBonus            int64  `db:"bonus" json:"left_bonus"`                            //更新后bonus
	LeftCanExchangeBonus int64  `db:"can_exchange_bonus" json:"left_can_exchange_bonus" ` //更新后可兑换bonus数量
	Recharge             int64  `db:"recharge" json:"recharge"`
	WithdrawMoney        int    `db:"withdraw_money" json:"withdraw_money"`
	WithdrawedMoney      int    `db:"withdrawed_money" json:"withdrawed_money"`
	Fxlb_yj_keti         int64  `db:"fxlb_yj_keti" json:"fxlb_yj_keti"` //分享裂变可提现金额
}

func GetLeftUserMoneyInfo(db *sqlx.DB, uid uint64) (*LeftUserMoneyInfo, error) {
	left := &LeftUserMoneyInfo{}
	err := rds.SqlxGetD(db, left, "select uid,channel,vip,win_cash, bonus, cash,can_exchange_bonus,recharge,withdraw_money,withdrawed_money from %s where uid = ?",
		sqlInfo.GetUserInfoTableName(uid), uid)

	return left, err
}

// 更新后物品的信息
type LeftItemInfo struct {
	LeftNum int64 `db:"num" json:"left_num"` // 更新后物品的数量
}

func AddSimpleCash(db *sqlx.DB, uid uint64, addNum int64) (int64, error) {
	if uid == 0 || addNum <= 0 {
		return 0, errors.New("param err")
	}
	tx, _ := db.Beginx()
	_, err := rds.SqlxTxExecD(tx, "update %s set cash=cash+?  where uid=?",
		sqlInfo.GetUserInfoTableName(uid), addNum, uid)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	left := &LeftUserMoneyInfo{}
	err = rds.SqlxTxGetD(tx, left, "select cash from %s where uid = ?",
		sqlInfo.GetUserInfoTableName(uid), uid)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return left.LeftCash, err
}
func SubSimpleWinCash(db *sqlx.DB, uid uint64, addNum int64) (*LeftUserMoneyInfo, error) {
	if uid == 0 || addNum <= 0 {
		return nil, errors.New("param err")
	}
	tx, _ := db.Beginx()
	_, err := rds.SqlxTxExecD(tx, "update %s set win_cash=win_cash-?  where uid=? and win_cash>=?",
		sqlInfo.GetUserInfoTableName(uid), addNum, uid, addNum)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	left := &LeftUserMoneyInfo{}
	err = rds.SqlxTxGetD(tx, left, "select uid,win_cash, bonus, cash,can_exchange_bonus,recharge from %s where uid = ?",
		sqlInfo.GetUserInfoTableName(uid), uid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return left, err
}
func AddSimpleBonus(db *sqlx.DB) (int64, error) {
	return 0, nil
}
func AddUserMoney(db *sqlx.DB, logtype int, inInfo UpdateUserMoneyInfo, logdb *sqlx.DB, fileds string, args ...interface{}) (*LeftUserMoneyInfo, error) {
	if inInfo.Uid == 0 || (inInfo.Cash == 0 && inInfo.WinCash == 0 && inInfo.Bonus == 0 && inInfo.CanExchangeBonus == 0 && inInfo.WithdrawedMoney == 0 && inInfo.WithdrawMoney == 0) {
		return nil, errors.New("param err")
	}
	var sqlStr strings.Builder
	sqlStr.WriteString("update ")
	sqlStr.WriteString(sqlInfo.GetUserInfoTableName(inInfo.Uid))
	sqlStr.WriteString(" set ")
	isFirst := true
	if inInfo.Cash != 0 {
		if !isFirst {
			sqlStr.WriteString(" , ")
		}
		sqlStr.WriteString(" cash = cash+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.Cash, 10))
		isFirst = false
	}
	if inInfo.WinCash != 0 {
		if !isFirst {
			sqlStr.WriteString(" , ")
		}
		sqlStr.WriteString(" win_cash = win_cash+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.WinCash, 10))
		isFirst = false
	}
	if inInfo.Bonus != 0 {
		if !isFirst {
			sqlStr.WriteString(" , ")
		}
		sqlStr.WriteString(" bonus = bonus+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.Bonus, 10))
		isFirst = false
	}
	if inInfo.CanExchangeBonus != 0 {
		if !isFirst {
			sqlStr.WriteString(" , ")
		}
		sqlStr.WriteString(" can_exchange_bonus = can_exchange_bonus+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.CanExchangeBonus, 10))
		isFirst = false
	}
	if inInfo.Recharge != 0 {
		if !isFirst {
			sqlStr.WriteString(" , ")
		}
		sqlStr.WriteString(" recharge=recharge+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.Recharge, 10))
		//这里，要把充值的钱转为vip经验
		sqlStr.WriteString(" , ")
		sqlStr.WriteString(" vip=vip+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.Recharge, 10))
		sqlStr.WriteString(" , ")
		sqlStr.WriteString(" recharge_count=recharge_count+1")
		isFirst = false
	}

	if inInfo.WithdrawMoney != 0 {
		if !isFirst {
			sqlStr.WriteString(" , ")
		}
		sqlStr.WriteString(" withdraw_money = withdraw_money+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.WithdrawMoney, 10))
		isFirst = false
	}

	if inInfo.WithdrawedMoney != 0 {
		if !isFirst {
			sqlStr.WriteString(" , ")
		}
		sqlStr.WriteString(" withdrawed_money = withdrawed_money+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.WithdrawedMoney, 10))
		isFirst = false
	}
	if inInfo.ExtraGift != 0 { //新增额外福利数据记录，不需要记录流水日志，在surprise礼包日志里面有
		if !isFirst {
			sqlStr.WriteString(" , ")
		}
		sqlStr.WriteString(" extra_gift = extra_gift+")
		sqlStr.WriteString(strconv.FormatInt(inInfo.ExtraGift, 10))
		isFirst = false
	}

	sqlStr.WriteString(" where uid =")
	sqlStr.WriteString(strconv.FormatInt(int64(inInfo.Uid), 10))
	tx, _ := db.Beginx()
	_, err := rds.SqlxExecTx(tx, sqlStr.String())
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	left := &LeftUserMoneyInfo{}
	err = rds.SqlxTxGetD(tx, left, "select uid,channel,vip,win_cash, bonus, cash,can_exchange_bonus,recharge,withdraw_money,withdrawed_money from %s where uid = ?",
		sqlInfo.GetUserInfoTableName(inInfo.Uid), inInfo.Uid)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	//记录充值项 记录

	if inInfo.Recharge != 0 {
		log_pool.Submit(func() {
			column := make(map[string]interface{}, 0)
			column["uid"] = inInfo.Uid
			column["recharge"] = inInfo.Recharge
			if ret, errQ := rds.SqlxNamedExec(db, "update user_charge set num=num+1,update_time=CURRENT_TIMESTAMP where uid=:uid and recharge=:recharge", column); errQ == nil {
				if af, errT := ret.RowsAffected(); af == 0 && errT == nil {
					if _, errQ := rds.SqlxNamedExec(db, "insert into user_charge(uid,recharge,num) values(:uid,:recharge,1)", column); errQ != nil {
						logger.Logger.Errorf("user_charge: insert failed! err:[%v]", errQ)
					}
				}
			} else {
				logger.Logger.Errorf("user_charge: insert failed! err:[%v]", errQ)
			}
		})
	}

	if inInfo.Cash != 0 || inInfo.WinCash != 0 || inInfo.Bonus != 0 || inInfo.CanExchangeBonus != 0 {
		//记录日志
		LogMoney(logdb, left, logtype, inInfo.Cash, inInfo.WinCash, inInfo.Bonus, inInfo.CanExchangeBonus, fileds, args...)
	}

	return left, err
}

/**增加物品*/
func AddItem(uid uint64, itemid int64, num int, logtype int, fileds string, args ...interface{}) (err error, left *LeftItemInfo) {
	if itemid == consts.ItemIdNewPlayerCoin { //新手嘉年华活动币
		err, left = sendAwardNewPlayerCoin(uid, num)
	}

	//记录日志
	LogItems(uid, left, logtype, itemid, int64(num), fileds, args...)
	return err, left
}

/**新手嘉年华活动币*/
func sendAwardNewPlayerCoin(uid uint64, num int) (err error, left *LeftItemInfo) {
	award_id := consts.ItemIdNewPlayerCoin
	left = &LeftItemInfo{}
	update_columns := map[string]interface{}{
		"uid":       uid,
		"award_id":  award_id,
		"award_num": num,
	}
	sqlres, err := rds.SqlxNamedExec(dbconn.NDB, "update task_user_award set award_num=award_num+:award_num where uid=:uid and  award_id=:award_id", update_columns)
	if err != nil {
		logger.Logger.Errorf("sendAwardNewPlayerCoin: update sql award failed! err:[%v]", err)
		return
	}
	if iret, _ := sqlres.RowsAffected(); iret <= 0 {
		logger.Logger.Errorf("sendAwardNewPlayerCoin: send award failed! uid:[%d],awardid:[%d],awardnum:[%d]", uid, award_id, num)
		return
	}
	//查看账户剩余
	var task_user_award_db_info TaskUserAwardDB
	err = rds.SqlxGet(dbconn.NDB, &task_user_award_db_info, "select * from task_user_award where uid=? and award_id=?", uid, award_id)
	if err != nil {
		logger.Logger.Errorf("func: select sql failed! err:[%v]", err)
		return nil, left
	}
	left.LeftNum = int64(task_user_award_db_info.AwardNum)
	return
}

type AddFxlbParams struct {
	Uid           uint64 //用户id
	FromUid       uint64 //邀请者id
	YqYj          int    //邀请佣金
	CzYj          int    //充值佣金
	DlYj          int    //代理佣金
	KeTi          int    //可提佣金
	WithdrawMoney int64  //提现审核中
	Remark        string //备注
	RoomId        int64  //房间id
}

/**增加分享裂变对应值*/
func AddFxlbItem(params *AddFxlbParams, logtype int) (left *LeftUserMoneyInfo, err error) {
	if (params.Uid == 0 && params.FromUid == 0) || (params.YqYj == 0 && params.CzYj == 0 && params.DlYj == 0 && params.KeTi == 0) {
		return nil, errors.New("param err")
	}
	update_columns := make(map[string]interface{}, 0)
	if params.Uid > 0 {
		update_columns["uid"] = params.Uid
	}
	add_yj := params.YqYj + params.CzYj + params.DlYj
	if add_yj > 0 {
		//先增加自己，给别人提供的未转化的佣金
		update_columns["fxlb_yj_wei"] = add_yj
		update_sql := fmt.Sprintf("update %s set fxlb_yj_wei=fxlb_yj_wei+:fxlb_yj_wei where uid=:uid", sqlInfo.GetUserInfoTableName(params.Uid))
		_, err = rds.SqlxNamedExec(dbconn.NDB, update_sql, update_columns)
		if err != nil {
			logger.Logger.Errorf("AddFxlbItem: 分享裂变-增加给【邀请者】的【未转化】的佣金失败! err:[%v],uid:[%v],from_uid:[%v],num:[%v]", err, params.Uid, params.FromUid, add_yj)
			return
		}
		//查看账户未转剩余
		left_fxlb_yj_wei := int64(0)
		select_sql := fmt.Sprintf("select fxlb_yj_wei from %s where uid=?", sqlInfo.GetUserInfoTableName(params.Uid))
		rows, err2 := rds.QueryX(dbconn.NDB, select_sql, params.Uid)
		if err2 != nil {
			logger.Logger.Errorf("AddFxlbItem: 分享裂变-查询账户未转剩余失败! err:[%v]", err2)
		}
		defer rows.Close()
		if rows.Next() {
			if err2 := rows.Scan(&left_fxlb_yj_wei); err2 != nil {
				logger.Logger.Errorf("AddFxlbItem: 分享裂变-解析账户未转剩余失败! err:[%v]", err2)
			}
		}
		//给邀请者，增加总佣金
		update_columns = make(map[string]interface{}, 0)
		update_columns["uid"] = params.FromUid
		update_columns["fxlb_yj"] = add_yj
		update_sql = fmt.Sprintf("update %s set fxlb_yj=fxlb_yj+:fxlb_yj where uid=:uid", sqlInfo.GetUserInfoTableName(params.FromUid))
		_, err = rds.SqlxNamedExec(dbconn.NDB, update_sql, update_columns)
		if err != nil {
			logger.Logger.Errorf("AddFxlbItem: 分享裂变-增加给【邀请者】的【总佣金】的佣金失败! err:[%v],uid:[%v],from_uid:[%v],num:[%v]", err, params.Uid, params.FromUid, add_yj)
			return
		}
		//查看邀请者账户剩余
		left_fxlb_yj := int64(0)
		select_sql = fmt.Sprintf("select fxlb_yj from %s where uid=?", sqlInfo.GetUserInfoTableName(params.FromUid))
		rows, err := rds.QueryX(dbconn.NDB, select_sql, params.FromUid)
		if err != nil {
			logger.Logger.Errorf("AddFxlbItem: 分享裂变-查询账户剩余失败! err:[%v]", err)
		}
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&left_fxlb_yj); err != nil {
				logger.Logger.Errorf("AddFxlbItem: 分享裂变-解析账户剩余失败! err:[%v]", err)
			}
		}

		//记录流水
		LogFxlbYj(&LogFxlbParams{
			Uid:       params.Uid,
			WeiZh:     int64(add_yj),
			LeftWeiZh: left_fxlb_yj_wei,
			Type:      logtype,
			ToUid:     params.FromUid,
			ToYj:      int64(add_yj),
			LeftToYj:  left_fxlb_yj,
			Remark:    params.Remark,
			RoomId:    params.RoomId,
		})

		//记录流水统计
		LogFxlbState(params)
	} else if params.FromUid == 0 && params.KeTi != 0 { //这是本人佣金提现
		update_columns["fxlb_yj_keti"] = params.KeTi
		update_columns["fxlb_yj_sh"] = params.WithdrawMoney
		update_sql := fmt.Sprintf("update %s set fxlb_yj_keti=fxlb_yj_keti+:fxlb_yj_keti,fxlb_yj_sh=fxlb_yj_sh+:fxlb_yj_sh where uid=:uid", sqlInfo.GetUserInfoTableName(params.Uid))
		_, err = rds.SqlxNamedExec(dbconn.NDB, update_sql, update_columns)
		if err != nil {
			logger.Logger.Errorf("AddFxlbItem: 分享裂变-扣除可提佣金失败! err:[%v],uid:[%v],num:[%v]", err, params.Uid, params.KeTi)
			return
		}

		//查看邀请者账户剩余
		left = &LeftUserMoneyInfo{}
		select_sql := fmt.Sprintf("select uid,channel,vip,win_cash, bonus, cash,can_exchange_bonus,recharge,withdraw_money,withdrawed_money,fxlb_yj_keti from %s where uid = ?", sqlInfo.GetUserInfoTableName(params.Uid))
		err = rds.SqlxGet(dbconn.NDB, left, select_sql, params.Uid)
		if err != nil {
			logger.Logger.Errorf("AddFxlbItem: 分享裂变-获取账户剩余信息失败! err:[%v],uid:[%v],num:[%v]", err, params.Uid, params.KeTi)
			return
		}

		//记录流水
		LogFxlbYj(&LogFxlbParams{
			Uid:      params.Uid,
			Type:     logtype,
			Keti:     params.KeTi,
			LeftKeti: left.Fxlb_yj_keti,
			Remark:   params.Remark,
			RoomId:   params.RoomId,
		})
	}
	return
}
