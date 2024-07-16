package tool

import (
	"strconv"
	"time"
)

/**文件件记录一些公用的时间函数*/

/**获取当天固定时间的时间戳*/
func TimeTodayBy(hour, min, sec int) (timestamp int64) {
	time_data := time.Now()

	time_data2 := time.Date(time_data.Year(), time_data.Month(), time_data.Day(), hour, min, sec, 0, time_data.Location())
	//fmt.Println(time_data2)
	timestamp = time_data2.Unix()
	//fmt.Println(timestamp2)

	return timestamp
}

/**获得距离第二天0点，还剩多少秒*/
func TimeSecFromTomorrow() (sec int64) {
	tomorrow := TimeTodayBy(0, 0, 0) + 86400
	return tomorrow - time.Now().Unix()
}

/**获得ymd的时间int*/
func TimeGetYmd() (ymd int) {
	now := time.Now().Local()
	ymd, _ = strconv.Atoi(now.Format("20060102"))
	return ymd
}

/**获取两个时间的天数差*/
func TimeDaysDiff(time_begin, time_end time.Time) (diff int) {
	diff = time_end.YearDay() - time_begin.YearDay()
	return
}

/**判断两个时间戳是否是同一个月*/
func TimeInSameMonth(t1, t2 int64) bool {
	time1 := time.Unix(t1, 0)
	time2 := time.Unix(t2, 0)

	return time1.Year() == time2.Year() && time1.Month() == time2.Month()
}
