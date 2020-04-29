package tools

import (
	"github.com/getsentry/sentry-go"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"time"
)

var LocalTz, err = time.LoadLocation("Asia/Krasnoyarsk")
var Now = func() time.Time { return time.Now().In(LocalTz) }

func init() {
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("init", "true")
			sentry.CaptureException(err)
		})
		helpers.Logger.Fatalw("Не удалось загрузить часовой пояс Красноярска", "err", err)
	}
}

// Название дней недели используемые в качестве ключей
var Weekdays = []string{
	"понедельник",
	"вторник",
	"среда",
	"четверг",
	"пятница",
	"суббота",
}
var FormattedWeekdays = map[string]string{
	"понедельник": "Понедельник",
	"вторник":     "Вторник",
	"среда":       "Среда",
	"четверг":     "Четверг",
	"пятница":     "Пятница",
	"суббота":     "Суббота",
}

const Sunday = "воскресенье"

// Возвращает time.Time с обнуленной датой
func TimeOnly(hour, min, sec int, loc *time.Location) time.Time {
	if loc == nil {
		loc = LocalTz
	}
	return time.Date(0, 0, 0, hour, min, sec, 0, loc)
}

// Находится ли дата check между start и end, либо равна одной из них
func IsTimeBetween(start, end, check time.Time) bool {
	return (check.After(start) || check.Equal(start)) && (check.Before(end) || check.Equal(end))
}

// Возвращает название дня недели переданного time.Time
func GetWeekdayName(now time.Time) string {
	weekday := now.Weekday()
	if weekday == 0 {
		return Sunday
	}

	return Weekdays[weekday-1]
}

// Возвращает название сегодняшнего дня недели
func TodayName() string {
	return GetWeekdayName(Now())
}

// Возвращает переданный time.Time с обнуленной датой
func RemoveDate(t time.Time) time.Time {
	return TimeOnly(t.Hour(), t.Minute(), t.Second(), LocalTz)
}
