package service

import (
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

	return 1000
}
