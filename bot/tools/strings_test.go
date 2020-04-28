package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoveWhitespaces(t *testing.T) {
	type testCase struct {
		v   string
		res string
	}
	testCases := []testCase{
		{
			"  текст текст текст           еще текст и еще текст ",
			"текст текст текст еще текст и еще текст",
		},
		{
			"текст текст текст",
			"текст текст текст",
		},
		{
			"немного текста  еще текст после двух пробелов",
			"немного текста еще текст после двух пробелов",
		},
		{
			"текст текст текстик\n\nтекст после 2-ух переносов строк      и еще текст после нескольких пробелов ",
			"текст текст текстик\n\nтекст после 2-ух переносов строк и еще текст после нескольких пробелов",
		},
	}

	for i, tcase := range testCases {
		if res := RemoveWhitespaces(tcase.v); res != tcase.res {
			t.Errorf("Тест кейс #%d завалился. %s != %s\n", i+1, res, tcase.res)
		}
	}
}

func TestSelectNonEmptyString(t *testing.T) {
	type testCase struct {
		s1  string
		s2  string
		res string
	}
	testCases := []testCase{
		{"", "nonempty", "nonempty"},
		{"nonempty1", "nonempty2", "nonempty1"},
		{"   	 	", "nonempty", "nonempty"},
		{"nonempty", "     			", "nonempty"},
		{"", "     			", ""},
	}

	for i, tcase := range testCases {
		assert.Equalf(t, tcase.res, SelectNonEmptyString(tcase.s1, tcase.s2), "testCase %d", i)
	}
}
