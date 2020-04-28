package handlers

import (
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/service"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
)

func Group(b *vk.Bot, u *vk.MessageNew) {
	command := common.GroupCommand

	data := service.NewData(u, command, common.UnknownErrorMessage, vk.NewKeyboard(true))
	err := service.GroupCommand(data)
	if !data.Validate(b, err) {
		common.SendHandlerErrToSentry(
			command, err, common.DefaultHandlerBreadcrumbs(u, data.Answer, data.K)...,
		)
		return
	}

	if _, err := b.SendKeyboardMessage(data.Answer, data.K, u.Message.PeerID); err != nil {
		common.SendHandlerErrToSentry(
			command, err, common.DefaultHandlerBreadcrumbs(u, data.Answer, data.K)...,
		)
	}
}
