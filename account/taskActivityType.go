package account

import "database/sql"

// 活动主体信息
type TaskActivityInfoDB struct {
	ActivityId int          `db:"activity_id" gorm:"primarykey;column:activity_id;type:int(11);not null;comment:活动id;" json:"activity_id" `
	Name       string       `db:"name" gorm:"column:name;type:varchar(255);comment:活动名称（这只是给后台看看用的）;" json:"name" `
	ActiviType int8         `db:"activi_type" gorm:"column:activi_type;type:tinyint(4);comment:活动类型（1：以用户时间为主，2：以活动时间为主）;" json:"activi_type" `
	Starttime  sql.NullTime `db:"starttime" gorm:"column:starttime;type:datetime;comment:活动开始时间（以活动为主才生效）;" json:"starttime" `
	Stoptime   sql.NullTime `db:"stoptime" gorm:"column:stoptime;type:datetime;comment:活动结束时间（以活动为主才生效）;" json:"stoptime" `
	Interval   int          `db:"interval" gorm:"column:interval;type:int(11);comment:活动持续时间（以用户为主才生效）;" json:"interval" `
	IsClose    int8         `db:"is_close" gorm:"column:is_close;type:tinyint(4);default:0;comment:活动是否关闭（0：正常，1：关闭）;" json:"is_close" `
	OnlineTime sql.NullTime `db:"online_time" gorm:"column:online_time;type:datetime;not null;default:CURRENT_TIMESTAMP;comment:上线时间;" json:"online_time" `
}

// 活动任务配置
type TaskActivityConfigDB struct {
	ActivityId int   `db:"activity_id" gorm:"column:activity_id;type:int(11);not null;default:0;index:task_activity_cofig_idx;comment:活动id;" json:"activity_id" `
	TaskId     int   `db:"task_id" gorm:"column:task_id;type:int(11);default:0;index:task_activity_cofig_idx;comment:任务id;" json:"task_id" `
	Dayidx     int   `db:"dayidx" gorm:"column:dayidx;type:int(11);default:0;comment:属于第几天;" json:"dayidx" `
	MaxNum     int   `db:"max_num" gorm:"column:max_num;type:int(11);default:0;comment:需要完成的值;" json:"max_num" `
	AwardId    int64 `db:"award_id" gorm:"column:award_id;type:bigint(20);default:0;comment:奖励id;" json:"award_id" `
	AwardNum   int   `db:"award_num" gorm:"column:award_num;type:int(11);default:0;comment:奖励数量;" json:"award_num" `
	SortId     int   `db:"sort_id" gorm:"column:sort_id;type:int(11);not null;default:0;comment:排序id;" json:"sort_id" `
}

// 活动的全部信息，包含开关，和任务数组
type TaskActivityInfoAll struct {
	ActivityId int                    // 活动id`
	ActiviType int8                   // 活动类型（1：以用户时间为主）
	Starttime  int64                  // 活动开始时间（以活动为主才生效）
	Stoptime   int64                  // 活动结束时间（以活动为主才生效）
	Interval   int                    // 活动持续时间（以用户为主才生效）
	TaskArr    []TaskActivityConfigDB //任务配置数组
	OnlineTime int64                  // 上线时间
}

// 用户挑战信息数据表结构体
type TaskUserProgressDB struct {
	Uid        int64        `db:"uid" gorm:"column:uid;type:bigint(20);not null;index:task_user_progress_id2;comment:用户id;" json:"uid" `
	ActivityId int          `db:"activity_id" gorm:"column:activity_id;type:int(11);index:task_user_progress_id2;comment:活动id;" json:"activity_id" `
	TaskId     int          `db:"task_id" gorm:"column:task_id;type:int(11);index:task_user_progress_id2;comment:任务id;" json:"task_id" `
	IsLock     int8         `db:"is_lock" gorm:"column:is_lock;type:tinyint(4);comment:任务是否锁定(0:不锁定，1：锁定）;" json:"is_lock" `
	CurNum     int          `db:"cur_num" gorm:"column:cur_num;type:int(11);comment:当前完成的值;" json:"cur_num" `
	Lasttime   sql.NullTime `db:"lasttime" gorm:"column:lasttime;type:datetime;comment:任务上次完成的时间;" json:"lasttime" `
	Dayidx     int          `db:"dayidx" gorm:"column:dayidx;type:int(11);default:0;comment:属于第几天;" json:"dayidx" `
	MaxNum     int          `db:"max_num" gorm:"column:max_num;type:int(11);default:0;comment:需要完成的值;" json:"max_num" `
	AwardId    int64        `db:"award_id" gorm:"column:award_id;type:bigint(20);default:0;comment:奖励id;" json:"award_id" `
	AwardNum   int          `db:"award_num" gorm:"column:award_num;type:int(11);default:0;comment:奖励数量;" json:"award_num" `
	Createtime sql.NullTime `db:"createtime" gorm:"column:createtime;type:datetime;default:CURRENT_TIMESTAMP;comment:创建时间;" json:"createtime" `
	AwardGeted int8         `db:"award_geted" gorm:"column:award_geted;type:tinyint(4);not null;default:0;comment:奖励是否已领取（0：未领取，1：已领取）;" json:"award_geted" `
	SortId     int          `db:"sort_id" gorm:"column:sort_id;type:int(11);not null;default:0;comment:排序id;" json:"sort_id" `
}

// 用户任务信息结构体
type TaskUserInfo struct {
	ActivityId int   // 活动id
	TaskId     int   // 任务id
	IsLock     int8  // 任务是否锁定(0:不锁定，1：锁定）
	CurNum     int   // 当前完成的值
	Lasttime   int64 // 任务上次完成的时间
	Dayidx     int   // 属于第几天

	MaxNum     int   // 需要完成的值
	AwardId    int64 // 奖励id
	AwardNum   int   // 奖励数量
	Createtime int64 // 创建时间
	AwardGeted int8  // 奖励是否已领取（0：未领取，1：已领取）
	SortId     int   //任务排序id
}

// 奖励结构体
type TaskUserAwardDB struct {
	Uid        int64 `db:"uid" gorm:"column:uid;type:int64;" json:"uid" `                       // 用户id
	ActivityId int   `db:"activity_id" gorm:"column:activity_id;type:int;" json:"activity_id" ` // 活动id
	AwardId    int64 `db:"award_id" gorm:"column:award_id;type:int64;" json:"award_id" `        // 奖励id
	AwardNum   int   `db:"award_num" gorm:"column:award_num;type:int;" json:"award_num" `       // 奖励数量
}

// 用户奖励信息
type TaskUserAward struct {
	AwardId  int64 // 奖励id
	AwardNum int   // 奖励数量
}

// 用户参加活动信息数据表结构体
type TaskActivityUserinfoDB struct {
	Uid        int64        `db:"uid" gorm:"column:uid;type:int64;" json:"uid" `                          // 用户id
	ActivityId int          `db:"activity_id" gorm:"column:activity_id;type:int;" json:"activity_id" `    // 活动id
	Starttime  sql.NullTime `db:"starttime" gorm:"column:starttime;type:sql.NullTime;" json:"starttime" ` // 开始时间
	Stoptime   sql.NullTime `db:"stoptime" gorm:"column:stoptime;type:sql.NullTime;" json:"stoptime" `    // 结束时间
	IsStop     int8         `db:"is_stop" gorm:"column:is_stop;type:int8;" json:"is_stop" `               // 是否完全结束（0：不是，1：是）
}

// 用户参加活动信息
type TaskActivityUserinfo struct {
	ActivityId  int   // 活动id
	Starttime   int64 // 开始时间
	Stoptime    int64 // 结束时间
	IsDeleteOld bool  //是否删除旧的信息
	IsStop      int8  // 是否完全结束（0：不是，1：是）
}
