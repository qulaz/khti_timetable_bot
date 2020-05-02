package main

import (
	"github.com/getsentry/sentry-go"
	"gitlab.com/qulaz/khti_timetable_bot/bot/common"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/handlers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
	"gitlab.com/qulaz/khti_timetable_bot/bot/parser"
	"gitlab.com/qulaz/khti_timetable_bot/vk"
	"log"
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

// Временная функция инициализации расписания в базе данных, чтобы на новой машине автоматически были все записи
// В дальнейшем заменится на реальный поиск расписания на сайте и запись его в бд
func loadTimetable() {
	t, err := parser.Parse("timetable.xls")
	if err != nil {
		log.Fatal(err)
	}
	if err := t.WriteInDB(); err != nil {
		log.Fatal(err)
	}
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
	loadTimetable()

	b, err := vk.CreateBot(vk.Settings{
		GroupID: helpers.Config.VK_GROUP_ID,
		Token:   helpers.Config.VK_GROUP_TOKEN,
	})
	if err != nil {
		helpers.Logger.Fatalf("Ошибка создания бота: %+v", err)
	}

	b.HandleMessage("Начать", handlers.Start)
	b.HandleCommand(common.StartCommand, handlers.Start)
	b.HandleCommand(common.MainCommand, handlers.Main)
	b.HandleCommand(common.GroupCommand, handlers.Group)
	b.HandleCommand(common.RingCommand, handlers.Ring)
	b.HandleCommand(common.TimetableCommand, handlers.Timetable)
	b.HandleCommand(common.WeekCommand, handlers.Week)
	b.HandleCommand(common.SettingsCommand, handlers.Settings)
	b.HandleMessageAllow(handlers.Allow)
	b.HandleMessageDeny(handlers.Deny)

	if err := b.Run(); err != nil {
		helpers.Logger.Fatalf("Ошибка запуска бота: %+v", err)
	}
}
