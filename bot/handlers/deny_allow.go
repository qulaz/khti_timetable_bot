package handlers

import (
	"github.com/getsentry/sentry-go"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"strconv"
)

func captureErr(err error, uid int) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{ID: strconv.Itoa(uid)})
		sentry.CaptureException(err)
	})
}

func Deny(b *vk.Bot, u *vk.MessageDeny) {
	helpers.Logger.Infow("message_deny", "user_id", u.UserID)

	user, err := db.GetUserByVkID(u.UserID)
	if err != nil {
		captureErr(err, u.UserID)
		return
	}

	if err := user.SetActive(false); err != nil {
		captureErr(err, u.UserID)
	}
}

func Allow(b *vk.Bot, u *vk.MessageAllow) {
	helpers.Logger.Infow("message_allow", "user_id", u.UserID)

	user, err := db.GetUserByVkID(u.UserID)
	if err != nil {
		captureErr(err, u.UserID)
		return
	}

	if err := user.SetActive(true); err != nil {
		captureErr(err, u.UserID)
	}
}
