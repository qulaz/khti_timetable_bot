package service

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
)

func buildSettingsKeyboard(u *db.UserModel) (string, *vk.Keyboard) {
	// Формат: номер группы, subscribedButtonLabel, subscribedState
	const SettingKeyboardHelp = "«Изменить группу» - смена установленной группы. Текущая группа: %s;\n" +
		"«%s» - %s оповещения об изменениях в расписании;"
	var (
		subscribedButtonLabel string
		subscribedButtonColor string
		subscribedState       string
	)

	if u.IsSubscribed {
		subscribedButtonLabel = "Выкл. оповещения о расписании"
		subscribedButtonColor = vk.COLOR_NEGATIVE
		subscribedState = "выключить"
	} else {
		subscribedButtonLabel = "Вкл. оповещения о расписании"
		subscribedButtonColor = vk.COLOR_POSITIVE
		subscribedState = "включить"
	}

	k := &vk.Keyboard{
		OneTime: false,
		Buttons: [][]vk.Button{
			{
				vk.TextButton{
					Action: vk.TextButtonAction{
						Type:    vk.TEXT_BUTTON,
						Label:   "Изменить группу",
						Payload: &vk.ButtonPayload{Command: common.StartCommand, Body: "reset"},
					},
					Color: vk.COLOR_PRIMARY,
				},
			},
			{
				vk.TextButton{
					Action: vk.TextButtonAction{
						Type:    vk.TEXT_BUTTON,
						Label:   subscribedButtonLabel,
						Payload: &vk.ButtonPayload{Command: common.SettingsCommand, Body: "расписание"},
					},
					Color: subscribedButtonColor,
				},
			},
			{
				vk.TextButton{
					Action: vk.TextButtonAction{
						Type:    vk.TEXT_BUTTON,
						Label:   "Назад",
						Payload: &vk.ButtonPayload{Command: common.MainCommand},
					},
					Color: vk.COLOR_SECONDARY,
				},
			},
		},
		Inline: false,
	}

	return fmt.Sprintf(SettingKeyboardHelp, u.Group.Code, subscribedButtonLabel, subscribedState), k
}

func SettingsCommand(d *Data) error {
	user, err := db.GetUserByVkID(d.u.Message.PeerID)
	if err != nil {
		return errors.WithStack(err)
	}

	switch body := d.u.Message.MessageBody; body {
	case "расписание":
		if err := user.SetSubscribe(!user.IsSubscribed); err != nil {
			return errors.WithStack(err)
		}

		d.Answer = "Готово!"
		_, d.K = buildSettingsKeyboard(user)
		return nil
	case "":
		d.Answer, d.K = buildSettingsKeyboard(user)
		d.Answer = "Подробнее о командах настроек:\n" + d.Answer
		return nil
	default:
		return common.IgnoreMessageError
	}
}
