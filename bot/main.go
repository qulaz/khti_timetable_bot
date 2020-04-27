package main

import (
	"github.com/getsentry/sentry-go"
	"gitlab.com/qulaz/khti_timetable_bot/bot/db"
	"gitlab.com/qulaz/khti_timetable_bot/bot/helpers"
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
		helpers.Logger.Errorf("PANIC!!: %+v", err)

		if err != nil {
			sentry.CurrentHub().Recover(err)
			closeApp()
		}
	}()
	defer closeApp()
}
