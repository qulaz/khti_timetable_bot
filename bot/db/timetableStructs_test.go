package db

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
)

func (suite *DBTestSuite) TestTimetable_WriteInDB() {
	timetable := &Timetable{}

	file := tools.LoadStringFromFile("db/testdata/test_timetable.json")
	err := timetable.FromJson(file)
	tools.Fatal(suite.T(), suite.NoError(err))
	err = timetable.WriteInDB()
	tools.Fatal(suite.T(), suite.NoError(err))

	for code := range timetable.Groups {
		_, err := GetGroupByGroupCode(code)
		suite.NoErrorf(err, "code %s", code)
	}

	res, err := GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(timetable, res)
}

func (suite *DBTestSuite) TestTimetable_GetStringifyedSchedule() {
	timetable, err := GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))
	weekdays := tools.Weekdays
	weekdays = append(weekdays, FullWeek)

	res, err := timetable.GetStringifyedSchedule("99-5", 1, weekdays[0])
	suite.Errorf(err, "res: %+v", res)

	for _, n := range []int{-2, -1, 0, 3, 4, 5, 55, 999, 5545} {
		_, err := timetable.GetStringifyedSchedule("58-1", n, weekdays[0])
		suite.Errorf(err, "weekNum: %d", n)
	}

	for _, n := range []string{"fef", "тест", "не день недели", "вт", "субота"} {
		res, err := timetable.GetStringifyedSchedule("58-1", 1, n)
		suite.Errorf(err, "dayName: %s;\nres: %s", n, res)
	}

	for code := range timetable.Groups {
		for weekNum := 0; weekNum < 2; weekNum++ {
			for _, dayName := range weekdays {
				s, err := timetable.GetStringifyedSchedule(code, weekNum+1, dayName)
				tools.Fatal(suite.T(), suite.NoError(err))

				tcase, _ := timetable.Groups[code].WeekSchedule[weekNum].ToString(dayName)
				suite.Equalf(tcase, s,
					"Расписание (%d, %s) для группы %q поученное через ф-ю не равно расписанию полученному вручную\n",
					weekNum+1, dayName, code,
				)
			}
		}

	}
}

func (suite *DBTestSuite) TestTimetable_GetWeekSchedule() {
	timetable, err := GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))

	for _, n := range []int{-2, -1, 0, 3, 4, 5, 55, 999, 5545} {
		_, err := timetable.GetWeekSchedule("58-1", n)
		suite.Errorf(err, "weekNum: %d", n)
	}

	_, err = timetable.GetWeekSchedule("99-5", 1)
	suite.Error(err)

	for code := range timetable.Groups {
		for weekNum := 0; weekNum < 2; weekNum++ {
			tcase := timetable.Groups[code].WeekSchedule[weekNum]
			res, err := timetable.GetWeekSchedule(code, weekNum+1)
			suite.NoErrorf(err, "group: %s", code)

			suite.Equal(tcase, res)
		}
	}
}

func (suite *DBTestSuite) TestWeekSchedule_GetDaySchedule() {
	timetable, err := GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))

	for code := range timetable.Groups {
		for weekNum := 0; weekNum < 2; weekNum++ {
			weekSchedule := timetable.Groups[code].WeekSchedule[weekNum]

			for _, dayName := range tools.Weekdays {
				_, err := weekSchedule.GetDaySchedule(dayName)
				suite.NoErrorf(err, "dayName: %s", dayName)
			}

			for _, dayName := range []string{tools.Sunday, FullWeek, "субота", "test", "тест"} {
				_, err := weekSchedule.GetDaySchedule(dayName)
				suite.Errorf(err, "dayName: %s", dayName)
			}
		}
	}
}

func (suite *DBTestSuite) TestWeekSchedule_ToString() {
	timetable, err := GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))

	weekdays := tools.Weekdays
	weekdays = append(weekdays, FullWeek)

	for code := range timetable.Groups {
		for weekNum := 0; weekNum < 2; weekNum++ {
			weekSchedule := timetable.Groups[code].WeekSchedule[weekNum]

			for _, dayName := range weekdays {
				_, err := weekSchedule.ToString(dayName)
				suite.NoErrorf(err, "dayName: %s", dayName)
			}

			for _, dayName := range []string{tools.Sunday, "субота", "test", "тест"} {
				_, err := weekSchedule.ToString(dayName)
				suite.Errorf(err, "dayName: %s", dayName)
			}
		}
	}
}
