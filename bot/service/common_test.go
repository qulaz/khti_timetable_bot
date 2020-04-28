package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getPrevAndNextOffset(t *testing.T) {
	type testCase struct {
		total      int
		limit      int
		offset     int
		prevOffset int
		nextOffset int
	}
	testCases := []testCase{
		{40, 27, 0, -1, 27},
		{10, 27, 0, -1, -1},
		{50, 27, 27, 0, -1},
		{80, 27, 27, 0, 54},
		{80, 27, 54, 27, -1},
		{250, 100, 0, -1, 100},
		{250, 100, 100, 0, 200},
		{250, 100, 200, 100, -1},
		{187, 50, 24, 0, 74},
		{33, 27, 0, -1, 27},
		{33, 27, 27, 0, -1},
	}

	for i, tcase := range testCases {
		prev, next := getPrevAndNextOffset(tcase.total, tcase.limit, tcase.offset)
		assert.Equalf(t, tcase.prevOffset, prev, "testCase %d", i+1)
		assert.Equalf(t, tcase.nextOffset, next, "testCase %d", i+1)
	}
}
