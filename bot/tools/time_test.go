package tools

import (
	"testing"
	"time"
)

func TestIsTimeBetween(t *testing.T) {
	type testCase struct {
		start time.Time
		end   time.Time
		check time.Time
		res   bool
	}

	testCases := []testCase{
		{
			TimeOnly(8, 15, 0, nil),
			TimeOnly(10, 5, 0, nil),
			TimeOnly(9, 30, 0, nil),
			true,
		},
		{
			TimeOnly(8, 15, 0, nil),
			TimeOnly(10, 5, 0, nil),
			TimeOnly(8, 0, 0, nil),
			false,
		},
		{
			TimeOnly(8, 15, 0, nil),
			TimeOnly(10, 5, 0, nil),
			TimeOnly(8, 15, 0, nil),
			true,
		},
		{
			TimeOnly(8, 15, 0, nil),
			TimeOnly(10, 5, 0, nil),
			TimeOnly(10, 5, 0, nil),
			true,
		},
		{
			TimeOnly(8, 15, 0, nil),
			TimeOnly(10, 5, 0, nil),
			TimeOnly(15, 40, 0, nil),
			false,
		},
		{
			TimeOnly(8, 15, 0, nil),
			TimeOnly(10, 5, 0, nil),
			TimeOnly(18, 22, 0, nil),
			false,
		},
	}

	for _, tcase := range testCases {
		if IsTimeBetween(tcase.start, tcase.end, tcase.check) != tcase.res {
			t.Errorf(
				"IsTimeBetween(start=%d:%d, end=%d:%d, check=%d:%d) != %v",
				tcase.start.Hour(), tcase.start.Minute(),
				tcase.end.Hour(), tcase.end.Minute(),
				tcase.check.Hour(), tcase.check.Minute(),
				tcase.res,
			)
		}
	}
}

func TestTodayName(t *testing.T) {
	type testCase struct {
		v   func() time.Time
		res string
	}
	testCases := []testCase{
		{
			func() time.Time { return time.Date(2020, 3, 2, 0, 0, 0, 0, time.UTC) },
			Weekdays[0],
		},
		{
			func() time.Time { return time.Date(2020, 3, 3, 0, 0, 0, 0, time.UTC) },
			Weekdays[1],
		},
		{
			func() time.Time { return time.Date(2020, 3, 4, 0, 0, 0, 0, time.UTC) },
			Weekdays[2],
		},
		{
			func() time.Time { return time.Date(2020, 3, 5, 0, 0, 0, 0, time.UTC) },
			Weekdays[3],
		},
		{
			func() time.Time { return time.Date(2020, 3, 6, 0, 0, 0, 0, time.UTC) },
			Weekdays[4],
		},
		{
			func() time.Time { return time.Date(2020, 3, 7, 0, 0, 0, 0, time.UTC) },
			Weekdays[5],
		},
		{
			func() time.Time { return time.Date(2020, 3, 8, 0, 0, 0, 0, time.UTC) },
			Sunday,
		},
	}

	for _, tcase := range testCases {
		Now = tcase.v
		if res := TodayName(); res != tcase.res {
			n := Now()
			t.Errorf(
				"%d.%s.%d определяется как %s, а это %s\n",
				n.Day(), n.Month(), n.Year(),
				res,
				tcase.res,
			)
		}
	}
}
