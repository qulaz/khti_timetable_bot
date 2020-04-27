package parser

import (
	"github.com/pkg/errors"
	"github.com/tealeg/xlsx"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
)

func xlsxParser(filePath string) (*db.Timetable, error) {
	timetable := &db.Timetable{
		Groups: make(map[string]db.Group),
	}

	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Ошбика открытия xlsx файла %s", filePath)
	}

	for _, sheet := range xlFile.Sheets {
		if err := processXlsxSheet(sheet, timetable); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return timetable, nil
}

func processXlsxSheet(sheet *xlsx.Sheet, t *db.Timetable) error {
	data := parseSheetData{
		t:           t,
		weekdayname: "",
		week:        -1,
		colGroupMap: make(map[int]string),
	}

	for rowNum, row := range sheet.Rows {
		for colNum, cell := range row.Cells {
			if err := data.processCell(rowNum, colNum, cell.Value); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}
