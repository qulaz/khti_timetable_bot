package db

import (
	"github.com/pkg/errors"
)

type GroupModel struct {
	Id   int
	Code string
}

// Получение модели группы из сканнера
func parseGroupModel(row scanner) (GroupModel, error) {
	var id int
	var code string

	if err := row.Scan(&id, &code); err != nil {
		return GroupModel{}, errors.WithStack(err)
	}

	return newGroup(id, code), nil
}

func newGroup(id int, code string) GroupModel {
	return GroupModel{Id: id, Code: code}
}

func createGroupIfNotExists(db queryable, groupCode string) error {
	_, err := db.Exec("INSERT INTO groups (code) values ($1) on conflict (code) do nothing;", groupCode)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Создание группы, если такой не существует
func CreateGroupIfNotExists(groupCode string) error {
	return createGroupIfNotExists(db, groupCode)
}

// Получение списка групп
func GetGroups(limit int, offset int) ([]GroupModel, error) {
	groups := make([]GroupModel, 0, limit)

	rows, err := db.Query("SELECT * FROM groups ORDER BY code LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return []GroupModel{}, errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		group, err := parseGroupModel(rows)

		if err != nil {
			return groups, errors.WithStack(err)
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// Получение общего колличества групп
func GroupsCount() (int, error) {
	var count int

	row := db.QueryRow("SELECT COUNT(*) as count FROM groups")
	if err := row.Scan(&count); err != nil {
		return 0, errors.WithStack(err)
	}

	return count, nil
}

// Получение группы по ее идентификатору
func GetGroup(groupId int) (GroupModel, error) {
	row := db.QueryRow("SELECT * FROM groups WHERE id=$1", groupId)
	group, err := parseGroupModel(row)

	if err != nil {
		return GroupModel{}, errors.WithStack(err)
	}

	return group, nil
}

// Получение группы по ее коду
func GetGroupByGroupCode(groupCode string) (GroupModel, error) {
	row := db.QueryRow("SELECT * FROM groups WHERE code=$1", groupCode)
	group, err := parseGroupModel(row)

	if err != nil {
		return GroupModel{}, errors.WithStack(err)
	}

	return group, nil
}
