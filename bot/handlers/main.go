package handlers

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/service"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
)

// Отправляет основную клавиатуру
func Main(b *vk.Bot, u *vk.MessageNew) {
	answer := "Основное меню" // возможно стоит придумать сообщение по-лучше
	k := service.MainKeyboard

	if _, err := b.SendKeyboardMessage(answer, k, u.Message.PeerID); err != nil {
		common.SendHandlerErrToSentry(
			common.MainCommand, err, common.DefaultHandlerBreadcrumbs(u, answer, k)...,
		)
	}
}
