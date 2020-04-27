# VK
[![coverage report](https://gitlab.com/qulaz/khti_timetable_bot/badges/master/coverage.svg?job=test_vk_library)](https://gitlab.com/qulaz/khti_timetable_bot/-/commits/master/vk)

Небольшая библиотека для создания чат-ботов во ВКонтакте.

## Пример использования
```go
package main

import (
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"log"
)

func StartHandler(b *vk.Bot, u *vk.MessageNew) {
	message := "Привет! Я бот, который поможет тебе следить за расписанием твоего ВУЗа"
	// Встроенные функции автоматически следят за соблюдением правил ВКонтакте относительно построения клавиатуры.
	// В случае, если Вы делаете что-то не так, они возвращают ошибку
	k := vk.NewKeyboard(false)
	k.AddTextButton("Сегодня", vk.COLOR_PRIMARY, &vk.ButtonPayload{Command: "/timetable", Body: "today"})
	k.AddTextButton("Завтра", vk.COLOR_PRIMARY, &vk.ButtonPayload{Command: "/timetable", Body: "tomorrow"})
	k.AddRow()
	k.AddTextButton("Первая неделя", vk.COLOR_PRIMARY, &vk.ButtonPayload{Command: "/timetable", Body: "first"})
	k.AddTextButton("Вторая неделя", vk.COLOR_PRIMARY, &vk.ButtonPayload{Command: "/timetable", Body: "second"})


	msgID, err := b.SendKeyboardMessage(message, k, u.Message.PeerID)
	if err != nil {
		b.Logger.Errorf("Ошибка при отправке сообщения: %s", err)
	}
	b.Logger.Debugf("Сообщение отправлено. Его ID: %d", msgID)
}

func TimetableHandler(b *vk.Bot, u *vk.MessageNew) {
	var response string
	// В u.Message.MessageBody хранится то, что мы передали в поле Body структуры ButtonPayload
	// при создании кнопки клавиатуры
	switch u.Message.MessageBody {
	case "today":
		response = "Расписание на сегодня"
	case "tomorrow":
		response = "Расписание на завтра"
	case "first":
		response = "Расписание на первую неделю"
	case "second":
		response = "Расписание на вторую неделю"
	}

	if _, err := b.SendTextMessage(response, u.Message.PeerID); err != nil {
		b.Logger.Errorf("Ошибка при отправке сообщения: %s", err)
	}
}

func main() {
	bot, err := vk.CreateBot(vk.Settings{
		GroupID: 123,
		Token:   "test_token",
	})
	if err != nil {
		log.Fatalf("Ошибка при создании бота: %s", err)
	}

	// Middleware позволяют логгировать апдейты, ограничивать доступ к командам и т.д.
	bot.AddMiddleware(vk.LoggingMiddleware)
	bot.AddMiddleware(vk.Middleware{
		// OnPreMessageNewUpdate запускается перед началом обработки нового сообщения. В случае, если из него вернется
		// false, обработка сообщения прекратится
		OnPreMessageNewUpdate: func(b *vk.Bot, u *vk.MessageNew) bool {
			adminID := 1
			if u.Message.MessageCommand == "/ban" && u.Message.PeerID == adminID {
				return true
			}
			return false
		},
	})

	bot.HandleCommand("/start", StartHandler)
	bot.HandleMessage("Начать", StartHandler)
	bot.HandleCommand("/timetable", TimetableHandler)

	if err := bot.Run(); err != nil {
		log.Fatalf("Ошибка при старте бота: %s", err)
	}
}
```
