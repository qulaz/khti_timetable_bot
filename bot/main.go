package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/robfig/cron/v3"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/handlers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/parser"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"os"
	"time"
)

func initApp() {
	helpers.LoadConfigFromEnv()
	helpers.InitLogger()
	db.InitDatabase()

	serverName, err := os.Hostname()
	if err != nil {
		serverName = "unknown"
	}

	var environment string
	if helpers.Config.IS_DEBUG {
		environment = "development"
	} else {
		environment = "production"
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              helpers.Config.SENTRY_DSN,
		AttachStacktrace: true,
		Environment:      environment,
		ServerName:       serverName,
	})
	if err != nil {
		helpers.Logger.Fatalf("sentry.Init: %s", err)
	}
}

func closeApp() {
	db.CloseDatabase()
	sentry.Flush(2 * time.Second)
	time.Sleep(time.Second * 5)
}

func main() {
	initApp()
	defer func() {
		err := recover()

		if err != nil {
			helpers.Logger.Errorf("PANIC!!: %+v", err)
			sentry.CurrentHub().Recover(err)
			closeApp()
		}
	}()
	defer closeApp()

	b, err := vk.CreateBot(vk.Settings{
		GroupID: helpers.Config.VK_GROUP_ID,
		Token:   helpers.Config.VK_GROUP_TOKEN,
		Logger:  helpers.Logger,
	})
	if err != nil {
		helpers.Logger.Fatalf("Ошибка создания бота: %+v", err)
	}

	// Регистрация обработчиков сообщений
	b.HandleMessage("Начать", handlers.Start)
	b.HandleMessage(`"/start"`, handlers.Start)
	b.HandleCommand(`/start`, handlers.Start)
	b.HandleCommand(common.StartCommand, handlers.Start)
	b.HandleCommand(common.MainCommand, handlers.Main)
	b.HandleCommand(common.GroupCommand, handlers.Group)
	b.HandleCommand(common.RingCommand, handlers.Ring)
	b.HandleCommand(common.TimetableCommand, handlers.Timetable)
	b.HandleCommand(common.WeekCommand, handlers.Week)
	b.HandleCommand(common.SettingsCommand, handlers.Settings)
	b.HandleMessageAllow(handlers.Allow)
	b.HandleMessageDeny(handlers.Deny)

	b.AddMiddleware(vk.Middleware{
		OnPostProcessUpdate: func(b *vk.Bot, u *vk.Update) bool {
			duration := time.Now().Sub(u.UpdateMeta.ReceiveTime)
			helpers.Logger.Infow(
				"Processed Update",
				"update_id", u.UpdateMeta.EventID,
				"process_time", duration.Milliseconds(),
				"update", u,
			)
			return true
		},
		OnPreMessageNewUpdate: func(b *vk.Bot, u *vk.MessageNew) bool {
			helpers.Logger.Infow(
				"Received message_new",
				"update_type", u.UpdateMeta.UpdateType,
				"update_id", u.UpdateMeta.EventID,
				"from_id", u.Message.FromID,
				"peer_id", u.Message.PeerID,
				"message_id", u.Message.ID,
				"update", u,
			)
			return true
		},
		OnPostMessageNewUpdate: func(b *vk.Bot, u *vk.MessageNew) bool {
			helpers.Logger.Infow(
				"Processed message",
				"update_type", u.UpdateMeta.UpdateType,
				"update_id", u.UpdateMeta.EventID,
				"from_id", u.Message.FromID,
				"peer_id", u.Message.PeerID,
				"message_id", u.Message.ID,
				"update", u,
				"handled", u.UpdateMeta.Handled,
			)
			return true
		},
		OnPreMessageAllowUpdate: func(b *vk.Bot, u *vk.MessageAllow) bool {
			helpers.Logger.Infow("Received message allow update",
				"update", u,
				"from_id", u.UserID,
			)
			return true
		},
		OnPreMessageDenyUpdate: func(b *vk.Bot, u *vk.MessageDeny) bool {
			helpers.Logger.Infow("Received message deny update",
				"update", u,
				"from_id", u.UserID,
			)
			return true
		},
	})

	// Добавление периодической задачи проверки расписания на предмет обновления (раз в 2 часа)
	c := cron.New()
	if _, err := c.AddFunc("0 1/2 * * *", func() {
		if err := parser.UpdateTimetable(b); err != nil {
			sentry.CaptureException(err)
			helpers.Logger.Errorw(
				"Ошибка в периодической задаче обновления расписания", "err", err,
			)
		}
	}); err != nil {
		helpers.Logger.Fatalf("Ошибка создания задачи обновления расписания: %+v", err)
	}

	// Добавление задачи, которая будет раз в 10 минут писать лог. Будет использоваться в качестве health check
	if _, err := c.AddFunc("*/10 * * * *", func() {
		helpers.Logger.Infow("Ping", "health", true)
	}); err != nil {
		helpers.Logger.Fatalf("Ошибка создания задачи health check: %+v", err)
	}

	c.Start()

	// Запуск этой задачи непосредственно перед запуском бота (для начальной инициализации расписания в базе данных)
	if err := parser.UpdateTimetable(b); err != nil {
		helpers.Logger.Fatal(err)
	}

	// Запуск бота
	if err := b.Run(); err != nil {
		helpers.Logger.Fatalf("Ошибка запуска бота: %+v", err)
	}
}
