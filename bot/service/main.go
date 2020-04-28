package service

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
)

var MainKeyboard = &vk.Keyboard{
	OneTime: false,
	Buttons: [][]vk.Button{
		{
			vk.TextButton{
				Action: vk.TextButtonAction{
					Type:    vk.TEXT_BUTTON,
					Label:   "Сегодня",
					Payload: &vk.ButtonPayload{Command: common.TimetableCommand, Body: "сегодня"},
				},
				Color: vk.COLOR_PRIMARY,
			},
			vk.TextButton{
				Action: vk.TextButtonAction{
					Type:    vk.TEXT_BUTTON,
					Label:   "Завтра",
					Payload: &vk.ButtonPayload{Command: common.TimetableCommand, Body: "завтра"},
				},
				Color: vk.COLOR_PRIMARY,
			},
		},
		{
			vk.TextButton{
				Action: vk.TextButtonAction{
					Type:    vk.TEXT_BUTTON,
					Label:   "Следующая пара",
					Payload: &vk.ButtonPayload{Command: common.TimetableCommand, Body: "следующая"},
				},
				Color: vk.COLOR_SECONDARY,
			},
		},
		{
			vk.TextButton{
				Action: vk.TextButtonAction{
					Type:    vk.TEXT_BUTTON,
					Label:   "Какая сейчас неделя?",
					Payload: &vk.ButtonPayload{Command: common.WeekCommand, Body: ""},
				},
				Color: vk.COLOR_SECONDARY,
			},
		},
		{
			vk.TextButton{
				Action: vk.TextButtonAction{
					Type:    vk.TEXT_BUTTON,
					Label:   "Сколько до звонка?",
					Payload: &vk.ButtonPayload{Command: common.RingCommand, Body: ""},
				},
				Color: vk.COLOR_SECONDARY,
			},
		},
		{
			vk.TextButton{
				Action: vk.TextButtonAction{
					Type:    vk.TEXT_BUTTON,
					Label:   "Первая неделя",
					Payload: &vk.ButtonPayload{Command: common.WeekCommand, Body: "1"},
				},
				Color: vk.COLOR_SECONDARY,
			},
			vk.TextButton{
				Action: vk.TextButtonAction{
					Type:    vk.TEXT_BUTTON,
					Label:   "Вторая неделя",
					Payload: &vk.ButtonPayload{Command: common.WeekCommand, Body: "2"},
				},
				Color: vk.COLOR_SECONDARY,
			},
		},
		{
			vk.TextButton{
				Action: vk.TextButtonAction{
					Type:    vk.TEXT_BUTTON,
					Label:   "Настройки",
					Payload: &vk.ButtonPayload{Command: common.SettingsCommand, Body: ""},
				},
				Color: vk.COLOR_SECONDARY,
			},
		},
	},
	Inline: false,
}

const (
	MainKeyboardHelp = "«Сегодня» - Расписание на сегодня;\n" +
		"«Завтра» - Расписание на завтра;\n" +
		"«Следующая пара» - Отправит следующую пару;\n" +
		"«Какая сейчас неделя?» - Номер текущей недели (первая, либо вторая);\n" +
		"«Сколько до звонка?» - Количество времени до начала/конца пары;\n" +
		"«Первая неделя» - Расписание на первую неделю;\n" +
		"«Вторая неделя» - Расписание на вторую неделю;\n" +
		"«Настройки» - Настройки бота"
	NonKeyboardMainHelp = "«/расписание номер_недели день_недели» - расписание на определенный день, " +
		"где номер_недели - это цифры 1 или 2 для первой и второй недели соответственно, " +
		"а день_недели - название дня недели (понедельник, вторник и тд.), " +
		"либо «полное» (без кавычек), чтобы показать расписание на всю неделю.\n" +
		"Пример:\n" +
		"> «/расписание 1 понедельник» - расписание на понедельник первой недели\n" +
		"> «/расписание 2 полное» - полное расписание второй недели\n\n" +
		"> «/неделя» - Номер текущей недели (первая, либо вторая);\n" +
		"> «/звонок» - Количество времени до начала/конца пары"
)
