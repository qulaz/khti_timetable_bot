package parser

import (
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"path/filepath"
	"regexp"
	"strings"
)

// Стейт, хранимый при обработке листа таблицы
type parseSheetData struct {
	t *db.Timetable
	// название дня недели
	weekdayname string
	// индекс массива текущей недели (0 или 1 для первой и второй недели соответственно)
	week int
	// ключ - номер колонки, значение - код группы
	colGroupMap map[int]string
}

// функция обработки ячейки листа таблицы
func (data *parseSheetData) processCell(rowNum, colNum int, cell string) error {
	cell = tools.RemoveWhitespaces(cell) // удаление двух и более пробелов между словами в ячейке
	lowerCell := strings.ToLower(cell)

	// Определение недели расписания
	if colNum == 0 && cell != "" {
		switch lowerCell {
		case "первая неделя":
			data.week = 0
		case "вторая неделя":
			data.week = 1
		}
	}

	// Определение дня недели
	if colNum == 1 && cell != "" {
		switch lowerCell {
		case "понедельник", "вторник", "среда", "четверг", "пятница", "суббота":
			data.weekdayname = lowerCell
		}
	}

	// проверяем относится ли текущая колонка к какой-нибудь группе
	if groupCode, ok := data.colGroupMap[colNum]; ok {
		if group, ok := data.t.Groups[groupCode]; ok {
			if strings.Contains(lowerCell, "военной подготовке") {
				// Военка занимает весь день
				schedule := []string{cell, cell, cell, cell, cell}
				group.WeekSchedule[data.week][data.weekdayname] = schedule
			} else {
				// Если этого дня недели нет в мапе - создаем слайс для него
				if _, scheduleOk := group.WeekSchedule[data.week][data.weekdayname]; !scheduleOk {
					group.WeekSchedule[data.week][data.weekdayname] = make([]string, 0, 6)
				}

				// Отфильтровываем лишнее, так как в дне максимум может быть только 5 пар и добавляем пару
				if len(group.WeekSchedule[data.week][data.weekdayname]) < 5 {
					group.WeekSchedule[data.week][data.weekdayname] = append(
						group.WeekSchedule[data.week][data.weekdayname], cell,
					)
				}
			}
		}
	}

	// Из шапки расписания забираем информацию о группах
	if rowNum < 8 {
		// Regexp, который выбирает ячейку с номером и названием группы.
		// Прим: 58-1 09.03.03 Прикладная информатик
		if matched, err := regexp.Match(`\d+-\d \d\d.\d\d.\d\d .+`, []byte(cell)); matched {
			if err != nil {
				return errors.WithStack(err)
			}

			groupCode := cell[:4]                          // Оставляем первые 4 символа - код группы (прим: 58-1)
			data.t.Groups[groupCode] = newGroup(groupCode) // Добавляем группу в расписание
			// Отмечаем, что в текущей колонке находится расписание для найденной группы
			data.colGroupMap[colNum] = groupCode
		}
	}

	return nil
}

func newGroup(code string) db.Group {
	g := db.Group{
		Code: code,
		// Первый элемент массива - map с рассписанием первой недели, второй - map со второй
		WeekSchedule: [2]db.WeekSchedule{},
	}
	g.WeekSchedule[0] = make(db.WeekSchedule)
	g.WeekSchedule[1] = make(db.WeekSchedule)
	return g
}

func Parse(filePath string) (*db.Timetable, error) {
	ext := filepath.Ext(filePath)

	if ext == ".xls" {
		t, err := xlsParser(filePath)
		if err != nil {
			return nil, err
		}

		return t, nil
	} else if ext == ".xlsx" {
		t, err := xlsxParser(filePath)
		if err != nil {
			return nil, err
		}

		return t, nil
	}

	return nil, errors.Errorf("Файл %s имеет неподдерживаемый формат %s", filePath, ext)
}
