package parser

import (
	"github.com/pkg/errors"
	"github.com/sergeilem/xls"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
)

func xlsParser(filePath string) (*db.Timetable, error) {
	timetable := &db.Timetable{
		Groups: make(map[string]db.Group),
	}

	xlFile, err := xls.Open(filePath, "utf-8")
	if err != nil {
		return nil, errors.Wrap(err, "ошибка открытия xls файла")
	}

	for sheetNum := 0; sheetNum < xlFile.NumSheets(); sheetNum++ {
		if sheet := xlFile.GetSheet(sheetNum); sheet != nil {
			if err := processXlsSheet(sheet, timetable); err != nil {
				return nil, errors.WithStack(err)
			}
		} else {
			return nil, errors.Wrapf(err, "ошибка при получении листа в файле %s; sheetNum=%d", filePath, sheetNum)
		}
	}

	return timetable, nil
}

// Обработка листа таблицы. Лист в таблице - расписание для одного курса.
func processXlsSheet(sheet *xls.WorkSheet, t *db.Timetable) error {
	data := parseSheetData{
		t:           t,
		weekdayname: "",
		week:        -1,
		colGroupMap: make(map[int]string),
	}

	for rowNum := 0; rowNum <= (int(sheet.MaxRow)); rowNum++ {
		if row := sheet.Row(rowNum); row != nil {
			for colNum := 0; colNum <= row.LastCol(); colNum++ {
				if err := data.processCell(rowNum, colNum, row.Col(colNum)); err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}

	return nil
}
