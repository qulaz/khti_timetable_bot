package service

import (
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
)

// Максимальное количество кнопок с группами в клавиатуре
const groupsLimit = 27 // 3 кнопки в ряду, 9 рядов, 10-й ряд - кнопки управления паггинацией

// Сборка клавиатуры с группами
func buildGroupKeyboard(limit, offset int) (*vk.Keyboard, error) {
	if limit > groupsLimit {
		return nil, errors.Errorf(
			"Переданное значение limit %d больше максимально допустимого %d", limit, groupsLimit,
		)
	}

	totalGroups, err := db.GroupsCount()
	if err != nil {
		return nil, errors.Wrap(err, "ошибка получения количества групп")
	}

	groups, err := db.GetGroups(limit, offset)
	if err != nil {
		return nil, errors.Wrapf(err,
			"ошибка получения списка групп. limit: %d; offset %d; total %d", limit, offset, totalGroups,
		)
	}

	prevOffset, nextOffset := getPrevAndNextOffset(totalGroups, limit, offset)

	k := vk.NewKeyboard(false)
	for i, group := range groups {
		if i%3 == 0 {
			if err := k.AddRow(); err != nil {
				helpers.Logger.Errorw(
					"ошибка добавления ряда в клавиатуру группы",
					"i", i, "group", group.Code, "err", err, "keyboard", k,
				)
				return nil, errors.Errorf("ошибка добавления ряда на итерации %d, группа %q", i, group.Code)
			}
		}
		if err := k.AddTextButton(
			group.Code,
			vk.COLOR_SECONDARY,
			&vk.ButtonPayload{Command: common.GroupCommand, Body: group.Code, Offset: 0},
		); err != nil {
			helpers.Logger.Errorw(
				"ошибка добавления кнопки в клавиатуру группы",
				"i", i, "group", group.Code, "err", err, "keyboard", k,
			)
			return nil, errors.Errorf("ошибка добавления кнопки на итерации %d, группа %q", i, group.Code)
		}
	}

	// В случае, если есть оффсеты - добавляем ряд пагинации
	if nextOffset > -1 || prevOffset > -1 {
		if err := k.AddRow(); err != nil {
			helpers.Logger.Errorw(
				"ошибка добавление ряда для пагинации в клавиатуру группы",
				"err", err, "keyboard", k,
			)
			return nil, errors.New("ошибка добавление ряда для пагинации")
		}
		if prevOffset > -1 {
			if err := k.AddTextButton(
				"< Назад",
				vk.COLOR_PRIMARY,
				&vk.ButtonPayload{Command: common.GroupCommand, Body: prevBody, Offset: prevOffset},
			); err != nil {
				helpers.Logger.Errorw(
					"ошибка добавления кнопки назад в клавиатуру группы",
					"err", err, "keyboard", k,
				)
				return nil, errors.New("ошибка добавления кнопки назад")
			}
		}
		if nextOffset > -1 {
			if err := k.AddTextButton(
				"Далее >",
				vk.COLOR_PRIMARY,
				&vk.ButtonPayload{Command: common.GroupCommand, Body: nextBody, Offset: nextOffset},
			); err != nil {
				helpers.Logger.Errorw(
					"ошибка добавления кнопки далее в клавиатуру группы",
					"err", err, "keyboard", k,
				)
				return nil, errors.New("ошибка добавления кнопки далее")
			}
		}
	}

	return k, nil
}
