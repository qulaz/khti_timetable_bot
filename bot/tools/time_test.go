package tools

import (
	"github.com/stretchr/testify/assert"
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

func TestDurationToHoursMinutesSeconds(t *testing.T) {
	type testCase struct {
		d       time.Duration
		hours   int
		minutes int
		seconds int
	}
	testCases := []testCase{
		{time.Second * 15, 0, 0, 15},
		{time.Minute + time.Second*14, 0, 1, 14},
		{time.Minute*21 + time.Second*14, 0, 21, 14},
		{time.Minute*61 + time.Second*14, 1, 1, 14},
		{time.Hour + time.Minute*61 + time.Second*14, 2, 1, 14},
	}

	for i, tcase := range testCases {
		hours, minutes, seconds := DurationToHoursMinutesSeconds(tcase.d)
		assert.Equalf(t, tcase.hours, hours, "testCase %d", i+1)
		assert.Equalf(t, tcase.minutes, minutes, "testCase %d", i+1)
		assert.Equalf(t, tcase.seconds, seconds, "testCase %d", i+1)
	}
}

func TestDurationToString(t *testing.T) {
	type testCase struct {
		d   time.Duration
		res string
	}
	testCases := []testCase{
		{time.Second * 5, "1 минута"},
		{time.Minute * 45, "45 минут"},
		{time.Minute*45 + time.Second*34, "46 минут"},
		{time.Hour + time.Minute*5 + time.Second*10, "1 час и 6 минут"},
		{time.Hour*2 + time.Minute*22 + time.Second*10, "2 часа и 23 минуты"},
	}

	for i, tcase := range testCases {
		assert.Equalf(t, tcase.res, DurationToString(tcase.d), "testCase %d", i+1)
	}
}
