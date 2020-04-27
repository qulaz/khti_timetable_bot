package db

import (
	"github.com/pkg/errors"
)

type TimetableModel struct {
	ID           int
	RawTimetable string
	*Timetable
}

func parseTimetableModel(row scanner) (TimetableModel, error) {
	var id int
	var rawTimetable string

	if err := row.Scan(&id, &rawTimetable); err != nil {
		return TimetableModel{}, errors.WithStack(err)
	}

	timetable := TimetableModel{ID: id, RawTimetable: rawTimetable, Timetable: &Timetable{}}
	if err := timetable.Timetable.FromJson(rawTimetable); err != nil {
		return TimetableModel{}, err
	}

	return timetable, nil
}

// Получение объекта расписания из базы данных
func GetTimetable() (*Timetable, error) {
	row := db.QueryRow("SELECT id, timetable FROM timetable ORDER BY id DESC LIMIT 1;")
	timetableModel, err := parseTimetableModel(row)
	if err != nil {
		return &Timetable{}, err
	}

	return timetableModel.Timetable, nil
}

// Имеется ли запись с расписанием в базе или нет
func isTimetableExists() (bool, error) {
	var count int

	row := db.QueryRow("SELECT COUNT(*) as count FROM timetable")
	if err := row.Scan(&count); err != nil {
		return false, errors.WithStack(err)
	}

	return count > 0, nil
}

func createTimetable(db queryable, rawTimetable string) error {
	_, err := db.Exec("INSERT INTO timetable (timetable) values ($1);", rawTimetable)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func updateTimetable(db queryable, rawTimetable string) error {
	exists, err := isTimetableExists()
	if err != nil {
		return errors.Wrap(err, "ошибка определения кол-ва записей в таблице timetable")
	}

	if exists {
		_, err = db.Exec("UPDATE timetable SET timetable = $1 WHERE id=(SELECT MAX(id) FROM timetable)", rawTimetable)
	} else {
		err = createTimetable(db, rawTimetable)
	}
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Обновление расписания, записанного в базе данных
func UpdateTimetable(rawTimetable string) error {
	return updateTimetable(db, rawTimetable)
}
