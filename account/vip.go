package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"za.game/lib/consts"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
	"za.game/lib/sqlInfo"
	"za.game/lib/tool"
)

/**vip数据表结构体*/
type VipConfigDB struct {
	Id            int          `db:"id" gorm:"primarykey;column:id;type:int;" json:"id" `                             //
	CreatedAt     sql.NullTime `db:"created_at" gorm:"column:created_at;type:sql.NullTime;" json:"created_at" `       //
	UpdatedAt     sql.NullTime `db:"updated_at" gorm:"column:updated_at;type:sql.NullTime;" json:"updated_at" `       //
	Level         int8         `db:"level" gorm:"column:level;type:int8;" json:"level" `                              // VIP等级
	NeedExp       int          `db:"need_exp" gorm:"column:need_exp;type:int;" json:"need_exp" `                      // 升级所需经验
	WithdrawNums  int8         `db:"withdraw_nums" gorm:"column:withdraw_nums;type:int8;" json:"withdraw_nums" `      // 每日提现次数
	WithdrawMoney int          `db:"withdraw_money" gorm:"column:withdraw_money;type:int;" json:"withdraw_money" `    // 每日提现金额
	DayAwardId    int64        `db:"day_award_id" gorm:"column:day_award_id;type:int64;" json:"day_award_id" `        // 每日奖品id
	DayAwardNum   int          `db:"day_award_num" gorm:"column:day_award_num;type:int;" json:"day_award_num" `       // 每日奖品数量
	WeekAwardId   int64        `db:"week_award_id" gorm:"column:week_award_id;type:int64;" json:"week_award_id" `     // 每周奖品id
	WeekAwardNum  int          `db:"week_award_num" gorm:"column:week_award_num;type:int;" json:"week_award_num" `    // 每周奖品数量
	MonthAwardId  int64        `db:"month_award_id" gorm:"column:month_award_id;type:int64;" json:"month_award_id" `  // 每月奖品id
	MonthAwardNum int          `db:"month_award_num" gorm:"column:month_award_num;type:int;" json:"month_award_num" ` // 每月奖品数量
	CashBankRate  int          `db:"cash_bank_rate" gorm:"column:cash_bank_rate;type:int;" json:"cash_bank_rate" `    // 存钱罐比例（千分比）
	CashBankTop   int          `db:"cash_bank_top" gorm:"column:cash_bank_top;type:int;" json:"cash_bank_top" `       // 存钱罐上限
}

/**vip信息结构体-用于内存中使用和redis中存储*/
type VipConfigRedis struct {
	Id            int   `db:"id" gorm:"primarykey;column:id;type:int;" json:"id" `                             //
	Level         int8  `db:"level" gorm:"column:level;type:int8;" json:"level" `                              // VIP等级
	NeedExp       int   `db:"need_exp" gorm:"column:need_exp;type:int;" json:"need_exp" `                      // 升级所需经验
	WithdrawNums  int8  `db:"withdraw_nums" gorm:"column:withdraw_nums;type:int8;" json:"withdraw_nums" `      // 每日提现次数
	WithdrawMoney int   `db:"withdraw_money" gorm:"column:withdraw_money;type:int;" json:"withdraw_money" `    // 每日提现金额
	DayAwardId    int64 `db:"day_award_id" gorm:"column:day_award_id;type:int8;" json:"day_award_id" `         // 每日奖品id
	DayAwardNum   int   `db:"day_award_num" gorm:"column:day_award_num;type:int;" json:"day_award_num" `       // 每日奖品数量
	WeekAwardId   int64 `db:"week_award_id" gorm:"column:week_award_id;type:int8;" json:"week_award_id" `      // 每周奖品id
	WeekAwardNum  int   `db:"week_award_num" gorm:"column:week_award_num;type:int;" json:"week_award_num" `    // 每周奖品数量
	MonthAwardId  int64 `db:"month_award_id" gorm:"column:month_award_id;type:int8;" json:"month_award_id" `   // 每月奖品id
	MonthAwardNum int   `db:"month_award_num" gorm:"column:month_award_num;type:int;" json:"month_award_num" ` // 每月奖品数量
	CashBankRate  int   `db:"cash_bank_rate" gorm:"column:cash_bank_rate;type:int;" json:"cash_bank_rate" `    // 存钱罐比例（千分比）
	CashBankTop   int   `db:"cash_bank_top" gorm:"column:cash_bank_top;type:int;" json:"cash_bank_top" `       // 存钱罐上限
}

/**用户相关信息*/
type VipUserInfo struct {
	Vip             int   //vip经验值
	DayGeted        int64 //每日奖励，下次可领取时间
	WeekGeted       int64 //每周奖励，下次可领取时间
	MonGeted        int64 //每月奖励，下次可领取时间
	VipCashBank1    int   // 可以领取的vip存钱罐奖励
	VipCashBank2    int   // 当天累计的vip存钱罐奖励
	VipCashBankTime int64 // 上一次累计的时间
	WithdrawNum     int   // 每日已提现次数
	WithdrawMoney   int   // 每日已提现额度
	WithdrawTime    int64 // 上次提现的时间
}

/**获得vip配置信息*/
func GetVipConfigData() (arrdbinfo []VipConfigRedis, iret int) {
	//获得活动配置信息
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()

	reply, err := redis.String(redisConn.Do("GET", consts.RedisVipConfig))

	arrdbinfo = make([]VipConfigRedis, 0)

	if reply != "" {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &arrdbinfo)
		if err != nil {
			logger.Logger.Errorf("getVipConfigData: parse config failed! err:[%v]", err)
		}
	}

	if err != nil {
		var arr_vip_config_db []VipConfigDB
		err = rds.SqlxSelect(dbconn.NDB, &arr_vip_config_db, "select id,created_at,updated_at,level,need_exp,withdraw_nums,withdraw_money,day_award_id,day_award_num,week_award_id,week_award_num,month_award_id,month_award_num,cash_bank_rate,cash_bank_top from vip_config")
		if err != nil {
			logger.Logger.Errorf("getVipConfigData: select sql failed! err:[%v]", err)
			return nil, consts.DataFindError.Code
		}
		//转换数据
		for _, val := range arr_vip_config_db {
			item := VipConfigRedis{
				Id:            val.Id,
				Level:         val.Level,
				NeedExp:       val.NeedExp,
				WithdrawNums:  val.WithdrawNums,
				WithdrawMoney: val.WithdrawMoney,
				DayAwardId:    val.DayAwardId,
				DayAwardNum:   val.DayAwardNum,
				WeekAwardId:   val.WeekAwardId,
				WeekAwardNum:  val.WeekAwardNum,
				MonthAwardId:  val.MonthAwardId,
				MonthAwardNum: val.MonthAwardNum,
				CashBankRate:  val.CashBankRate,
				CashBankTop:   val.CashBankTop,
			}
			arrdbinfo = append(arrdbinfo, item)
		}

		//存入redis
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(arrdbinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("getActivityInfo: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", consts.RedisVipConfig, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("getActivityInfo: set config reids failed! err:[%v]", err)
		}
	}

	//给数据排序
	sort.Slice(arrdbinfo, func(i, j int) bool {
		return arrdbinfo[i].NeedExp < arrdbinfo[j].NeedExp
	})

	return arrdbinfo, 0
}

/**根据vip经验，获取vip等级信息*/
func GetVipItem(vipnum int, arrdbinfo []VipConfigRedis) (vip_item *VipConfigRedis) {
	//用二分法查找信息(注意用二分法需要把数组排序，这里在获取的时候就排序过了，所以这里不排序)
	left, right := 0, len(arrdbinfo)-1
	for left <= right {
		mid := (left + right) / 2

		// 当前范围的下界
		currentLow := arrdbinfo[mid].NeedExp
		// 下一个范围的下界
		var nextLow int
		if mid+1 < len(arrdbinfo) {
			nextLow = arrdbinfo[mid+1].NeedExp
		} else {
			nextLow = int(^uint(0) >> 1) // 代表正无穷大
		}

		if vipnum >= currentLow && vipnum < nextLow {
			return &arrdbinfo[mid]
		} else if vipnum < currentLow {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return nil // 如果未找到，返回 -1 表示未找到对应的 level
}

/**获取用户vip功能信息*/
func GetVipUserInfo(uid uint64, uidstr string) (userinfo VipUserInfo, err error) {
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()

	redis_key := consts.RedisVipUser + uidstr
	reply, err := redis.String(redisConn.Do("GET", redis_key))

	userinfo = VipUserInfo{}

	if reply != "" {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &userinfo)
		if err != nil {
			logger.Logger.Errorf("getVipUserInfo: parse config failed! err:[%v]", err)
		}
	}

	if err != nil {
		//获取vip信息
		userInfodb := sqlInfo.UserInfo{}
		err = rds.SqlxGetD(dbconn.NDB, &userInfodb, "select vip,vip_cash_bank1,vip_cash_bank2,vip_cash_bank_time from %s where uid = ?",
			sqlInfo.GetUserInfoTableName(uid), uid)
		if err != nil {
			logger.Logger.Errorf("getVipUserInfo: get user info failed! err:[%v]", err)
			return
		}
		//转换数据
		userinfo.Vip = int(userInfodb.Vip)
		userinfo.VipCashBank1 = userInfodb.VipCashBank1
		userinfo.VipCashBank2 = userInfodb.VipCashBank2
		userinfo.VipCashBankTime = userInfodb.VipCashBankTime.Unix()
		//获取vip 功能信息
		day_geted, week_geted, mon_geted, withdraw_num, withdraw_money, withdraw_time := sql.NullTime{}, sql.NullTime{}, sql.NullTime{}, 0, 0, sql.NullTime{}
		var rows *sqlx.Rows
		rows, err = rds.Query(dbconn.NDB, "select day_geted,week_geted,mon_geted,withdraw_num,withdraw_money,withdraw_time from vip_user_info where uid=%v", uid)
		if err != nil {
			logger.Logger.Errorf("getVipUserInfo: get user vip info failed! err:[%v]", err)
			return
		}
		defer rows.Close()
		time_now := time.Now()
		if rows.Next() {
			if err = rows.Scan(&day_geted, &week_geted, &mon_geted, &withdraw_num, &withdraw_money, &withdraw_time); err != nil {
				logger.Logger.Errorf("getVipUserInfo: prase data failed! err:[%v]", err)
				return
			}
		} else { //没有数据，则插入数据
			insert_columns := map[string]interface{}{
				"uid":        uid,
				"day_geted":  time_now,
				"week_geted": time_now,
				"mon_geted":  time_now,
			}
			_, err = rds.SqlxNamedExec(dbconn.NDB, "insert into vip_user_info(uid,day_geted,week_geted,mon_geted) values(:uid,:day_geted,:week_geted,:mon_geted)", insert_columns)
			if err != nil {
				logger.Logger.Errorf("getVipUserInfo: insert sql failed! err:[%v]", err)
				return
			}
		}
		if day_geted.Valid {
			userinfo.DayGeted = day_geted.Time.Unix()
		} else {
			userinfo.DayGeted = time_now.Unix()
		}
		if week_geted.Valid {
			userinfo.WeekGeted = week_geted.Time.Unix()
		} else {
			userinfo.WeekGeted = time_now.Unix()
		}
		if mon_geted.Valid {
			userinfo.MonGeted = mon_geted.Time.Unix()
		} else {
			userinfo.MonGeted = time_now.Unix()
		}

		userinfo.WithdrawNum = withdraw_num
		userinfo.WithdrawMoney = withdraw_money
		if withdraw_time.Valid {
			userinfo.WithdrawTime = withdraw_time.Time.Unix()
		}
	}

	//判断上次累计的时间是否是今天，如果不是，则要把前天的数据删除
	err = updateVipBankRefresh(uid, &userinfo)
	if err != nil {
		logger.Logger.Errorf("getVipUserInfo: update user info failed! err:[%v]", err)
		return
	}
	if reply == "" {
		//存入redis
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("getVipUserInfo: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("getVipUserInfo: set config reids failed! err:[%v]", err)
		}
	}

	return
}

/**跨天更新用户vip存钱罐信息*/
func updateVipBankRefresh(uid uint64, userinfo *VipUserInfo) (err error) {
	time_now := time.Now()
	day_diff := tool.TimeDaysDiff(time.Unix(userinfo.VipCashBankTime, 0), time_now)
	if (userinfo.VipCashBank1 > 0 || userinfo.VipCashBank2 > 0) && day_diff != 0 { //处理vip存钱罐信息
		update_columns := map[string]interface{}{
			"uid":                uid,
			"vip_cash_bank_time": time_now,
		}
		update_sql := ""
		if day_diff == 1 { //相差一天，则把昨日累计放入可提取
			update_sql = "update " + sqlInfo.GetUserInfoTableName(uid) + " set vip_cash_bank1=vip_cash_bank2,vip_cash_bank2=0,vip_cash_bank_time=:vip_cash_bank_time where uid=:uid"
			userinfo.VipCashBank1 = userinfo.VipCashBank2
			userinfo.VipCashBank2 = 0
			if userinfo.VipCashBank1 > 0 {
				AddRedDot(uid, strconv.FormatUint(uid, 10), consts.RedDot_VIP, 1, false, false)
			}
		} else { //相差不止一天，则全都清空
			update_sql = "update " + sqlInfo.GetUserInfoTableName(uid) + " set vip_cash_bank1=0,vip_cash_bank2=0,vip_cash_bank_time=:vip_cash_bank_time where uid=:uid"
			userinfo.VipCashBank1 = 0
			userinfo.VipCashBank2 = 0
			AddRedDot(uid, strconv.FormatUint(uid, 10), consts.RedDot_VIP, 0, true, false)
		}
		_, err = rds.SqlxNamedExec(dbconn.NDB, update_sql, update_columns)
		if err != nil {
			logger.Logger.Errorf("updateVipCashBank: update vip_cash_bank1 failed!uid:[%v], err:[%v]", uid, err)
			return
		}

		userinfo.VipCashBankTime = time_now.Unix()
	}

	//提现已经跨天了，则清空原提现数据
	if userinfo.WithdrawTime > 0 && (userinfo.WithdrawNum > 0 || userinfo.WithdrawMoney > 0) && tool.TimeDaysDiff(time.Unix(userinfo.WithdrawTime, 0), time_now) != 0 {
		update_columns := map[string]interface{}{
			"uid":            uid,
			"withdraw_num":   0,
			"withdraw_money": 0,
			"withdraw_time":  time_now,
		}
		_, err = rds.SqlxNamedExec(dbconn.NDB, "update vip_user_info set withdraw_num=:withdraw_num,withdraw_money=:withdraw_money,withdraw_time=:withdraw_time where uid=:uid", update_columns)
		if err != nil {
			logger.Logger.Errorf("updateVipCashBank: update sql failed! err:[%v]", err)
			return
		}
		userinfo.WithdrawNum = 0
		userinfo.WithdrawMoney = 0
		userinfo.WithdrawTime = time_now.Unix()
	}
	return
}

/**修改vip提现次数和提现金额
*@param uid: 用户id
*@param withdraw_num: 提现次数
*@param withdraw_money: 提现额度*/
func UpdateVipWithdrawInfo(uid uint64, withdraw_num, withdraw_money int) (err error) {
	update_columns := map[string]interface{}{
		"uid":            uid,
		"withdraw_num":   withdraw_num,
		"withdraw_money": withdraw_money,
	}
	updata_sql_patch := ""
	if withdraw_num > 0 || withdraw_money > 0 { //如果提现增加，则记录时间
		update_columns["withdraw_time"] = time.Now()
		updata_sql_patch = ",withdraw_time=:withdraw_time"
	}
	update_sql := fmt.Sprintf("update vip_user_info set withdraw_num=withdraw_num+:withdraw_num,withdraw_money=withdraw_money+:withdraw_money%s where uid=:uid and withdraw_num+withdraw_num>=0 and withdraw_money+withdraw_money>=0", updata_sql_patch)
	_, err = rds.SqlxNamedExec(dbconn.NDB, update_sql, update_columns)
	if err != nil {
		logger.Logger.Errorf("updateVipWithdrawInfo: update sql failed! err:[%v],sql:[%v]", err, update_sql)
		return
	}

	//删除缓存
	//删除redis
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()

	redis_key := consts.RedisVipUser + strconv.FormatUint(uid, 10)
	_, err = redisConn.Do("DEL", redis_key)
	if err != nil {
		logger.Logger.Errorf("updateVipWithdrawInfo: delete redis failed! uid:[%v],err:[%v]", uid, err)
	}
	return
}
