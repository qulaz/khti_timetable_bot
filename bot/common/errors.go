package common

import (
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"strconv"
)

// Текст, отправляемый пользователям в случае какой-то неизвестной ошибки
const UnknownErrorMessage = "Неизвестная ошибка. Попробуйте повторить позже."

var UnknownError = errors.New(UnknownErrorMessage)

// Шорткат для отправки сообщения об ошибке пользователю
func SendErrorMessageToUser(b *vk.Bot, message string, eventID string, peerID int) {
	if _, err := b.SendTextMessage(message, peerID); err != nil {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("event_id", eventID)
			scope.SetTag("error_sender", "ture")
			scope.SetUser(sentry.User{ID: strconv.Itoa(peerID)})
			scope.AddBreadcrumb(
				&sentry.Breadcrumb{
					Data: map[string]interface{}{
						"answer_message": message,
						"peer_id":        peerID,
					},
				},
				1,
			)

			sentry.CaptureException(err)
		})
	}
}
