package parser

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"testing"
)

func TestParse__wrong_file_type(t *testing.T) {
	_, err := Parse("parser/testdata/timetable.json")
	assert.Error(t, err)
}

func TestParse__nonexistent_file(t *testing.T) {
	_, err := Parse("parser/testdata/table.xls")
	assert.Error(t, err)
}

func TestParse__xls(t *testing.T) {
	timetable, err := Parse("parser/testdata/timetable.xls")
	tools.Fatal(t, assert.NoError(t, err))

	j, err := timetable.ToJson()
	tools.Fatal(t, assert.NoError(t, err))
	f := tools.LoadStringFromFile("parser/testdata/timetable.json")
	assert.Equal(t, f, j)
}

func TestParse__xlsx(t *testing.T) {
	timetable, err := Parse("parser/testdata/timetable.xlsx")
	tools.Fatal(t, assert.NoError(t, err))

	j, err := timetable.ToJson()
	tools.Fatal(t, assert.NoError(t, err))
	f := tools.LoadStringFromFile("parser/testdata/timetable.json")
	assert.Equal(t, f, j)
}

func TestParse__xls_and_xlsx_parsers_the_same(t *testing.T) {
	xlsTimetable, err := Parse("parser/testdata/timetable.xls")
	tools.Fatal(t, assert.NoError(t, err))
	xlsxTimetable, err := Parse("parser/testdata/timetable.xlsx")
	tools.Fatal(t, assert.NoError(t, err))

	assert.Equal(t, xlsTimetable, xlsxTimetable)

	xlsJ, err := xlsTimetable.ToJson()
	tools.Fatal(t, assert.NoError(t, err))
	xlsxJ, err := xlsxTimetable.ToJson()
	tools.Fatal(t, assert.NoError(t, err))
	f := tools.LoadStringFromFile("parser/testdata/timetable.json")

	assert.Equal(t, xlsJ, f)
	assert.Equal(t, xlsxJ, f)
}
