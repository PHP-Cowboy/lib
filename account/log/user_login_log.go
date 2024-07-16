package log

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
	"za.game/lib/rds"
	"za.game/lib/tool"
)

type UserLoginLog struct {
	ID                 uint64    `db:"id"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
	Uid                int       `db:"uid"`
	Nickname           string    `db:"nickname"`
	ChannelId          int       `db:"channel_id"`
	Channel            string    `db:"channel"`
	Assets             int       `db:"assets"`
	ReferralCommission int       `db:"referral_commission"`
	Ip                 string    `db:"ip"`
	Device             string    `db:"device"`
	Version            string    `db:"version"`
	LoginMode          int       `db:"login_mode"`
	RegTime            time.Time `db:"reg_time"`
}

const (
	LoginMode = iota
	LoginModeGuest
	LoginModePhone
)

func (t *UserLoginLog) TableName() string {
	return "user_login_log_" + time.Now().Local().Format(tool.TimeFormatYm)
}

func (t *UserLoginLog) Save(db *sqlx.DB, fields, values string, args map[string]interface{}) (err error) {
	sql := fmt.Sprintf("insert into %s (%s) values (%s)", t.TableName(), fields, values)
	_, err = rds.SqlxNamedExec(db, sql, args)
	return
}
