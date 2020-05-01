package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsStringInSlice(t *testing.T) {
	type testCase struct {
		value []string
		s     string
		res   bool
	}
	testCases := []testCase{
		{[]string{"string", "one more string", "mb one more", "test", "smth"}, "smth", true},
		{[]string{"string", "one more string", "mb one more", "test", "smth"}, "none", false},
		{[]string{}, "", false},
		{[]string{"da", "net", "random", ""}, "", true},
	}

	for i, tcase := range testCases {
		assert.Equalf(t, tcase.res, IsStringInSlice(tcase.s, tcase.value), "testCase %d", i+1)
	}
}
