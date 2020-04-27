package tools

import "testing"

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
