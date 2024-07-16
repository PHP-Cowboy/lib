/**
=======================================
	本文件，主要为任务功能
=======================================
*/

package account

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"git.dev666.cc/external/breezedup/goserver/core/logger"
	"github.com/alitto/pond"
	"github.com/gomodule/redigo/redis"
	"za.game/lib/consts"
	"za.game/lib/dbconn"
	"za.game/lib/rds"
	"za.game/lib/tool"
)

// 任务管理池，用来异步执行完成任务
var map_task_pool map[int]*pond.WorkerPool
var task_workers int = 2 //线程数量

func init() {
	//创建一个线程池，工作线程（1），最多任务（100）。
	map_task_pool = make(map[int]*pond.WorkerPool, 0)
	for i := 0; i < task_workers; i++ {
		map_task_pool[i] = pond.New(1, 100)
	}
}

/**完成任务接口*/
func FinishTask(uidstr string, uid uint64, arr_taskid []int, arr_finish []int) {

	//异步投递任务，同一个用户只能投递到同一个线程中，防止数据冲突
	workid := int(uid) % task_workers
	map_task_pool[workid].Submit(func() {
		//获得开着的活动主体
		map_activity, err := GetActivityInfo()
		if err != nil {
			return
		}

		actids := []int{}
		for _, val := range map_activity {
			actids = append(actids, val.ActivityId)
		}
		//活动用户活动信息
		map_activity_userinfo, err := GetUserAvticityInfo(uidstr, uid, actids)
		if err != nil {
			return
		}
		//因为这里已经100%用户有任务了，所以不用取获取活动任务信息，再去创建用户任务了
		//直接获得用户信息即可
		map_userinfo, err := GetUserTaskAll(uidstr, uid, map_activity_userinfo)
		if err != nil {
			return
		}

		FinishUserTask(uidstr, uid, map_userinfo, arr_taskid, arr_finish)
	})
}

/**获得开着的活动配置信息*/
func GetActivityInfo() (map_info map[int]TaskActivityInfoAll, err error) {
	//获得活动配置信息
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()

	reply, err := redis.String(redisConn.Do("GET", consts.RedisTaskActivityInfo))

	arrdbinfo := make([]TaskActivityInfoDB, 0)

	if reply != "" {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &arrdbinfo)
		if err != nil {
			logger.Logger.Errorf("getActivityInfo: parse config failed! err:[%v]", err)
		}
	}

	if err != nil {
		err = rds.SqlxSelect(dbconn.NDB, &arrdbinfo, "select activity_id,`name`,activi_type,starttime,stoptime,`interval`,is_close,online_time from task_activity_info where is_close=0")
		if err != nil {
			logger.Logger.Errorf("func: select sql failed! err:[%v]", err)
			return
		}

		//存入redis
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(arrdbinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("getActivityInfo: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", consts.RedisTaskActivityInfo, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("getActivityInfo: set config reids failed! err:[%v]", err)
		}
	}

	time_now := time.Now().Unix()
	map_info = make(map[int]TaskActivityInfoAll, len(arrdbinfo))
	for _, val := range arrdbinfo {
		if val.ActiviType == consts.TaskActivityTypeActivityCentric && val.Stoptime.Time.Unix() <= time_now { //周期活动，并且已结束，跳过
			continue
		}
		item := TaskActivityInfoAll{
			ActivityId: val.ActivityId,
			ActiviType: val.ActiviType,
			Starttime:  val.Starttime.Time.Unix(),
			Stoptime:   val.Stoptime.Time.Unix(),
			Interval:   val.Interval,
			TaskArr:    make([]TaskActivityConfigDB, 0),
			OnlineTime: val.OnlineTime.Time.Unix(),
		}
		map_info[item.ActivityId] = item
	}

	return map_info, nil
}

/**获得活动对应的全部任务信息*/
func GetActivityTask(map_activity map[int]TaskActivityInfoAll) (err error) {
	if len(map_activity) == 0 {
		return
	}
	//获得活动配置信息
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()

	reply, err := redis.String(redisConn.Do("GET", consts.RedisTaskActivityConfig))

	arr_task_activity_info_db := make([]TaskActivityConfigDB, 0)

	if reply != "" {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &arr_task_activity_info_db)
		if err != nil {
			logger.Logger.Errorf("getActivityTask: parse config failed! err:[%v]", err)
		}
	}

	if err != nil {
		where_sql := ""
		for _, val := range map_activity {
			if where_sql != "" {
				where_sql += ","
			}
			where_sql += strconv.Itoa(val.ActivityId)
		}
		where_sql = "select activity_id,task_id,dayidx,max_num,award_id,award_num,sort_id from task_activity_config where activity_id in (" + where_sql + ")"
		err = rds.SqlxSelect(dbconn.NDB, &arr_task_activity_info_db, where_sql)
		if err != nil {
			logger.Logger.Errorf("getActivityTask: select sql failed! err:[%v]", err)
			return
		}

		//存入redis
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(arr_task_activity_info_db)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("getActivityTask: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", consts.RedisTaskActivityConfig, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("getActivityTask: set config reids failed! err:[%v]", err)
		}
	}

	for _, val := range arr_task_activity_info_db {
		if activity, ok := map_activity[val.ActivityId]; ok {
			activity.TaskArr = append(activity.TaskArr, val)
			map_activity[val.ActivityId] = activity
		}
	}
	return
}

/**获得用户参加活动的信息*/
func GetUserAvticityInfo(uidstr string, uid uint64, actids []int) (map_userinfo map[int]TaskActivityUserinfo, err error) {
	if len(actids) == 0 {
		return
	}
	//获得活动配置信息
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()

	redis_key := consts.RedisTaskActivityUserInfo + uidstr
	//logger.Logger.Errorf("GetUserAvticityInfo: uid:[%v]", uidstr)
	reply, err := redis.String(redisConn.Do("GET", redis_key))

	map_userinfo = make(map[int]TaskActivityUserinfo, 0)

	if reply != "" {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &map_userinfo)
		if err != nil {
			logger.Logger.Errorf("GetUserAvticityInfo: parse config failed! err:[%v]", err)
		}
	}

	//redis中没有数据，则查询数据库，并同步
	if err != nil {
		var arr_task_activity_userinfo_db []TaskActivityUserinfoDB
		where_sql := ""
		for _, ActivityId := range actids {
			if where_sql != "" {
				where_sql += ","
			}
			where_sql += strconv.Itoa(ActivityId)
		}
		where_sql = "select uid,activity_id,starttime,stoptime from task_activity_userinfo where uid = '" + uidstr + "' and activity_id in (" + where_sql + ")"
		err = rds.SqlxSelect(dbconn.NDB, &arr_task_activity_userinfo_db, where_sql)
		if err != nil {
			logger.Logger.Errorf("GetUserAvticityInfo: select sql failed! err:[%v]", err)
			return
		}

		//解析层内存能用的数据
		for _, val := range arr_task_activity_userinfo_db {
			item := TaskActivityUserinfo{
				ActivityId:  val.ActivityId,
				Starttime:   val.Starttime.Time.Unix(),
				Stoptime:    val.Stoptime.Time.Unix(),
				IsDeleteOld: false,
				IsStop:      val.IsStop,
			}
			map_userinfo[val.ActivityId] = item
		}

		//同步redis
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(map_userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("GetUserAvticityInfo: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("GetUserAvticityInfo: set config reids failed! err:[%v]", err)
		}
	}

	return map_userinfo, nil
}

/*
*创建用户参加活动信息
 */
func CreateUserActivityInfo(uidstr string, uid uint64, map_activity map[int]TaskActivityInfoAll, map_userinfo map[int]TaskActivityUserinfo) (err error) {
	if len(map_activity) == 0 {
		return
	}
	//如果活动没参加过，则创建
	not_exist_id := []int{}
	for _, val := range map_activity {
		if _, ok := map_userinfo[val.ActivityId]; !ok {
			not_exist_id = append(not_exist_id, val.ActivityId)
		}
	}
	insert_values := []interface{}{}
	insert_values_str := ""
	time_now := time.Now()
	bsync := false
	for _, activity_id := range not_exist_id {
		info := map_activity[activity_id]
		//生成插入数据
		insert_columns := map[string]interface{}{
			"uid":         uid,
			"activity_id": activity_id,
			"starttime":   time_now,
			"stoptime":    time.Unix(time_now.Unix()+int64(info.Interval), 0),
			"is_stop":     0,
		}
		insert_values = append(insert_values, insert_columns)
		//同步放一份到内存里，省的再次读取，并再最后同步redis
		item := TaskActivityUserinfo{
			ActivityId:  activity_id,
			Starttime:   time_now.Unix(),
			Stoptime:    time_now.Unix() + int64(info.Interval),
			IsDeleteOld: false,
			IsStop:      0,
		}

		map_userinfo[activity_id] = item

		if insert_values_str != "" {
			insert_values_str += ","
		}
		insert_values_str += "(:uid,:activity_id,:starttime,:stoptime,:is_stop)"
	}
	if len(insert_values) > 0 {
		insert_sql := "insert into task_activity_userinfo(uid,activity_id,starttime,stoptime,is_stop) values" + insert_values_str
		_, err = rds.SqlxNamedExec(dbconn.NDB, insert_sql, insert_values)
		if err != nil {
			logger.Logger.Errorf("createUserActivityInfo: insert sql failed! err:[%v]", err)
			return
		}
		bsync = true
	}

	//检测，如果有新的活动，且是新的周期，则同步数据
	for activity_id, val := range map_userinfo {
		activity_type := map_activity[activity_id].ActiviType
		if val.Stoptime < time_now.Unix() && activity_type == consts.TaskActivityTypeActivityCentric { //已经过期的周期活动，修改信息
			stoptime := time.Unix(time_now.Unix()+int64(map_activity[activity_id].Interval), 0)
			update_columns := map[string]interface{}{
				"uid":         uid,
				"activity_id": activity_id,
				"starttime":   time_now,
				"stoptime":    stoptime,
			}
			_, err = rds.SqlxNamedExec(dbconn.NDB, "update task_activity_userinfo set starttime=:starttime,stoptime=:stoptime where uid=:uid and activity_id=:activity_id", update_columns)
			if err != nil {
				logger.Logger.Errorf("createUserActivityInfo: update sql failed! err:[%v]", err)
				return
			}
			val.Starttime = time_now.Unix()
			val.Stoptime = stoptime.Unix()
			val.IsDeleteOld = true
			map_userinfo[activity_id] = val
			bsync = true
		}
	}

	if bsync {
		//把数据同步一遍redis
		redisConn := dbconn.RedisPool.Get()
		defer redisConn.Close()
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(map_userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("createUserActivityInfo: config to json failed! err:[%v]", err)
			return nil
		}

		//设置缓存数据
		redis_key := consts.RedisTaskActivityUserInfo + uidstr
		//logger.Logger.Errorf("GetUserAvticityInfo: uid:[%v]", uidstr)
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("createUserActivityInfo: set config reids failed! err:[%v]", err)
			return nil
		}
	}

	return nil
}

// 获取用户奖励信息
func GetUserAward(uidstr string, uid uint64, map_activity_userinfo map[int]TaskActivityUserinfo) (map_userinfo map[int]TaskUserAward, err error) {
	if len(map_activity_userinfo) == 0 {
		return
	}
	//获得活动配置信息
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()

	redis_key := consts.RedisTaskUserAward + uidstr
	reply, err := redis.String(redisConn.Do("GET", redis_key))

	map_userinfo = make(map[int]TaskUserAward, 0)

	if reply != "" {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &map_userinfo)
		if err != nil {
			logger.Logger.Errorf("GetUserAward: parse config failed! err:[%v]", err)
		}
	}

	//redis中没有数据，则查询数据库，并同步
	if err != nil {
		where_sql := ""
		time_now := time.Now().Unix()
		for _, val := range map_activity_userinfo {
			if val.Stoptime <= time_now { //活动过期，则不查询
				continue
			}
			if where_sql != "" {
				where_sql += ","
			}
			where_sql += strconv.Itoa(val.ActivityId)
		}
		if where_sql == "" {
			return map_userinfo, nil
		}
		where_sql = "select uid,activity_id,award_id,award_num from task_user_award where uid = '" + uidstr + "' and activity_id in (" + where_sql + ")"
		var arr_task_user_award_db []TaskUserAwardDB
		err = rds.SqlxSelect(dbconn.NDB, &arr_task_user_award_db, where_sql)
		if err != nil {
			logger.Logger.Errorf("GetUserAward: select sql failed! err:[%v]", err)
			return
		}

		for _, val := range arr_task_user_award_db {
			item := TaskUserAward{
				AwardId:  val.AwardId,
				AwardNum: val.AwardNum,
			}
			map_userinfo[val.ActivityId] = item
		}

		//同步redis
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(map_userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("GetUserAward: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("GetUserAward: set config reids failed! err:[%v]", err)
		}
	}

	return map_userinfo, nil
}

/*
*创建用户参加奖励信息
 */
func CreateUserActivityAwardInfo(uidstr string, uid uint64, map_activity map[int]TaskActivityUserinfo, map_userinfo map[int]TaskUserAward) (err error) {
	if len(map_activity) == 0 {
		return
	}
	//如果活动没参加过，则创建
	not_exist_id := []int{}
	time_now := time.Now().Unix()
	for _, val := range map_activity {
		if val.Stoptime <= time_now { //用户活动时间过期，则不进行新建操作
			continue
		}
		if _, ok := map_userinfo[val.ActivityId]; !ok {
			not_exist_id = append(not_exist_id, val.ActivityId)
		}
	}
	insert_values := []interface{}{}
	insert_values_str := ""
	for _, activity_id := range not_exist_id {
		//生成插入数据
		awardid := int64(0)
		if activity_id == consts.AvtId_NewPlayer { //是新手嘉年华
			awardid = consts.ItemIdNewPlayerCoin
		}
		insert_columns := map[string]interface{}{
			"uid":         uid,
			"activity_id": activity_id,
			"award_id":    awardid,
			"award_num":   0,
		}
		insert_values = append(insert_values, insert_columns)
		//同步放一份到内存里，省的再次读取，并再最后同步redis
		item := TaskUserAward{
			AwardId:  awardid,
			AwardNum: 0,
		}

		map_userinfo[activity_id] = item

		if insert_values_str != "" {
			insert_values_str += ","
		}
		insert_values_str += "(:uid,:activity_id,:award_id,:award_num)"
	}
	if len(insert_values) > 0 {
		insert_sql := "insert into task_user_award(uid,activity_id,award_id,award_num) values" + insert_values_str
		_, err = rds.SqlxNamedExec(dbconn.NDB, insert_sql, insert_values)
		if err != nil {
			logger.Logger.Errorf("createUserActivityAwardInfo: insert sql failed! err:[%v]", err)
			return
		}

		//把数据同步一遍redis
		redisConn := dbconn.RedisPool.Get()
		defer redisConn.Close()
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(map_userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("createUserActivityAwardInfo: config to json failed! err:[%v]", err)
			return nil
		}

		//设置缓存数据
		redis_key := consts.RedisTaskUserAward + uidstr
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("createUserActivityAwardInfo: set config reids failed! err:[%v]", err)
			return nil
		}
	}

	return nil
}

/**获得用户全部活动任务的进度信息*/
func GetUserTaskAll(uidstr string, uid uint64, map_activity_userinfo map[int]TaskActivityUserinfo) (map_userinfo map[int][]TaskUserInfo, err error) {
	if len(map_activity_userinfo) == 0 {
		return
	}
	//获得活动配置信息
	redisConn := dbconn.RedisPool.Get()
	defer redisConn.Close()

	redis_key := consts.RedisTaskUserInfo + uidstr
	reply, err := redis.String(redisConn.Do("GET", redis_key))

	map_userinfo = make(map[int][]TaskUserInfo, 0)

	if reply != "" {
		//解析redis数据到map结构中
		err = json.Unmarshal([]byte(reply), &map_userinfo)
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: parse config failed! err:[%v]", err)
		}
	}

	//redis中没有数据，则查询数据库，并同步
	if err != nil {
		where_sql := ""
		time_now := time.Now().Unix()
		for _, val := range map_activity_userinfo {
			if val.Stoptime <= time_now { //过期的活动不查询数据
				continue
			}
			if where_sql != "" {
				where_sql += ","
			}
			where_sql += strconv.Itoa(val.ActivityId)
		}
		if where_sql == "" {
			return map_userinfo, nil
		}
		where_sql = "select uid,activity_id,task_id,is_lock,cur_num,lasttime,dayidx,max_num,award_id,award_num,createtime,award_geted,sort_id from task_user_progress where uid = '" + uidstr + "' and activity_id in (" + where_sql + ")"
		var arr_task_user_progress_db []TaskUserProgressDB
		err = rds.SqlxSelect(dbconn.NDB, &arr_task_user_progress_db, where_sql)
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: select sql failed! err:[%v]", err)
			return
		}

		//解析层内存能用的数据
		for _, val := range arr_task_user_progress_db {
			item := TaskUserInfo{
				ActivityId: val.ActivityId,
				TaskId:     val.TaskId,
				IsLock:     val.IsLock,
				CurNum:     val.CurNum,
				Lasttime:   val.Lasttime.Time.Unix(),
				Dayidx:     val.Dayidx,
				MaxNum:     val.MaxNum,
				AwardId:    val.AwardId,
				AwardNum:   val.AwardNum,
				Createtime: val.Createtime.Time.Unix(),
				AwardGeted: val.AwardGeted,
				SortId:     val.SortId,
			}
			if _, ok := map_userinfo[val.TaskId]; !ok {
				map_userinfo[val.TaskId] = make([]TaskUserInfo, 0)
			}
			map_userinfo[val.TaskId] = append(map_userinfo[val.TaskId], item)
		}

		//同步redis
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(map_userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: set config reids failed! err:[%v]", err)
		}
	}

	return map_userinfo, nil
}

/*
*创建用户活动任务信息
@param map_userinfo: 用户任务信息map[taskid][]taskinfo
*/
func CreateUserActivityTask(uidstr string, uid uint64, map_activity map[int]TaskActivityInfoAll, map_activity_userinfo map[int]TaskActivityUserinfo, map_userinfo map[int][]TaskUserInfo) (err error) {
	if len(map_activity_userinfo) == 0 {
		return
	}

	//获取要创建的活动
	// activity_ids := []int{}
	// for _, arr_info := range map_userinfo {
	// 	for _, info := range arr_info {
	// 		if _, ok := map_activity_userinfo[info.ActivityId]; !ok {
	// 			activity_ids = append(activity_ids, info.ActivityId)
	// 		}
	// 	}
	// }
	activity_ids := []int{}
	time_now := time.Now().Unix()
	for _, act := range map_activity {
		if _, ok := map_activity_userinfo[act.ActivityId]; ok {
			if map_activity_userinfo[act.ActivityId].Stoptime <= time_now { //用户过期的活动，不创建新任务
				continue
			}
		}
		//遍历所有任务
		bexists := false
		for _, arr_info := range map_userinfo {
			for _, info := range arr_info {
				if info.ActivityId == act.ActivityId {
					bexists = true
					break
				}
			}
			if bexists {
				break
			}
		}
		if bexists {
			continue
		}
		//走到这一步，说明当前活动没有任务
		activity_ids = append(activity_ids, act.ActivityId)
	}

	//先删除需要删除的老的任务
	where_sql := ""
	for _, val := range map_activity_userinfo {
		if val.IsDeleteOld {
			if where_sql != "" {
				where_sql += ","
			}
			where_sql += strconv.Itoa(val.ActivityId)
			activity_ids = append(activity_ids, val.ActivityId)
		}
	}

	if where_sql != "" {
		delete_columns := map[string]interface{}{
			"uid":       uid,
			"where_sql": where_sql,
		}
		_, err = rds.SqlxNamedExec(dbconn.NDB, "delete task_user_progress where uid=:uid and activity_id in (:where_sql)", delete_columns)
		if err != nil {
			logger.Logger.Errorf("createUserActivityTask: delete sql failed! err:[%v]", err)
			return
		}
	}

	insert_values := []interface{}{}
	insert_values_str := ""
	rd_count := 0
	for _, activity_id := range activity_ids {
		info := map_activity[activity_id]
		for _, taskinfo := range info.TaskArr {
			//生成插入数据
			is_lock := int8(1)
			cur_num := 0
			db_lasttime := sql.NullTime{}
			lasttime := int64(0)
			if taskinfo.Dayidx == 0 || taskinfo.Dayidx == 1 {
				is_lock = 0                                      //如果没有天数限制，或者有天数限制，则直接解锁第一天
				if taskinfo.TaskId == consts.TaskId_LoginAccum { //是登陆的任务，默认完成
					cur_num = 1
					db_lasttime.Time = time.Now()
					db_lasttime.Valid = true
					lasttime = db_lasttime.Time.Unix()
					rd_count += 1
				}
			}
			insert_columns := map[string]interface{}{
				"uid":         uid,
				"activity_id": activity_id,
				"task_id":     taskinfo.TaskId,
				"is_lock":     is_lock,
				"cur_num":     cur_num,
				"lasttime":    db_lasttime,
				"dayidx":      taskinfo.Dayidx,
				"max_num":     taskinfo.MaxNum,
				"award_id":    taskinfo.AwardId,
				"award_num":   taskinfo.AwardNum,
				"sort_id":     taskinfo.SortId,
			}
			insert_values = append(insert_values, insert_columns)
			//同步放一份到内存里，省的再次读取，并再最后同步redis
			item := TaskUserInfo{
				ActivityId: activity_id,
				TaskId:     taskinfo.TaskId,
				IsLock:     is_lock,
				CurNum:     cur_num,
				Lasttime:   lasttime,
				Dayidx:     taskinfo.Dayidx,
				MaxNum:     taskinfo.MaxNum,
				AwardId:    taskinfo.AwardId,
				AwardNum:   taskinfo.AwardNum,
				Createtime: time.Now().Unix(),
				SortId:     taskinfo.SortId,
			}
			if _, ok := map_userinfo[taskinfo.TaskId]; !ok {
				map_userinfo[taskinfo.TaskId] = make([]TaskUserInfo, 0)
			}
			map_userinfo[taskinfo.TaskId] = append(map_userinfo[taskinfo.TaskId], item)
		}

		if insert_values_str != "" {
			insert_values_str += ","
		}
		insert_values_str += "(:uid,:activity_id,:task_id,:is_lock,:cur_num,:lasttime,:dayidx,:max_num,:award_id,:award_num,:sort_id)"
	}

	if len(insert_values) > 0 {
		insert_sql := "insert into task_user_progress(uid,activity_id,task_id,is_lock,cur_num,lasttime,dayidx,max_num,award_id,award_num,sort_id) values" + insert_values_str
		_, err = rds.SqlxNamedExec(dbconn.NDB, insert_sql, insert_values)
		if err != nil {
			logger.Logger.Errorf("createUserActivityTask: insert sql failed! err:[%v]", err)
			return
		}

		//把数据同步一遍redis
		redisConn := dbconn.RedisPool.Get()
		defer redisConn.Close()
		jsonTxt := ""

		jsonTxt, err = tool.JsonString(map_userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("createUserActivityTask: config to json failed! err:[%v]", err)
			return nil
		}

		//设置缓存数据
		redis_key := consts.RedisTaskUserInfo + uidstr
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: set config reids failed! err:[%v]", err)
			return nil
		}

		//有任务完成，增加红点
		if rd_count > 0 {
			AddRedDot(uid, uidstr, consts.RedDot_ActNewPlayer, rd_count, false, true)
		}
	}

	return nil
}

/*
*任务解锁(目前只有时间解锁，如果以后有新的解锁，则加条件模块)
@param map_userinfo: 用户任务信息map[taskid][]taskinfo
*/
func UnlockUserTask(uidstr string, uid uint64, map_userinfo map[int][]TaskUserInfo) (bok bool) {
	//循环便利任务，查看哪些任务可以被解锁
	unlock_list := []*TaskUserInfo{} //可以解锁的数组[[活动id，任务id]]
	rd_count := 0
	blogin := false //登陆任务只能完成一次
	for taskid, arr_task := range map_userinfo {
		for idx, taskinfo := range arr_task {
			if taskinfo.IsLock == 1 { //任务是锁定的，检测是否可以解锁
				if taskinfo.Dayidx > 1 { //任务是按天数解锁的
					day_diff := tool.TimeDaysDiff(time.Unix(taskinfo.Createtime, 0), time.Now()) + 1
					if day_diff >= taskinfo.Dayidx { //可以解锁
						map_userinfo[taskid][idx].IsLock = 0 //设置内存解锁

						if taskinfo.TaskId == consts.TaskId_LoginAccum { //解锁每日登陆的时候，起始就触发了完成
							map_userinfo[taskid][idx].Lasttime = time.Now().Unix()
							if !blogin {
								map_userinfo[taskid][idx].CurNum = taskinfo.Dayidx
								blogin = true
								rd_count += 1
							}
						}

						unlock_list = append(unlock_list, &map_userinfo[taskid][idx])
					}
				}

				//预留，其他解锁类型处理
			} else {
				if taskinfo.TaskId == consts.TaskId_LoginAccum && !blogin { //查看进入是否已经完成过登陆任务
					day_diff := tool.TimeDaysDiff(time.Unix(taskinfo.Lasttime, 0), time.Now()) + 1
					if day_diff == 1 {
						blogin = true
					}
				}
			}
		}
	}

	if len(unlock_list) > 0 { //同步数据库，同步redis
		for i := 0; i < len(unlock_list); i++ {
			taskitem := unlock_list[i]
			lasttime := sql.NullTime{}
			if taskitem.Lasttime > 0 {
				lasttime.Time = time.Unix(taskitem.Lasttime, 0)
				lasttime.Valid = true
			}
			update_columns := map[string]interface{}{
				"uid":         uid,
				"activity_id": taskitem.ActivityId,
				"task_id":     taskitem.TaskId,
				"is_lock":     0,
				"dayidx":      taskitem.Dayidx,
				"cur_num":     taskitem.CurNum,
				"lasttime":    lasttime,
			}
			sqlres, err := rds.SqlxNamedExec(dbconn.NDB, "update task_user_progress set is_lock=:is_lock,cur_num=:cur_num,lasttime=:lasttime where uid=:uid and activity_id=:activity_id and task_id=:task_id and dayidx=:dayidx", update_columns)
			if err != nil {
				logger.Logger.Errorf("UnlockUserTask: update sql failed! err:[%v],uid:[%d]", err, uid)
				return false
			}
			rows, _ := sqlres.RowsAffected()
			if rows <= 0 {
				logger.Logger.Errorf("UnlockUserTask: update sql has no change! err:[%v],uid:[%d]", err, uid)
				return false
			}
		}
		//同步redis
		redisConn := dbconn.RedisPool.Get()
		defer redisConn.Close()
		jsonTxt := ""

		jsonTxt, err := tool.JsonString(map_userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		redis_key := consts.RedisTaskUserInfo + uidstr
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: set config reids failed! err:[%v]", err)
		}

		//有任务完成，增加红点
		if rd_count > 0 {
			AddRedDot(uid, uidstr, consts.RedDot_ActNewPlayer, rd_count, false, true)
		}
	}

	return true
}

/**任务完成任务
*@param arr_taskid: 完成的任务数组
*@param arr_finish: 完成的对应数量*/
func FinishUserTask(uidstr string, uid uint64, map_userinfo map[int][]TaskUserInfo, arr_taskid []int, arr_finish []int) (bok bool) {
	//找到对应任务，修改他的当前值
	finish_list := []*TaskUserInfo{} //可以完成的任务
	bfinish_login := false           //是否今日已完成登陆，一天只能完成一次（新手嘉年华）
	for idx, req_taskid := range arr_taskid {
		req_finish_num := arr_finish[idx]
		if _, ok := map_userinfo[req_taskid]; ok {
			for i := 0; i < len(map_userinfo[req_taskid]); i++ {
				task_item := &map_userinfo[req_taskid][i]
				if task_item.ActivityId == consts.AvtId_NewPlayer && task_item.TaskId == consts.TaskId_LoginAccum {
					if bfinish_login {
						continue
					}
					//新手嘉年华，累计登陆
					time_last := time.Unix(task_item.Lasttime, 0)
					time_now := time.Now()
					day_diff := tool.TimeDaysDiff(time_last, time_now)
					if day_diff == 0 { //今天已经有完成过一个登陆，则所有属于新手嘉年华的累计登陆任务都不能完成
						bfinish_login = true
						continue
					}
				}
				if checkCanDoTask(task_item, req_finish_num) { //可以增加本次完成的数值
					task_item.CurNum += req_finish_num
					if task_item.ActivityId == consts.AvtId_NewPlayer && task_item.TaskId == consts.TaskId_LoginAccum { //累计登陆，完成就是最大值
						task_item.CurNum = task_item.MaxNum
						bfinish_login = true
					}
					if task_item.CurNum > task_item.MaxNum {
						task_item.CurNum = task_item.MaxNum
					}
					task_item.Lasttime = time.Now().Unix()
					finish_list = append(finish_list, task_item)
				}
			}
		}
	}

	if len(finish_list) > 0 { //同步数据库，同步redis
		red_count := 0
		for i := 0; i < len(finish_list); i++ {
			task_item := finish_list[i]
			lasttime := sql.NullTime{}
			if task_item.Lasttime > 0 {
				lasttime.Time = time.Unix(task_item.Lasttime, 0)
				lasttime.Valid = true
			}
			update_columns := map[string]interface{}{
				"uid":         uid,
				"activity_id": task_item.ActivityId,
				"task_id":     task_item.TaskId,
				"cur_num":     task_item.CurNum,
				"dayidx":      task_item.Dayidx,
				"lasttime":    lasttime,
			}
			sqlres, err := rds.SqlxNamedExec(dbconn.NDB, "update task_user_progress set cur_num=:cur_num,lasttime=:lasttime where uid=:uid and activity_id=:activity_id and task_id=:task_id and dayidx=:dayidx", update_columns)
			if err != nil {
				logger.Logger.Errorf("UnlockUserTask: update sql failed! err:[%v],uid:[%d]", err, uid)
				return false
			}
			rows, _ := sqlres.RowsAffected()
			if rows <= 0 {
				logger.Logger.Errorf("FinishUserTask: update sql has no change! err:[%v],uid:[%d],taskid:[%v],actid:[%v],day:[%v]", err, uid, task_item.TaskId, task_item.ActivityId, task_item.Dayidx)
				return false
			}

			if task_item.CurNum >= task_item.MaxNum {
				red_count += 1
			}
		}
		//同步redis
		redisConn := dbconn.RedisPool.Get()
		defer redisConn.Close()
		jsonTxt := ""

		jsonTxt, err := tool.JsonString(map_userinfo)

		//转换出错直接返回map，不继续操作
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: config to json failed! err:[%v]", err)
		}

		//设置缓存数据
		redis_key := consts.RedisTaskUserInfo + uidstr
		_, err = redisConn.Do("Set", redis_key, jsonTxt, "EX", 5*60)
		if err != nil {
			logger.Logger.Errorf("getUserTaskAll: set config reids failed! err:[%v]", err)
		}

		if red_count > 0 {
			//增加红点
			AddRedDot(uid, uidstr, consts.RedDot_ActNewPlayer, red_count, false, true)
		}
	}

	return true
}

/**检测某个任务是否可以完成*/
func checkCanDoTask(task_item *TaskUserInfo, finish_num int) bool {
	if task_item.IsLock == 0 && task_item.CurNum < task_item.MaxNum { //以解锁，且未完成，则可以触发
		if task_item.TaskId == consts.TaskId_LoginAccum { //如果是【累计登陆任务】，则需要查看是否是累计
			time_last := time.Unix(task_item.Lasttime, 0)
			time_now := time.Now()
			day_diff := tool.TimeDaysDiff(time_last, time_now)
			if day_diff == 0 { //今天已经登陆过，不需要再完成
				return false
			}
		} else if task_item.TaskId == consts.TaskId_RechargeSingle { //单笔充值是否达到
			if finish_num < task_item.MaxNum { //没达到要求
				return false
			}
		}

		return true
	}
	if task_item.TaskId == consts.TaskId_RechargeAccum && task_item.CurNum < task_item.MaxNum { //累计充值任务，不需要解锁也能完成
		return true
	}
	return false
}
