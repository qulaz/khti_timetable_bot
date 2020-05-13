package db

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
)

func (suite *DBTestSuite) TestGetTimetable() {
	file := tools.LoadStringFromFile("db/testdata/test_timetable.json")
	testCase := &Timetable{}
	err := testCase.FromJson(file)
	tools.Fatal(suite.T(), suite.NoError(err))

	res, err := GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(testCase, res)
}

func (suite *DBTestSuite) TestUpdateTimetable() {
	file := tools.LoadStringFromFile("db/testdata/test_timetable.json")
	testCase := &Timetable{}
	err := testCase.FromJson(file)
	tools.Fatal(suite.T(), suite.NoError(err))
	testCase.Groups["58-1"].WeekSchedule[0]["понедельник"] = []string{"", "", "", "", ""}
	rawTestCase, err := testCase.ToJson()
	tools.Fatal(suite.T(), suite.NoError(err))

	err = UpdateTimetable(rawTestCase)
	tools.Fatal(suite.T(), suite.NoError(err))

	res, err := GetTimetable()
	tools.Fatal(suite.T(), suite.NoError(err))
	suite.Equal(testCase, res)
}
