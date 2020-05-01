package tools

import (
	"time"
)

// Возвращает номер недели (1 или 2)
func GetWeekNum(t time.Time) int {
	// Определяем первый учебный день. Обычно считается, что учеба начинается с первой недели
	firstSeptember := time.Date(getStudyYear(t), 9, 1, 0, 0, 0, 0, LocalTz)
	if int(firstSeptember.Weekday()) == 0 { // воскресенье
		firstSeptember = time.Date(getStudyYear(t), 9, 2, 0, 0, 0, 0, LocalTz)
	}

	_, firstWeek := firstSeptember.ISOWeek()
	_, currWeek := t.ISOWeek()

	if firstWeek > currWeek {
		return ((firstWeek - currWeek) % 2) + 1
	} else {
		return ((currWeek - firstWeek) % 2) + 1
	}
}

// Возвращает номер текущей недели (1 или 2)
//
// Но я не до конца уверен, что данный алгоритм будет правильно работать для всех годов
func GetCurrentWeekNum() int {
	return GetWeekNum(Now())
}

func getStudyYear(t time.Time) int {
	if int(t.Month()) < 9 {
		return t.Year() - 1
	} else {
		return t.Year()
	}
}

// Находит в расписании на день пары-окна и заменяет их с пустой строки на прочерк
func WindowAnalyzer(s []string) []string {
	var (
		prev  bool
		lastI int
	)

	// Определяет индекс последней пары в слайсе
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != "" {
			lastI = i
			break
		}
	}

	for i, lesson := range s {
		if lesson != "" {
			prev = true
		}
		if lesson == "" && prev && i < lastI {
			s[i] = "-"
		}
	}

	return s
}
