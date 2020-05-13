package service

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/mocks"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"time"
)

func (suite *ServiceTestSuite) TestWeekCommand_weekNum() {
	t := suite.T()

	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = ""

	type testCase struct {
		now time.Time
		res string
	}
	testCases := []testCase{
		{
			time.Date(2020, 5, 1, 6, 40, 35, 0, tools.LocalTz),
			"Текущая неделя - первая",
		},
		{
			time.Date(2020, 4, 27, 18, 40, 35, 0, tools.LocalTz),
			"Текущая неделя - первая",
		},
		{
			time.Date(2020, 4, 30, 14, 50, 12, 0, tools.LocalTz),
			"Текущая неделя - первая",
		},
		{
			time.Date(2020, 3, 23, 15, 15, 15, 0, tools.LocalTz),
			"Текущая неделя - вторая",
		},
		{
			time.Date(2020, 5, 5, 15, 15, 15, 0, tools.LocalTz),
			"Текущая неделя - вторая",
		},
	}

	for i, tcase := range testCases {
		tools.Now = func() time.Time { return tcase.now }
		data := NewData(mocks.StartMessage, common.WeekCommand, "", MainKeyboard)
		err := WeekCommand(data)
		assert.NoErrorf(t, err, "testCase %d", i+1)
		assert.Equalf(t, tcase.res, data.Answer, "testCase %d", i+1)
	}
}

func (suite *ServiceTestSuite) TestWeekCommand() {
	t := suite.T()

	mocks.InitStartMocks()
	mocks.StartMessage.Message.MessageBody = ""

	type testCase struct {
		body string
		res  string
	}
	testCases := []testCase{
		{
			"1",
			"Выбери день недели",
		},
		{
			"2",
			"Выбери день недели",
		},
		{
			"0",
			"",
		},
		{
			"что-то еще",
			"",
		},
	}

	for i, tcase := range testCases {
		mocks.StartMessage.Message.MessageBody = tcase.body
		data := NewData(mocks.StartMessage, common.WeekCommand, "", MainKeyboard)
		err := WeekCommand(data)
		if tcase.res == "" {
			assert.Errorf(t, err, "testCase %d", i+1)
		} else {
			assert.NoErrorf(t, err, "testCase %d", i+1)
			assert.Equalf(t, tcase.res, data.Answer, "testCase %d", i+1)
		}
	}
}
