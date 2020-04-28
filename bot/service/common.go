package service

import (
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/tools"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
)

const (
	prevBody = "prev"
	nextBody = "next"
)

type Data struct {
	u       *vk.MessageNew
	command string

	Answer string
	K      *vk.Keyboard
}

// Проверка данных и обработка ошибок. Если возвращается true - значит все хорошо, данные можно отправлять клиенту,
// если false, то с данными что-то не так
func (d Data) Validate(b *vk.Bot, err error) bool {
	if err != nil {
		common.SendErrorMessageToUser(
			b, tools.SelectNonEmptyString(d.Answer, common.UnknownErrorMessage),
			d.u.UpdateMeta.EventID, d.u.Message.PeerID,
		)
		common.SendHandlerErrToSentry(d.command, err, common.DefaultHandlerBreadcrumbs(d.u, d.Answer, d.K)...)
		return false
	}

	// Пустая клавиатура не отображается в вк при отправке
	if d.K == nil {
		d.K = vk.NewKeyboard(true)
	}

	if d.Answer == "" || d.Answer == common.UnknownErrorMessage {
		common.SendErrorMessageToUser(b, common.UnknownErrorMessage, d.u.UpdateMeta.EventID, d.u.Message.PeerID)
		common.SendHandlerErrToSentry(
			d.command, errors.New("Попытка отправки пустого сообщения"),
			common.DefaultHandlerBreadcrumbs(d.u, d.Answer, d.K)...,
		)

		return false
	}

	return true
}

func NewData(u *vk.MessageNew, command, answer string, k *vk.Keyboard) *Data {
	return &Data{
		u:       u,
		command: command,
		Answer:  answer,
		K:       k,
	}
}

// Определяет предыдущий и следующий оффсеты. В случае, если таковых нет, возвращается -1
func getPrevAndNextOffset(total, limit, offset int) (prevOffset, nextOffset int) {
	if nextOffset = offset + limit; nextOffset >= total {
		nextOffset = -1
	}

	if offset < 0 {
		prevOffset = -1
	} else {
		if prevOffset = offset - limit; prevOffset < 0 {
			if offset == 0 {
				prevOffset = -1
			} else {
				prevOffset = 0
			}
		}
	}

	return
}
