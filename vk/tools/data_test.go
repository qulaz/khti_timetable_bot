package tools

import "testing"

func TestIsStringInSlice(t *testing.T) {
	type testCase struct {
		value []string
		s     string
		res   bool
	}
	testCases := []testCase{
		{[]string{"string", "one more string", "mb one more", "test", "smth"}, "smth", true},
		{[]string{"string", "one more string", "mb one more", "test", "smth"}, "none", false},
	}

	for _, test := range testCases {
		if v := IsStringInSlice(test.s, test.value); v != test.res {
			t.Errorf("Value: %+v; s: %q; Expected: %v, Got: %v", test.value, test.s, test.res, v)
		}
	}
}

func TestSliceOfIntsToString(t *testing.T) {
	type testCase struct {
		value []int
		sep   string
		res   string
	}
	testCases := []testCase{
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, ",", "1,2,3,4,5,6,7,8,9"},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, "", "123456789"},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, " ", "1 2 3 4 5 6 7 8 9"},
		{[]int{}, ",", ""},
	}

	for _, test := range testCases {
		if v := SliceOfIntsToString(test.value, test.sep); v != test.res {
			t.Errorf("Value: %+v; sep: %q; Expected: %v, Got: %v", test.value, test.sep, test.res, v)
		}
	}
}
