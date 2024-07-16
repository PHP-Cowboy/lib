package log

import (
	"time"
	"za.game/lib/sqlInfo"
)

type RoomFundsFlowLog struct {
	sqlInfo.RoomFundsFlowLog
}

func (t *RoomFundsFlowLog) TableName() string {
	return "room_funds_flow_log_" + time.Now().Format("200601")
}

func (t *RoomFundsFlowLog) GetTableNameByTime(date time.Time) string {
	return "room_funds_flow_log_" + date.Format("200601")
}
