package valog

import (
	"fmt"
	"time"
)

type Log struct {
	logID uint
}

var OBLog *Log = &Log{logID: 0}

func (this *Log) LogMessage(strLog string) {
	fmt.Println("【LOG】", GetNowTime(), strLog)
}

func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (this *Log) LogMessageV2(logType int32, strLog string) {
	fmt.Println(fmt.Sprintf("【LOG - %d】 %s %s", logType, GetNowTime(), strLog))
}
