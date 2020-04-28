package common

import (
	"github.com/getsentry/sentry-go"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"strconv"
)

// Структура для "хлебных крошек" Sentry
type Breadcrumb struct {
	Key   string
	Value interface{}
}

// Шаблонная отправка ошибок на Sentry
func SendHandlerErrToSentry(command string, err error, breadcrumbs ...Breadcrumb) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetTag("command", command)

		breadLimit := len(breadcrumbs)
		for _, breadcrumb := range breadcrumbs {
			scope.AddBreadcrumb(
				&sentry.Breadcrumb{
					Data: map[string]interface{}{breadcrumb.Key: breadcrumb.Value},
				},
				breadLimit,
			)

			if u, ok := breadcrumb.Value.(*vk.MessageNew); ok {
				scope.SetUser(sentry.User{ID: strconv.Itoa(u.Message.PeerID)})
				scope.SetTag("event_id", u.UpdateMeta.EventID)
			}
		}

		sentry.CaptureException(err)
	})
}

// Стандартные "хлебные крошки" для хэндлеров новых сообщений
func DefaultHandlerBreadcrumbs(u *vk.MessageNew, answer string, k *vk.Keyboard) []Breadcrumb {
	breadcrumbs := make([]Breadcrumb, 0, 3)

	breadcrumbs = append(breadcrumbs, Breadcrumb{"answer_message", answer})
	breadcrumbs = append(breadcrumbs, Breadcrumb{"keyboard", k})
	breadcrumbs = append(breadcrumbs, Breadcrumb{"message", u})

	return breadcrumbs
}
