package account

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

func OperationAccount(db *sqlx.DB, params Account) (err error) {
	info := UserInfo{}

	infoObj := new(UserInfo)

	info, err = infoObj.GetOneByUid(db, "id, uid, channel, vip, win_cash, cash, bonus, recharge, bet, transport, win", params.Uid)
	if err != nil {
		return
	}

	amount := info.WinCash
	cash := info.Cash
	freeze := info.Bonus
	recharge := info.Recharge

	tx, err := db.Beginx()

	if err != nil {
		return err
	}
	switch params.Cate {
	case 1, 3, 4, 8, 9:
		err = GetUpOtherAccount(amount, cash, freeze, recharge, params, tx)
	case 2:
		err = GetUpVipAccount()
	case 5:
		err = GetUpRechargeAccount(amount, cash, freeze, recharge, params, tx)
	case 6:
		err = GetUpSettleAccount()
	case 7:
		err = GetUpStoreAccount(amount, cash, freeze, recharge, params, tx)
	case 11:
		err = GetExchangeAccount(amount, cash, freeze, recharge, params, tx)
	default:
		err = errors.New("类型有误")
	}

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return
}

// 其他操作账户
func GetUpOtherAccount(amount, cash, freeze, recharge int64, params Account, tx *sqlx.Tx) (err error) {

	if params.Money < 1 && params.Cash < 1 && params.Store < 1 && params.Recharge < 1 {
		return errors.New("金额必须是正整数")
	}
	if params.Store > 0 {
		params.OptionType = 1
		params.IsBank = 1
		err = GetUpStoreAccount(amount, cash, freeze, recharge, params, tx)
		if err != nil {
			return err
		}
	}

	if params.Money > 0 || params.Cash > 0 || params.Recharge > 0 {
		err = GetAccountLog(amount, cash, freeze, recharge, params, tx)
		if err != nil {
			return err
		}
	}
	return err
}

func GetUpVipAccount() (err error) {
	return
}

func GetUpRechargeAccount(amount, cash, freeze, recharge int64, params Account, tx *sqlx.Tx) (err error) {
	if params.Money < 1 {
		return errors.New("金额必须是正整数")
	}

	if params.Store > 0 {
		params.OptionType = 1
		params.IsBank = 1
		err = GetUpStoreAccount(amount, cash, freeze, recharge, params, tx)
		if err != nil {
			return err
		}
	}

	err = GetAccountLog(amount, cash, freeze, recharge, params, tx)
	if err != nil {
		return err
	}
	return err
}

func GetUpSettleAccount() (err error) {
	return
}

func GetUpStoreAccount(amount, cash, freeze, recharge int64, params Account, tx *sqlx.Tx) (err error) {
	if params.Store < 1 {
		return errors.New("金额必须是正整数")
	}
	now := time.Now().Local()

	userInfoObj := new(UserInfo)

	userInfo, err := userInfoObj.GetOneByUidTx(tx, "id,uid,bonus,can_exchange_bonus", params.Uid)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	total, draw := userInfo.Bonus, userInfo.CanExchangeBonus

	if params.OptionType == 1 {
		if params.IsBank == 1 {
			bank := userInfo.Bonus + params.Store
			total = bank
			err = userInfoObj.UpdateTxBonusById(tx, bank, params.Uid, userInfo.ID)

			if err != nil {
				return err
			}
		} else if params.IsBank == 2 {
			withdrawal := userInfo.CanExchangeBonus + params.Store
			if withdrawal > userInfo.Bonus {
				withdrawal = userInfo.Bonus
			}
			draw = withdrawal
			err = userInfoObj.UpdateTxCanExchangeBonusById(tx, withdrawal, params.Uid, userInfo.ID)

			if err != nil {
				return err
			}
		}
	} else if params.OptionType == 2 {
		if userInfo.Bonus < 1 || userInfo.CanExchangeBonus < 1 {
			return errors.New("储钱罐数据异常")
		}
		bank := userInfo.Bonus - userInfo.CanExchangeBonus
		err = userInfoObj.UpdateTxBonusAndCanExchangeBonusById(tx, bank, 0, params.Uid, userInfo.ID)
		if err != nil {
			return err
		}
		err = GetAccountLog(amount, cash, freeze, recharge, params, tx)
		if err != nil {
			return err
		}
		draw = 0
	}
	//
	logObj := new(StoreBankLog)
	err = logObj.SaveTx(
		tx,
		"uuid,uid,ymd,total,draw,amount,type,category,created_at",
		":uuid,:uid,:ymd,:total,:draw,:amount,:type,:category,:created_at",
		map[string]interface{}{
			"uuid":       params.UUID,
			"uid":        params.Uid,
			"ymd":        params.Ymd,
			"total":      total,
			"draw":       draw,
			"amount":     params.Store,
			"type":       params.Cate,
			"category":   params.OptionType,
			"created_at": now,
		},
	)
	return err
}

func GetExchangeAccount(amount, cash, freeze, recharge int64, params Account, tx *sqlx.Tx) (err error) {
	if params.Money < 1 {
		return errors.New("金额必须是正整数")
	}

	err = GetAccountLog(amount, cash, freeze, recharge, params, tx)
	if err != nil {
		return err
	}
	return err

}
func GetAccountLog(amount, cash, bonus, recharge int64, params Account, tx *sqlx.Tx) (err error) {
	now := time.Now().Local()

	prizeObj := new(UserPrize)

	err = prizeObj.SaveTx(
		tx,
		"uuid,uid,game_id,ymd,type,win_cash,cash,bonus,category,tax,msg_time,created_at",
		":uuid,:uid,:game_id,:ymd,:type,:win_cash,:cash,:bonus,:category,:tax,:msg_time,:created_at",
		map[string]interface{}{
			"uuid":       params.UUID,
			"uid":        params.Uid,
			"game_id":    params.GameId,
			"ymd":        params.Ymd,
			"type":       params.Cate,
			"win_cash":   params.Money,
			"cash":       params.Cash,
			"bonus":      params.Store,
			"category":   params.OptionType,
			"tax":        params.Tax,
			"msg_time":   params.MsgTime,
			"created_at": now,
		},
	)

	if err != nil {
		return err
	}
	userInfoObj := new(UserInfo)
	accountObj := new(UserAccount)
	if params.Cate == 1 || params.Cate == 4 || params.Cate == 3 || params.Cate == 8 || params.Cate == 9 {
		money := amount + params.Money
		betAmount := cash + params.Cash
		recharge += params.Recharge
		optionMoney := params.Money + params.Cash

		err = userInfoObj.UpdateTxMoneyByUid(tx, money, betAmount, recharge, params.Uid)

		if err != nil {
			return err
		}

		err = accountObj.SaveTx(
			tx,
			"uid,ymd,uuid,win_cash,cash,bonus,`option`,tax,type,game_id,room,desk,created_at,updated_at",
			":uid,:ymd,:uuid,:win_cash,:cash,:bonus,:option,:tax,:type,:game_id,:room,:desk,:created_at,:updated_at",
			map[string]interface{}{
				"uid":        params.Uid,
				"ymd":        params.Ymd,
				"uuid":       params.UUID,
				"win_cash":   money,
				"cash":       betAmount,
				"bonus":      bonus,
				"option":     optionMoney,
				"tax":        params.Tax,
				"type":       params.Cate,
				"game_id":    params.GameId,
				"room":       0,
				"desk":       0,
				"created_at": now,
				"updated_at": now,
			},
			params.Uid,
		)

		return err
	}

	if params.Cate == 2 {

	}

	if params.Cate == 5 {
		newAccount := amount + params.Money
		newBetAmount := cash + params.Cash
		newFreeze := bonus + params.Store
		recharge += params.Recharge

		err = userInfoObj.UpdateTxAccountByUid(tx, newAccount, newBetAmount, newFreeze, recharge, params.Uid)

		if err != nil {
			return err
		}
		err = accountObj.SaveTx(
			tx,
			"uid,ymd,uuid,win_cash,cash,bonus,option,tax,type,game_id,room,desk,created_at,updated_at",
			":uid,:ymd,:uuid,:win_cash,:cash,:bonus,:option,:tax,:type,:game_id,:room,:desk,:created_at,:updated_at",
			map[string]interface{}{
				"uid":        params.Uid,
				"ymd":        params.Ymd,
				"uuid":       params.UUID,
				"win_cash":   newAccount,
				"cash":       cash,
				"bonus":      bonus,
				"option":     params.Money,
				"tax":        0,
				"type":       params.Cate,
				"game_id":    params.GameId,
				"room":       0,
				"desk":       0,
				"created_at": now,
				"updated_at": now,
			},
			params.Uid,
		)

		return err
	}

	if params.Cate == 6 {

	}

	if params.Cate == 7 {
		betAmount := cash + params.Store
		err = userInfoObj.UpdateTxCashByUid(tx, betAmount, params.Uid)
		if err != nil {
			return err
		}
		err = accountObj.SaveTx(
			tx,
			"uid,ymd,uuid,win_cash,cash,bonus,option,tax,type,game_id,room,desk,created_at,updated_at",
			":uid,:ymd,:uuid,:win_cash,:cash,:bonus,:option,:tax,:type,:game_id,:room,:desk,:created_at,:updated_at",
			map[string]interface{}{
				"uid":        params.Uid,
				"ymd":        params.Ymd,
				"uuid":       params.UUID,
				"win_cash":   amount,
				"cash":       betAmount,
				"bonus":      bonus,
				"option":     params.Store,
				"tax":        0,
				"type":       params.Cate,
				"game_id":    params.GameId,
				"room":       0,
				"desk":       0,
				"created_at": now,
				"updated_at": now,
			},
			params.Uid,
		)
		return err
	}

	if params.Cate == 11 {
		newAccount := amount
		logAccount := params.Money
		if params.OptionType == 1 {
			newAccount = amount + params.Money
		}

		if params.OptionType == 2 {
			newAccount = amount - params.Money
			logAccount = logAccount * -1
		}

		err = userInfoObj.UpdateTxWinCashByUid(tx, newAccount, params.Uid)

		if err != nil {
			return err
		}

		err = accountObj.SaveTx(
			tx,
			"uid,ymd,uuid,win_cash,cash,bonus,option,tax,type,game_id,room,desk,created_at,updated_at",
			":uid,:ymd,:uuid,:win_cash,:cash,:bonus,:option,:tax,:type,:game_id,:room,:desk,:created_at,:updated_at",
			map[string]interface{}{
				"uid":        params.Uid,
				"ymd":        params.Ymd,
				"uuid":       params.UUID,
				"win_cash":   newAccount,
				"cash":       cash,
				"bonus":      bonus,
				"option":     logAccount,
				"tax":        0,
				"type":       params.Cate,
				"game_id":    params.GameId,
				"room":       0,
				"desk":       0,
				"created_at": now,
				"updated_at": now,
			},
			params.Uid,
		)
		return err

	}
	return err
}
