package db

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"reflect"
	"testing"
)

func TestGetTimetable(t *testing.T) {
	PrepareTestDatabase()

	file := tools.LoadStringFromFile("db/testdata/test_timetable.json")
	testCase := &Timetable{}
	if err := testCase.FromJson(file); err != nil {
		t.Fatalf("Ошибка при конвертировании json в структуру расписания: %+v\n", err)
	}

	res, err := GetTimetable()
	if err != nil {
		t.Fatalf("Ошибка при получении расписания из базы данных: %+v\n", err)
	}

	if !reflect.DeepEqual(testCase, res) {
		t.Fatalf("Ожидаемое расписание не равно полученному из базы данных")
	}
}

func TestUpdateTimetable(t *testing.T) {
	PrepareTestDatabase()

	file := tools.LoadStringFromFile("db/testdata/test_timetable.json")
	testCase := &Timetable{}
	if err := testCase.FromJson(file); err != nil {
		t.Fatalf("Ошибка при конвертировании json в структуру расписания: %+v\n", err)
	}
	testCase.Groups["58-1"].WeekSchedule[0]["понедельник"] = []string{"", "", "", "", ""}
	rawTestCase, err := testCase.ToJson()
	if err != nil {
		t.Fatalf("Ошибка приведения расписания к json-строке: %+v\n", err)
	}

	if err := UpdateTimetable(rawTestCase); err != nil {
		t.Fatalf("Ошибка при обновлении расписания: %+v\n", err)
	}

	res, err := GetTimetable()
	if err != nil {
		t.Fatalf("Ошибка при получении расписания из базы данных: %+v\n", err)
	}

	if !reflect.DeepEqual(testCase, res) {
		t.Fatalf("Ожидаемое обновленное расписание не равно полученному из базы данных\n")
	}
}
