package parser

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"os"
	"path"
	"runtime"
	"testing"
)

func init() {
	// Тесты должны запускаться из дериктории в котрой находится этот файл
	_, filename, _, _ := runtime.Caller(0)

	dir := path.Join(path.Dir(filename))
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestParse__wrong_file_type(t *testing.T) {
	_, err := Parse("testdata/timetable.json")
	assert.Error(t, err)
}

func TestParse__nonexistent_file(t *testing.T) {
	_, err := Parse("testdata/table.xls")
	assert.Error(t, err)
}

func TestParse__xls(t *testing.T) {
	timetable, err := Parse("testdata/timetable.xls")
	tools.Fatal(t, assert.NoError(t, err))

	j, err := timetable.ToJson()
	tools.Fatal(t, assert.NoError(t, err))
	f := tools.LoadStringFromFile("testdata/timetable.json")
	assert.Equal(t, f, j)
}

func TestParse__xlsx(t *testing.T) {
	timetable, err := Parse("testdata/timetable.xlsx")
	tools.Fatal(t, assert.NoError(t, err))

	j, err := timetable.ToJson()
	tools.Fatal(t, assert.NoError(t, err))
	f := tools.LoadStringFromFile("testdata/timetable.json")
	assert.Equal(t, f, j)
}

func TestParse__xls_and_xlsx_parsers_the_same(t *testing.T) {
	xlsTimetable, err := Parse("testdata/timetable.xls")
	tools.Fatal(t, assert.NoError(t, err))
	xlsxTimetable, err := Parse("testdata/timetable.xlsx")
	tools.Fatal(t, assert.NoError(t, err))

	assert.Equal(t, xlsTimetable, xlsxTimetable)

	xlsJ, err := xlsTimetable.ToJson()
	tools.Fatal(t, assert.NoError(t, err))
	xlsxJ, err := xlsxTimetable.ToJson()
	tools.Fatal(t, assert.NoError(t, err))
	f := tools.LoadStringFromFile("testdata/timetable.json")

	assert.Equal(t, xlsJ, f)
	assert.Equal(t, xlsxJ, f)
}
