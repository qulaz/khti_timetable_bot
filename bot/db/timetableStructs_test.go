package db

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"reflect"
	"testing"
)

func TestTimetable_WriteInDB(t *testing.T) {
	timetable := &Timetable{}

	file := tools.LoadStringFromFile("db/testdata/test_timetable.json")
	if err := timetable.FromJson(file); err != nil {
		t.Fatalf("Ошибка при конвертировании json-строки в структуру расписания: %+v\n", err)
	}

	if err := timetable.WriteInDB(); err != nil {
		t.Fatalf("Ошибка записи расписания в базу данных: %+v\n", err)
	}

	for code, _ := range timetable.Groups {
		if _, err := GetGroupByGroupCode(code); err != nil {
			t.Errorf("Не найдена группа с кодом %q: %+v\n", code, err)
		}
	}

	res, err := GetTimetable()
	if err != nil {
		t.Fatalf("Ошибка получения расписания из базы данных: %+v\n", err)
	}

	if !reflect.DeepEqual(timetable, res) {
		t.Fatalf("Ожидаемое расписание не равно полученному из базы данных\n")
	}
}

func TestTimetable_GetStringifyedSchedule(t *testing.T) {
	PrepareTestDatabase()

	timetable, err := GetTimetable()
	if err != nil {
		t.Fatalf("Ошибка получения расписания из базы данных: %+v\n", err)
	}
	weekdays := tools.Weekdays
	weekdays = append(weekdays, FullWeek)

	if res, err := timetable.GetStringifyedSchedule("99-5", 1, weekdays[0]); err == nil {
		t.Errorf("Ф-я отработала без ошибок с переданной несуществующей группой. Результат ее работы: %+v\n", res)
	}

	for _, n := range []int{-2, -1, 0, 3, 4, 5, 55, 999, 5545} {
		if _, err := timetable.GetStringifyedSchedule("58-1", n, weekdays[0]); err == nil {
			t.Errorf("Ф-я приняла значение weekNum == %d\n", n)
		}
	}

	for _, n := range []string{"fef", "тест", "не день недели", "вт", "субота"} {
		if res, err := timetable.GetStringifyedSchedule("58-1", 1, n); err == nil {
			t.Errorf("Ф-я приняла значение dayName == %s: %+v\n", n, res)
		}
	}

	for code, _ := range timetable.Groups {
		for weekNum := 0; weekNum < 2; weekNum++ {
			for _, dayName := range weekdays {
				s, err := timetable.GetStringifyedSchedule(code, weekNum+1, dayName)
				if err != nil {
					t.Fatalf("%+v\n", err)
				}

				tcase, _ := timetable.Groups[code].WeekSchedule[weekNum].ToString(dayName)

				if s != tcase {
					t.Errorf(
						"Расписание (%d, %s) для группы %q поученное через ф-ю не равно расписанию полученному вручную\n",
						weekNum+1, dayName, code,
					)
				}
			}
		}

	}
}

func TestTimetable_GetWeekSchedule(t *testing.T) {
	PrepareTestDatabase()

	timetable, err := GetTimetable()
	if err != nil {
		t.Fatalf("Ошибка получения расписания из базы данных: %+v\n", err)
	}

	for _, n := range []int{-2, -1, 0, 3, 4, 5, 55, 999, 5545} {
		if _, err := timetable.GetWeekSchedule("58-1", n); err == nil {
			t.Errorf("Ф-я приняла значение weekNum == %d\n", n)
		}
	}

	if res, err := timetable.GetWeekSchedule("99-5", 1); err == nil {
		t.Errorf("Ф-я отработала без ошибок с переданной несуществующей группой. Результат ее работы: %+v\n", res)
	}

	for code, _ := range timetable.Groups {
		for weekNum := 0; weekNum < 2; weekNum++ {
			tcase := timetable.Groups[code].WeekSchedule[weekNum]
			res, err := timetable.GetWeekSchedule(code, weekNum+1)
			if err != nil {
				t.Errorf("Ошибка при получении расписания (%d) для группы %q: %+v\n", weekNum+1, code, err)
			}
			if !reflect.DeepEqual(tcase, res) {
				t.Errorf(
					"WeekSchedule (%d) полученный вручную для группы %q не равен WeekSchedule "+
						"полученному через timetable.GetWeekSchedule\n",
					weekNum+1, code,
				)
			}
		}
	}
}

func TestWeekSchedule_GetDaySchedule(t *testing.T) {
	PrepareTestDatabase()

	timetable, err := GetTimetable()
	if err != nil {
		t.Fatalf("Ошибка получения расписания из базы данных: %+v\n", err)
	}

	for code, _ := range timetable.Groups {
		for weekNum := 0; weekNum < 2; weekNum++ {
			weekSchedule := timetable.Groups[code].WeekSchedule[weekNum]

			for _, dayName := range tools.Weekdays {
				if _, err := weekSchedule.GetDaySchedule(dayName); err != nil {
					t.Errorf("Ошибка при получении дня %q: %+v\n", dayName, err)
				}
			}

			for _, dayName := range []string{tools.Sunday, FullWeek, "субота", "test", "тест"} {
				if d, err := weekSchedule.GetDaySchedule(dayName); err == nil {
					t.Errorf("Ф-я приняла значение dayName == %q: %+v\n", dayName, d)
				}
			}
		}
	}
}

func TestWeekSchedule_ToString(t *testing.T) {
	PrepareTestDatabase()

	timetable, err := GetTimetable()
	if err != nil {
		t.Fatalf("Ошибка получения расписания из базы данных: %+v\n", err)
	}

	weekdays := tools.Weekdays
	weekdays = append(weekdays, FullWeek)

	for code, _ := range timetable.Groups {
		for weekNum := 0; weekNum < 2; weekNum++ {
			weekSchedule := timetable.Groups[code].WeekSchedule[weekNum]

			for _, dayName := range weekdays {
				if _, err := weekSchedule.ToString(dayName); err != nil {
					t.Errorf("Ошибка получения расписания (%d, %s): %+v\n", weekNum+1, dayName, err)
				}
			}

			for _, dayName := range []string{tools.Sunday, "субота", "test", "тест"} {
				if d, err := weekSchedule.ToString(dayName); err == nil {
					t.Errorf("Ф-я приняла значение dayName == %q: %+v\n", dayName, d)
				}
			}
		}
	}
}
