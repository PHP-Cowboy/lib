package game

import (
	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/jmoiron/sqlx"
	"log"
	"za.game/lib/sqlInfo"
)

type Roomview struct {
	sqlInfo.RoomConfig
}

type RoomViewBind struct {
	Field   string `json:"field"`
	Comment string `json:"comment"`
}

type ColumnComment struct {
	ColumnName    string `db:"COLUMN_NAME" json:"column_name"`
	ColumnComment string `db:"COLUMN_COMMENT" json:"column_comment"`
}

var RoomviewSlice = []string{
	//"Svrid",
	"OnlineMax", "OnlineMin",
	"id", "GameId", "RoomId", "RoomIndex", "Base", "MinEntry", "MaxEntry", "RoomName", "RoomType", "RoomSwitch", "RoomWelfare", "RoomWelfareDesc",
	"Desc", "Tax", "BonusDiscount", "AiSwitch", "AiLimit", "ExtData", "PoolID", "RechargeLimit", "pExtData", "RechargeCount", "CurPoolValue",
}

func GetColumnList(db *sqlx.DB) (list []string, err error) {

	// 查询表的字段名
	query := `SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`

	err = db.Select(&list, query, "game", "roomview")

	if err != nil {
		logger.Logger.Errorf("查询字段名时出错: %v", err)
	}

	return
}

func GetRoomViewColumnComment(db *sqlx.DB) (dataList []ColumnComment, err error) {
	// 查询指定表的字段及其注释
	tableName := "roomview"
	query := `SELECT COLUMN_NAME, COLUMN_COMMENT FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? `

	// 执行查询
	err = db.Select(&dataList, query, tableName)
	if err != nil {
		logger.Logger.Errorf("Query failed err:[%v]", err)
		return
	}

	return
}

func GetRoomViewBindList(db *sqlx.DB) (list []RoomViewBind, err error) {
	// 查询表的字段名
	query := `SHOW FULL COLUMNS FROM roomview`

	rows, err := db.Queryx(query)

	if err != nil {
		logger.Logger.Errorf("SHOW FULL COLUMNS err:[%v]", err)
		return
	}

	defer rows.Close()

	// 获取列名
	columnNames, err := rows.Columns()
	if err != nil {
		log.Fatalf("获取列名出错: %v", err)
	}

	// 创建一个 map 切片来存储结果
	var results []map[string]interface{}

	// 遍历每一行
	for rows.Next() {
		// 创建一个映射来存储当前行的数据
		rowData := make(map[string]interface{}, len(columnNames))

		// 创建一个用于 Scan 的空接口切片
		scanArgs := make([]interface{}, len(columnNames))
		for i := range scanArgs {
			var v interface{}
			scanArgs[i] = &v
		}

		// 使用 Scan 读取当前行的值
		if err = rows.Scan(scanArgs...); err != nil {
			log.Fatalf("读取行数据出错: %v", err)
		}

		// 将读取的值放入映射中
		for i, value := range scanArgs {
			rowData[columnNames[i]] = *value.(*interface{})
		}

		// 将映射添加到结果切片中
		results = append(results, rowData)
	}

	// 检查是否有错误发生（例如，在遍历结束后）
	if err = rows.Err(); err != nil {
		log.Fatalf("读取数据出错: %v", err)
	}

	list = make([]RoomViewBind, 0, len(results))
	// 输出结果
	for _, r := range results {
		for k, v := range r {
			if k != "Field" && k != "Comment" {
				continue
			}

			tmp := RoomViewBind{Field: k}

			switch v.(type) {
			case []uint8:
				tmp.Comment = string(v.([]uint8))
			}

			list = append(list, tmp)
		}
	}

	return
}
