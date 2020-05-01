package service

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"time"
)

// Расписание звонков
var LessonTimes = [][]time.Time{
	{
		tools.TimeOnly(8, 30, 0, tools.LocalTz), // начало пары
		tools.TimeOnly(10, 5, 0, tools.LocalTz), // конец пары
	},
	{
		tools.TimeOnly(10, 15, 0, tools.LocalTz),
		tools.TimeOnly(11, 50, 0, tools.LocalTz),
	},
	{
		tools.TimeOnly(12, 00, 0, tools.LocalTz),
		tools.TimeOnly(13, 35, 0, tools.LocalTz),
	},
	{
		tools.TimeOnly(14, 10, 0, tools.LocalTz),
		tools.TimeOnly(15, 45, 0, tools.LocalTz),
	},
	{
		tools.TimeOnly(15, 55, 0, tools.LocalTz),
		tools.TimeOnly(17, 30, 0, tools.LocalTz),
	},
}

// Определяет идет ли в данный момент занятие или нет
func IsLesson() bool {
	n := tools.RemoveDate(tools.Now())

	for _, lessonDuration := range LessonTimes {
		if tools.IsTimeBetween(lessonDuration[0], lessonDuration[1], n) {
			return true
		}
	}

	return false
}

// Возвращает порядковый номер текущей пары. Нумерация начинается с единицы.
// Крайние случаи: -1 - пары еще не начались, 999 - пару уже закончились, 1000 - не удалось определить текущую пару
func CurrentLessonNum() int {
	n := tools.RemoveDate(tools.Now())

	// Текущее время меньше старта занятий
	if firstLessonStart := LessonTimes[0][0]; n.Before(firstLessonStart) {
		return -1
	}

	// Текущее время позже конца занятий
	if lastLessonStart := LessonTimes[len(LessonTimes)-1][1]; n.After(lastLessonStart) {
		return 999
	}

	isLesson := IsLesson()
	for i, lessonDuration := range LessonTimes {
		// текущее время находится между началом и концом i-й пары
		isiLessonNow := tools.IsTimeBetween(lessonDuration[0], lessonDuration[1], n)
		// i-я пара не последняя
		isiNotLastLesson := i < len(LessonTimes)-1
		// текущее время находится между концом прошлой пары и началом i-й пары, т.е. сейчас перемена
		isiBreak := !isLesson && tools.IsTimeBetween(lessonDuration[1], LessonTimes[i+1][0], n)

		// Если сейчас идет i-я пара или перемена после не последней пары - возвращаем i+1 - номер текущей пары
		if isiLessonNow || (isiNotLastLesson && isiBreak) {
			return i + 1
		}
	}

	helpers.Logger.Errorw("Не нашлась идущая сейчас пара CurrentLessonNum()", "now", n)
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("func", "CurrentLessonNum")
		scope.AddBreadcrumb(
			&sentry.Breadcrumb{Data: map[string]interface{}{"now": n}}, 1,
		)

		sentry.CaptureException(errors.New("Не нашлась текущая пара CurrentLessonNum()"))
	})

	return 1000
}

// Возвращает время до конца/начала пары в зависимости от текущего времени
func TimeToRing() time.Duration {
	now := tools.RemoveDate(tools.Now())
	currentLessonNum := CurrentLessonNum()
	isLesson := IsLesson()

	if currentLessonNum == 1000 {
		return -1
	}

	// Пары еще не начались / уже закончились
	if currentLessonNum == -1 || currentLessonNum == 999 || (!isLesson && currentLessonNum == 5) {
		currentLessonNum = 1
	}

	if isLesson {
		endLessonTime := LessonTimes[currentLessonNum-1][1]
		return endLessonTime.Sub(now)
	} else {
		// Если пары еще не начались/уже закончились - отдаем время до первой пары
		if currentLessonNum == 1 && (now.Before(LessonTimes[0][0]) || now.After(LessonTimes[len(LessonTimes)-1][1])) {
			firstLessonTime := LessonTimes[0][0]

			// Добавляем 1 день ко времени начала первой пары, если она начнется на следующий день
			if now.Hour() >= LessonTimes[len(LessonTimes)-1][1].Hour() {
				firstLessonTime = firstLessonTime.Add(time.Hour * 24)
			}

			return firstLessonTime.Sub(now)
		}

		startLessonTime := LessonTimes[currentLessonNum][0]
		return startLessonTime.Sub(now)
	}
}

func RingCommand(d *Data) error {
	ringDuration := TimeToRing()
	if ringDuration < 0 {
		helpers.Logger.Errorw("TimeToRing вернул отрицательный Duration",
			"now", tools.Now(), "res", ringDuration,
		)
		return errors.Errorf(
			"TimeToRing вернул отрицательный Duration.\nnow %v;\nres %d", tools.Now(), ringDuration,
		)
	}
	strDuration := tools.DurationToString(ringDuration)

	if IsLesson() {
		d.Answer = fmt.Sprintf("До конца пары осталось %s", strDuration)
	} else {
		d.Answer = fmt.Sprintf("До начала пары осталось %s", strDuration)
	}

	return nil
}
