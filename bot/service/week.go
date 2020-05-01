package service

import (
	"fmt"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"strconv"
	"strings"
)

func buildWeekKeyboard(weekNum int) *vk.Keyboard {
	k := vk.NewKeyboard(false)

	for i, dayName := range tools.Weekdays {
		if i%2 == 0 {
			k.AddRow()
		}
		k.AddTextButton(
			strings.Title(strings.ToLower(dayName)),
			vk.COLOR_PRIMARY,
			&vk.ButtonPayload{Command: common.TimetableCommand, Body: fmt.Sprintf("%d %s", weekNum, dayName)},
		)
	}

	k.AddRow()
	k.AddTextButton(
		"Полное расписание недели",
		vk.COLOR_PRIMARY,
		&vk.ButtonPayload{Command: common.TimetableCommand, Body: fmt.Sprintf("%d %s", weekNum, db.FullWeek)},
	)
	k.AddRow()
	k.AddTextButton(
		"Назад",
		vk.COLOR_SECONDARY,
		&vk.ButtonPayload{Command: common.MainCommand},
	)

	return k
}

func WeekCommand(d *Data) error {
	switch body := d.u.Message.MessageBody; body {
	// Какая сейчас неделя?
	case "":
		weekNum := tools.GetCurrentWeekNum()
		if weekNum == 1 {
			d.Answer = "Текущая неделя - первая"
			return nil
		} else if weekNum == 2 {
			d.Answer = "Текущая неделя - вторая"
			return nil
		}
	// Первая/вторая неделя. Нужны для отправки клавиатуры с днями недели
	case "1", "2":
		weekNum, _ := strconv.Atoi(body)

		d.Answer = "Выбери день недели"
		d.K = buildWeekKeyboard(weekNum)
		return nil
	}

	return common.IgnoreMessageError
}
