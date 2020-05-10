package service

import (
	"github.com/pkg/errors"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"time"
)

var sleepDuration = time.Second * 10

func userModelsToIDs(users []*db.UserModel) []int {
	IDs := make([]int, 0, len(users))
	for _, user := range users {
		IDs = append(IDs, user.VkID)
	}
	return IDs
}

func SendNotifyAboutTimetableUpdate(b *vk.Bot) error {
	helpers.Logger.Info("Отправка уведомлений пользователям об обновлении расписания...")

	users, err := db.GetSubscribedUsers()
	if err != nil {
		return errors.WithStack(err)
	}

	for i := 0; i < len(users); i += 100 {
		end := i + 100
		if end > len(users) {
			end = len(users)
		}

		IDs := userModelsToIDs(users[i:end])
		if _, err := b.SendMultipleTextMessages("Обновилось расписание!", IDs...); err != nil {
			return errors.WithStack(err)
		}
		time.Sleep(sleepDuration)
	}

	return nil
}
