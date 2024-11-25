package utils

import (
	"time"
)

// 时间格式常量
const (
	TimeFormat      = "2006-01-02 15:04:05"
	DateFormat      = "2006-01-02"
	TimeFormatNoSec = "2006-01-02 15:04"
)

// 格式化时间
func FormatTime(t time.Time) string {
	return t.Format(TimeFormat)
}

// 格式化日期
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}

// 解析时间字符串
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(TimeFormat, timeStr)
}

// 获取当天开始时间
func GetDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// 获取当天结束时间
func GetDayEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// 检查时间是否在范围内
func IsTimeInRange(t, start, end time.Time) bool {
	return (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end))
}
