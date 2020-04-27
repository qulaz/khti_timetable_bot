package tools

import (
	"time"
)

var LocalTz, _ = time.LoadLocation("Asia/Krasnoyarsk")
var Now = func() time.Time { return time.Now().In(LocalTz) }

var Weekdays = []string{
	"понедельник",
	"вторник",
	"среда",
	"четверг",
	"пятница",
	"суббота",
}
var FormatedWeekdays = map[string]string{
	"понедельник": "Понедельник",
	"вторник":     "Вторник",
	"среда":       "Среда",
	"четверг":     "Четверг",
	"пятница":     "Пятница",
	"суббота":     "Суббота",
}

const Sunday = "воскресенье"

func TimeOnly(hour, min, sec int, loc *time.Location) time.Time {
	if loc == nil {
		loc = LocalTz
	}
	return time.Date(0, 0, 0, hour, min, sec, 0, loc)
}

func IsTimeBetween(start, end, check time.Time) bool {
	return (check.After(start) || check.Equal(start)) && (check.Before(end) || check.Equal(end))
}

func GetWeekdayName(now time.Time) string {
	weekday := now.Weekday()
	if weekday == 0 {
		return Sunday
	}

	return Weekdays[weekday-1]
}

func TodayName() string {
	return GetWeekdayName(Now())
}

func RemoveDate(t time.Time) time.Time {
	return TimeOnly(t.Hour(), t.Minute(), t.Second(), time.UTC)
}
