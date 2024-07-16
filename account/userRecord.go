package account

import (
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"time"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
)

type UserRecordLog struct {
	Uid        uint64         `db:"uid"  json:"uid" `
	RecordData sql.NullString `db:"record_data" json:"record_data" `
	UpdatedDt  time.Time      `db:"updated_at" json:"updated_dt" `
}

func (t *UserRecordLog) TableName() string {
	return "user_record_log"
}

type UserRecordItem struct {
	Cash          int       `json:"c" ` //win+cash
	Bonus         int       `json:"b" `
	Type          int       `json:"p" ` //类型 如果大于1000，则是游戏流水，千位以上是游戏id
	recordTime    time.Time `json:"-" `
	RecordTimeStr string    `json:"t" `
}
type UserRecord struct {
	UserID       uint64
	RecentRecord []UserRecordItem
}

var RecordMax = 100

// 获取用户属性信息
func GetUserRecordLog(userID uint64) (string, error) {
	u := UserRecordLog{}
	err := rds.SqlxGetForErrNoRows(dbconn.LogDB, &u, "SELECT record_data FROM user_record_log where uid=?", userID)
	if err != nil {
		return "", err
	}
	if u.RecordData.Valid {
		return u.RecordData.String, nil
	}
	return "", nil
}

// 获取用户属性信息
func GetUserPropertiesFromDB(userID uint64) (*UserRecord, error) {
	u := UserRecordLog{}
	err := rds.SqlxGetForErrNoRows(dbconn.LogDB, &u, "SELECT record_data FROM user_record_log where uid=?", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if _, err = rds.SqlxExec(dbconn.LogDB, "insert into user_record_log(uid) values (?)", userID); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	// 解析最近交易记录
	var recentRecord []UserRecordItem
	if u.RecordData.Valid {
		if err = json.Unmarshal([]byte(u.RecordData.String), &recentRecord); err != nil {
			return nil, err
		}
	}
	for i := 0; i < len(recentRecord); i++ {
		if recordTime, errT := time.Parse("2006-01-02 15:04:05", recentRecord[i].RecordTimeStr); errT == nil {
			recentRecord[i].recordTime = recordTime
		}
	}
	// 构造用户属性对象
	userProps := &UserRecord{
		UserID:       userID,
		RecentRecord: recentRecord,
	}
	return userProps, nil
}

// 更新用户最近100条流水
func UpdateRecentRecord(userID uint64, newTransaction ...UserRecordItem) error {
	// 检查用户是否存在
	userProps, err := GetUserPropertiesFromDB(userID)
	if err != nil {
		return err
	}
	for k, v := range newTransaction {
		if v.recordTime.IsZero() {
			if recordTime, errT := time.Parse("2006-01-02 15:04:05", newTransaction[k].RecordTimeStr); errT == nil {
				newTransaction[k].recordTime = recordTime
			}
		}
	}
	// 添加新交易到流水中
	userProps.RecentRecord = append(userProps.RecentRecord, newTransaction...)

	// 保留最新的100条流水，按时间戳排序
	sort.Slice(userProps.RecentRecord, func(i, j int) bool {
		return userProps.RecentRecord[i].recordTime.After(userProps.RecentRecord[j].recordTime)
	})
	if len(userProps.RecentRecord) > RecordMax {
		userProps.RecentRecord = userProps.RecentRecord[:RecordMax]
	}

	if b, errJ := json.Marshal(userProps.RecentRecord); errJ == nil {
		_, errJ = rds.SqlxExec(dbconn.LogDB, "update user_record_log set record_data = ?,updated_dt=CURRENT_TIMESTAMP where uid= ?", string(b), userID)
		if errJ != nil {
			return errJ
		}
	} else {
		return errJ
	}

	return nil
}
