package db

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"strings"
)

const FullWeek = "полное"

// Структура расписания
type Timetable struct {
	// ключ - код группы ("58-1", "19-1" и тд), значение - инстанс типа Group
	Groups map[string]Group `json:"groups"`
}

// Запись текущего расписания в базу данных. Записывается как непосредственно само расписание, так и актуализируется
// список студенческих групп
func (t *Timetable) WriteInDB() error {
	// Обновление происходит в транзакции. В случае любой ошибки транзакция отменяется
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Добавление новых групп в базу данных
	for _, group := range t.Groups {
		if err := createGroupIfNotExists(tx, group.Code); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Сериализуем расписание в JSON и обновляем его в базе даннх
	rawTimetable, err := t.ToJson()
	if err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}
	if err := updateTimetable(tx, rawTimetable); err != nil {
		tx.Rollback()
		return errors.WithStack(err)
	}

	return errors.WithStack(tx.Commit())
}

// Получение расписания на неделю по коду группы и номеру недели. Номер недели (weekNum) принимает только числа 1 и 2,
// что соответствует первой и второй недели.
func (t *Timetable) GetWeekSchedule(groupCode string, weekNum int) (WeekSchedule, error) {
	if weekNum != 1 && weekNum != 2 {
		return nil, errors.New("параметр weekNum может принимать только числа 1 или 2")
	}

	if c, ok := t.Groups[groupCode]; ok {
		return c.WeekSchedule[weekNum-1], nil
	}

	return nil, errors.New("не найдена группа с переданным кодом")
}

// Возвращает оформленное строковое представление расписание группы на указанный день недели, либо на всю неделю
//
//  groupCode: номер группы ("58-1", "19-1" и тд.)
//
//  dayName: название дня недели на русском языке в нижнем регистре (из слайса tools.Weekdays), либо значение константы
//  FullWeek, чтобы вернуть расписание на всю неделю
//
//  weekNum: номер недели. Принимает значения 1 или 2 которые соответствую первой и второй недели
func (t *Timetable) GetStringifyedSchedule(groupCode string, weekNum int, dayName string) (string, error) {
	schedule, err := t.GetWeekSchedule(groupCode, weekNum)
	if err != nil {
		return "", errors.WithStack(err)
	}

	s, err := schedule.ToString(dayName)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return s, nil
}

// Сериализация расписания в json формат
func (t *Timetable) ToJson() (string, error) {
	j, err := json.Marshal(t)

	if err != nil {
		return "", errors.Wrap(err, "ошибка преобразования расписания в json-строку")
	}

	return string(j), err
}

// Десериализация из json-строки
func (t *Timetable) FromJson(j string) error {
	if err := json.Unmarshal([]byte(j), t); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Расписание на день. Тип представляет из себя слайс строк с названиями предметов: [ 2 пара 3 пара  ]
type DaySchedule []string

// Расписание на неделю. Тип представляет из себя map ключ в котором - название дня недели, а значение инстанс
// DaySchedule: map[понедельник:[ 2 пара 3 пара  ] вторник:[1 пара  3 пара  ] ...]
type WeekSchedule map[string]DaySchedule

// Возвращает DaySchedule по названию дня недели.
//
//  dayName: название дня недели на русском языке в нижнем регистре (из слайса tools.Weekdays)
func (s WeekSchedule) GetDaySchedule(dayName string) (DaySchedule, error) {
	if dayName == tools.Sunday {
		return DaySchedule{}, errors.New("Воскресенье - выходной!")
	}

	if d, ok := s[dayName]; ok {
		return d, nil
	}

	return DaySchedule{}, errors.Errorf("Неверный день недели: %q", dayName)
}

// Преобразование расписания на день в форматированную строку
func (d DaySchedule) ToString() string {
	tmpDay := make([]string, len(d), len(d))
	for i, subject := range d {
		if subject == "" {
			d[i] = "-"
		}
		tmpDay[i] = fmt.Sprintf("%d пара:\n%s", i+1, d[i])
	}
	return strings.Join(tmpDay, "\n--------------------\n")
}

// Преобразование расписания на день/неделю в форматированную строку
//
//  dayName: название дня недели на русском языке в нижнем регистре (из слайса tools.Weekdays), либо значение константы
//  FullWeek, чтобы вернуть расписание на всю неделю
func (s WeekSchedule) ToString(dayName string) (string, error) {
	if d, err := s.GetDaySchedule(dayName); err == nil {
		return fmt.Sprintf("> %s:\n%s", tools.FormattedWeekdays[dayName], d.ToString()), nil
	}

	if dayName != FullWeek {
		return "", errors.Errorf("Неверный день недели: %q", dayName)
	}

	tmp := make([]string, 0, 6)
	for _, dayName := range tools.Weekdays {
		if d, ok := s[dayName]; ok {
			tmp = append(tmp, fmt.Sprintf("> %s:\n%s", tools.FormattedWeekdays[dayName], d.ToString()))
		} else {
			return "", errors.Errorf("не найден день недели в расписании: %s", dayName)
		}
	}

	return strings.Join(tmp, "\n\n"), nil
}

type Group struct {
	// Код группы ("58-1", "19-1" и тд.)
	Code string `json:"code"`
	// Расписание на 2 недели для группы
	WeekSchedule [2]WeekSchedule `json:"schedule"`
}
