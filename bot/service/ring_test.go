package service

import (
	"github.com/stretchr/testify/assert"
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
