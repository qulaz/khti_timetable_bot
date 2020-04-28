package service

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
)

const startMessage = "Привет! Я бот Хакасского Технического Института. Я буду сообщать тебе об изменениях " +
	"в расписании, важных новостях института и твоей группы.\n\nДля начала выбери группу, в которой ты учишься"

var RegisteredUserTryingToUseStartCommandError = errors.New("Registered user trying to use start command")

func StartCommand(d *Data) error {
	var err error

	isReset := d.u.Message.MessageBody == "reset"

	// Если пользователь уже существует и это не ресет - игнорируем сообщение
	if _, err := db.GetUserByVkID(d.u.Message.PeerID); err == nil && !isReset {
		return RegisteredUserTryingToUseStartCommandError
	}

	d.K, err = buildGroupKeyboard(groupsLimit, 0)
	if err != nil {
		return errors.Wrap(err, "ошибка сборки клавиатуры групп")
	}

	if isReset {
		d.Answer = "Выбери группу, в которой ты учишься"
	} else {
		d.Answer = startMessage
	}

	if !d.u.ClientInfo.IsKeyboardSupported() {
		d.Answer = fmt.Sprintf(
			"%s\n\nК сожалению твой клиент ВКонтакте не поддерживает кнопки для ботов, поэтому все взаимодействие "+
				"будет проходить с помощью команд.\nЧтобы выбрать свою группу, введи команду «‎%s номер_группы». "+
				"Например команда «‎%s 58-1» отнесет тебя к группе 58-1. Все команды вводятся без кавычек!",
			d.Answer,
			common.GroupCommand,
			common.GroupCommand,
		)
	}

	return nil
}
