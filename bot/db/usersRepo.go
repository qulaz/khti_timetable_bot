package db

import (
	"github.com/pkg/errors"
)

// Модель пользователя
type UserModel struct {
	ID           int
	VkID         int
	Group        GroupModel
	IsActive     bool
	IsSubscribed bool
}

// Установка статуса активности пользователя
func (u *UserModel) SetActive(b bool) error {
	_, err := db.Exec("UPDATE users SET is_active = $1 WHERE vk_id=$2", b, u.VkID)
	if err != nil {
		return errors.WithStack(err)
	}
	u.IsActive = b
	return nil
}

// Смена статуса подписки пользователя на обновления расписания
func (u *UserModel) SetSubscribe(b bool) error {
	_, err := db.Exec("UPDATE users SET is_subscribed = $1 WHERE vk_id=$2", b, u.VkID)
	if err != nil {
		return errors.WithStack(err)
	}
	u.IsSubscribed = b
	return nil
}

// Смена группы пользователя
func (u *UserModel) ChangeGroup(groupCode string) error {
	group, err := GetGroupByGroupCode(groupCode)
	if err != nil {
		return errors.Wrapf(err, "нет группы с кодом: %s", groupCode)
	}

	_, err = db.Exec("UPDATE users SET group_id = $1 WHERE vk_id=$2", group.Id, u.VkID)
	if err != nil {
		return errors.WithStack(err)
	}
	u.Group = newGroup(group.Id, group.Code)
	return nil
}

func parseUserModel(row scanner) (*UserModel, error) {
	var id int
	var vkId int
	var isActive bool
	var isSubscribed bool
	var groupId int
	var groupCode string

	if err := row.Scan(&id, &vkId, &isActive, &isSubscribed, &groupId, &groupCode); err != nil {
		return &UserModel{}, errors.WithStack(err)
	}

	return &UserModel{
		ID:           id,
		VkID:         vkId,
		Group:        newGroup(groupId, groupCode),
		IsActive:     isActive,
		IsSubscribed: isSubscribed,
	}, nil
}

// Создание пользователя
func CreateUser(vkID int, groupID int) error {
	_, err := db.Exec("INSERT INTO users(vk_id, group_id) VALUES ($1, $2);", vkID, groupID)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Получение пользователя по его VK ID
func GetUserByVkID(vkID int) (*UserModel, error) {
	row := db.QueryRow(
		"SELECT users.id, users.vk_id, users.is_active, users.is_subscribed, groups.id, groups.code "+
			"FROM users JOIN groups ON groups.id = users.group_id WHERE vk_id=$1",
		vkID,
	)
	user, err := parseUserModel(row)

	if err != nil {
		return &UserModel{}, errors.WithStack(err)
	}

	return user, nil
}

// Получение списка пользователей
func GetUsers(limit, offset int) ([]*UserModel, error) {
	users := make([]*UserModel, 0, limit)

	rows, err := db.Query(
		"SELECT users.id, users.vk_id, users.is_active, users.is_subscribed, groups.id, groups.code "+
			"FROM users JOIN groups ON groups.id = users.group_id LIMIT $1 OFFSET $2;",
		limit,
		offset,
	)
	if err != nil {
		return users, errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		user, err := parseUserModel(rows)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Получение списка пользователей, которые подписаны на уведомления об изменениях в расписании
func GetSubscribedUsers() ([]*UserModel, error) {
	users := make([]*UserModel, 0, 20)

	rows, err := db.Query(
		"SELECT users.id, users.vk_id, users.is_active, users.is_subscribed, groups.id, groups.code " +
			"FROM users JOIN groups ON groups.id = users.group_id " +
			"WHERE users.is_subscribed = true AND users.is_active = true;",
	)
	if err != nil {
		return users, errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		user, err := parseUserModel(rows)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}
