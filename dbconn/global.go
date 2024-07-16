package dbconn

import (
	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"za.game/lib/rds"
)

var (
	Err       error
	CONFIG    *Config
	RedisPool *redis.Pool
	GameDB    *sqlx.DB
	NDB       *sqlx.DB
	PayDB     *sqlx.DB
	LogDB     *sqlx.DB
)

func InitApp() {

	//初始化配置文件
	InitConfig("./dbconfig.toml")

	if CONFIG == nil {
		logger.Error("dbconn::InitApp: no config dbconfig.toml")
		return
	}

	//初始化游戏数据库

	GameDB, Err = rds.InitSqlDB(CONFIG.MysqlGame.User, CONFIG.MysqlGame.Passwd, CONFIG.MysqlGame.Host, CONFIG.MysqlGame.Port, CONFIG.MysqlGame.DbName, 30, 300, 59)
	if Err != nil {
		logger.Logger.Error("InitApp: init GameDB failed! err:[%v]", Err)
	} else {
		logger.Logger.Info("InitApp: init GameDB Success!")
	}

	//初始用户化数据库
	NDB, Err = rds.InitSqlDB(CONFIG.MysqlNdb.User, CONFIG.MysqlNdb.Passwd, CONFIG.MysqlNdb.Host, CONFIG.MysqlNdb.Port, CONFIG.MysqlNdb.DbName, 30, 300, 59)
	if Err != nil {
		logger.Logger.Error("InitApp: init NDB failed! err:[%v]", Err)
	} else {
		logger.Logger.Info("InitApp: init NDB Success!")
	}

	//初始化支付数据库
	PayDB, Err = rds.InitSqlDB(CONFIG.MysqlPay.User, CONFIG.MysqlPay.Passwd, CONFIG.MysqlPay.Host, CONFIG.MysqlPay.Port, CONFIG.MysqlPay.DbName, 30, 300, 59)
	if Err != nil {
		logger.Logger.Error("InitApp: init PayDB failed! err:[%v]", Err)
	} else {
		logger.Logger.Info("InitApp: init PayDB Success!")
	}

	//初始化日志数据库
	LogDB, Err = rds.InitSqlDB(CONFIG.MysqlLog.User, CONFIG.MysqlLog.Passwd, CONFIG.MysqlLog.Host, CONFIG.MysqlLog.Port, CONFIG.MysqlLog.DbName, 30, 300, 59)
	if Err != nil {
		logger.Logger.Error("InitApp: init PayDB failed! err:[%v]", Err)
	} else {
		logger.Logger.Info("InitApp: init PayDB Success!")
	}

	//初始化redis
	RedisPool, Err = rds.InitRedis(CONFIG.Redis.User, CONFIG.Redis.Pwd, CONFIG.Redis.Host, CONFIG.Redis.Port, CONFIG.Redis.DB, 4, 50, 2000)
	if Err != nil {
		logger.Logger.Error("InitApp: init RedisPool failed! err:[%v]", Err)
	} else {
		logger.Logger.Info("InitApp: init RedisPool Success!")
	}
}
