package rds

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	uiltLog "za.game/za.log"
)

var sqlLog *uiltLog.ZapLog

func InitDB(user, pwd, host, port, dbName string, min, max, lifeTime int) (*gorm.DB, error) {
	var err error
	dsn := user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8&parseTime=true&loc=Local"

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxOpenConns(max)
	sqlDB.SetMaxIdleConns(min)
	sqlDB.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second)
	return DB, nil
}

func InitSqlDB(user, pwd, host, port, dbName string, min, max, lifeTime int) (*sqlx.DB, error) {
	var err error
	dsn := user + ":" + pwd + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8&parseTime=true&loc=Local"
	sqlDB, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(max)                                      // 设置最大的并发连接数（in-use + idle）
	sqlDB.SetMaxIdleConns(min)                                      // 设置最大的空闲连接数（idle）
	sqlDB.SetConnMaxLifetime(time.Duration(lifeTime) * time.Second) // 设置连接的最大生命周期
	return sqlDB, nil
}

// const sqlLogDebug = true
//
// // sql 执行日志输出
//
//	func SqlLogOutput(db any, query string, errlog error, values ...any) {
//		if sqlLogDebug {
//			var (
//				interpolateQueryFunc func(string, ...interface{}) (string, error)
//				//dbType               string
//			)
//
//			switch v := db.(type) {
//			case *sql.DB:
//				interpolateQueryFunc = v.InterpolateQuery
//				//dbType = "*sql.DB"
//			case *sql.Tx:
//				interpolateQueryFunc = v.InterpolateQuery
//				//dbType = "*sql.Tx"
//			default:
//
//			}
//			interpolatedQuery, err := interpolateQueryFunc(query, values...)
//			sqlErrLog(interpolatedQuery, err, values...)
//
//		} else {
//			sqlErrLog(query, errlog, values...)
//		}
//	}
func sqlErrLog(query string, errlog error, values ...any) {
	if sqlLog == nil {
		sqlLog = &uiltLog.ZapLog{}
		sqlLog.InitDBWarnLog("./txtlog/")
	}
	modName := "sqlx"
	errLog := fmt.Sprint(errlog.Error(), "| ", query, values)
	uiltLog.DBWarnString(modName, "warn", errLog)
}
func ZapErrLog(modName, lv string, values ...any) {
	if sqlLog == nil {
		sqlLog = &uiltLog.ZapLog{}
		sqlLog.InitDBWarnLog("./txtlog/")
	}
	errLog := fmt.Sprint(values)
	uiltLog.DBWarnString(modName, lv, errLog)
}
func RedisErrLog(values ...any) {
	ZapErrLog("redis", "warn", values...)
}

// 大量字段insert的时候用
func InsertByStruct(db *sqlx.DB, table string, value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %T", value)
	}

	typ := val.Type()

	var cols []string
	var valCols []string
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		col := field.Name

		// 检查字段标签，跳过带有 db:"-" 的字段
		if dbTag := field.Tag.Get("insert_db"); dbTag == "-" {
			continue
		}
		dbTag := field.Tag.Get("db")
		if dbTag == "-" {
			continue
		}
		if dbTag != "" {
			col = dbTag
		}

		cols = append(cols, col)
		valCols = append(valCols, ":"+col)
	}

	insertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table,
		strings.Join(cols, ","), strings.Join(valCols, ","))

	_, err := SqlxNamedExec(db, insertStmt, value)
	return err
}

// 大量字段insert的时候用
func InsertByStructTx(tx *sqlx.Tx, table string, value any) error {
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %T", value)
	}

	typ := val.Type()

	var cols []string
	var placeholders []string
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		col := field.Name

		// 检查字段标签，跳过带有 db:"-" 的字段
		if dbTag := field.Tag.Get("db"); dbTag == "-" {
			continue
		}

		cols = append(cols, col)
		placeholders = append(placeholders, "?")
	}

	insertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table,
		strings.Join(cols, ", "), strings.Join(placeholders, ", "))

	_, err := tx.Exec(insertStmt, value)
	return err
}

// 大量字段insert的时候用
func InsertBySliceTx(tx *sql.Tx, table string, values ...any) error {
	for value := range values {
		val := reflect.ValueOf(value)
		if val.Kind() != reflect.Struct {
			return fmt.Errorf("expected a struct, got %T", value)
		}

		typ := val.Type()

		var cols []string
		var placeholders []string
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			col := field.Name

			// 检查字段标签，跳过带有 db:"-" 的字段
			if dbTag := field.Tag.Get("db"); dbTag == "-" {
				continue
			}

			cols = append(cols, col)
			placeholders = append(placeholders, "?")
		}

		insertStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table,
			strings.Join(cols, ", "), strings.Join(placeholders, ", "))

		_, err := tx.Exec(insertStmt, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func SqlxGet(db *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	err := db.Get(dest, query, args...)
	if err != nil {
		sqlErrLog(query, err, args)
	}
	return err
}

// 不记录sql未查到记录错误
func SqlxGetForErrNoRows(db *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	err := db.Get(dest, query, args...)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			sqlErrLog(query, err, args)
		}
	}
	return err
}

func SqlxSelect(db *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	err := db.Select(dest, query, args...)
	if err != nil {
		sqlErrLog(query, err, args)
	}
	return err
}

// 不记录sql未查到记录错误
func SqlxSelectForErrNoRows(db *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	err := db.Select(dest, query, args...)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			sqlErrLog(query, err, args)
		}
	}
	return err
}
func SqlxExec(db *sqlx.DB, query string, args ...any) (sql.Result, error) {
	t, err := db.Exec(query, args...)
	if err != nil {
		sqlErrLog(query, err, args)
	}
	return t, err
}

func SqlxNamedExec(db *sqlx.DB, query string, arg interface{}) (sql.Result, error) {
	t, err := db.NamedExec(query, arg)
	if err != nil {
		sqlErrLog(query, err, arg)
	}
	return t, err
}

// tx
func SqlxTxGet(db *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	err := db.Get(dest, query, args...)
	if err != nil {
		sqlErrLog(query, err, args)
	}
	return err
}

func SqlxTxSelect(db *sqlx.Tx, dest interface{}, query string, args ...interface{}) error {
	err := db.Select(dest, query, args...)
	if err != nil {
		sqlErrLog(query, err, args)
	}
	return err
}
func SqlxExecTx(db *sqlx.Tx, query string, args ...any) (sql.Result, error) {
	t, err := db.Exec(query, args...)
	if err != nil {
		sqlErrLog(query, err)
	}
	return t, err
}

func SqlxTxExec(tx *sql.Tx, query string, args ...any) (sql.Result, error) {
	t, err := tx.Exec(query, args...)
	if err != nil {
		sqlErrLog(query, err, args)
	}
	return t, err
}
func SqlxNamedExecTx(db *sqlx.Tx, query string, arg interface{}) (sql.Result, error) {
	t, err := db.NamedExec(query, arg)
	if err != nil {
		sqlErrLog(query, err)
	}
	return t, err
}

// 动态表名,因为有预处理 需要表名
func SqlxGetD(db *sqlx.DB, dest interface{}, query string, table string, args ...interface{}) error {
	sqlStr := fmt.Sprintf(query, table)
	err := db.Get(dest, sqlStr, args...)
	if err != nil {
		sqlErrLog(sqlStr, err, args)
	}
	return err
}
func SqlxSelectD(db *sqlx.DB, dest interface{}, query string, table string, args ...interface{}) error {
	sqlStr := fmt.Sprintf(query, table)
	err := db.Select(dest, sqlStr, args...)
	if err != nil {
		sqlErrLog(sqlStr, err, args)
	}
	return err
}
func SqlxExecD(db *sqlx.DB, query string, table string, args ...any) (sql.Result, error) {
	sqlStr := fmt.Sprintf(query, table)
	t, err := db.Exec(sqlStr, args...)
	if err != nil {
		sqlErrLog(sqlStr, err, args)
	}
	return t, err
}

func SqlxNamedExecD(db *sqlx.DB, query string, table string, arg interface{}) (sql.Result, error) {
	sqlStr := fmt.Sprintf(query, table)
	t, err := db.NamedExec(sqlStr, arg)
	if err != nil {
		sqlErrLog(sqlStr, err, arg)
	}
	return t, err
}

// tx
func SqlxTxGetD(db *sqlx.Tx, dest interface{}, query string, table string, args ...interface{}) error {
	sqlStr := fmt.Sprintf(query, table)
	err := db.Get(dest, sqlStr, args...)
	if err != nil {
		sqlErrLog(sqlStr, err, args)
	}
	return err
}
func SqlxTxSelectD(db *sqlx.Tx, dest interface{}, query string, table string, args ...interface{}) error {
	sqlStr := fmt.Sprintf(query, table)
	err := db.Select(dest, sqlStr, args...)
	if err != nil {
		sqlErrLog(sqlStr, err, args)
	}
	return err
}
func SqlxTxExecD(tx *sqlx.Tx, query string, table string, args ...any) (sql.Result, error) {
	sqlStr := fmt.Sprintf(query, table)
	t, err := tx.Exec(sqlStr, args...)
	if err != nil {
		sqlErrLog(sqlStr, err, args)
	}
	return t, err
}

/**查询数据，返回数据*/
func Query(db *sqlx.DB, format string, args ...interface{}) (*sqlx.Rows, error) {
	query := fmt.Sprintf(format, args...)
	// 执行查询语句
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

/**参数化查询语句*/
func QueryX(db *sqlx.DB, sql string, args ...interface{}) (*sqlx.Rows, error) {
	// 执行查询语句
	rows, err := db.Queryx(sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
