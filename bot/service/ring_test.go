package service

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/mocks"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"testing"
	"time"
)

func TestCurrentLessonNum(t *testing.T) {
	type testCase struct {
		v   func() time.Time
		res int
	}
	testCases := []testCase{
		{
			func() time.Time { return tools.TimeOnly(9, 30, 14, nil) },
			1,
		},
		{
			func() time.Time { return tools.TimeOnly(4, 25, 26, nil) },
			-1,
		},
		{
			func() time.Time { return tools.TimeOnly(8, 00, 45, nil) },
			-1,
		},
		{
			func() time.Time { return tools.TimeOnly(8, 30, 0, nil) },
			1,
		},
		{
			func() time.Time { return tools.TimeOnly(13, 50, 34, nil) },
			3,
		},
		{
			func() time.Time { return tools.TimeOnly(15, 10, 55, nil) },
			4,
		},
		{
			func() time.Time { return tools.TimeOnly(10, 10, 19, nil) },
			1,
		},
		{
			func() time.Time { return tools.TimeOnly(18, 50, 23, nil) },
			999,
		},
		{
			func() time.Time { return tools.TimeOnly(23, 50, 41, nil) },
			999,
		},
		{
			func() time.Time { return tools.TimeOnly(5, 50, 0, nil) },
			-1,
		},
		{
			func() time.Time { return tools.TimeOnly(0, 0, 0, nil) },
			-1,
		},
		{
			func() time.Time { return tools.TimeOnly(9, 40, 0, nil) },
			1,
		},
		{
			func() time.Time { return tools.TimeOnly(10, 5, 0, nil) },
			1,
		},
		{
			func() time.Time { return tools.TimeOnly(17, 30, 0, nil) },
			5,
		},
		{
			func() time.Time { return tools.TimeOnly(17, 30, 1, nil) },
			999,
		},
	}

	for i, tcase := range testCases {
		tools.Now = tcase.v // подмена ф-ии определения текущего времени на нашу, отдающую статику
		assert.Equalf(t, tcase.res, CurrentLessonNum(), "testCase %d", i+1)
	}
}

func TestIsLesson(t *testing.T) {
	type testCase struct {
		v   func() time.Time
		res bool
	}
	testCases := []testCase{
		{
			func() time.Time { return tools.TimeOnly(9, 30, 0, nil) },
			true,
		},
		{
			func() time.Time { return tools.TimeOnly(4, 25, 0, nil) },
			false,
		},
		{
			func() time.Time { return tools.TimeOnly(8, 00, 0, nil) },
			false,
		},
		{
			func() time.Time { return tools.TimeOnly(8, 30, 0, nil) },
			true,
		},
		{
			func() time.Time { return tools.TimeOnly(13, 50, 0, nil) },
			false,
		},
		{
			func() time.Time { return tools.TimeOnly(15, 10, 0, nil) },
			true,
		},
		{
			func() time.Time { return tools.TimeOnly(10, 10, 0, nil) },
			false,
		},
		{
			func() time.Time { return tools.TimeOnly(18, 50, 0, nil) },
			false,
		},
		{
			func() time.Time { return tools.TimeOnly(23, 50, 0, nil) },
			false,
		},
		{
			func() time.Time { return tools.TimeOnly(5, 50, 0, nil) },
			false,
		},
		{
			func() time.Time { return tools.TimeOnly(0, 0, 0, nil) },
			false,
		},
	}

	for i, tcase := range testCases {
		tools.Now = tcase.v // подмена ф-ии определения текущего времени на нашу, отдающую статику
		assert.Equalf(t, tcase.res, IsLesson(), "testCase %d", i)
	}
}

func TestTimeToRing(t *testing.T) {
	type testCase struct {
		v   func() time.Time
		res time.Duration
	}
	testCases := []testCase{
		{
			func() time.Time { return time.Date(2020, 3, 26, 8, 0, 34, 0, tools.LocalTz) },
			//"До начала пары осталось 29 минут",
			time.Minute*29 + time.Second*26,
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 0, 0, 23, 0, tools.LocalTz) },
			//"До начала пары осталось 8:29",
			time.Hour*8 + time.Minute*29 + time.Second*37,
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 17, 31, 47, 0, tools.LocalTz) },
			//"До начала пары осталось 14:58",
			time.Hour*14 + time.Minute*58 + time.Second*13,
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 0, 12, 29, 0, tools.LocalTz) },
			//"До начала пары осталось 8:17",
			time.Hour*8 + time.Minute*17 + time.Second*31,
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 11, 50, 1, 0, tools.LocalTz) },
			//"До начала пары осталось 9 минут",
			time.Minute*9 + time.Second*59,
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 10, 0, 15, 0, tools.LocalTz) },
			//"До конца пары осталось 4 минут",
			time.Minute*4 + time.Second*45,
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 14, 10, 40, 0, tools.LocalTz) },
			//"До конца пары осталось 1:34",
			time.Hour + time.Minute*34 + time.Second*20,
		},
	}

	for i, tcase := range testCases {
		tools.Now = tcase.v
		assert.Equalf(t, tcase.res, TimeToRing(), "testCase %d", i)
	}
}

func (suite *ServiceTestSuite) TestRingCommand() {
	t := suite.T()

	common.TestInits()
	mocks.InitStartMocks()
	type testCase struct {
		now    func() time.Time
		answer string
	}
	testCases := []testCase{
		{
			func() time.Time { return time.Date(2020, 3, 26, 8, 0, 34, 0, tools.LocalTz) },
			"До начала пары осталось 30 минут",
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 0, 0, 23, 0, tools.LocalTz) },
			"До начала пары осталось 8 часов и 30 минут",
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 17, 31, 47, 0, tools.LocalTz) },
			"До начала пары осталось 14 часов и 59 минут",
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 0, 12, 29, 0, tools.LocalTz) },
			"До начала пары осталось 8 часов и 18 минут",
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 11, 50, 1, 0, tools.LocalTz) },
			"До начала пары осталось 10 минут",
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 10, 0, 15, 0, tools.LocalTz) },
			"До конца пары осталось 5 минут",
		},
		{
			func() time.Time { return time.Date(2020, 3, 26, 14, 10, 40, 0, tools.LocalTz) },
			"До конца пары осталось 1 час и 35 минут",
		},
	}

	for i, tcase := range testCases {
		tools.Now = tcase.now
		data := NewData(mocks.StartMessage, common.RingCommand, common.UnknownErrorMessage, MainKeyboard)
		err := RingCommand(data)
		assert.NoError(t, err)
		assert.Equalf(t, tcase.answer, data.Answer, "testCase %d", i+1)
	}
}
