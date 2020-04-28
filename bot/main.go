package main

import (
	"github.com/getsentry/sentry-go"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/handlers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
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
	})
	if err != nil {
		helpers.Logger.Fatalf("Ошибка создания бота: %+v", err)
	}

	b.HandleMessage("Начать", handlers.Start)
	b.HandleCommand(common.StartCommand, handlers.Start)
	b.HandleCommand(common.GroupCommand, handlers.Group)

	if err := b.Run(); err != nil {
		helpers.Logger.Fatalf("Ошибка запуска бота: %+v", err)
	}
}
