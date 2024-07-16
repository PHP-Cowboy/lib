package tool

import (
	"gorm.io/gorm"
)

type SPResult struct {
	LeftValue   int `gorm:"column:@leftValue"`
	SuccessFlag int `gorm:"column:@successFlag"`
}
type Result struct {
	Value int
}

func AddAmount(db *gorm.DB, userID int64, addVal uint) (leftVal SPResult, e error) {
	var spResult SPResult
	var result Result
	// 调用存储过程
	if err := db.Raw("CALL AddUserAmount(?, ?, @leftValue, @successFlag)", userID, addVal).Scan(&result).Error; err != nil {
		return spResult, err
	}
	// 调用存储过程并扫描输出参数
	if err := db.Raw("SELECT @leftValue, @successFlag").Scan(&spResult).Error; err != nil {
		return spResult, err
	}
	//todo 写一个异步日志

	return spResult, nil
}
