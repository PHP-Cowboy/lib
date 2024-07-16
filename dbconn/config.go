package dbconn

import (
	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/BurntSushi/toml"
)

type Config struct {
	MysqlGame struct {
		DbName string `toml:"dbName"`
		User   string `toml:"user"`
		Passwd string `toml:"passwd"`
		Port   string `toml:"port"`
		Host   string `toml:"host"`
	} `toml:"MysqlGame"`

	MysqlNdb struct {
		DbName string `toml:"dbName"`
		User   string `toml:"user"`
		Passwd string `toml:"passwd"`
		Port   string `toml:"port"`
		Host   string `toml:"host"`
	} `toml:"MysqlNdb"`

	MysqlLog struct {
		DbName string `toml:"dbName"`
		User   string `toml:"user"`
		Passwd string `toml:"passwd"`
		Port   string `toml:"port"`
		Host   string `toml:"host"`
	} `toml:"MysqlLog"`

	MysqlPay struct {
		DbName string `toml:"dbName"`
		User   string `toml:"user"`
		Passwd string `toml:"passwd"`
		Port   string `toml:"port"`
		Host   string `toml:"host"`
	} `toml:"MysqlPay"`

	Redis struct {
		Host string `toml:"host"`
		Port string `toml:"port"`
		User string `toml:"user"`
		Pwd  string `toml:"pwd"`
		DB   int    `toml:"DB"`
	} `toml:"redis"`
}

func InitConfig(config_path string) {
	// if _, err := toml.DecodeFile("./dbconfig.toml", &CONFIG); err != nil {
	// 	logger.Logger.Errorf("decode" + err.Error())
	// }
	if _, err := toml.DecodeFile(config_path, &CONFIG); err != nil {
		logger.Logger.Errorf("InitConfig: decode config falied! err:[%v],path:[%v]", err, config_path)
	}
}
