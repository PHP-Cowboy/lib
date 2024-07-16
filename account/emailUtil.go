package account

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/jmoiron/sqlx"
	"za.game/lib/consts"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
)

type EmailUser struct {
	ID          uint64     `db:"id" gorm:"primarykey;column:id;type:int(11);";json:"id"`
	EmailId     uint32     `db:"email_id" gorm:"column:email_id;type:int(11);not null;default:0;";json:"email_id"`          //邮件ID
	Type        uint8      `db:"type" gorm:"column:type;type:tinyint(4);not null;default:1;"`                               //邮件类型(1=普通邮件,2=赠送退款邮件)
	Uid         uint64     `db:"uid" gorm:"column:uid;type:bigint(20);not null;default:0;";json:"uid"`                      //用户ID
	Title       string     `db:"title" gorm:"column:title;type:varchar(64);not null;default'';";json:"title"`               //邮件标题
	Msg         string     `db:"msg" gorm:"column:msg;type:varchar(1024);not null;default:'';";json:"msg"`                  //邮件内容
	Status      uint8      `db:"status" gorm:"column:status;type:tinyint(4);not null;default:0;";json:"status"`             //读状态(0=未读,1=已读)
	ReadTime    *time.Time `db:"read_time" gorm:"column:read_time;";json:"read_time"`                                       //读取时间
	GetTime     *time.Time `db:"get_time" gorm:"column:get_time;";json:"get_time"`                                          //兑换时间
	IsAnnex     uint8      `db:"is_annex" gorm:"column:is_annex;type:tinyint(4);not null;default:0;";json:"is_annex"`       //是否有附件
	AnnexIds    string     `db:"annex_ids" gorm:"column:annex_ids;type:varchar(128);not null;default:0;";json:"annex_ids"`  //附件IDS(多个用,分割)
	InsertTime  int64      `db:"insert_time" gorm:"column:insert_time;type:int(11);not null;default:0;";json:"insert_time"` //插入时间戳
	CreatedAt   *time.Time `db:"created_at" gorm:"column:created_at;";json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" gorm:"column:updated_at;";json:"updated_at"`
	HasAttach   uint8      `db:"has_attach" gorm:"type:tinyint(4);not null;default:0;"` //是否有附件(1=有,0=否)
	Attachments string     `db:"attachments" gorm:"type:varchar(255);"`                 //附件信息
}

//func BaseEmailUser() *gorm.DB {
//	return config.GameDB.Model(new(EmailUser))
//}

func GetEmailUserTableName(uid uint64) string {
	return fmt.Sprintf("email_0%d", uid%10)
}

func (t *EmailUser) TableName(uid uint64) string {
	return fmt.Sprintf("email_0%d", uid%10)
}

func (t *EmailUser) Save(fields, values string, mp map[string]interface{}, uid uint64) (err error) {
	sql := fmt.Sprintf("insert into %s(%s) values(%s)", t.TableName(uid), fields, values)

	_, err = rds.SqlxNamedExec(dbconn.GameDB, sql, mp)
	return
}

func (t *EmailUser) BatchSaveTx(tx *sqlx.Tx, users []*EmailUser, uid uint64) (err error) {
	sql := fmt.Sprintf(
		"insert into %s(email_id,uid,title,msg,status,read_time,get_time,is_annex,annex_ids,insert_time,created_at,updated_at) "+
			"values(:email_id,:uid,:title,:msg,:status,:read_time,:get_time,:is_annex,:annex_ids,:insert_time,:created_at,:updated_at)",
		t.TableName(uid),
	)

	_, err = rds.SqlxNamedExecTx(tx, sql, users)
	return
}

type Attachment struct {
	ItemId int `json:"item_id"`
	Nums   int `json:"nums"`
}

type SendEmailParams struct {
	Uid         uint64
	Title       string
	Msg         string
	Attachments []Attachment
	Type        int
}

// 向用户发送邮件
func SendEmail(param *SendEmailParams) (err error) {
	now := time.Now().Local()

	emailObj := new(EmailUser)

	HasAttach := 0
	if len(param.Attachments) > 0 {
		HasAttach = 1
	}

	Attachments := ""
	var b []byte
	b, err = json.Marshal(param.Attachments)

	if err != nil {
		logger.Logger.Errorf("SendEmail: to json failed! err:[%v]", err)
		return
	}

	if param.Type == 0 {
		param.Type = 1
	}

	Attachments = string(b)

	err = emailObj.Save(
		"uid,title,msg,insert_time,created_at,updated_at,has_attach,attachments,`type`",
		":uid,:title,:msg,:insert_time,:created_at,:updated_at,:has_attach,:attachments,:type",
		map[string]interface{}{
			"uid":         param.Uid,
			"title":       param.Title,
			"msg":         param.Msg,
			"insert_time": now.Unix(),
			"created_at":  now,
			"updated_at":  now,
			"has_attach":  HasAttach,
			"attachments": Attachments,
			"type":        param.Type,
		},
		param.Uid,
	)

	if err != nil {
		logger.Logger.Errorf("SendEmail: send email failed! err:[%v]", err)
		return
	}

	//添加红点，新邮件需要添加计数2：读取和领取
	if HasAttach == 1 {
		AddRedDot(param.Uid, strconv.FormatUint(param.Uid, 10), consts.RedDot_Email, 2, false, true)
	} else {
		AddRedDot(param.Uid, strconv.FormatUint(param.Uid, 10), consts.RedDot_Email, 1, false, true)
	}

	return
}
