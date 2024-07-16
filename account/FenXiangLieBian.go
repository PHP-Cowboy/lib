/**分享裂变*/

package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/gomodule/redigo/redis"
	"za.game/lib/consts"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
	"za.game/lib/sqlInfo"
	"za.game/lib/tool"
)

/************************************** class ***********************************************************/

/**分享裂变等级，数据表结构*/
type FxlbLevelConfigDB struct {
	Level         int     `db:"level" gorm:"primarykey;column:level;type:int(11);not null;comment:等级;" json:"level" `
	NeedExp       int64   `db:"need_exp" gorm:"column:need_exp;type:bigint(20);not null;default:0;comment:升级所需经验;" json:"need_exp" `
	YqYj          int     `db:"yq_yj" gorm:"column:yq_yj;type:int(11);not null;default:0;comment:邀请佣金/人;" json:"yq_yj" `
	CzYjRate      float64 `db:"cz_yj_rate" gorm:"column:cz_yj_rate;type:decimal(4,4);not null;default:0.0000;comment:充值佣金比例;" json:"cz_yj_rate" `
	DlYjRate      float64 `db:"dl_yj_rate" gorm:"column:dl_yj_rate;type:decimal(4,4);not null;default:0.0000;comment:代理佣金比例;" json:"dl_yj_rate" `
	Layer         int     `db:"layer" gorm:"column:layer;type:int(11);not null;default:0;comment:层级计算;" json:"layer" `
	WithdrawNum   int     `db:"withdraw_num" gorm:"column:withdraw_num;type:int(11);not null;default:0;comment:每日提现次数;" json:"withdraw_num" `
	WithdrawMoney int     `db:"withdraw_money" gorm:"column:withdraw_money;type:int(11);not null;default:0;comment:每日提现金额;" json:"withdraw_money" `
}

/**分享裂变，用户绑定信息*/
type FxlbUserBindData struct {
	From_uid  uint64 //上级id
	From_code string //上级邀请码
}

/**分享裂变，用户提现信息*/
type FxlbUserWithdrawData struct {
	Num   int   //今日提现次数
	Time  int64 //提现时间
	Money int   //今日提现额度
}

/**分享裂变，公共配置*/
type FxlbCommonConfig struct {
	Is_open    int     //是否开启
	Yqyx_user  int     //邀请有效人数 >= N 才能提现
	Yj_zh_rate float64 //佣金转换比例
}

/************************************** public***********************************************************/

/**获取分享裂变开关以及通用配置信息*/
func GetFXLBCommonConfig() (err error, map_fxlb *FxlbCommonConfig) {
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()
	redis_str := consts.FXLB_dict
	reply, err := redis.String(redisConn.Do("GET", redis_str))

	map_fxlb = &FxlbCommonConfig{}

	if err == nil {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &map_fxlb)
		if err != nil {
			logger.Logger.Debugf("GetFXLBCommonConfig: 分享裂变-解析json失败！错误:[%v]", err)
		}
	}

	//读取redis失败，则从数据库获取
	if err != nil {
		//查询数据
		code, value := "", ""
		rows, err := rds.QueryX(dbconn.GameDB, "select code,value from dict where type_code = 'fxlb'")
		if err != nil {
			logger.Logger.Debugf("GetFXLBCommonConfig: 分享裂变-获取数据表 dict 失败！错误:[%v]", err)
			return err, nil
		}
		defer rows.Close()
		for rows.Next() {
			if err := rows.Scan(&code, &value); err != nil {
				logger.Logger.Errorf("GetFXLBCommonConfig: 分享裂变-解析数据表 dict 失败! err:[%v]", err)
				return err, nil
			}
			if code == "is_open" {
				map_fxlb.Is_open, _ = strconv.Atoi(value)
			} else if code == "yqyx_user" {
				map_fxlb.Yqyx_user, _ = strconv.Atoi(value)
			} else if code == "yj_zh_rate" {
				map_fxlb.Yj_zh_rate, _ = strconv.ParseFloat(value, 64)
			}
		}

		jsonTxt := ""

		jsonTxt, err = tool.JsonString(map_fxlb)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Debugf("GetFXLBCommonConfig: 分享裂变-转为json失败！错误:[%v]", err)
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", redis_str, jsonTxt, "EX", tool.GetRedisExpireTime())
		if err != nil {
			logger.Logger.Debugf("GetFXLBCommonConfig: 分享裂变-设置缓存数据失败！错误:[%v]", err)
		}
	}
	return nil, map_fxlb
}

/**获得用户分享裂变绑定信息*/
func GetUserFxlbBindInfo(uid uint64) (err error, info FxlbUserBindData) {
	//由于分享裂变的数据，随着下级而变动，可能会变动频繁，所以这里不做redis缓存
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()
	redis_str := consts.Fxlb_bind_info + strconv.Itoa(int(uid))
	reply, err := redis.String(redisConn.Do("GET", redis_str))

	info = FxlbUserBindData{}

	if err == nil {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &info)
		if err != nil {
			logger.Logger.Debugf("FenXiangLieBianServices: 分享裂变-解析redis绑定信息失败！错误:[%v]", err)
		}
	}

	//读取redis失败，则从数据库获取
	if err != nil {
		//查询字典数据
		from_uid, from_code := uint64(0), ""
		rows, err := rds.QueryX(dbconn.NDB, "select from_uid,from_code from fxlb_user_relation where uid=?", uid)
		if err != nil {
			logger.Logger.Errorf("FenXiangLieBianServices:  分享裂变-获取数据库绑定信息失败! err:[%v]", err)
			return err, info
		}
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&from_uid, &from_code); err != nil {
				logger.Logger.Errorf("FenXiangLieBianServices: 分享裂变-解析数据库绑定信息失败! err:[%v]", err)
				return err, info
			}
		}
		info.From_uid = from_uid
		info.From_code = from_code
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(info)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Debugf("FenXiangLieBianServices: 分享裂变-绑定信息转为json失败！错误:[%v]", err)
			return nil, info
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", redis_str, jsonTxt, "EX", tool.GetRedisExpireTime())
		if err != nil {
			logger.Logger.Debugf("FenXiangLieBianServices: 分享裂变-绑定信息设置缓存数据失败！错误:[%v]", err)
		}
	}
	return nil, info
}

/**获得分享裂变等级配置信息*/
func GetFxlbLevelConfig() (err error, info []FxlbLevelConfigDB) {
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()
	redis_str := consts.Fxlb_level_config
	reply, err := redis.String(redisConn.Do("GET", redis_str))

	info = []FxlbLevelConfigDB{}

	if err == nil {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &info)
		if err != nil {
			logger.Logger.Debugf("(s *InitServices) getFxlbLevelConfig: 分享裂变-解析json失败！错误:[%v]", err)
		}
	}

	//读取redis失败，则从数据库获取
	if err != nil {
		//查询字典数据
		err = rds.SqlxSelect(dbconn.NDB, &info, "select level,need_exp,yq_yj,cz_yj_rate,dl_yj_rate,layer,withdraw_num,withdraw_money from fxlb_level_config")
		if err != nil {
			logger.Logger.Errorf("(s *InitServices) getFxlbLevelConfig: 分享裂变-获取等级配置失败! err:[%v]", err)
			return
		}

		jsonTxt := ""

		jsonTxt, err = tool.JsonString(info)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Debugf("(s *InitServices) getFxlbLevelConfig: 分享裂变-转为json失败！错误:[%v]", err)
			return nil, info
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", redis_str, jsonTxt, "EX", tool.GetRedisExpireTime())
		if err != nil {
			logger.Logger.Debugf("(s *InitServices) getFxlbLevelConfig: 分享裂变-设置缓存数据失败！错误:[%v]", err)
		}
	}
	return nil, info
}

/**获取用户分享裂变经验值*/
func GetUserFxlbJyValue(uid uint64) (err error, value int64) {
	//由于分享裂变的数据，随着下级而变动，可能会变动频繁，所以这里不做redis缓存
	fxlb_jy := int64(0)
	select_sql := fmt.Sprintf("select fxlb_jy from %s where uid=?", sqlInfo.GetUserInfoTableName(uid))
	rows, err := rds.QueryX(dbconn.NDB, select_sql, uid)
	if err != nil {
		logger.Logger.Errorf("getUserFxlbJyValue: 分享裂变-获取用户信息失败! err:[%v],uid:[%v]", err, uid)
		return
	}
	defer rows.Close()
	if rows.Next() {
		if err = rows.Scan(&fxlb_jy); err != nil {
			logger.Logger.Errorf("getUserFxlbJyValue: 分享裂变-解析用户数据失败! err:[%v],uid:[%v]", err, uid)
			return
		}
	}
	value = fxlb_jy
	return
}

/**根据vip经验，获取vip等级信息*/
func GetFxlbVipItem(vipnum int64, arrdbinfo []FxlbLevelConfigDB) (vip_item *FxlbLevelConfigDB) {
	//用二分法查找信息(注意用二分法需要把数组排序，这里在获取的时候就排序过了，所以这里不排序)
	left, right := 0, len(arrdbinfo)-1
	for left <= right {
		mid := (left + right) / 2

		// 当前范围的下界
		currentLow := arrdbinfo[mid].NeedExp
		// 下一个范围的下界
		var nextLow int64
		if mid+1 < len(arrdbinfo) {
			nextLow = arrdbinfo[mid+1].NeedExp
		} else {
			nextLow = int64(^uint(0) >> 1) // 代表正无穷大
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

/**获得用户分享裂变提现信息*/
func GetUserFxlbWithdrawInfo(uid uint64) (err error, info FxlbUserWithdrawData) {
	//由于分享裂变的数据，随着下级而变动，可能会变动频繁，所以这里不做redis缓存
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()
	redis_str := consts.Fxlb_withdraw_info + strconv.Itoa(int(uid))
	reply, err := redis.String(redisConn.Do("GET", redis_str))

	info = FxlbUserWithdrawData{}

	if err == nil {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &info)
		if err != nil {
			logger.Logger.Debugf("GetUserFxlbWithdrawInfo: 分享裂变-解析redis提现信息失败！错误:[%v],uid:[%v]", err, uid)
		}
	}

	//读取redis失败，则从数据库获取
	if err != nil {
		//查询数据
		withdraw_num, withdraw_time, withdraw_money := 0, sql.NullTime{}, 0
		rows, err := rds.QueryX(dbconn.NDB, "select withdraw_num,withdraw_time,withdraw_money from fxlb_user_relation where uid=?", uid)
		if err != nil {
			logger.Logger.Errorf("GetUserFxlbWithdrawInfo:  分享裂变-获取数据库提现信息失败! err:[%v],uid:[%v]", err, uid)
			return err, info
		}
		defer rows.Close()
		if rows.Next() {
			if err := rows.Scan(&withdraw_num, &withdraw_time, &withdraw_money); err != nil {
				logger.Logger.Errorf("GetUserFxlbWithdrawInfo: 分享裂变-解析数据库提现信息失败! err:[%v],uid:[%v]", err, uid)
				return err, info
			}
		}
		info.Num = withdraw_num
		info.Money = withdraw_money
		time_now := time.Now()
		if withdraw_time.Valid {
			if tool.TimeDaysDiff(time_now, withdraw_time.Time) != 0 {
				//跨天了，重置提现次数
				info.Num = 0
				info.Money = 0
				info.Time = time_now.Unix()
			} else {
				info.Time = withdraw_time.Time.Unix()
			}
		} else {
			info.Time = time.Now().Unix()
		}
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(info)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Debugf("GetUserFxlbWithdrawInfo: 分享裂变-提现信息转为json失败！错误:[%v],uid:[%v]", err, uid)
			return nil, info
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", redis_str, jsonTxt, "EX", tool.GetRedisExpireTime())
		if err != nil {
			logger.Logger.Debugf("GetUserFxlbWithdrawInfo: 分享裂变-提现信息设置缓存数据失败！错误:[%v],uid:[%v]", err, uid)
		}
	}
	return nil, info
}
